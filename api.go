package amo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/NomNes/go-errors-sentry"
)

func (a *AmoCrm) getHost() string {
	return "https://" + a.subdomain + ".amocrm.ru"
}

type ErrorRes struct {
	Hint   string `json:"hint"`
	Title  string `json:"title"`
	Type   string `json:"type"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func (a *AmoCrm) request(method, path string, jsonBody interface{}, r interface{}, auth bool) error {
	if a.Debug {
		log.Println("amo start request", method, path)
	}
	if auth {
		err := a.Restore()
		if err != nil {
			return errors.Wrap(err)
		}
	}
	var br io.Reader
	if jsonBody != nil {
		b, err := json.Marshal(jsonBody)
		if err != nil {
			return errors.Wrap(err)
		}
		br = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, a.getHost()+path, br)
	if err != nil {
		return errors.Wrap(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "amoCRM-API-Library/1.0")
	if auth {
		if d := a.Storage.Get(); d != nil {
			req.Header.Set("Authorization", fmt.Sprintf("%s %s", d.TokenType, d.AccessToken))
		}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err)
	}

	if a.Debug {
		log.Println("amo response", method, path, res.StatusCode, string(body))
	}
	if res.StatusCode != 200 {
		var errJson ErrorRes
		err := json.Unmarshal(body, &errJson)
		if err != nil {
			return errors.Wrap(err)
		}
		return errors.New(fmt.Sprintf("%d %s\n%s\n%s\n%s", errJson.Status, errJson.Title, errJson.Hint, errJson.Detail, errJson.Type))
	}
	return errors.Wrap(json.Unmarshal(body, &r))
}

func (a *AmoCrm) post(path string, jsonBody interface{}, r interface{}, auth bool) error {
	return errors.Wrap(a.request(http.MethodPost, path, jsonBody, &r, auth))
}

func (a *AmoCrm) get(path string, r interface{}, auth bool) error {
	return errors.Wrap(a.request(http.MethodGet, path, nil, &r, auth))
}
