package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"strconv"
	"time"
)

// PrepareLoadFullTestCase
// Load Full TestCase from Database
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareLoadFullTestCase(testCaseUuidToLoad string, gCPAuthenticatedUser string) (responseMessage *fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse) {

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

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "Problem when Loading from database",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse{
			AckNackResponse:  ackNackResponse,
			DetailedTestCase: nil,
		}

		return responseMessage
	}

	defer txn.Commit(context.Background())

	// Load Domains that User has access to
	var domainAndAuthorizations []DomainAndAuthorizationsStruct
	domainAndAuthorizations, err = fenixCloudDBObject.PrepareLoadUsersDomains(gCPAuthenticatedUser)

	// If user doesn't have access to any domains then exit with warning in log
	if len(domainAndAuthorizations) == 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "02a762dc-932b-4edd-a8b5-a6d0a53ba36b",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Warning("User doesn't have access to any domains")

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     fmt.Sprintf("User %s doesn't have access to any domains", gCPAuthenticatedUser),
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse{
			AckNackResponse:  ackNackResponse,
			DetailedTestCase: nil,
		}

		return responseMessage

	}

	// Load the TestCase
	var fullTestCaseMessage *fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage
	fullTestCaseMessage, err = fenixCloudDBObject.loadFullTestCase(txn, testCaseUuidToLoad, domainAndAuthorizations)

	// Error when retrieving TestCase
	if err != nil {
		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "Problem when Loading from database",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse{
			AckNackResponse:  ackNackResponse,
			DetailedTestCase: nil,
		}

		return responseMessage
	}

	// TestCase
	if fullTestCaseMessage == nil {

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "TestCase couldn't be found in Database or the user doesn't have access to the TestCase",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse{
			AckNackResponse:  ackNackResponse,
			DetailedTestCase: nil,
		}

		return responseMessage
	}

	// Load Users TestData
	err = fenixCloudDBObject.loadUsersTestDataForTestCase(
		txn,
		gCPAuthenticatedUser,
		fullTestCaseMessage)

	// Error when retrieving Users TestData
	if err != nil {
		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "Problem when Loading Users Testdata from database",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse{
			AckNackResponse:  ackNackResponse,
			DetailedTestCase: nil,
		}

		return responseMessage
	}

	// Create response message
	var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
	ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
	}

	responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse{
		AckNackResponse:  ackNackResponse,
		DetailedTestCase: fullTestCaseMessage,
	}

	return responseMessage
}

/*
SELECT TC."TestCaseBasicInformationAsJsonb", TC."TestInstructionsAsJsonb", "TestInstructionContainersAsJsonb"
FROM "FenixBuilder"."TestCases" TC
WHERE TC."TestCaseUuid" = '1f969ca4-e279-431a-b588-491f6f62d41e'
ORDER BY TC."TestCaseVersion" DESC
LIMIT 1;

*/

