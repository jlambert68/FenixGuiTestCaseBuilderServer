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

func (fenixCloudDBObject *FenixCloudDBObjectStruct) SavePublishedTemplateRepositoryConnectionParametersCommitOrRoleBack(
	dbTransactionReference *pgx.Tx,
	doCommitNotRoleBackReference *bool) {

	dbTransaction := *dbTransactionReference
	doCommitNotRoleBack := *doCommitNotRoleBackReference

	if doCommitNotRoleBack == true {
		dbTransaction.Commit(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "8c325f12-afa1-40cb-8ab8-1d51c8ba07bc",
		}).Debug("Doing Commit for SQL  in 'SavePublishedTemplateRepositoryConnectionParametersCommitOrRoleBack'")

	} else {
		dbTransaction.Rollback(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "750776ab-c11e-4e36-aa5e-a6bec8d212f0",
		}).Info("Doing Rollback for SQL  in 'SavePublishedTemplateRepositoryConnectionParametersCommitOrRoleBack'")

	}
}

// PrepareSavePublishedTemplateRepositoryConnectionParameters
// Do initial preparations to be able to save all published Template Repository Connection Parameters
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareSavePublishedTemplateRepositoryConnectionParameters(
	domainUuid string,
	templateRepositoriesConnectionParameters []*fenixTestCaseBuilderServerGrpcApi.TemplateRepositoryConnectionParameters,
	signedMessageByWorkerServiceAccountMessage *fenixTestCaseBuilderServerGrpcApi.SignedMessageByWorkerServiceAccountMessage) (
	err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "06b34589-a6c7-4fd0-ab31-3ce66e4cae3d",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'PrepareSavePublishedTemplateRepositoryConnectionParameters'")

		return err

	}

	// After all stuff is done, then Commit or Rollback depending on result
	var doCommitNotRoleBack bool

	// Standard is to do a Rollback
	doCommitNotRoleBack = false

	// When leaving then do the actual commit or rollback
	defer fenixCloudDBObject.SavePublishedTemplateRepositoryConnectionParametersCommitOrRoleBack(
		&txn,
		&doCommitNotRoleBack)

	// Save all published Template Repository Connection Parameters
	err = fenixCloudDBObject.SavePublishedTemplateRepositoryConnectionParameters(
		txn,
		domainUuid,
		templateRepositoriesConnectionParameters,
		signedMessageByWorkerServiceAccountMessage)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":    "d2bb1db2-3b10-4978-85e9-ebd99b9ef8c2",
			"error": err,
		}).Error("Problem when saving all published Template Repository Connection Parameters")

		return err
	}

	doCommitNotRoleBack = true

	return err
}

// Save all published Template Repository Connection Parameters
func (fenixCloudDBObject *FenixCloudDBObjectStruct) SavePublishedTemplateRepositoryConnectionParameters(
	dbTransaction pgx.Tx,
	domainUuid string,
	templateRepositoriesConnectionParameters []*fenixTestCaseBuilderServerGrpcApi.TemplateRepositoryConnectionParameters,
	signedMessageByWorkerServiceAccountMessage *fenixTestCaseBuilderServerGrpcApi.SignedMessageByWorkerServiceAccountMessage) (
	err error) {

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

	// Verify Signed message to secure that sending worker is using correct Service Account
	var verificationOfSignatureSucceeded bool
	verificationOfSignatureSucceeded, err = fenixCloudDBObject.verifySignatureFromWorker(
		signedMessageByWorkerServiceAccountMessage,
		domainBaseData)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":  "618180da-8de6-454d-b489-50eb24a7a41e",
			"err": err,
		}).Info("Got some problem when verifying Signature")

		return err
	}

	// The signature couldn't be verified correctly
	if verificationOfSignatureSucceeded == false {
		common_config.Logger.WithFields(logrus.Fields{
			"id": "5624bb59-7ce9-4643-ad62-e36bd3ba319f",
		}).Warning("The correctness of the signature couldn't be verified")

		err = errors.New("the correctness of the signature couldn't be verified")

		return err
	}

	// Delete old data in database for published Template Repository Connection Parameters
	err = fenixCloudDBObject.performDeleteCurrentTemplateRepositoryConnectionParameters(
		dbTransaction,
		domainUuid)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "33ac892a-37d6-43c6-93cc-d335019044d2",
			"DomainName": domainBaseData.domainName,
			"DomainUUID": domainBaseData.domainUUID,
			"error":      err,
		}).Error("Got some problem when deleting old data for published Template Repository Connection Parameters in database")

		return err
	}

	// Save published Template Repository Connection Parameters
	err = fenixCloudDBObject.performSaveTemplateRepositoryConnectionParameters(
		dbTransaction,
		templateRepositoriesConnectionParameters,
		domainBaseData)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "1dfbacc3-9521-46d1-9219-bd63f38d17cc",
			"DomainName": domainBaseData.domainName,
			"DomainUUID": domainBaseData.domainUUID,
			"error":      err,
		}).Error("Couldn't save published Template Repository Connection Parameters to CloudDB")

		return err
	}

	return err
}

