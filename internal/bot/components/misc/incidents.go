package misc

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/codemicro/dgo-toolkit/route"
	"github.com/golang/freetype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"strconv"
	"time"
)

//go:embed lgballtMoment.png
var incidentImage []byte

//go:embed liberationsans.ttf
var liberationSans []byte

func makeIncidentImage(daysSince string, output io.Writer) error {

	sourceImage := bytes.NewBuffer(incidentImage)
	rawImage, _, err := image.Decode(sourceImage)
	if err != nil {
		return err
	}

	b := rawImage.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), rawImage, b.Min, draw.Src)

	const (
		x = 170
		y = 250
	)

	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	parsedFont, err := freetype.ParseFont(liberationSans)
	if err != nil {
		return err
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(parsedFont)
	c.SetFontSize(75)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)

	if _, err = c.DrawString(daysSince, point); err != nil {
		return err
	}

	return png.Encode(output, img)
}

func getTimeSinceLastIncident() (time.Duration, error, bool) {

	dat, err := ioutil.ReadFile("lastIncident")
	if err != nil {
		return 0, err, false
	}

	i, err := strconv.ParseInt(string(dat), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("corrupted incident date: %s", err.Error()), true
	}
	previousTime := time.Unix(i, 0)

	return time.Since(previousTime), nil, false
}

func (*Misc) SinceLastIncident(ctx *route.MessageContext) error {

	timeSince, err, userErr := getTimeSinceLastIncident()
	if err != nil {
		if userErr {
			return ctx.SendErrorMessage(err.Error())
		} else {
			return err
		}
	}

	daysSince := int64(timeSince.Hours()) / 24
	daysSinceString := strconv.FormatInt(daysSince, 10)

	bb := new(bytes.Buffer)
	err = makeIncidentImage(daysSinceString, bb)
	if err != nil {
		return err
	}

	_, err = ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
		Content:         fmt.Sprintf("%d days since last incident", daysSince),
		Files:           []*discordgo.File{
			{
				Name:        "daysSince.png",
				ContentType: "image/png",
				Reader:      bb,
			},
		},
		AllowedMentions: &discordgo.MessageAllowedMentions{},
	})
	return err
}

func (*Misc) ResetSinceLastIncident(ctx *route.MessageContext) error {

	timeSince, err, userErr := getTimeSinceLastIncident()
	if err != nil {
		if userErr {
			return ctx.SendErrorMessage(err.Error())
		} else {
			return err
		}
	}

	daysSince := int64(timeSince.Hours()) / 24

	currentTime := time.Now().Unix()
	asBytes := []byte(strconv.FormatInt(currentTime, 10))

	if err := ioutil.WriteFile("lastIncident", asBytes, 0644); err != nil {
		return err
	}

	_, err = ctx.SendMessageString(ctx.Message.ChannelID, fmt.Sprintf("Counter reset :(\nIt's been %d days since the last incident", daysSince))
	return err
}