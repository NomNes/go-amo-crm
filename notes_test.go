package amo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNotes(t *testing.T) {
	auth(t)

	lead := amo.NewLead(&Lead{
		Name: "test lead",
	})
	Convey("Add Lead", t, func() {
		err := lead.Save()
		So(err, ShouldBeNil)
		So(lead.Id, ShouldBeGreaterThan, 0)
		So(lead.UpdatedAt, ShouldBeGreaterThan, 0)
	})
	Convey("Add Note", t, func() {
		text := "Test note text"
		note := amo.NewNote(&Note{
			EntityId: lead.Id,
			NoteType: NoteTypeCommon,
			Params:   NoteParams{Text: &text},
		})
		err := note.Save(LeadsEntity)
		So(err, ShouldBeNil)
		So(note.Id, ShouldBeGreaterThan, 0)
		So(note.UpdatedAt, ShouldBeGreaterThan, 0)
	})
}
