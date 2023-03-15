// payment_machine is designed to create an ACH file to debit people's bank
// accounts who's refunds did not fund, and have not paid their
// tax preparation fees.
package main

import (
	"errors"
	"time"

	"github.com/dstroot/payment_machine/database"
	"github.com/gchaincl/dotsql"
	// Notice that we're loading the MSSQL driver anonymously, aliasing its
	// package qualifier to _ so none of its exported names are visible
	// to our code. Under the hood, the driver registers itself as being
	// available to the database/sql package.
	_ "github.com/denisenkom/go-mssqldb" // https://github.com/denisenkom/go-mssqldb
)

// getBankHolidays reads the SQL database for a list of bank holidays and
// loads them into a map for a quick and easy way to check if a date is a
// bank holiday.
func getBankHolidays(sql *dotsql.DotSql) (bankHolidayMap map[string]bool, err error) {

	// Make a map of bank holidays
	bankHolidayMap = make(map[string]bool)

	// Get the refund transaction fees
	rows, err := sql.Query(database.DB, "GET_BANK_HOLIDAYS")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bankHolday time.Time

	// iterate the holidays and load them into the map
	var i = 0
	for rows.Next() {
		i++

		err1 := rows.Scan(&bankHolday)
		if err1 != nil {
			return nil, err1
		}

		date := bankHolday.Format(dateFormat)
		bankHolidayMap[date] = true
	}
	err2 := rows.Err()
	if err2 != nil {
		return nil, err2
	}

	if i == 0 {
		return nil, errors.New("no bank holidays loaded")
	}

	return bankHolidayMap, nil

}
