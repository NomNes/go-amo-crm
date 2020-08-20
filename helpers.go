package amo

import "encoding/json"

func reMarshal(s, d interface{}) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &d)
}
