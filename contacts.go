package amo

import "time"

const (
	contactsEntity = "contacts"
)

type Contact struct {
	Id                int                 `json:"id"`
	Name              string              `json:"name"`
	FirstName         string              `json:"first_name"`
	LastName          string              `json:"last_name"`
	ResponsibleUserId int                 `json:"responsible_user_id"`
	GroupId           int                 `json:"group_id"`
	CreatedBy         int                 `json:"created_by"`
	UpdatedBy         int                 `json:"updated_by"`
	ClosestTaskAt     *int                `json:"closest_task_at"`
	AccountId         int                 `json:"account_id"`
	CustomFieldValues []CustomFieldValues `json:"custom_fields_values"`
	Embedded          ContactEmbedded     `json:"_embedded"`
	EntityTime
	client *AmoCrm
}

type ContactEmbedded struct {
	Tags            []Tag         `json:"tags,omitempty"`
	Companies       []IdLink      `json:"companies,omitempty"`
	Leads           []IdLink      `json:"leads,omitempty"`
	CatalogElements []interface{} `json:"catalog_elements,omitempty"` // TODO
}

// GetContacts return slice of Contact
// Available with: WithLeads, WithCustomers, WithCatalogElements
// Available order: OrderByCreatedAt, OrderByUpdatedAt, OrderById
func (a *AmoCrm) GetContacts(limit, page int, with []string, query string, order map[string]string) ([]Contact, *Pages, error) {
	var contacts []Contact
	paged, err := a.getItems([]string{contactsEntity}, &entitiesQuery{Limit: limit, Page: page, With: with, Query: query, Order: order}, &contacts)
	for i := range contacts {
		contacts[i].client = a
	}
	return contacts, paged, err
}

// GetContact return Contact by id
// Available with: WithLeads, WithCustomers, WithCatalogElements
func (a *AmoCrm) GetContact(id int, with []string) (*Contact, error) {
	var contact *Contact
	err := a.getItem([]string{contactsEntity}, &id, &entitiesQuery{With: with}, &contact)
	if err != nil {
		return nil, err
	}
	contact.client = a
	return contact, nil
}

// NewContact return Contact for current AmoCrm client
func (a *AmoCrm) NewContact(contact *Contact) *Contact {
	if contact == nil {
		contact = &Contact{}
	}
	contact.client = a
	return contact
}

// GetContactsTags return slice of Tag for contacts
func (a *AmoCrm) GetContactsTags(limit, page int) ([]Tag, *Pages, error) {
	return a.getTags(contactsEntity, limit, page)
}

// GetContactsTags return slice of CustomField for contacts
func (a *AmoCrm) GetContactsCustomFields(limit, page int) ([]CustomField, *Pages, error) {
	return a.getCustomFields(contactsEntity, limit, page)
}

// Save add or update Contact
func (c *Contact) Save() error {
	if c.Id > 0 {
		c.SetUpdatedAtTime(time.Now())
	}
	return c.client.addItem([]string{contactsEntity}, c, c.Id > 0, nil)
}
