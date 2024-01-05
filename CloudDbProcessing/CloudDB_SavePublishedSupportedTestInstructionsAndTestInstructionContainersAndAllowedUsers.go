package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"FenixGuiTestCaseBuilderServer/messagesToWorkerServer"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
	"github.com/sirupsen/logrus"
	"time"
)

func (fenixCloudDBObject *FenixCloudDBObjectStruct) saveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersCommitOrRoleBack(
	dbTransactionReference *pgx.Tx,
	doCommitNotRoleBackReference *bool) {

	dbTransaction := *dbTransactionReference
	doCommitNotRoleBack := *doCommitNotRoleBackReference

	if doCommitNotRoleBack == true {
		dbTransaction.Commit(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "675d793a-403c-4626-a70f-bd7ccd747090",
		}).Debug("Doing Commit for SQL  in 'saveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersCommitOrRoleBack'")

	} else {
		dbTransaction.Rollback(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "d0a3d9d6-6ee3-423a-aba5-81ad53be07d3",
		}).Info("Doing Rollback for SQL  in 'saveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersCommitOrRoleBack'")

	}
}

// PrepareSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers
// Do initial preparations to be able to save all supported TestInstructions, TestInstructionContainers and Allowed Users
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct,
	signedMessageByWorkerServiceAccountMessage *fenixTestCaseBuilderServerGrpcApi.SignedMessageByWorkerServiceAccountMessage) (
	err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "d0a3d9d6-6ee3-423a-aba5-81ad53be07d3",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'PrepareSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")

		return err

	}

	// After all stuff is done, then Commit or Rollback depending on result
	var doCommitNotRoleBack bool

	// Standard is to do a Rollback
	doCommitNotRoleBack = false

	// When leaving then do the actual commit or rollback
	defer fenixCloudDBObject.saveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersCommitOrRoleBack(
		&txn,
		&doCommitNotRoleBack)

	// Save  all supported TestInstructions, TestInstructionContainers and Allowed Users
	err = fenixCloudDBObject.saveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
		txn,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage,
		signedMessageByWorkerServiceAccountMessage)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":    "657a20ff-d1b4-4568-a665-6a519547702a",
			"error": err,
		}).Error("Problem when saving supported TestInstructions, TestInstructionContainers and Allowed Users")

		return err
	}

	doCommitNotRoleBack = true

	return err
}

// Save all supported TestInstructions, TestInstructionContainers and Allowed Users
func (fenixCloudDBObject *FenixCloudDBObjectStruct) saveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	dbTransaction pgx.Tx,
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct,
	signedMessageByWorkerServiceAccountMessage *fenixTestCaseBuilderServerGrpcApi.SignedMessageByWorkerServiceAccountMessage) (
	err error) {

	// Verify that Domain exists in database
	var domainBaseData *domainBaseDataStruct
	domainBaseData, err = fenixCloudDBObject.verifyDomainExistsInDatabase(
		dbTransaction,
		string(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID))

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "b248368c-efdf-475c-b2e3-8c4643a11c9d",
			"domainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"error":      err,
		}).Error("Domain does not exist in database or some error occurred when calling database")

		return err
	}

	// Get saved message hash for Domain
	var savedMessageHash string
	savedMessageHash, err = fenixCloudDBObject.prepareLoadSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessageHash(
		string(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID))

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "5986b01f-d584-470a-908e-6f8898fd71e1",
			"ConnectorsDomainUUID": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID,
			"error":                err,
		}).Error("Couldn't get saved Message Hash from CloudDB")

		return err

	}

	// When the saved Message Hash is equal to the incoming Message Hash then nothing is change, which is the base case
	if savedMessageHash == testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.
		TestInstructionsAndTestInstructionsContainersAndUsersMessageHash {

		return nil
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

	// When there is no Message Hash in database then just save the message in the database
	// or when this is a new 'baseline' for the domains supported TestInstructions, TestInstructionContainers and Allowed Users
	if savedMessageHash == "" ||
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.
			ForceNewBaseLineForTestInstructionsAndTestInstructionContainers == true {

		if savedMessageHash == "" {
			common_config.Logger.WithFields(logrus.Fields{
				"id":         "d8c8ef69-49f7-464e-b51f-23b5ca59bca9",
				"domainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
				"domainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			}).Info("No Message Hash found in database, so supported TestInstructions, TestInstructionContainers and Allowed Users will be saved")
		}

		// New forced 'baseline' for the domains supported TestInstructions, TestInstructionContainers and Allowed Users
		if testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.
			ForceNewBaseLineForTestInstructionsAndTestInstructionContainers == true {
			common_config.Logger.WithFields(logrus.Fields{
				"id":         "ca054b92-f093-438a-bfb1-be5438ca3f33",
				"domainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
				"domainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			}).Info("New forced 'baseline' for the domains supported TestInstructions, TestInstructionContainers and Allowed Users")
		}

		// Save supported TestInstructions, TestInstructionContainers and Allowed Users in Database due to New forced 'baseline'
		err = fenixCloudDBObject.performSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
			dbTransaction, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage,
			domainBaseData)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"id":         "056dfe6d-fa17-4e2b-a08b-453e597d033e",
				"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
				"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
				"error":      err,
			}).Error("Couldn't save supported TestInstructions, TestInstructionContainers and Allowed Users to CloudDB")

			return err
		}

		return err
	}

	// Verify changes to TestInstructions, TestInstructionContainers and Allowed Users separately
	err = fenixCloudDBObject.verifyChangesToTestInstructionsAndTestInstructionContainersAndAllowedUsersSeparately(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage,
		domainBaseData)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "9e17b75b-2461-49d1-ba6d-37d47db67e39",
			"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			"error":      err,
		}).Error("Got some problem when Verifying changes to TestInstructions, TestInstructionContainers and Allowed Users separately")

		return err
	}

	return err
}

