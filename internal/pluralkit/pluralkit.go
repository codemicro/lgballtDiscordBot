package pluralkit

import (
	"encoding/json"
	"fmt"
	"github.com/codemicro/lgballtDiscordBot/internal/buildInfo"
	"github.com/codemicro/lgballtDiscordBot/internal/config"
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
		go requestSpinner()
	}
}

type trackedRequest struct {
	responseNotifier chan completedRequest
	request          *http.Request
}

type completedRequest struct {
	response *http.Response
	err      error
}

func makeRequest(req *http.Request) (*http.Response, error) {
	responseNotifier := make(chan completedRequest)

	requestQueue <- trackedRequest{
		responseNotifier: responseNotifier,
		request:          req,
	}

	completed := <-responseNotifier
	return completed.response, completed.err
}

func requestSpinner() {
	for {
		rq := <-requestQueue
		rq.request.Header["User-Agent"] = userAgent
		resp, err := client.Do(rq.request)
		rq.responseNotifier <- completedRequest{
			response: resp,
			err:      err,
		}
		time.Sleep(time.Millisecond * time.Duration(config.PkApi.MinRequestDelay))
	}
}

func parseJsonResponse(resp *http.Response, output interface{}) error {
	respBodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(respBodyContent))
	return json.Unmarshal(respBodyContent, output)
}
