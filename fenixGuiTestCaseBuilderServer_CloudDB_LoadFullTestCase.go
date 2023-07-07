package main

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

// Load Full TestCase from Database
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) prepareLoadFullTestCase(testCaseUuidToLoad string) (responseMessage *fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"id":    "f5ccddd6-cf8f-4eed-bfcb-1db8a757fb0b",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'prepareSaveFullTestCase'")

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

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse{
			AckNackResponse:  ackNackResponse,
			DetailedTestCase: nil,
		}

		return responseMessage
	}

	defer txn.Commit(context.Background())

	// Load the TestCase
	var fullTestCaseMessage *fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage
	fullTestCaseMessage, err = fenixGuiTestCaseBuilderServerObject.LoadFullTestCase(txn, testCaseUuidToLoad)

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

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse{
			AckNackResponse:  ackNackResponse,
			DetailedTestCase: nil,
		}

		return responseMessage
	}

	// TestCase
	if fullTestCaseMessage == nil {

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "TestCase couldn't be found in Database",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiTestCaseBuilderServerObject.getHighestFenixTestDataProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse{
			AckNackResponse:  ackNackResponse,
			DetailedTestCase: nil,
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

	responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse{
		AckNackResponse:  ackNackResponse,
		DetailedTestCase: fullTestCaseMessage,
	}

	return responseMessage
}

/*
SELECT TC."TestCaseBasicInformationAsJsonb", TC."TestInstructionsAsJsonb", "TestInstructionContainersAsJsonb"
FROM "FenixBuilder"."TestCases" TC
WHERE TC."TestCaseUuid" = '1f969ca4-e279-431a-b588-491f6f62d41e'
ORDER BY TC."TestCaseVersion" DESC
LIMIT 1;

*/

// Load All Domains and their address information
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) LoadFullTestCase(
	dbTransaction pgx.Tx, testCaseUuidToLoad string) (
	fullTestCaseMessage *fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage, err error) {

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TC.\"TestCaseBasicInformationAsJsonb\", " +
		"TC.\"TestInstructionsAsJsonb\", \"TestInstructionContainersAsJsonb\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestCases\" TC "
	sqlToExecute = sqlToExecute + fmt.Sprintf("WHERE TC.\"TestCaseUuid\" = '%s' ", testCaseUuidToLoad)
	sqlToExecute = sqlToExecute + "ORDER BY TC.\"TestCaseVersion\" DESC "
	sqlToExecute = sqlToExecute + "LIMIT 1; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "01b246fb-effe-4348-9a5c-830604e6daf6",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'LoadFullTestCase'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "784c6f8d-fd77-44e0-9f2b-17e8438ad749",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	var (
		tempTestCaseBasicInformation             fenixTestCaseBuilderServerGrpcApi.TestCaseBasicInformationMessage
		tempMatureTestInstructions               fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionsMessage
		tempMatureTestInstructionContainers      fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionContainersMessage
		tempTestCaseBasicInformationAsString     string
		tempTestInstructionsAsString             string
		tempTestInstructionContainersAsString    string
		tempTestCaseBasicInformationAsByteArray  []byte
		tempTestInstructionsAsByteArray          []byte
		tempTestInstructionContainersAsByteArray []byte
		//tempTestCaseBasicInformationAsJsonb := protojson.Format(fullTestCaseMessage.TestCaseBasicInformation)
		//tempTestInstructionsAsJsonb := protojson.Format(fullTestCaseMessage.MatureTestInstructions)
		//tempTestInstructionContainersAsJsonb := protojson.Format(fullTestCaseMessage.MatureTestInstructionContainers)
	)

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&tempTestCaseBasicInformationAsString,
			&tempTestInstructionsAsString,
			&tempTestInstructionContainersAsString,
		)

		if err != nil {

			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "2a81dfda-4937-4d9e-9827-7191eb7ac7de",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert json-strings into byte-arrays
		tempTestCaseBasicInformationAsByteArray = []byte(tempTestCaseBasicInformationAsString)
		tempTestInstructionsAsByteArray = []byte(tempTestInstructionsAsString)
		tempTestInstructionContainersAsByteArray = []byte(tempTestInstructionContainersAsString)

		// Convert json-byte-arrys into proto-messages
		err = protojson.Unmarshal(tempTestCaseBasicInformationAsByteArray, &tempTestCaseBasicInformation)
		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":    "d315ea2b-8263-4ad8-9b96-d62da4acf35f",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestCaseBasicInformationAsByteArray' into proto-message")

			return nil, err
		}

		err = protojson.Unmarshal(tempTestInstructionsAsByteArray, &tempMatureTestInstructions)
		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":    "441a35b3-5139-4046-8aeb-a986a84827df",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestInstructionsAsByteArray' into proto-message")

			return nil, err
		}

		err = protojson.Unmarshal(tempTestInstructionContainersAsByteArray, &tempMatureTestInstructionContainers)
		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":    "f4738b27-4c49-448b-b49b-f6cf08508f12",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestInstructionContainersAsByteArray' into proto-message")

			return nil, err
		}

		// Add the different parts into full TestCase-message
		fullTestCaseMessage = &fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage{
			TestCaseBasicInformation:        &tempTestCaseBasicInformation,
			MatureTestInstructions:          &tempMatureTestInstructions,
			MatureTestInstructionContainers: &tempMatureTestInstructionContainers,
			MessageHash:                     "",
		}

	}

	return fullTestCaseMessage, err

}
