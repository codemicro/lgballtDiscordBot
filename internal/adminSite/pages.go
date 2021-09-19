package adminSite

import (
	"crypto/sha256"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/adminSite/templates"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

func (w *webApp) index(ctx *fiber.Ctx) error {

	nextURL := ctx.Query("next")

	sess := getSession(ctx)
	auth := getAuth(ctx)

	if auth.HasAuth {
		return ctx.Redirect("/services")
	}

	// determine a value for `state`
	var state string
	{
		sessionID := sess.ID()
		hashedBytes := sha256.Sum256([]byte(sessionID))

		s := oauthState{CheckValue: fmt.Sprintf("%x", hashedBytes), NextURL: nextURL}

		state = s.String()
		
		sess.Set(stateKey, state)

		if err := sess.Save(); err != nil {
			return err
		}

	}

	oauthURL := oauthConf.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("redirect_uri", redirectURI),
		oauth2.SetAuthURLParam("prompt", "none"),
	)

	ctx.Set(fiber.HeaderCacheControl, "no-store")

	if nextURL != "" {
		return ctx.Redirect(oauthURL)
	}

	return ctx.Type("html").SendString(templates.RenderPage(&templates.IndexPage{DiscordLoginURL: oauthURL}))
}

func (w *webApp) serviceListing(ctx *fiber.Ctx) error {
	
	sess := getSession(ctx)
	auth := getAuth(ctx)

	var links []templates.ActionButton
	if auth.IsAdmin {
		links = append(links, templates.ActionButton{Title: "Bio manager", Location: "/bio"})
	}

	username := sess.Get(userNameKey).(string)

	return ctx.Type("html").SendString(templates.RenderPage(&templates.ServicesPage{Name: username, Actions: links}))
}