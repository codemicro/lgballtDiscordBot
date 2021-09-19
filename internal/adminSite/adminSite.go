package adminSite

import (
	"github.com/codemicro/lgballtDiscordBot/internal/bot/common"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/state"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

func Start(state *state.State) error {

	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			log.Error().Err(err).Str("area", "fiber-handler").Send()

			var message string
			if code == 500 {
				message = "Internal server error - check back later"
			} else {
				message = err.Error()
			}

			return ctx.Status(code).JSON(map[string]string{"status": "error", "message": message})
		},
		DisableStartupMessage: !config.DebugMode,
	})

	setupWebApp(app)

	go func() {
		if err := app.Listen(config.AdminSite.ServeAddress); err != nil {
			log.Error().Err(err).Str("area", "fiber-listen").Send()
		}
	}()

	go func() {
		state.WaitUntilShutdownTrigger()
		_ = app.Shutdown()
		state.FinishGoroutine()
	}()

	log.Info().Msgf("Running admin site HTTP server at %s", config.AdminSite.ServeAddress)

	return nil
}

type auth struct {
	HasAuth bool
	IsAdmin bool
}

type webApp struct {
	session *session.Store
}

func getAuth(ctx *fiber.Ctx) *auth {
	return ctx.Locals("auth").(*auth)
}

func getSession(ctx *fiber.Ctx) *session.Session {
	return ctx.Locals("session").(*session.Session)
}

func setupWebApp(app *fiber.App) {
	wa := new(webApp)

	wa.session = session.New(session.Config{
		CookieHTTPOnly: true,
	})
	wa.session.RegisterType(oauth2.Token{})

	app.Use(func(ctx *fiber.Ctx) error {

		au := &auth{}

		sess, err := wa.session.Get(ctx)
		if err != nil {
			return err
		}

		if x := sess.Get(hasAuthKey); x != nil {
			if v, ok := x.(bool); ok && v {
				au.HasAuth = true
			}
		}

		ctx.Locals("auth", au)
		ctx.Locals("session", sess)

		return ctx.Next()
	})

	app.Get("/", wa.index)
	app.Get("/auth/inbound", wa.authInbound)

	app.Use(func(ctx *fiber.Ctx) error {

		// user must be logged in

		auth := getAuth(ctx)
		if !auth.HasAuth {
			return ctx.Redirect("/")
		}

		sess := getSession(ctx)
		var guildRoles []string
		{
			gri := sess.Get(guildRolesKey)
			if gri == nil {
				return ctx.Redirect("/")
			}
			guildRoles = gri.([]string)
		}

		if common.IsAdmin(guildRoles) {
			auth.IsAdmin = true
		}

		return ctx.Next()
	})

	app.Get("/services", wa.serviceListing)

	app.Use(func(ctx *fiber.Ctx) error {
		// user must be admin

		auth := getAuth(ctx)

		if !auth.IsAdmin {
			return fiber.ErrForbidden
		}

		return ctx.Next()
	})

	app.Get("/bio", wa.bioUIDSearch)
	app.Get("/bio/view", wa.bioView)
}