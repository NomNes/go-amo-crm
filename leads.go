package amo

import "time"

const LeadsEntity = "leads"

type Lead struct {
	Id                     int                 `json:"id"`
	Name                   string              `json:"name"`
	Price                  int                 `json:"price,omitempty"`
	ResponsibleUserId      int                 `json:"responsible_user_id,omitempty"`
	GroupId                int                 `json:"group_id,omitempty"`
	StatusId               int                 `json:"status_id,omitempty"`
	PipelineId             int                 `json:"pipeline_id,omitempty"`
	LossReasonId           *int                `json:"loss_reason_id,omitempty"`
	SourceId               *int                `json:"source_id,omitempty"`
	CreatedBy              int                 `json:"created_by,omitempty"`
	UpdatedBy              int                 `json:"updated_by,omitempty"`
	ClosedAt               *int                `json:"closed_at,omitempty"`
	ClosestTaskAt          *int                `json:"closest_task_at,omitempty"`
	IsDeleted              bool                `json:"is_deleted,omitempty"`
	Score                  *int                `json:"score,omitempty"`
	AccountId              int                 `json:"account_id,omitempty"`
	CustomFieldValues      []CustomFieldValues `json:"custom_fields_values,omitempty"`
	Embedded               LeadEmbedded        `json:"_embedded,omitempty"`
	IsPriceModifiedByRobot *bool               `json:"is_price_modified_by_robot,omitempty"`
	EntityTime
	client *AmoCrm
}

type LeadEmbedded struct {
	Tags            []Tag         `json:"tags"`
	Companies       []IdLink      `json:"companies"`
	CatalogElements []interface{} `json:"catalog_elements,omitempty"` // TODO
	LossReason      []interface{} `json:"loss_reason,omitempty"`      // TODO
	Contacts        []LeadContact `json:"contacts,omitempty"`
}

type LeadContact struct {
	Id     int  `json:"id"`
	IsMain bool `json:"is_main"`
}

// GetLeads return slice of Lead
// Available with: WithCatalogElements, WithIsPriceModifiedByRobot, WithLossReason, WithContacts, WithOnlyDeleted
// Available order: OrderByCreatedAt, OrderByUpdatedAt, OrderById
func (a *AmoCrm) GetLeads(limit, page int, with []string, query string, order map[string]string) ([]Lead, *Pages, error) {
	var leads []Lead
	paged, err := a.getItems([]string{LeadsEntity}, &entitiesQuery{Limit: limit, Page: page, With: with, Query: query, Order: order}, &leads)
	for i := range leads {
		leads[i].client = a
	}
	return leads, paged, err
}

// GetLead return Lead by id
// Available with: WithCatalogElements, WithIsPriceModifiedByRobot, WithLossReason, WithContacts, WithOnlyDeleted
func (a *AmoCrm) GetLead(id int, with []string) (*Lead, error) {
	var lead *Lead
	err := a.getItem([]string{LeadsEntity}, &id, &entitiesQuery{With: with}, &lead)
	if err != nil {
		return nil, err
	}
	lead.client = a
	return lead, nil
}

// NewLead return Lead for current AmoCrm client
func (a *AmoCrm) NewLead(contact *Lead) *Lead {
	if contact == nil {
		contact = &Lead{}
	}
	contact.client = a
	return contact
}

// GetLeadsTags return slice of Tag for contacts
func (a *AmoCrm) GetLeadsTags(limit, page int) ([]Tag, *Pages, error) {
	return a.getTags(LeadsEntity, limit, page)
}

// GetLeadsCustomFields return slice of CustomField for contacts
func (a *AmoCrm) GetLeadsCustomFields(limit, page int) ([]CustomField, *Pages, error) {
	return a.getCustomFields(LeadsEntity, limit, page)
}

// Save add or update Lead
func (l *Lead) Save() error {
	if l.Id > 0 {
		l.SetUpdatedAtTime(time.Now())
	}
	return l.client.addItem([]string{LeadsEntity}, l, l.Id > 0, nil)
}

func (l *Lead) LinkContact(contact *Contact, isMain bool) error {
	return l.client.link(LeadsEntity, l.Id, []linkData{{
		Id:   contact.Id,
		Type: ContactsEntity,
		Metadata: &linkMetadata{
			IsMain: isMain,
		},
	}})
}
