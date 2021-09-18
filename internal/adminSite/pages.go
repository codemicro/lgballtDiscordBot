package adminSite

import (
	"crypto/sha256"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

func (w *webApp) index(ctx *fiber.Ctx) error {

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
		state = fmt.Sprintf("%x", hashedBytes)
		sess.Set(stateKey, state)

		if err := sess.Save(); err != nil {
			return err
		}

	}

	oauthURL := oauthConf.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("redirect_uri", redirectURI),
		oauth2.SetAuthURLParam("prompt", "consent"),
	)

	ctx.Set(fiber.HeaderCacheControl, "no-store")

	return ctx.Type("html").SendString(fmt.Sprintf("<a href='%s'>%s</a>", oauthURL, oauthURL))
}

func (w *webApp) serviceListing(ctx *fiber.Ctx) error {
	
	sess := getSession(ctx)

	// TODO: don't show bio tool link if not admin

	username := sess.Get(userNameKey).(string)

	return ctx.Type("html").SendString(fmt.Sprintf("<h1>Hello <b>%s</b></h1><br><a href='bio'>Bio tool</a>", username))
}