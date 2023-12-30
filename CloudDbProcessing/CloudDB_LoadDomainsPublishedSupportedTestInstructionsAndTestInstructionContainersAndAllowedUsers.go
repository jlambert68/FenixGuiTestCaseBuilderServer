package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"github.com/jackc/pgx/v4"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
	"github.com/sirupsen/logrus"
	"time"
)

// Do initial preparations to be able to load all supported TestInstructions, TestInstructionContainers and Allowed Users for a list of domains
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareLoadDomainsSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	domainList []string) (
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages []*TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct,
	err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "1c67ca25-dad3-490c-b666-bff1c423d33f",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'prepareLoadDomainsSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")

		return nil, err

	}

	defer txn.Commit(context.Background())

	// Load  all supported TestInstructions, TestInstructionContainers and Allowed Users for a list ofdomains
	var supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessages []*supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessageStruct
	supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessages, err = fenixCloudDBObject.
		loadDomainsSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
			txn,
			domainList)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "7e3e96ae-2e16-407f-956a-f6c66d3ddb93",
			"error":      err,
			"domainList": domainList,
		}).Error("Couldn't load all supported TestInstructions, TestInstructionContainers and Allowed Users from CloudDB")

		return testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages, err
	}

	// Loop all messages received from database and convert into 'fenixCloudDBObjects'
	for _, supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage := range supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessages {

		// Convert into 'fenixCloudDBObject'
		var testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
			TestInstructionsAndTestInstructionsContainersStruct

		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage, err = fenixCloudDBObject.
			convertIntoTestInstructionsAndTestInstructionContainersFromGrpcBuilderMessage(
				supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"id":         "0fe873a5-02b0-40e6-9689-1b2ae2576169",
				"error":      err,
				"domainList": domainList,
			}).Error("Couldn't convert into a 'testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage'")

			return nil, err
		}

		// Append converted 'fenixCloudDBObject' into list of objects
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages = append(
			testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages,
			testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)
	}

	return testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages, err
}

// Load  all supported TestInstructions, TestInstructionContainers and Allowed Users for a list of domains
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadDomainsSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	dbTransaction pgx.Tx,
	domainList []string) (
	supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessages []*supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessageStruct,
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "550b6189-e1a4-48ae-b3ed-978e7b9232fe",
	}).Debug("Entering: loadDomainsSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "8e3a4ed6-c618-453f-8ca3-ddbf714c4fca",
		}).Debug("Exiting: loadDomainsSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers()")
	}()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT * "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"SupportedTIAndTICAndAllowedUsers\" "
	sqlToExecute = sqlToExecute + "WHERE \"domainuuid\" IN " + common_config.GenerateSQLINArray(domainList) + " "
	sqlToExecute = sqlToExecute + ";"

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "e69ad9a8-ba39-404b-99a3-27b781e145a2",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return nil, err
	}

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
				"Id":           "21b3b638-937c-4e3c-b437-0cfd450de40e",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert into message to be returned
		var supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage *supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessageStruct
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

		// Append message to list of messages
		supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessages = append(
			supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessages,
			supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage)

	}

	return supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessages, err
}
