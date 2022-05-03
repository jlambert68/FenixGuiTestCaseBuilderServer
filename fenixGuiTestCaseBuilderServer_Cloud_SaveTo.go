package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

// Prepare to Save Pinned TestInstructions and TestInstructionContainers to CloudDB
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) prepareSaveMerkleHashMerkleTreeAndTestDataRowsToCloudDB(pinnedTestInstructionsAndTestContainersMessage *fenixTestCaseBuilderServerGrpcApi.PinnedTestInstructionsAndTestContainersMessage) (returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())
	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"id":    "306edce0-7a5a-4a0f-992b-5c9b69b0bcc6",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin' for user: ", pinnedTestInstructionsAndTestContainersMessage.UserId)

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "Problem when saving to database",
			ErrorCodes: errorCodes,
		}

		return returnMessage
	}
	defer txn.Commit(context.Background())

	// Save Pinned TestInstructions- and TestInstructionContainer-data
	err = fenixGuiTestCaseBuilderServerObject.savePinnedTestInstructionsAndTestContainersToCloudDB(txn, pinnedTestInstructionsAndTestContainersMessage)
	if err != nil {

		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"id":    "07b91f77-db17-484f-8448-e53375df94ce",
			"error": err,
		}).Error("Couldn't Save Pinned TestInstructions and pre-created TestInstructionContainer to CloudDB for user: ", pinnedTestInstructionsAndTestContainersMessage.UserId)

		// Stop process in and outgoing messages
		// TODO implement stopping gRPC-api
		// fenixGuiTestCaseBuilderServerObject.stateProcessIncomingAndOutgoingMessage = true

		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"id": "348629ad-c358-4043-81ca-ff5f73b579c5",
		}).Error("Stop process for in- and outgoing messages")

		// Rollback any SQL transactions
		txn.Rollback(context.Background())

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "Problem when saving to database",
			ErrorCodes: errorCodes,
		}

		return returnMessage

	}

	return nil
}

// Save Pinned TestInstructions and TestInstructionContainers to CloudDB
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) savePinnedTestInstructionsAndTestContainersToCloudDB(dbTransaction pgx.Tx, pinnedTestInstructionsAndTestContainersMessage *fenixTestCaseBuilderServerGrpcApi.PinnedTestInstructionsAndTestContainersMessage) (err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "9d4e401a-edbf-4a45-bd34-8d3c13eeaffb",
	}).Debug("Entering: savePinnedTestInstructionsAndTestContainersToCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "e0f4ded9-c140-40cf-95a9-c366daa49e07",
		}).Debug("Exiting: savePinnedTestInstructionsAndTestContainersToCloudDB()")
	}()

	// Get current user
	currentUserUuid := pinnedTestInstructionsAndTestContainersMessage.UserId

	// Get a common dateTimeStamp to use
	currentDataTimeStamp := fenixSyncShared.GenerateDatetimeTimeStampForDB()

	var dataRowToBeInsertedMultiType []interface{}
	var dataRowsToBeInsertedMultiType [][]interface{}

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""

	/*
				create table "PinnedTestInstructionsAndPreCreatedTestInstructionContainers"
				(
		    "UserId"     varchar   not null,
		    "PinnedUuid" uuid      not null,
		    "PinnedName" uuid      not null,
		    "PinnedType" int      not null,
		    "TimeStamp"  timestamp not null,
	*/
	// Create Delete Statement for removing users all pinned TestInstructions and TestInstructionsContainers
	sqlToExecute = sqlToExecute + "DELETE FROM \"" + usedDBSchema + "\".\"PinnedTestInstructionsAndPreCreatedTestInstructionContainers\" "
	sqlToExecute = sqlToExecute + "WHERE \"UserId\" = '" + currentUserUuid + "'; "

	// Create Insert Statement for users pinned TestInstructions
	// Data to be inserted in the DB-table
	dataRowsToBeInsertedMultiType = nil

	for _, pinnedTestInstructionMessage := range pinnedTestInstructionsAndTestContainersMessage.PinnedTestInstructionMessages {

		dataRowToBeInsertedMultiType = nil

		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, currentUserUuid)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, pinnedTestInstructionMessage.TestInstructionUuid)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, pinnedTestInstructionMessage.TestInstructionName)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, 1) // 1 = TestInstructionType
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, currentDataTimeStamp)

		dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)
	}

	// Create Insert Statement for users pinned TestInstructionContainers
	// Data to be inserted in the DB-table

	for _, pinnedTestInstructionContainerMessage := range pinnedTestInstructionsAndTestContainersMessage.PinnedTestInstructionContainerMessages {

		dataRowToBeInsertedMultiType = nil

		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, currentUserUuid)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, pinnedTestInstructionContainerMessage.TestInstructionContainerUuid)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, pinnedTestInstructionContainerMessage.TestInstructionContainerName)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, 2) // 2 = TestInstructionContainerType
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, currentDataTimeStamp)

		dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)
	}

	sqlToExecute = sqlToExecute + "INSERT INTO \"" + usedDBSchema + "\".\"PinnedTestInstructionsAndPreCreatedTestInstructionContainers\" "
	sqlToExecute = sqlToExecute + "(\"UserId\", \"PinnedUuid\", \"PinnedName\", \"PinnedType\", \"TimeStamp\") "
	sqlToExecute = sqlToExecute + fenixGuiTestCaseBuilderServerObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
	sqlToExecute = sqlToExecute + ";"

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {
		return err
	}

	// Log response from CloudDB
	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id":                       "dcb110c2-822a-4dde-8bc6-9ebbe9fcbdb0",
		"comandTag.Insert()":       comandTag.Insert(),
		"comandTag.Delete()":       comandTag.Delete(),
		"comandTag.Select()":       comandTag.Select(),
		"comandTag.Update()":       comandTag.Update(),
		"comandTag.RowsAffected()": comandTag.RowsAffected(),
		"comandTag.String()":       comandTag.String(),
	}).Debug("Return data for SQL executed in database")

	// No errors occurred
	return nil

}

