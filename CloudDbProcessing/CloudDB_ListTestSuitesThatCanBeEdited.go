package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

// PrepareListTestSuitesThatCanBeEdited
// List all TestSuites from Database that the user can edit
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareListTestSuitesThatCanBeEdited(
	gCPAuthenticatedUser string,
	testSuiteUpdatedMinTimeStamp time.Time,
	testSuiteExecutionUpdatedMinTimeStamp time.Time) (
	responseMessage *fenixTestCaseBuilderServerGrpcApi.ListTestSuitesResponseMessage) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "e19adf43-69f5-4b12-a975-d95a2a351d11",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'PrepareListTestSuitesThatCanBeEdited'")

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

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestSuitesResponseMessage{
			AckNackResponse:           ackNackResponse,
			BasicTestSuiteInformation: nil,
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
			"id":                   "71af59dc-68d9-41c7-80c8-e0eae1ddede4",
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

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestSuitesResponseMessage{
			AckNackResponse:           ackNackResponse,
			BasicTestSuiteInformation: nil,
		}

		return responseMessage

	}

	// Load the TestSuites
	var testSuitesThatCanBeEditedResponse []*fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage
	testSuitesThatCanBeEditedResponse, err = fenixCloudDBObject.listTestSuitesThatCanBeEdited(
		txn,
		domainAndAuthorizations,
		testSuiteUpdatedMinTimeStamp)

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
			AckNack:                      false,
			Comments:                     "Problem when Loading from database",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestSuitesResponseMessage{
			AckNackResponse:           ackNackResponse,
			BasicTestSuiteInformation: nil,
		}

		return responseMessage
	}

	// No TestSuites
	if testSuitesThatCanBeEditedResponse == nil {

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "TestSuites couldn't be found in Database or the user doesn't have access to the TestSuites",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestSuitesResponseMessage{
			AckNackResponse:           ackNackResponse,
			BasicTestSuiteInformation: nil,
		}

		return responseMessage
	}

	// Create a slice with all TestSuiteUuid to be used for finding execution status
	var testSuiteUuidSlice []string

	for _, tempTestSuite := range testSuitesThatCanBeEditedResponse {
		testSuiteUuidSlice = append(testSuiteUuidSlice, tempTestSuite.NonEditableInformation.GetTestSuiteUuid())
	}

	// Load the latest Execution Status for TestSuites
	var testSuitesLatestExecutionStatusMap map[string]*fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage
	testSuitesLatestExecutionStatusMap = make(map[string]*fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage)
	testSuitesLatestExecutionStatusMap, err = fenixCloudDBObject.loadLatestExecutionStatusForTestSuites(
		txn,
		testSuiteUuidSlice,
		testSuiteExecutionUpdatedMinTimeStamp)

	// Error when retrieving TestSuiteExecution-status
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
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(
				common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestSuitesResponseMessage{
			AckNackResponse:           ackNackResponse,
			BasicTestSuiteInformation: nil,
		}

		return responseMessage
	}

	// Load the latest OK Execution Status for TestCase
	var testSuitesLatestFinishedOkExecutionStatusMap map[string]*fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage
	testSuitesLatestFinishedOkExecutionStatusMap = make(map[string]*fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage)
	testSuitesLatestFinishedOkExecutionStatusMap, err = fenixCloudDBObject.loadLatestFinishedOkExecutionStatusForTestSuites(
		txn,
		testSuiteUuidSlice)

	// Error when retrieving TestCaseExecution-status
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
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(
				common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestSuitesResponseMessage{
			AckNackResponse:           ackNackResponse,
			BasicTestSuiteInformation: nil,
		}

		return responseMessage
	}

	// Merge Execution status into full TestSuiteList
	var foundInMap bool
	var changesAreMade bool
	for testCaseIndex, tempTestSuite := range testSuitesThatCanBeEditedResponse {

		// Reset 'changesAreMade'
		changesAreMade = false

		// Latest Execution Status Information
		var tempTestSuiteFromStatus *fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage
		tempTestSuiteFromStatus, foundInMap = testSuitesLatestExecutionStatusMap[tempTestSuite.NonEditableInformation.GetTestSuiteUuid()]

		// TestSuiteExecution-status wasn't found in Map which indicates that there are no executions for the TestSuite
		if foundInMap == false {

		} else {
			// Add Latest Status information
			tempTestSuite.LatestTestSuiteExecutionStatus = tempTestSuiteFromStatus.LatestTestSuiteExecutionStatus
			tempTestSuite.LatestTestSuiteExecutionStatusInsertTimeStamp = tempTestSuiteFromStatus.
				LatestTestSuiteExecutionStatusInsertTimeStamp

			// Indicate tha changes are done
			changesAreMade = true

		}

		// Latest Finished OK Execution Status Information
		var tempTestSuiteFromFinishedStatus *fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage
		tempTestSuiteFromFinishedStatus, foundInMap = testSuitesLatestFinishedOkExecutionStatusMap[tempTestSuite.NonEditableInformation.GetTestSuiteUuid()]

		// TestCaseExecution-status wasn't found in Map which indicates that there are no Finished OK executions for the TestSuite
		if foundInMap == false {

		} else {
			// Add Latest Finished OK Status information
			tempTestSuite.LatestFinishedOkTestSuiteExecutionStatusInsertTimeStamp = tempTestSuiteFromFinishedStatus.
				LatestTestSuiteExecutionStatusInsertTimeStamp

			// Indicate tha changes are done
			changesAreMade = true

		}

		// Save back the TestSuite into the Slice when changes are done
		if changesAreMade == true {
			testSuitesThatCanBeEditedResponse[testCaseIndex] = tempTestSuite
		}

	}

	// Create response message
	var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
	ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
	}

	responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestSuitesResponseMessage{
		AckNackResponse:           ackNackResponse,
		BasicTestSuiteInformation: testSuitesThatCanBeEditedResponse,
	}

	return responseMessage
}

