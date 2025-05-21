package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"time"
)

func (fenixCloudDBObject *FenixCloudDBObjectStruct) savePublishedSupportedTestCaseMetaDataParametersCommitOrRoleBack(
	dbTransactionReference *pgx.Tx,
	doCommitNotRoleBackReference *bool) {

	dbTransaction := *dbTransactionReference
	doCommitNotRoleBack := *doCommitNotRoleBackReference

	if doCommitNotRoleBack == true {
		dbTransaction.Commit(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "d0a7e400-8ff8-456d-9d01-15e2812aecbf",
		}).Debug("Doing Commit for SQL  in 'savePublishedSupportedTestCaseMetaDataParametersCommitOrRoleBack'")

	} else {
		dbTransaction.Rollback(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "1faa0ad9-299a-47b2-a8bb-961e88b65883",
		}).Info("Doing Rollback for SQL  in 'savePublishedSupportedTestCaseMetaDataParametersCommitOrRoleBack'")

	}
}

// PrepareSavePublishedSupportedTestCaseMetaDataParameters
// Do initial preparations to be able to save all published SupportedTestCaseMetaData Parameters
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareSavePublishedSupportedTestCaseMetaDataParameters(
	domainUuid string,
	supportedMetaDataAsJson string,
	messageSignatureData *fenixTestCaseBuilderServerGrpcApi.MessageSignatureDataMessage,
	reCreatedMessageHashThatWasSigned string) (
	err error) {

	// Begin SQL TransactionConnectorPublish
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "0a5fb682-c274-4bb4-a049-c4d6426bcbd8",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'PrepareSavePublishedSupportedTestCaseMetaDataParameters'")

		return err

	}

	// After all stuff is done, then Commit or Rollback depending on result
	var doCommitNotRoleBack bool

	// Standard is to do a Rollback
	doCommitNotRoleBack = false

	// When leaving then do the actual commit or rollback
	defer fenixCloudDBObject.savePublishedSupportedTestCaseMetaDataParametersCommitOrRoleBack(
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
			"id":    "9d881429-7985-41cc-ba9a-6b19a19d99f7",
			"error": err,
		}).Error("Problem when verifying signature")

		return err
	}

	// Save all published SupportedTestCaseMetaData Parameters
	err = fenixCloudDBObject.savePublishedSupportedTestCaseMetaDataParameters(
		txn,
		domainUuid,
		supportedMetaDataAsJson)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":    "d2bb1db2-3b10-4978-85e9-ebd99b9ef8c2",
			"error": err,
		}).Error("Problem when saving all published SupportedTestCaseMetaData Parameters")

		return err
	}

	doCommitNotRoleBack = true

	return err
}

// Save all published SupportedTestCaseMetaData Parameters
func (fenixCloudDBObject *FenixCloudDBObjectStruct) savePublishedSupportedTestCaseMetaDataParameters(
	dbTransaction pgx.Tx,
	domainUuid string,
	supportedMetaDataAsJson string) (
	err error) {

	// Verify that Domain exists in database
	var domainBaseData *domainBaseDataStruct
	domainBaseData, err = fenixCloudDBObject.verifyDomainExistsInDatabase(
		dbTransaction,
		domainUuid)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "ec13bcca-65eb-4b18-a3f7-698cb3571a0e",
			"domainUuid": domainUuid,
			"error":      err,
		}).Error("Domain does not exist in database or some error occurred when calling database")

		return err
	}

	// Delete old data in database for published SupportedTestCaseMetaData Parameters
	err = fenixCloudDBObject.performDeleteCurrentSupportedTestCaseMetaDataParameters(
		dbTransaction,
		domainUuid)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "33ac892a-37d6-43c6-93cc-d335019044d2",
			"DomainName": domainBaseData.domainName,
			"DomainUUID": domainBaseData.domainUUID,
			"error":      err,
		}).Error("Got some problem when deleting old data for published SupportedTestCaseMetaDataParameters in database")

		return err
	}

	// Save published SupportedTestCaseMetaData Parameters
	err = fenixCloudDBObject.performSaveSupportedTestCaseMetaDataParameters(
		dbTransaction,
		supportedMetaDataAsJson,
		domainBaseData)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "c4b06572-6a12-4edd-a1b2-8155a81c7fce",
			"DomainName": domainBaseData.domainName,
			"DomainUUID": domainBaseData.domainUUID,
			"error":      err,
		}).Error("Couldn't save published SupportedTestCaseMetaData-parameters to CloudDB")

		return err
	}

	return err
}

