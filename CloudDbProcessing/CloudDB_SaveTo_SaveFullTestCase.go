package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/shared_code"
	"strconv"
	"time"

	//"database/sql/driver"
	//"encoding/json"
	//"errors"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
)

func (fenixCloudDBObject *FenixCloudDBObjectStruct) saveFullTestCaseCommitOrRoleBack(
	dbTransactionReference *pgx.Tx,
	doCommitNotRoleBackReference *bool) {

	dbTransaction := *dbTransactionReference
	doCommitNotRoleBack := *doCommitNotRoleBackReference

	if doCommitNotRoleBack == true {
		dbTransaction.Commit(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "c2a383a8-07a3-46c6-9a7d-f327a7f7450b",
		}).Debug("Doing Commit for SQL  in 'saveFullTestCaseCommitOrRoleBack'")

	} else {
		dbTransaction.Rollback(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "43150f43-8b9f-4ddf-a444-e18788910833",
		}).Info("Doing Rollback for SQL  in 'saveFullTestCaseCommitOrRoleBack'")

	}
}

// PrepareSaveFullTestCase
// Do initial preparations to be able to save the TestCase
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareSaveFullTestCase(
	fullTestCaseMessage *fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage) (
	returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
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
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage
	}

	// After all stuff is done, then Commit or Rollback depending on result
	var doCommitNotRoleBack bool

	// Standard is to do a Rollback
	doCommitNotRoleBack = false

	// When leaving then do the actual commit or rollback
	defer fenixCloudDBObject.saveFullTestCaseCommitOrRoleBack(
		&txn,
		&doCommitNotRoleBack)

	// Extract Domain that Owns the TestCase
	var ownerDomainForTestCase domainForTestCaseOrTestSuiteStruct
	ownerDomainForTestCase = fenixCloudDBObject.extractOwnerDomainFromTestCase(fullTestCaseMessage)

	// Extract all Domains that exist within all TestInstructions and TestInstructionContainers in the TestCase
	var allDomainsWithinTestCase []domainForTestCaseOrTestSuiteStruct
	allDomainsWithinTestCase = fenixCloudDBObject.extractAllDomainsWithinTestCase(fullTestCaseMessage)

	// Load Users all Domains
	var usersDomainsAndAuthorizations []DomainAndAuthorizationsStruct
	usersDomainsAndAuthorizations, err = fenixCloudDBObject.concatenateUsersDomainsAndDomainOpenToEveryOneToUse(
		txn, fullTestCaseMessage.TestCaseBasicInformation.GetUserIdentification().GetGCPAuthenticatedUser())
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "20737f52-eaad-4818-ae51-60ff6aa74a79",
			"error": err,
		}).Error("Got some problem when loading Users Domains")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     err.Error(),
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage
	}

	// Verify that User is allowed to Save TestCase
	var userIsAllowedToSaveTestCase bool
	var authorizationValueForOwnerDomain int64
	var authorizationValueForAllDomainsInTestCase int64
	userIsAllowedToSaveTestCase, authorizationValueForOwnerDomain,
		authorizationValueForAllDomainsInTestCase, err = fenixCloudDBObject.verifyThatUserIsAllowedToSaveTestCase(
		txn, ownerDomainForTestCase, allDomainsWithinTestCase, usersDomainsAndAuthorizations)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "6211ffda-79e4-46dc-82e8-1f0a8a666213",
			"error": err,
		}).Error("Some technical database problem when trying to verify if user is allowed to save TestCase")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     err.Error(),
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage
	}

	// User is not allowed to save TestCase
	if userIsAllowedToSaveTestCase == false {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "0b646cac-c7d4-4c05-81c9-158d2fbfc9b9",
			"error": err,
		}).Error("User is not allowed to save TestCase in database")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "User is not allowed to save TestCase in database",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage
	}

	// Save the TestCase
	returnMessage, err = fenixCloudDBObject.saveFullTestCase(
		txn, fullTestCaseMessage, authorizationValueForOwnerDomain, authorizationValueForAllDomainsInTestCase)

	if err != nil {
		return returnMessage
	}

	// Save the Users TestData for the TestCase
	returnMessage, err = fenixCloudDBObject.saveUsersTestDataForTestCase(
		txn,
		fullTestCaseMessage,
		fullTestCaseMessage.TestCaseBasicInformation.GetUserIdentification().GetGCPAuthenticatedUser())

	if err != nil {
		return returnMessage
	}

	doCommitNotRoleBack = true

	return returnMessage
}