// Verify that Domain exists in database
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifyDomainExistsInDatabase(
	dbTransaction pgx.Tx,
	domainUUID string) (
	domainBaseData *domainBaseDataStruct,
	err error) {

	domainBaseData, err = fenixCloudDBObject.loadDomainBaseData(dbTransaction, domainUUID)

	return domainBaseData, err
}

// Do the actual save for all supported TestInstructions, TestInstructionContainers and Allowed Users to database
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	dbTransaction pgx.Tx,
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct,
	domainBaseData *domainBaseDataStruct) (
	err error) {

	// Do the actual save for all supported TestInstructions, TestInstructionContainers to database
	err = fenixCloudDBObject.performSaveSupportedTestInstructionsAndTestInstructionContainers(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage,
		domainBaseData)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":  "4aff6581-f81c-4a8e-ad00-7b966a62a3ee",
			"err": err,
		}).Error("Got problems when saving supported TestInstructions, TestInstructionContainers to database")

		return err
	}

	// Do the actual save for Allowed Users to database
	err = fenixCloudDBObject.performSaveSupportedAllowedUsers(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage,
		domainBaseData)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":  "4ed84bb6-3f34-49c3-9b15-a867c5ce724d",
			"err": err,
		}).Error("Got problems when saving supported Allowed Users to Database")

		return err
	}

	return err
}