// Delete old data in database for published Template Repository Connection Parameters
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performDeleteCurrentTemplateRepositoryConnectionParameters(
	dbTransaction pgx.Tx,
	connectorsDomainUUID string) (
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "a557a438-a6ef-4d25-91df-4e28b4e49942",
	}).Debug("Entering: performDeleteCurrentTemplateRepositoryConnectionParameters()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "09a72d8a-301d-45b4-ab1f-07d2d5026694",
		}).Debug("Exiting: performDeleteCurrentTemplateRepositoryConnectionParameters()")
	}()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "DELETE FROM \"FenixBuilder\".\"TemplateRepositoryConnectionParameters\" TRCP "
	sqlToExecute = sqlToExecute + "WHERE TRCP.\"DomainUuid\" = '" + connectorsDomainUUID + "' "
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "a40dbadf-907d-458b-a160-9a7be4adea77",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'performDeleteCurrentTemplateRepositoryConnectionParameters'")
	}

	// Query DB
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(ctx, sqlToExecute)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "5931baa8-ad97-4df7-83f4-144fdb056f8b",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Log response from CloudDB
	common_config.Logger.WithFields(logrus.Fields{
		"Id":                       "3aac589c-d02b-4025-ba98-99ce847c40a4",
		"comandTag.Insert()":       comandTag.Insert(),
		"comandTag.Delete()":       comandTag.Delete(),
		"comandTag.Select()":       comandTag.Select(),
		"comandTag.Update()":       comandTag.Update(),
		"comandTag.RowsAffected()": comandTag.RowsAffected(),
		"comandTag.String()":       comandTag.String(),
	}).Debug("Return data for SQL executed in database")

	return err
}

// Do the actual save for published Template Repository Connection Parameters
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performSaveTemplateRepositoryConnectionParameters(
	dbTransaction pgx.Tx,
	templateRepositoriesConnectionParameters []*fenixTestCaseBuilderServerGrpcApi.TemplateRepositoryConnectionParameters,
	domainBaseData *domainBaseDataStruct) (
	err error) {

	// Data to be inserted in the DB-table
	var dataRowToBeInsertedMultiType []interface{}
	var dataRowsToBeInsertedMultiType [][]interface{}
	dataRowsToBeInsertedMultiType = nil

	// Exist if now users are specified
	if len(templateRepositoriesConnectionParameters) == 0 {
		return err
	}

	var tempTimestampToBeUsed string
	tempTimestampToBeUsed = common_config.GenerateDatetimeFromTimeInputForDB(time.Now())

	// Loop TemplateRepositoriesConnectionParameters
	for _, templateRepositoryConnectionParameters := range templateRepositoriesConnectionParameters {

		dataRowToBeInsertedMultiType = nil

		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, domainBaseData.domainUUID)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, domainBaseData.domainName)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, templateRepositoryConnectionParameters.GetRepositoryApiUrl())
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, templateRepositoryConnectionParameters.GetRepositoryOwner())
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, templateRepositoryConnectionParameters.GetRepositoryName())
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, templateRepositoryConnectionParameters.GetRepositoryPath())
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, templateRepositoryConnectionParameters.GetGitHubApiKey())
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTimestampToBeUsed)

		dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)

	}

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "INSERT INTO \"FenixBuilder\".\"TemplateRepositoryConnectionParameters\" "
	sqlToExecute = sqlToExecute + "(\"DomainUuid\", \"DomainName\", \"RepositoryApiUrl\", \"RepositoryOwner\", " +
		"\"RepositoryName\", \"RepositoryPath\", \"GitHubApiKey\", \"UpdateTimeStamp\")"
	sqlToExecute = sqlToExecute + fenixCloudDBObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "d5e52603-e976-4abb-a80e-0b7fab2960d3",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'performSaveTemplateRepositoryConnectionParameters'")
	}

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "77520b7b-1edd-4b1d-a6b3-ee9be55dd521",
			"sqlToExecute": sqlToExecute,
		}).Error("Got some problem when executing SQL within 'performSaveTemplateRepositoryConnectionParameters'")

		return err
	}

	// Log response from CloudDB
	common_config.Logger.WithFields(logrus.Fields{
		"Id":                       "66dec9fd-1f25-4162-beba-5e7116918189",
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