// Struct used when extracting the Owner Domain for a TestCase or a TestSuite
type domainForTestCaseOrTestSuiteStruct struct {
	domainUuid string
	domainName string
}

// Extract Domain that Owns the TestCase
func (fenixCloudDBObject *FenixCloudDBObjectStruct) extractOwnerDomainFromTestCase(
	fullTestCaseMessage *fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage) (
	ownerDomainForTestCase domainForTestCaseOrTestSuiteStruct) {

	// Extract the Owner Domain Uuid
	ownerDomainForTestCase.domainUuid = fullTestCaseMessage.GetTestCaseBasicInformation().
		GetBasicTestCaseInformation().GetNonEditableInformation().GetDomainUuid()

	// Extract the Owner Domain Name
	ownerDomainForTestCase.domainName = fullTestCaseMessage.GetTestCaseBasicInformation().
		GetBasicTestCaseInformation().GetNonEditableInformation().GetDomainName()

	return ownerDomainForTestCase
}

// Extract all Domains that exist within all TestInstructions and TestInstructionContainers in the TestCase
func (fenixCloudDBObject *FenixCloudDBObjectStruct) extractAllDomainsWithinTestCase(
	fullTestCaseMessage *fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage) (
	allDomainsWithinTestCase []domainForTestCaseOrTestSuiteStruct) {

	var tempDomainsMap map[string]string
	var existsInDomainsMap bool
	tempDomainsMap = make(map[string]string)

	// Extract the Domain for each TestInstruction
	for _, tempMatureTestInstruction := range fullTestCaseMessage.GetMatureTestInstructions().
		GetMatureTestInstructions() {

		var tempDomainsWithinTestCase domainForTestCaseOrTestSuiteStruct
		tempDomainsWithinTestCase = domainForTestCaseOrTestSuiteStruct{
			domainUuid: tempMatureTestInstruction.GetBasicTestInstructionInformation().GetNonEditableInformation().
				GetDomainUuid(),
			domainName: tempMatureTestInstruction.GetBasicTestInstructionInformation().GetNonEditableInformation().
				GetDomainName(),
		}

		// Check if the Domain already exists in 'tempDomainsMap'
		_, existsInDomainsMap = tempDomainsMap[tempDomainsWithinTestCase.domainUuid]

		// Only store the Domain is missing in map
		if existsInDomainsMap == false {

			// Add to Map
			tempDomainsMap[tempDomainsWithinTestCase.domainUuid] = tempDomainsWithinTestCase.domainUuid

			// Add Domain to slice of alla Domains within TestCase
			allDomainsWithinTestCase = append(allDomainsWithinTestCase, tempDomainsWithinTestCase)
		}
	}

	// Extract the Domain for each TestInstructionContainer
	for _, tempMatureTestInstructionContainer := range fullTestCaseMessage.GetMatureTestInstructionContainers().
		GetMatureTestInstructionContainers() {

		var tempDomainsWithinTestCase domainForTestCaseOrTestSuiteStruct
		tempDomainsWithinTestCase = domainForTestCaseOrTestSuiteStruct{
			domainUuid: tempMatureTestInstructionContainer.GetBasicTestInstructionContainerInformation().
				GetNonEditableInformation().GetDomainUuid(),
			domainName: tempMatureTestInstructionContainer.GetBasicTestInstructionContainerInformation().
				GetNonEditableInformation().GetDomainName(),
		}

		// Check if the Domain already exists in 'tempDomainsMap'
		_, existsInDomainsMap = tempDomainsMap[tempDomainsWithinTestCase.domainUuid]

		// Only store the Domain is missing in map
		if existsInDomainsMap == false {

			// Add to Map
			tempDomainsMap[tempDomainsWithinTestCase.domainUuid] = tempDomainsWithinTestCase.domainUuid

			// Add Domain to slice of alla Domains within TestCase
			allDomainsWithinTestCase = append(allDomainsWithinTestCase, tempDomainsWithinTestCase)
		}
	}

	return allDomainsWithinTestCase
}

