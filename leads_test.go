package amo

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLeads(t *testing.T) {
	auth(t)

	var lastId int

	Convey("GetLeads", t, func() {
		leads, _, err := amo.GetLeads(0, 1, []string{"leads"}, "", nil)
		So(err, ShouldBeNil)
		So(len(leads), ShouldBeGreaterThan, 0)
		lastId = leads[len(leads)-1].Id
	})

	Convey("GetLead", t, func() {
		contact, err := amo.GetLead(lastId, nil)
		So(err, ShouldBeNil)
		So(contact, ShouldNotBeNil)
	})

	lead := amo.NewLead(&Lead{
		Name: "test lead",
	})
	Convey("Add Lead", t, func() {
		err := lead.Save()
		So(err, ShouldBeNil)
		So(lead.Id, ShouldBeGreaterThan, 0)
		So(lead.UpdatedAt, ShouldBeGreaterThan, 0)
	})
	Convey("Update Lead", t, func() {
		oldUpdated := lead.UpdatedAt
		time.Sleep(time.Second)
		id := lead.Id
		newName := "updated test lead"
		lead.Name = newName
		err := lead.Save()
		So(err, ShouldBeNil)
		So(lead.Id, ShouldEqual, id)
		So(lead.Name, ShouldEqual, newName)
		So(lead.UpdatedAt, ShouldBeGreaterThan, oldUpdated)
	})
	Convey("Leads Tags", t, func() {
		_, _, err := amo.GetLeadsTags(0, 1)
		So(err, ShouldBeNil)
	})
	Convey("Leads Custom Fields", t, func() {
		_, _, err := amo.GetLeadsCustomFields(0, 1)
		So(err, ShouldBeNil)
	})
}
