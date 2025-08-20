package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/shared_code"
	"strings"

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
			"id": "5a680f7c-9960-48d4-99b6-8fc0228b2fd8",
		}).Debug("Doing Commit for SQL  in 'saveFullTestSuiteCommitOrRoleBack'")

	} else {
		dbTransaction.Rollback(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "3b5d95a1-6c35-45f4-add9-6860ffd3ceb0",
		}).Info("Doing Rollback for SQL  in 'saveFullTestSuiteCommitOrRoleBack'")

	}
}

// PrepareSaveFullTestSuite
// Do initial preparations to be able to save the TestSuite
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareSaveFullTestSuite(
	gRPCTestSuiteMessage *fenixTestCaseBuilderServerGrpcApi.SaveFullTestSuiteMessageRequest) (
	returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse) {

	// Extract full TestSuiteMessage
	var fullTestSuiteMessage *fenixTestCaseBuilderServerGrpcApi.FullTestSuiteMessage
	fullTestSuiteMessage = gRPCTestSuiteMessage.GetTestSuite()

	// Extract UserIdentification
	var userIdentification *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage
	userIdentification = gRPCTestSuiteMessage.GetUserIdentification()

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "b7270c79-7a71-41ae-b2b8-ae0a0f45904c",
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

	// Extract Domain that Owns the TestSuite
	var ownerDomainForTestSuite domainForTestCaseOrTestSuiteStruct
	ownerDomainForTestSuite = fenixCloudDBObject.extractOwnerDomainFromTestSuite(fullTestSuiteMessage)

	// Extract all TestCaseUuid's from TestSuite
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
				"id":    "d5ec1a39-b6ab-4691-8c69-e73bff0dd28e",
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
		txn, userIdentification.GetGCPAuthenticatedUser())
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "59eb5990-3067-4b11-b38b-a854a6cc22b1",
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
		txn, ownerDomainForTestSuite, allDomainsWithinTestCase, usersDomainsAndAuthorizations)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "5509decd-a0e9-4d4c-a954-38f0c72a4c41",
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
			"id":    "41aaf01e-4f4c-4d76-8b39-541ae6cf3b39",
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

	// Load PreViews for all TestCases in TestSuite
	var tempTestCasesPreview []*fenixTestCaseBuilderServerGrpcApi.TestCasePreviewMessage
	tempTestCasesPreview, err = fenixCloudDBObject.loadTestCasesPreviewForTestSuite(txn, testCaseUuidsInTestSuite)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "734a3f81-cc80-4393-880a-95c3eaad6863",
			"error": err,
		}).Error("Got some problem when loading TestCasesPreView from database")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "Got some problem when loading TestCasesPreView from database",
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage
	}

	// Save the TestSuite
	returnMessage, err = fenixCloudDBObject.saveFullTestSuite(
		txn,
		gRPCTestSuiteMessage,
		authorizationValueForOwnerDomain,
		authorizationValueForAllDomainsInTestSuite,
		tempTestCasesPreview)

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
	fullTestSuiteMessage *fenixTestCaseBuilderServerGrpcApi.FullTestSuiteMessage) (
	ownerDomainForTestSuite domainForTestCaseOrTestSuiteStruct) {

	// Extract the Owner Domain Uuid
	ownerDomainForTestSuite.domainUuid = fullTestSuiteMessage.GetTestSuiteBasicInformation().GetDomainUuid()

	// Extract the Owner Domain Name
	ownerDomainForTestSuite.domainName = fullTestSuiteMessage.GetTestSuiteBasicInformation().GetDomainName()

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
			"Id":                      "7c585333-5f74-4d82-8083-9335d10e4b39",
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
			"Id":                        "c11f3d87-8ef5-46d1-aa82-c9a28d462ec8",
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
	gRPCTestSuiteMessage *fenixTestCaseBuilderServerGrpcApi.SaveFullTestSuiteMessageRequest,
	authorizationValueForOwnerDomain int64,
	authorizationValueForAllDomainsInTestSuite int64,
	tempTestCasesPreview []*fenixTestCaseBuilderServerGrpcApi.TestCasePreviewMessage) (
	returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse,
	err error) {

	// Extract full TestSuiteMessage
	var fullTestSuiteMessage *fenixTestCaseBuilderServerGrpcApi.FullTestSuiteMessage
	fullTestSuiteMessage = gRPCTestSuiteMessage.GetTestSuite()

	// Extract UserIdentification
	var userIdentification *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage
	userIdentification = gRPCTestSuiteMessage.GetUserIdentification()

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
	tempInsertedByUserIdOnComputer := userIdentification.GetUserIdOnComputer()
	tempInsertedByGCPAuthenticatedUser := userIdentification.GetGCPAuthenticatedUser()
	tempTestSuiteType := fullTestSuiteMessage.GetTestSuiteType().GetTestSuiteType()
	tempTestSuiteTypeName := fullTestSuiteMessage.GetTestSuiteType().GetTestSuiteTypeName()
	tempTestSuiteDescription := fullTestSuiteMessage.GetTestSuiteBasicInformation().GetTestSuiteDescription()
	tempTestSuiteExecutionEnvironment := fullTestSuiteMessage.GetTestSuiteBasicInformation().GetTestSuiteExecutionEnvironment()

	tempTestSuiteIsDeleted := false
	tempTestSuiteDeletedInsertTimeStamp := "NULL"
	// TODO fixa så delete i framtiden kan leva vidare när man spara ny versionen så ska Delete-date i framtiden kopieras vidare

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

	// Initiate 'TestSuiteTestData' if nil
	if fullTestSuiteMessage.GetTestSuiteTestData() == nil {
		fullTestSuiteMessage.TestSuiteTestData = &fenixTestCaseBuilderServerGrpcApi.UsersChosenTestDataForTestSuiteMessage{}
	}
	tempTestSuiteTestDataAsJsonb := protojson.Format(fullTestSuiteMessage.GetTestSuiteTestData())

	// Initiate 'TestSuitePreview' if nil
	if fullTestSuiteMessage.GetTestSuitePreview() == nil {
		fullTestSuiteMessage.TestSuitePreview = &fenixTestCaseBuilderServerGrpcApi.TestSuitePreviewMessage{}

	} else {

		var tempTestSuiteStructureObjects []*fenixTestCaseBuilderServerGrpcApi.TestSuitePreviewStructureMessage_TestSuiteStructureObjectMessage
		var tempSelectedTestSuiteMetaDataValuesMap map[string]*fenixTestCaseBuilderServerGrpcApi.
			TestSuitePreviewStructureMessage_SelectedTestSuiteMetaDataValueMessage
		tempSelectedTestSuiteMetaDataValuesMap = make(map[string]*fenixTestCaseBuilderServerGrpcApi.
			TestSuitePreviewStructureMessage_SelectedTestSuiteMetaDataValueMessage)

		var tempTestSuitePreview *fenixTestCaseBuilderServerGrpcApi.TestSuitePreviewMessage
		tempTestSuitePreview = &fenixTestCaseBuilderServerGrpcApi.TestSuitePreviewMessage{
			TestSuitePreview: &fenixTestCaseBuilderServerGrpcApi.TestSuitePreviewStructureMessage{
				TestSuiteName:                      tempTestSuiteName,
				DomainThatOwnTheTestSuite:          tempDomainUuid,
				TestSuiteDescription:               fullTestSuiteMessage.GetTestSuiteBasicInformation().GetTestSuiteDescription(),
				TestSuiteStructureObjects:          tempTestSuiteStructureObjects,
				TestSuiteUuid:                      tempTestSuiteUuid,
				TestSuiteVersion:                   strconv.Itoa(int(nexTestSuiteVersion)),
				LastSavedByUserOnComputer:          tempInsertedByUserIdOnComputer,
				LastSavedByUserGCPAuthorization:    tempInsertedByGCPAuthenticatedUser,
				LastSavedTimeStamp:                 insertTimeStamp,
				SelectedTestSuiteMetaDataValuesMap: tempSelectedTestSuiteMetaDataValuesMap,
			},
			TestSuitePreviewHash: "",
		}

		// Generate 'SelectedTestSuiteMetaDataValuesMap'
		for _, tempMetaDataGroupMessagePtr := range fullTestSuiteMessage.TestSuiteMetaData.MetaDataGroupsMap {

			// Get the 'MetaDataGroupMessage' from ptr
			tempMetaDataGroupMessage := *tempMetaDataGroupMessagePtr

			// Loop 'MetaDataInGroupMap'
			for SelectedTestSuiteMetaDataValuesMap, tempMetaDataInGroupMessagePtr := range tempMetaDataGroupMessage.MetaDataInGroupMap {

				switch tempMetaDataInGroupMessagePtr.GetSelectType() {

				case fenixTestCaseBuilderServerGrpcApi.MetaDataSelectTypeEnum_MetaDataSelectType_SingleSelect:

					// Create 'SelectedTestSuiteMetaDataValueMessage' to be stored in Map
					var tempSelectedTestSuiteMetaDataValueMessage *fenixTestCaseBuilderServerGrpcApi.
						TestSuitePreviewStructureMessage_SelectedTestSuiteMetaDataValueMessage

					tempSelectedTestSuiteMetaDataValueMessage = &fenixTestCaseBuilderServerGrpcApi.
						TestSuitePreviewStructureMessage_SelectedTestSuiteMetaDataValueMessage{
						OwnerDomainUuid:   tempSelectedTestSuiteMetaDataValueMessage.GetOwnerDomainUuid(),
						OwnerDomainName:   tempSelectedTestSuiteMetaDataValueMessage.GetOwnerDomainName(),
						MetaDataGroupName: tempMetaDataInGroupMessagePtr.GetMetaDataGroupName(),
						MetaDataName:      tempMetaDataInGroupMessagePtr.GetMetaDataName(),
						MetaDataNameValue: tempMetaDataInGroupMessagePtr.GetSelectedMetaDataValueForSingleSelect(),
						SelectType:        tempSelectedTestSuiteMetaDataValueMessage.GetSelectType(),
						IsMandatory:       tempSelectedTestSuiteMetaDataValueMessage.GetIsMandatory(),
					}

					// Add value to map
					tempSelectedTestSuiteMetaDataValuesMap[SelectedTestSuiteMetaDataValuesMap] = tempSelectedTestSuiteMetaDataValueMessage

				case fenixTestCaseBuilderServerGrpcApi.MetaDataSelectTypeEnum_MetaDataSelectType_MultiSelect:

					// Loop all selected values and add to map
					for _, tempSelectedMetaDataValue := range tempMetaDataInGroupMessagePtr.GetSelectedMetaDataValuesForMultiSelect() {

						// Create 'SelectedTestSuiteMetaDataValueMessage' to be stored in Map
						var tempSelectedTestSuiteMetaDataValueMessage *fenixTestCaseBuilderServerGrpcApi.
							TestSuitePreviewStructureMessage_SelectedTestSuiteMetaDataValueMessage

						tempSelectedTestSuiteMetaDataValueMessage = &fenixTestCaseBuilderServerGrpcApi.
							TestSuitePreviewStructureMessage_SelectedTestSuiteMetaDataValueMessage{
							OwnerDomainUuid:   tempSelectedTestSuiteMetaDataValueMessage.GetOwnerDomainUuid(),
							OwnerDomainName:   tempSelectedTestSuiteMetaDataValueMessage.GetOwnerDomainName(),
							MetaDataGroupName: tempMetaDataInGroupMessagePtr.GetMetaDataGroupName(),
							MetaDataName:      tempMetaDataInGroupMessagePtr.GetMetaDataName(),
							MetaDataNameValue: tempSelectedMetaDataValue,
							SelectType:        tempSelectedTestSuiteMetaDataValueMessage.GetSelectType(),
							IsMandatory:       tempSelectedTestSuiteMetaDataValueMessage.GetIsMandatory(),
						}

						// Add value to map
						tempSelectedTestSuiteMetaDataValuesMap[SelectedTestSuiteMetaDataValuesMap] = tempSelectedTestSuiteMetaDataValueMessage

					}

				default:
					common_config.Logger.WithFields(logrus.Fields{
						"Id":                   "6972f543-5fba-414e-b9d2-079a374d0f48",
						"fullTestSuiteMessage": fullTestSuiteMessage,
						"tempMetaDataInGroupMessagePtr.GetSelectType()": tempMetaDataInGroupMessagePtr.GetSelectType(),
					}).Error("Unknown SelectType in 'MetaDataInGroupMap'")

					// Set Error codes to return message
					var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
					var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

					errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
					errorCodes = append(errorCodes, errorCode)

					// Create Return message
					returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
						AckNack:    false,
						Comments:   "Problem when Saving TestSuite to database",
						ErrorCodes: errorCodes,
						ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
							CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
					}

					return returnMessage, err

				}
			}
		}

		// Add 'tempSelectedTestSuiteMetaDataValuesMap' to 'tempTestSuitePreview'
		tempTestSuitePreview.TestSuitePreview.SelectedTestSuiteMetaDataValuesMap = tempSelectedTestSuiteMetaDataValuesMap

		// Calculate 'TestSuitePreviewHash'
		var tempTestSuitePreviewHash string
		tempJson := protojson.Format(tempTestSuitePreview)
		tempTestSuitePreviewHash = common_config.HashSingleValue(tempJson)

		tempTestSuitePreview.TestSuitePreviewHash = tempTestSuitePreviewHash

		// finish Preview-structure to be saved

		// Add TestCases Preview-object to TestSuite
		fullTestSuiteMessage.TestSuitePreview = tempTestSuitePreview

	}

	// Generate json from gRCP-object
	tempTestSuitePreviewAsJsonb := protojson.Format(fullTestSuiteMessage.GetTestSuitePreview())

	// Initiate 'TestSuiteImplementedFunctionsMap' if nil
	if fullTestSuiteMessage.GetTestSuiteImplementedFunctionsMap() == nil {
		fullTestSuiteMessage.TestSuiteImplementedFunctionsMap = make(map[int32]bool)
	}
	tempTestSuiteImplementedFunctionsAsByteArray, err := json.Marshal(fullTestSuiteMessage.GetTestSuiteImplementedFunctionsMap())
	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"Id": "4b297436-3ce2-4c58-93a6-d70cb7da13c6",
			"fullTestSuiteMessage.GetTestSuiteImplementedFunctionsMap()": fullTestSuiteMessage.GetTestSuiteImplementedFunctionsMap(),
		}).Error("Problem generating json-byte-array from 'TestSuiteImplementedFunctionsMap'")

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
	tempTestSuiteImplementedFunctionsAsJsonb := string(tempTestSuiteImplementedFunctionsAsByteArray)

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
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestSuiteDeletedInsertTimeStamp)

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestCasesInTestSuiteAsJsonb)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestSuitePreviewAsJsonb)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestSuiteMetaDataAsJsonb)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestSuiteTestDataAsJsonb)

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, fmt.Sprintf("%d", tempTestSuiteType))
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, fmt.Sprintf("%s", tempTestSuiteTypeName))

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestSuiteDescription)
	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestSuiteExecutionEnvironment)

	dataRowToBeInsertedMultiType = append(dataRowToBeInsertedMultiType, tempTestSuiteImplementedFunctionsAsJsonb)

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

	// A quick fix to be able to use NULL for TimeStamp
	var sqlValuesToInsert string
	sqlValuesToInsert = fenixCloudDBObject.generateSQLInsertValues(dataRowsToBeInsertedMultiType)
	sqlValuesToInsert = strings.ReplaceAll(sqlValuesToInsert, "'NULL'", "NULL")

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "INSERT INTO \"FenixBuilder\".\"TestSuites\" "
	sqlToExecute = sqlToExecute + "(\"DomainUuid\", \"DomainName\", \"TestSuiteUuid\", \"TestSuiteName\", " +
		"\"TestSuiteVersion\", \"TestSuiteHash\", " +
		"\"CanListAndViewTestSuiteAuthorizationLevelOwnedByDomain\", " +
		"\"CanListAndViewTestSuiteAuthorizationLevelHavingTiAndTicWith\", " +
		"\"InsertTimeStamp\", \"InsertedByUserIdOnComputer\", \"InsertedByGCPAuthenticatedUser\", " +
		"\"TestSuiteIsDeleted\", \"DeletedInsertedTimeStamp\", " +
		" \"TestCasesInTestSuite\", \"TestSuitePreview\", \"TestSuiteMetaData\", \"TestSuiteTestData\", " +
		"\"TestSuiteType\", \"TestSuiteTypeName\", \"TestSuiteDescription\", \"TestSuiteExecutionEnvironment\", " +
		"\"TestSuiteImplementedFunctions\") "
	sqlToExecute = sqlToExecute + sqlValuesToInsert
	sqlToExecute = sqlToExecute + ";"

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "7d960b15-a49b-408b-9f6c-f8bb19af1d4e",
			"sqlToExecute": sqlToExecute,
			"err":          err.Error(),
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
		"Id":                       "900b2c42-3c6e-4afa-8605-35e0e07b4e12",
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
	sqlToExecute = sqlToExecute + "SELECT Distinct ON (\"TestCaseUuid\")  \"UniqueCounter\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestCases\" "
	sqlToExecute = sqlToExecute + "WHERE \"TestCaseUuid\" IN " + fenixCloudDBObject.generateSQLINArray(testCasesUuid) + " "
	sqlToExecute = sqlToExecute + "ORDER BY \"TestCaseUuid\", \"UniqueCounter\" DESC "
	sqlToExecute = sqlToExecute + ") "
	sqlToExecute = sqlToExecute + "SELECT t.\"TestInstructionsAsJsonb\", t.\"TestInstructionContainersAsJsonb\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestCases\" t "
	sqlToExecute = sqlToExecute + "WHERE t.\"UniqueCounter\" IN (SELECT * FROM uniquecounters) "
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "5c36b6c4-84bf-4b3a-a851-28bc03106496",
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
			"Id":           "2e2176c5-fc44-412c-a0d3-b8c6298179ac",
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
				"Id":           "fdb86499-bbc6-427e-96df-768e79cfb443",
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
				"Id":    "bba83c43-226b-472f-85d2-41623186b99b",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestInstructionsAsByteArray' into proto-message")

			return nil, nil, err
		}

		// Convert json-byte-arrays into proto-messages - tempTestInstructionContainersAsByteArray
		var tempTestInstructionContainers fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionContainersMessage
		err = protojson.Unmarshal(tempTestInstructionContainersAsByteArray, &tempTestInstructionContainers)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "9ef87c74-b39c-4ef7-b968-afc2dc6d79eb",
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