// Do the actual save for all supported TestInstructions, TestInstructionContainers to database
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performSaveSupportedTestInstructionsAndTestInstructionContainers(
	dbTransaction pgx.Tx,
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct,
	domainBaseData *domainBaseDataStruct) (
	err error) {

	// Create Insert Statement for TestCaseExecution that will be put on ExecutionQueue
	// Data to be inserted in the DB-table
	var dataRowToBeInsertedMultiType []interface{}
	var dataRowsToBeInsertedMultiType [][]interface{}
	dataRowsToBeInsertedMultiType = nil
	dataRowToBeInsertedMultiType = nil

	var tempsupportedtiandticandallowedusersmessageasjsonbAsByteString []byte
	tempsupportedtiandticandallowedusersmessageasjsonbAsByteString, err = json.Marshal(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

	var tempTimestampToBeUsed string
	tempTimestampToBeUsed = common_config.GenerateDatetimeFromTimeInputForDB(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.MessageCreationTimeStamp)

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, string(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID))
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, string(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName))
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.TestInstructionsAndTestInstructionsContainersAndUsersMessageHash)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.TestInstructions.TestInstructionsHash)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.TestInstructionContainers.TestInstructionContainersHash)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.AllowedUsers.AllowedUsersHash)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, string(tempsupportedtiandticandallowedusersmessageasjsonbAsByteString))
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTimestampToBeUsed)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTimestampToBeUsed)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, domainBaseData.bitNumberValue)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, domainBaseData.bitNumberValue)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, domainBaseData.bitNumberValue)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, domainBaseData.bitNumberValue)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, domainBaseData.bitNumberValue)

	dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "INSERT INTO \"FenixBuilder\".\"SupportedTIAndTICAndAllowedUsers\"  "
	sqlToExecute = sqlToExecute + "(\"domainuuid\", \"domainname\", \"messagehash\", \"testinstructionshash\", " +
		"\"testinstructioncontainershash\", \"allowedusershash\", \"supportedtiandticandallowedusersmessageasjsonb\", " +
		"\"updatedtimestamp\", \"lastpublishedtimestamp\", " +
		"canlistandviewtestcaseownedbythisdomain, canbuildandsavetestcaseownedbythisdomain, " +
		"canlistandviewtestcasehavingtiandticfromthisdomain, canlistandviewtestcasehavingtiandticfromthisdomainextended, " +
		"canbuildandsavetestcasehavingtiandticfromthisdomain) "
	sqlToExecute = sqlToExecute + fenixCloudDBObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "3d66b86b-5d5b-45d3-8290-cb2338f2b8bf",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'performSaveSupportedTestInstructionsAndTestInstructionContainers'")
	}

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "c8a6c3da-c83d-42ef-bcf7-f56adc03361a",
			"sqlToExecute": sqlToExecute,
		}).Error("Got some problem when executing SQL within 'performSaveSupportedTestInstructionsAndTestInstructionContainers'")

		return err
	}

	// Log response from CloudDB
	common_config.Logger.WithFields(logrus.Fields{
		"Id":                       "10e10916-108f-48dc-96da-a84c4e7df835",
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

// Do the actual save for Allowed Users to database
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performSaveSupportedAllowedUsers(
	dbTransaction pgx.Tx,
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct,
	domainBaseData *domainBaseDataStruct) (
	err error) {

	// Create Insert Statement for TestCaseExecution that will be put on ExecutionQueue
	// Data to be inserted in the DB-table
	var dataRowToBeInsertedMultiType []interface{}
	var dataRowsToBeInsertedMultiType [][]interface{}
	dataRowsToBeInsertedMultiType = nil

	// Loop Allowed User
	for _, allowedUser := range testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.AllowedUsers.AllowedUsers {

		dataRowToBeInsertedMultiType = nil

		var tempUniqueIdHash string // concat(DomainUUID, UserIdOnComputer, GCPAuthenticatedUser)
		var tempUniqueIdHashValuesSlice []string
		tempUniqueIdHashValuesSlice = []string{
			string(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID),
			allowedUser.UserIdOnComputer,
			allowedUser.GCPAuthenticatedUser}

		// Convert 'CanListAndViewTestCaseOwnedByThisDomain'-bool into value based on 'domainBaseData.bitNumberValue'
		var tempCanListAndViewTestCaseOwnedByThisDomain int64
		if allowedUser.UserAuthorizationRights.CanListAndViewTestCaseOwnedByThisDomain == true {
			tempCanListAndViewTestCaseOwnedByThisDomain = domainBaseData.bitNumberValue
		}

		// Convert 'CanBuildAndSaveTestCaseOwnedByThisDomain'-bool into value based on 'domainBaseData.bitNumberValue'
		var tempCanBuildAndSaveTestCaseOwnedByThisDomain int64
		if allowedUser.UserAuthorizationRights.CanBuildAndSaveTestCaseOwnedByThisDomain == true {
			tempCanBuildAndSaveTestCaseOwnedByThisDomain = domainBaseData.bitNumberValue
		}

		// Convert 'CanListAndViewTestCaseHavingTIandTICFromThisDomain'-bool into value based on 'domainBaseData.bitNumberValue'
		var tempCanListAndViewTestCaseHavingTIandTICFromThisDomain int64
		if allowedUser.UserAuthorizationRights.CanListAndViewTestCaseHavingTIandTICFromThisDomain == true {
			tempCanListAndViewTestCaseHavingTIandTICFromThisDomain = domainBaseData.bitNumberValue
		}

		// Convert 'CanListAndViewTestCaseHavingTIandTICFromThisDomainExtended'-bool into value based on 'domainBaseData.bitNumberValue'
		var tempCanListAndViewTestCaseHavingTIandTICFromThisDomainExtendedn int64
		if allowedUser.UserAuthorizationRights.CanListAndViewTestCaseHavingTIandTICFromThisDomainExtended == true {
			tempCanListAndViewTestCaseHavingTIandTICFromThisDomainExtendedn = domainBaseData.bitNumberValue
		}

		// Convert 'CanBuildAndSaveTestCaseHavingTIandTICFromThisDomain'-bool into value based on 'domainBaseData.bitNumberValue'
		var tempCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain int64
		if allowedUser.UserAuthorizationRights.CanBuildAndSaveTestCaseHavingTIandTICFromThisDomain == true {
			tempCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain = domainBaseData.bitNumberValue
		}

		// Hash slice
		tempUniqueIdHash = fenixSyncShared.HashValues(tempUniqueIdHashValuesSlice, true)

		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempUniqueIdHash)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, string(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID))
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, string(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName))
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, allowedUser.UserIdOnComputer)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, allowedUser.GCPAuthenticatedUser)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, allowedUser.UserEmail)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, allowedUser.UserFirstName)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, allowedUser.UserLastName)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempCanListAndViewTestCaseOwnedByThisDomain)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempCanBuildAndSaveTestCaseOwnedByThisDomain)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempCanListAndViewTestCaseHavingTIandTICFromThisDomain)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempCanListAndViewTestCaseHavingTIandTICFromThisDomainExtendedn)
		dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain)

		dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)

	}

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "INSERT INTO \"FenixDomainAdministration\".\"allowedusers\" "
	sqlToExecute = sqlToExecute + "(\"uniqueidhash\", \"domainuuid\", \"domainname\", \"useridoncomputer\", " +
		"\"gcpauthenticateduser\", \"useremail\", \"userfirstname\", \"userlastname\" ," +
		"canlistandviewtestcaseownedbythisdomain, canbuildandsavetestcaseownedbythisdomain, " +
		"canlistandviewtestcasehavingtiandticfromthisdomain, canlistandviewtestcasehavingtiandticfromthisdomainextended, " +
		"canbuildandsavetestcasehavingtiandticfromthisdomain)"
	sqlToExecute = sqlToExecute + fenixCloudDBObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "377a3597-ab10-41bb-a2b0-b4c0e2722b58",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'performSaveSupportedAllowedUsers'")
	}

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "fd94e10e-3f4e-4dc8-a634-adbf94ac6c35",
			"sqlToExecute": sqlToExecute,
		}).Error("Got some problem when executing SQL within 'performSaveSupportedAllowedUsers'")

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