// Load all TestSuites that the user can edit
func (fenixCloudDBObject *FenixCloudDBObjectStruct) listTestSuitesThatCanBeEdited(
	dbTransaction pgx.Tx,
	domainAndAuthorizations []DomainAndAuthorizationsStruct,
	testSuiteUpdatedMinTimeStamp time.Time) (
	testSuitesThatCanBeEditedByUser []*fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage,
	err error) {

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
			"Id":    "73b8d307-001f-472e-80b3-1105284e9b6d",
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
			"Id":    "40cf39b1-6277-4958-8f3d-3e3785db3290",
			"Error": err,
			"tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain": tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain,
		}).Error("Couldn't convert into string representation")

		return nil, err
	}

	// Delete timestamp
	var deleteTimeStampAsString string
	deleteTimeStampAsString = time.Now().Format("2006-01-02 00:00:00")

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT ts1.\"DomainUuid\", ts1.\"DomainName\", ts1.\"TestSuiteUuid\", " +
		"ts1.\"TestSuiteName\", ts1.\"TestSuiteVersion\", ts1.\"InsertTimeStamp\",  ts1.\"TestSuiteExecutionEnvironment\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestSuites\" ts1 "
	sqlToExecute = sqlToExecute + "WHERE (ts1.\"CanListAndViewTestSuiteAuthorizationLevelOwnedByDomain\" & " + tempCanListAndViewTestSuiteOwnedByThisDomainAsString + ")"
	sqlToExecute = sqlToExecute + "= ts1.\"CanListAndViewTestSuiteAuthorizationLevelOwnedByDomain\" "
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "(ts1.\"CanListAndViewTestSuiteAuthorizationLevelHavingTiAndTicWith\" & " + tempCanListAndViewTestSuiteHavingTIandTICfromThisDomainAsString + ")"
	sqlToExecute = sqlToExecute + "= ts1.\"CanListAndViewTestSuiteAuthorizationLevelHavingTiAndTicWith\" "
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "ts1.\"InsertTimeStamp\" IS NOT NULL AND " +
		"ts1.\"TestSuiteVersion\" = (" +
		"SELECT MAX(ts2.\"TestSuiteVersion\") " +
		"FROM \"FenixBuilder\".\"TestSuites\" ts2 " +
		"WHERE ts2.\"TestSuiteUuid\" = ts1.\"TestSuiteUuid\") AND "
	sqlToExecute = sqlToExecute + "ts1.\"InsertTimeStamp\" > '" +
		common_config.GenerateDatetimeFromTimeInputForDB(testSuiteUpdatedMinTimeStamp) + "' AND "
	sqlToExecute = sqlToExecute + "ts1.\"DeleteTimestamp\" > '" + deleteTimeStampAsString + "' "
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "446f21ba-22c8-4499-ac38-81118f2c7476",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'listTestSuitesThatCanBeEdited'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "fe8b2ba0-7960-4639-8a8d-f0f21a2a7826",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	var (
		tempDomainUuid                    string
		tempDomainName                    string
		tempTestSuiteUuid                 string
		tempTestSuiteName                 string
		tempTestSuiteVersion              int
		tempInsertTimeStampAsTimeStamp    time.Time
		tempTestSuiteExecutionEnvironment string
	)

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&tempDomainUuid,
			&tempDomainName,
			&tempTestSuiteUuid,
			&tempTestSuiteName,
			&tempTestSuiteVersion,
			&tempInsertTimeStampAsTimeStamp,
			&tempTestSuiteExecutionEnvironment,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "0402e085-9d27-455d-a3b4-52c90b3d43b8",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Create 'TestSuiteThatCanBeEditedByUser'-object
		var tempTestSuiteThatCanBeEditedByUser fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage
		tempTestSuiteThatCanBeEditedByUser = fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage{
			NonEditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage_NonEditableBasicInformationMessage{
				TestSuiteUuid:                 tempTestSuiteUuid,
				DomainUuid:                    tempDomainUuid,
				DomainName:                    tempDomainName,
				TestSuiteVersion:              uint32(tempTestSuiteVersion),
				TestSuiteExecutionEnvironment: tempTestSuiteExecutionEnvironment,
			},
			EditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage_EditableBasicInformationMessage{
				TestSuiteName:        tempTestSuiteName,
				TestSuiteDescription: "",
			},
			LatestTestSuiteExecutionStatus:                          0,
			LatestTestSuiteExecutionStatusInsertTimeStamp:           nil,
			LatestFinishedOkTestSuiteExecutionStatusInsertTimeStamp: nil,
			LastSavedTimeStamp:                                      timestamppb.New(tempInsertTimeStampAsTimeStamp),
			TestSuitePreview: &fenixTestCaseBuilderServerGrpcApi.
				TestSuitePreviewMessage{
				TestSuitePreview: &fenixTestCaseBuilderServerGrpcApi.TestSuitePreviewStructureMessage{
					TestSuiteUuid:                 "",
					TestSuiteName:                 "",
					TestSuiteVersion:              "",
					DomainUuidThatOwnTheTestSuite: "",
					DomainNameThatOwnTheTestSuite: "",
					TestSuiteDescription:          "",
					TestSuiteStructureObjects: &fenixTestCaseBuilderServerGrpcApi.
						TestSuitePreviewStructureMessage_TestSuiteStructureObjectMessage{},
					LastSavedByUserOnComputer:          "",
					LastSavedByUserGCPAuthorization:    "",
					LastSavedTimeStamp:                 "",
					SelectedTestSuiteMetaDataValuesMap: nil,
				},
				TestSuitePreviewHash: "",
			},
		}

		// Add to slice of TestCases
		testSuitesThatCanBeEditedByUser = append(testSuitesThatCanBeEditedByUser, &tempTestSuiteThatCanBeEditedByUser)

	}

	return testSuitesThatCanBeEditedByUser, err

}

