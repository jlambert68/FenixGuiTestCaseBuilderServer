package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
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
	err = fenixCloudDBObject.verifyDomainExistsInDatabase(dbTransaction, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainHash)
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
			}).Info("No Message Hash found in database, so supported TestInstructions, TestInstructionContainers and Allowed Users will be saved")
		}

		if testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.
			ForceNewBaseLineForTestInstructionsAndTestInstructionContainers == true {
			common_config.Logger.WithFields(logrus.Fields{
				"id":         "d8c8ef69-49f7-464e-b51f-23b5ca59bca9",
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

	return err
}

// Do the actual save for all supported TestInstructions, TestInstructionContainers and Allowed Users to database
func (fenixCloudDBObject *FenixCloudDBObjectStruct) performSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	dbTransaction pgx.Tx,
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	err error) {

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
			testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID)

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
		}).Error("Got some problem when verifying changes to TestInstructions")

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
		}).Error("Got some problem when verifying changes to TestInstructions")

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

	return correctNewChangesFoundInTestInstructionContainers, err
}

// Verify changes in Allowed Users
func (fenixCloudDBObject *FenixCloudDBObjectStruct) verifyChangesToAllowedUsers(
	dbTransaction pgx.Tx,
	allowedUsers *TestInstructionAndTestInstuctionContainerTypes.AllowedUsersStruct,
	testInstructionsAndTestInstructionContainersAndAllowedUsersMessageSavedInDB *TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	correctNewChangesFoundInAllowedUsers bool, err error) {

	// No Allowed Users
	if allowedUsers != nil {

		return correctNewChangesFoundInAllowedUsers, err
	}

	return correctNewChangesFoundInAllowedUsers, err
}
