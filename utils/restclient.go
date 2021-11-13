package restclient

import (
	"bytes"
	"deni/notesbot"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

// Post sends a post request to the URL with the body
func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}
	request.Header = headers
	return Client.Do(request)
}

func SendNotes(client HTTPClient, url string, notes []notesbot.Note) (*http.Response, error) {
	jsonBytes, _ := json.Marshal(notes)
	r := ioutil.NopCloser(bytes.NewReader(jsonBytes))
	request, _ := http.NewRequest("POST", url, r)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	return client.Do(request)
}