/*
-- Latest Execution
SELECT tce1."TestCaseUuid", tce1."TestCaseVersion", tce1."TestCaseExecutionStatus", tce1."ExecutionStatusUpdateTimeStamp"
FROM "FenixExecution"."TestCasesUnderExecution" tce1
WHERE tce1."TestCaseUuid" IN ('4eebed04-39a9-4ad9-ae67-51c9de984486',  '653a43f7-2a7f-445b-aa3f-8cb4d1797369',  'df92530f-3a43-4e3d-9578-f300d996e2e5',  '2d4128e9-04f3-4615-aaf2-330193fcdbd7',  '94652b04-4b47-48c8-b2e4-c8029fe24d64',  'e2a33173-2c6f-4030-838d-c735efaafa6e',  'e642c8d3-ed74-4a42-b002-14a18993904f',  '9c820672-4873-477e-869c-a6f5d48cfbea',  'c675329a-a9b9-4279-ae09-699fc77fc337',  '1a9456e1-fa8b-45e6-b54a-d549af771559',  'c95c4b9b-8a3a-48ea-b407-65864e6d89c1',  '4bc7fa7e-345a-4230-9c95-b552de62a089',  '384421a1-2825-488f-bc82-c028cd1767c1',  'c1b2281c-d354-4541-ae04-79f818ce0a50',  '37a261ee-3e9c-4534-847a-0fedf47f1134',  '714de3cd-e8c1-43ca-8aa2-fe61002a1c6a',  '06caf869-91a7-4351-bfb4-c2b11ad16178',  '6555c2e3-f8cb-4ff0-a9d7-2481880fe6cf',  '0b04b7ae-35b2-461a-be2a-1d0bf6febc5b',  '0631e989-13c3-47ba-b46d-c00282b4e912',  '686b18d6-7b7a-4e04-b314-6a0845bf54fb',  '4b8d2ef5-55e1-44a9-93ed-5275f9578191',  'a9339c07-a1f6-4906-95f9-a3a20a6cf347',  '63e8c105-fe1d-4896-9248-d2d3bf704d1c',  'd187ec5e-0299-44db-b710-ad0e38d02d54',  '30d183ba-a8a0-4801-853d-c6c9653f3ad6',  '8d7f06e9-f0fd-4bde-99bb-62d4573954c6',  'e38d4907-53d2-49a0-a3c2-e386bdbddd31',  '09e2d149-d244-4868-bb04-3819c2af0f79',  'bf5f27c3-6c3a-4f2d-88ec-914458103d4e',  '69b713a3-3118-4145-90be-dd7683509d04',  'd6cbf68d-565b-4c4b-8216-7bf74ad13841',  '06c9ba86-ae41-42b0-9d19-201b1abb0e27',  '30aa8167-a158-4062-b889-847a9c7b47ae',  'd473d15e-6b78-4400-9a45-4ce41ef12bca',  '4dd31556-ef4b-43c7-bb32-fb23ab836d09',  '7af057a0-cae1-43f7-8f5b-16484393cc37',  '213f4225-2961-4d0e-bb05-3d0c4bf9948c',  '30d16af9-a3ee-45dd-bbef-18a9e3cbdb2a',  '56a803e1-ff46-4d18-bd21-d6aa90c4cc2d',  'e3a270c9-b461-4d67-89cb-558e0c95c0f6',  'b8a55500-af80-423c-96b5-fca6120b9a33',  'c385b56a-9947-416e-a93a-c1d3afed48c5',  '1db106f0-7a66-4be5-aac1-1b04c3960b61',  'f25ff490-fb6b-4b86-bf90-c02abadc250f',  'b0367c1c-b419-4c22-a5c0-f004dd1f9fd0',  '0bb1a81d-1c63-4d86-be5d-97e3f0182c3f')
  AND tce1."ExecutionStatusUpdateTimeStamp" = (
    SELECT MAX(tce2."ExecutionStatusUpdateTimeStamp")
    FROM "FenixExecution"."TestCasesUnderExecution" tce2
    WHERE tce2."TestCaseUuid" = tce1."TestCaseUuid"
      --AND tce2."TestCaseVersion" = tce1."TestCaseVersion"
  );

*/