// Verify changes to TestInstructions, TestInstructionContainers and Allowed Users separately
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifyChangesToTestInstructionsAndTestInstructionContainersAndAllowedUsersSeparately(
	dbTransaction pgx.Tx,
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct,
	domainBaseData *domainBaseDataStruct) (
	err error) {

	// Load saved message for supported TestInstructions, TestInstructionContainers and Allowed Users from database
	// To be used when comparing TestInstructions, TestInstructionContainers and Allowed Users

	var supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage *supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessageStruct
	supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage, err = fenixCloudDBObject.
		loadDomainSpecificPublishedSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage(
			dbTransaction,
			string(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID))

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "dc62c173-d641-46c1-ba48-437885ac6b02",
			"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			"error":      err,
		}).Error("Got some problem when loading save message for supportedTestInstructions, TestInstructionContainers and Allowed Users from database")

		return err

	}

	// Convert loaded message from database into message to be use below
	var testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct
	testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB, err = fenixCloudDBObject.
		convertIntoTestInstructionsAndTestInstructionContainersFromGrpcBuilderMessage(
			supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage)

	// Verify changes in TestInstructions
	var correctNewChangesFoundInTestInstructions bool
	correctNewChangesFoundInTestInstructions, err = fenixCloudDBObject.verifyChangesToTestInstructions(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.TestInstructions,
		testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "1651733e-b147-4b9c-952a-1aba60618e88",
			"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			"error":      err,
		}).Error("Got some problem when verifying changes to TestInstructions")

		return err
	}

	// Verify changes in TestInstructionContainers
	var correctNewChangesFoundInTestInstructionContainers bool
	correctNewChangesFoundInTestInstructionContainers, err = fenixCloudDBObject.verifyChangesToTestInstructionContainers(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.TestInstructionContainers,
		testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "1651733e-b147-4b9c-952a-1aba60618e88",
			"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			"error":      err,
		}).Error("Got some problem when verifying changes to TestInstructionContainers")

		return err
	}

	// Verify changes in Allowed Users
	var correctNewChangesFoundInAllowedUsers bool
	correctNewChangesFoundInAllowedUsers, err = fenixCloudDBObject.verifyChangesToAllowedUsers(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.AllowedUsers,
		testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "1651733e-b147-4b9c-952a-1aba60618e88",
			"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			"error":      err,
		}).Error("Got some problem when verifying changes to Allowed Users")

		return err
	}

	// Found correct changes, so update supported TestInstructions, TestInstructionContainers and Allowed Users in database
	// Can only be done when ForceNewBaseLineForTestInstructionsAndTestInstructionContainers==true
	if (correctNewChangesFoundInTestInstructions == true || correctNewChangesFoundInTestInstructionContainers ||
		correctNewChangesFoundInAllowedUsers) && testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.
		ForceNewBaseLineForTestInstructionsAndTestInstructionContainers == false {
		common_config.Logger.WithFields(logrus.Fields{
			"id": "83536474-a037-4b67-85fc-2818f9181e38",
			"correctNewChangesFoundInTestInstructions":          correctNewChangesFoundInTestInstructions,
			"correctNewChangesFoundInTestInstructionContainers": correctNewChangesFoundInTestInstructionContainers,
			"correctNewChangesFoundInAllowedUsers":              correctNewChangesFoundInAllowedUsers,
		}).Info("Found correct changes, so update supported TestInstructions, TestInstructionContainers and Allowed Users in database")
	}

	// Delete old data in database for Supported TestInstructions, TestInstructionContainers And Allowed Users
	err = fenixCloudDBObject.performDeleteCurrentSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "d65868ff-251f-47ec-a3c2-c8cdc3fac90c",
			"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			"error":      err,
		}).Error("Got some problem when deleting old data for supported TestInstructions, TestInstructionContainers and Allowed Users in database")

		return err
	}

	// Save new message with supported TestInstructions, TestInstructionContainers and Allowed Users in database
	err = fenixCloudDBObject.performSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage,
		domainBaseData)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "3df8e2a0-efb2-40a4-9fa2-58eb23ffa911",
			"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			"error":      err,
		}).Error("Got some problem when saving supported TestInstructions, TestInstructionContainers and Allowed Users in database")

		return err
	}

	return err
}

