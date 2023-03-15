--name: GET_PAYMENT_RECORDS
IF DATENAME(dw,GETDATE()) = 'Thursday'
    -- This SQL should get all unsent ACH payments *including* retries
    SELECT P.COLLECTION_PAYMENT_ID, I.PROVIDER_ORDER_ID, I.FIRST_NAME, I.MIDDLE_INITIAL, I.LAST_NAME, I.RTN, I.DAN, P.PAYMENT_AMOUNT
         FROM COLLECTION_PAYMENT P
            JOIN COLLECTION_ITEM I ON P.REFUND_TRANSACTION_ID = I.REFUND_TRANSACTION_ID
              WHERE cast(P.CREATION_DATE as DATE) < cast(getDate() as DATE) -- Previous transactions up to midnight
                AND P.PAYMENT_TYPE_ID = 1 -- ACH
                AND I.RTN <> ''  -- added so we don't grab blank rtns
                AND I.DAN <> '0'   -- so we don't grab bad bank accts
                AND P.CUSTOMER_INITIATED_IND = 1 --1 for customer-initiated
                AND (P.ACH_SENT_IND IS NULL OR P.ACH_SENT_IND = 0)
ELSE
    -- This SQL should ignore ACH retries
    SELECT P.COLLECTION_PAYMENT_ID, I.PROVIDER_ORDER_ID, I.FIRST_NAME, I.MIDDLE_INITIAL, I.LAST_NAME, I.RTN, I.DAN, P.PAYMENT_AMOUNT
         FROM COLLECTION_PAYMENT P
            JOIN COLLECTION_ITEM I ON P.REFUND_TRANSACTION_ID = I.REFUND_TRANSACTION_ID
              WHERE cast(P.CREATION_DATE as DATE) < cast(getDate() as DATE) -- Previous transactions up to midnight
                AND P.PAYMENT_TYPE_ID = 1 -- ACH
                AND I.RTN <> ''  -- added so we don't grab blank rtns
                AND I.DAN <> '0'   -- so we don't grab bad bank accts
                AND P.CUSTOMER_INITIATED_IND = 1 -- 1 for customer-initiated
                AND (P.ACH_RETRY_IND IS NULL OR P.ACH_RETRY_IND = 0) -- do not pick up retries
                AND (P.ACH_SENT_IND IS NULL OR P.ACH_SENT_IND = 0);

--name: UPDATE_PAYMENT_RECORD
UPDATE COLLECTION_PAYMENT
    SET ACH_SENT_IND = 1,
        MODIFY_DATE = GETDATE()
        WHERE COLLECTION_PAYMENT_ID = ?;
