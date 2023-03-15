// payment_machine is designed to create an ACH file to debit people's bank
// accounts who's refunds did not fund, and have not paid their
// tax preparation fees.
package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dstroot/payment_machine/database"
	u "github.com/dstroot/utility"
	"github.com/gchaincl/dotsql"
	"github.com/spf13/viper"
)

/**
 * Constants
 */

const (
	timeFormat = "2006-01-02T15-04-05.000"
	dateFormat = "060102"
	shortTime  = "1504"
	debit      = "27"
	credit     = "22"
)

/**
 * Global Variables
 */

var (
	debug bool // global debug
)

func main() {
	initialize()

	err := setupDatabase()
	u.Check(err)
	defer database.DB.Close()

	// load queries
	sql, err := dotsql.LoadFromFile("sql/sql.sql")
	u.Check(err)

	// get bank holidays
	bankHolidayMap, err := getBankHolidays(sql)
	u.Check(err)

	// log bank holdays
	if debug {
		log.Println("================================================")
		log.Println("Bank Holidays")
		log.Println("================================================")
		for key, value := range bankHolidayMap {
			log.Println("Holiday Date:", key, "Value:", value)
		}
	}

	// get settlement date
	t := time.Now()
	settlementDate := u.CalcSettlementDate(t, bankHolidayMap).Format(dateFormat)
	report.SettlementDate = settlementDate

	// prepare update stmt
	update, err := sql.Prepare(database.DB, "UPDATE_PAYMENT_RECORD")
	u.Check(err)

	// get payments to process
	payments, err := getAchRecords(sql)
	u.Check(err)

	if len(payments) > 0 {

		// create output file
		f, err := os.Create(u.MakeFileName(".ach", viper.GetString("directory")))
		u.Check(err)
		defer f.Close()

		// make a write buffer
		w := bufio.NewWriter(f)
		defer w.Flush()

		/**
		 * Make File Header Record
		 * http://links.hwtreasurysolution.com/documents/NACHA-FORMAT.pdf
		 */

		// create 1 record
		rec, err := createOneRec(t)
		u.Check(err)
		if len(rec) != 96 {
			report.Errors++
			log.Printf("error - record must be 96 bytes long")
		}

		// write 1 record
		_, err = w.WriteString(rec)
		u.Check(err)
		report.TotalFileRecords++

		/**
		 * *************************************************************
		 * *************************************************************
		 * *************************************************************
		 *
		 * Make Debit batch
		 *
		 * *************************************************************
		 * *************************************************************
		 * *************************************************************
		 */

		/**
		 * Make Batch Header Record
		 * http://links.hwtreasurysolution.com/documents/NACHA-FORMAT.pdf
		 */

		// create 5 record
		rec, err = createFiveRec(settlementDate)
		u.Check(err)
		if len(rec) != 96 {
			report.Errors++
			log.Printf("error - record must be 96 bytes long")
		}

		// write 5 record
		_, err = w.WriteString(rec)
		u.Check(err)
		report.TotalFileRecords++

		// process payments
		for i := range payments {

			// create 6 record
			rec, err = createSixRec(payments[i], i)
			u.Check(err)
			if len(rec) != 96 {
				report.Errors++
				log.Printf("error - record must be 96 bytes long")
			}

			// write 6 record
			_, err = w.WriteString(rec)
			u.Check(err)
			report.TotalBatchRecords++
			report.TotalEntryRecords++
			report.TotalFileRecords++

			// create 7 record
			rec, err = createSevenRec(payments[i], i)
			u.Check(err)
			if len(rec) != 96 {
				report.Errors++
				log.Printf("error - record must be 96 bytes long")
			}

			// write 7 record
			_, err = w.WriteString(rec)
			u.Check(err)
			report.TotalBatchRecords++
			report.TotalEntryRecords++
			report.TotalFileRecords++

			// Mark record as ACH sent
			_, err = update.Exec(payments[i].paymentID)
			u.Check(err)

		}

		/**
		 * Make Batch Control Record
		 * http://links.hwtreasurysolution.com/documents/NACHA-FORMAT.pdf
		 */

		// create 8 record
		rec, err = createEightRec()
		u.Check(err)
		if len(rec) != 96 {
			report.Errors++
			log.Printf("error - record must be 96 bytes long")
		}

		// write 8 record
		_, err = w.WriteString(rec)
		u.Check(err)
		report.TotalFileRecords++

		if viper.GetBool("balancedFile") {

			/**
			 * *************************************************************
			 * *************************************************************
			 * *************************************************************
			 *
			 * Make Credit batch
			 *
			 * *************************************************************
			 * *************************************************************
			 * *************************************************************
			 */

			/**
			 * Make Batch Header Record
			 * http://links.hwtreasurysolution.com/documents/NACHA-FORMAT.pdf
			 */

			// create 5 record
			rec, err = createFiveRec(settlementDate)
			u.Check(err)
			if len(rec) != 96 {
				report.Errors++
				log.Printf("error - record must be 96 bytes long")
			}

			// write 5 record
			_, err = w.WriteString(rec)
			u.Check(err)
			report.TotalFileRecords++
			report.TotalBatchRecords = 0 // New batch, reset count
			report.BatchEntryHash = 0    // New batch, reset hash

			/**
			 * Make ACH 6 Record
			 * http://links.hwtreasurysolution.com/documents/NACHA-FORMAT.pdf
			 */

			// ENTRY DETAIL RECORD USES RECORD TYPE CODE 6.
			recordTypeCode := "6" // Len 1, Position 01-01

			// TRANSACTION CODE.  "22" is a credit to checking account
			transactionCode := credit // Len 2, Position 02-03

			// FIRST 8 DIGITS OF THE RECEIVER’S BANK TRANSIT
			// ROUTING NUMBER AT THE FINANCIAL INSTITUTION
			// WHERE THE RECEIVER'S ACCOUNT IS MAINTAINED.
			ReceivingRtn := viper.GetString("rtn")[:len(viper.GetString("rtn"))-1] // Len 8, Position 04-11

			// LAST DIGIT OF RECEIVER'S BANK TRANSIT ROUTING NUMBER.
			checkDigit := viper.GetString("rtn")[len(viper.GetString("rtn"))-1:] // Len 1, Position 12-12

			// Add up RTN for the file hash
			hash, err9 := strconv.Atoi(ReceivingRtn)
			u.Check(err9)

			report.BatchEntryHash = report.BatchEntryHash + hash
			report.FileEntryHash = report.FileEntryHash + hash

			// THIS IS THE RECEIVER’S BANK ACCOUNT NUMBER. IF
			// THE ACCOUNT NUMBER EXCEEDS 17 POSITIONS, ONLY
			// USE THE LEFT MOST 17 CHARACTERS. ANY SPACES
			// WITHIN THE ACCOUNT NUMBER SHOULD BE OMITTED
			// WHEN PREPARING THE ENTRY. THIS FIELD MUST
			// BE LEFT JUSTIFIED.
			paddedAccountNum, err10 := u.Padding(viper.GetString("dan"), 17, "left", " ") // Len 17, Position 13-29
			u.Check(err10)

			// THE AMOUNT OF THE TRANSACTION. FOR
			// PRENOTIFICATIONS, THE AMOUNT MUST BE ZERO.
			report.BatchCreditAmount = report.BatchDebitAmount
			report.FileCreditAmount = report.FileDebitAmount
			paddedCreditAmount, err11 := u.Padding(strconv.Itoa(report.BatchCreditAmount), 10, "right", "0") // Len 12, Position 21-32
			u.Check(err11)

			// THIS IS AN IDENTIFYING NUMBER BY WHICH THE
			// IDENTIFICATION NUMBER RECEIVER IS KNOWN TO THE ORIGINATOR. IT IS
			// INCLUDED FOR FURTHER IDENTIFICATION AND
			// DESCRIPTIVE PURPOSES.
			idenNumber, err12 := u.Padding("", 15, "left", " ") // Len 15, Position 40-54
			u.Check(err12)

			paddedName, err13 := u.Padding(viper.GetString("companyName"), 22, "left", " ")
			u.Check(err13)

			// DISCRETIONARY DATA - THIS FIELD MUST BE LEFT BLANK.
			discretionaryData := "  " // Len 2, Position 77-78

			// IF PPD OR CCD, ENTER 0 IN THIS FIELD TO INDICATE
			// NO ADDENDA RECORD WILL FOLLOW. IF AN ADDENDA
			// DOES FOLLOW THIS DETAIL RECORD, ENTER 1 TO
			// INDICATE A '7' RECORD WILL FOLLOW.
			adendaInd := "0" // Len 1, Position 79-79

			// THE TRACE NUMBER IS A MEANS FOR THE 80-94 CHARACTERS
			// TO IDENTIFY THE INDIVIDUAL ENTRIES. THE FIRST 8 POSITIONS
			// OF THE FIELD SHOULD BE OUR BANK TRANSIT ROUTING NUMBER
			// (WITHOUT THE CHECK DIGIT).
			//
			// THE REMAINDER OF THE FIELD MUST BE A UNIQUE NUMBER, ASSIGNED IN
			// ASCENDING ORDER FOR EACH ENTRY. TRACE NUMBERS MAY BE DUPLICATED
			// ACROSS DIFFERENT FILES.

			// First,truncate the check digit from our RTN
			ourRtn := viper.GetString("rtn")[:len(viper.GetString("rtn"))-1]

			// Then a sequential number
			paddedSeqential, err14 := u.Padding("1", 7, "right", "0")
			u.Check(err14)

			// Assemble the trace number
			traceNum := ourRtn + paddedSeqential // Len 15, Position 80-94

			// write 6 record to the file
			res, err14 := w.WriteString(recordTypeCode + transactionCode + ReceivingRtn + checkDigit + paddedAccountNum + paddedCreditAmount + idenNumber + paddedName + discretionaryData + adendaInd + traceNum + "\r\n")
			u.Check(err14)
			if res != 96 { // <- have to allow for line endings
				report.Errors++
				log.Printf("error - record must be 96 bytes long")
			}

			report.TotalBatchRecords++
			report.TotalEntryRecords++
			report.TotalFileRecords++

			/**
			 * Make Batch Control Record
			 * http://links.hwtreasurysolution.com/documents/NACHA-FORMAT.pdf
			 */

			// create 8 record
			rec, err = createEightRec()
			u.Check(err)
			if len(rec) != 96 {
				report.Errors++
				log.Printf("error - record must be 96 bytes long")
			}

			// write 8 record
			_, err = w.WriteString(rec)
			u.Check(err)
			report.TotalFileRecords++
		}

		/**
		 * Make File Control Record
		 * http://links.hwtreasurysolution.com/documents/NACHA-FORMAT.pdf
		 */

		// create 9 record
		numberOfRecordsToAdd, rec, err := createNineRec()
		u.Check(err)
		if len(rec) != 96 {
			report.Errors++
			log.Printf("error - record must be 96 bytes long")
		}

		// write 9 record
		_, err = w.WriteString(rec)
		u.Check(err)
		report.TotalFileRecords++

		/**
		 * Make Filler Records
		 * http://links.hwtreasurysolution.com/documents/NACHA-FORMAT.pdf
		 */

		// FILLER RECORDS ARE ALL 9'S
		fillerRecord, err := u.Padding("", 94, "right", "9") // Len 94, Position 01-94
		u.Check(err)

		for i := 0; i < numberOfRecordsToAdd; i++ {
			res, err := w.WriteString(fillerRecord + "\r\n")
			u.Check(err)
			if res != 96 {
				report.Errors++
				log.Printf("error - record must be 96 bytes long")
			}
			report.TotalFileRecords++
		}

		w.Flush()
		f.Close()
		report.Files++
	}

	// Log running duration
	log.Printf("Running for: %s\n", time.Since(start))
	writeReport()
	log.Println("Goodbye!")
}
