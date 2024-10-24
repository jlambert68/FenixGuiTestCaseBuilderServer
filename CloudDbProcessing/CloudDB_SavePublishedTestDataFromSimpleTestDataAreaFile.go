package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

func (fenixCloudDBObject *FenixCloudDBObjectStruct) SavePublishedTestDataFromSimpleTestDataAreaFileCommitOrRoleBack(
	dbTransactionReference *pgx.Tx,
	doCommitNotRoleBackReference *bool) {

	dbTransaction := *dbTransactionReference
	doCommitNotRoleBack := *doCommitNotRoleBackReference

	if doCommitNotRoleBack == true {
		dbTransaction.Commit(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "87ac20f0-d5b1-4d6b-9002-f88877f99aa8",
		}).Debug("Doing Commit for SQL  in 'SavePublishedTestDataFromSimpleTestDataAreaFileCommitOrRoleBack'")

	} else {
		dbTransaction.Rollback(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "5491410b-4d18-4ca4-ac4f-70be1a0b740c",
		}).Info("Doing Rollback for SQL  in 'SavePublishedTestDataFromSimpleTestDataAreaFileCommitOrRoleBack'")

	}
}

// PrepareSavePublishedTestDataFromSimpleTestDataAreaFile
// Do initial preparations to be able to save the TestData from 'simple' TestDataArea-files
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareSavePublishedTestDataFromSimpleTestDataAreaFile(
	domainUuid string,
	templateRepositoriesConnectionParameters []*fenixTestCaseBuilderServerGrpcApi.TestDataFromOneSimpleTestDataAreaFileMessage,
	messageSignatureData *fenixTestCaseBuilderServerGrpcApi.MessageSignatureDataMessage,
	reCreatedMessageHashThatWasSigned string) (
	err error) {

	// Begin SQL TransactionConnectorPublish
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "9a410f9a-0185-4a6f-a1ce-b6410d8a7679",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'PrepareSavePublishedTestDataFromSimpleTestDataAreaFile'")

		return err

	}

	// After all stuff is done, then Commit or Rollback depending on result
	var doCommitNotRoleBack bool

	// Standard is to do a Rollback
	doCommitNotRoleBack = false

	// When leaving then do the actual commit or rollback
	defer fenixCloudDBObject.SavePublishedTestDataFromSimpleTestDataAreaFileCommitOrRoleBack(
		&txn,
		&doCommitNotRoleBack)

	// Verify that the signature was produced by correct private key
	err = fenixCloudDBObject.validateSignedMessage(
		txn,
		domainUuid,
		messageSignatureData,
		reCreatedMessageHashThatWasSigned)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":    "5766f524-5500-43b3-920f-dc7de4fe4848",
			"error": err,
		}).Error("Problem when verifying signature")

		return err
	}

	// Save the TestData from 'simple' TestDataArea-files
	err = fenixCloudDBObject.savePublishedTestDataFromSimpleTestDataAreaFile(
		txn,
		domainUuid,
		templateRepositoriesConnectionParameters)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":    "c0f2bdf5-11f9-4156-aa89-2bc8d417d3d0",
			"error": err,
		}).Error("Problem when saving TestDataFrom 'simple' TestDataArea-files")

		return err
	}

	doCommitNotRoleBack = true

	return err
}