// Verify that User is allowed to Save TestCase
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifyThatUserIsAllowedToSaveTestCase(
	dbTransaction pgx.Tx,
	ownerDomainForTestCase domainForTestCaseOrTestSuiteStruct,
	allDomainsWithinTestCase []domainForTestCaseOrTestSuiteStruct,
	usersDomainsAndAuthorizations []DomainAndAuthorizationsStruct) (
	userIsAllowedToSaveTestCase bool,
	authorizationValueForOwnerDomain int64,
	authorizationValueForAllDomainsInTestCase int64,
	err error) {

	// List Authorization value for 'OwnerDomain' from database
	authorizationValueForOwnerDomain, err = fenixCloudDBObject.loadAuthorizationValueBasedOnDomainList(
		dbTransaction, []domainForTestCaseOrTestSuiteStruct{ownerDomainForTestCase})

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":                     "cdcab7ce-11e8-467a-b23e-16c8fad5bfa1",
			"Error":                  err,
			"ownerDomainForTestCase": ownerDomainForTestCase,
		}).Error("Couldn't load Authorization vale based on Owner Domain")

		return false,
			0,
			0,
			err
	}

	// List Authorization value for all domains within TestCase from Database
	authorizationValueForAllDomainsInTestCase, err = fenixCloudDBObject.loadAuthorizationValueBasedOnDomainList(
		dbTransaction, allDomainsWithinTestCase)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":                       "cdcab7ce-11e8-467a-b23e-16c8fad5bfa1",
			"Error":                    err,
			"allDomainsWithinTestCase": allDomainsWithinTestCase,
		}).Error("Couldn't load Authorization vale based on Owner Domain")

		return false,
			0,
			0,
			err
	}

	// Calculate the Authorization requirements
	var tempCalculatedDomainAndAuthorizations DomainAndAuthorizationsStruct
	var domainList []string
	for _, domainAndAuthorization := range usersDomainsAndAuthorizations {
		// Add to DomainList
		domainList = append(domainList, domainAndAuthorization.DomainUuid)

		// Calculate the Authorization requirements for...
		// CanBuildAndSaveTestCaseOrTestSuiteOwnedByThisDomain
		tempCalculatedDomainAndAuthorizations.CanBuildAndSaveTestCaseOrTestSuiteOwnedByThisDomain =
			tempCalculatedDomainAndAuthorizations.CanBuildAndSaveTestCaseOrTestSuiteOwnedByThisDomain +
				domainAndAuthorization.CanBuildAndSaveTestCaseOrTestSuiteOwnedByThisDomain

		// CanBuildAndSaveTestCaseOrTestSuiteHavingTIandTICFromThisDomain
		tempCalculatedDomainAndAuthorizations.CanBuildAndSaveTestCaseOrTestSuiteHavingTIandTICFromThisDomain =
			tempCalculatedDomainAndAuthorizations.CanBuildAndSaveTestCaseOrTestSuiteHavingTIandTICFromThisDomain +
				domainAndAuthorization.CanBuildAndSaveTestCaseOrTestSuiteHavingTIandTICFromThisDomain
	}

	// Check if User can Save TestCase due to 'CanBuildAndSaveTestCaseOrTestSuiteOwnedByThisDomain'
	var userCanBuildAndSaveTestCaseOwnedByThisDomain bool
	userCanBuildAndSaveTestCaseOwnedByThisDomain =
		(tempCalculatedDomainAndAuthorizations.CanBuildAndSaveTestCaseOrTestSuiteOwnedByThisDomain & authorizationValueForOwnerDomain) ==
			authorizationValueForOwnerDomain

	// Check if User canSave TestCase due to 'CanBuildAndSaveTestCaseOrTestSuiteHavingTIandTICFromThisDomain'
	var userCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain bool
	userCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain =
		(tempCalculatedDomainAndAuthorizations.CanBuildAndSaveTestCaseOrTestSuiteHavingTIandTICFromThisDomain & authorizationValueForAllDomainsInTestCase) ==
			authorizationValueForAllDomainsInTestCase

	// Are both control 'true'
	userIsAllowedToSaveTestCase = userCanBuildAndSaveTestCaseOwnedByThisDomain && userCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain

	return userIsAllowedToSaveTestCase,
		authorizationValueForOwnerDomain,
		authorizationValueForAllDomainsInTestCase,
		err
}

