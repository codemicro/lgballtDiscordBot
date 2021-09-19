package adminSite

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/adminSite/templates"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/gofiber/fiber/v2"
	"strings"
)

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

		page.ShowSearchResults = true
		page.SearchResults = resultsList

	}

	return ctx.Type("html").SendString(templates.RenderPage(page))
}

func (w *webApp) bioView(ctx *fiber.Ctx) error {

	userID := ctx.Query("user")
	memberID := ctx.Query("member")

	ub := new(db.UserBio)
	ub.UserId = userID
	ub.SysMemberID = memberID
	found, err := ub.Populate()
	if err != nil {
		return err
	}
	if !found {
		return fiber.ErrNotFound
	}

	var sb strings.Builder

	sb.WriteString(ub.UserId)
	sb.WriteRune('\n')
	sb.WriteString(ub.SysMemberID)
	sb.WriteRune('\n')
	sb.WriteString(fmt.Sprintf("%#v", ub.BioData))
	sb.WriteRune('\n')
	sb.WriteString(ub.ImageURL)

	return ctx.Type("txt").SendString(sb.String())
}