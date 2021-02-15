package verification

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"math"
	"net/http"
	"regexp"
	"time"
)

var (
	client              = new(http.Client)
	redditUsernameRegex = regexp.MustCompile(`(?m)\/?u(?:ser)?\/([\w-]+)`)
)

const redditAgeDateForm = "January 2, 2006"

func findRedditUsername(messageBody string) string {
	sub := redditUsernameRegex.FindStringSubmatch(messageBody)
	if len(sub) >= 2 {
		return sub[1]
	}
	return ""
}

func getRedditAccountAge(username string) string {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://www.reddit.com/u/%s/", username), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:86.0) Gecko/20100101 Firefox/86.0"}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return ""
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	ageString := doc.Find("#profile--id-card--highlight-tooltip--cakeday").First().Text()
	t, _ := time.Parse(redditAgeDateForm, ageString)

	since := time.Since(t)
	days := int64(math.Abs(since.Hours() / float64(24)))
	months := days / 30
	years := months / 12

	months -= years * 12
	days -= (years * 12 * 30) + (months * 30)

	return fmt.Sprintf("%dy%dm%dd", years, months, days)
}
