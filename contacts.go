package go_amo_crm

import (
	"fmt"
)

type Contact struct {
	Id                int     `json:"id"`
	Name              string  `json:"name"`
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	ResponsibleUserId int     `json:"responsible_user_id"`
	CreatedBy         int     `json:"created_by"`
	CreatedAt         int     `json:"created_at"`
	UpdatedAt         int     `json:"updated_at"`
	AccountId         int     `json:"account_id"`
	UpdatedBy         int     `json:"updated_by"`
	GroupId           int     `json:"group_id"`
	Company           Company `json:"company"`
	ClosestTaskAt     int     `json:"closest_task_at"`
}

type Company struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type cr struct {
	Embedded struct {
		Items []Contact `json:"items"`
	} `json:"_embedded"`
}

func (a *AmoCrm) GetContacts() ([]Contact, error) {
	var r cr
	err := a.get("/api/v2/contacts", &r)
	if err != nil {
		return nil, err
	}
	return r.Embedded.Items, nil
}

func (a *AmoCrm) GetContact(id int) (*Contact, error) {
	var r cr
	err := a.get(fmt.Sprintf("/api/v2/contacts?id=%d", id), &r)
	if err != nil {
		return nil, err
	}
	if len(r.Embedded.Items) > 0 {
		return &r.Embedded.Items[0], nil
	}
	return nil, nil
}

type AddContact struct {
	Id                int                     `json:"id,omitempty"`
	Name              string                  `json:"name,omitempty"`
	FirstName         string                  `json:"first_name,omitempty"`
	LastName          string                  `json:"last_name,omitempty"`
	CreatedAt         int                     `json:"created_at,omitempty"`
	UpdatedAt         int                     `json:"updated_at,omitempty"`
	ResponsibleUserId int                     `json:"responsible_user_id,omitempty"`
	CreatedBy         int                     `json:"created_by,omitempty"`
	CompanyName       string                  `json:"company_name,omitempty"`
	Tags              string                  `json:"tags,omitempty"`
	CustomFields      []AddContactCustomField `json:"custom_fields,omitempty"`
}

type AddContactCustomField struct {
	Id     int                          `json:"id"`
	Values []AddContactCustomFieldValue `json:"values"`
}

type AddContactCustomFieldValue struct {
	Value string `json:"value"`
	Enum  string `json:"enum"`
}

func (a *AmoCrm) postContacts(action string, contact []AddContact) ([]int, error) {
	var r struct {
		Embedded struct {
			Items []struct {
				Id int `json:"id"`
			} `json:"items"`
		} `json:"_embedded"`
	}
	err := a.post("/api/v2/contacts", map[string]interface{}{
		action: contact,
	}, &r)
	if err != nil {
		return nil, err
	}

	if len(r.Embedded.Items) > 0 {
		var result []int
		for _, item := range r.Embedded.Items {
			result = append(result, item.Id)
		}
		return result, nil
	}
	return nil, nil
}

func (a *AmoCrm) AddContacts(contact []AddContact) ([]int, error) {
	return a.postContacts("add", contact)
}

func (a *AmoCrm) UpdateContacts(contact []AddContact) ([]int, error) {
	return a.postContacts("update", contact)
}