// Load TestCase's-Preview for TestSuites
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadTestCasesPreviewForTestSuite(
	dbTransaction pgx.Tx,
	testCaseUuids []string) (
	tempTestCasesPreview []*fenixTestCaseBuilderServerGrpcApi.TestCasePreviewMessage,
	err error) {

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT DISTINCT ON (tc.\"TestCaseUuid\") " +
		"tc.\"TestCaseUuid\", tc.\"TestCaseVersion\", tc.\"TestCasePreview\" "
	sqlToExecute = sqlToExecute + "FROM   \"FenixBuilder\".\"TestCases\" t "
	sqlToExecute = sqlToExecute + "WHERE tc.\"TestCaseUuid\" IN  "
	sqlToExecute = sqlToExecute + common_config.GenerateSQLINArray(testCaseUuids) + " "
	sqlToExecute = sqlToExecute + "ORDER  BY tc.\"TestCaseUuid\", tc.\"TestCaseVersion\" DESC "
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "fbb9a0b4-ec7c-4674-a97c-0e3047a976de",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'listTestCasesThatCanBeEdited'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "4ad9875a-4dce-4ae1-b318-e10fc08803c8",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	var (
		tempTestCaseUuid               string
		tempTestCaseVersion            int
		tempTestCasePreviewAsString    string
		tempTestCasePreviewAsByteArray []byte
	)

	// Extract data from DB result set
	for rows.Next() {

		var tempTestCasePreview fenixTestCaseBuilderServerGrpcApi.TestCasePreviewMessage

		err = rows.Scan(
			&tempTestCaseUuid,
			&tempTestCaseVersion,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "d61ef5bb-f6e2-4571-86a1-6a4383e640c4",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert json-string into byte-array
		tempTestCasePreviewAsByteArray = []byte(tempTestCasePreviewAsString)

		// Convert json-byte-arrays into proto-messages
		err = protojson.Unmarshal(tempTestCasePreviewAsByteArray, &tempTestCasePreview)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "8e4e83a1-a8d1-41ae-9c5f-f7091ece8b15",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestCasePreviewAsByteArray' into proto-message")

			return nil, err
		}

		// Add to slice of TestCases
		tempTestCasesPreview = append(tempTestCasesPreview, &tempTestCasePreview)

	}

	return tempTestCasesPreview, err

}