// Verify changes in TestInstructions
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifyChangesToTestInstructions(
	dbTransaction pgx.Tx,
	testInstructionsMessage *TestInstructionAndTestInstuctionContainerTypes.TestInstructionsStruct,
	testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	correctNewChangesFoundInTestInstructions bool, err error) {

	// No TestInstructions
	if testInstructionsMessage == nil {

		return correctNewChangesFoundInTestInstructions, err
	}

	// There aren't any changes then just return
	if testInstructionsMessage.TestInstructionsHash !=
		testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB.TestInstructions.TestInstructionsHash {

		return correctNewChangesFoundInTestInstructions, err
	}

	// Now we know that there are changes, find which TestInstruction that are changed, added or removed
	// Loop existing TestInstruction that were stored in database
	for tempTestInstructionUUIDFromDB, tempTestInstructionFromDB := range testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB.TestInstructions.TestInstructionsMap {
		// Check if it exists in new message
		publishedTestInstructionFromConnector, existInMap := testInstructionsMessage.TestInstructionsMap[tempTestInstructionUUIDFromDB]

		// If it doesn't exit in new message then just continue
		if existInMap == false {
			continue
		}

		// If it has the same hash then just continue
		if publishedTestInstructionFromConnector.TestInstructionVersionsHash == tempTestInstructionFromDB.TestInstructionVersionsHash {
			continue
		}

		// If the hash is different then process the TestInstructionVersions
		if publishedTestInstructionFromConnector.TestInstructionVersionsHash != tempTestInstructionFromDB.TestInstructionVersionsHash {

			// Loop existing TestInstructionVersion that were stored in database
			for versionCounter, tempTestInstructionVersionFromDB := range tempTestInstructionFromDB.TestInstructionVersions {

				// If the version hash is different then something is wrong
				if tempTestInstructionVersionFromDB.TestInstructionInstanceVersionHash !=
					publishedTestInstructionFromConnector.TestInstructionVersions[versionCounter].
						TestInstructionInstanceVersionHash {

					var tempTestInstructionUUID string
					var tempMajorVersionNumber int
					var tempMinorVersionNumber int

					tempTestInstructionUUID = string(publishedTestInstructionFromConnector.
						TestInstructionVersions[versionCounter].TestInstructionInstance.TestInstruction.TestInstructionUUID)
					tempMajorVersionNumber = publishedTestInstructionFromConnector.
						TestInstructionVersions[versionCounter].TestInstructionInstance.TestInstruction.MajorVersionNumber
					tempMinorVersionNumber = publishedTestInstructionFromConnector.
						TestInstructionVersions[versionCounter].TestInstructionInstance.TestInstruction.MinorVersionNumber

					common_config.Logger.WithFields(logrus.Fields{
						"id":                  "8c91a76d-887a-477d-b5ea-168b0d0c6061",
						"TestInstructionUUID": tempTestInstructionUUID,
						"MajorVersionNumber":  tempMajorVersionNumber,
						"MinorVersionNumber":  tempMinorVersionNumber,
						"versionCounter":      versionCounter,
					}).Error("New TestInstructionVersion is not the same as existing, which was expected")

					err = errors.New(fmt.Sprintf("New TestInstructionVersion has not the same hash as existing, "+
						"which was expected. "+
						"TestInstructionUUID=%s, MajorVersionNumber=%d, MinorVersionNumber=%d, SlicePosition=%d ",
						tempTestInstructionUUID,
						tempMajorVersionNumber,
						tempMinorVersionNumber,
						versionCounter))

					return false, err

				}
			}
			// Verify if there are more TestInstructionVersions in new published TestInstructionVersions
			if len(publishedTestInstructionFromConnector.TestInstructionVersions) > len(tempTestInstructionFromDB.TestInstructionVersions) {

				// There are more TestInstructionVersions in new published TestInstructionVersions so accept that
				correctNewChangesFoundInTestInstructions = true

				return correctNewChangesFoundInTestInstructions, err

			} else {
				// This should happen because the versions-Hash can't be the same when there are more versions in new TestInstruction

				var tempTestInstructionUUID string
				tempTestInstructionUUID = string(publishedTestInstructionFromConnector.
					TestInstructionVersions[0].TestInstructionInstance.TestInstruction.TestInstructionUUID)

				common_config.Logger.WithFields(logrus.Fields{
					"id":                  "7dc54367-7290-4641-a43a-dfb4ca30a728",
					"TestInstructionUUID": tempTestInstructionUUID,
				}).Error("This should happen because the versions-Hash can't be the same when there are more versions in new TestInstruction")

				err = errors.New(fmt.Sprintf("This should happen because the versions-Hash can't be the same when "+
					"there are more versions in new TestInstruction. "+
					"TestInstructionUUID=%s",
					tempTestInstructionUUID))

				return false, err

			}
		}
	}

	return correctNewChangesFoundInTestInstructions, err
}

