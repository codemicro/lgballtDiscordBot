package adminSite

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/lgballtDiscordBot/internal/adminSite/templates"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

var oauthConf = &oauth2.Config{
	ClientID:     config.AdminSite.ClientID,
	ClientSecret: config.AdminSite.ClientSecret,
	Scopes:       []string{"guilds", "identify"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://discord.com/api/oauth2/authorize",
		TokenURL: "https://discord.com/api/oauth2/token",
	},
}

var redirectURI = fmt.Sprintf("%s/auth/inbound", config.AdminSite.VisibleURL)

const stateKey = "oauth2.state"
const requestCodeKey = "oauth2.code"
const oauthTokenKey = "oauth2.token"

const guildIDsKey = "guild.ids"
const guildRolesKey = "guild.role.ids"

const userIDKey = "user.id"
const userNameKey = "user.name"

const hasAuthKey = "auth.has"

// for use in oauth2 interactions with Discord
type oauthState struct {
	NextURL string `json:"n"`
	CheckValue string `json:"c"`
}

func (o *oauthState) String() string {
	dat, _ := json.Marshal(o)
	return string(dat)
}

func oauthStateFromString(x string) (*oauthState, error) {
	o := new(oauthState)
	err := json.Unmarshal([]byte(x), o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (w *webApp) authInbound(ctx *fiber.Ctx) error {

	sess := getSession(ctx)

	code := ctx.Query("code")
	receivedState := ctx.Query("state")

	storedState := sess.Get(stateKey)
	if storedState == nil || code == "" || receivedState == "" {
		return fiber.ErrBadRequest
	}

	storedStateString := storedState.(string)

	if receivedState != storedStateString {
		return fiber.ErrBadRequest
	}

	decodedState, err := oauthStateFromString(storedStateString)
	if err != nil {
		return err
	}

	// all should be ok
	sess.Set(requestCodeKey, code)

	// attempt oauth exchange
	backgroundContext := context.Background()
	token, err := oauthConf.Exchange(backgroundContext, code,
		oauth2.SetAuthURLParam("redirect_uri", redirectURI),
		oauth2.SetAuthURLParam("grant_type", "authorization_code"),
		oauth2.SetAuthURLParam("code", code),
		oauth2.SetAuthURLParam("client_id", oauthConf.ClientID),
		oauth2.SetAuthURLParam("client_secret", oauthConf.ClientSecret),
	)
	if err != nil {
		return err
	}

	sess.Set(oauthTokenKey, token)

	// get guild listing
	var guildIDs []string
	var guildRoles []string
	{
		dg, err := discordgo.New("Bearer " + token.AccessToken)
		if err != nil {
			return err
		}

		dgBot, err := discordgo.New("Bot " + config.Token)
		if err != nil {
			return err
		}

		me, err := dg.User("@me")
		if err != nil {
			return err
		}

		sess.Set(userIDKey, me.ID)
		sess.Set(userNameKey, me.Username)

		guilds, err := dg.UserGuilds(100, "0", "")
		if err != nil {
			return err
		}

		for _, guild := range guilds {
			guildIDs = append(guildIDs, guild.ID)

			if guild.ID == config.MainGuildID {

				guildMember, err := dgBot.GuildMember(config.MainGuildID, me.ID)
				if err != nil {
					return err
				}

				guildRoles = guildMember.Roles

			}

		}

	}

	sess.Set(guildRolesKey, guildRoles)
	sess.Set(guildIDsKey, guildIDs)

	sess.Set(hasAuthKey, true)

	if err := sess.Save(); err != nil {
		return err
	}

	if decodedState.NextURL != "" {
		return ctx.Redirect(decodedState.NextURL)
	}

	return ctx.Redirect("/services")
}

func (w* webApp) authLogout(ctx *fiber.Ctx) error {
	sess := getSession(ctx)
	err := sess.Destroy()
	if err != nil {
		return err
	}
	return ctx.Type("html").SendString(templates.RenderPage(&templates.FeedbackPage{
		WasSuccess:        true,
		Message:           "Signed out!",
		NextURL:           "/",
		RedirectTimeoutMs: 1000,
	}))
}