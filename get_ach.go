// payment_machine is designed to create an ACH file to debit people's bank
// accounts who's refunds did not fund, and have not paid their
// tax preparation fees.
package main

import (
	"github.com/dstroot/payment_machine/database"
	"github.com/gchaincl/dotsql"
	// Notice that we're loading the MSSQL driver anonymously, aliasing its
	// package qualifier to _ so none of its exported names are visible
	// to our code. Under the hood, the driver registers itself as being
	// available to the database/sql package.
	_ "github.com/denisenkom/go-mssqldb" // https://github.com/denisenkom/go-mssqldb
)

// Payment defines a payment record
type Payment struct {
	paymentID       int
	providerOrderID string
	firstName       string
	middleInitial   string
	lastName        string
	rtn             string
	dan             string
	amount          string
	ssn             string
}

// GetAchRecords reads the sql database for a list of payment records to process
// and loads them into a slice of payments.  The way we can check the length of
// the slice before we create a file and write the records.
func getAchRecords(sql *dotsql.DotSql) (payments []Payment, err error) {

	// make slice payments of Payment struct, zero length
	payments = make([]Payment, 0)
	var p Payment

	// Get records to process
	rows, err := sql.Query(database.DB, "GET_PAYMENT_RECORDS")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// interate the records and fill the slice
	for rows.Next() {

		// scan the row into p
		err1 := rows.Scan(
			&p.paymentID,
			&p.providerOrderID,
			&p.firstName,
			&p.middleInitial,
			&p.lastName,
			&p.rtn,
			&p.dan,
			&p.amount)
		if err1 != nil {
			return nil, err1
		}

		// append to slice
		payments = append(payments, p)
	}
	err2 := rows.Err()
	if err2 != nil {
		return nil, err2
	}

	return payments, nil
}