// Load All Domains and their address information
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadFullTestCase(
	dbTransaction pgx.Tx,
	testCaseUuidToLoad string,
	domainAndAuthorizations []DomainAndAuthorizationsStruct) (
	fullTestCaseMessage *fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage, err error) {

	// Generate a Domains list and Calculate the Authorization requirements
	var tempCalculatedDomainAndAuthorizations DomainAndAuthorizationsStruct
	var domainList []string
	for _, domainAndAuthorization := range domainAndAuthorizations {
		// Add to DomainList
		domainList = append(domainList, domainAndAuthorization.DomainUuid)

		// Calculate the Authorization requirements for...
		// TestCaseAuthorizationLevelOwnedByDomain
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOwnedByThisDomain =
			tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOwnedByThisDomain +
				domainAndAuthorization.CanListAndViewTestCaseOwnedByThisDomain

		// TestCaseAuthorizationLevelHavingTiAndTicWithDomain
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain =
			tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain +
				domainAndAuthorization.CanListAndViewTestCaseHavingTIandTICFromThisDomain
	}

	// Convert Values into string for TestCaseAuthorizationLevelOwnedByDomain
	var tempCanListAndViewTestCaseOwnedByThisDomainAsString string
	tempCanListAndViewTestCaseOwnedByThisDomainAsString = strconv.FormatInt(
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain, 10)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":    "938e364a-5ab7-4248-bb39-217b283407b8",
			"Error": err,
			"tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain": tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain,
		}).Error("Couldn't convert into string representation")

		return nil, err
	}

	// Convert Values into string for TestCaseAuthorizationLevelOwnedByDomain
	var tempCanListAndViewTestCaseHavingTIandTICfromThisDomainAsString string
	tempCanListAndViewTestCaseHavingTIandTICfromThisDomainAsString = strconv.FormatInt(
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain, 10)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":    "938e364a-5ab7-4248-bb39-217b283407b8",
			"Error": err,
			"tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain": tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain,
		}).Error("Couldn't convert into string representation")

		return nil, err
	}

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT TC.\"TestCaseBasicInformationAsJsonb\", " +
		"TC.\"TestInstructionsAsJsonb\", \"TestInstructionContainersAsJsonb\"," +
		"TC.\"TestCaseHash\", TC.\"TestCaseExtraInformationAsJsonb\", TC.\"TestCaseTemplateFilesAsJsonb\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestCases\" TC "
	sqlToExecute = sqlToExecute + fmt.Sprintf("WHERE TC.\"TestCaseUuid\" = '%s' ", testCaseUuidToLoad)
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "(TC.\"CanListAndViewTestCaseAuthorizationLevelOwnedByDomain\" & " + tempCanListAndViewTestCaseOwnedByThisDomainAsString + ")"
	sqlToExecute = sqlToExecute + "= TC.\"CanListAndViewTestCaseAuthorizationLevelOwnedByDomain\" "
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "(TC.\"CanListAndViewTestCaseAuthorizationLevelHavingTiAndTicWithDomain\" & " + tempCanListAndViewTestCaseHavingTIandTICfromThisDomainAsString + ")"
	sqlToExecute = sqlToExecute + "= TC.\"CanListAndViewTestCaseAuthorizationLevelHavingTiAndTicWithDomain\" "
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "TC.\"TestCaseVersion\" = (SELECT MAX(\"TestCaseVersion\") "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestCases\" tc2 "
	sqlToExecute = sqlToExecute + fmt.Sprintf("WHERE tc2.\"TestCaseUuid\" = '%s' ", testCaseUuidToLoad)
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "tc2.\"TestCaseIsDeleted\" = false )"
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "01b246fb-effe-4348-9a5c-830604e6daf6",
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
			"Id":           "784c6f8d-fd77-44e0-9f2b-17e8438ad749",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	var (
		tempTestCaseBasicInformation             fenixTestCaseBuilderServerGrpcApi.TestCaseBasicInformationMessage
		tempMatureTestInstructions               fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionsMessage
		tempMatureTestInstructionContainers      fenixTestCaseBuilderServerGrpcApi.MatureTestInstructionContainersMessage
		tempTestCaseExtraInformation             fenixTestCaseBuilderServerGrpcApi.TestCaseExtraInformationMessage
		tempTestCaseTemplateFilesMessage         fenixTestCaseBuilderServerGrpcApi.TestCaseTemplateFilesMessage
		tempTestCaseBasicInformationAsString     string
		tempTestInstructionsAsString             string
		tempTestInstructionContainersAsString    string
		tempTestCaseExtraInformationAsString     string
		tempTestCaseTemplateFilesAsString        string
		tempTestCaseBasicInformationAsByteArray  []byte
		tempTestInstructionsAsByteArray          []byte
		tempTestInstructionContainersAsByteArray []byte
		tempTestCaseExtraInformationAsByteArray  []byte
		tempTestCaseTemplateFilesAsByteArray     []byte

		testCaseHash string
	)

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&tempTestCaseBasicInformationAsString,
			&tempTestInstructionsAsString,
			&tempTestInstructionContainersAsString,
			&testCaseHash,
			&tempTestCaseExtraInformationAsString,
			&tempTestCaseTemplateFilesAsString,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "2a81dfda-4937-4d9e-9827-7191eb7ac7de",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert json-strings into byte-arrays
		tempTestCaseBasicInformationAsByteArray = []byte(tempTestCaseBasicInformationAsString)
		tempTestInstructionsAsByteArray = []byte(tempTestInstructionsAsString)
		tempTestInstructionContainersAsByteArray = []byte(tempTestInstructionContainersAsString)
		tempTestCaseExtraInformationAsByteArray = []byte(tempTestCaseExtraInformationAsString)
		tempTestCaseTemplateFilesAsByteArray = []byte(tempTestCaseTemplateFilesAsString)

		// Convert json-byte-arrays into proto-messages
		err = protojson.Unmarshal(tempTestCaseBasicInformationAsByteArray, &tempTestCaseBasicInformation)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "d315ea2b-8263-4ad8-9b96-d62da4acf35f",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestCaseBasicInformationAsByteArray' into proto-message")

			return nil, err
		}

		err = protojson.Unmarshal(tempTestInstructionsAsByteArray, &tempMatureTestInstructions)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "441a35b3-5139-4046-8aeb-a986a84827df",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestInstructionsAsByteArray' into proto-message")

			return nil, err
		}

		err = protojson.Unmarshal(tempTestInstructionContainersAsByteArray, &tempMatureTestInstructionContainers)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "f4738b27-4c49-448b-b49b-f6cf08508f12",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestInstructionContainersAsByteArray' into proto-message")

			return nil, err
		}

		err = protojson.Unmarshal(tempTestCaseExtraInformationAsByteArray, &tempTestCaseExtraInformation)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "3b02c709-cb70-43c6-982e-19eb58395a24",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestCaseExtraInformationAsByteArray' into proto-message")

			return nil, err
		}

		err = protojson.Unmarshal(tempTestCaseTemplateFilesAsByteArray, &tempTestCaseTemplateFilesMessage)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "f262f8cb-16ea-4e12-8f25-ab73306fdeba",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestCaseTemplateFilesAsByteArray' into proto-message")

			return nil, err
		}

		// Add the different parts into full TestCase-message
		fullTestCaseMessage = &fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage{
			TestCaseBasicInformation:        &tempTestCaseBasicInformation,
			MatureTestInstructions:          &tempMatureTestInstructions,
			MatureTestInstructionContainers: &tempMatureTestInstructionContainers,
			MessageHash:                     testCaseHash,
			TestCaseExtraInformation:        &tempTestCaseExtraInformation,
			TestCaseTemplateFiles:           &tempTestCaseTemplateFilesMessage,
		}

	}

	return fullTestCaseMessage, err

}

