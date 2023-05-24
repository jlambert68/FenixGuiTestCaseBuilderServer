package main

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
)

// Load BasicInformation for TestCase to be able to populate the TestCaseExecution
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) prepareSaveFullTestCase(fullTestCaseMessage *fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage) (returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse) {

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

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "Problem when saving to database",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiTestCaseBuilderServerObject.getHighestFenixTestDataProtoFileVersion()),
		}

		return returnMessage
	}

	defer txn.Commit(context.Background())

	// Save the TestCase
	returnMessage, err = fenixGuiTestCaseBuilderServerObject.saveFullTestCase(txn, fullTestCaseMessage)

	return returnMessage
}

// Save TestCase in Execution-queue
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) saveFullTestCase(dbTransaction pgx.Tx, fullTestCaseMessage *fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage) (returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse, err error) {

	usedDBSchema := "FenixBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	nexTestCaseVersion, err := fenixGuiTestCaseBuilderServerObject.getNexTestCaseVersion(fullTestCaseMessage.TestCaseBasicInformation.BasicTestCaseInformation.NonEditableInformation.TestCaseUuid)
	if err != nil {
		if err != nil {

			// Set Error codes to return message
			var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
			var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

			errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
			errorCodes = append(errorCodes, errorCode)

			// Create Return message
			returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Problem when getting next TestCaseVersion from database",
				ErrorCodes:                   errorCodes,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiTestCaseBuilderServerObject.getHighestFenixTestDataProtoFileVersion()),
			}
		}

		return returnMessage, err

	}

	// Set Next TestCaseVersion in TestCase
	fullTestCaseMessage.TestCaseBasicInformation.BasicTestCaseInformation.NonEditableInformation.TestCaseVersion = nexTestCaseVersion

	// Extract column data to be added to data-row
	tempDomainUuid := fullTestCaseMessage.TestCaseBasicInformation.BasicTestCaseInformation.NonEditableInformation.DomainUuid
	tempDomainName := fullTestCaseMessage.TestCaseBasicInformation.BasicTestCaseInformation.NonEditableInformation.DomainName
	tempTestCaseUuid := fullTestCaseMessage.TestCaseBasicInformation.BasicTestCaseInformation.NonEditableInformation.TestCaseUuid
	tempTestCaseName := fullTestCaseMessage.TestCaseBasicInformation.BasicTestCaseInformation.EditableInformation.TestCaseName
	tempTestCaseVersion := fullTestCaseMessage.TestCaseBasicInformation.BasicTestCaseInformation.NonEditableInformation.TestCaseVersion
	tempTestCaseBasicInformationAsJsonb := protojson.Format(fullTestCaseMessage.TestCaseBasicInformation)
	tempTestInstructionsAsJsonb := protojson.Format(fullTestCaseMessage.MatureTestInstructions)
	tempTestInstructionContainersAsJsonb := protojson.Format(fullTestCaseMessage.MatureTestInstructionContainers)

	var dataRowToBeInsertedMultiType []interface{}
	var dataRowsToBeInsertedMultiType [][]interface{}

	// Create Insert Statement for TestCaseExecution that will be put on ExecutionQueue
	// Data to be inserted in the DB-table
	dataRowsToBeInsertedMultiType = nil

	dataRowToBeInsertedMultiType = nil

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempDomainUuid)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempDomainName)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestCaseUuid)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestCaseName)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestCaseVersion)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestCaseBasicInformationAsJsonb)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestInstructionsAsJsonb)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestInstructionContainersAsJsonb)

	dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "INSERT INTO \"" + usedDBSchema + "\".\"TestCases\" "
	sqlToExecute = sqlToExecute + "(\"DomainUuid\", \"DomainName\", \"TestCaseUuid\", \"TestCaseName\", \"TestCaseVersion\", " +
		"\"TestCaseBasicInformationAsJsonb\", \"TestInstructionsAsJsonb\", \"TestInstructionContainersAsJsonb\") "
	sqlToExecute = sqlToExecute + fenixGuiTestCaseBuilderServerObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
	sqlToExecute = sqlToExecute + ";"

	// Execute Query CloudDB
	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "Problem when Loading TestCase Basic Information from database",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiTestCaseBuilderServerObject.getHighestFenixTestDataProtoFileVersion()),
		}
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
	return &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiTestCaseBuilderServerObject.getHighestFenixTestDataProtoFileVersion()),
	}, nil

}

// See https://www.alexedwards.net/blog/using-postgresql-jsonb
// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a myAttrStruct) Value() (driver.Value, error) {

	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *myAttrStruct) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

type myAttrStruct struct {
	fenixTestCaseBuilderServerGrpcApi.BasicTestCaseInformationMessage
}
