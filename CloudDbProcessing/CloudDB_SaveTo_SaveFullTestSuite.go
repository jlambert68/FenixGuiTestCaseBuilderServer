package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"fmt"
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

func (fenixCloudDBObject *FenixCloudDBObjectStruct) saveFullTestSuiteCommitOrRoleBack(
	dbTransactionReference *pgx.Tx,
	doCommitNotRoleBackReference *bool) {

	dbTransaction := *dbTransactionReference
	doCommitNotRoleBack := *doCommitNotRoleBackReference

	if doCommitNotRoleBack == true {
		dbTransaction.Commit(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "44efb853-20b0-4341-be64-bd8bf4897275",
		}).Debug("Doing Commit for SQL  in 'saveFullTestSuiteCommitOrRoleBack'")

	} else {
		dbTransaction.Rollback(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "9161f901-615a-4623-b161-9cbd99a7ffd6",
		}).Info("Doing Rollback for SQL  in 'saveFullTestSuiteCommitOrRoleBack'")

	}
}

// PrepareSaveFullTestCasePrepareSaveFullTestSuite
// Do initial preparations to be able to save the TestSuite
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareSaveFullTestSuite(
	fullTestSuiteMessage *fenixTestCaseBuilderServerGrpcApi.FullTestSuiteMessage) (
	returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "de9f81f1-c076-4bd3-8f3d-c951fd99bf38",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'prepareSaveFullTestSuite'")

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
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(
					common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage
	}

	// After all stuff is done, then Commit or Rollback depending on result
	var doCommitNotRoleBack bool

	// Standard is to do a Rollback
	doCommitNotRoleBack = false

	// When leaving then do the actual commit or rollback
	defer fenixCloudDBObject.saveFullTestSuiteCommitOrRoleBack(
		&txn,
		&doCommitNotRoleBack)

	// Extract Domain that Owns the TestCase
	var ownerDomainForTestCase domainForTestCaseOrTestSuiteStruct
	ownerDomainForTestCase = fenixCloudDBObject.extractOwnerDomainFromTestSuite(fullTestSuiteMessage)

	// Extract all TestCaseUuid from TestSuite
	var testCaseUuidsInTestSuite []string
	if fullTestSuiteMessage.TestCasesInTestSuite != nil {
		for _, tempTestCasesInTestSuite := range fullTestSuiteMessage.TestCasesInTestSuite.TestCasesInTestSuite {
			testCaseUuidsInTestSuite = append(testCaseUuidsInTestSuite, tempTestCasesInTestSuite.TestCaseUuid)
		}
	}

	var testInstructionsInTestSuite *[]*fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionsMessage
	var testInstructionContainersInTestSuite *[]*fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionContainersMessage
	if len(testCaseUuidsInTestSuite) > 0 {

		testInstructionsInTestSuite, testInstructionContainersInTestSuite,
			err = fenixCloudDBObject.loadTestCasesTIAndTICBelongingToTestSuite(txn, testCaseUuidsInTestSuite)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"id":    "979adb86-c8b9-44ba-8cd0-82059fd8d7c3",
				"error": err,
			}).Error("Got some problem when loading TestInstructions and TestInstructionContainers from TestSuite")

			// Set Error codes to return message
			var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
			var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

			errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
			errorCodes = append(errorCodes, errorCode)

			// Create Return message
			returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   err.Error(),
				ErrorCodes: errorCodes,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
			}

			return returnMessage
		}

	}

	// Extract all Domains that exist within all TestInstructions and TestInstructionContainers in the TestCase
	var allDomainsWithinTestCase []domainForTestCaseOrTestSuiteStruct
	allDomainsWithinTestCase = fenixCloudDBObject.extractAllDomainsWithinTestSuite(
		testInstructionsInTestSuite,
		testInstructionContainersInTestSuite)

	// Load Users all Domains
	var usersDomainsAndAuthorizations []DomainAndAuthorizationsStruct
	usersDomainsAndAuthorizations, err = fenixCloudDBObject.concatenateUsersDomainsAndDomainOpenToEveryOneToUse(
		txn, fullTestSuiteMessage.GetUserIdentification().GetGCPAuthenticatedUser())
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "f27bf395-e3fb-4106-b84a-9d46bc377e81",
			"error": err,
		}).Error("Got some problem when loading Users Domains")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   err.Error(),
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage
	}

	// Verify that User is allowed to Save TestSuite
	var userIsAllowedToSaveTestSuite bool
	var authorizationValueForOwnerDomain int64
	var authorizationValueForAllDomainsInTestSuite int64
	userIsAllowedToSaveTestSuite, authorizationValueForOwnerDomain,
		authorizationValueForAllDomainsInTestSuite, err = fenixCloudDBObject.verifyThatUserIsAllowedToSaveTestSuite(
		txn, ownerDomainForTestCase, allDomainsWithinTestCase, usersDomainsAndAuthorizations)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "c307dbf4-c2c8-4012-966f-7e70f077810d",
			"error": err,
		}).Error("Some technical database problem when trying to verify if user is allowed to save TestSuite")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   err.Error(),
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage
	}

	// User is not allowed to save TestSuite
	if userIsAllowedToSaveTestSuite == false {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "3ec4fb53-63a9-4c06-85bd-f6f0ba16cb20",
			"error": err,
		}).Error("User is not allowed to save TestSuite in database")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "User is not allowed to save TestSuite in database",
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage
	}

	// Save the TestSuite
	returnMessage, err = fenixCloudDBObject.saveFullTestSuite(
		txn, fullTestSuiteMessage, authorizationValueForOwnerDomain, authorizationValueForAllDomainsInTestSuite)

	if err != nil {
		return returnMessage
	}

	/*
		// Save the Users TestData for the TestSuite
		returnMessage, err = fenixCloudDBObject.saveTestDataForTestSuite(
			txn,
			fullTestSuiteMessage,
			fullTestSuiteMessage.GetUserIdentification().GetGCPAuthenticatedUser())

		if err != nil {
			return returnMessage
		}

	*/

	doCommitNotRoleBack = true

	return returnMessage
}

