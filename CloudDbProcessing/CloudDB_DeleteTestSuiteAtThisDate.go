package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"time"
)

func (fenixCloudDBObject *FenixCloudDBObjectStruct) saveDeleteTestSuiteAtThisDateCommitOrRoleBack(
	dbTransactionReference *pgx.Tx,
	doCommitNotRoleBackReference *bool) {

	dbTransaction := *dbTransactionReference
	doCommitNotRoleBack := *doCommitNotRoleBackReference

	if doCommitNotRoleBack == true {
		dbTransaction.Commit(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "1af9d744-a4d9-4c60-a2b9-a68ca70ce679",
		}).Debug("Doing Commit for SQL  in 'saveDeleteTestSuiteAtThisDateCommitOrRoleBack'")

	} else {
		dbTransaction.Rollback(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "de4b06ab-6a6c-4499-9058-eb3a2a3d9795",
		}).Info("Doing Rollback for SQL  in 'saveDeleteTestSuiteAtThisDateCommitOrRoleBack'")

	}
}

// PrepareDeleteTestSuiteAtThisDate
// Do initial preparations to be able to update delete of TestSuite
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareDeleteTestSuiteAtThisDate(
	deleteTestSuiteAtThisDateRequest *fenixTestCaseBuilderServerGrpcApi.DeleteTestSuiteAtThisDateRequest) (
	err error) {

	// Begin SQL TransactionConnectorPublish
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "1775c6d4-6aef-4fb5-8799-9010f591bfe7",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'PrepareDeleteTestSuiteAtThisDate'")

		return err

	}

	// After all stuff is done, then Commit or Rollback depending on result
	var doCommitNotRoleBack bool

	// Standard is to do a Rollback
	doCommitNotRoleBack = false

	// When leaving then do the actual commit or rollback
	defer fenixCloudDBObject.saveDeleteTestSuiteAtThisDateCommitOrRoleBack(
		&txn,
		&doCommitNotRoleBack)

	// Extract Domain that Owns the TestSuite
	var ownerDomainForTestSuite domainForTestCaseOrTestSuiteStruct
	ownerDomainForTestSuite = domainForTestCaseOrTestSuiteStruct{
		domainUuid: deleteTestSuiteAtThisDateRequest.GetDeleteThisTestSuiteAtThisDate().GetDomainUuid(),
		domainName: deleteTestSuiteAtThisDateRequest.GetDeleteThisTestSuiteAtThisDate().GetDomainName(),
	}

	// Load Full TestSuite from Database
	var tempDetailedTestSuiteFromDatabase *fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse
	tempDetailedTestSuiteFromDatabase = fenixCloudDBObject.PrepareLoadFullTestSuite(
		deleteTestSuiteAtThisDateRequest.GetDeleteThisTestSuiteAtThisDate().GetTestSuiteUuid(),
		deleteTestSuiteAtThisDateRequest.GetUserIdentification().GetGCPAuthenticatedUser())

	if tempDetailedTestSuiteFromDatabase == nil {

		err = errors.New(fmt.Sprintf("didn't find any existing TestSuite in Database with TestSuiteUuid = '%s' and TestSuiteVersion = '%d'",
			deleteTestSuiteAtThisDateRequest.GetDeleteThisTestSuiteAtThisDate().GetTestSuiteUuid(),
			deleteTestSuiteAtThisDateRequest.GetDeleteThisTestSuiteAtThisDate().GetTestSuiteVersion()))

		return err
	}

	// Extract all TestCaseUuid from TestSuite
	var testCaseUuidsInTestSuite []string
	if tempDetailedTestSuiteFromDatabase.GetDetailedTestSuite().TestCasesInTestSuite.TestCasesInTestSuite != nil {
		for _, tempTestCasesInTestSuite := range tempDetailedTestSuiteFromDatabase.GetDetailedTestSuite().
			TestCasesInTestSuite.TestCasesInTestSuite {
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

			errorMsg := errors.New("Got some problem when loading TestInstructions and TestInstructionContainers from TestSuite. " + err.Error())

			return errorMsg
		}

	}

	// Extract all Domains that exist within all TestInstructions and TestInstructionContainers in the TestSuite
	var allDomainsWithinTestSuite []domainForTestCaseOrTestSuiteStruct
	allDomainsWithinTestSuite = fenixCloudDBObject.extractAllDomainsWithinTestSuite(
		testInstructionsInTestSuite,
		testInstructionContainersInTestSuite)

	// Load Users all Domains
	var usersDomainsAndAuthorizations []DomainAndAuthorizationsStruct
	usersDomainsAndAuthorizations, err = fenixCloudDBObject.concatenateUsersDomainsAndDomainOpenToEveryOneToUse(
		txn, deleteTestSuiteAtThisDateRequest.GetUserIdentification().GetGCPAuthenticatedUser())
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "f61737d7-2327-445c-be12-ea2b7b4aefc9",
			"error": err,
		}).Error("Got some problem when loading Users Domains")

		errorMsg := errors.New("Got some problem when loading Users Domains. " + err.Error())

		return errorMsg

	}

	// Verify that User is allowed to Save TestSuite
	var userIsAllowedToSaveTestSuite bool
	//var authorizationValueForOwnerDomain int64
	//var authorizationValueForAllDomainsInTestSuite int64
	userIsAllowedToSaveTestSuite, _, _,
		err = fenixCloudDBObject.verifyThatUserIsAllowedToSaveTestSuite(
		txn, ownerDomainForTestSuite, allDomainsWithinTestSuite, usersDomainsAndAuthorizations)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "c307dbf4-c2c8-4012-966f-7e70f077810d",
			"error": err,
		}).Error("Some technical database problem when trying to verify if user is allowed to save TestSuite")

		errorMsg := errors.New("Some technical database problem when trying to verify if user is allowed to save TestSuite. " + err.Error())

		return errorMsg
	}

	// User is not allowed to save TestSuite
	if userIsAllowedToSaveTestSuite == false {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "9347b3bd-752c-4b96-87ee-d75152d5ca8f",
			"error": err,
		}).Error("User is not allowed to save TestSuite in database")

		errorMsg := errors.New("User is not allowed to save TestSuite in database. " + err.Error())

		return errorMsg

	}

	// Delete the TestSuite
	err = fenixCloudDBObject.performDeleteTestSuite(
		txn,
		deleteTestSuiteAtThisDateRequest)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":                               "2cafb253-6988-4d8b-9e13-3c1a61f5e3f6",
			"error":                            err,
			"deleteTestSuiteAtThisDateRequest": deleteTestSuiteAtThisDateRequest,
		}).Error("Problem when updating Delete date for TestSuite")

		err = errors.New("Problem when updating Delete date for TestSuite. " + err.Error())

		return err
	}

	doCommitNotRoleBack = true

	return err
}

