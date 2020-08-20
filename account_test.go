package amo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAccount(t *testing.T) {
	auth(t)
	Convey("Get Account", t, func() {
		_, err := amo.GetAccount(nil)
		So(err, ShouldBeNil)
	})
}
