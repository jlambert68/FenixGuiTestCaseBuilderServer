package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
	"github.com/sirupsen/logrus"
)

// PrepareSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers
// Do initial preparations to be able to save all supported TestInstructions, TestInstructionContainers and Allowed Users
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "d0a3d9d6-6ee3-423a-aba5-81ad53be07d3",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'PrepareSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")

		return err

	}

	defer txn.Commit(context.Background())

	// Save  all supported TestInstructions, TestInstructionContainers and Allowed Users
	err = fenixCloudDBObject.saveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
		txn,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":    "657a20ff-d1b4-4568-a665-6a519547702a",
			"error": err,
		}).Error("Problem when saving supported TestInstructions, TestInstructionContainers and Allowed Users")

		return err
	}

	return err
}

// Save all supported TestInstructions, TestInstructionContainers and Allowed Users
func (fenixCloudDBObject *FenixCloudDBObjectStruct) saveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	dbTransaction pgx.Tx,
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	err error) {

	// Verify that Domain exists in database
	err = fenixCloudDBObject.verifyDomainExistsInDatabase(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "b248368c-efdf-475c-b2e3-8c4643a11c9d",
			"domainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"error":      err,
		}).Error("Domain does not exist in database")

		return err
	}

	// Get saved message hash for Domain
	var savedMessageHash string
	savedMessageHash, err = fenixCloudDBObject.prepareLoadSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessageHash(
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash)

	if err != nil {

		common_config.Logger.WithFields(logrus.Fields{
			"id":         "5986b01f-d584-470a-908e-6f8898fd71e1",
			"domainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"error":      err,
		}).Error("Couldn't get saved Message Hash from CloudDB")

		return err

	}

	// When the saved Message Hash is equal to the incoming Message Hash then nothing is change, which is the base case
	if savedMessageHash == testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.
		TestInstructionsAndTestInstructionsContainersAndUsersMessageHash {

		return nil
	}

	// When there is no Message Hash then just save the message in the database
	// or when this is a new 'baseline' for the domains supported TestInstructions, TestInstructionContainers and Allowed Users
	if savedMessageHash == "" ||
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.
			ForceNewBaseLineForTestInstructionsAndTestInstructionContainers == true {

		if savedMessageHash == "" {
			common_config.Logger.WithFields(logrus.Fields{
				"id":         "d8c8ef69-49f7-464e-b51f-23b5ca59bca9",
				"domainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
				"domainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			}).Info("No Message Hash found in database, so supported TestInstructions, TestInstructionContainers and Allowed Users will be saved")
		}

		if testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.
			ForceNewBaseLineForTestInstructionsAndTestInstructionContainers == true {
			common_config.Logger.WithFields(logrus.Fields{
				"id":         "ca054b92-f093-438a-bfb1-be5438ca3f33",
				"domainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
				"domainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			}).Info("New forced 'baseline' for the domains supported TestInstructions, TestInstructionContainers and Allowed Users")
		}

		// Save supported TestInstructions, TestInstructionContainers and Allowed Users in Database
		err = fenixCloudDBObject.performSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
			dbTransaction, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"id":         "056dfe6d-fa17-4e2b-a08b-453e597d033e",
				"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
				"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
				"error":      err,
			}).Error("Couldn't save supported TestInstructions, TestInstructionContainers and Allowed Users to CloudDB")

			return err
		}

		return err
	}

	// Verify changes to TestInstructions, TestInstructionContainers and Allowed Users separately
	err = fenixCloudDBObject.verifyChangesToTestInstructionsAndTestInstructionContainersAndAllowedUsersSeparately(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "9e17b75b-2461-49d1-ba6d-37d47db67e39",
			"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			"error":      err,
		}).Error("Got some problem when Verifying changes to TestInstructions, TestInstructionContainers and Allowed Users separately")

		return err
	}

	return err
}

// Verify that Domain exists in database
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifyDomainExistsInDatabase(
	dbTransaction pgx.Tx,
	domainUUID string) (
	err error) {

	fmt.Println(" **** Verify that Domain exists in database ****")

	return err
}

// Do the actual save for all supported TestInstructions, TestInstructionContainers and Allowed Users to database
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	dbTransaction pgx.Tx,
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	err error) {

	fmt.Println(" **** Do the actual save for all supported TestInstructions, TestInstructionContainers and Allowed Users to database ****")

	return err
}