// Extract Domain that Owns the TestSuite
func (fenixCloudDBObject *FenixCloudDBObjectStruct) extractOwnerDomainFromTestSuite(
	fullTestCaseMessage *fenixTestCaseBuilderServerGrpcApi.FullTestSuiteMessage) (
	ownerDomainForTestSuite domainForTestCaseOrTestSuiteStruct) {

	// Extract the Owner Domain Uuid
	ownerDomainForTestSuite.domainUuid = fullTestCaseMessage.GetTestSuiteBasicInformation().GetDomainUuid()

	// Extract the Owner Domain Name
	ownerDomainForTestSuite.domainName = fullTestCaseMessage.GetTestSuiteBasicInformation().GetDomainName()

	return ownerDomainForTestSuite
}

// Extract all Domains that exist within all TestInstructions and TestInstructionContainers in the TestSuite
func (fenixCloudDBObject *FenixCloudDBObjectStruct) extractAllDomainsWithinTestSuite(
	testInstructionsInTestSuite *[]*fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionsMessage,
	testInstructionContainersInTestSuite *[]*fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionContainersMessage) (
	allDomainsWithinTestSuite []domainForTestCaseOrTestSuiteStruct) {

	var tempDomainsMap map[string]string
	var existsInDomainsMap bool
	tempDomainsMap = make(map[string]string)

	if testInstructionsInTestSuite != nil {

		// Loop TestInstructions for each TestCase
		for _, tempMatureTestInstructionsPerTestCase := range *testInstructionsInTestSuite {

			// Loop TestInstructions in TestCase
			for _, tempMatureTestInstruction := range tempMatureTestInstructionsPerTestCase.GetMatureTestInstructions() {

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

					// Add Domain to slice of alla Domains within TestSuite
					allDomainsWithinTestSuite = append(allDomainsWithinTestSuite, tempDomainsWithinTestCase)
				}
			}
		}
	}

	if testInstructionContainersInTestSuite != nil {

		// Loop TestInstructionContainers for each TestCase
		for _, tempMatureTestInstructionContainersPerTestCase := range *testInstructionContainersInTestSuite {

			// Loop TestInstructionContainers in TestCase
			for _, tempMatureTestInstructionContainer := range tempMatureTestInstructionContainersPerTestCase.GetMatureTestInstructionContainers() {

				// Extract the Domain for each TestInstructionContainer
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

					// Add Domain to slice of alla Domains within TestSuite
					allDomainsWithinTestSuite = append(allDomainsWithinTestSuite, tempDomainsWithinTestCase)
				}
			}
		}
	}

	return allDomainsWithinTestSuite
}

