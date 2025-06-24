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

// This date is used as delete date when the TestSuite is not deleted
const testSuiteNotDeletedDate = "2068-11-18"

// PrepareLoadFullTestSuite
// Load Full TestSuite from Database
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareLoadFullTestSuite(
	testSuiteUuidToLoad string,
	gCPAuthenticatedUser string) (
	responseMessage *fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "b9867b4a-dbf1-4d17-ad80-04265d1d4964",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'prepareSaveFullTestSuite'")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "Problem when Loading from database",
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
			AckNackResponse:   ackNackResponse,
			DetailedTestSuite: nil,
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
			"id":                   "4d865b21-b923-438c-bd63-250034583c88",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Warning("User doesn't have access to any domains")

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack: false,
			Comments: fmt.Sprintf("User %s doesn't have access to any domains",
				gCPAuthenticatedUser),
			ErrorCodes: nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
			AckNackResponse:   ackNackResponse,
			DetailedTestSuite: nil,
		}

		return responseMessage

	}

	// Load the TestSuite
	var fullTestSuiteMessage *fenixTestCaseBuilderServerGrpcApi.FullTestSuiteMessage
	fullTestSuiteMessage, err = fenixCloudDBObject.loadFullTestSuite(txn, testSuiteUuidToLoad, domainAndAuthorizations)

	// Error when retrieving TestSuite
	if err != nil {
		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "Problem when Loading from database",
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
			AckNackResponse:   ackNackResponse,
			DetailedTestSuite: nil,
		}

		return responseMessage
	}

	// TestCase
	if fullTestSuiteMessage == nil {

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "TestCase couldn't be found in Database or the user doesn't have access to the TestSuite",
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
			AckNackResponse:   ackNackResponse,
			DetailedTestSuite: nil,
		}

		return responseMessage
	}

	// Load Users TestData
	err = fenixCloudDBObject.loadUsersTestDataForTestSuite(
		txn,
		gCPAuthenticatedUser,
		fullTestSuiteMessage)

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
			AckNack:    false,
			Comments:   "Problem when Loading Users Testdata from database",
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
			AckNackResponse:   ackNackResponse,
			DetailedTestSuite: nil,
		}

		return responseMessage
	}

	// Create response message
	var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
	ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:    true,
		Comments:   "",
		ErrorCodes: nil,
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
			CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
	}

	responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
		AckNackResponse:   ackNackResponse,
		DetailedTestSuite: fullTestSuiteMessage,
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