// Verify changes to TestInstructions, TestInstructionContainers and Allowed Users separately
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifyChangesToTestInstructionsAndTestInstructionContainersAndAllowedUsersSeparately(
	dbTransaction pgx.Tx,
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	err error) {

	// Load saved message for supported TestInstructions, TestInstructionContainers and Allowed Users from database
	// To be used when comparing TestInstructions, TestInstructionContainers and Allowed Users
	var testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct
	testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB, err = fenixCloudDBObject.
		loadDomainSpecificPublishedSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage(
			dbTransaction,
			string(testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID))

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "dc62c173-d641-46c1-ba48-437885ac6b02",
			"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			"error":      err,
		}).Error("Got some problem when loading save message for supportedTestInstructions, TestInstructionContainers and Allowed Users from database")

		return err

	}

	// Verify changes in TestInstructions
	var correctNewChangesFoundInTestInstructions bool
	correctNewChangesFoundInTestInstructions, err = fenixCloudDBObject.verifyChangesToTestInstructions(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.TestInstructions,
		testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "1651733e-b147-4b9c-952a-1aba60618e88",
			"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			"error":      err,
		}).Error("Got some problem when verifying changes to TestInstructions")

		return err
	}

	// Verify changes in TestInstructionContainers
	var correctNewChangesFoundInTestInstructionContainers bool
	correctNewChangesFoundInTestInstructionContainers, err = fenixCloudDBObject.verifyChangesToTestInstructionContainers(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.TestInstructionContainers,
		testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "1651733e-b147-4b9c-952a-1aba60618e88",
			"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			"error":      err,
		}).Error("Got some problem when verifying changes to TestInstructionContainers")

		return err
	}

	// Verify changes in Allowed Users
	var correctNewChangesFoundInAllowedUsers bool
	correctNewChangesFoundInAllowedUsers, err = fenixCloudDBObject.verifyChangesToAllowedUsers(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.AllowedUsers,
		testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "1651733e-b147-4b9c-952a-1aba60618e88",
			"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			"error":      err,
		}).Error("Got some problem when verifying changes to Allowed Users")

		return err
	}

	// Found correct changes, so update supported TestInstructions, TestInstructionContainers and Allowed Users in database
	if correctNewChangesFoundInTestInstructions == true || correctNewChangesFoundInTestInstructionContainers ||
		correctNewChangesFoundInAllowedUsers {
		common_config.Logger.WithFields(logrus.Fields{
			"id": "83536474-a037-4b67-85fc-2818f9181e38",
			"correctNewChangesFoundInTestInstructions":          correctNewChangesFoundInTestInstructions,
			"correctNewChangesFoundInTestInstructionContainers": correctNewChangesFoundInTestInstructionContainers,
			"correctNewChangesFoundInAllowedUsers":              correctNewChangesFoundInAllowedUsers,
		}).Info("Found correct changes, so update supported TestInstructions, TestInstructionContainers and Allowed Users in database")
	}

	// Save new message with supported TestInstructions, TestInstructionContainers and Allowed Users in database
	err = fenixCloudDBObject.performSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
		dbTransaction,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":         "3df8e2a0-efb2-40a4-9fa2-58eb23ffa911",
			"DomainHash": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash,
			"DomainName": testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName,
			"error":      err,
		}).Error("Got some problem when saving supported TestInstructions, TestInstructionContainers and Allowed Users in database")

		return err
	}

	return err
}