// Verify that User is allowed to Save TestSuite
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifyThatUserIsAllowedToSaveTestSuite(
	dbTransaction pgx.Tx,
	ownerDomainForTestSuite domainForTestCaseOrTestSuiteStruct,
	allDomainsWithinTestSuite []domainForTestCaseOrTestSuiteStruct,
	usersDomainsAndAuthorizations []DomainAndAuthorizationsStruct) (
	userIsAllowedToSaveTestSuite bool,
	authorizationValueForOwnerDomain int64,
	authorizationValueForAllDomainsInTestSuite int64,
	err error) {

	// List Authorization value for 'OwnerDomain' from database
	authorizationValueForOwnerDomain, err = fenixCloudDBObject.loadAuthorizationValueBasedOnDomainList(
		dbTransaction, []domainForTestCaseOrTestSuiteStruct{ownerDomainForTestSuite})

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":                      "cdcab7ce-11e8-467a-b23e-16c8fad5bfa1",
			"Error":                   err,
			"ownerDomainForTestSuite": ownerDomainForTestSuite,
		}).Error("Couldn't load Authorization vale based on Owner Domain")

		return false,
			0,
			0,
			err
	}

	// List Authorization value for all domains within TestSuite from Database

	authorizationValueForAllDomainsInTestSuite, err = fenixCloudDBObject.loadAuthorizationValueBasedOnDomainList(
		dbTransaction, allDomainsWithinTestSuite)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":                        "7b192bbf-286b-42e5-acdc-41adf57c48e6",
			"Error":                     err,
			"allDomainsWithinTestSuite": allDomainsWithinTestSuite,
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
		// CanBuildAndSaveTestSuiteOwnedByThisDomain
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
		(tempCalculatedDomainAndAuthorizations.CanBuildAndSaveTestCaseOrTestSuiteHavingTIandTICFromThisDomain & authorizationValueForAllDomainsInTestSuite) ==
			authorizationValueForAllDomainsInTestSuite

	// Are both control 'true'
	userIsAllowedToSaveTestSuite = userCanBuildAndSaveTestCaseOwnedByThisDomain && userCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain

	return userIsAllowedToSaveTestSuite,
		authorizationValueForOwnerDomain,
		authorizationValueForAllDomainsInTestSuite,
		err
}