// Verify changes in TestInstructionContainers
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifyChangesToTestInstructionContainers(
	dbTransaction pgx.Tx,
	testInstructionContainersMessage *TestInstructionAndTestInstuctionContainerTypes.TestInstructionContainersStruct,
	testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	correctNewChangesFoundInTestInstructionContainers bool, err error) {

	// No TestInstructionContainers
	if testInstructionContainersMessage == nil {

		return correctNewChangesFoundInTestInstructionContainers, err
	}

	// There aren't any changes then just return
	if testInstructionContainersMessage.TestInstructionContainersHash !=
		testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB.TestInstructionContainers.TestInstructionContainersHash {

		return correctNewChangesFoundInTestInstructionContainers, err
	}

	// Now we know that there are changes, find which TestInstruction that are changed, added or removed
	// Loop existing TestInstruction that were stored in database
	for tempTestInstructionContainerUUIDFromDB, tempTestInstructionContainerFromDB := range testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB.TestInstructionContainers.TestInstructionContainersMap {

		// Check if it exists in new message
		publishedTestInstructionContainerFromConnector, existInMap := testInstructionContainersMessage.
			TestInstructionContainersMap[tempTestInstructionContainerUUIDFromDB]

		// If it doesn't exit in new message then just continue
		if existInMap == false {
			continue
		}

		// If it has the same hash then just continue
		if publishedTestInstructionContainerFromConnector.TestInstructionContainerVersionsHash ==
			tempTestInstructionContainerFromDB.TestInstructionContainerVersionsHash {
			continue
		}

		// If the hash is different then process the TestInstructionVersions
		if publishedTestInstructionContainerFromConnector.TestInstructionContainerVersionsHash !=
			tempTestInstructionContainerFromDB.TestInstructionContainerVersionsHash {

			// Loop existing TestInstructionContainerVersion that were stored in database
			for versionCounter, tempTestInstructionContainerVersionFromDB := range tempTestInstructionContainerFromDB.TestInstructionContainerVersions {

				// If the version hash is different then something is wrong
				if tempTestInstructionContainerVersionFromDB.TestInstructionContainerInstanceHash !=
					publishedTestInstructionContainerFromConnector.TestInstructionContainerVersions[versionCounter].
						TestInstructionContainerInstanceHash {

					var tempTestInstructionContainerUUID string
					var tempMajorVersionNumber int
					var tempMinorVersionNumber int

					tempTestInstructionContainerUUID = string(publishedTestInstructionContainerFromConnector.
						TestInstructionContainerVersions[versionCounter].TestInstructionContainerInstance.
						TestInstructionContainer.TestInstructionContainerUUID)
					tempMajorVersionNumber = publishedTestInstructionContainerFromConnector.
						TestInstructionContainerVersions[versionCounter].TestInstructionContainerInstance.
						TestInstructionContainer.MajorVersionNumber
					tempMinorVersionNumber = publishedTestInstructionContainerFromConnector.
						TestInstructionContainerVersions[versionCounter].TestInstructionContainerInstance.
						TestInstructionContainer.MinorVersionNumber

					common_config.Logger.WithFields(logrus.Fields{
						"id":                           "160988ea-7f7c-4138-8190-e0f6d87e4a75",
						"TestInstructionContainerUUID": tempTestInstructionContainerUUID,
						"MajorVersionNumber":           tempMajorVersionNumber,
						"MinorVersionNumber":           tempMinorVersionNumber,
						"versionCounter":               versionCounter,
					}).Error("New TestInstructionContainerVersion is not the same as existing, which was expected")

					err = errors.New(fmt.Sprintf("New TestInstructionContainerVersion has not the same hash as existing, "+
						"which was expected. "+
						"TestInstructionUUID=%s, MajorVersionNumber=%d, MinorVersionNumber=%d, SlicePosition=%d ",
						tempTestInstructionContainerUUID,
						tempMajorVersionNumber,
						tempMinorVersionNumber,
						versionCounter))

					return false, err

				}
			}
			// Verify if there are more TestInstructionContainerVersions in new published TestInstructionVersions
			if len(publishedTestInstructionContainerFromConnector.TestInstructionContainerVersions) >
				len(tempTestInstructionContainerFromDB.TestInstructionContainerVersions) {

				// There are more TestInstructionVersions in new published TestInstructionVersions so accept that
				correctNewChangesFoundInTestInstructionContainers = true

				return correctNewChangesFoundInTestInstructionContainers, err

			} else {
				// This should happen because the versions-Hash can't be the same when there are more versions in new TestInstructionContainer

				var tempTestInstructionContainerVersionUUID string
				tempTestInstructionContainerVersionUUID = string(publishedTestInstructionContainerFromConnector.
					TestInstructionContainerVersions[0].TestInstructionContainerInstance.TestInstructionContainer.TestInstructionContainerUUID)

				common_config.Logger.WithFields(logrus.Fields{
					"id":                  "ea5745aa-e943-44b3-b7d8-ea753535968c",
					"TestInstructionUUID": tempTestInstructionContainerVersionUUID,
				}).Error("This should happen because the versions-Hash can't be the same when there are more versions in new TestInstructionContainer")

				err = errors.New(fmt.Sprintf("This should happen because the versions-Hash can't be the same when "+
					"there are more versions in new TestInstructionContainer. "+
					"TestInstructionUUID=%s",
					tempTestInstructionContainerVersionUUID))

				return false, err

			}
		}
	}

	return correctNewChangesFoundInTestInstructionContainers, err
}

// Verify changes in Allowed Users
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifyChangesToAllowedUsers(
	dbTransaction pgx.Tx,
	allowedUsers *TestInstructionAndTestInstuctionContainerTypes.AllowedUsersStruct,
	testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	correctNewChangesFoundInAllowedUsers bool, err error) {

	// There aren't any changes then just return
	if allowedUsers.AllowedUsersHash !=
		testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB.AllowedUsers.AllowedUsersHash {

		return correctNewChangesFoundInAllowedUsers, err
	}

	// Hash is changed so accept the new User Setup
	correctNewChangesFoundInAllowedUsers = true

	return correctNewChangesFoundInAllowedUsers, err
}

type domainBaseDataStruct struct {
	domainUUID          string
	domainName          string
	workerAddressToDial string
	bitNumberValue      int64
}

