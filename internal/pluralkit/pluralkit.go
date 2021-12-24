package pluralkit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/codemicro/lgballtDiscordBot/internal/buildInfo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"go.uber.org/ratelimit"
)

var (
	userAgent   = fmt.Sprintf("r/LGBallT Discord bot v%s (%s)", buildInfo.Version, config.PkApi.ContactEmail)
	ratelimiter = ratelimit.New(2) // per second

	ErrorCouldNotUnmarshal = errors.New("pluralkit: could not unmarshal error response")
)

type Error struct {
	Code           ErrorCode `json:"code"`
	Message        string    `json:"message"`
	RetryAfter     int       `json:"retry_after,omitempty"`
	HTTPStatusCode int       `json:"-"`
}

func (err *Error) Error() string {
	return fmt.Sprintf("pluralkit: the PK API returned an error response: %s (status: %d, HTTP: %d)", err.Message, err.Code, err.HTTPStatusCode)
}

var responseCache = cache.New(3*time.Minute, 5*time.Minute)

// orchestrateRequest sends a request to the PluralKit API and unmarshals the response, returning a *pluralkit.Error if
// required.
func orchestrateRequest(url string, output interface{}) error {

	if x, found := responseCache.Get(url); found {
		log.Debug().Str("url", url).Msg("PK API cache hit")
		apiResp := x.(*string)
		return json.Unmarshal([]byte(*apiResp), output)
	}

	_ = ratelimiter.Take()

	log.Debug().Str("url", url).Msg("running PK API request")

	var respBodyContent string

	err := requests.URL(url).
		AddValidator(func(resp *http.Response) error {
			if int(resp.StatusCode/100) != 2 { // if status code != 2xx

				var responseBody string
				if err := requests.ToString(&responseBody)(resp); err != nil {
					return err
				}

				e := new(Error)
				err := json.Unmarshal([]byte(responseBody), e)
				if err != nil {
					return ErrorCouldNotUnmarshal
				}
				e.HTTPStatusCode = resp.StatusCode

				log.Debug().Err(e).Msg("PK error response")

				return e
			}
			return nil
		}).
		UserAgent(userAgent).
		ToString(&respBodyContent).
		Fetch(context.Background())

	if err != nil {
		return err
	}

	// if we get here, that means we probably got a request we can use back
	// hence we can cache it
	responseCache.Set(url, &respBodyContent, cache.DefaultExpiration)

	// parse response and return error or nil
	return json.Unmarshal([]byte(respBodyContent), output)
}
