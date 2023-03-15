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

func TestGetBankHolidays(t *testing.T) {
	Convey("Getting Bank Holidays", t, func() {
		Convey("should return a map of holiday dates", func() {

			initialize()

			err1 := setupDatabase()
			So(err1, ShouldEqual, nil)
			defer database.DB.Close()

			// Load queries
			sql, err2 := dotsql.LoadFromFile("sql/sql.sql")
			So(err2, ShouldEqual, nil)

			// test load holidays
			holidays, err3 := getBankHolidays(sql)
			So(err3, ShouldEqual, nil)
			So(len(holidays), ShouldNotEqual, 0)

			const timeFormat = "060102"

			// define input and expected result
			type testpair struct {
				date   string
				result bool
			}

			// define test data
			var tests = []testpair{
				{"160531", false}, // non holiday
				{"160430", false}, // non holiday
				{"160315", false}, // non holiday
				{"160704", true},  // holiday
				{"170904", true},  // holiday
			}

			// run tests
			for _, testData := range tests {
				_, found := holidays[testData.date]
				So(found, ShouldEqual, testData.result)
			}
		})
	})
}
