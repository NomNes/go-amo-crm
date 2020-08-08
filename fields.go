package amo

import "github.com/NomNes/go-errors-sentry"

type Field struct {
	Id          int               `json:"id"`
	Name        string            `json:"name"`
	FieldType   int               `json:"field_type"`
	Sort        int               `json:"sort"`
	IsMultiple  bool              `json:"is_multiple"`
	IsSystem    bool              `json:"is_system"`
	IsEditable  bool              `json:"is_editable"`
	IsRequired  bool              `json:"is_required"`
	IsDeletable bool              `json:"is_deletable"`
	IsVisible   bool              `json:"is_visible"`
	Enums       map[string]string `json:"enums"`
	Code        string            `json:"code"`
}

type CustomFields struct {
	Contacts  map[string]Field `json:"contacts"`
	Leads     map[string]Field `json:"leads"`
	Companies map[string]Field `json:"companies"`
	Customers map[string]Field `json:"customers"`
}

func (a *AmoCrm) GetFields() (*CustomFields, error) {
	account, err := a.GetAccount([]string{"custom_fields"})
	if err != nil {
		return nil, errors.Wrap(err)
	}
	return account.Embedded.CustomFields, nil
}