// Save the full TestSuite to CloudDB
func (fenixCloudDBObject *FenixCloudDBObjectStruct) saveFullTestSuite(
	dbTransaction pgx.Tx,
	fullTestSuiteMessage *fenixTestCaseBuilderServerGrpcApi.FullTestSuiteMessage,
	authorizationValueForOwnerDomain int64,
	authorizationValueForAllDomainsInTestSuite int64) (
	returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse,
	err error) {

	nexTestSuiteVersion, err := fenixCloudDBObject.getNexTestSuiteVersion(fullTestSuiteMessage.GetTestSuiteBasicInformation().GetTestSuiteUuid())
	if err != nil {
		if err != nil {

			// Set Error codes to return message
			var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
			var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

			errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
			errorCodes = append(errorCodes, errorCode)

			// Create Return message
			returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   "Problem when getting next TestSuiteVersion from database",
				ErrorCodes: errorCodes,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
			}
		}

		return returnMessage, err

	}

	// Check if Next TestSuiteVersion number is the same as set in TestSuite to be saved
	if nexTestSuiteVersion != fullTestSuiteMessage.GetTestSuiteBasicInformation().GetTestSuiteVersion() {
		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack: false,
			Comments: fmt.Sprintf("TestSuiteVersion in TestSuite to be saved, '%d' is not the same as the next TestSuiteVersion in database, '%d'",
				fullTestSuiteMessage.GetTestSuiteBasicInformation().GetTestSuiteVersion(),
				nexTestSuiteVersion),
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage, err
	}

	// Extract column data to be added to data-row
	tempDomainUuid := fullTestSuiteMessage.GetTestSuiteBasicInformation().GetDomainUuid()
	tempDomainName := fullTestSuiteMessage.GetTestSuiteBasicInformation().GetDomainName()
	tempTestSuiteUuid := fullTestSuiteMessage.GetTestSuiteBasicInformation().GetTestSuiteUuid()
	tempTestSuiteName := fullTestSuiteMessage.GetTestSuiteBasicInformation().GetTestSuiteName()
	tempTestSuiteVersion := fullTestSuiteMessage.GetTestSuiteBasicInformation().GetTestSuiteVersion()
	tempTestSuiteHash := fullTestSuiteMessage.GetMessageHash()
	// CanListAndViewTestSuiteAuthorizationLevelOwnedByDomain
	// CanListAndViewTestSuiteAuthorizationLevelHavingTiAndTicWith
	insertTimeStamp := shared_code.GenerateDatetimeTimeStampForDB()
	tempInsertedByUserIdOnComputer := fullTestSuiteMessage.GetUserIdentification().GetUserIdOnComputer()
	tempInsertedByGCPAuthenticatedUser := fullTestSuiteMessage.GetUserIdentification().GetGCPAuthenticatedUser()

	tempTestSuiteIsDeleted := false

	// Initiate 'TestCasesInTestSuite' if nil
	if fullTestSuiteMessage.GetTestCasesInTestSuite() == nil {
		fullTestSuiteMessage.TestCasesInTestSuite = &fenixTestCaseBuilderServerGrpcApi.TestCasesInTestSuiteMessage{}
	}
	tempTestCasesInTestSuiteAsJsonb := protojson.Format(fullTestSuiteMessage.GetTestCasesInTestSuite())
	// TestSuitePreviewAsJsonb - below

	// Initiate 'TestSuiteMetaData' if nil
	if fullTestSuiteMessage.GetTestSuiteMetaData() == nil {
		fullTestSuiteMessage.TestSuiteMetaData = &fenixTestCaseBuilderServerGrpcApi.UserSpecifiedTestSuiteMetaDataMessage{}
	}
	tempTestSuiteMetaDataAsJsonb := protojson.Format(fullTestSuiteMessage.GetTestSuiteMetaData())

	// Initiate 'TestSuitePreview' if nil
	if fullTestSuiteMessage.GetTestSuitePreview() == nil {
		fullTestSuiteMessage.TestSuitePreview = &fenixTestCaseBuilderServerGrpcApi.TestSuitePreviewMessage{}

	} else {

		// finish Preview-structure to be saved
		fullTestSuiteMessage.TestSuitePreview.TestSuitePreview.TestSuiteVersion = strconv.Itoa(int(nexTestSuiteVersion))
		fullTestSuiteMessage.TestSuitePreview.TestSuitePreview.LastSavedTimeStamp = insertTimeStamp
	}
	tempTestSuitePreviewAsJsonb := protojson.Format(fullTestSuiteMessage.GetTestSuitePreview())

	var dataRowToBeInsertedMultiType []interface{}
	var dataRowsToBeInsertedMultiType [][]interface{}

	// Create Insert Statement for TestSuite that will be saved in Database
	// Data to be inserted in the DB-table
	dataRowsToBeInsertedMultiType = nil

	dataRowToBeInsertedMultiType = nil

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempDomainUuid)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempDomainName)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestSuiteUuid)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestSuiteName)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestSuiteVersion)

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestSuiteHash)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, authorizationValueForOwnerDomain)

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, authorizationValueForAllDomainsInTestSuite)

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, insertTimeStamp)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempInsertedByUserIdOnComputer)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempInsertedByGCPAuthenticatedUser)

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestSuiteIsDeleted)

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestCasesInTestSuiteAsJsonb)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestSuitePreviewAsJsonb)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestSuiteMetaDataAsJsonb)

	dataRowsToBeInsertedMultiType = append(dataRowsToBeInsertedMultiType, dataRowToBeInsertedMultiType)

	/*
		/*
		    "DomainUuid"                                                  uuid not null,
		    "DomainName"                                                  varchar not null,
		    "TestSuiteUuid"                                               uuid not null,
		    "TestSuiteName"                                               varchar not null,
		    "TestSuiteVersion"                                            integer not null
		,
		    "TestSuiteHash"                                               varchar not null,

		    "CanListAndViewTestSuiteAuthorizationLevelOwnedByDomain"      bigint not null,
		    "CanListAndViewTestSuiteAuthorizationLevelHavingTiAndTicWith" bigint not null,


		    "InsertTimeStamp"                                             timestamp not null,
		    "InsertedByUserIdOnComputer"                                  varchar not null,
		    "InsertedByGCPAuthenticatedUser"                              varchar not null,

		    "TestSuiteIsDeleted"                                          boolean not null,

		    "DeleteTimestamp"                                             timestamp default '2068-11-18 00:00:00'::timestamp without time zone,
		    "DeletedInsertedTimeStamp"                                    timestamp,
		    "DeletedByUserIdOnComputer"                                   varchar,
		    "DeletedByGCPAuthenticatedUser"                               varchar,

		    "TestCasesInTestSuite"                                        jsonb   not null,
		    "TestSuitePreview"                                            jsonb   not null,
		    "TestSuiteMetaData"                                           jsonb   not null,
		    "UniqueCounter"                                               serial,

	*/

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "INSERT INTO \"FenixBuilder\".\"TestSuites\" "
	sqlToExecute = sqlToExecute + "(\"DomainUuid\", \"DomainName\", \"TestSuiteUuid\", \"TestSuiteName\", " +
		"\"TestSuiteVersion\", \"TestSuiteHash\", " +
		"\"CanListAndViewTestSuiteAuthorizationLevelOwnedByDomain\", " +
		"\"CanListAndViewTestSuiteAuthorizationLevelHavingTiAndTicWith\", " +
		"\"InsertTimeStamp\", \"InsertedByUserIdOnComputer\", \"InsertedByGCPAuthenticatedUser\", " +
		"\"TestSuiteIsDeleted\", " +
		" \"TestCasesInTestSuite\", \"TestSuitePreview\", \"TestSuiteMetaData\") "
	sqlToExecute = sqlToExecute + fenixCloudDBObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
	sqlToExecute = sqlToExecute + ";"

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "7867cae8-2869-4898-be0b-bcc3ba4570e2",
			"sqlToExecute": sqlToExecute,
		}).Error("Problem when Saving TestSuite to database")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "Problem when Saving TestCase to database",
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage, err
	}

	// Log response from CloudDB
	common_config.Logger.WithFields(logrus.Fields{
		"Id":                       "e44546c3-922c-4d9b-a1e4-94bb655825d9",
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
		AckNack:    true,
		Comments:   "",
		ErrorCodes: nil,
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
			CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
	}, nil

}

