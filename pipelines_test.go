package amo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPipelines(t *testing.T) {
	auth(t)

	var lastId int

	Convey("GetPipelines", t, func() {
		pipelines, _, err := amo.GetPipelines()
		So(err, ShouldBeNil)
		So(len(pipelines), ShouldBeGreaterThan, 0)
		lastId = pipelines[len(pipelines)-1].Id
	})

	Convey("GetPipeline", t, func() {
		contact, err := amo.GetPipeline(lastId)
		So(err, ShouldBeNil)
		So(contact, ShouldNotBeNil)
	})

	pipeline := amo.NewPipeline(&Pipeline{
		Name: "new test pipeline",
		Sort: 100,
		Embedded: PipelineEmbedded{
			Statuses: []PipelineStatus{
				{Name: "test status", Sort: 10, Color: "#fffd7f"},
			},
		},
	})
	Convey("Add Pipeline", t, func() {
		err := pipeline.Save()
		So(err, ShouldBeNil)
		So(pipeline.Id, ShouldBeGreaterThan, 0)
	})
	Convey("Update Pipeline", t, func() {
		id := pipeline.Id
		newName := "updated test pipeline"
		pipeline.Name = newName
		err := pipeline.Save()
		So(err, ShouldBeNil)
		So(pipeline.Id, ShouldEqual, id)
		So(pipeline.Name, ShouldEqual, newName)
	})
	// Convey("Delete Pipeline", t, func() {
	// 	err := pipeline.Delete()
	// 	So(err, ShouldBeNil)
	// })
}