// Save the TestData from 'simple' TestDataArea-files
func (fenixCloudDBObject *FenixCloudDBObjectStruct) savePublishedTestDataFromSimpleTestDataAreaFile(
	dbTransaction pgx.Tx,
	domainUuid string,
	testDataFromSimpleTestDataAreaFileMessage []*fenixTestCaseBuilderServerGrpcApi.TestDataFromOneSimpleTestDataAreaFileMessage) (
	err error) {

	var testAreaUuid string

	// Verify that Domain exists in database
	var domainBaseData *domainBaseDataStruct
	domainBaseData, err = fenixCloudDBObject.verifyDomainExistsInDatabase(
		dbTransaction,
		domainUuid)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "822bf930-5bfa-4d96-8136-d7acf8445a17",
			"domainUuid": domainUuid,
			"error":      err,
		}).Error("Domain does not exist in database or some error occurred when calling database")

		return err
	}

	// Ensure that DomainUUID for all TestDataAreas use the same as the sender Domains UUID
	var errorMessage, errorMessageToAdd string
	for _, testDataArea := range testDataFromSimpleTestDataAreaFileMessage {
		if testDataArea.TestDataDomainUuid != domainUuid {

			errorMessageToAdd = fmt.Sprintf("TestDataArea '%s'('%s') has DomainUuid='%s' but expected '%s'",
				testDataArea.GetTestDataAreaName(),
				testDataArea.GetTestDataAreaUuid(),
				testDataArea.GetTestDataDomainUuid(),
				domainUuid)

			if len(errorMessage) > 0 {
				errorMessage = errorMessage + ", " + errorMessageToAdd
			} else {
				errorMessage = errorMessageToAdd
			}
		}
	}

	// Check if there was any problem
	if len(errorMessage) > 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"id":           "9321ee8b-bc73-4d71-be05-18329e1d56ac",
			"errorMessage": errorMessage,
		}).Warning("Not correct DomainUUID in TestDataFile")

		err = errors.New(errorMessage)

		return err
	}

	// Loop alla TestDataAreas and delete existing and add the new one
	for _, testDataArea := range testDataFromSimpleTestDataAreaFileMessage {

		// Extract TestAreaUuid from incoming message
		testAreaUuid = testDataArea.GetTestDataAreaUuid()

		// Delete old data in database for published Template Repository Connection Parameters
		err = fenixCloudDBObject.performDeleteCurrentTestDataFromSimpleTestDataAreaFile(
			dbTransaction,
			domainUuid,
			testAreaUuid)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"id":         "2d6f5cc1-0ba6-4fd1-a0a1-5b3bbb757a8c",
				"DomainName": domainBaseData.domainName,
				"DomainUUID": domainBaseData.domainUUID,
				"error":      err,
			}).Error("Got some problem when deleting old data for published Template Repository Connection Parameters in database")

			return err
		}

		// Save published Template Repository Connection Parameters
		err = fenixCloudDBObject.performSaveTestDataFromSimpleTestDataAreaFile(
			dbTransaction,
			testDataArea,
			domainBaseData)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"id":         "607e5e8f-430f-4416-b8f8-9beb05c8b0b7",
				"DomainName": domainBaseData.domainName,
				"DomainUUID": domainBaseData.domainUUID,
				"error":      err,
			}).Error("Couldn't save published Template Repository Connection Parameters to CloudDB")

			return err
		}
	}

	return err
}

// Delete old data in database for published Template Repository Connection Parameters
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performDeleteCurrentTestDataFromSimpleTestDataAreaFile(
	dbTransaction pgx.Tx,
	domainUUID string,
	testDataAreaUuid string) (
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "5be68fed-384d-4de9-8a66-d45e33a81bb9",
	}).Debug("Entering: performDeleteCurrentTestDataFromSimpleTestDataAreaFile()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "faa786c9-d3f8-41bf-b4c8-2c9329e99c22",
		}).Debug("Exiting: performDeleteCurrentTestDataFromSimpleTestDataAreaFile()")
	}()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "DELETE FROM \"FenixBuilder\".\"TestDataFromSimpleTestDataAreaFile\" TDSAF "
	sqlToExecute = sqlToExecute + "WHERE TDSAF.\"TestDataDomainUuid\" = '" + domainUUID + "' AND "
	sqlToExecute = sqlToExecute + "TDSAF.\"TestDataAreaUuid\" = '" + testDataAreaUuid + "' "
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "0d3dcc9a-4f96-4c49-872c-eee84219be82",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'performDeleteCurrentTestDataFromSimpleTestDataAreaFile'")
	}

	// Query DB
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(ctx, sqlToExecute)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "1c9ad006-c113-4cab-a3fb-a7804a408971",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Log response from CloudDB
	common_config.Logger.WithFields(logrus.Fields{
		"Id":                       "3abef697-cfd4-4759-8c38-2b74f9b13b5a",
		"comandTag.Insert()":       comandTag.Insert(),
		"comandTag.Delete()":       comandTag.Delete(),
		"comandTag.Select()":       comandTag.Select(),
		"comandTag.Update()":       comandTag.Update(),
		"comandTag.RowsAffected()": comandTag.RowsAffected(),
		"comandTag.String()":       comandTag.String(),
	}).Debug("Return data for SQL executed in database")

	return err
}