// Load Users saved TestData for a specific TestCase
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadUsersTestDataForTestCase(
	dbTransaction pgx.Tx,
	gCPAuthenticatedUser string,
	fullTestCaseMessage *fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage) (
	err error) {

	testCaseUuid := fullTestCaseMessage.GetTestCaseBasicInformation().GetBasicTestCaseInformation().GetNonEditableInformation().TestCaseUuid
	testCaseVersion := fullTestCaseMessage.GetTestCaseBasicInformation().GetBasicTestCaseInformation().GetNonEditableInformation().TestCaseVersion

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT utdftc.\"TestData\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"UsersTestDataForTestCase\" utdftc "
	sqlToExecute = sqlToExecute + fmt.Sprintf("WHERE utdftc.\"TestCaseUuid\" = '%s' ", testCaseUuid)
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + fmt.Sprintf("utdftc.\"TestCaseVersion\" = %d ", testCaseVersion)
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + fmt.Sprintf("utdftc.\"GcpAuthenticatedUser\" = '%s' ", gCPAuthenticatedUser)
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
			"Id":           "a09022fe-e880-45f8-ba70-5c893a97e38e",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	var (
		tempTestDataAsString    string
		tempTestDataAsByteArray []byte
	)

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&tempTestDataAsString,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "c562fc51-20a9-4d5e-aa04-7653469dd399",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// Convert json-strings into byte-arrays
		tempTestDataAsByteArray = []byte(tempTestDataAsString)

		// Convert json-byte-arrays into proto-messages
		var testData fenixTestCaseBuilderServerGrpcApi.UsersChosenTestDataForTestCaseMessage
		err = protojson.Unmarshal(tempTestDataAsByteArray, &testData)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "a0378d91-3d7f-4de9-ae4b-79f5b709cb39",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestDataAsByteArray' into proto-message")

			return err
		}

		// Store TestData into full TestCase-message
		fullTestCaseMessage.TestCaseTestData = &testData

	}

	return err

}