// Load Authorization value based on Domain List
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadAuthorizationValueBasedOnDomainList(
	dbTransaction pgx.Tx,
	domainList []domainForTestCaseOrTestSuiteStruct) (
	authorizationValue int64,
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "29f64855-9c6a-49f1-ade0-bdf7020deb7e",
	}).Debug("Entering: loadAuthorizationValueBasedOnDomainList()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":                 "47c7f7b3-9dce-4230-89c8-dd66605b5f7a",
			"authorizationValue": authorizationValue,
		}).Debug("Exiting: loadAuthorizationValueBasedOnDomainList()")
	}()

	// only process of there are any domains in the list
	if len(domainList) == 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":         "2ae322c3-82a3-4671-93ba-de24819e709c",
			"domainList": domainList,
		}).Debug("domainList is empty")

		return 0, err
	}

	// Convert into 'pure' string array
	var tempDomainUuidSlice []string
	for _, tempDomain := range domainList {
		tempDomainUuidSlice = append(tempDomainUuidSlice, tempDomain.domainUuid)
	}

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT SUM(authvalue.bitnumbervalue) "
	sqlToExecute = sqlToExecute + "FROM \"FenixDomainAdministration\".domains dom,  " +
		"\"FenixDomainAdministration\".domainbitpositionenum authvalue "
	sqlToExecute = sqlToExecute + "WHERE dom.domain_uuid  IN " + common_config.GenerateSQLINArray(tempDomainUuidSlice)
	sqlToExecute = sqlToExecute + " AND "
	sqlToExecute = sqlToExecute + "dom.bitnumbername = authvalue.bitnumbername "
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "4a118124-af61-4699-8d59-ce31e7611f7e",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadUsersDomains'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "9177c263-04d0-411d-8ac2-148279038fb3",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return 0, err
	}

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&authorizationValue,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "a0d518f7-20b5-4a2d-8e60-ec1d0c975409",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return 0, err
		}

		break

	}

	return authorizationValue, err
}

// Save the full TestCase to CloudDB
func (fenixCloudDBObject *FenixCloudDBObjectStruct) saveFullTestCase(
	dbTransaction pgx.Tx,
	fullTestCaseMessage *fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage,
	authorizationValueForOwnerDomain int64,
	authorizationValueForAllDomainsInTestCase int64) (
	returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse,
	err error) {

	usedDBSchema := "FenixBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	nexTestCaseVersion, err := fenixCloudDBObject.getNexTestCaseVersion(fullTestCaseMessage.TestCaseBasicInformation.BasicTestCaseInformation.NonEditableInformation.TestCaseUuid)
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
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
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
	tempDomainHash := fullTestCaseMessage.MessageHash
	tempTestCaseExtraInformationAsJsonb := protojson.Format(fullTestCaseMessage.TestCaseExtraInformation)
	tempTestCaseTemplateFilesAsJsonb := protojson.Format(fullTestCaseMessage.TestCaseTemplateFiles)
	timeStampToUse := shared_code.GenerateDatetimeTimeStampForDB()

	// finish Preview-structure to be saved
	fullTestCaseMessage.TestCasePreview.TestCasePreview.TestCaseVersion = strconv.Itoa(int(nexTestCaseVersion))
	fullTestCaseMessage.TestCasePreview.TestCasePreview.LastSavedTimeStamp = timeStampToUse
	testCasePreviewAsJsonb := protojson.Format(fullTestCaseMessage.TestCasePreview)

	// Convert TestCaseMetaData into jsonb
	testCaseMetaDataAsJsonb := protojson.Format(fullTestCaseMessage.TestCaseMetaData)

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
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempDomainHash)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestCaseExtraInformationAsJsonb)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, authorizationValueForOwnerDomain)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, authorizationValueForAllDomainsInTestCase)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, false)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestCaseTemplateFilesAsJsonb)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, timeStampToUse)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCasePreviewAsJsonb)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, testCaseMetaDataAsJsonb)

	dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "INSERT INTO \"" + usedDBSchema + "\".\"TestCases\" "
	sqlToExecute = sqlToExecute + "(\"DomainUuid\", \"DomainName\", \"TestCaseUuid\", \"TestCaseName\", \"TestCaseVersion\", " +
		"\"TestCaseBasicInformationAsJsonb\", \"TestInstructionsAsJsonb\", \"TestInstructionContainersAsJsonb\", " +
		"\"TestCaseHash\", \"TestCaseExtraInformationAsJsonb\", \"CanListAndViewTestCaseAuthorizationLevelOwnedByDomain\", " +
		"\"CanListAndViewTestCaseAuthorizationLevelHavingTiAndTicWithDomai\", \"TestCaseIsDeleted\", " +
		"\"TestCaseTemplateFilesAsJsonb\", \"InsertTimeStamp\", \"TestCasePreview\", \"TestCaseMetaData\")  "
	sqlToExecute = sqlToExecute + fenixCloudDBObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
	sqlToExecute = sqlToExecute + ";"

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "9f116364-c38b-4bd8-bea9-3ec307c76dfa",
			"sqlToExecute": sqlToExecute,
		}).Error("Problem when Saving TestCase to database")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "Problem when Saving TestCase to database",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage, err
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

}

