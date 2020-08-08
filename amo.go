package amo

import "github.com/NomNes/go-errors-sentry"

type AmoCrm struct {
	subdomain    string
	clientId     string
	clientSecret string
	redirectUri  string
	authorized   bool
	Storage
}

func (a *AmoCrm) Setup(subdomain, clientId, clientSecret, redirectUri string) {
	a.subdomain = subdomain
	a.clientId = clientId
	a.clientSecret = clientSecret
	a.redirectUri = redirectUri
}

func (a *AmoCrm) Auth(code string) error {
	d, err := a.getToken(code)
	if err != nil {
		return errors.Wrap(err)
	}
	err = a.Storage.Set(d)
	if err != nil {
		return errors.Wrap(err)
	}
	return nil
}

func (a *AmoCrm) Restore() error {
	d := a.Storage.Get()
	d, err := a.actualizeToken(d)
	if err != nil {
		return errors.Wrap(err)
	}
	return a.Storage.Set(d)
}
