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

// PrepareListTestCasesThatCanBeEdited
// List all TestCases from Database that the user can edit
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareListTestCasesThatCanBeEdited(
	gCPAuthenticatedUser string) (
	responseMessage *fenixTestCaseBuilderServerGrpcApi.ListTestCasesThatCanBeEditedResponseMessage) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "060eb7e1-a915-443a-9e9c-81d7ae13bd38",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'PrepareListTestCasesThatCanBeEdited'")

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

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestCasesThatCanBeEditedResponseMessage{
			AckNackResponse:                ackNackResponse,
			TestCasesThatCanBeEditedByUser: nil,
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
			"id":                   "e4c6807b-1b78-4d20-b219-7d47b030dea3",
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

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestCasesThatCanBeEditedResponseMessage{
			AckNackResponse:                ackNackResponse,
			TestCasesThatCanBeEditedByUser: nil,
		}

		return responseMessage

	}

	// Load the TestCase
	var testCasesThatCanBeEditedResponse []*fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage
	testCasesThatCanBeEditedResponse, err = fenixCloudDBObject.listTestCasesThatCanBeEdited(
		txn,
		domainAndAuthorizations)

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

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestCasesThatCanBeEditedResponseMessage{
			AckNackResponse:                ackNackResponse,
			TestCasesThatCanBeEditedByUser: nil,
		}

		return responseMessage
	}

	// TestCase
	if testCasesThatCanBeEditedResponse == nil {

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

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestCasesThatCanBeEditedResponseMessage{
			AckNackResponse:                ackNackResponse,
			TestCasesThatCanBeEditedByUser: nil,
		}

		return responseMessage
	}

	// Create a slice with all TestCaseUuid to be used for finding execution status
	var testCaseUuidSlice []string

	for _, tempTestCase := range testCasesThatCanBeEditedResponse {
		testCaseUuidSlice = append(testCaseUuidSlice, tempTestCase.TestCaseUuid)
	}

	// Load the latest Execution Status for TestCase
	var testCasesLatestExecutionStatusMap map[string]*fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage
	testCasesLatestExecutionStatusMap = make(map[string]*fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage)
	testCasesLatestExecutionStatusMap, err = fenixCloudDBObject.loadLatestExecutionStatusForTestCases(
		txn,
		testCaseUuidSlice)

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

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestCasesThatCanBeEditedResponseMessage{
			AckNackResponse:                ackNackResponse,
			TestCasesThatCanBeEditedByUser: nil,
		}

		return responseMessage
	}

	// Load the latest OK Execution Status for TestCase
	var testCasesLatestFinishedOkExecutionStatusMap map[string]*fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage
	testCasesLatestFinishedOkExecutionStatusMap = make(map[string]*fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage)
	testCasesLatestFinishedOkExecutionStatusMap, err = fenixCloudDBObject.loadLatestFinishedOkExecutionStatusForTestCases(
		txn,
		testCaseUuidSlice)

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

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestCasesThatCanBeEditedResponseMessage{
			AckNackResponse:                ackNackResponse,
			TestCasesThatCanBeEditedByUser: nil,
		}

		return responseMessage
	}

	// Merge Execution status into full TestCaseList
	var foundInMap bool
	var changesAreMade bool
	for testCaseIndex, temptestCase := range testCasesThatCanBeEditedResponse {

		// Reset 'changesAreMade'
		changesAreMade = false

		// Latest Execution Status Information
		var temptestCaseFromStatus *fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage
		temptestCaseFromStatus, foundInMap = testCasesLatestExecutionStatusMap[temptestCase.TestCaseUuid]

		// TestCaseExecution-status wasn't found in Map which indicates that there are no executions for the TestCase
		if foundInMap == false {

		} else {
			// Add Latest Status information
			temptestCase.LatestTestCaseExecutionStatus = temptestCaseFromStatus.LatestTestCaseExecutionStatus
			temptestCase.LatestTestCaseExecutionStatusInsertTimeStamp = temptestCaseFromStatus.
				LatestTestCaseExecutionStatusInsertTimeStamp

			// Indicate tha changes are done
			changesAreMade = true

		}

		// Latest Finished OK Execution Status Information
		var temptestCaseFromFinishedStatus *fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage
		temptestCaseFromFinishedStatus, foundInMap = testCasesLatestFinishedOkExecutionStatusMap[temptestCase.TestCaseUuid]

		// TestCaseExecution-status wasn't found in Map which indicates that there are no Finished OK executions for the TestCase
		if foundInMap == false {

		} else {
			// Add Latest Finished OK Status information
			temptestCase.LatestFinishedOkTestCaseExecutionStatusInsertTimeStamp = temptestCaseFromFinishedStatus.
				LatestFinishedOkTestCaseExecutionStatusInsertTimeStamp

			// Indicate tha changes are done
			changesAreMade = true

		}

		// Save back the TestCase into the Slice when changes are done
		if changesAreMade == true {
			testCasesThatCanBeEditedResponse[testCaseIndex] = temptestCase
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

	responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestCasesThatCanBeEditedResponseMessage{
		AckNackResponse:                ackNackResponse,
		TestCasesThatCanBeEditedByUser: testCasesThatCanBeEditedResponse,
	}

	return responseMessage
}

// Load all TestCases that the user can edit
func (fenixCloudDBObject *FenixCloudDBObjectStruct) listTestCasesThatCanBeEdited(
	dbTransaction pgx.Tx,
	domainAndAuthorizations []DomainAndAuthorizationsStruct) (
	testCasesThatCanBeEditedByUser []*fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage,
	err error) {

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
			"Id":    "9e68a788-dec3-4473-b9c6-9b752301da41",
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
			"Id":    "2b2e13cc-b485-4777-9cfa-c271ec05b65f",
			"Error": err,
			"tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain": tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain,
		}).Error("Couldn't convert into string representation")

		return nil, err
	}
	/*
		SELECT tc1."DomainUuid", tc1."DomainName", tc1."TestCaseUuid", tc1."TestCaseName", tc1."TestCaseVersion", tc1."InsertTimeStamp"
		FROM "FenixBuilder"."TestCases" TC1
		WHERE tc1."TestCaseIsDeleted"  = false AND tc1."InsertTimeStamp" IS NOT NULL  AND tc1."TestCaseVersion" = (SELECT MAX(tc2."TestCaseVersion") FROM "FenixBuilder"."TestCases" tc2 WHERE tc2."TestCaseUuid" = tc1."TestCaseUuid") ;
	*/
	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT tc1.\"DomainUuid\", tc1.\"DomainName\", tc1.\"TestCaseUuid\", " +
		"tc1.\"TestCaseName\", tc1.\"TestCaseVersion\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestCases\" tc1 "
	sqlToExecute = sqlToExecute + "WHERE (tc1.\"CanListAndViewTestCaseAuthorizationLevelOwnedByDomain\" & " + tempCanListAndViewTestCaseOwnedByThisDomainAsString + ")"
	sqlToExecute = sqlToExecute + "= tc1.\"CanListAndViewTestCaseAuthorizationLevelOwnedByDomain\" "
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "(tc1.\"CanListAndViewTestCaseAuthorizationLevelHavingTiAndTicWithDomai\" & " + tempCanListAndViewTestCaseHavingTIandTICfromThisDomainAsString + ")"
	sqlToExecute = sqlToExecute + "= tc1.\"CanListAndViewTestCaseAuthorizationLevelHavingTiAndTicWithDomai\" "
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "tc1.\"TestCaseIsDeleted\"  = false AND tc1.\"InsertTimeStamp\" IS NOT NULL " +
		"AND tc1.\"TestCaseVersion\" = (" +
		"SELECT MAX(tc2.\"TestCaseVersion\") " +
		"FROM \"FenixBuilder\".\"TestCases\" tc2 " +
		"WHERE tc2.\"TestCaseUuid\" = tc1.\"TestCaseUuid\") "
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
			"Id":           "490b7218-eecf-4cbe-90f7-eab91870f4bb",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Extract data from DB result set
	for rows.Next() {

		var tempTestCaseThatCanBeEditedByUser fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage

		err = rows.Scan(
			&tempTestCaseThatCanBeEditedByUser.DomainUuid,
			&tempTestCaseThatCanBeEditedByUser.DomainName,
			&tempTestCaseThatCanBeEditedByUser.TestCaseUuid,
			&tempTestCaseThatCanBeEditedByUser.TestCaseName,
			&tempTestCaseThatCanBeEditedByUser.TestCaseVersion,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "0402e085-9d27-455d-a3b4-52c90b3d43b8",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Add to slice of TestCases
		testCasesThatCanBeEditedByUser = append(testCasesThatCanBeEditedByUser, &tempTestCaseThatCanBeEditedByUser)

	}

	return testCasesThatCanBeEditedByUser, err

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

// Load the latest Execution Status for TestCases
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadLatestExecutionStatusForTestCases(
	dbTransaction pgx.Tx,
	testCaseUuidSlice []string) (
	testCasesLatestExecutionStatusMap map[string]*fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage,
	err error) {

	testCasesLatestExecutionStatusMap = make(map[string]*fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage)

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT tce1.\"TestCaseUuid\", tce1.\"TestCaseVersion\", " +
		"tce1.\"TestCaseExecutionStatus\", tce1.\"ExecutionStatusUpdateTimeStamp\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestCasesUnderExecution\" tce1 "
	sqlToExecute = sqlToExecute + "WHERE tce1.\"TestCaseUuid\" IN "
	sqlToExecute = sqlToExecute + fenixCloudDBObject.generateSQLINArray(testCaseUuidSlice) + " "
	sqlToExecute = sqlToExecute + "AND tce1.\"ExecutionStatusUpdateTimeStamp\" = "
	sqlToExecute = sqlToExecute + "(SELECT MAX(tce2.\"ExecutionStatusUpdateTimeStamp\") "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestCasesUnderExecution\" tce2 "
	sqlToExecute = sqlToExecute + "WHERE tce2.\"TestCaseUuid\" = tce1.\"TestCaseUuid\") "
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "ff2c0040-4560-4d45-83a6-bad1e18f0f0c",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadLatestExecutionStatusForTestCases'")
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
		tempInsertTimeStampAsTimeStamp time.Time
	)

	// Extract data from DB result set
	for rows.Next() {

		var tempTestCaseThatCanBeEditedByUser fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage

		err = rows.Scan(
			&tempTestCaseThatCanBeEditedByUser.TestCaseUuid,
			&tempTestCaseThatCanBeEditedByUser.TestCaseVersion,
			&tempTestCaseThatCanBeEditedByUser.LatestTestCaseExecutionStatus,
			&tempInsertTimeStampAsTimeStamp,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "cc07db8c-a30b-4059-a4c0-2ef72bcb71b0",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert DataTime into gRPC-version
		tempTestCaseThatCanBeEditedByUser.LatestTestCaseExecutionStatusInsertTimeStamp = timestamppb.New(tempInsertTimeStampAsTimeStamp)

		// Add to map of TestCases execution data
		testCasesLatestExecutionStatusMap[tempTestCaseThatCanBeEditedByUser.TestCaseUuid] = &tempTestCaseThatCanBeEditedByUser

	}

	return testCasesLatestExecutionStatusMap, err

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

// Load the latest  Finished Execution Status for TestCases
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadLatestFinishedOkExecutionStatusForTestCases(
	dbTransaction pgx.Tx,
	testCaseUuidSlice []string) (
	testCasesLatestExecutionStatusMap map[string]*fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage,
	err error) {

	testCasesLatestExecutionStatusMap = make(map[string]*fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage)

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT tce1.\"TestCaseUuid\", tce1.\"TestCaseVersion\", " +
		"tce1.\"TestCaseExecutionStatus\", tce1.\"ExecutionStatusUpdateTimeStamp\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestCasesUnderExecution\" tce1 "
	sqlToExecute = sqlToExecute + "WHERE tce1.\"TestCaseUuid\" IN "
	sqlToExecute = sqlToExecute + fenixCloudDBObject.generateSQLINArray(testCaseUuidSlice) + " "
	sqlToExecute = sqlToExecute + "AND tce1.\"ExecutionStatusUpdateTimeStamp\" = "
	sqlToExecute = sqlToExecute + "(SELECT MAX(tce2.\"ExecutionStatusUpdateTimeStamp\") "
	sqlToExecute = sqlToExecute + "FROM \"FenixExecution\".\"TestCasesUnderExecution\" tce2 "
	sqlToExecute = sqlToExecute + "WHERE tce2.\"TestCaseUuid\" = tce1.\"TestCaseUuid\" "
	sqlToExecute = sqlToExecute + "AND tce2.\"TestCaseExecutionStatus\" IN (5, 6))  "
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "4e53e9c5-7dad-431b-a01b-fe92c1c1b94b",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadLatestFinishedOkExecutionStatusForTestCases'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "14730d8d-1d10-45c0-81bf-088b01987d05",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	var (
		tempInsertTimeStampAsTimeStamp time.Time
	)

	// Extract data from DB result set
	for rows.Next() {

		var tempTestCaseThatCanBeEditedByUser fenixTestCaseBuilderServerGrpcApi.TestCaseThatCanBeEditedByUserMessage

		err = rows.Scan(
			&tempTestCaseThatCanBeEditedByUser.TestCaseUuid,
			&tempTestCaseThatCanBeEditedByUser.TestCaseVersion,
			&tempTestCaseThatCanBeEditedByUser.LatestTestCaseExecutionStatus,
			&tempInsertTimeStampAsTimeStamp,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "4c28b79c-cf64-4127-97c9-30f249526c60",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert DataTime into gRPC-version
		tempTestCaseThatCanBeEditedByUser.LatestTestCaseExecutionStatusInsertTimeStamp = timestamppb.New(tempInsertTimeStampAsTimeStamp)

		// Add to map of TestCases execution data
		testCasesLatestExecutionStatusMap[tempTestCaseThatCanBeEditedByUser.TestCaseUuid] = &tempTestCaseThatCanBeEditedByUser

	}

	return testCasesLatestExecutionStatusMap, err

}
