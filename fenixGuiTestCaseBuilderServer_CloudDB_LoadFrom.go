package main

import (
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

// ****************************************************************************************************************
// Load data from CloudDB into memory structures
//
// Load TestInstructions and pre-created TestInstructionContainers for Client
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObject_struct) loadClientsTestInstructionsFromCloudDB(userID string, cloudDBTestInstructionItems *[]*fenixTestCaseBuilderServerGrpcApi.TestInstructionMessage) (err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "61b8b021-9568-463e-b867-ac1ddb10584d",
	}).Debug("Entering: loadClientsTestInstructionsFromCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "78a97c41-a098-4122-88d2-01ed4b6c4844",
		}).Debug("Exiting: loadClientsTestInstructionsFromCloudDB()")
	}()

	/* Example
	SELECT FGB_TI."DomainUuid", FGB_TI."DomainName", FGB_TI."SystemUuid", FGB_TI."SystemName", FGB_TI."TestInstructionUuid", FGB_TI."TestInstructionName", FGB_TI."TestInstructionTypeUuid", FGB_TI."TestInstructionTypeName", FGB_TI."TestInstructionDescription", FGB_TI."TestInstructionMouseOverText", FGB_TI."Deprecated", FGB_TI."Enabled", FGB_TI."MajorVersionNumber", FGB_TI."MinorVersionNumber", FGB_TI."UpdatedTimeStamp"
	FROM "FenixGuiBuilder"."TestInstructions" FGB_TI;

	    "DomainUuid"                   uuid      not null,
	    "DomainName"                   varchar   not null,
	    "TestInstructionUuid"          uuid      not null (Key)
	    "TestInstructionName"          varchar   not null,
	    "TestInstructionTypeUuid"      uuid      not null,
	    "TestInstructionTypeName"      varchar   not null,
	    "TestInstructionDescription"   varchar   not null,
	    "TestInstructionMouseOverText" varchar   not null,
	    "Deprecated"                   boolean   not null,
	    "Enabled"                      boolean   not null,
	    "MajorVersionNumber"           integer   not null,
	    "MinorVersionNumber"           integer   not null,
	    "UpdatedTimeStamp"             timestamp not null

	*/

	usedDBSchema := "TestInstructions" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT * "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestInstructions\" FGB_TI;"

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "2f130d7e-f8aa-466f-b29d-0fb63608c1a6",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Variables to used when extract data from result set
	var cloudDBTestInstructionItem fenixTestCaseBuilderServerGrpcApi.TestInstructionMessage

	// Extract data from DB result set
	for rows.Next() {
		err := rows.Scan(
			&cloudDBTestInstructionItem.DomainUuid,
			&cloudDBTestInstructionItem.DomainName,
			&cloudDBTestInstructionItem.TestInstructionUuid,
			&cloudDBTestInstructionItem.TestInstructionName,
			&cloudDBTestInstructionItem.TestInstructionTypeUuid,
			&cloudDBTestInstructionItem.TestInstructionTypeName,
			&cloudDBTestInstructionItem.TestInstructionDescription,
			&cloudDBTestInstructionItem.TestInstructionMouseOverText,
			&cloudDBTestInstructionItem.Deprecated,
			&cloudDBTestInstructionItem.Enabled,
			&cloudDBTestInstructionItem.MajorVersionNumber,
			&cloudDBTestInstructionItem.MinorVersionNumber,
			&cloudDBTestInstructionItem.UpdatedTimeStamp,
		)

		if err != nil {
			return err
		}

		// Add values to the object that is pointed to by variable in function
		*cloudDBTestInstructionItems = append(*cloudDBTestInstructionItems, &cloudDBTestInstructionItem)

	}

	// No errors occurred
	return nil

}

// ****************************************************************************************************************
// Load data from CloudDB into memory structures
//
// Load pre-created TestInstructionContainerContainers for Client
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObject_struct) loadClientsTestInstructionContainersFromCloudDB(userID string, cloudDBTestInstructionContainerItems *[]*fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerMessage) (err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "216ddb15-6413-4051-a6a9-f479e6dd429b",
	}).Debug("Entering: loadClientsTestInstructionContainersFromCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "40ccee29-d32e-4674-9a3a-fd4403b55d32",
		}).Debug("Exiting: loadClientsTestInstructionContainersAndTestInstructionContainersContainersFromCloudDB()")
	}()

	/* Example
	SELECT FGB_TI."DomainUuid", FGB_TI."DomainName", FGB_TI."SystemUuid", FGB_TI."SystemName", FGB_TI."TestInstructionContainerUuid", FGB_TI."TestInstructionContainerName", FGB_TI."TestInstructionContainerTypeUuid", FGB_TI."TestInstructionContainerTypeName", FGB_TI."TestInstructionContainerDescription", FGB_TI."TestInstructionContainerMouseOverText", FGB_TI."Deprecated", FGB_TI."Enabled", FGB_TI."MajorVersionNumber", FGB_TI."MinorVersionNumber", FGB_TI."UpdatedTimeStamp"
	FROM "FenixGuiBuilder"."TestInstructionContainers" FGB_TI;

	    "DomainUuid"                   uuid      not null,
	    "DomainName"                   varchar   not null,
	    "TestInstructionContainerUuid"          uuid      not null (Key)
	    "TestInstructionContainerName"          varchar   not null,
	    "TestInstructionContainerTypeUuid"      uuid      not null,
	    "TestInstructionContainerTypeName"      varchar   not null,
	    "TestInstructionContainerDescription"   varchar   not null,
	    "TestInstructionContainerMouseOverText" varchar   not null,
	    "Deprecated"                   boolean   not null,
	    "Enabled"                      boolean   not null,
	    "MajorVersionNumber"           integer   not null,
	    "MinorVersionNumber"           integer   not null,
	    "UpdatedTimeStamp"             timestamp not null

	*/

	usedDBSchema := "TestInstructionContainers" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT * "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestInstructionContainers\" FGB_TI;"

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "2f130d7e-f8aa-466f-b29d-0fb63608c1a6",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Variables to used when extract data from result set
	var cloudDBTestInstructionContainerItem fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerMessage

	// Extract data from DB result set
	for rows.Next() {
		err := rows.Scan(
			&cloudDBTestInstructionContainerItem.DomainUuid,
			&cloudDBTestInstructionContainerItem.DomainName,
			&cloudDBTestInstructionContainerItem.TestInstructionContainerUuid,
			&cloudDBTestInstructionContainerItem.TestInstructionContainerName,
			&cloudDBTestInstructionContainerItem.TestInstructionContainerTypeUuid,
			&cloudDBTestInstructionContainerItem.TestInstructionContainerTypeName,
			&cloudDBTestInstructionContainerItem.TestInstructionContainerDescription,
			&cloudDBTestInstructionContainerItem.TestInstructionContainerMouseOverText,
			&cloudDBTestInstructionContainerItem.Deprecated,
			&cloudDBTestInstructionContainerItem.Enabled,
			&cloudDBTestInstructionContainerItem.MajorVersionNumber,
			&cloudDBTestInstructionContainerItem.MinorVersionNumber,
			&cloudDBTestInstructionContainerItem.UpdatedTimeStamp,
		)

		if err != nil {
			return err
		}

		// Add values to the object that is pointed to by variable in function
		*cloudDBTestInstructionContainerItems = append(*cloudDBTestInstructionContainerItems, &cloudDBTestInstructionContainerItem)

	}

	// No errors occurred
	return nil

}
