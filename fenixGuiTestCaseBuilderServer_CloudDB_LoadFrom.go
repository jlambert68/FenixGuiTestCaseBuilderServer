package main

import (
	"context"
	"fmt"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// ****************************************************************************************************************
// Load data from CloudDB
//
/*
// Load TestInstructions and pre-created TestInstructionContainers for Client
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) loadClientsTestInstructionsFromCloudDB(userID string, cloudDBTestInstructionItems *[]*fenixTestCaseBuilderServerGrpcApi.TestInstructionMessage) (err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "61b8b021-9568-463e-b867-ac1ddb10584d",
	}).Debug("Entering: loadClientsTestInstructionsFromCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "78a97c41-a098-4122-88d2-01ed4b6c4844",
		}).Debug("Exiting: loadClientsTestInstructionsFromCloudDB()")
	}()

	/* Example
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

	//* /

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

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
	var tempTimeStamp time.Time

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
			&tempTimeStamp,
		)

		// Convert TimeStamp into proto-format for TimeStamp
		cloudDBTestInstructionItem.UpdatedTimeStamp = timestamppb.New(tempTimeStamp)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                         "e7925b78-327c-40ad-9144-ae4a8a6f35f5",
				"Error":                      err,
				"sqlToExecute":               sqlToExecute,
				"cloudDBTestInstructionItem": cloudDBTestInstructionItem,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// Add values to the object that is pointed to by variable in function
		*cloudDBTestInstructionItems = append(*cloudDBTestInstructionItems, &cloudDBTestInstructionItem)

	}

	// No errors occurred
	return nil

}
*/
// Load TestInstructions and pre-created TestInstructionContainers for Client
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) loadClientsImmatureTestInstructionsFromCloudDB(userID string, cloudDBImmatureTestInstructionItems *[]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage) (err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "38fbd4e2-cfe8-405c-84ce-1667c2292c58",
	}).Debug("Entering: loadClientsImmatureTestInstructionsFromCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "6acee974-1280-48f5-9c4f-886aeff58863",
		}).Debug("Exiting: loadClientsImmatureTestInstructionsFromCloudDB()")
	}()

	var (
		basicTestInstructionInformation            fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage
		basicTestInstructionInformationSQLCount    int64
		immatureTestInstructionInformation         fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage
		immatureTestInstructionInformationSQLCount int64
		//immatureSubTestCaseModel                   fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
		//immatureSubTestCaseModelSQLCount           int64
	)

	ImmatureTestInstructionMessageMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage)

	/* Example
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

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	// **** BasicTestInstructionInformation ****
	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT * "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"BasicTestInstructionInformation\" BTII_TI "
	sqlToExecute = sqlToExecute + "ORDER BY \"TestInstructionUuid\" ASC;"

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
	//var basicTestInstructionInformation fenixTestCaseBuilderServerGrpcApi.TestInstructionMessage
	var tempTimeStamp time.Time

	// Get number of rows for 'basicTestInstructionInformation'
	basicTestInstructionInformationSQLCount = rows.CommandTag().RowsAffected()
	var (
		nonEditableInformation    fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_NonEditableBasicInformationMessage
		editableInformation       fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_EditableBasicInformationMessage
		invisibleBasicInformation fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_InvisibleBasicInformationMessage
	)

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		nonEditableInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_NonEditableBasicInformationMessage{}
		editableInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_EditableBasicInformationMessage{}
		invisibleBasicInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_InvisibleBasicInformationMessage{}

		err := rows.Scan(
			// NonEditableInformation
			&nonEditableInformation.DomainUuid,
			&nonEditableInformation.DomainName,
			&nonEditableInformation.TestInstructionUuid,
			&nonEditableInformation.TestInstructionName,
			&nonEditableInformation.TestInstructionTypeUuid,
			&nonEditableInformation.TestInstructionTypeName,
			&nonEditableInformation.Deprecated,
			&nonEditableInformation.MajorVersionNumber,
			&nonEditableInformation.MinorVersionNumber,
			&tempTimeStamp,
			&nonEditableInformation.TestInstructionColor,
			&nonEditableInformation.TCRuleDeletion,
			&nonEditableInformation.TCRuleSwap,

			// EditableInformation
			&editableInformation.TestInstructionDescription,
			&editableInformation.TestInstructionMouseOverText,

			// InvisibleBasicInformation
			&invisibleBasicInformation.Enabled,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "e7925b78-327c-40ad-9144-ae4a8a6f35f5",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// Convert TimeStamp into proto-format for TimeStamp
		nonEditableInformation.UpdatedTimeStamp = timestamppb.New(tempTimeStamp)

		// Add 'basicTestInstructionInformation' to map
		testInstructionUuid := nonEditableInformation.TestInstructionUuid
		immatureTestInstructionMessage, existsInMap := ImmatureTestInstructionMessageMap[testInstructionUuid]
		// testInstructionUuid shouldn't exist in map. If so then there is a problem
		if existsInMap == true {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                  "4713a8c8-c9d0-4315-9341-27365d64cdc8",
				"testInstructionUuid": testInstructionUuid,
				"sqlToExecute":        sqlToExecute,
			}).Fatal("TestInstructionUuid shouldn't exist in map. If so then there is a problem")

		}

		// Create 'basicTestInstructionInformation' of the parts
		basicTestInstructionInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage{
			NonEditableInformation:    &nonEditableInformation,
			EditableInformation:       &editableInformation,
			InvisibleBasicInformation: &invisibleBasicInformation,
		}

		// Create 'immatureTestInstructionMessage' and add 'BasicTestInstructionInformation'
		immatureTestInstructionMessage = fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage{BasicTestInstructionInformation: &basicTestInstructionInformation}

		// Save immatureTestInstructionMessage in map
		ImmatureTestInstructionMessageMap[testInstructionUuid] = immatureTestInstructionMessage

	}

	// **** immatureTestInstructionInformation ****
	sqlToExecute = ""
	sqlToExecute = sqlToExecute + "SELECT * "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"ImmatureTestInstructionInformation\" ITII_TI "
	sqlToExecute = sqlToExecute + "ORDER BY \"TestInstructionUuid\" ASC,  \"DropZoneUuid\" ASC, \"TestInstructionAttributeGuid\" ASC; "

	// Query DB
	rows, err = fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "b3ef4fec-9097-46c4-8ff6-85a758967e46",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Get number of rows for 'immatureTestInstructionInformation'
	immatureTestInstructionInformationSQLCount = rows.CommandTag().RowsAffected()

	// Create map to store ImmatureTestInstructionInformationMessages
	immatureTestInstructionInformationMessagesMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage)

	// Temp variables used when extracting data
	var domainUuid, previousDomainUuid string
	var domainName string
	var testInstructionUuid, previousTestInstructionUuid string
	var testInstructionName string
	var tempTestInstructionAttributeType string
	// First Row in TestData
	var firstRowInSQLRespons bool
	firstRowInSQLRespons = true

	var (
		availableDropZone, previousAvailableDropZone fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
		availableDropZones                           []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
	)

	var (
		dropZonePreSetTestInstructionAttribute, previousDropZonePreSetTestInstructionAttribute fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
		dropZonePreSetTestInstructionAttributes                                                []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	)

	var firstImmatureElementUuid string

	var dataStateChange uint8

	// Clear previous variables
	previousDomainUuid = ""
	previousTestInstructionUuid = ""

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		availableDropZone = fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage{}
		dropZonePreSetTestInstructionAttribute = fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}

		err := rows.Scan(

			// temp-data which is not stored in object
			&domainUuid,
			&domainName,
			&testInstructionUuid,
			&testInstructionName,

			// DropZone-data
			&availableDropZone.DropZoneUuid,
			&availableDropZone.DropZoneName,
			&availableDropZone.DropZoneDescription,
			&availableDropZone.DropZoneMouseOver,
			&availableDropZone.DropZoneColor,

			// DropZoneAttributes-data
			&tempTestInstructionAttributeType,
			&dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid,
			&dropZonePreSetTestInstructionAttribute.TestInstructionAttributeName,
			&dropZonePreSetTestInstructionAttribute.AttributeValueAsString,
			&dropZonePreSetTestInstructionAttribute.AttributeValueUuid,

			// Reference to first element in element-model
			&firstImmatureElementUuid,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "9f0618f2-ca04-41e9-aeef-60cd1874f6b7",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// Convert 'tempTestInstructionAttributeType' into gRPC-type
		dropZonePreSetTestInstructionAttribute.TestInstructionAttributeType = fenixTestCaseBuilderServerGrpcApi.TestInstructionAttributeTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestInstructionAttributeTypeEnum_value[tempTestInstructionAttributeType])

		// Handle the correct order of building together the full object
		dataStateChange = 0

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		dataStateChangeFound :=
			firstRowInSQLRespons == true &&
				domainUuid != previousDomainUuid &&
				testInstructionUuid != previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 1
		}

		// All UUIDs are changed and this is not the first row [dataStateChange=2]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid != previousDomainUuid &&
				testInstructionUuid != previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 2
		}

		// Only DropZonePreSetTestInstructionAttributeUuid is changed and this is not the first row [dataStateChange=3]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid == previousDomainUuid &&
				testInstructionUuid == previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid == previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 3
		}

		// Only AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=4]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid == previousDomainUuid &&
				testInstructionUuid == previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 4
		}

		// Only TestInstructionUuid, AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=5]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid == previousDomainUuid &&
				testInstructionUuid != previousTestInstructionUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 5
		}

		// Act on which 'dataStateChange' that was achieved
		switch dataStateChange {

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		case 1:

		// All UUIDs are changed and this is not the first row [dataStateChange=2]
		// Only TestInstructionUuid, AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=5]
		case 2, 5:
			// New DropZone so add the previous DropZone-attributes to the DropZone-array
			dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, &previousDropZonePreSetTestInstructionAttribute)

			// Add attributes to previousDropZone
			previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributes

			// Add previousAvailableDropZone to array of DropZone
			availableDropZones = append(availableDropZones, &previousAvailableDropZone)

			// Add the availableDropZones to the ImmatureTestInstructionInformationMessage-map
			immatureTestInstructionInformation.AvailableDropZones = availableDropZones
			immatureTestInstructionInformationMessagesMap[previousTestInstructionUuid] = immatureTestInstructionInformation

			// Create fresh versions of variables
			immatureTestInstructionInformation = fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage{}
			availableDropZones = []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage{}
			dropZonePreSetTestInstructionAttributes = []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}

		// Only DropZonePreSetTestInstructionAttributeUuid is changed and this is not the first row [dataStateChange=3]
		case 3:
			// Add the DropZone attribute to the array for attributes
			dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, &previousDropZonePreSetTestInstructionAttribute)

		// Only AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=4]
		case 4:
			// New DropZone so add the previous DropZone-attributes to the DropZone-array
			dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, &previousDropZonePreSetTestInstructionAttribute)

			// Add attributes to previousDropZone
			previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributes

			// Add previousAvailableDropZone to array of DropZone
			availableDropZones = append(availableDropZones, &previousAvailableDropZone)

			// Create fresh versions of variables
			dropZonePreSetTestInstructionAttributes = []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}

			// Something is wrong in the ordering of the testdata or the testdata itself
		default:
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                                     "ca46deb0-a788-4c68-aefb-27bb7ccaad0d",
				"domainUuid":                             domainUuid,
				"previousDomainUuid":                     previousDomainUuid,
				"testInstructionUuid":                    testInstructionUuid,
				"previousTestInstructionUuid":            previousTestInstructionUuid,
				"availableDropZone.DropZoneUuid":         availableDropZone.DropZoneUuid,
				"previousAvailableDropZone.DropZoneUuid": previousAvailableDropZone.DropZoneUuid,
				"dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid":         dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid,
				"previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid": previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid,
			}).Fatal("Something is wrong in the ordering of the testdata or the testdata itself  --> Should bot happen")

		}

		// Move actual values into previous-variables
		previousDomainUuid = domainUuid
		previousTestInstructionUuid = testInstructionUuid
		previousAvailableDropZone = availableDropZone
		previousDropZonePreSetTestInstructionAttribute = dropZonePreSetTestInstructionAttribute

	}

	// Handle last row from database
	// Add the previous DropZone-attributes to the DropZone-array
	dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, &previousDropZonePreSetTestInstructionAttribute)

	// Add attributes to previousDropZone
	previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributes

	// Add previousAvailableDropZone to array of DropZone
	availableDropZones = append(availableDropZones, &previousAvailableDropZone)

	// Add the availableDropZones to the ImmatureTestInstructionInformationMessage-map
	immatureTestInstructionInformation.AvailableDropZones = availableDropZones
	immatureTestInstructionInformationMessagesMap[previousTestInstructionUuid] = immatureTestInstructionInformation

	// Add 'basicTestInstructionInformation' to map
	immatureTestInstructionMessage, existsInMap := ImmatureTestInstructionMessageMap[testInstructionUuid]
	// testInstructionUuid shouldn't exist in map. If so then there is a problem
	if existsInMap == false {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":                  "4894cf70-08fd-401b-b076-d643ea721abb",
			"testInstructionUuid": testInstructionUuid,
		}).Fatal("TestInstructionUuid should exist in map. If not then there is a problem")
	}

	// Immature part to 'immatureTestInstructionMessage'
	immatureTestInstructionMessage.ImmatureTestInstructionInformation = &immatureTestInstructionInformation

	// Store the information back in the map
	ImmatureTestInstructionMessageMap[testInstructionUuid] = immatureTestInstructionMessage

	fmt.Println(basicTestInstructionInformationSQLCount)
	fmt.Println(immatureTestInstructionInformationSQLCount)

	// No errors occurred
	return nil

}

/*
// ****************************************************************************************************************
// Load data from CloudDB into memory structures
//
// Load pre-created TestInstructionContainerContainers for Client
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) loadClientsTestInstructionContainersFromCloudDB(userID string, cloudDBTestInstructionContainerItems *[]*fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerMessage) (err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "f91a7e85-d5df-42f5-80ff-a65b8350467f",
	}).Debug("Entering: loadClientsTestInstructionContainersFromCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "40ccee29-d32e-4674-9a3a-fd4403b55d32",
		}).Debug("Exiting: loadClientsTestInstructionContainersFromCloudDB()")
	}()

	/* Example

	   "DomainUuid"                            uuid      not null,
	   "DomainName"                            varchar   not null,
	   "TestInstructionContainerUuid"          uuid      not null
	   "TestInstructionContainerName"          varchar   not null,
	   "TestInstructionContainerTypeUuid"      uuid      not null,
	   "TestInstructionContainerTypeName"      varchar   not null,
	   "TestInstructionContainerDescription"   varchar   not null,
	   "TestInstructionContainerMouseOverText" varchar   not null,
	   "Deprecated"                            boolean   not null,
	   "Enabled"                               boolean   not null,
	   "MajorVersionNumber"                    integer   not null,
	   "MinorVersionNumber"                    integer   not null,
	   "UpdatedTimeStamp"                      timestamp not null,
	   "ChildrenIsParallelProcessed"           boolean   not null



	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT * "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestInstructionContainers\" FGB_TIC;"

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "b54c3ae1-9d96-4f00-9bc3-2c1a1712b91a",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Variables to used when extract data from result set
	var tempTimeStamp time.Time
	var childrenIsParallelProcessed bool

	// Extract data from DB result set
	for rows.Next() {

		// Define for every loop because otherwise the same object is referenced in array
		var cloudDBTestInstructionContainerItem fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerMessage

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
			&tempTimeStamp,
			&childrenIsParallelProcessed,
		)

		// Convert TimeStamp into proto-format for TimeStamp
		cloudDBTestInstructionContainerItem.UpdatedTimeStamp = timestamppb.New(tempTimeStamp)

		// Convert 'childrenIsParallelProcessed' into Proto-message-format
		if childrenIsParallelProcessed == true {
			// Children executed in Parallel
			cloudDBTestInstructionContainerItem.TestInstructionContainerExecutionType = fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerExecutionTypeEnum_PARALLELLED_PROCESSED

		} else {
			// Children executed in Serial
			cloudDBTestInstructionContainerItem.TestInstructionContainerExecutionType = fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerExecutionTypeEnum_SERIAL_PROCESSED

		}

		if err != nil {
			return err
		}

		// TODO Load children
		cloudDBTestInstructionContainerItem.TestInstructionContainerChildren = nil

		// Add values to the object that is pointed to by variable in function
		*cloudDBTestInstructionContainerItems = append(*cloudDBTestInstructionContainerItems, &cloudDBTestInstructionContainerItem)

	}

	// No errors occurred
	return nil

}
*/
/*
// Load TestInstructions and pre-created TestInstructionContainers for Client
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) loadPinnedClientsTestInstructionsFromCloudDB(userID string, cloudDBTestInstructionItems *[]*fenixTestCaseBuilderServerGrpcApi.TestInstructionMessage) (err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "61aed366-5342-4f33-8bde-99edf990d143",
	}).Debug("Entering: loadPinnedClientsTestInstructionsFromCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "7586abec-3f6a-4fcf-97fc-73c50543a18c",
		}).Debug("Exiting: loadPinnedClientsTestInstructionsFromCloudDB()")
	}()

	/* Example
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



	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT FGB_TI.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestInstructions\" FGB_TI "
	sqlToExecute = sqlToExecute + "WHERE FGB_TI.\"TestInstructionUuid\" IN (SELECT FGB_PTIC.\"PinnedUuid\" "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"PinnedTestInstructionsAndPreCreatedTestInstructionContainers\" FGB_PTIC "
	sqlToExecute = sqlToExecute + "WHERE FGB_PTIC.\"UserId\" = '" + userID + "' AND FGB_PTIC.\"PinnedType\" = 1);"

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "b6295054-3c4e-427e-b2e2-55bf69d89a20",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Variables to used when extract data from result set
	var cloudDBTestInstructionItem fenixTestCaseBuilderServerGrpcApi.TestInstructionMessage
	var tempTimeStamp time.Time

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
			&tempTimeStamp,
		)

		// Convert TimeStamp into proto-format for TimeStamp
		cloudDBTestInstructionItem.UpdatedTimeStamp = timestamppb.New(tempTimeStamp)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                         "e7925b78-327c-40ad-9144-ae4a8a6f35f5",
				"Error":                      err,
				"sqlToExecute":               sqlToExecute,
				"cloudDBTestInstructionItem": cloudDBTestInstructionItem,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// Add values to the object that is pointed to by variable in function
		*cloudDBTestInstructionItems = append(*cloudDBTestInstructionItems, &cloudDBTestInstructionItem)

	}

	// No errors occurred
	return nil

}
*/
/*
// ****************************************************************************************************************
// Load data from CloudDB into memory structures
//
// Load pre-created TestInstructionContainerContainers for Client
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) loadPinnedClientsTestInstructionContainersFromCloudDB(userID string, cloudDBTestInstructionContainerItems *[]*fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerMessage) (err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "99ed695b-a52f-4260-9b45-49a9c33f9470",
	}).Debug("Entering: loadPinnedClientsTestInstructionContainersFromCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "5b6eb21d-21cc-4457-9576-d1b6357c851e",
		}).Debug("Exiting: loadPinnedClientsTestInstructionContainersFromCloudDB()")
	}()

	/* Example

	   "DomainUuid"                            uuid      not null,
	   "DomainName"                            varchar   not null,
	   "TestInstructionContainerUuid"          uuid      not null
	   "TestInstructionContainerName"          varchar   not null,
	   "TestInstructionContainerTypeUuid"      uuid      not null,
	   "TestInstructionContainerTypeName"      varchar   not null,
	   "TestInstructionContainerDescription"   varchar   not null,
	   "TestInstructionContainerMouseOverText" varchar   not null,
	   "Deprecated"                            boolean   not null,
	   "Enabled"                               boolean   not null,
	   "MajorVersionNumber"                    integer   not null,
	   "MinorVersionNumber"                    integer   not null,
	   "UpdatedTimeStamp"                      timestamp not null,
	   "ChildrenIsParallelProcessed"           boolean   not null



	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT FGB_TIC.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"TestInstructionContainers\" FGB_TIC "
	sqlToExecute = sqlToExecute + "WHERE FGB_TIC.\"TestInstructionContainerUuid\" IN (SELECT FGB_PTIC.\"PinnedUuid\" "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"PinnedTestInstructionsAndPreCreatedTestInstructionContainers\" FGB_PTIC "
	sqlToExecute = sqlToExecute + "WHERE FGB_PTIC.\"UserId\" = '" + userID + "' AND FGB_PTIC.\"PinnedType\" = 2);"

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "fb560958-9082-483c-950c-95267d40f507",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Variables to used when extract data from result set
	var tempTimeStamp time.Time
	var childrenIsParallelProcessed bool

	// Extract data from DB result set
	for rows.Next() {

		// Define for every loop because otherwise the same object is referenced in array
		var cloudDBTestInstructionContainerItem fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerMessage

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
			&tempTimeStamp,
			&childrenIsParallelProcessed,
		)

		// Convert TimeStamp into proto-format for TimeStamp
		cloudDBTestInstructionContainerItem.UpdatedTimeStamp = timestamppb.New(tempTimeStamp)

		// Convert 'childrenIsParallelProcessed' into Proto-message-format
		if childrenIsParallelProcessed == true {
			// Children executed in Parallel
			cloudDBTestInstructionContainerItem.TestInstructionContainerExecutionType = fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerExecutionTypeEnum_PARALLELLED_PROCESSED

		} else {
			// Children executed in Serial
			cloudDBTestInstructionContainerItem.TestInstructionContainerExecutionType = fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerExecutionTypeEnum_SERIAL_PROCESSED

		}

		if err != nil {
			return err
		}

		// TODO Load children
		cloudDBTestInstructionContainerItem.TestInstructionContainerChildren = nil

		// Add values to the object that is pointed to by variable in function
		*cloudDBTestInstructionContainerItems = append(*cloudDBTestInstructionContainerItems, &cloudDBTestInstructionContainerItem)

	}

	// No errors occurred
	return nil

}
*/