// When row is found the Domain exists and is allowed to use Fenix
// Functions also returns Service Account that should be used by calling Worker
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadDomainBaseData(
	dbTransaction pgx.Tx,
	domainUUID string) (
	domainBaseData *domainBaseDataStruct,
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "13b68d1e-1ecd-4bc6-8b72-cdf50391942c",
	}).Debug("Entering: loadDomainBaseData()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "57e50f73-3524-4cfd-8c43-dc52324a0140",
		}).Debug("Exiting: loadDomainBaseData()")
	}()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT fdad.domain_uuid, fdad.domain_name, fdad.workeraddress, dbpe.bitnumbervalue "
	sqlToExecute = sqlToExecute + "FROM \"FenixDomainAdministration\".\"domains\" fdad, " +
		"\"FenixDomainAdministration\".\"domainbitpositionenum\" dbpe "
	sqlToExecute = sqlToExecute + "WHERE fdad.activated = true AND fdad.deleted = false AND "
	sqlToExecute = sqlToExecute + "fdad.domain_uuid = '" + domainUUID + "' AND "
	sqlToExecute = sqlToExecute + "fdad.bitnumbername = dbpe.bitnumbername "
	sqlToExecute = sqlToExecute + ";"

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "940b1066-507f-4ff4-bec7-51fdc33698a1",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return nil, err
	}

	var rowsCounter int

	// Extract data from DB result set
	for rows.Next() {

		// Temp parameters
		var tempDomainUUID string
		var tempDomainName string
		var tempWorkerAddressToDial string
		var tempBitNumberValue int64

		err = rows.Scan(
			&tempDomainUUID,
			&tempDomainName,
			&tempWorkerAddressToDial,
			&tempBitNumberValue,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "33359dae-cda8-45f3-a57e-8dab751be154",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		domainBaseData = &domainBaseDataStruct{
			domainUUID:          tempDomainUUID,
			domainName:          tempDomainName,
			workerAddressToDial: tempWorkerAddressToDial,
			bitNumberValue:      tempBitNumberValue,
		}

		// Add to row counter; Max = 1
		rowsCounter = rowsCounter + 1

	}

	// Check how many rows that were found
	if rowsCounter > 1 {
		// Shouldn't happen

		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "427c2341-4d82-45a7-aeaf-dfbc6352e307",
			"rowsCounter":  rowsCounter,
			"sqlToExecute": sqlToExecute,
		}).Error("More than 1 row was found in database")

		newErrorMessage := errors.New("More than 1 row was found in database for domain=" + domainUUID)

		return nil, newErrorMessage

	} else if rowsCounter == 0 {
		// Domain is not allowed to use Fenix

		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "3b319f30-1ddd-4cab-b043-f86899d9df89",
			"rowsCounter":  rowsCounter,
			"sqlToExecute": sqlToExecute,
			"domainUUID":   domainUUID,
		}).Warning("Domain is not allowed to use Fenix")

		newErrorMessage := errors.New("Domain is not allowed to use Fenix; DomainUUID=" + domainUUID)

		return nil, newErrorMessage

	}

	return domainBaseData, err
}

// Do the signature verification of signature received from Worker
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifySignatureFromWorker(
	signedMessageByWorkerServiceAccountMessage *fenixTestCaseBuilderServerGrpcApi.SignedMessageByWorkerServiceAccountMessage,
	domainBaseData *domainBaseDataStruct) (
	verificationOfSignatureSucceeded bool,
	err error) {

	// Set up temporary variable used when calling Worker over gRPC
	var tempMessagesToWorkerServerObject *messagesToWorkerServer.MessagesToWorkerServerObjectStruct
	tempMessagesToWorkerServerObject = &messagesToWorkerServer.MessagesToWorkerServerObjectStruct{Logger: common_config.Logger}

	// Call Worker over gRPC
	var signMessageResponse *fenixExecutionWorkerGrpcApi.SignMessageResponse
	signMessageResponse, err = tempMessagesToWorkerServerObject.SendBuilderServerAskWorkerToSignMessage(
		signedMessageByWorkerServiceAccountMessage.GetMessageToBeSigned(),
		domainBaseData.workerAddressToDial)

	// Got some problem when doing gRPC-call to WorkerServer
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":             "63858aac-491f-4162-9cbb-ed9d4b1c5ba6",
			"error":          err,
			"domainBaseData": domainBaseData,
		}).Error("Got a problem when calling WorkerServer over gRPC to verify signature")

		return false, err
	}

	// Got some problem when doing gRPC-call to WorkerServer
	if signMessageResponse.GetAckNackResponse().GetAckNack() == false {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                  "09580908-dfbf-4f86-adb5-ca64e9a610b8",
			"error":               err,
			"domainBaseData":      domainBaseData,
			"signMessageResponse": signMessageResponse,
		}).Error("Got a problem when calling WorkerServer over gRPC to verify signature")

		var newError error
		newError = errors.New(signMessageResponse.GetAckNackResponse().GetComments())

		return false, newError
	}

	// Verify recreated signature with signature produced by Worker when sending published TI, TIC and Allowed Users
	if signMessageResponse.GetSignedMessageByWorkerServiceAccount().
		GetHashOfSignature() != signedMessageByWorkerServiceAccountMessage.HashOfSignature {
		return false, err
	}

	// Verify recreated KeyId with KeyId produced by Worker
	if signMessageResponse.GetSignedMessageByWorkerServiceAccount().
		GetHashedKeyId() != signedMessageByWorkerServiceAccountMessage.HashedKeyId {
		return false, err
	}

	// Success in signature verification
	return true, err
}

