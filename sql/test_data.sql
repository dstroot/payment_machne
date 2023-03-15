---------------------------------------------
-- Table created to contain unfunded records
-- for collection.
---------------------------------------------

IF EXISTS (SELECT * FROM SYS.TABLES WHERE NAME='COLLECTION_ITEM')
BEGIN
  DROP TABLE dbo.COLLECTION_ITEM;
END;

CREATE TABLE dbo.COLLECTION_ITEM
(
  COLLECTION_ITEM_ID      BIGINT IDENTITY (1,1)  NOT NULL, -- SK/PK
  REFUND_TRANSACTION_ID   BIGINT             NOT NULL ,  -- RAL/Refund transaction ID
  PROVIDER_ORDER_ID         VARCHAR(16)        NOT NULL ,  -- optional for Intuit (EFE#)
  FIRST_NAME              VARCHAR(56)        NOT NULL ,
  MIDDLE_INITIAL          VARCHAR(1)         NULL ,
  LAST_NAME               VARCHAR(56)        NOT NULL ,
  ADDRESS1                VARCHAR(56)        NOT NULL ,
  ADDRESS2                VARCHAR(56)        NULL ,
  CITY                    VARCHAR(56)        NOT NULL , -- City/locality
  STATE                   VARCHAR(2)         NOT NULL ,  -- State/province/territory
  ZIPCODE                 VARCHAR(5)         NOT NULL ,  -- More international
  ZIPCODE_PLUS            VARCHAR(4)         NULL ,
  COUNTRY                 VARCHAR(56)        NULL ,
  EMAIL                   VARCHAR(56)        NOT NULL ,
  MOBILE_PHONE            VARCHAR(15)        NULL ,
  AMOUNT_OWED             NUMERIC(8,2)       NOT NULL ,  -- Limited by ACH field size to 10
  PAYMENT_URL             VARCHAR(150)       NULL ,
  RTN                     VARCHAR(9)         NOT NULL ,  -- must be 9 digits
  DAN                     VARCHAR(17)        NOT NULL , -- up to 17 characters
  HASH                    VARCHAR(60)        NULL ,
  TOKEN_EXPIRES           BIGINT             NULL , -- was DATETIME, now LINUX format
  PAID_IND                BIT                NOT NULL DEFAULT 0 ,
  CREATION_DATE           DATETIME           NOT NULL DEFAULT (CURRENT_TIMESTAMP) ,
  MODIFY_DATE             DATETIME           NOT NULL DEFAULT (CURRENT_TIMESTAMP) ,
  CREATED_BY              VARCHAR(56)        NULL DEFAULT SUSER_SNAME() ,
  MODIFIED_BY             VARCHAR(56)        NULL DEFAULT SUSER_SNAME() ,
  TAX_AMOUNT_OWED         NUMERIC(8,2)       NULL ,
  PROCESS_TYPE_ID         INT                NOT NULL -- REFERENCES PROCESS,
);

ALTER TABLE dbo.COLLECTION_ITEM
  ADD CONSTRAINT [PK_COLLECTION_ITEM] PRIMARY KEY  CLUSTERED ([COLLECTION_ITEM_ID] ASC);


CREATE NONCLUSTERED INDEX [IDX_REFUND_TRANSACTION_ID] ON [COLLECTION_ITEM]
(
  [REFUND_TRANSACTION_ID]         ASC
);


CREATE NONCLUSTERED INDEX [IDX_PROVIDER_ORDER_ID] ON [COLLECTION_ITEM]
(
  [PROVIDER_ORDER_ID]         ASC
);


CREATE NONCLUSTERED INDEX [IDX_RTN] ON [COLLECTION_ITEM]
(
  [RTN]         ASC
);


CREATE NONCLUSTERED INDEX [IDX_CREATION_DATE] ON [COLLECTION_ITEM]
(
  [CREATION_DATE]         ASC
);

INSERT INTO dbo.[COLLECTION_ITEM] (REFUND_TRANSACTION_ID,PROVIDER_ORDER_ID,FIRST_NAME,MIDDLE_INITIAL,LAST_NAME,ADDRESS1,CITY,STATE,ZIPCODE,EMAIL,MOBILE_PHONE,AMOUNT_OWED,RTN,DAN,PROCESS_TYPE_ID) VALUES
  (23229117,'EFE432XQW14JDR0J','Dan','J','stroot','1223 Main Street','New Brunswick','CA','92101','dan.stroot@sbtpg.com','19494634044',49.95,'231374945','1234ABC12345',10),
  (24145448,'EFE432XQX14JDONH','Daniel','X','Stroot','1223 Main Street','New Brunswick','CA','92101','dan.stroot@sbtpg.com','19494634044',39.95,'123456789','1234ABC12345',10),
  (24274941,'EFE432XQX14JDONH','Daniel','X','Stroot','1223 Main Street','New Brunswick','CA','92101','dan.stroot@sbtpg.com','19494634044',39.95,'123456789','1234ABC12345',10),
  (23430362,'EFE432XQX14JDONH','Daniel','X','Stroot','1223 Main Street','New Brunswick','CA','92101','dan.stroot@sbtpg.com','19494634044',39.95,'123456789','1234ABC12345',10),
  (23390200,'EFE432XQX14JDONH','Daniel','X','Stroot','1223 Main Street','New Brunswick','CA','92101','dan.stroot@sbtpg.com','19494634044',39.95,'123456789','1234ABC12345',10);


---------------------------------------------
-- Table created to store payments
-- for each collection item.
---------------------------------------------

IF EXISTS (SELECT * FROM SYS.TABLES WHERE NAME='COLLECTION_PAYMENT')
BEGIN
  DROP TABLE dbo.COLLECTION_PAYMENT;
END;

CREATE TABLE dbo.COLLECTION_PAYMENT
(
  COLLECTION_PAYMENT_ID   BIGINT IDENTITY (1,1) NOT NULL ,
  REFUND_TRANSACTION_ID   BIGINT          NOT NULL,
  PAYMENT_AMOUNT          NUMERIC(8,2)    NOT NULL ,  -- Limited by ACH field size to 10
  PAYMENT_DATE            DATETIME        NOT NULL ,
  PAYMENT_TYPE_ID         INT             NOT NULL ,   -- (1) Bank ACH, (2) Credit Card See payment type table
  PAYMENT_CC_JSON         VARCHAR(MAX)    NULL ,  -- JSON response from stripe
  CUSTOMER_INITIATED_IND  BIT             NOT NULL DEFAULT 0,  -- t/f, cust approved bank ACH?
  ACH_SENT_IND            BIT             NULL ,  -- t/f ACH processed
  CREATION_DATE           DATETIME        NOT NULL DEFAULT (CURRENT_TIMESTAMP) ,
  MODIFY_DATE             DATETIME        NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  CREATED_BY              VARCHAR(56)     NULL DEFAULT SUSER_SNAME(),
  MODIFIED_BY             VARCHAR(56)     NULL DEFAULT SUSER_SNAME(),
  ACH_REJECTED_IND        BIT             NOT NULL DEFAULT 0,
  ACH_RETRY_IND           BIT             NOT NULL DEFAULT 0
);

INSERT INTO COLLECTION_PAYMENT (REFUND_TRANSACTION_ID, PAYMENT_AMOUNT, PAYMENT_DATE, PAYMENT_TYPE_ID, PAYMENT_CC_JSON, CUSTOMER_INITIATED_IND, ACH_SENT_IND, CREATION_DATE) VALUES
  (23229117,37.88,'2016-03-08 23:05:10',1,NULL,1,0,'2016-03-09 07:59:59'),
  (24145448,32.46,'2016-03-08 23:08:00',1,NULL,1,0,'2016-03-09 23:05:10'),
  (24274941,59.53,'2016-03-08 23:17:00',1,NULL,0,0,'2016-03-14 07:59:59'),
  (23430362,76.30,'2016-03-08 23:22:25',2,Null,1,0,'2016-03-09 07:59:59'),
  (23390200,78.36,'2016-03-08 23:32:57',2,Null,1,0,'2016-03-09 07:59:59');
--   (23660902,77.74,'2016-03-08 23:37:28',1,NULL,1,0),
--   (24164276,32.46,'2016-03-08 23:46:53',2,Null,1,0),
--   (23855341,77.74,'2016-03-08 23:53:44',2,Null,1,0),
--   (24189766,37.79,'2016-03-08 23:53:53',1,NULL,1,0),
--   (23725861,124.29,'2016-03-08 23:54:12',1,NULL,1,0),
--   (24226710,72.70,'2016-03-09 00:02:44',1,NULL,1,0),
--   (23658206,38.35,'2016-03-09 00:04:41',1,NULL,1,0),
--   (24252093,124.29,'2016-03-09 00:15:51',2,null,1,0),
--   (24240525,99.34,'2016-03-09 00:17:14',1,NULL,1,0),
--   (24321116,31.34,'2016-03-09 00:17:43',2,null,1,0),
--   (23906982,122.81,'2016-03-09 00:23:38',1,NULL,1,0),
--   (24184522,78.36,'2016-03-09 00:30:06',1,NULL,1,0),
--   (23679593,78.36,'2016-03-09 00:34:51',2,null,1,0),
--   (23817591,82.25,'2016-03-09 00:35:01',1,NULL,1,0),
--   (24190183,77.02,'2016-03-09 00:35:19',1,NULL,1,0),
--   (23788430,87.79,'2016-03-09 00:38:42',1,NULL,1,0),
--   (24223148,37.88,'2016-03-09 00:39:01',1,NULL,1,0),
--   (23173716,77.02,'2016-03-09 00:46:11',1,NULL,1,0),
--   (24026138,77.02,'2016-03-09 00:51:07',2,null,1,0),
--   (24143481,31.71,'2016-03-09 00:55:30',1,NULL,1,0),
--   (23623828,30.29,'2016-03-09 00:59:27',1,NULL,1,0),
--   (23918173,32.46,'2016-03-09 01:04:22',1,NULL,1,0),
--   (23082715,32.01,'2016-03-09 01:11:47',1,NULL,1,0),
--   (24152949,37.88,'2016-03-09 01:16:33',1,NULL,1,0),
--   (23872718,82.69,'2016-03-09 01:17:08',2,null,1,0),
--   (24225192,75.58,'2016-03-09 01:18:08',1,NULL,1,0),
--   (24205387,21.64,'2016-03-09 01:32:32',1,NULL,1,0),
--   (24242751,74.43,'2016-03-09 01:33:23',1,NULL,1,0),
--   (22960570,76.48,'2016-03-09 01:37:06',1,NULL,1,0),
--   (24283790,105.34,'2016-03-09 01:46:21',1,NULL,1,0),
--   (24150964,59.53,'2016-03-09 01:54:12',2,null,1,0),
--   (24230552,67.24,'2016-03-09 02:31:31',1,NULL,1,0),
--   (23360453,36.39,'2016-03-09 02:37:33',1,NULL,1,0),
--   (24234029,43.43,'2016-03-09 02:38:42',1,NULL,1,0),
--   (23717438,184.28,'2016-03-09 02:51:13',1,NULL,1,0),
--   (22827855,37.88,'2016-03-09 03:50:41',1,NULL,1,0),
--   (23055693,121.47,'2016-03-09 03:59:31',1,NULL,1,0),
--   (23998985,113.86,'2016-03-09 04:14:42',1,NULL,1,0),
--   (23532067,76.30,'2016-03-09 04:26:59',1,NULL,1,0),
--   (24243869,38.31,'2016-03-09 05:29:40',1,NULL,1,0),
--   (23515026,37.88,'2016-03-09 06:23:52',1,NULL,1,0),
--   (23267191,32.46,'2016-03-09 06:33:40',1,NULL,1,0),
--   (23907935,37.09,'2016-03-09 07:37:52',1,NULL,1,0),
--   (24221031,37.88,'2016-03-09 07:40:21',1,NULL,1,0);

---------------------------------------------
-- Bank Holiday Table
-- https://www.frbservices.org/holidayschedules/
---------------------------------------------
IF EXISTS (SELECT * FROM SYS.TABLES WHERE NAME='BANK_HOLIDAY')
BEGIN
  DROP TABLE OLTP_SYS.dbo.BANK_HOLIDAY;
END;

CREATE TABLE OLTP_SYS.dbo.BANK_HOLIDAY
(
  HOLIDAY_DATE    DATE        NOT NULL,
  HOLIDAY_NAME    VARCHAR(56) NOT NULL,
  CREATION_DATE   DATETIME    NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  MODIFY_DATE     DATETIME    NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  CREATED_BY      VARCHAR(56) NOT NULL DEFAULT SUSER_SNAME(),
  MODIFIED_BY     VARCHAR(56) NOT NULL DEFAULT SUSER_SNAME()
);

ALTER TABLE OLTP_SYS.dbo.BANK_HOLIDAY
ADD CONSTRAINT [PK_HOLIDAY_DATE] PRIMARY KEY CLUSTERED ([HOLIDAY_DATE] ASC);

INSERT INTO OLTP_SYS.dbo.BANK_HOLIDAY (HOLIDAY_DATE, HOLIDAY_NAME) VALUES
  ('2016-01-01', 'New Years Day'),
  ('2016-01-18', 'Martin Luther King, Jr. Day'),
  ('2016-02-15', 'Presidents Day'),
  ('2016-05-30', 'Memorial Day'),
  ('2016-07-04', 'Independence Day'),
  ('2016-09-05', 'Labor Day'),
  ('2016-10-10', 'Columbus Day'),
  ('2016-11-11', 'Veterans Day'),
  ('2016-11-24', 'Thanksgiving Day'),
  ('2016-12-26', 'Christmas Day'),
  ('2017-01-01', 'New Years Day'),
  ('2017-01-16', 'Martin Luther King, Jr. Day'),
  ('2017-02-20', 'Presidents Day'),
  ('2017-05-29', 'Memorial Day'),
  ('2017-07-04', 'Independence Day'),
  ('2017-09-04', 'Labor Day'),
  ('2017-10-09', 'Columbus Day'),
  ('2017-11-23', 'Thanksgiving Day'),
  ('2017-12-25', 'Christmas Day'),
  ('2018-01-01', 'New Years Day'),
  ('2018-01-15', 'Martin Luther King, Jr. Day'),
  ('2018-02-19', 'Presidents Day'),
  ('2018-05-28', 'Memorial Day'),
  ('2018-07-04', 'Independence Day'),
  ('2018-09-03', 'Labor Day'),
  ('2018-10-08', 'Columbus Day'),
  ('2018-11-12', 'Veterans Day'),
  ('2018-11-22', 'Thanksgiving Day'),
  ('2018-12-25', 'Christmas Day');
