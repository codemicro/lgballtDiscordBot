package pluralkit

import (
	"encoding/json"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/buildInfo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	userAgent    = []string{fmt.Sprintf("r/LGBallT Discord bot v%s (%s)", buildInfo.Version, config.PkApi.ContactEmail)}
	requestQueue = make(chan trackedRequest, 1024)

	client = new(http.Client)
)

func init() {
	for i := 0; i < config.PkApi.NumWorkers; i += 1 {
		go requestWorker()
	}
}

// trackedRequest represents a pending request
type trackedRequest struct {
	responseNotifier chan completedRequest
	request          *http.Request
}

// completedRequest represents a response and corresponding error as a result of a HTTP request
type completedRequest struct {
	response *http.Response
	err      error
}

type ApiError struct {
	StatusCode   int
	ResponseBody []byte
}

func (err *ApiError) Error() string {
	return fmt.Sprintf("pluralkit: the PK API returned a non-okay status code, %d", err.StatusCode)
}

func newApiError(statusCode int, responseBody []byte) *ApiError {
	return &ApiError{
		StatusCode:   statusCode,
		ResponseBody: responseBody,
	}
}

var responseCache = cache.New(3*time.Minute, 5*time.Minute)

// orchestrateRequest takes various parameters, makes a request and returns an error. output should be a variable that
// can be used to unmarshal response JSON into. isStatusCode should be a function that returns true if a status code is
// received that does not indicate a failed request. errorsByStatusCode is a map of errors that should be returned in
// the event a specific status code is returned from the API.
func orchestrateRequest(url string, output interface{}, isStatusCodeOk func(int) bool,
	errorsByStatusCode map[int]error) error {

	if x, found := responseCache.Get(url); found {
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

	// check status code map
	for code, err := range errorsByStatusCode {
		if resp.StatusCode == code {
			return err
		}
	}

	respBodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// check status function
	if !isStatusCodeOk(resp.StatusCode) {
		return newApiError(resp.StatusCode, respBodyContent)
	}

	// if we get here, that means we probably got a request we can use back
	// hence we can cache it
	responseCache.Set(url, &respBodyContent, cache.DefaultExpiration)

	// parse response and return error or nil
	return json.Unmarshal(respBodyContent, output)
}

// sendRequest adds a http.Request to the request queue and returns a response and corresponding error when the response
// is received.
func sendRequest(req *http.Request) (*http.Response, error) {
	responseNotifier := make(chan completedRequest)

	requestQueue <- trackedRequest{
		responseNotifier: responseNotifier,
		request:          req,
	}

	completed := <-responseNotifier
	return completed.response, completed.err
}

// requestWorker is a function that should be run as a goroutine. This actually does HTTP request dispatch and
// rate limiting.
func requestWorker() {
	for rq := range requestQueue {
		rq.request.Header["User-Agent"] = userAgent
		resp, err := client.Do(rq.request)
		rq.responseNotifier <- completedRequest{
			response: resp,
			err:      err,
		}
		time.Sleep(time.Millisecond * time.Duration(config.PkApi.MinRequestDelay))
	}
}
