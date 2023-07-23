package main

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)

// Load Full TestCase from Database
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) prepareLoadTestCaseHashes(testCaseUuids *[]string) (responseMessage *fenixTestCaseBuilderServerGrpcApi.TestCasesHashResponse) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"id":    "41b572ec-306e-4865-ac38-a738047703cc",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'prepareLoadTestCaseHashes'")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "Problem when Loading from database",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiTestCaseBuilderServerObject.getHighestFenixTestDataProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.TestCasesHashResponse{
			AckNack:         ackNackResponse,
			TestCasesHashes: nil,
		}

		return responseMessage
	}

	defer txn.Commit(context.Background())

	// Load the list with TestCase Uuid's and their Hashes
	var testCasesHashMessage []*fenixTestCaseBuilderServerGrpcApi.TestCasesHashResponse_TestCasesHashMessage
	testCasesHashMessage, err = fenixGuiTestCaseBuilderServerObject.loadTestCasesHash(txn, testCaseUuids)

	// Error when retrieving TestCase
	if err != nil {
		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "Problem when Loading from database",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiTestCaseBuilderServerObject.getHighestFenixTestDataProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.TestCasesHashResponse{
			AckNack:         ackNackResponse,
			TestCasesHashes: nil,
		}

		return responseMessage
	}

	// TestCase
	if testCasesHashMessage == nil {

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "No TestCases couldn't be found in Database",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiTestCaseBuilderServerObject.getHighestFenixTestDataProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.TestCasesHashResponse{
			AckNack:         ackNackResponse,
			TestCasesHashes: nil,
		}

		return responseMessage
	}

	// Create response message
	var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
	ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiTestCaseBuilderServerObject.getHighestFenixTestDataProtoFileVersion()),
	}

	responseMessage = &fenixTestCaseBuilderServerGrpcApi.TestCasesHashResponse{
		AckNack:         ackNackResponse,
		TestCasesHashes: testCasesHashMessage,
	}

	return responseMessage
}

/*
BEGIN;
CREATE TEMPORARY TABLE  tempraryTableName ON COMMIT DROP AS
SELECT TC."TestCaseUuid", MAX(TC."UniqueCounter") "uniqueCounter"
FROM "FenixBuilder"."TestCases" TC
WHERE TC."TestCaseUuid" IN  ('2178a9f7-9c21-4bdf-8a24-d07178b4ab99', 'e41c409e-b58f-40d9-954c-f8f38fa93974')
group by TC."TestCaseUuid";

SELECT TC."TestCaseUuid", TC."TestCaseHash"
FROM "FenixBuilder"."TestCases" TC
WHERE TC."UniqueCounter" IN (SELECT temp."uniqueCounter" FROM tempraryTableName temp);
COMMIT;

*/

// Load Hashes for requested TestCasUuid's
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) loadTestCasesHash(
	dbTransaction pgx.Tx, testCaseUuids *[]string) (
	testCasesHashMessages []*fenixTestCaseBuilderServerGrpcApi.TestCasesHashResponse_TestCasesHashMessage, err error) {

	// Generate unique number for temporary table namne
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 100000
	randomNumber := rand.Intn(max-min+1) + min
	var randomNumberAsString string
	randomNumberAsString = strconv.Itoa(randomNumber)
	var tempraryTableName = "TEMP_TABLE_" + randomNumberAsString

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "CREATE TEMPORARY TABLE " + tempraryTableName + " ON COMMIT DROP AS "
	sqlToExecute = sqlToExecute + "SELECT TC.\"TestCaseUuid\", MAX(TC.\"UniqueCounter\") \"uniqueCounter\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestCases\" TC "
	sqlToExecute = sqlToExecute + "WHERE TC.\"TestCaseUuid\" IN " + ""
	sqlToExecute = sqlToExecute + "GROUP BY TC.\"TestCaseUuid\" "
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "fec062e5-2c22-4a28-ba6b-510670389c9c",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadTestCasesHash'")
	}

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "112ee469-252c-47ca-b0f3-5cd81297705b",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Log response from CloudDB
	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id":                       "dff4d8fd-3c9b-456f-95b4-fe75a0103a5d",
		"comandTag.Insert()":       comandTag.Insert(),
		"comandTag.Delete()":       comandTag.Delete(),
		"comandTag.Select()":       comandTag.Select(),
		"comandTag.Update()":       comandTag.Update(),
		"comandTag.RowsAffected()": comandTag.RowsAffected(),
		"comandTag.String()":       comandTag.String(),
	}).Debug("Return data for SQL executed in database")

	/*SELECT TC."TestCaseUuid", TC."TestCaseHash"
	FROM "FenixBuilder"."TestCases" TC
	WHERE TC."UniqueCounter" IN (SELECT temp."uniqueCounter" FROM tempraryTableName temp);

	*/

	sqlToExecute = ""
	sqlToExecute = sqlToExecute + "SELECT TC.\"TestCaseUuid\", TC.\"TestCaseHash\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestCases\" TC "
	sqlToExecute = sqlToExecute + "WHERE TC.\"TestCaseExecutionUuid\"  IN "
	sqlToExecute = sqlToExecute + "(SELECT temp.\"uniqueCounter\" FROM \"" + tempraryTableName + "\" temp);"
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "663ef807-1b47-4ebc-a374-e54e36162954",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadTestCasesHash'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "b609397e-cd70-4f6c-9b11-7dfe63f41c87",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Extract data from DB result set
	for rows.Next() {

		var testCasesHashMessage fenixTestCaseBuilderServerGrpcApi.TestCasesHashResponse_TestCasesHashMessage

		err = rows.Scan(

			&testCasesHashMessage.TestCaseUuid,
			&testCasesHashMessage.TestCaseHash,
		)

		if err != nil {

			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "764b7cb3-2426-4f50-bf12-0a6a80f0b2ab",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Add Queue-message to slice of messages
		testCasesHashMessages = append(testCasesHashMessages, &testCasesHashMessage)
	}

	return testCasesHashMessages, err

}
