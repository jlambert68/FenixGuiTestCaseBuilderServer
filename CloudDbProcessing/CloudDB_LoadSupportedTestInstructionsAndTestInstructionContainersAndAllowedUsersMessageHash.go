package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

// Do initial preparations to be able to load Message Hash for supported TestInstructions, TestInstructionContainers and Allowed Users
func (fenixCloudDBObject *FenixCloudDBObjectStruct) prepareLoadSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessageHash(
	domainUUID string) (
	messageHash string,
	err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "af0498aa-dec6-4649-a58d-aeaeaff3bc61",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'prepareLoadSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessageHash'")

		return "", err

	}

	defer txn.Commit(context.Background())

	// Load  the Message Hash for a specific Domain for supported TestInstructions, TestInstructionContainers and Allowed Users for all domains
	messageHash, err = fenixCloudDBObject.
		loadSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessageHash(txn, domainUUID)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "f4df9aad-dab6-4678-978d-be375bde28d9",
			"error": err,
		}).Error("Couldn't load Message Hash for Domain regarding supported TestInstructions, TestInstructionContainers and Allowed Users from CloudDB")

		return messageHash, err
	}

	return messageHash, err
}

// Load Message Hash for supported TestInstructions, TestInstructionContainers and Allowed Users for specific domain
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessageHash(
	dbTransaction pgx.Tx,
	domainUUID string) (
	messageHash string,
	err error) {

	fmt.Println(" **** Load Message Hash for supported TestInstructions, TestInstructionContainers and Allowed Users for specific domain ****")

	return messageHash, err
}
