package amo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func (a *AmoCrm) request(method, path string, jsonBody interface{}, r interface{}) error {
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
	if d := a.Storage.Get(); d != nil {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", d.TokenType, d.AccessToken))
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err)
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

func (a *AmoCrm) post(path string, jsonBody interface{}, r interface{}) error {
	return errors.Wrap(a.request(http.MethodPost, path, jsonBody, &r))
}

func (a *AmoCrm) get(path string, r interface{}) error {
	return errors.Wrap(a.request(http.MethodGet, path, nil, &r))
}
