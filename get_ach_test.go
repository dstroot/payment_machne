package main

import (
	"testing"

	"github.com/dstroot/payment_machine/database"
	"github.com/gchaincl/dotsql"
	. "github.com/smartystreets/goconvey/convey"
	// Notice that we're loading the MSSQL driver anonymously, aliasing its
	// package qualifier to _ so none of its exported names are visible
	// to our code. Under the hood, the driver registers itself as being
	// available to the database/sql package.
	_ "github.com/denisenkom/go-mssqldb" // https://github.com/denisenkom/go-mssqldb
)

func TestGetAchRecords(t *testing.T) {

	Convey("function getACHRecords", t, func() {
		Convey("should return a slice of records", func() {

			initialize()

			err1 := setupDatabase()
			So(err1, ShouldEqual, nil)
			defer database.DB.Close()

			// Load queries
			sql, err2 := dotsql.LoadFromFile("sql/sql.sql")
			So(err2, ShouldEqual, nil)

			// test load payments
			payments, err3 := getAchRecords(sql)
			So(err3, ShouldEqual, nil)
			So(payments, ShouldNotBeNil)

		})
	})
}
