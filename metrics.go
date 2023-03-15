package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dstroot/payment_machine/database"
	"github.com/dstroot/utility"
	"github.com/pkg/errors"
)

var (
	report     Metrics            // global metrics
	start      = time.Now().UTC() // global start
	buildstamp = "not set"
	githash    = "not set"
)

// Metrics holds our metrics for reporting
type Metrics struct {
	Program           string
	Buildstamp        string
	GitHash           string
	GoVersion         string
	RunTime           string
	SettlementDate    string
	Files             int
	TotalBatches      int
	TotalBatchRecords int
	TotalEntryRecords int
	TotalFileRecords  int
	BatchDebitAmount  int
	BatchCreditAmount int
	FileDebitAmount   int
	FileCreditAmount  int
	BatchEntryHash    int
	FileEntryHash     int
	Errors            int
	DBconnections     int
}

// status writes a JSON object with the current metrics
func writeReport() (err error) {
	report.RunTime = fmt.Sprintf("%v", utility.RoundDuration(time.Since(start), time.Second))
	report.DBconnections = database.DB.Stats().OpenConnections
	res, err := json.MarshalIndent(report, "", "    ")
	if err != nil {
		return errors.Wrap(err, "error marshalling json")
	}
	fmt.Println(string(res))
	return nil
}