// Load the latest Execution Status for TestSuites
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadLatestExecutionStatusForTestSuites(
	dbTransaction pgx.Tx,
	testSuiteUuidSlice []string,
	testSuiteExecutionUpdatedMinTimeStamp time.Time) (
	testSuitesLatestExecutionStatusMap map[string]*fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage,
	err error) {

	testSuitesLatestExecutionStatusMap = make(map[string]*fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage)

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT tsefl1.\"TestSuiteUuid\", tsefl1.\"TestSuiteVersion\", " +
		"tsefl1.\"TestSuiteExecutionStatus\", tsefl1.\"ExecutionStatusUpdateTimeStamp\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestSuitesExecutionsForListings\" tsefl1 "
	sqlToExecute = sqlToExecute + "WHERE tsefl1.\"TestSuiteUuid\" IN "
	sqlToExecute = sqlToExecute + fenixCloudDBObject.generateSQLINArray(testSuiteUuidSlice) + " "
	sqlToExecute = sqlToExecute + "AND tsefl1.\"ExecutionStatusUpdateTimeStamp\" = "
	sqlToExecute = sqlToExecute + "(SELECT MAX(tsefl2.\"ExecutionStatusUpdateTimeStamp\") "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestSuitesExecutionsForListings\" tsefl2 "
	sqlToExecute = sqlToExecute + "WHERE tsefl2.\"TestSuiteUuid\" = tsefl1.\"TestSuiteUuid\") AND "
	sqlToExecute = sqlToExecute + "tsefl1.\"ExecutionStatusUpdateTimeStamp\" > '" + common_config.GenerateDatetimeFromTimeInputForDB(testSuiteExecutionUpdatedMinTimeStamp) + "' "
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "5384a17f-955e-430e-9fa2-1dda17cee590",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadLatestExecutionStatusForTestSuites'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "2d3cdcf9-5b07-4ee1-99bb-baab7bf12c23",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	var (
		tempTestSuiteUuid                  string
		tempTestSuiteVersion               uint32
		tempLatestTestSuiteExecutionStatus int32
		tempInsertTimeStampAsTimeStamp     time.Time
	)

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&tempTestSuiteUuid,
			&tempTestSuiteVersion,
			&tempLatestTestSuiteExecutionStatus,
			&tempInsertTimeStampAsTimeStamp,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "10f01a37-466e-4840-ad74-879742f20211",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Create the 'tempTestSuiteThatCanBeEditedByUser'object
		var tempTestSuiteThatCanBeEditedByUser fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage
		tempTestSuiteThatCanBeEditedByUser = fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage{
			NonEditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage_NonEditableBasicInformationMessage{
				TestSuiteUuid:                 tempTestSuiteUuid,
				DomainUuid:                    "",
				DomainName:                    "",
				TestSuiteVersion:              tempTestSuiteVersion,
				TestSuiteExecutionEnvironment: "",
			},
			EditableInformation:                                     nil,
			LatestTestSuiteExecutionStatus:                          fenixTestCaseBuilderServerGrpcApi.TestSuiteExecutionStatusEnum(tempLatestTestSuiteExecutionStatus),
			LatestTestSuiteExecutionStatusInsertTimeStamp:           timestamppb.New(tempInsertTimeStampAsTimeStamp),
			LatestFinishedOkTestSuiteExecutionStatusInsertTimeStamp: nil,
			LastSavedTimeStamp:                                      nil,
			TestSuitePreview:                                        nil,
		}

		// Add to map of TestCases execution data
		testSuitesLatestExecutionStatusMap[tempTestSuiteThatCanBeEditedByUser.NonEditableInformation.TestSuiteUuid] = &tempTestSuiteThatCanBeEditedByUser

	}

	return testSuitesLatestExecutionStatusMap, err

}

