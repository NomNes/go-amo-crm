package amo

import (
	"errors"
	"time"
)

type TokenData struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	Timestamp    time.Time `json:"timestamp,omitempty"`
}

func (a *AmoCrm) getToken(code string) (*TokenData, error) {
	var d *TokenData
	err := a.post("/oauth2/access_token", map[string]string{
		"client_id":     a.clientId,
		"client_secret": a.clientSecret,
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  "https://dubai-realty.com",
	}, &d, false)
	if err != nil {
		return nil, err
	}
	d.Timestamp = time.Now().Add(-time.Second * 30)
	return d, nil
}

func (a *AmoCrm) refreshToken(d *TokenData) (*TokenData, error) {
	err := a.post("/oauth2/access_token", map[string]string{
		"client_id":     a.clientId,
		"client_secret": a.clientSecret,
		"grant_type":    "refresh_token",
		"refresh_token": d.RefreshToken,
		"redirect_uri":  "https://dubai-realty.com",
	}, &d, false)
	if err != nil {
		return d, err
	}
	d.Timestamp = time.Now().Add(-time.Second * 30)
	return d, nil
}

func (a *AmoCrm) actualizeToken(d *TokenData, force bool) (*TokenData, error) {
	if d == nil || d.RefreshToken == "" || d.AccessToken == "" {
		return nil, errors.New("token not found")
	}
	var err error
	if !force {
		if time.Now().AddDate(0, -3, 0).After(d.Timestamp) {
			return nil, errors.New("token expired")
		}
		e := d.Timestamp.Add(time.Duration(int64(time.Second) * int64(d.ExpiresIn)))
		if !time.Now().After(e) {
			return d, nil
		}
	}
	d, err = a.refreshToken(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}
