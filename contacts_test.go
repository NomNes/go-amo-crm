package amo

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestContacts(t *testing.T) {
	auth(t)

	var lastId int

	Convey("GetContacts", t, func() {
		contacts, _, err := amo.GetContacts(0, 1, []string{"leads"}, "", nil)
		So(err, ShouldBeNil)
		So(len(contacts), ShouldBeGreaterThan, 0)
		lastId = contacts[len(contacts)-1].Id
	})

	Convey("GetContact", t, func() {
		contact, err := amo.GetContact(lastId, nil)
		So(err, ShouldBeNil)
		So(contact, ShouldNotBeNil)
	})

	contact := amo.NewContact(&Contact{
		Name: "test user",
		CustomFieldValues: []CustomFieldValues{
			{Id: 188231, Values: []FieldValue{{Value: "test@example.com", EnumCode: "WORK"}}},
		},
	})
	Convey("Add Contact", t, func() {
		err := contact.Save()
		So(err, ShouldBeNil)
		So(contact.Id, ShouldBeGreaterThan, 0)
		So(contact.UpdatedAt, ShouldBeGreaterThan, 0)
	})
	Convey("Update Contact", t, func() {
		oldUpdated := contact.UpdatedAt
		time.Sleep(time.Second)
		id := contact.Id
		newName := "updated test user"
		contact.Name = newName
		err := contact.Save()
		So(err, ShouldBeNil)
		So(contact.Id, ShouldEqual, id)
		So(contact.Name, ShouldEqual, newName)
		So(contact.UpdatedAt, ShouldBeGreaterThan, oldUpdated)
	})
	Convey("Contacts Tags", t, func() {
		_, _, err := amo.GetContactsTags(0, 1)
		So(err, ShouldBeNil)
	})
	Convey("Contacts Custom Fields", t, func() {
		_, _, err := amo.GetContactsCustomFields(0, 1)
		So(err, ShouldBeNil)
	})
}
