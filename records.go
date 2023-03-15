// payment_machine is designed to create an ACH file to debit people's bank
// accounts who's refunds did not fund, and have not paid their
// tax preparation fees.
package main

import (
	"math"
	"strconv"
	"strings"
	"time"

	u "github.com/dstroot/utility"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func createOneRec(t time.Time) (rec string, err error) {

	//THIS IS THE FIRST POSITION OF ALL RECORD FORMATS.
	// THE CODE IS UNIQUE FOR EACH RECORD TYPE.
	// THE FILE HEADER RECORD RECORD IS TYPE CODE 1.
	recordTypeCode := "1" // Len 1, Position 01-01

	// PRIORITY CODES ARE NOT USED AT THIS TIME; THIS FIELD
	// MUST CONTAIN 01
	priorityCode := "01" // Len 2, Position 02-03

	// ENTER YOUR PNC BANK TRANSIT/ROUTING NUMBER
	// 04-13 PRECEDED BY A BLANK SPACE I.E. B999999999
	immediateDestination := " " + viper.GetString("rtn") // Len 10, Position 04-13

	// THIS FIELD IDENTIFIES THE ORGANIZATION OR COMPANY
	// ORIGINATING THE FILE. THE FIELD BEGINS WITH A
	// NUMBER, USUALLY '1' AND THE ORIGINATOR'S 9-DIGIT TAX
	// ID WILL FOLLOW. IF THE FIELD CANNOT BE POPULATED
	// WITH 10 DIGITS, A BLANK AND 9 DIGITS CAN BE USED.
	immediateOrigin := viper.GetString("tin") // Len 10, Position 14-23

	// DATE WHEN THE ORIGINATOR CREATED THE FILE. THE DATE
	// MUST BE IN "YYMMDD" FORMAT.
	fileDate := t.Format(dateFormat) // Len 6, Position 24-29

	// TIME WHEN THE ORIGINATOR CREATED THE FILE. THE TIME
	// MUST BE IN "HHMM" FORMAT.
	fileTime := t.Format(shortTime) // Len 4, Position 30-33

	// THIS PROVIDES A MEANS FOR AN ORIGINATOR TO
	// DISTINGUISH BETWEEN MULTIPLE FILES CREATED ON THE
	// SAME DATE. ONLY UPPERCASE, A-Z AND NUMBERS, 0-9 ARE
	// PERMITTED.
	fileIDModifier := "A" // Len 1, Position 34-34

	// THIS FIELD INDICATES THE NUMBER OF CHARACTERS
	// CONTAINED IN EACH RECORD. THE VALUE 094 IS USED.
	recordSize := "094" // Len 2, Position 35-37

	// THIS BLOCKING FACTOR DEFINES THE NUMBER OF PHYSICAL
	// RECORDS WITHIN A FILE. THE VALUE 10 MUST BE USED.
	blockingFactor := "10" // Len 2, Position 38-39

	// THIS FIELD MUST CONTAIN 1.
	formatCode := "1" // Len 1, Position 40-40

	// BANK NAME
	immediateDesinationName, err := u.Padding("FEDERAL RESERVE BANK", 23, "left", " ") // Len 23, Position 41-63
	if err != nil {
		return "", errors.Wrap(err, "error creating immediateDesinationName")
	}

	// THIS FIELD IDENTIFIES THE ORIGINATOR OF THE FILE. THE
	// NAME OF THE ORIGINATING COMPANY SHOULD BE USED.
	immediateOriginName, err := u.Padding(viper.GetString("bankName"), 23, "left", " ") // Len 23, Position 64-86
	if err != nil {
		return "", errors.Wrap(err, "error creating immediateOriginName")
	}

	// BLANKS FILL THIS FIELD.
	referenceCode := viper.GetString("rtn")[:len(viper.GetString("rtn"))-1] // Len 8, Position 87-94

	// create 1 record
	rec = recordTypeCode +
		priorityCode +
		immediateDestination +
		immediateOrigin +
		fileDate +
		fileTime +
		fileIDModifier +
		recordSize +
		blockingFactor +
		formatCode +
		immediateDesinationName +
		immediateOriginName +
		referenceCode + "\r\n"

	return rec, nil
}

func createFiveRec(settlementDate string) (rec string, err error) {

	// BATCH HEADER RECORD USES RECORD TYPE CODE 5.
	recordTypeCode := "5" // Len 1, Position 01-01

	// THE SERVICE CLASS CODE DEFINES THE TYPE OF ENTRIES
	// CONTAINED IN THE BATCH. CODE TRANSACTION TYPE:
	// 200 ACH DEBITS AND CREDITS
	// 220 ACH CREDITS ONLY
	// 225 ACH DEBITS ONLY

	serviceClassCode := "200" // Len 3, Position 02-05

	// THIS FIELD IDENTIFIES THE COMPANY THAT HAS THE RELATIONSHIP
	// WITH THE RECEIVERS OF THE ACH TRANSACTIONS. IN ACCORDANCE
	// WITH FEDERAL REGULATION E, MOST RECEIVING FINANCIAL INSTITUTIONS
	// WILL DISPLAY THIS FIELD ON THEIR CUSTOMER'S BANK STATEMENT.
	companyName, err := u.Padding(viper.GetString("companyName"), 16, "left", " ") // Len 16, Position 05-20
	if err != nil {
		return "", errors.Wrap(err, "error creating companyName")
	}

	// REFERENCE INFORMATION FOR USE BY THE ORIGINATOR.
	companyDiscretionaryData, err := u.Padding(viper.GetString("bankName"), 20, "left", " ") // Len 20, Position 21-40
	if err != nil {
		return "", errors.Wrap(err, "error creating companyDiscretionaryData")
	}

	// THIS FIELD IDENTIFIES THE ORIGINATOR OF THE TRANSACTION VIA
	// THE ORIGINATOR'S FEDERAL TAX ID (IRS EIN). THIS FIELD BEGINS
	// WITH THE NUMBER 1, FOLLOWED BY THE COMPANY'S 9-DIGIT TAX ID
	// (WITHOUT A HYPHEN.)
	companyID := viper.GetString("tin") // Len 10, Position 41-50

	// THIS FIELD DEFINES THE TYPE OF ACH ENTRIES CONTAINED IN THE
	// BATCH. ENTER: PPD (PREARRANGED PAYMENTS AND DEPOSITS) FOR
	// CONSUMER TRANSACTIONS DESTINED TO AN INDIVIDUAL or CCD (CASH
	// CONCENTRATION OR DISBURSEMENT) FOR CORPORATE TRANSACTIONS.
	entryClassCode := "WEB" // Len 3, Position 51-53

	// THIS FIELD IS USED BY THE ORIGINATOR TO PROVIDE A DESCRIPTION
	// OF THE TRANSACTION FOR THE RECEIVER. FOR EXAMPLE, PAYROLL OR
	// DIVIDEND, ETC. IN ACCORDANCE WITH REGULATION E, MOST RECEIVING
	// BANKS WILL DISPLAY THIS FIELD ON THEIR BANK STATEMENT.
	// NOTE We use the Intuit Call Center Phone Number
	companyEntryDescription, err := u.Padding(viper.GetString("companyEntryDesc"), 10, "left", " ") // Len 10, Position 54-63
	if err != nil {
		return "", errors.Wrap(err, "error creating companyEntryDescription")
	}

	// THIS FIELD IS USED BY THE ORIGINATOR TO PROVIDE A
	// DESCRIPTIVE DATE FOR THE RECEIVER. THIS IS SOLELY FOR
	// DESCRIPTIVE PURPOSES AND WILL NOT BE USED TO CALCULATE
	// SETTLEMENT OR USED FOR POSTING PURPOSES.
	// MANY RECEIVING FINANCIAL INSTITUTIONS WILL DISPLAY THIS FIELD
	// ON THE CONSUMER'S BANK STATEMENT.
	companyDescriptiveDate := settlementDate // Len 6, Position 64-69

	// THIS REPRESENTS THE DATE ON WHICH THE ORIGINATOR INTENDS
	// A BATCH OF ENTRIES TO BE SETTLED.
	effectiveEntryDate := settlementDate // Len 6, Position 70-75

	// THIS FIELD MUST BE LEFT BLANK.
	blankEntry := "   " // Len 3, Position 76-78

	// THIS FIELD MUST CONTAIN 1
	originatorStatusCode := "1" // Len 1, Position 79-79

	// ENTER THE FIRST 8 DIGITS OF YOUR PNC BANK ABA OR TRANSIT
	// ROUTING NUMBER.
	originatingDfiID := viper.GetString("rtn")[:len(viper.GetString("rtn"))-1] // Len 8, Position 80-87

	// USED BY THE ORIGINATOR TO ASSIGN A NUMBER
	// IN ASCENDING SEQUENCE TO EACH BATCH IN THE FILE.
	report.TotalBatches++
	batchNumber, err := u.Padding(strconv.Itoa(report.TotalBatches), 7, "right", "0") // Len 7, Position 88-94
	if err != nil {
		return "", errors.Wrap(err, "error creating batchNumber")
	}

	// create 5 record
	rec = recordTypeCode +
		serviceClassCode +
		companyName +
		companyDiscretionaryData +
		companyID +
		entryClassCode +
		companyEntryDescription +
		companyDescriptiveDate +
		effectiveEntryDate +
		blankEntry +
		originatorStatusCode +
		originatingDfiID +
		batchNumber + "\r\n"

	return rec, nil
}

func createSixRec(p Payment, i int) (rec string, err error) {

	// ENTRY DETAIL RECORD USES RECORD TYPE CODE 6.
	recordTypeCode := "6" // Len 1, Position 01-01

	// TRANSACTION CODE.  "27" is a debit to checking account
	transactionCode := debit // Len 2, Position 02-03

	// FIRST 8 DIGITS OF THE RECEIVER’S BANK TRANSIT
	// ROUTING NUMBER AT THE FINANCIAL INSTITUTION
	// WHERE THE RECEIVER'S ACCOUNT IS MAINTAINED.
	ReceivingRtn := p.rtn[:len(p.rtn)-1] // Len 8, Position 04-11

	// LAST DIGIT OF RECEIVER'S BANK TRANSIT ROUTING NUMBER.
	checkDigit := p.rtn[len(p.rtn)-1:] // Len 1, Position 12-12

	// Add up RTN for the file hash
	hash, err := strconv.Atoi(ReceivingRtn)
	if err != nil {
		return "", errors.Wrap(err, "error creating hash")
	}

	report.BatchEntryHash = report.BatchEntryHash + hash
	report.FileEntryHash = report.FileEntryHash + hash

	// THIS IS THE RECEIVER’S BANK ACCOUNT NUMBER. IF
	// THE ACCOUNT NUMBER EXCEEDS 17 POSITIONS, ONLY
	// USE THE LEFT MOST 17 CHARACTERS. ANY SPACES
	// WITHIN THE ACCOUNT NUMBER SHOULD BE OMITTED
	// WHEN PREPARING THE ENTRY. THIS FIELD MUST
	// BE LEFT JUSTIFIED.
	paddedAccountNum, err := u.Padding(p.dan, 17, "left", " ") // Len 17, Position 13-29
	if err != nil {
		return "", errors.Wrap(err, "error creating paddedAccountNum")
	}

	// THE AMOUNT OF THE TRANSACTION. FOR
	// PRENOTIFICATIONS, THE AMOUNT MUST BE ZERO.

	// First, remove decimal from amount string by splitting
	// on the period and recombining the two parts
	newAmount := ""
	splitAmount := strings.Split(p.amount, ".")
	if splitAmount[1] == "" {
		newAmount = splitAmount[0] + "00"
	} else {
		newAmount = splitAmount[0] + splitAmount[1]
	}

	// Add up debit amounts
	newAmountInt, err := strconv.Atoi(newAmount)
	if err != nil {
		return "", errors.Wrap(err, "error creating newAmountInt")
	}

	// Keep debits tally
	report.BatchDebitAmount = report.BatchDebitAmount + newAmountInt
	report.FileDebitAmount = report.FileDebitAmount + newAmountInt

	// Left zero fill amount if necessary: length 10
	paddedAmount, err := u.Padding(newAmount, 10, "right", "0") // Len 10, Position 30-39
	if err != nil {
		return "", errors.Wrap(err, "error creating paddedAmount")
	}

	// THIS IS AN IDENTIFYING NUMBER BY WHICH THE
	// IDENTIFICATION NUMBER RECEIVER IS KNOWN TO THE ORIGINATOR. IT IS
	// INCLUDED FOR FURTHER IDENTIFICATION AND
	// DESCRIPTIVE PURPOSES.

	// The idenNumber is the Intuit Order Number, but with
	// the "EFE" chopped off so it will fit, or the TaxSlayer Number
	// with the EFIN removed.

	// TaxSlayer Format is:
	// Efin + Year + Julian day + 7 digit counter
	// 53003220170140024141

	idenNumber := "" // Len 15, Position 40-54
	trimmed := strings.TrimSpace(p.providerOrderID)
	if trimmed[:3] == "EFE" {
		idenNumber, err = u.Padding(trimmed[3:len(trimmed)], 15, "left", " ")
		if err != nil {
			return "", errors.Wrap(err, "error creating idenNumber")
		}
	} else {
		idenNumber, err = u.Padding(trimmed[6:len(trimmed)], 15, "left", " ")
		if err != nil {
			return "", errors.Wrap(err, "error creating idenNumber")
		}
	}

	// THIS IS THE NAME IDENTIFYING THE RECEIVER
	// NAME OF THE TRANSACTION.

	// make uppercase, truncate or pad name: length 22
	// Len 22, Position 55-76
	name := strings.ToUpper(p.firstName + " " + p.lastName)
	if len(name) > 22 {
		name = name[:22]
	}
	paddedName, err := u.Padding(name, 22, "left", " ")
	if err != nil {
		return "", errors.Wrap(err, "error creating paddedName")
	}

	// DISCRETIONARY DATA - THIS FIELD SHOULD BE LEFT BLANK.
	discretionaryData := "S " // Len 2, Position 77-78

	// IF PPD OR CCD, ENTER 0 IN THIS FIELD TO INDICATE
	// NO ADDENDA RECORD WILL FOLLOW. IF AN ADDENDA
	// DOES FOLLOW THIS DETAIL RECORD, ENTER 1 TO
	// INDICATE A '7' RECORD WILL FOLLOW.
	adendaInd := "1" // Len 1, Position 79-79

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
	paddedSeqential, err := u.Padding(strconv.Itoa(i+1), 7, "right", "0")
	if err != nil {
		return "", errors.Wrap(err, "error creating paddedSeqential")
	}

	// Assemble the trace number
	traceNum := ourRtn + paddedSeqential // Len 15, Position 80-94

	// create 6 record
	rec = recordTypeCode +
		transactionCode +
		ReceivingRtn +
		checkDigit +
		paddedAccountNum +
		paddedAmount +
		idenNumber +
		paddedName +
		discretionaryData +
		adendaInd +
		traceNum + "\r\n"

	return rec, nil
}

func createSevenRec(p Payment, i int) (rec string, err error) {
	/**
	 * Make ACH 7 Record
	 * http://links.hwtreasurysolution.com/documents/NACHA-FORMAT.pdf
	 */

	// ADDENDA RECORD USES RECORD TYPE CODE 7.
	recordTypeCode := "7" // Len 1, Position 01-01

	// ADENDA CODE.  Always "05"
	addendaCode := "05" // Len 2, Position 02-03

	// MAY CONTAIN ALPHAMERIC INFORMATION/TEXT THAT
	// FURTHER IDENTIFIES THE PURPOSE OF THE
	// TRANSACTION. THIS FIELD MAY CONTAIN INVOICE OR
	// REFERENCE NUMBERS TO HELP THE RECEIVER APPLY
	// THE TRANSACTION IN THEIR ACCOUNTING SYSTEM. 04-83
	pmtInformation := viper.GetString("companyName") + " " + viper.GetString("companyPhone") + " REF:" + p.providerOrderID

	// left justify pmtInformation
	paddedPmtInformation, err := u.Padding(pmtInformation, 80, "left", " ")
	u.Check(err)

	// ADDENDA SEQUENCE NUMBER A NUMBER CONSECUTIVELY ASSIGNED TO EACH
	// ADDENDA RECORD FOLLOWING A SEQUENCE NUMBER
	// ENTRY DETAIL RECORD. FOR CCD+ OR PPD+, THIS
	// NUMBER WILL ALWAYS BE 0001.
	addendaSequence := "0001" // Len 4, Position 84-87

	// Then a sequential number
	paddedSeqential, err1 := u.Padding(strconv.Itoa(i+1), 7, "right", "0")
	u.Check(err1)

	// create 7 record
	rec = recordTypeCode +
		addendaCode +
		paddedPmtInformation +
		addendaSequence +
		paddedSeqential + "\r\n"

	return rec, nil
}

func createEightRec() (rec string, err error) {

	// BATCH CONTROL RECORD USES RECORD TYPE CODE 8
	recordTypeCode := "8" // Len 1, Position 01-01

	// THE SERVICE CLASS CODE DEFINES THE TYPE OF ENTRIES
	// CONTAINED IN THE BATCH. CODE TRANSACTION TYPE:
	// 200 ACH DEBITS AND CREDITS
	// 220 ACH CREDITS ONLY
	// 225 ACH DEBITS ONLY

	serviceClassCode := "200" // Len 3, Position 02-05

	// COUNT IS A TALLY OF EACH TYPE ‘6’ RECORD AND IF
	// USED, ALSO EACH ADDENDA WITHIN THE BATCH
	count, err := u.Padding(strconv.Itoa(report.TotalBatchRecords), 6, "right", "0") // Len 6, Position 05-10
	u.Check(err)

	// FOR EACH ORIGINATED TRANSACTION, YOU HAVE
	// GENERATED A TYPE ‘6’ OR ENTRY DETAIL RECORD. ON THE
	// ENTRY DETAIL RECORD THERE IS A RECEIVING DEPOSITORY
	// FINANANCIAL INSTITUTION (RDFI) IDENTIFICATION (TRANSIT
	// ROUTING NUMBER) LOCATED IN POSITIONS 4 THROUGH 11.
	// THE FIRST 8 DIGITS OF EACH RDFI’s TRANSIT ROUTING
	// NUMBER IS TREATED AS A NUMBER.
	//
	// ALL TRANSIT ROUTING NUMBERS WITHIN THE BATCH ARE
	// ADDED TOGETHER FOR THE ENTRY HASH ON THE TYPE '8',
	// BATCH CONTROL RECORD. ALL TRANSIT ROUTING NUMBERS
	// WITHIN EACH FILE ARE ADDED TOGETHER TO CALCULATE THE
	// VALUE OF THE ENTRY HASH ON THE TYPE '9', FILE CONTROL
	// RECORD. (NOTE: DO NOT INCLUDE THE CHECK DIGIT OF THE
	// TRANSIT ROUTING NUMBER, POSITION 12, IN THIS
	// CALCULATION.)
	//
	// THE ENTRY HASH CALCULATIION CHECK IS
	// USED IN THE PNC BANK FILE EDITING PROCESS TO HELP
	// ENSURE DATA INTEGRITY OF THE BATCH AND FILE
	// GENERATED BY YOUR PROCESSING.
	//
	// EXAMPLE: 04300009
	//        + 03100005
	//       -----------
	//          07400014
	//
	// IN THIS EXAMPLE, THERE ARE ONLY TWO ENTRY DETAIL
	// RECORDS. THE TOTAL OF THE TWO TRANSIT ROUTING
	// NUMBERS IS LESS THAN TEN DIGITS, SO ADD ENOUGH ZEROS
	// TO THE FRONT OF THE NUMBER TO MAKE THE NUMBER TEN
	// DIGITS SO THAT 0007400014 IS THE ENTRY HASH.
	//
	// IF THE SUM OF THE RDFI TRANSIT ROUTING NUMBERS IS A
	// NUMBER GREATER THAN TEN DIGITS, REMOVE OR DROP THE
	// NUMBER OF DIGITS FROM THE LEFT SIDE OF THE NUMBER
	// UNTIL ONLY TEN DIGITS REMAIN. FOR EXAMPLE, IF THE SUM
	// OF THE TRANSIT ROUTING NUMBERS IS 234567898765,
	// REMOVE THE “23” FOR A HASH OF 4567898765.

	entryHash := strconv.Itoa(report.BatchEntryHash)
	if len(entryHash) > 10 {
		// truncate to 10 characters
		entryHash = entryHash[len(entryHash)-10:]
	} else {
		// fill to 10 characters
		entryHash, err = u.Padding(entryHash, 10, "right", "0") // Len 10, Position 11-20
		u.Check(err)
	}

	// SUM TOTAL OF ALL DEBIT AMOUNTS WITHIN BATCH’S
	// DOLLAR AMOUNT TYPE ‘6’ RECORD.

	totalDebitAmount, err := u.Padding(strconv.Itoa(report.BatchDebitAmount), 12, "right", "0") // Len 12, Position 21-32
	u.Check(err)

	// SUM TOTAL OF ALL CREDIT AMOUNTS WITHIN BATCH’S
	// DOLLAR AMOUNT TYPE ‘6’ RECORD.
	totalCreditAmount, err := u.Padding("", 12, "right", "0") // Len 12, Position 33-44
	u.Check(err)

	// THIS FIELD IDENTIFIES THE ORIGINATOR OF THE TRANSACTION VIA
	// THE ORIGINATOR'S FEDERAL TAX ID (IRS EIN). THIS FIELD BEGINS
	// WITH THE NUMBER 1, FOLLOWED BY THE COMPANY'S 9-DIGIT TAX ID
	// (WITHOUT A HYPHEN.)
	companyID := viper.GetString("tin") // Len 10, Position 41-50

	// BLANK
	messageAuthCode, err := u.Padding("", 19, "right", " ") // Len 19, Position 55-73
	u.Check(err)

	// BLANK
	reserved, err := u.Padding("", 6, "right", " ") // Len 6, Position 74-79
	u.Check(err)

	// ENTER THE FIRST 8 DIGITS OF YOUR PNC BANK ABA OR TRANSIT
	// ROUTING NUMBER.
	originatingDfiID := viper.GetString("rtn")[:len(viper.GetString("rtn"))-1] // Len 8, Position 80-87

	// USED BY THE ORIGINATOR TO ASSIGN A NUMBER
	// IN ASCENDING SEQUENCE TO EACH BATCH IN THE FILE.
	batchNumber, err := u.Padding(strconv.Itoa(report.TotalBatches), 7, "right", "0") // Len 7, Position 88-94
	u.Check(err)

	// create 8 record
	rec = recordTypeCode +
		serviceClassCode +
		count +
		entryHash +
		totalDebitAmount +
		totalCreditAmount +
		companyID +
		messageAuthCode +
		reserved +
		originatingDfiID +
		batchNumber + "\r\n"

	return rec, nil
}

func createNineRec() (numberOfRecordsToAdd int, rec string, err error) {

	// FILE CONTROL RECORD USES RECORD TYPE CODE 9
	recordTypeCode := "9" // Len 1, Position 01-01

	// BATCH COUNT
	batchCount, err := u.Padding(strconv.Itoa(report.TotalBatches), 6, "right", "0") // Len 6, Position 02-07
	u.Check(err)

	// NUMBER OF PHYSICAL BLOCKS IN THE FILE, INCLUDING FILE HEADER AND FILE CONTROL RECORDS.
	numBlocks, remainder := math.Modf(float64(report.TotalFileRecords+1) / float64(10))

	if remainder > 0 {
		numBlocks = numBlocks + 1 // because we will pad out file
		numberOfRecordsToAdd = (10 - int(remainder*10))
	}

	blockCount, err := u.Padding(strconv.Itoa(int(numBlocks)), 6, "right", "0") // Len 6, Position 08-13
	u.Check(err)

	// SUM OF ALL ‘6’ RECORDS AND ALSO '7' RECORDS, IF USED.
	entryCount, err := u.Padding(strconv.Itoa(report.TotalEntryRecords), 8, "right", "0") // Len 8, Position 14-21
	u.Check(err)

	// SUM OF ALL RECEIVING DEPOSITORY FINANCIAL INSTITUTION IDS IN EACH ‘6’ RECORD.
	// IF SUM IS MORE THAN 10 POSITIONS, TRUNCATE LEFTMOST NUMBERS.
	entryHash := strconv.Itoa(report.FileEntryHash)
	if len(entryHash) > 10 {
		// truncate to 10 characters
		entryHash = entryHash[len(entryHash)-10:]
	} else {
		// fill to 10 characters
		entryHash, err = u.Padding(entryHash, 10, "right", "0") // Len 10, Position 11-20
		u.Check(err)
	}

	// SUM TOTAL OF ALL DEBIT AMOUNTS WITHIN THE FILE
	// DOLLAR AMOUNT TYPE ‘6’ RECORD.
	fileDebitAmount, err := u.Padding(strconv.Itoa(report.FileDebitAmount), 12, "right", "0") // Len 12, Position 21-32
	u.Check(err)

	// SUM TOTAL OF ALL CREDIT AMOUNTS WITHIN THE FILE
	// DOLLAR AMOUNT TYPE ‘6’ RECORD.
	fileCreditAmount, err := u.Padding(strconv.Itoa(report.FileCreditAmount), 12, "right", "0") // Len 12, Position 33-44
	u.Check(err)

	// BLANKS
	reserved, err := u.Padding("", 39, "right", " ") // Len 39, Position 56-94
	u.Check(err)

	// create 9 record
	rec = recordTypeCode +
		batchCount +
		blockCount +
		entryCount +
		entryHash +
		fileDebitAmount +
		fileCreditAmount +
		reserved + "\r\n"

	return numberOfRecordsToAdd, rec, nil
}