// Verify changes in TestInstructions
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifyChangesToTestInstructions(
	dbTransaction pgx.Tx,
	testInstructionsMessage *TestInstructionAndTestInstuctionContainerTypes.TestInstructionsStruct,
	testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	correctNewChangesFoundInTestInstructions bool, err error) {

	// No TestInstructions
	if testInstructionsMessage != nil {

		return correctNewChangesFoundInTestInstructions, err
	}

	// There aren't any changes then just return
	if testInstructionsMessage.TestInstructionsHash !=
		testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB.TestInstructions.TestInstructionsHash {

		return correctNewChangesFoundInTestInstructions, err
	}

	// Now we know that there are changes, find which TestInstruction that are changed, added or removed
	// Loop existing TestInstruction that were stored in database
	for tempTestInstructionUUIDFromDB, tempTestInstructionFromDB := range testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB.TestInstructions.TestInstructionsMap {
		// Check if it exists in new message
		publishedTestInstructionFromConnector, existInMap := testInstructionsMessage.TestInstructionsMap[tempTestInstructionUUIDFromDB]

		// If it doesn't exit in new message then just continue
		if existInMap == false {
			continue
		}

		// If it has the same hash then just continue
		if publishedTestInstructionFromConnector.TestInstructionVersionsHash == tempTestInstructionFromDB.TestInstructionVersionsHash {
			continue
		}

		// If the hash is different then process the TestInstructionVersions
		if publishedTestInstructionFromConnector.TestInstructionVersionsHash != tempTestInstructionFromDB.TestInstructionVersionsHash {

			// Loop existing TestInstructionVersion that were stored in database
			for versionCounter, tempTestInstructionVersionFromDB := range tempTestInstructionFromDB.TestInstructionVersions {

				// If the version hash is different then something is wrong
				if tempTestInstructionVersionFromDB.TestInstructionInstanceVersionHash !=
					publishedTestInstructionFromConnector.TestInstructionVersions[versionCounter].
						TestInstructionInstanceVersionHash {

					var tempTestInstructionUUID string
					var tempMajorVersionNumber int
					var tempMinorVersionNumber int

					tempTestInstructionUUID = string(publishedTestInstructionFromConnector.
						TestInstructionVersions[versionCounter].TestInstructionInstance.TestInstruction.TestInstructionUUID)
					tempMajorVersionNumber = publishedTestInstructionFromConnector.
						TestInstructionVersions[versionCounter].TestInstructionInstance.TestInstruction.MajorVersionNumber
					tempMinorVersionNumber = publishedTestInstructionFromConnector.
						TestInstructionVersions[versionCounter].TestInstructionInstance.TestInstruction.MinorVersionNumber

					common_config.Logger.WithFields(logrus.Fields{
						"id":                  "8c91a76d-887a-477d-b5ea-168b0d0c6061",
						"TestInstructionUUID": tempTestInstructionUUID,
						"MajorVersionNumber":  tempMajorVersionNumber,
						"MinorVersionNumber":  tempMinorVersionNumber,
						"versionCounter":      versionCounter,
					}).Error("New TestInstructionVersion is not the same as existing, which was expected")

					err = errors.New(fmt.Sprintf("New TestInstructionVersion has not the same hash as existing, "+
						"which was expected. "+
						"TestInstructionUUID=%s, MajorVersionNumber=%d, MinorVersionNumber=%d, SlicePosition=%d ",
						tempTestInstructionUUID,
						tempMajorVersionNumber,
						tempMinorVersionNumber,
						versionCounter))

					return false, err

				}
			}
			// Verify if there are more TestInstructionVersions in new published TestInstructionVersions
			if len(publishedTestInstructionFromConnector.TestInstructionVersions) > len(tempTestInstructionFromDB.TestInstructionVersions) {

				// There are more TestInstructionVersions in new published TestInstructionVersions so accept that
				correctNewChangesFoundInTestInstructions = true

				return correctNewChangesFoundInTestInstructions, err

			} else {
				// This should happen because the versions-Hash can't be the same when there are more versions in new TestInstruction

				var tempTestInstructionUUID string
				tempTestInstructionUUID = string(publishedTestInstructionFromConnector.
					TestInstructionVersions[0].TestInstructionInstance.TestInstruction.TestInstructionUUID)

				common_config.Logger.WithFields(logrus.Fields{
					"id":                  "7dc54367-7290-4641-a43a-dfb4ca30a728",
					"TestInstructionUUID": tempTestInstructionUUID,
				}).Error("This should happen because the versions-Hash can't be the same when there are more versions in new TestInstruction")

				err = errors.New(fmt.Sprintf("This should happen because the versions-Hash can't be the same when "+
					"there are more versions in new TestInstruction. "+
					"TestInstructionUUID=%s",
					tempTestInstructionUUID))

				return false, err

			}
		}
	}

	return correctNewChangesFoundInTestInstructions, err
}

