package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
	"github.com/sirupsen/logrus"
)

// Do initial preparations to be able to load all supported TestInstructions, TestInstructionContainers and Allowed Users for all domains
func (fenixCloudDBObject *FenixCloudDBObjectStruct) prepareLoadAllSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers() (
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages []*TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct,
	err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "d0a3d9d6-6ee3-423a-aba5-81ad53be07d3",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'prepareLoadAllSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")

		return nil, err

	}

	defer txn.Commit(context.Background())

	// Load  all supported TestInstructions, TestInstructionContainers and Allowed Users for all domains
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages, err = fenixCloudDBObject.
		loadAllSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(txn)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "7e3e96ae-2e16-407f-956a-f6c66d3ddb93",
			"error": err,
		}).Error("Couldn't load all supported TestInstructions, TestInstructionContainers and Allowed Users from CloudDB")

		return testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages, err
	}

	return testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages, err
}

// Load  all supported TestInstructions, TestInstructionContainers and Allowed Users for a specific domain
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadAllSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	dbTransaction pgx.Tx) (
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages []*TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct,
	err error) {

	fmt.Println(" **** Load  all supported TestInstructions, TestInstructionContainers and Allowed Users for a specific domain ****")

	return testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages, err
}
