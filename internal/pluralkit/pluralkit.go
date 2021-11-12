package pluralkit

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/buildInfo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
	"go.uber.org/ratelimit"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	userAgent   = []string{fmt.Sprintf("r/LGBallT Discord bot v%s (%s)", buildInfo.Version, config.PkApi.ContactEmail)}
	client      = new(http.Client)
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
		apiResp := x.(*[]byte)
		return json.Unmarshal(*apiResp, output)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := sendRequest(req)
	defer func() {
		if resp != nil {
			if resp.Body != nil {
				_ = resp.Body.Close() // goroutine leaks begone!
			}
		}
	}()
	if err != nil {
		return err
	}

	respBodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 { // TODO: This could cause problems in cases where the API returns codes like 204, which it does do. We don't use that at the moment, so I'm taking the lazy option.
		e := new(Error)
		err := json.Unmarshal(respBodyContent, e)
		if err != nil {
			return ErrorCouldNotUnmarshal
		}
		e.HTTPStatusCode = resp.StatusCode
		return e
	}

	// if we get here, that means we probably got a request we can use back
	// hence we can cache it
	responseCache.Set(url, &respBodyContent, cache.DefaultExpiration)

	// parse response and return error or nil
	return json.Unmarshal(respBodyContent, output)
}

// sendRequest sends an HTTP request while obeying a rate limit.
func sendRequest(req *http.Request) (*http.Response, error) {
	_ = ratelimiter.Take()
	log.Debug().Str("url", req.URL.String()).Msg("running PK API request")
	req.Header["User-Agent"] = userAgent
	return client.Do(req)
}