// Generates all "VALUES('xxx', 'yyy')..." for insert statements
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) generateSQLInsertValues(testdata [][]interface{}) (sqlInsertValuesString string) {

	sqlInsertValuesString = ""

	// Loop over both rows and values
	for rowCounter, rowValues := range testdata {
		if rowCounter == 0 {
			// Only add 'VALUES' for first row
			sqlInsertValuesString = sqlInsertValuesString + "VALUES("
		} else {
			sqlInsertValuesString = sqlInsertValuesString + ",("
		}

		for valueCounter, value := range rowValues {
			switch valueType := value.(type) {

			case bool:
				sqlInsertValuesString = sqlInsertValuesString + fmt.Sprint(value)

			case int:
				sqlInsertValuesString = sqlInsertValuesString + fmt.Sprint(value)
			case string:

				sqlInsertValuesString = sqlInsertValuesString + "'" + fmt.Sprint(value) + "'"

			default:
				fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
					"id": "53539786-cbb6-418d-8752-c2e337b9e962",
				}).Fatal("Unhandled type, %valueType", valueType)
			}

			// After the last value then add ')'
			if valueCounter == len(rowValues)-1 {
				sqlInsertValuesString = sqlInsertValuesString + ") "
			} else {
				// Not last value, so Add ','
				sqlInsertValuesString = sqlInsertValuesString + ", "
			}

		}

	}

	return sqlInsertValuesString
}

// Generates incoming values in the following form:  "('monkey', 'tiger'. 'fish')"
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) generateSQLINArray(testdata []string) (sqlInsertValuesString string) {

	// Create a list with '' as only element if there are no elements in array
	if len(testdata) == 0 {
		sqlInsertValuesString = "('')"

		return sqlInsertValuesString
	}

	sqlInsertValuesString = "("

	// Loop over both rows and values
	for counter, value := range testdata {

		if counter == 0 {
			// Only used for first row
			sqlInsertValuesString = sqlInsertValuesString + "'" + value + "'"

		} else {

			sqlInsertValuesString = sqlInsertValuesString + ", '" + value + "'"
		}
	}

	sqlInsertValuesString = sqlInsertValuesString + ") "

	return sqlInsertValuesString
}