/*
-- Latest OK Execution
SELECT tce1."TestCaseUuid", tce1."TestCaseVersion", tce1."TestCaseExecutionStatus", tce1."ExecutionStatusUpdateTimeStamp"
FROM "FenixExecution"."TestCasesUnderExecution" tce1
WHERE tce1."TestCaseUuid" IN ('4eebed04-39a9-4ad9-ae67-51c9de984486',  '653a43f7-2a7f-445b-aa3f-8cb4d1797369',  'df92530f-3a43-4e3d-9578-f300d996e2e5',  '2d4128e9-04f3-4615-aaf2-330193fcdbd7',  '94652b04-4b47-48c8-b2e4-c8029fe24d64',  'e2a33173-2c6f-4030-838d-c735efaafa6e',  'e642c8d3-ed74-4a42-b002-14a18993904f',  '9c820672-4873-477e-869c-a6f5d48cfbea',  'c675329a-a9b9-4279-ae09-699fc77fc337',  '1a9456e1-fa8b-45e6-b54a-d549af771559',  'c95c4b9b-8a3a-48ea-b407-65864e6d89c1',  '4bc7fa7e-345a-4230-9c95-b552de62a089',  '384421a1-2825-488f-bc82-c028cd1767c1',  'c1b2281c-d354-4541-ae04-79f818ce0a50',  '37a261ee-3e9c-4534-847a-0fedf47f1134',  '714de3cd-e8c1-43ca-8aa2-fe61002a1c6a',  '06caf869-91a7-4351-bfb4-c2b11ad16178',  '6555c2e3-f8cb-4ff0-a9d7-2481880fe6cf',  '0b04b7ae-35b2-461a-be2a-1d0bf6febc5b',  '0631e989-13c3-47ba-b46d-c00282b4e912',  '686b18d6-7b7a-4e04-b314-6a0845bf54fb',  '4b8d2ef5-55e1-44a9-93ed-5275f9578191',  'a9339c07-a1f6-4906-95f9-a3a20a6cf347',  '63e8c105-fe1d-4896-9248-d2d3bf704d1c',  'd187ec5e-0299-44db-b710-ad0e38d02d54',  '30d183ba-a8a0-4801-853d-c6c9653f3ad6',  '8d7f06e9-f0fd-4bde-99bb-62d4573954c6',  'e38d4907-53d2-49a0-a3c2-e386bdbddd31',  '09e2d149-d244-4868-bb04-3819c2af0f79',  'bf5f27c3-6c3a-4f2d-88ec-914458103d4e',  '69b713a3-3118-4145-90be-dd7683509d04',  'd6cbf68d-565b-4c4b-8216-7bf74ad13841',  '06c9ba86-ae41-42b0-9d19-201b1abb0e27',  '30aa8167-a158-4062-b889-847a9c7b47ae',  'd473d15e-6b78-4400-9a45-4ce41ef12bca',  '4dd31556-ef4b-43c7-bb32-fb23ab836d09',  '7af057a0-cae1-43f7-8f5b-16484393cc37',  '213f4225-2961-4d0e-bb05-3d0c4bf9948c',  '30d16af9-a3ee-45dd-bbef-18a9e3cbdb2a',  '56a803e1-ff46-4d18-bd21-d6aa90c4cc2d',  'e3a270c9-b461-4d67-89cb-558e0c95c0f6',  'b8a55500-af80-423c-96b5-fca6120b9a33',  'c385b56a-9947-416e-a93a-c1d3afed48c5',  '1db106f0-7a66-4be5-aac1-1b04c3960b61',  'f25ff490-fb6b-4b86-bf90-c02abadc250f',  'b0367c1c-b419-4c22-a5c0-f004dd1f9fd0',  '0bb1a81d-1c63-4d86-be5d-97e3f0182c3f')
  AND tce1."ExecutionStatusUpdateTimeStamp" = (
    SELECT MAX(tce2."ExecutionStatusUpdateTimeStamp")
    FROM "FenixExecution"."TestCasesUnderExecution" tce2
    WHERE tce2."TestCaseUuid" = tce1."TestCaseUuid"
      AND tce2."TestCaseExecutionStatus" IN (5, 6)
  );
*/

