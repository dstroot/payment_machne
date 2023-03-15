// handlers.article_test.go

package main

import (
	"fmt"
	"testing"

	"github.com/dstroot/payment_machine/database"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestConfigure(t *testing.T) {
	Convey("When setting up the environment", t, func() {
		Convey("It should read the config file", func() {
			configure()
			So(viper.GetString("mssql.host"), ShouldEqual, "yms0wtdqkl.database.windows.net")
		})
	})
}

func TestInitialize(t *testing.T) {
	Convey("Initialize", t, func() {
		Convey("Should initialize our configuration", func() {
			fmt.Printf("\n\n")
			err := initialize()
			So(err, ShouldEqual, nil)
			So(debug, ShouldEqual, true)
		})
	})
}

func TestSetupDatabase(t *testing.T) {
	Convey("Good configuration", t, func() {
		Convey("It can connect to a database", func() {

			// Connect to database
			err1 := setupDatabase()
			So(err1, ShouldEqual, nil)
			defer database.DB.Close()

			// Ping database
			err2 := database.DB.Ping()
			So(err2, ShouldEqual, nil)
		})
	})

}
