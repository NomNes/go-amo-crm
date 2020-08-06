package go_amo_crm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (a *AmoCrm) getHost() string {
	return "https://" + a.subdomain + ".amocrm.com"
}

type ErrorRes struct {
	Hint   string `json:"hint"`
	Title  string `json:"title"`
	Type   string `json:"type"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

func (a *AmoCrm) post(path string, jsonBody interface{}, r interface{}) error {
	b, err := json.Marshal(jsonBody)
	if err != nil {
		return err
	}
	res, err := http.Post(a.getHost()+path, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		var errJson ErrorRes
		err := json.Unmarshal(body, &errJson)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("%d %s\n%s\n%s\n%s", errJson.Status, errJson.Title, errJson.Hint, errJson.Detail, errJson.Type))
	}
	log.Println(string(body))
	return json.Unmarshal(body, &r)
}

type TokenData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	Timestamp    time.Time `json:"timestamp,omitempty"`
}

func (a *AmoCrm) getToken(code string) (d *TokenData, err error) {
	err = a.post("/oauth2/access_token", map[string]string{
		"client_id":     a.clientId,
		"client_secret": a.clientSecret,
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  "https://dubai-realty.com",
	}, &d)
	if err != nil {
		return
	}
	d.Timestamp = time.Now().Add(-time.Second * 30)
	return
}

func (a *AmoCrm) refreshToken(d *TokenData) (*TokenData, error) {
	err := a.post("/oauth2/access_token", map[string]string{
		"client_id":     a.clientId,
		"client_secret": a.clientSecret,
		"grant_type":    "authorization_code",
		"refresh_token": d.RefreshToken,
		"redirect_uri":  "https://dubai-realty.com",
	}, &d)
	if err != nil {
		return d, err
	}
	d.Timestamp = time.Now().Add(-time.Second * 30)
	return d, nil
}

func (a *AmoCrm) actualizeToken(d *TokenData) (*TokenData, error) {
	if d.RefreshToken == "" || d.AccessToken == "" {
		return nil, errors.New("token not found")
	}
	if time.Now().AddDate(0, -3, 0).After(d.Timestamp) {
		return nil, errors.New("token expired")
	}
	var err error
	if time.Now().Add(-time.Second * time.Duration(d.ExpiresIn)).After(d.Timestamp) {
		d, err = a.refreshToken(d)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}