// Load the latest  Finished Execution Status for TestSuites
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadLatestFinishedOkExecutionStatusForTestSuites(
	dbTransaction pgx.Tx,
	testSuiteUuidSlice []string) (
	testSuitesLatestExecutionStatusMap map[string]*fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage,
	err error) {

	testSuitesLatestExecutionStatusMap = make(map[string]*fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage)

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT tsefl1.\"TestSuiteUuid\", tsefl1.\"TestSuiteVersion\", " +
		"tsefl1.\"TestSuiteExecutionStatus\", tsefl1.\"ExecutionStatusUpdateTimeStamp\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestSuitesExecutionsForListings\" tsefl1 "
	sqlToExecute = sqlToExecute + "WHERE tsefl1.\"TestSuiteUuid\" IN "
	sqlToExecute = sqlToExecute + fenixCloudDBObject.generateSQLINArray(testSuiteUuidSlice) + " "
	sqlToExecute = sqlToExecute + "AND tsefl1.\"ExecutionStatusUpdateTimeStamp\" = "
	sqlToExecute = sqlToExecute + "(SELECT MAX(tsefl2.\"ExecutionStatusUpdateTimeStamp\") "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestSuitesExecutionsForListings\" tsefl2 "
	sqlToExecute = sqlToExecute + "WHERE tsefl2.\"TestSuiteUuid\" = tsefl1.\"TestSuiteUuid\" AND "
	sqlToExecute = sqlToExecute + "tsefl2.\"TestSuiteExecutionStatus\" IN (5, 6))  "
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "431dfb42-6341-4a0a-8585-d3a83fd75724",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadLatestFinishedOkExecutionStatusForTestSuites'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "33a6a5c1-e405-4e4a-a7e8-9f9b1497e284",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	var (
		tempTestSuiteUuid                  string
		tempTestSuiteVersion               uint32
		tempLatestTestSuiteExecutionStatus int32
		tempInsertTimeStampAsTimeStamp     time.Time
	)

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&tempTestSuiteUuid,
			&tempTestSuiteVersion,
			&tempLatestTestSuiteExecutionStatus,
			&tempInsertTimeStampAsTimeStamp,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "725909f4-0f6b-4410-8e96-c986bbef2972",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Create the 'tempTestSuiteThatCanBeEditedByUser'-object
		var tempTestSuiteThatCanBeEditedByUser fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage
		tempTestSuiteThatCanBeEditedByUser = fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage{
			NonEditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestSuiteInformationMessage_NonEditableBasicInformationMessage{
				TestSuiteUuid:                 tempTestSuiteUuid,
				DomainUuid:                    "",
				DomainName:                    "",
				TestSuiteVersion:              tempTestSuiteVersion,
				TestSuiteExecutionEnvironment: "",
			},
			EditableInformation:                                     nil,
			LatestTestSuiteExecutionStatus:                          fenixTestCaseBuilderServerGrpcApi.TestSuiteExecutionStatusEnum(tempLatestTestSuiteExecutionStatus),
			LatestTestSuiteExecutionStatusInsertTimeStamp:           nil,
			LatestFinishedOkTestSuiteExecutionStatusInsertTimeStamp: timestamppb.New(tempInsertTimeStampAsTimeStamp),
			LastSavedTimeStamp:                                      nil,
			TestSuitePreview:                                        nil,
		}

		// Add to map of TestCases execution data
		testSuitesLatestExecutionStatusMap[tempTestSuiteThatCanBeEditedByUser.NonEditableInformation.TestSuiteUuid] = &tempTestSuiteThatCanBeEditedByUser

	}

	return testSuitesLatestExecutionStatusMap, err

}
