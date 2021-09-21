package adminSite

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/adminSite/templates"
	"github.com/codemicro/lgballtDiscordBot/internal/bot/components/bios"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"strings"
)

func getBioFromRequest(ctx *fiber.Ctx) (*db.UserBio, bool, error) {

	userID := ctx.Query("user")
	memberID := ctx.Query("member")

	ub := new(db.UserBio)
	ub.UserId = userID
	ub.SysMemberID = memberID
	found, err := ub.Populate()
	if err != nil {
		return nil, false, err
	}

	return ub, found, nil
}

func (w *webApp) bioUIDSearch(ctx *fiber.Ctx) error {

	page := &templates.BioSearchPage{}

	searchQuery := ctx.Query("q")
	if searchQuery != "" {

		biosByUser, err := db.GetBiosForAccount(searchQuery)
		if err != nil {
			return err
		}
		biosBySystem, err := db.GetBiosForSystem(searchQuery)
		if err != nil {
			return err
		}
		biosByMember, err := db.GetBiosForMember(searchQuery)
		if err != nil {
			return err
		}

		resultsList := append(biosByUser, biosByMember...)
		resultsList = append(resultsList, biosBySystem...)

		if len(resultsList) == 1 {
			return ctx.Redirect(templates.ViewURL(resultsList[0]))
		}

		page.ShowSearchResults = true
		page.SearchResults = resultsList

	}

	return ctx.Type("html").SendString(templates.RenderPage(page))
}

func (w *webApp) bioView(ctx *fiber.Ctx) error {

	ub, found, err := getBioFromRequest(ctx)
	if err != nil {
		return err
	}
	if !found {
		return fiber.ErrNotFound
	}

	return ctx.Type("html").SendString(templates.RenderPage(&templates.BioViewPage{Bio: *ub}))
}

func (w *webApp) bioEditField(ctx *fiber.Ctx) error {

	dontCache(ctx)

	field := ctx.Query("field")
	if field == "" {
		return fiber.ErrBadRequest
	}

	{
		var cf string
		for _, ix := range config.BioFields {
			if strings.EqualFold(ix, field) {
				cf = ix
				goto ok
			}
		}
		return fiber.ErrBadRequest
	ok:
		field = cf
	}

	ub, found, err := getBioFromRequest(ctx)
	if err != nil {
		return err
	}
	if !found {
		return fiber.ErrNotFound
	}

	if ctx.Method() == "POST" {
		// do something
		newContent := ctx.FormValue("new")
		if newContent == "" {
			delete(ub.BioData, field)
		} else {
			ub.BioData[field] = newContent
		}

		if bios.ShouldRemoveBio(ub) {
			err = ub.Delete()
		} else {
			err = ub.Save()
		}

		if err != nil {
			return err
		}

		nextURL := fmt.Sprintf(
			"/bio/view?user=%s&member=%s",
			url.QueryEscape(ctx.Query("user")),
			url.QueryEscape(ctx.Query("member")),
		)

		return ctx.Type("html").SendString(templates.RenderPage(&templates.FeedbackPage{
			WasSuccess:        true,
			Message:           "Saved successfully!",
			NextURL:           nextURL,
			RedirectTimeoutMs: 3000,
		}))
	}

	return ctx.Type("html").SendString(templates.RenderPage(&templates.BioEditPage{
		FieldName:      field,
		InitialContent: ub.BioData[field],
	}))
}

func (w *webApp) bioEditImage(ctx *fiber.Ctx) error {

	dontCache(ctx)

	ub, found, err := getBioFromRequest(ctx)
	if err != nil {
		return err
	}
	if !found {
		return fiber.ErrNotFound
	}

	if ctx.Method() == "POST" {
		// do something
		newContent := ctx.FormValue("new")

		// validate URL
		_, err := url.ParseRequestURI(newContent)
		if err != nil {
			return ctx.Type("html").SendString(templates.RenderPage(&templates.FeedbackPage{
				WasSuccess:        false,
				Message:           "Invalid URL format.",
				NextURL:           ctx.OriginalURL(),
				RedirectTimeoutMs: 3000,
			}))
		}

		ub.ImageURL = newContent

		if bios.ShouldRemoveBio(ub) {
			err = ub.Delete()
		} else {
			err = ub.Save()
		}

		if err != nil {
			return err
		}

		nextURL := fmt.Sprintf(
			"/bio/view?user=%s&member=%s",
			url.QueryEscape(ctx.Query("user")),
			url.QueryEscape(ctx.Query("member")),
		)

		return ctx.Type("html").SendString(templates.RenderPage(&templates.FeedbackPage{
			WasSuccess:        true,
			Message:           "Saved successfully!",
			NextURL:           nextURL,
			RedirectTimeoutMs: 1000,
		}))
	}

	return ctx.Type("html").SendString(templates.RenderPage(&templates.BioEditPage{
		FieldName:      "image URL",
		InitialContent: ub.ImageURL,
	}))
}