package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v4"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
	"github.com/sirupsen/logrus"
	"time"
)

// Do initial preparations to be able to load all supported TestInstructions, TestInstructionContainers and Allowed Users for all domains
func (fenixCloudDBObject *FenixCloudDBObjectStruct) prepareLoadDomainSpecificPublishedSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage(
	domainUUID string) (
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct,
	err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "d0a3d9d6-6ee3-423a-aba5-81ad53be07d3",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'prepareLoadDomainSpecificPublishedSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage'")

		return nil, err

	}

	defer txn.Commit(context.Background())

	// Load all supported TestInstructions, TestInstructionContainers and Allowed Users for a specific domain
	var supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage *supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessageStruct
	supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage, err = fenixCloudDBObject.
		loadDomainSpecificPublishedSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage(
			txn, domainUUID)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "7e3e96ae-2e16-407f-956a-f6c66d3ddb93",
			"error":      err,
			"domainUUID": domainUUID,
		}).Error("Couldn't load all supported TestInstructions, TestInstructionContainers and Allowed Users from CloudDB")

		return testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage, err
	}

	// Convert into 'fenixCloudDBObject'
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage, err = fenixCloudDBObject.
		convertIntoTestInstructionsAndTestInstructionContainersFromGrpcBuilderMessage(
			supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "0d3d5912-4ed4-4dcd-a734-8008bd659862",
			"error":      err,
			"domainUUID": domainUUID,
		}).Error("Couldn't convert into a 'testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage'")

		return testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage, err
	}

	return testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage, err
}

type supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessageStruct struct {
	domainUUID                                     string
	domainName                                     string
	messageHash                                    string
	testInstructionsHash                           string
	testInstructionContainersHash                  string
	allowedUsersHash                               string
	supportedTIAndTICAndAllowedUsersMessageAsJsonb string
	updatedTimeStamp                               time.Time
	lastPublishedTimeStamp                         time.Time
}

// Load all supported TestInstructions, TestInstructionContainers and Allowed Users for a specific domain
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadDomainSpecificPublishedSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage(
	dbTransaction pgx.Tx,
	domainUUID string) (
	supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage *supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessageStruct,
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "6e5e9f62-b57a-40e2-8cd4-19b0c8ab1a7e",
	}).Debug("Entering: loadDomainSpecificPublishedSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "b76c3b72-9f1e-4512-8aea-ac1ec740e6e8",
		}).Debug("Exiting: loadDomainSpecificPublishedSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage()")
	}()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT * "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"SupportedTIAndTICAndAllowedUsers\" "
	sqlToExecute = sqlToExecute + "WHERE \"domainuuid\" = '" + domainUUID + "'"
	sqlToExecute = sqlToExecute + ";"

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "20f0c8a3-68ee-4a2d-8ef1-24951ef15495",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return nil, err
	}

	var rowsCounter int

	var (
		tempDomainUUID                                     string
		tempDomainName                                     string
		tempMessageHash                                    string
		tempTestInstructionsHash                           string
		tempTestInstructionContainersHash                  string
		tempAllowedUsersHash                               string
		tempSupportedTIAndTICAndAllowedUsersMessageAsJsonb string
		tempUpdatedTimeStamp                               time.Time
		tempLastPublishedTimeStamp                         time.Time
	)

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&tempDomainUUID,
			&tempDomainName,
			&tempMessageHash,
			&tempTestInstructionsHash,
			&tempTestInstructionContainersHash,
			&tempAllowedUsersHash,
			&tempSupportedTIAndTICAndAllowedUsersMessageAsJsonb,
			&tempUpdatedTimeStamp,
			&tempLastPublishedTimeStamp,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "26b85bcb-596a-415e-9195-b24d067c1da2",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert into message to be returned
		supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage = &supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessageStruct{
			domainUUID:                    tempDomainUUID,
			domainName:                    tempDomainName,
			messageHash:                   tempMessageHash,
			testInstructionsHash:          tempTestInstructionsHash,
			testInstructionContainersHash: tempTestInstructionContainersHash,
			allowedUsersHash:              tempAllowedUsersHash,
			supportedTIAndTICAndAllowedUsersMessageAsJsonb: tempSupportedTIAndTICAndAllowedUsersMessageAsJsonb,
			updatedTimeStamp:       tempUpdatedTimeStamp,
			lastPublishedTimeStamp: tempLastPublishedTimeStamp,
		}

		// Add to row counter; Max = 1
		rowsCounter = rowsCounter + 1

	}

	// Check how many rows that were found
	if rowsCounter > 1 {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "427c2341-4d82-45a7-aeaf-dfbc6352e307",
			"rowsCounter":  rowsCounter,
			"sqlToExecute": sqlToExecute,
		}).Error("More than 1 row was found in database")

		newErrorMessage := errors.New("More than 1 row was found in database for domain=" + domainUUID)

		return nil, newErrorMessage
	}

	return supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage, err
}

// Convert Database message into message to be used for returning to TesterGUI regarding supported TestInstructions,
// TestInstructionContainers and Allowed Users for a specific domain
func (fenixCloudDBObject *FenixCloudDBObjectStruct) convertIntoTestInstructionsAndTestInstructionContainersFromGrpcBuilderMessage(
	supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage *supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessageStruct) (
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct,
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "7be30c64-e243-4111-80b2-feaaf5db9a66",
	}).Debug("Entering: convertIntoTestInstructionsAndTestInstructionContainersFromGrpcBuilderMessage()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "cc1d80c2-3032-49f7-bf37-f4f9fb27c68d",
		}).Debug("Exiting: convertIntoTestInstructionsAndTestInstructionContainersFromGrpcBuilderMessage()")
	}()

	// Convert json-strings into byte-arrays
	tempTestCaseBasicInformationAsByteArray := []byte(supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage.supportedTIAndTICAndAllowedUsersMessageAsJsonb)

	// Convert json-byte-arrays into struct-messages
	err = json.Unmarshal(tempTestCaseBasicInformationAsByteArray, &testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":    "31db6991-fc80-4958-bfdc-c034a8008aca",
			"Error": err,
		}).Error("Something went wrong when converting 'tempTestCaseBasicInformationAsByteArray' into struct-message")

		return nil, err
	}

	return testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage, err
}
