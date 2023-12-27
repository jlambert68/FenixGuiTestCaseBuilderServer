package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"time"
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

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "368aba80-5eaa-4847-84e1-96150cb22d68",
	}).Debug("Entering: loadSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessageHash()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "0837db94-32b1-417e-9ad0-823f06536caa",
		}).Debug("Exiting: loadSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessageHash()")
	}()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT STITICAU.\"messagehash\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"SupportedTIAndTICAndAllowedUsers\" STITICAU "
	sqlToExecute = sqlToExecute + "WHERE STITICAU.\"domainuuid\" = '" + domainUUID + "' "
	sqlToExecute = sqlToExecute + ";"

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "94d20f62-2963-4893-8dbd-8dc79ba4d1ea",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return "", err
	}

	var rowsCounter int

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&messageHash,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "d15964db-198e-4f7d-9e99-0eb47871198c",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return "", err
		}

		// Add to row counter; Max = 1
		rowsCounter = rowsCounter + 1

	}

	// Check how many rows that were found
	if rowsCounter > 1 {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "de645faa-8f61-4234-9493-10761ce47452",
			"rowsCounter":  rowsCounter,
			"sqlToExecute": sqlToExecute,
		}).Error("More than 1 row was found in database")

		newErrorMessage := errors.New("More than 1 row was found in database for domain=" + domainUUID)

		return "", newErrorMessage
	}

	return messageHash, err
}
