package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
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

	// Load  all supported TestInstructions, TestInstructionContainers and Allowed Users for a specific domain
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
	sqlToExecute = sqlToExecute + "SELECT STITICAU.* "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"SupportedTIAndTICAndAllowedUsers\" STITICAU "
	sqlToExecute = sqlToExecute + "WHERE STITICAU.\"DomainUUID\" = '" + domainUUID + "'"
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

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage.domainUUID,
			&supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage.domainName,
			&supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage.messageHash,
			&supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage.testInstructionsHash,
			&supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage.testInstructionContainersHash,
			&supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage.allowedUsersHash,
			&supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage.supportedTIAndTICAndAllowedUsersMessageAsJsonb,
			&supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage.updatedTimeStamp,
			&supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersDbMessage.lastPublishedTimeStamp,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "26b85bcb-596a-415e-9195-b24d067c1da2",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
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
		"Id": "6e5e9f62-b57a-40e2-8cd4-19b0c8ab1a7e",
	}).Debug("Entering: loadDomainSpecificPublishedSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "b76c3b72-9f1e-4512-8aea-ac1ec740e6e8",
		}).Debug("Exiting: loadDomainSpecificPublishedSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage()")
	}()

	// Convert json-strings into byte-arrays
	tempTestCaseBasicInformationAsByteArray = []byte(tempTestCaseBasicInformationAsString)
	tempTestInstructionsAsByteArray = []byte(tempTestInstructionsAsString)
	tempTestInstructionContainersAsByteArray = []byte(tempTestInstructionContainersAsString)
	tempTestCaseExtraInformationAsByteArray = []byte(tempTestCaseExtraInformationAsString)

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

	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage = &TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct{
		TestInstructions:          nil,
		TestInstructionContainers: nil,
		AllowedUsers:              nil,
		MessageCreationTimeStamp:  time.Time{},
		TestInstructionsAndTestInstructionsContainersAndUsersMessageHash: "",
		ForceNewBaseLineForTestInstructionsAndTestInstructionContainers:  false,
		ConnectorsDomain: TestInstructionAndTestInstuctionContainerTypes.ConnectorsDomainStruct{},
	}

	return testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage, err
}