// Save Users TestData for a TestCase to CloudDB
func (fenixCloudDBObject *FenixCloudDBObjectStruct) saveUsersTestDataForTestCase(
	dbTransaction pgx.Tx,
	fullTestCaseMessage *fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage,
	gcpAuthenticatedUser string) (
	returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse,
	err error) {

	// Extract column data to be added to data-row
	tempDomainUuid := fullTestCaseMessage.TestCaseBasicInformation.BasicTestCaseInformation.NonEditableInformation.DomainUuid
	tempDomainName := fullTestCaseMessage.TestCaseBasicInformation.BasicTestCaseInformation.NonEditableInformation.DomainName
	tempTestCaseUuid := fullTestCaseMessage.TestCaseBasicInformation.BasicTestCaseInformation.NonEditableInformation.TestCaseUuid
	tempTestCaseName := fullTestCaseMessage.TestCaseBasicInformation.BasicTestCaseInformation.EditableInformation.TestCaseName
	tempTestCaseVersion := fullTestCaseMessage.TestCaseBasicInformation.BasicTestCaseInformation.NonEditableInformation.TestCaseVersion
	tempTestDataAsJsonb := protojson.Format(fullTestCaseMessage.GetTestCaseTestData())
	tempTestDataHash := fullTestCaseMessage.GetTestCaseTestData().GetHashOfThisMessageWithEmptyHashField()

	var dataRowToBeInsertedMultiType []interface{}
	var dataRowsToBeInsertedMultiType [][]interface{}

	// Create Insert Statement
	// Data to be inserted in the DB-table
	dataRowsToBeInsertedMultiType = nil

	dataRowToBeInsertedMultiType = nil

	timeStampToUse := shared_code.GenerateDatetimeTimeStampForDB()

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempDomainUuid)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempDomainName)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestCaseUuid)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestCaseName)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestCaseVersion)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestDataAsJsonb)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestDataHash)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, false)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, timeStampToUse)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, gcpAuthenticatedUser)

	dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "INSERT INTO \"FenixBuilder\".\"UsersTestDataForTestCase\" "
	sqlToExecute = sqlToExecute + "(\"DomainUuid\", \"DomainName\", \"TestCaseUuid\", \"TestCaseName\", \"TestCaseVersion\", " +
		"\"TestData\", \"TestDataHash\", \"TestCaseIsDeleted\", \"InsertedTimeStamp\", \"GcpAuthenticatedUser\")  "
	sqlToExecute = sqlToExecute + fenixCloudDBObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
	sqlToExecute = sqlToExecute + ";"

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "f2956cdc-02e7-4e3b-b820-abe12608260d",
			"sqlToExecute": sqlToExecute,
		}).Error("Problem when Saving Users TestData for TestCase to database")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "Problem inserting Users TestData for TestCase",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage, err
	}

	// Log response from CloudDB
	common_config.Logger.WithFields(logrus.Fields{
		"Id":                       "2540118c-ad2f-4d68-bceb-abbed93d4891",
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

}