// Load All TestCases TestInstructions and TestInstructionContainers belonging to a TestSuite
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadTestCasesTIAndTICBelongingToTestSuite(
	dbTransaction pgx.Tx,
	testCasesUuid []string) (
	_ *[]*fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionsMessage,
	_ *[]*fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionContainersMessage,
	err error) {

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "WITH uniquecounters AS ( "
	sqlToExecute = sqlToExecute + "SELECT Distinct ON (\"TestCaseUuid\")  \"UniqueCounter "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestCases "
	sqlToExecute = sqlToExecute + "WHERE \"TestCaseUuid\" IN " + fenixCloudDBObject.generateSQLINArray(testCasesUuid) + " "
	sqlToExecute = sqlToExecute + "ORDER BY \"TestCaseUuid\", \"UniqueCounter\" DESC "
	sqlToExecute = sqlToExecute + ") "
	sqlToExecute = sqlToExecute + "SELECT t.\"TestInstructionsAsJsonb\", t.\"TestInstructionContainersAsJsonb\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestCases\" t "
	sqlToExecute = sqlToExecute + "WHERE t.\"UniqueCounter\" IN (SELECT * FROM uniquecounters); "
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "2ddf3a3a-98e1-4dc4-91b7-3fe7b7da34eb",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadFullTestCase'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "0dee732b-15dc-4b6c-99e8-b555cceb82a2",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, nil, err
	}

	var (
		tempTestInstructionsAsJson               string
		tempTestInstructionContainersAsJson      string
		tempTestInstructionsAsByteArray          []byte
		tempTestInstructionContainersAsByteArray []byte
		allTestInstructions                      []*fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionsMessage
		allTestInstructionContainers             []*fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionContainersMessage
	)

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&tempTestInstructionsAsJson,
			&tempTestInstructionContainersAsJson,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "c562fc51-20a9-4d5e-aa04-7653469dd399",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, nil, err
		}

		// Convert json-strings into byte-arrays
		tempTestInstructionsAsByteArray = []byte(tempTestInstructionsAsJson)
		tempTestInstructionContainersAsByteArray = []byte(tempTestInstructionContainersAsJson)

		// Convert json-byte-arrays into proto-messages - tempTestInstructionsAsByteArray
		var tempTestInstructions fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionsMessage
		err = protojson.Unmarshal(tempTestInstructionsAsByteArray, &tempTestInstructions)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "3055f791-3b56-4458-a778-eb71226a94b8",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestInstructionsAsByteArray' into proto-message")

			return nil, nil, err
		}

		// Convert json-byte-arrays into proto-messages - tempTestInstructionContainersAsByteArray
		var tempTestInstructionContainers fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionContainersMessage
		err = protojson.Unmarshal(tempTestInstructionContainersAsByteArray, &tempTestInstructionContainers)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "37d5c732-ed9a-4cea-aa46-c8a12a0f076b",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestInstructionContainersAsByteArray' into proto-message")

			return nil, nil, err
		}

		// Add TestInstructions to slice of all TestInstructions
		allTestInstructions = append(allTestInstructions, &tempTestInstructions)

		// Add TestInstructionContainerss to slice of all TestInstructionContainers
		allTestInstructionContainers = append(allTestInstructionContainers, &tempTestInstructionContainers)

	}

	return &allTestInstructions, &allTestInstructionContainers, err

}
