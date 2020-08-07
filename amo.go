package amo

type AmoCrm struct {
	subdomain    string
	clientId     string
	clientSecret string
	redirectUri  string
	authorized   bool
	storage      Storage
}

type Settings struct {
	Subdomain    string
	ClientId     string
	ClientSecret string
	Storage      *Storage
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
	err = a.storage.Set(d)
	if err != nil {
		return err
	}
	return nil
}

func (a *AmoCrm) Restore() error {
	d := a.storage.Get()
	d, err := a.actualizeToken(d)
	if err != nil {
		return err
	}
	return a.storage.Set(d)
}
