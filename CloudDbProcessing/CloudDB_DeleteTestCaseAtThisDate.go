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

func (fenixCloudDBObject *FenixCloudDBObjectStruct) saveDeleteTestCaseAtThisDateCommitOrRoleBack(
	dbTransactionReference *pgx.Tx,
	doCommitNotRoleBackReference *bool) {

	dbTransaction := *dbTransactionReference
	doCommitNotRoleBack := *doCommitNotRoleBackReference

	if doCommitNotRoleBack == true {
		dbTransaction.Commit(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "db07a734-57f4-45b7-bc91-7a36fae5b7c1",
		}).Debug("Doing Commit for SQL  in 'saveDeleteTestCaseAtThisDateCommitOrRoleBack'")

	} else {
		dbTransaction.Rollback(context.Background())

		common_config.Logger.WithFields(logrus.Fields{
			"id": "e73291cc-debb-40cb-afc0-beb524437e86",
		}).Info("Doing Rollback for SQL  in 'saveDeleteTestCaseAtThisDateCommitOrRoleBack'")

	}
}

// PrepareDeleteTestCaseAtThisDate
// Do initial preparations to be able to save the TestData from 'simple' TestDataArea-files
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareDeleteTestCaseAtThisDate(
	deleteTestCaseAtThisDateRequest *fenixTestCaseBuilderServerGrpcApi.DeleteTestCaseAtThisDateRequest) (
	err error) {

	// Begin SQL TransactionConnectorPublish
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "901f54ef-e547-44b1-93c1-74e6cfa86d20",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'PrepareDeleteTestCaseAtThisDate'")

		return err

	}

	// After all stuff is done, then Commit or Rollback depending on result
	var doCommitNotRoleBack bool

	// Standard is to do a Rollback
	doCommitNotRoleBack = false

	// When leaving then do the actual commit or rollback
	defer fenixCloudDBObject.saveDeleteTestCaseAtThisDateCommitOrRoleBack(
		&txn,
		&doCommitNotRoleBack)

	// Extract Domain that Owns the TestCase
	var ownerDomainForTestCase domainForTestCaseStruct
	ownerDomainForTestCase = domainForTestCaseStruct{
		domainUuid: deleteTestCaseAtThisDateRequest.GetDeleteThisTestCaseAtThisDate().GetDomainUuid(),
		domainName: deleteTestCaseAtThisDateRequest.GetDeleteThisTestCaseAtThisDate().GetDomainName(),
	}

	// Load Full TestCase from Database
	var tempDetailedTestCaseFromDatabase *fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse
	tempDetailedTestCaseFromDatabase = fenixCloudDBObject.PrepareLoadFullTestCase(
		deleteTestCaseAtThisDateRequest.GetDeleteThisTestCaseAtThisDate().GetTestCaseUuid(),
		deleteTestCaseAtThisDateRequest.GetUserIdentification().GetGCPAuthenticatedUser())

	// Extract all Domains that exist within all TestInstructions and TestInstructionContainers in the TestCase
	var allDomainsWithinTestCase []domainForTestCaseStruct
	allDomainsWithinTestCase = fenixCloudDBObject.extractAllDomainsWithinTestCase(tempDetailedTestCaseFromDatabase.DetailedTestCase)

	// Load Users all Domains
	var usersDomainsAndAuthorizations []DomainAndAuthorizationsStruct
	usersDomainsAndAuthorizations, err = fenixCloudDBObject.concatenateUsersDomainsAndDomainOpenToEveryOneToUse(
		txn, deleteTestCaseAtThisDateRequest.UserIdentification.GetGCPAuthenticatedUser())
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "d95403de-4470-4ebc-bcb2-c03844ba9ca6",
			"error": err,
		}).Error("Got some problem when loading Users Domains")

		return err
	}

	// Verify that User is allowed to Save TestCase
	var userIsAllowedToSaveTestCase bool
	userIsAllowedToSaveTestCase, _, _, err = fenixCloudDBObject.verifyThatUserIsAllowedToSaveTestCase(
		txn, ownerDomainForTestCase, allDomainsWithinTestCase, usersDomainsAndAuthorizations)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "1f2d4b36-6bb3-43e7-a894-6899c9a4371a",
			"error": err,
		}).Error("Some technical database problem when trying to verify if user is allowed to Delete TestCase")

		return err
	}

	// User is not allowed to save TestCase
	if userIsAllowedToSaveTestCase == false {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "58c31c4d-aa82-4c2f-8863-1701dae578b3",
			"error": err,
		}).Error("User is not allowed to delete TestCase in database")

		err = errors.New("user is not allowed to save TestCase in database")

		return err

	}

	// Verify that latest saved TestCase doesn't have a higher TestInstructionVersion
	if tempDetailedTestCaseFromDatabase.GetDetailedTestCase().GetTestCaseBasicInformation().GetBasicTestCaseInformation().
		GetNonEditableInformation().GetTestCaseVersion() !=
		deleteTestCaseAtThisDateRequest.GetDeleteThisTestCaseAtThisDate().TestCaseVersion {

		err = errors.New("there is a new TestCase saved in the database. Can't delete TestCase")

		return err
	}

	// Delete the TestCase
	err = fenixCloudDBObject.performDeleteTestCase(
		txn,
		deleteTestCaseAtThisDateRequest)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":                              "d1b07b1d-b308-436f-ae9c-fb042e2944bb",
			"error":                           err,
			"deleteTestCaseAtThisDateRequest": deleteTestCaseAtThisDateRequest,
		}).Error("Problem when updating Delete date for TestData")

		err = errors.New("Problem when updating Delete date for TestData. " + err.Error())

		return err
	}

	doCommitNotRoleBack = true

	return err
}

