package main

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	// Notice that we're loading the MSSQL driver anonymously, aliasing its
	// package qualifier to _ so none of its exported names are visible
	// to our code. Under the hood, the driver registers itself as being
	// available to the database/sql package.
	_ "github.com/denisenkom/go-mssqldb" // https://github.com/denisenkom/go-mssqldb
)

func TestCreateOneRec(t *testing.T) {
	Convey("function createOneRec", t, func() {
		Convey("should return an ACH 1 record", func() {

			t := time.Now()

			rec, err := createOneRec(t)

			So(err, ShouldEqual, nil)
			So(rec, ShouldNotBeNil)
			So(len(rec), ShouldEqual, 96)
		})
	})
}

func TestCreateFiveRec(t *testing.T) {
	Convey("function createFiveRec", t, func() {
		Convey("should return an ACH 5 record", func() {

			date := "060102"

			rec, err := createFiveRec(date)

			So(err, ShouldEqual, nil)
			So(rec, ShouldNotBeNil)
			So(len(rec), ShouldEqual, 96)
		})
	})
}

func TestCreateSixRec(t *testing.T) {
	Convey("function createSixRec", t, func() {
		Convey("should return an ACH 6 record", func() {

			p := Payment{
				paymentID:       1,
				providerOrderID: "EFE432XQW14JDR0J",
				firstName:       "Dan",
				middleInitial:   "J",
				lastName:        "Stroot",
				rtn:             "231374945",
				dan:             "1234ABC12345",
				amount:          "37.88",
				ssn:             "123456789",
			}

			i := 1

			rec, err := createSixRec(p, i)

			So(err, ShouldEqual, nil)
			So(rec, ShouldNotBeNil)
			So(len(rec), ShouldEqual, 96)
		})
	})
}

func TestCreateSevenRec(t *testing.T) {
	Convey("function createSevenRec", t, func() {
		Convey("should return an ACH 7 record", func() {

			p := Payment{
				paymentID:       1,
				providerOrderID: "EFE432XQW14JDR0J",
				firstName:       "Dan",
				middleInitial:   "J",
				lastName:        "Stroot",
				rtn:             "231374945",
				dan:             "1234ABC12345",
				amount:          "37.88",
				ssn:             "123456789",
			}

			i := 1

			rec, err := createSevenRec(p, i)

			So(err, ShouldEqual, nil)
			So(rec, ShouldNotBeNil)
			So(len(rec), ShouldEqual, 96)
		})
	})
}

func TestCreateEightRec(t *testing.T) {
	Convey("function createEightRec", t, func() {
		Convey("should return an ACH 8 record", func() {

			rec, err := createEightRec()

			So(err, ShouldEqual, nil)
			So(rec, ShouldNotBeNil)
			So(len(rec), ShouldEqual, 96)
		})
	})
}

func TestCreateNineRec(t *testing.T) {
	Convey("function createNineRec", t, func() {
		Convey("should return an ACH 9 record", func() {

			number, rec, err := createNineRec()

			So(err, ShouldEqual, nil)
			So(rec, ShouldNotBeNil)
			So(number, ShouldNotBeNil)
			So(len(rec), ShouldEqual, 96)
		})
	})
}