// Delete old data in database for published SupportedTestCaseMetaData Parameters
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performDeleteCurrentSupportedTestCaseMetaDataParameters(
	dbTransaction pgx.Tx,
	connectorsDomainUUID string) (
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "f35ab1b0-eee8-4702-be97-64e02a55ef94",
	}).Debug("Entering: performDeleteCurrentSupportedTestCaseMetaDataParameters()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "0053358a-abdc-4fa0-924f-6e5ded0c600a",
		}).Debug("Exiting: performDeleteCurrentSupportedTestCaseMetaDataParameters()")
	}()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "DELETE FROM \"FenixBuilder\".\"SupportedTestCaseMetaData\" STCMD "
	sqlToExecute = sqlToExecute + "WHERE STCMD.\"DomainUuid\" = '" + connectorsDomainUUID + "' "
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "a8b07184-b72c-4e92-8b8e-790a7237f2d3",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'performDeleteCurrentSupportedTestCaseMetaDataParameters'")
	}

	// Query DB
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(ctx, sqlToExecute)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "46785732-b5dc-498e-aaa7-ed6c4b68896a",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Log response from CloudDB
	common_config.Logger.WithFields(logrus.Fields{
		"Id":                       "00bc447b-7f1b-4f76-8be5-7094cf4eb7ff",
		"comandTag.Insert()":       comandTag.Insert(),
		"comandTag.Delete()":       comandTag.Delete(),
		"comandTag.Select()":       comandTag.Select(),
		"comandTag.Update()":       comandTag.Update(),
		"comandTag.RowsAffected()": comandTag.RowsAffected(),
		"comandTag.String()":       comandTag.String(),
	}).Debug("Return data for SQL executed in database")

	return err
}

// Do the actual save for published SupportedTestCaseMetaData Parameters
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performSaveSupportedTestCaseMetaDataParameters(
	dbTransaction pgx.Tx,
	supportedMetaDataAsJson string,
	domainBaseData *domainBaseDataStruct) (
	err error) {

	// Data to be inserted in the DB-table
	var dataRowToBeInsertedMultiType []interface{}
	var dataRowsToBeInsertedMultiType [][]interface{}
	dataRowsToBeInsertedMultiType = nil

	// Exist if now users are specified
	if len(supportedMetaDataAsJson) == 0 {

		common_config.Logger.WithFields(logrus.Fields{
			"Id":                      "d16bc9ec-d6c5-4d84-9598-bc0a1258693f",
			"supportedMetaDataAsJson": supportedMetaDataAsJson,
		}).Debug("json must have a value, can't be empty")

		err = errors.New("json must have a value, can't be empty")

		return err
	}

	var tempTimestampToBeUsed string
	tempTimestampToBeUsed = common_config.GenerateDatetimeFromTimeInputForDB(time.Now())

	dataRowToBeInsertedMultiType = nil

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, domainBaseData.domainUUID)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, domainBaseData.domainName)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, supportedMetaDataAsJson)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTimestampToBeUsed)

	dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "INSERT INTO \"FenixBuilder\".\"SupportedTestCaseMetaData\" "
	sqlToExecute = sqlToExecute + "(\"DomainUuid\", \"DomainName\", \"SupportedTestCaseMetaData\", \"UpdateTimeStamp\") "
	sqlToExecute = sqlToExecute + fenixCloudDBObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "cfba1303-eca4-45e1-8903-597bafcce489",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'performSaveSupportedTestCaseMetaDataParameters'")
	}

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "2321d844-b6bb-4938-8a04-da0181d358ce",
			"sqlToExecute": sqlToExecute,
		}).Error("Got some problem when executing SQL within 'performSaveSupportedTestCaseMetaDataParameters'")

		return err
	}

	// Log response from CloudDB
	common_config.Logger.WithFields(logrus.Fields{
		"Id":                       "c222fb53-1671-4338-b9e1-9097d54c8c1e",
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