// Do the actual save for published TestData from oen TestDataArea
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performDeleteTestCase(
	dbTransaction pgx.Tx,
	deleteTestCaseAtThisDateRequest *fenixTestCaseBuilderServerGrpcApi.DeleteTestCaseAtThisDateRequest) (
	err error) {

	var tempTimestampToBeUsed string
	tempTimestampToBeUsed = common_config.GenerateDatetimeFromTimeInputForDB(time.Now())

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "UPDATE \"FenixBuilder\".\"TestCases\" "
	sqlToExecute = sqlToExecute + fmt.Sprintf("SET \"DeleteTimestamp\" = '%s', ", deleteTestCaseAtThisDateRequest.
		GetDeleteThisTestCaseAtThisDate().GetDeletedDate())
	sqlToExecute = sqlToExecute + fmt.Sprintf("\"DeletedInsertedTImeStamp\" = '%s' ", tempTimestampToBeUsed)
	sqlToExecute = sqlToExecute + fmt.Sprintf("\"DeletedByGCPAuthenticatedUser\" = '%s' ", deleteTestCaseAtThisDateRequest.
		GetUserIdentification().GetGCPAuthenticatedUser())
	sqlToExecute = sqlToExecute + fmt.Sprintf("WHERE \"TestCaseUuid\" = '%s' AND ", deleteTestCaseAtThisDateRequest.
		GetDeleteThisTestCaseAtThisDate().GetTestCaseUuid())
	sqlToExecute = sqlToExecute + fmt.Sprintf("\"TestCaseVersion\" = '%s' ", deleteTestCaseAtThisDateRequest.
		GetDeleteThisTestCaseAtThisDate().GetTestCaseVersion())
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "86593199-00af-43a8-a834-2c5bc67818c5",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'performDeleteTestCase'")
	}

	// Execute Query CloudDB
	comandTag, err := dbTransaction.Exec(context.Background(), sqlToExecute)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "81dbc7cf-c95d-46ae-8bca-23919827eb40",
			"sqlToExecute": sqlToExecute,
		}).Error("Got some problem when executing SQL within 'performDeleteTestCase'")

		return err
	}

	// Log response from CloudDB
	common_config.Logger.WithFields(logrus.Fields{
		"Id":                       "48ce081e-6a68-41d5-9486-cd9ed41bd9c9",
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
