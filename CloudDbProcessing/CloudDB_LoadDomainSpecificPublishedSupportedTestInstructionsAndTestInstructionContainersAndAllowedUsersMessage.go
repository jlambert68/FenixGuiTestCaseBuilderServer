package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"github.com/jackc/pgx/v4"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
	"github.com/sirupsen/logrus"
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
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage, err = fenixCloudDBObject.
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

	return testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage, err
}

// Load  all supported TestInstructions, TestInstructionContainers and Allowed Users for a specific domain
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadDomainSpecificPublishedSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage(
	dbTransaction pgx.Tx,
	domainUUID string) (
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct,
	err error) {

	return testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage, err
}
