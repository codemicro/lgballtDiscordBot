package adminSite

import (
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/db"
	"github.com/gofiber/fiber/v2"
	"net/url"
	"strconv"
	"strings"
)

func (w *webApp) bioUIDSearch(ctx *fiber.Ctx) error {

	var resp strings.Builder

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

		resp.WriteString("<p>You queried: ")
		resp.WriteString(searchQuery)
		resp.WriteString("</p><br>")

		if len(resultsList) != 0 {
			resp.WriteString("Found ")
			resp.WriteString(strconv.Itoa(len(resultsList)))
			resp.WriteString(" results<Br>")
			for _, res := range resultsList {

				viewURL := fmt.Sprintf("/bio/view?user=%s", url.PathEscape(res.UserId))

				if res.SysMemberID != "" {
					viewURL += "&member=" + url.PathEscape(res.SysMemberID)
				}

				resp.WriteString("<a href='")
				resp.WriteString(viewURL)
				resp.WriteString("'>")
				resp.WriteString(res.UserId)
				resp.WriteString(" - ")
				resp.WriteString(res.SysMemberID)
				resp.WriteString("</a><br>")
			}
		} else {
			resp.WriteString("not found :(")
		}
	}

	resp.WriteString("<br><form action=''>User/System/Member ID: <input type='text' name='q'><br><input type='submit'></form>")

	return ctx.Type("html").SendString(resp.String())
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