// Do the actual save for published TestData from oen TestDataArea
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performSaveTestDataFromSimpleTestDataAreaFile(
	dbTransaction pgx.Tx,
	testDataFromOneSimpleTestDataAreaFileMessage *fenixTestCaseBuilderServerGrpcApi.TestDataFromOneSimpleTestDataAreaFileMessage,
	domainBaseData *domainBaseDataStruct) (
	err error) {

	// Data to be inserted in the DB-table
	var dataRowToBeInsertedMultiType []interface{}
	var dataRowsToBeInsertedMultiType [][]interface{}
	dataRowsToBeInsertedMultiType = nil

	// Exist if now users are specified
	if testDataFromOneSimpleTestDataAreaFileMessage == nil {
		return err
	}

	var tempTimestampToBeUsed string
	tempTimestampToBeUsed = common_config.GenerateDatetimeFromTimeInputForDB(time.Now())

	dataRowToBeInsertedMultiType = nil

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, domainBaseData.domainUUID)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, domainBaseData.domainName)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testDataFromOneSimpleTestDataAreaFileMessage.GetTestDataDomainTemplateName())
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testDataFromOneSimpleTestDataAreaFileMessage.GetTestDataAreaUuid())
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testDataFromOneSimpleTestDataAreaFileMessage.GetTestDataAreaName())
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testDataFromOneSimpleTestDataAreaFileMessage.GetTestDataFileSha256Hash())
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testDataFromOneSimpleTestDataAreaFileMessage.GetImportantDataInFileSha256Hash())
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTimestampToBeUsed)

	tempHeadersForTestDataFromOneSimpleTestDataAreaFileAsJsonb := protojson.Format(testDataFromOneSimpleTestDataAreaFileMessage)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempHeadersForTestDataFromOneSimpleTestDataAreaFileAsJsonb)

	dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "INSERT INTO \"FenixBuilder\".\"TestDataFromSimpleTestDataAreaFile\" "
	sqlToExecute = sqlToExecute + "(\"TestDataDomainUuid\", \"TestDataDomainName\", \"TestDataDomainTemplateName\"," +
		" \"TestDataAreaUuid\", \"TestDataAreaName\", " +
		"\"TestDataFileSha256Hash\", \"ImportantDataInFileSha256Hash\", \"InsertedTimeStamp\", " +
		"\"TestDataFromOneSimpleTestDataAreaFileFullMessage\") "
	sqlToExecute = sqlToExecute + fenixCloudDBObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "e6d81bd5-f00b-4e15-a120-783f690ce66d",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'performSaveTestDataFromSimpleTestDataAreaFile'")
	}

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "fba3f8bf-d43a-417d-83da-68dfa32d3c44",
			"sqlToExecute": sqlToExecute,
		}).Error("Got some problem when executing SQL within 'performSaveTestDataFromSimpleTestDataAreaFile'")

		return err
	}

	// Log response from CloudDB
	common_config.Logger.WithFields(logrus.Fields{
		"Id":                       "825433fb-b219-4b3a-b7d0-859791c1210d",
		"comandTag.Insert()":       comandTag.Insert(),
		"comandTag.Delete()":       comandTag.Delete(),
		"comandTag.Select()":       comandTag.Select(),
		"comandTag.Update()":       comandTag.Update(),
		"comandTag.RowsAffected()": comandTag.RowsAffected(),
		"comandTag.String()":       comandTag.String(),
		"sqlToExecute":             sqlToExecute,
	}).Debug("Return data for SQL executed in database")

	// No errors occurred
	return err
}