// Load the TestSuite with all its data
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadFullTestSuite(
	dbTransaction pgx.Tx,
	testSuiteUuidToLoad string,
	domainAndAuthorizations []DomainAndAuthorizationsStruct) (
	fullTestSuiteMessage *fenixTestCaseBuilderServerGrpcApi.FullTestSuiteMessage, err error) {

	// Generate a Domains list and Calculate the Authorization requirements
	var tempCalculatedDomainAndAuthorizations DomainAndAuthorizationsStruct
	var domainList []string
	for _, domainAndAuthorization := range domainAndAuthorizations {
		// Add to DomainList
		domainList = append(domainList, domainAndAuthorization.DomainUuid)

		// Calculate the Authorization requirements for...
		// TestCaseAuthorizationLevelOwnedByDomain
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteOwnedByThisDomain =
			tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteOwnedByThisDomain +
				domainAndAuthorization.CanListAndViewTestCaseOrTestSuiteOwnedByThisDomain

		// TestCaseAuthorizationLevelHavingTiAndTicWithDomain
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain =
			tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain +
				domainAndAuthorization.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain
	}

	// Convert Values into string for TestSuiteAuthorizationLevelOwnedByDomain
	var tempCanListAndViewTestSuiteOwnedByThisDomainAsString string
	tempCanListAndViewTestSuiteOwnedByThisDomainAsString = strconv.FormatInt(
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain, 10)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":    "04ce498d-d440-4eaf-8f0a-8ade191a47ba",
			"Error": err,
			"tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain": tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain,
		}).Error("Couldn't convert into string representation")

		return nil, err
	}

	// Convert Values into string for TestSuiteAuthorizationLevelOwnedByDomain
	var tempCanListAndViewTestSuiteHavingTIandTICfromThisDomainAsString string
	tempCanListAndViewTestSuiteHavingTIandTICfromThisDomainAsString = strconv.FormatInt(
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain, 10)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":    "9f094b06-34e4-46d0-bb25-1f2bbedb1d2f",
			"Error": err,
			"tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain": tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain,
		}).Error("Couldn't convert into string representation")

		return nil, err
	}

	/*
		sqlToExecute = sqlToExecute + 	"SELECT TS.\"DomainUuid\", TS.\"DomainName\", TS.\"TestSuiteUuid\", TS.\"TestSuiteName\", " +
		sqlToExecute = sqlToExecute + 				"TS.\"TestSuiteVersion\", TS.\"TestSuiteHash\", " +
		sqlToExecute = sqlToExecute + 			"TS.\"InsertTimeStamp\", TS.\"InsertedByUserIdOnComputer\", TS.\"InsertedByGCPAuthenticatedUser\", " +
		sqlToExecute = sqlToExecute + 			"TS.\"TestCasesInTestSuite\", TS.\"TestSuitePreview\", TS.\"TestSuiteMetaData\" "

		WITH uniquecounters AS (
		    SELECT Distinct ON ("TestSuiteUuid")  "UniqueCounter"
		    FROM "FenixBuilder"."TestSuites"
		    WHERE "TestSuiteUuid" =  '3e903bc2-57bc-4083-bc0d-184ddb6522d4'


		    ORDER BY "TestSuiteUuid", "UniqueCounter" DESC
		)
		SELECT t.*
		FROM "FenixBuilder"."TestSuites" t
		WHERE t."UniqueCounter" IN (SELECT * FROM uniquecounters) AND
		      t."TestSuiteIsDeleted" = false;

	*/

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "WITH uniquecounters AS ( "
	sqlToExecute = sqlToExecute + "SELECT Distinct ON (\"TestSuiteUuid\")  \"UniqueCounter\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestSuites\" "
	sqlToExecute = sqlToExecute + fmt.Sprintf("WHERE \"TestSuiteUuid\" = '%s' ", testSuiteUuidToLoad)
	sqlToExecute = sqlToExecute + "ORDER BY \"TestSuiteUuid\", \"UniqueCounter\" DESC "
	sqlToExecute = sqlToExecute + ") "

	sqlToExecute = sqlToExecute + "SELECT TS.\"DomainUuid\", TS.\"DomainName\", TS.\"TestSuiteUuid\", TS.\"TestSuiteName\", " +
		"TS.\"TestSuiteVersion\", TS.\"TestSuiteHash\", TS.\"DeleteTimestamp\"" +
		"TS.\"InsertTimeStamp\", TS.\"InsertedByUserIdOnComputer\", TS.\"InsertedByGCPAuthenticatedUser\", " +
		"TS.\"TestCasesInTestSuite\", TS.\"TestSuitePreview\", TS.\"TestSuiteMetaData\" "

	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestSuites\" TS "
	sqlToExecute = sqlToExecute + fmt.Sprintf("WHERE TS.\"TestSuiteUuid\" = '%s' ", testSuiteUuidToLoad)
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "(TS.\"CanListAndViewTestSuiteAuthorizationLevelOwnedByDomain\" & " +
		tempCanListAndViewTestSuiteOwnedByThisDomainAsString + ")"
	sqlToExecute = sqlToExecute + "= TS.\"CanListAndViewTestCaseAuthorizationLevelOwnedByDomain\" "
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "(TS.\"CanListAndViewTestSuiteAuthorizationLevelHavingTiAndTicWith\" & " +
		tempCanListAndViewTestSuiteHavingTIandTICfromThisDomainAsString + ")"
	sqlToExecute = sqlToExecute + "= TS.\"CanListAndViewTestSuiteAuthorizationLevelHavingTiAndTicWith\" "
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "TS.\"UniqueCounter\" IN (SELECT * FROM uniquecounters) AND "
	sqlToExecute = sqlToExecute + "TS.\"TestSuiteIsDeleted\" = false "
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "3b331ac4-7944-4037-ba1a-aa8cb4d38074",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadFullTestSuite'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "c66b4f71-8b00-482f-82d2-a8f1e622f08b",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	var (
		tempDomainUuid                     string
		tempDomainName                     string
		tempTestSuiteUuid                  string
		tempTestSuiteName                  string
		tempTestSuiteVersion               int
		tempTestSuiteHash                  string
		tempDeleteTimestamp                string
		tempInsertTimeStampAsString        string
		tempInsertedByUserIdOnComputer     string
		tempInsertedByGCPAuthenticatedUser string
		tempTestCasesInTestSuiteAsJson     string
		tempTestSuitePreviewAsJson         string
		tempTestSuiteMetaDataAsJson        string

		tempTestCasesInTestSuiteAsByteArray []byte
		tempTestSuitePreviewAsByteArray     []byte
		tempTestSuiteMetaDataAsByteArray    []byte

		tempTestCasesInTestSuiteAsGrpc fenixTestCaseBuilderServerGrpcApi.TestCasesInTestSuiteMessage
		tempTestSuitePreviewAsGrpc     fenixTestCaseBuilderServerGrpcApi.TestSuitePreviewMessage
		tempTestSuiteMetaDataAsGrpc    fenixTestCaseBuilderServerGrpcApi.UserSpecifiedTestSuiteMetaDataMessage
	)

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&tempDomainUuid,
			&tempDomainName,
			&tempTestSuiteUuid,
			&tempTestSuiteName,
			&tempTestSuiteVersion,
			&tempTestSuiteHash,
			&tempDeleteTimestamp,
			&tempInsertTimeStampAsString,
			&tempInsertedByUserIdOnComputer,
			&tempInsertedByGCPAuthenticatedUser,
			&tempTestCasesInTestSuiteAsJson,
			&tempTestSuitePreviewAsJson,
			&tempTestSuiteMetaDataAsJson,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "89c7db75-5ff1-4bf3-a1bb-5470520af682",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert json-strings into byte-arrays
		tempTestCasesInTestSuiteAsByteArray = []byte(tempTestCasesInTestSuiteAsJson)
		tempTestSuitePreviewAsByteArray = []byte(tempTestSuitePreviewAsJson)
		tempTestSuiteMetaDataAsByteArray = []byte(tempTestSuiteMetaDataAsJson)

		// Convert json-byte-arrays into proto-messages
		err = protojson.Unmarshal(tempTestCasesInTestSuiteAsByteArray, &tempTestCasesInTestSuiteAsGrpc)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "2e73987c-1693-4673-b366-0b56efcfbc09",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestCasesInTestSuiteAsByteArray' into proto-message")

			return nil, err
		}

		err = protojson.Unmarshal(tempTestSuitePreviewAsByteArray, &tempTestSuitePreviewAsGrpc)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "eafb06bb-1257-4c49-b57c-99edb3baa81d",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestSuitePreviewAsByteArray' into proto-message")

			return nil, err
		}

		err = protojson.Unmarshal(tempTestSuiteMetaDataAsByteArray, &tempTestSuiteMetaDataAsGrpc)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "4643c3d5-87bf-4a5f-a171-b01a9d3c0d4b",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestSuiteMetaDataAsByteArray' into proto-message")

			return nil, err
		}

		// Add the different parts into full TestSuite-message
		fullTestSuiteMessage = &fenixTestCaseBuilderServerGrpcApi.FullTestSuiteMessage{
			TestSuiteBasicInformation: &fenixTestCaseBuilderServerGrpcApi.TestSuiteBasicInformationMessage{
				DomainUuid:           tempDomainUuid,
				DomainName:           tempDomainName,
				TestSuiteUuid:        tempTestSuiteUuid,
				TestSuiteVersion:     uint32(tempTestSuiteVersion),
				TestSuiteName:        tempTestSuiteName,
				TestSuiteDescription: tempTestSuiteName,
			},
			TestSuiteTestData:    nil,
			TestSuitePreview:     &tempTestSuitePreviewAsGrpc,
			TestSuiteMetaData:    &tempTestSuiteMetaDataAsGrpc,
			TestCasesInTestSuite: &tempTestCasesInTestSuiteAsGrpc,
			DeletedDate:          d,
			SupportedTestSuiteDataToBeStoredAsJsonString: "",
			MessageHash: "",
		}

	}

	return fullTestSuiteMessage, err

}

/*
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


*/
