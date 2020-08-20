package amo

import (
	"sync"
)

const VERSION = 4

type AmoCrm struct {
	subdomain    string
	clientId     string
	clientSecret string
	redirectUri  string
	Authorized   bool
	Storage
	sync.Mutex
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
		return err
	}
	err = a.Storage.Set(d)
	if err != nil {
		return err
	}
	return nil
}

func (a *AmoCrm) Restore(force bool) error {
	a.Lock()
	defer a.Unlock()
	d := a.Storage.Get()
	d, err := a.actualizeToken(d, force)
	if err != nil {
		return err
	}
	return a.Storage.Set(d)
}