// Do the actual update of 'delete' of the TestSuite
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performDeleteTestSuite(
	dbTransaction pgx.Tx,
	deleteTestSuiteAtThisDateRequest *fenixTestCaseBuilderServerGrpcApi.DeleteTestSuiteAtThisDateRequest) (
	err error) {

	var tempTimestampToBeUsed string
	tempTimestampToBeUsed = common_config.GenerateDatetimeFromTimeInputForDB(time.Now())

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "UPDATE \"FenixBuilder\".\"TestSuites\" "
	sqlToExecute = sqlToExecute + fmt.Sprintf("SET \"DeleteTimestamp\" = '%s', ", deleteTestSuiteAtThisDateRequest.
		GetDeleteThisTestSuiteAtThisDate().GetDeletedDate())
	sqlToExecute = sqlToExecute + fmt.Sprintf("\"DeleteTimeStamp\" = '%s', ", tempTimestampToBeUsed)
	sqlToExecute = sqlToExecute + fmt.Sprintf("\"DeletedByGCPAuthenticatedUser\" = '%s', ", deleteTestSuiteAtThisDateRequest.
		GetUserIdentification().GetGCPAuthenticatedUser())
	sqlToExecute = sqlToExecute + fmt.Sprintf("\"TestSuiteIsDeleted\" = %s ", "true ")
	sqlToExecute = sqlToExecute + fmt.Sprintf("WHERE \"TestSuiteUuid\" = '%s' AND ", deleteTestSuiteAtThisDateRequest.
		GetDeleteThisTestSuiteAtThisDate().GetTestSuiteUuid())
	sqlToExecute = sqlToExecute + fmt.Sprintf("\"TestSuiteVersion\" = '%d' ", deleteTestSuiteAtThisDateRequest.
		GetDeleteThisTestSuiteAtThisDate().GetTestSuiteVersion())
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "e0a42dc9-9b57-4b32-8542-682078c9cec4",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'performDeleteTestSuite'")
	}

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "2f1a941b-bc78-494b-bf19-1630a88f1e5e",
			"sqlToExecute": sqlToExecute,
		}).Error("Got some problem when executing SQL within 'performDeleteTestSuite'")

		return err
	}

	// Log response from CloudDB
	common_config.Logger.WithFields(logrus.Fields{
		"Id":                       "e3a169fd-c2fb-478d-8756-5aa9d5e5643d",
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
