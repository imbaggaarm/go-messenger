package go_messenger

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

const kAccessToken = "access_token"

// Graph Api for facebook messenger api
type GraphApi struct {
	AccessToken string
	ApiVersion  string
	GraphUrl    string
}

type Bot struct {
	GraphApi
}

func (self *Bot) sendRaw(payload interface{}) (*http.Response, error) {
	requestEndpoint := self.GraphUrl + "/me/messages"

	client := &http.Client{
		Timeout: time.Second + 10,
	}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(payload)

	req, _ := http.NewRequest(http.MethodPost, requestEndpoint, body)
	req.Header.Add("Content-Type", "application/json")
	// Add access token to request params
	q := req.URL.Query()
	q.Add(kAccessToken, self.AccessToken)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error request:", err.Error())
		return resp, err
	}

	json.NewEncoder(body).Encode(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Error http response -> %v", resp)
		log.Println("Error http " + strconv.Itoa(resp.StatusCode) + " -> " + body.String())
	}

	return resp, err
}