// Verify changes in TestInstructionContainers
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifyChangesToTestInstructionContainers(
	dbTransaction pgx.Tx,
	testInstructionContainersMessage *TestInstructionAndTestInstuctionContainerTypes.TestInstructionContainersStruct,
	testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	correctNewChangesFoundInTestInstructionContainers bool, err error) {

	// No TestInstructionContainers
	if testInstructionContainersMessage != nil {

		return correctNewChangesFoundInTestInstructionContainers, err
	}

	// There aren't any changes then just return
	if testInstructionContainersMessage.TestInstructionContainersHash !=
		testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB.TestInstructionContainers.TestInstructionContainersHash {

		return correctNewChangesFoundInTestInstructionContainers, err
	}

	// Now we know that there are changes, find which TestInstruction that are changed, added or removed
	// Loop existing TestInstruction that were stored in database
	for tempTestInstructionContainerUUIDFromDB, tempTestInstructionContainerFromDB := range testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB.TestInstructionContainers.TestInstructionContainersMap {

		// Check if it exists in new message
		publishedTestInstructionContainerFromConnector, existInMap := testInstructionContainersMessage.
			TestInstructionContainersMap[tempTestInstructionContainerUUIDFromDB]

		// If it doesn't exit in new message then just continue
		if existInMap == false {
			continue
		}

		// If it has the same hash then just continue
		if publishedTestInstructionContainerFromConnector.TestInstructionContainerVersionsHash ==
			tempTestInstructionContainerFromDB.TestInstructionContainerVersionsHash {
			continue
		}

		// If the hash is different then process the TestInstructionVersions
		if publishedTestInstructionContainerFromConnector.TestInstructionContainerVersionsHash !=
			tempTestInstructionContainerFromDB.TestInstructionContainerVersionsHash {

			// Loop existing TestInstructionContainerVersion that were stored in database
			for versionCounter, tempTestInstructionContainerVersionFromDB := range tempTestInstructionContainerFromDB.TestInstructionContainerVersions {

				// If the version hash is different then something is wrong
				if tempTestInstructionContainerVersionFromDB.TestInstructionContainerInstanceHash !=
					publishedTestInstructionContainerFromConnector.TestInstructionContainerVersions[versionCounter].
						TestInstructionContainerInstanceHash {

					var tempTestInstructionContainerUUID string
					var tempMajorVersionNumber int
					var tempMinorVersionNumber int

					tempTestInstructionContainerUUID = string(publishedTestInstructionContainerFromConnector.
						TestInstructionContainerVersions[versionCounter].TestInstructionContainerInstance.
						TestInstructionContainer.TestInstructionContainerUUID)
					tempMajorVersionNumber = publishedTestInstructionContainerFromConnector.
						TestInstructionContainerVersions[versionCounter].TestInstructionContainerInstance.
						TestInstructionContainer.MajorVersionNumber
					tempMinorVersionNumber = publishedTestInstructionContainerFromConnector.
						TestInstructionContainerVersions[versionCounter].TestInstructionContainerInstance.
						TestInstructionContainer.MinorVersionNumber

					common_config.Logger.WithFields(logrus.Fields{
						"id":                           "160988ea-7f7c-4138-8190-e0f6d87e4a75",
						"TestInstructionContainerUUID": tempTestInstructionContainerUUID,
						"MajorVersionNumber":           tempMajorVersionNumber,
						"MinorVersionNumber":           tempMinorVersionNumber,
						"versionCounter":               versionCounter,
					}).Error("New TestInstructionContainerVersion is not the same as existing, which was expected")

					err = errors.New(fmt.Sprintf("New TestInstructionContainerVersion has not the same hash as existing, "+
						"which was expected. "+
						"TestInstructionUUID=%s, MajorVersionNumber=%d, MinorVersionNumber=%d, SlicePosition=%d ",
						tempTestInstructionContainerUUID,
						tempMajorVersionNumber,
						tempMinorVersionNumber,
						versionCounter))

					return false, err

				}
			}
			// Verify if there are more TestInstructionContainerVersions in new published TestInstructionVersions
			if len(publishedTestInstructionContainerFromConnector.TestInstructionContainerVersions) >
				len(tempTestInstructionContainerFromDB.TestInstructionContainerVersions) {

				// There are more TestInstructionVersions in new published TestInstructionVersions so accept that
				correctNewChangesFoundInTestInstructionContainers = true

				return correctNewChangesFoundInTestInstructionContainers, err

			} else {
				// This should happen because the versions-Hash can't be the same when there are more versions in new TestInstructionContainer

				var tempTestInstructionContainerVersionUUID string
				tempTestInstructionContainerVersionUUID = string(publishedTestInstructionContainerFromConnector.
					TestInstructionContainerVersions[0].TestInstructionContainerInstance.TestInstructionContainer.TestInstructionContainerUUID)

				common_config.Logger.WithFields(logrus.Fields{
					"id":                  "ea5745aa-e943-44b3-b7d8-ea753535968c",
					"TestInstructionUUID": tempTestInstructionContainerVersionUUID,
				}).Error("This should happen because the versions-Hash can't be the same when there are more versions in new TestInstructionContainer")

				err = errors.New(fmt.Sprintf("This should happen because the versions-Hash can't be the same when "+
					"there are more versions in new TestInstructionContainer. "+
					"TestInstructionUUID=%s",
					tempTestInstructionContainerVersionUUID))

				return false, err

			}
		}
	}

	return correctNewChangesFoundInTestInstructionContainers, err
}

// Verify changes in Allowed Users
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifyChangesToAllowedUsers(
	dbTransaction pgx.Tx,
	allowedUsers *TestInstructionAndTestInstuctionContainerTypes.AllowedUsersStruct,
	testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	correctNewChangesFoundInAllowedUsers bool, err error) {

	// There aren't any changes then just return
	if allowedUsers.AllowedUsersHash !=
		testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB.AllowedUsers.AllowedUsersHash {

		return correctNewChangesFoundInAllowedUsers, err
	}

	// Hash is changed so accept the new User Setup
	correctNewChangesFoundInAllowedUsers = true

	return correctNewChangesFoundInAllowedUsers, err
}
