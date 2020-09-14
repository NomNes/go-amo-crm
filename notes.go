package amo

import (
	"fmt"
	"time"
)

const (
	notesEntity = "notes"

	NoteTypeCommon                 = "common"
	NoteTypeCallIn                 = "call_in"
	NoteTypeCallOut                = "call_out"
	NoteTypeServiceMessage         = "service_message"
	NoteTypeExtendedServiceMessage = "extended_service_message"
	NoteTypeMessageCashier         = "message_cashier"
	NoteTypeInvoicePaid            = "invoice_paid"
	NoteTypeGeolocation            = "geolocation"
	NoteTypeSmsIn                  = "sms_in"
	NoteTypeSmsOut                 = "sms_out"
)

type Note struct {
	Id                int        `json:"id"`
	EntityId          int        `json:"entity_id"`
	CreatedBy         int        `json:"created_by"`
	UpdatedBy         int        `json:"updated_by"`
	ResponsibleUserId int        `json:"responsible_user_id"`
	GroupId           int        `json:"group_id"`
	NoteType          string     `json:"note_type"`
	Params            NoteParams `json:"params"`
	AccountId         int        `json:"account_id"`
	client            *AmoCrm
	EntityTime
}

type NoteParams struct {
	Text         *string `json:"text,omitempty"`
	OriginalName *string `json:"original_name,omitempty"`
	Attachment   *string `json:"attachment,omitempty"`
	Uniq         *string `json:"uniq,omitempty"`
	Duration     *int    `json:"duration,omitempty"`
	Source       *string `json:"source,omitempty"`
	Link         *string `json:"link,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	Service      *string `json:"service,omitempty"`
	Status       *string `json:"status,omitempty"`
	IconUrl      *string `json:"icon_url,omitempty"`
	Address      *string `json:"address,omitempty"`
	Longitude    *string `json:"longitude,omitempty"`
	Latitude     *string `json:"latitude,omitempty"`
}

// GetNotes return slice of Note for passed entity type
// Available order: OrderByUpdatedAt, OrderById
func (a *AmoCrm) GetNotes(entityType string, entityId int, limit, page int, order map[string]string) ([]Note, *Pages, error) {
	var notes []Note
	e := []string{entityType}
	if entityId > 0 {
		e = append(e, fmt.Sprintf("%d", entityId))
	}
	e = append(e, notesEntity)
	paged, err := a.getItems(e, &entitiesQuery{Limit: limit, Page: page, Order: order}, &notes)
	for i := range notes {
		notes[i].client = a
	}
	return notes, paged, err
}

// GetNote return Note by id
func (a *AmoCrm) GetNote(entityType string, entityId int, id int) (*Note, error) {
	var note *Note
	e := []string{entityType}
	if entityId > 0 {
		e = append(e, fmt.Sprintf("%d", entityId))
	}
	e = append(e, notesEntity)
	err := a.getItem(e, &id, nil, &note)
	note.client = a
	return note, err
}

// NewNote return Note for current AmoCrm client
func (a *AmoCrm) NewNote(note *Note) *Note {
	if note == nil {
		note = &Note{}
	}
	note.client = a
	return note
}

// Save add or update Contact
func (n *Note) Save(entityType string) error {
	if n.Id > 0 {
		n.SetUpdatedAtTime(time.Now())
	}
	return n.client.addItem([]string{entityType, notesEntity}, n, n.Id > 0, nil)
}
