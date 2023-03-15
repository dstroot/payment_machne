package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWriteReport(t *testing.T) {
	Convey("Function writereport", t, func() {
		Convey("Should return status JSON", func() {

			err := writeReport()
			So(err, ShouldBeNil)
		})
	})
}
