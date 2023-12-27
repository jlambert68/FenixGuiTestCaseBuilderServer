package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
	"github.com/sirupsen/logrus"
	"time"
)

// PrepareSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers
// Do initial preparations to be able to save all supported TestInstructions, TestInstructionContainers and Allowed Users
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "d0a3d9d6-6ee3-423a-aba5-81ad53be07d3",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'PrepareSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")

		return err

	}

	defer txn.Commit(context.Background())

	// Save  all supported TestInstructions, TestInstructionContainers and Allowed Users
	err = fenixCloudDBObject.saveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
		txn,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":    "657a20ff-d1b4-4568-a665-6a519547702a",
			"error": err,
		}).Error("Problem when saving supported TestInstructions, TestInstructionContainers and Allowed Users")

		return err
	}

	return err
}

// Save all supported TestInstructions, TestInstructionContainers and Allowed Users
func (fenixCloudDBObject *FenixCloudDBObjectStruct) saveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	dbTransaction pgx.Tx,
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	err error) {

	// Verify that Domain exists in database
	var domainServiceAccountRelation *domainServiceAccountRelationStruct
	domainServiceAccountRelation, err = fenixCloudDBObject.verifyDomainExistsInDatabase(
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

	fmt.Println("************************************* domainServiceAccountRelation *************************************")
	fmt.Println(domainServiceAccountRelation)

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
	fmt.Println("************************************* Verify Signed message *************************************")

	// When there is no Message Hash then just save the message in the database
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
			dbTransaction, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

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
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

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
	domainServiceAccountRelation *domainServiceAccountRelationStruct,
	err error) {

	domainServiceAccountRelation, err = fenixCloudDBObject.loadDomainBaseData(dbTransaction, domainUUID)

	return domainServiceAccountRelation, err
}

// Do the actual save for all supported TestInstructions, TestInstructionContainers and Allowed Users to database
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	dbTransaction pgx.Tx,
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	err error) {

	fmt.Println(" **** Do the actual save for all supported TestInstructions, TestInstructionContainers and Allowed Users to database ****")
	// Create Insert Statement for TestCaseExecution that will be put on ExecutionQueue
	// Data to be inserted in the DB-table
	var dataRowToBeInsertedMultiType []interface{}
	var dataRowsToBeInsertedMultiType [][]interface{}
	dataRowsToBeInsertedMultiType = nil
	dataRowToBeInsertedMultiType = nil

	var tempsupportedtiandticandallowedusersmessageasjsonbAsByteString []byte
	tempsupportedtiandticandallowedusersmessageasjsonbAsByteString, err = json.Marshal(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.TestInstructionsAndTestInstructionsContainersAndUsersMessageHash)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.TestInstructions.TestInstructionsHash)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.TestInstructionContainers.TestInstructionContainersHash)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.AllowedUsers.AllowedUsersHash)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, string(tempsupportedtiandticandallowedusersmessageasjsonbAsByteString))
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.MessageCreationTimeStamp)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.MessageCreationTimeStamp)

	dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "INSERT INTO \"" + usedDBSchema + "\".\"TestCases\" "
	sqlToExecute = sqlToExecute + "(\"DomainUuid\", \"DomainName\", \"TestCaseUuid\", \"TestCaseName\", \"TestCaseVersion\", " +
		"\"TestCaseBasicInformationAsJsonb\", \"TestInstructionsAsJsonb\", \"TestInstructionContainersAsJsonb\", " +
		"\"TestCaseHash\", \"TestCaseExtraInformationAsJsonb\") "
	sqlToExecute = sqlToExecute + fenixCloudDBObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
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
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}
	}

	// Log response from CloudDB
	common_config.Logger.WithFields(logrus.Fields{
		"Id":                       "bea64662-3a70-4a5b-9e92-26d130983f63",
		"comandTag.Insert()":       comandTag.Insert(),
		"comandTag.Delete()":       comandTag.Delete(),
		"comandTag.Select()":       comandTag.Select(),
		"comandTag.Update()":       comandTag.Update(),
		"comandTag.RowsAffected()": comandTag.RowsAffected(),
		"comandTag.String()":       comandTag.String(),
		"sqlToExecute":             sqlToExecute,
	}).Debug("Return data for SQL executed in database")

	// No errors occurred
	return &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
	}, nil

	return err
}

// Verify changes to TestInstructions, TestInstructionContainers and Allowed Users separately
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifyChangesToTestInstructionsAndTestInstructionContainersAndAllowedUsersSeparately(
	dbTransaction pgx.Tx,
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
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
	if correctNewChangesFoundInTestInstructions == true || correctNewChangesFoundInTestInstructionContainers ||
		correctNewChangesFoundInAllowedUsers {
		common_config.Logger.WithFields(logrus.Fields{
			"id": "83536474-a037-4b67-85fc-2818f9181e38",
			"correctNewChangesFoundInTestInstructions":          correctNewChangesFoundInTestInstructions,
			"correctNewChangesFoundInTestInstructionContainers": correctNewChangesFoundInTestInstructionContainers,
			"correctNewChangesFoundInAllowedUsers":              correctNewChangesFoundInAllowedUsers,
		}).Info("Found correct changes, so update supported TestInstructions, TestInstructionContainers and Allowed Users in database")
	}

	// Save new message with supported TestInstructions, TestInstructionContainers and Allowed Users in database
	err = fenixCloudDBObject.performSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

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
	if testInstructionsMessage != nil {

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
	if testInstructionContainersMessage != nil {

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

type domainServiceAccountRelationStruct struct {
	domainUUID                 string
	domainName                 string
	serviceAccountUsedByWorker string
}

// When row is found the Domain exists and is allowed to use Fenix
// Functions also returns Service Account that should be used by calling Worker
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadDomainBaseData(
	dbTransaction pgx.Tx,
	domainUUID string) (
	domainServiceAccountRelation *domainServiceAccountRelationStruct,
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
	sqlToExecute = sqlToExecute + "SELECT domain_uuid, domain_name, callingworkerserviceaccountname "
	sqlToExecute = sqlToExecute + "FROM \"FenixDomainAdministration\".\"domains\" "
	sqlToExecute = sqlToExecute + "WHERE activated = true AND deleted = false AND "
	sqlToExecute = sqlToExecute + "domain_uuid = '" + domainUUID + "'"
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

		var tempDomainUUID string
		var tempDomainName string
		var tempDomainServiceAccountRelation string

		err = rows.Scan(
			&tempDomainUUID,
			&tempDomainName,
			&tempDomainServiceAccountRelation,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "33359dae-cda8-45f3-a57e-8dab751be154",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		domainServiceAccountRelation = &domainServiceAccountRelationStruct{
			domainUUID:                 tempDomainUUID,
			domainName:                 tempDomainName,
			serviceAccountUsedByWorker: tempDomainServiceAccountRelation,
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

	return domainServiceAccountRelation, err
}