// Delete old data in database for Supported TestInstructions, TestInstructionContainers And Allowed Users
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performDeleteCurrentSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	dbTransaction pgx.Tx,
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	err error) {

	// Do the actual Delete for all supported TestInstructions, TestInstructionContainers to database
	err = fenixCloudDBObject.performDeleteCurrentSupportedTestInstructionsAndTestInstructionContainers(
		dbTransaction,
		string(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID))

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":                   "f85b3257-071b-4dd2-b371-95d1cb95e89b",
			"err":                  err,
			"connectorsDomainUUID": string(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID),
		}).Error("Got problems when deleting supported TestInstructions, TestInstructionContainers from database")

		return err
	}

	// Do the actual Delete for Allowed Users to database
	err = fenixCloudDBObject.performDeleteCurrentAllowedUsers(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":  "59e80faa-bcb6-4d33-b54a-7a4e854dba96",
			"err": err,
		}).Error("Got problems when deleting supported Allowed Users from Database")

		return err
	}

	return err
}

// Delete old data in database for Supported TestInstructions, TestInstructionContainers
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performDeleteCurrentSupportedTestInstructionsAndTestInstructionContainers(
	dbTransaction pgx.Tx,
	connectorsDomainUUID string) (
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "24827f66-869a-4a20-9e62-e9a6ae85a609",
	}).Debug("Entering: performDeleteCurrentSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "90b01dc1-3eca-4b96-bdcb-c8e9f44ca054",
		}).Debug("Exiting: performDeleteCurrentSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers()")
	}()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "DELETE FROM \"FenixBuilder\".\"SupportedTIAndTICAndAllowedUsers\" STITICAU "
	sqlToExecute = sqlToExecute + "WHERE STITICAU.\"domainuuid\" = '" + connectorsDomainUUID + "' "
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "dc1801a8-df60-463d-b2be-40ed1c05f018",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'deleteTestInstructionMessagesReceivedByWrongInstanceFromDatabaseInCloudDB'")
	}

	// Query DB
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(ctx, sqlToExecute)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "7d047414-29b9-4280-aef2-481598486932",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Log response from CloudDB
	common_config.Logger.WithFields(logrus.Fields{
		"Id":                       "266ceec9-27b1-481b-914d-e127c5ca3f0f",
		"comandTag.Insert()":       comandTag.Insert(),
		"comandTag.Delete()":       comandTag.Delete(),
		"comandTag.Select()":       comandTag.Select(),
		"comandTag.Update()":       comandTag.Update(),
		"comandTag.RowsAffected()": comandTag.RowsAffected(),
		"comandTag.String()":       comandTag.String(),
	}).Debug("Return data for SQL executed in database")

	return err
}

// Delete old data in database for Supported TestInstructions, TestInstructionContainers
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performDeleteCurrentAllowedUsers(
	dbTransaction pgx.Tx,
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "24827f66-869a-4a20-9e62-e9a6ae85a609",
	}).Debug("Entering: performDeleteCurrentAllowedUsers()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "34ce838b-88fc-404a-8ab0-8ee09c313b80",
		}).Debug("Exiting: performDeleteCurrentAllowedUsers()")
	}()

	// Loop Allowed User
	var tempUniqueIdHashesSlice []string
	for _, allowedUser := range testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.AllowedUsers.AllowedUsers {

		var tempUniqueIdHash string // concat(DomainUUID, UserIdOnComputer, GCPAuthenticatedUser)
		var tempUniqueIdHashValuesSlice []string
		tempUniqueIdHashValuesSlice = []string{
			string(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID),
			allowedUser.UserIdOnComputer,
			allowedUser.GCPAuthenticatedUser}

		// Hash slice
		tempUniqueIdHash = fenixSyncShared.HashValues(tempUniqueIdHashValuesSlice, true)

		// Add to slice of UniqueIdHashes
		tempUniqueIdHashesSlice = append(tempUniqueIdHashesSlice, tempUniqueIdHash)
	}

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "DELETE FROM \"FenixDomainAdministration\".\"allowedusers\" au "
	sqlToExecute = sqlToExecute + "WHERE au.\"uniqueidhash\" IN " + fenixCloudDBObject.generateSQLINArray(tempUniqueIdHashesSlice)
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "6c90f1a5-cabf-4315-ad98-c7c4f5cd4e33",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'performDeleteCurrentAllowedUsers'")
	}

	// Query DB
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(ctx, sqlToExecute)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "0caed30f-502a-4b46-8d2f-4ba52578bbf1",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Log response from CloudDB
	common_config.Logger.WithFields(logrus.Fields{
		"Id":                       "e4445ded-3e9c-4aed-bf9b-902bc102a612",
		"comandTag.Insert()":       comandTag.Insert(),
		"comandTag.Delete()":       comandTag.Delete(),
		"comandTag.Select()":       comandTag.Select(),
		"comandTag.Update()":       comandTag.Update(),
		"comandTag.RowsAffected()": comandTag.RowsAffected(),
		"comandTag.String()":       comandTag.String(),
	}).Debug("Return data for SQL executed in database")

	return err
}
