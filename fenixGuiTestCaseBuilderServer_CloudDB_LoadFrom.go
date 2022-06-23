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

// Load TestInstructions for Client
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) loadClientsImmatureTestInstructionsFromCloudDB(userID string) (cloudDBImmatureTestInstructionItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage, err error) {

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

	// **** BasicTestInstructionInformation **** **** BasicTestInstructionInformation **** **** BasicTestInstructionInformation ****
	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT BTI.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"BasicTestInstructionInformation\" BTI "
	sqlToExecute = sqlToExecute + "ORDER BY BTI.\"DomainUuid\" ASC,  BTI.\"TestInstructionTypeUuid\" ASC, BTI.\"TestInstructionUuid\" ASC; "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "2f130d7e-f8aa-466f-b29d-0fb63608c1a6",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Variables to used when extract data from result set
	//var basicTestInstructionInformation fenixTestCaseBuilderServerGrpcApi.TestInstructionMessage
	var tempTimeStamp time.Time

	// Get number of rows for 'basicTestInstructionInformation'
	basicTestInstructionInformationSQLCount = rows.CommandTag().RowsAffected()
	var (
		nonEditableInformation      fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_NonEditableBasicInformationMessage
		editableInformation         fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_EditableBasicInformationMessage
		invisibleBasicInformation   fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_InvisibleBasicInformationMessage
		immatureElementModelMessage fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
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

			return nil, err
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

		// Create 'immatureTestInstructionMessage' and add 'BasicTestInstructionInformation' and a small part of 'ImmatureSubTestCaseModel'
		immatureTestInstructionMessage = fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage{
			BasicTestInstructionInformation: &basicTestInstructionInformation,
			ImmatureSubTestCaseModel:        &immatureElementModelMessage}

		// Save immatureTestInstructionMessage in map
		ImmatureTestInstructionMessageMap[testInstructionUuid] = immatureTestInstructionMessage

	}

	// **** immatureTestInstructionInformation **** **** immatureTestInstructionInformation **** **** immatureTestInstructionInformation ****
	sqlToExecute = ""
	sqlToExecute = sqlToExecute + "SELECT * "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"ImmatureTestInstructionInformation\" ITII "
	sqlToExecute = sqlToExecute + "ORDER BY ITII.\"DomainUuid\" ASC, ITII.\"TestInstructionUuid\" ASC,  ITII.\"DropZoneUuid\" ASC, ITII.\"TestInstructionAttributeGuid\" ASC; "

	// Query DB
	rows, err = fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "b3ef4fec-9097-46c4-8ff6-85a758967e46",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
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

			return nil, err
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

		// Set to be first row
		firstRowInSQLRespons = false

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

	// Add 'firstImmatureElementUuid'
	immatureTestInstructionMessage.ImmatureSubTestCaseModel.FirstImmatureElementUuid = firstImmatureElementUuid

	// Store the information back in the map
	ImmatureTestInstructionMessageMap[testInstructionUuid] = immatureTestInstructionMessage

	// ***************************************************************************************************

	// **** ImmatureElementModelMessage **** **** ImmatureElementModelMessage **** **** ImmatureElementModelMessage ****

	sqlToExecute = ""
	sqlToExecute = sqlToExecute + "SELECT IEM.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"BasicTestInstructionInformation\" BTI, "
	sqlToExecute = sqlToExecute + "\"" + usedDBSchema + "\".\"ImmatureElementModelMessage\" IEM "
	sqlToExecute = sqlToExecute + "WHERE BTI.\"TestInstructionUuid\" = IEM.\"ImmatureElementUuid\" "
	sqlToExecute = sqlToExecute + "ORDER BY BTI.\"DomainUuid\" ASC, BTI.\"TestInstructionUuid\" ASC; "

	// Query DB
	rows, err = fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "b3ef4fec-9097-46c4-8ff6-85a758967e46",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Get number of rows for 'immatureTestInstructionInformation'
	immatureTestInstructionInformationSQLCount = rows.CommandTag().RowsAffected()

	// Create map to store ImmatureTestInstructionInformationMessages
	//immatureTestInstructionInformationMessagesMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage)

	// Temp variables used when extracting data
	var tempImmatureElementModelDomainUuid, previousTempImmatureDomainUuid string
	var tempImmatureElementModelDomainName string
	var tempTestCaseModelElementTypeAsString string
	//var previousOriginalElementUuid string
	//var testInstructionUuid, previousTestInstructionUuid string
	//var testInstructionName string
	//var tempTestInstructionAttributeType string
	// First Row in TestData
	//var firstRowInSQLRespons bool
	firstRowInSQLRespons = true

	var (
	//availableDropZone, previousAvailableDropZone fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
	//availableDropZones                           []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
	)

	var (
	//dropZonePreSetTestInstructionAttribute, previousDropZonePreSetTestInstructionAttribute fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	//dropZonePreSetTestInstructionAttributes                                                []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	)

	//var immatureElementModelMessage fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
	var immatureElementModelElement, previousImmatureElementModelElement fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage
	var immatureElementModelElements []*fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage

	//var firstImmatureElementUuid string

	//var dataStateChange uint8

	// Clear previous variables
	previousDomainUuid = ""
	previousTestInstructionUuid = ""

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		immatureElementModelElement = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage{}

		err = rows.Scan(

			// temp-data which is not stored in object
			&tempImmatureElementModelDomainUuid,
			&tempImmatureElementModelDomainName,

			// ImmatureElementModel

			&immatureElementModelElement.OriginalElementUuid,
			&immatureElementModelElement.OriginalElementName,
			&immatureElementModelElement.PreviousElementUuid,
			&immatureElementModelElement.NextElementUuid,
			&immatureElementModelElement.FirstChildElementUuid,
			&immatureElementModelElement.ParentElementUuid,
			&tempTestCaseModelElementTypeAsString,
			&immatureElementModelElement.CurrentElementModelElement,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "808377e2-70ec-4894-bb17-7d92321caaa2",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert 'tempTestCaseModelElementTypeAsString' into gRPC-type
		immatureElementModelElement.TestCaseModelElementType = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum_value[tempTestCaseModelElementTypeAsString])

		// Handle the correct order of building together the full object
		dataStateChange = 0

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		dataStateChangeFound :=
			firstRowInSQLRespons == true &&
				tempImmatureElementModelDomainUuid != previousTempImmatureDomainUuid &&
				immatureElementModelElement.OriginalElementUuid != previousImmatureElementModelElement.OriginalElementUuid
		if dataStateChangeFound == true {
			dataStateChange = 1
		}

		// All UUIDs are changed and this is not the first row [dataStateChange=2]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				tempImmatureElementModelDomainUuid != previousTempImmatureDomainUuid &&
				immatureElementModelElement.OriginalElementUuid != previousImmatureElementModelElement.OriginalElementUuid
		if dataStateChangeFound == true {
			dataStateChange = 2
		}

		// Only immatureElementModelElement.OriginalElementUuid is changed and this is not the first row [dataStateChange=3]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid == previousDomainUuid &&
				tempImmatureElementModelDomainUuid == previousTempImmatureDomainUuid &&
				immatureElementModelElement.OriginalElementUuid != previousImmatureElementModelElement.OriginalElementUuid
		if dataStateChangeFound == true {
			dataStateChange = 3
		}

		// Act on which 'dataStateChange' that was achieved
		switch dataStateChange {

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		case 1:

			// All UUIDs are changed and this is not the first row [dataStateChange=2]
			// Only immatureElementModelElement.OriginalElementUuid is changed and this is not the first row [dataStateChange=3]
		case 2, 5:
			// New ElementModelElement so add the previous ElementModelElement to the ElementModelElements-array
			immatureElementModelElements = append(immatureElementModelElements, &previousImmatureElementModelElement)

			// Add immatureElementModelElements to 'immatureTestInstructionMessage.ImmatureSubTestCaseModel' which can be found in map
			immatureTestInstructionMessage, existsInMap = ImmatureTestInstructionMessageMap[previousImmatureElementModelElement.OriginalElementUuid]
			// testInstructionUuid shouldn exist in map. If not so then there is a problem
			if existsInMap == false {
				fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
					"Id": "c2d2448e-9e86-4947-a5f5-7787a72e7ef8",
					"previousImmatureElementModelElement.OriginalElementUuid": previousImmatureElementModelElement.OriginalElementUuid,
				}).Fatal("TestInstructionUuid should exist in map. If not then there is a problem")
			}

			immatureElementModelMessage.TestCaseModelElements = immatureElementModelElements
			immatureTestInstructionMessage.ImmatureSubTestCaseModel = &immatureElementModelMessage

			// Create fresh versions of variables
			immatureElementModelMessage = fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage{}

		// Only immatureElementModelElement.OriginalElementUuid is changed and this is not the first row [dataStateChange=3]
		case 3:
			// New ElementModelElement so add the previous ElementModelElement to the ElementModelElements-array
			immatureElementModelElements = append(immatureElementModelElements, &previousImmatureElementModelElement)

			// Something is wrong in the ordering of the testdata or the testdata itself
		default:
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                          "e464567a-79ab-49e7-9519-ce187b34458d",
				"domainUuid":                  domainUuid,
				"previousDomainUuid":          previousDomainUuid,
				"testInstructionUuid":         testInstructionUuid,
				"previousTestInstructionUuid": previousTestInstructionUuid,
			}).Fatal("Something is wrong in the ordering of the testdata or the testdata itself  --> Should bot happen")

			// Move actual values into previous-variables
			previousImmatureElementModelElement = immatureElementModelElement

		}

		// Move previous values to current
		previousImmatureElementModelElement = immatureElementModelElement
		previousTempImmatureDomainUuid = tempImmatureElementModelDomainUuid

		// Set to not be first row
		firstRowInSQLRespons = false

	}
	// Handle last row from database

	// New ElementModelElement so add the previous ElementModelElement to the ElementModelElements-array
	immatureElementModelElements = append(immatureElementModelElements, &immatureElementModelElement)

	// Add immatureElementModelElements to 'immatureTestInstructionMessage.ImmatureSubTestCaseModel' which can be found in map
	immatureTestInstructionMessage, existsInMap = ImmatureTestInstructionMessageMap[immatureElementModelElement.OriginalElementUuid]
	// testInstructionUuid shouldn exist in map. If not so then there is a problem
	if existsInMap == false {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "c2d2448e-9e86-4947-a5f5-7787a72e7ef8",
			"immatureElementModelElement.OriginalElementUuid": immatureElementModelElement.OriginalElementUuid,
		}).Fatal("TestInstructionUuid should exist in map. If not then there is a problem")
	}

	immatureElementModelMessage.TestCaseModelElements = immatureElementModelElements
	immatureTestInstructionMessage.ImmatureSubTestCaseModel = &immatureElementModelMessage
	ImmatureTestInstructionMessageMap[immatureElementModelElement.OriginalElementUuid] = immatureTestInstructionMessage

	// Loop all ImmatureTestInstructionMessage and create gRPC-response
	var allImmatureTestInstructionMessage []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage

	for _, value := range ImmatureTestInstructionMessageMap { // Order not specified
		allImmatureTestInstructionMessage = append(allImmatureTestInstructionMessage, &value)
	}

	cloudDBImmatureTestInstructionItems = allImmatureTestInstructionMessage

	fmt.Println(basicTestInstructionInformationSQLCount)
	fmt.Println(immatureTestInstructionInformationSQLCount)

	// No errors occurred
	return cloudDBImmatureTestInstructionItems, nil

}

// Load TestInstructionContainers for Client
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) loadClientsImmatureTestInstructionContainersFromCloudDB(userID string) (cloudDBImmatureTestInstructionContainerItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage, err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "68b965ea-234c-425b-b525-1f8b7154850b",
	}).Debug("Entering: loadClientsImmatureTestInstructionContainersFromCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "12021bfa-154f-48f2-bd8c-0809e1877fd4",
		}).Debug("Exiting: loadClientsImmatureTestInstructionContainersFromCloudDB()")
	}()

	var (
		//	basicTestInstructionContainerInformation            fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage
		basicTestInstructionContainerInformationSQLCount    int64
		immatureTestInstructionContainerInformation         fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage
		immatureTestInstructionContainerInformationSQLCount int64
		//immatureSubTestCaseModel                   fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
		//immatureSubTestCaseModelSQLCount           int64
	)

	ImmatureTestInstructionContainerMessageMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage)

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

	// **** BasicTestInstructionContainerInformation **** **** BasicTestInstructionContainerInformation **** **** BasicTestInstructionContainerInformation ****
	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT BTIC.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"BasicTestInstructionContainerInformation\" BTIC "
	sqlToExecute = sqlToExecute + "ORDER BY BTIC.\"DomainUuid\" ASC,  BTIC.\"TestInstructionContainerTypeUuid\" ASC, BTIC.\"TestInstructionContainerUuid\" ASC; "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "b944c506-4ded-4f5e-98c4-06f272d16e1a",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Variables to used when extract data from result set
	//var basicTestInstructionContainerInformation fenixTestCaseBuilderServerGrpcApi.TestInstructionMessage
	var tempTimeStamp time.Time
	var tempTestInstructionContainerExecutionType string

	// Get number of rows for 'basicTestInstructionContainerInformation'
	basicTestInstructionContainerInformationSQLCount = rows.CommandTag().RowsAffected()
	var (
		nonEditableInformation                    fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_NonEditableBasicInformationMessage
		editableInformation                       fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableBasicInformationMessage
		invisibleBasicInformation                 fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_InvisibleBasicInformationMessage
		editableTestInstructionContainerAttribute fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableTestInstructionContainerAttributesMessage
		//immatureElementModelMessage                        fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
		//immatureTestInstructionContainerInformationMessage fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage
	)

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		nonEditableInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_NonEditableBasicInformationMessage{}
		editableInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableBasicInformationMessage{}
		invisibleBasicInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_InvisibleBasicInformationMessage{}
		editableTestInstructionContainerAttribute = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableTestInstructionContainerAttributesMessage{}

		err := rows.Scan(
			// NonEditableInformation
			&nonEditableInformation.DomainUuid,
			&nonEditableInformation.DomainName,
			&nonEditableInformation.TestInstructionContainerUuid,
			&nonEditableInformation.TestInstructionContainerName,
			&nonEditableInformation.TestInstructionContainerTypeUuid,
			&nonEditableInformation.TestInstructionContainerTypeName,
			&nonEditableInformation.Deprecated,
			&nonEditableInformation.MajorVersionNumber,
			&nonEditableInformation.MinorVersionNumber,
			&tempTimeStamp,
			&nonEditableInformation.TestInstructionContainerColor,
			&nonEditableInformation.TCRuleDeletion,
			&nonEditableInformation.TCRuleSwap,

			// EditableInformation
			&editableInformation.TestInstructionContainerDescription,
			&editableInformation.TestInstructionContainerMouseOverText,

			// InvisibleBasicInformation
			&invisibleBasicInformation.Enabled,

			// EditableTestInstructionContainerAttribute
			&tempTestInstructionContainerExecutionType,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "7d082f7c-f987-44e7-97b7-c3c1652955c3",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert TimeStamp into proto-format for TimeStamp
		nonEditableInformation.UpdatedTimeStamp = timestamppb.New(tempTimeStamp)

		// Convert 'tempTestInstructionContainerExecutionType' gRPC-type
		editableTestInstructionContainerAttribute.TestInstructionContainerExecutionType = fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerExecutionTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerExecutionTypeEnum_value[tempTestInstructionContainerExecutionType])

		// Add 'basicTestInstructionContainerInformation' to map
		testInstructionContainerUuid := nonEditableInformation.TestInstructionContainerUuid
		immatureTestInstructionContainerMessage, existsInMap := ImmatureTestInstructionContainerMessageMap[testInstructionContainerUuid]
		// testInstructionContainerUuid shouldn't exist in map. If so then there is a problem
		if existsInMap == true {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                           "58cd4928-e4b5-4faf-9724-047c1cbc82a1",
				"testInstructionContainerUuid": testInstructionContainerUuid,
				"sqlToExecute":                 sqlToExecute,
			}).Fatal("TestInstructionContainerUuid shouldn't exist in map. If so then there is a problem")

		}

		// Create 'basicTestInstructionContainerInformation' of the parts
		basicTestInstructionContainerInformation := fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage{
			NonEditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_NonEditableBasicInformationMessage{
				DomainUuid:                       nonEditableInformation.DomainUuid,
				DomainName:                       nonEditableInformation.DomainName,
				TestInstructionContainerUuid:     nonEditableInformation.TestInstructionContainerUuid,
				TestInstructionContainerName:     nonEditableInformation.TestInstructionContainerName,
				TestInstructionContainerTypeUuid: nonEditableInformation.TestInstructionContainerTypeUuid,
				TestInstructionContainerTypeName: nonEditableInformation.TestInstructionContainerTypeName,
				Deprecated:                       nonEditableInformation.Deprecated,
				MajorVersionNumber:               nonEditableInformation.MajorVersionNumber,
				MinorVersionNumber:               nonEditableInformation.MinorVersionNumber,
				UpdatedTimeStamp:                 nonEditableInformation.UpdatedTimeStamp,
				TestInstructionContainerColor:    nonEditableInformation.TestInstructionContainerColor,
				TCRuleDeletion:                   nonEditableInformation.TCRuleDeletion,
				TCRuleSwap:                       nonEditableInformation.TCRuleSwap,
			},
			EditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableBasicInformationMessage{
				TestInstructionContainerDescription:   editableInformation.TestInstructionContainerDescription,
				TestInstructionContainerMouseOverText: editableInformation.TestInstructionContainerMouseOverText,
			},
			InvisibleBasicInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_InvisibleBasicInformationMessage{
				Enabled: invisibleBasicInformation.Enabled},
			EditableTestInstructionContainerAttributes: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableTestInstructionContainerAttributesMessage{
				TestInstructionContainerExecutionType: editableTestInstructionContainerAttribute.TestInstructionContainerExecutionType},
		}

		immatureTestInstructionContainerInformationMessage := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage{}
		immatureElementModelMessage := fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage{}

		// Create 'immatureTestInstructionContainerMessage' and add 'BasicTestInstructionInformation' and a small part of 'ImmatureSubTestCaseModel'
		immatureTestInstructionContainerMessage = fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage{
			BasicTestInstructionContainerInformation:    &basicTestInstructionContainerInformation,
			ImmatureTestInstructionContainerInformation: &immatureTestInstructionContainerInformationMessage,
			ImmatureSubTestCaseModel:                    &immatureElementModelMessage}

		// Save immatureTestInstructionContainerMessage in map
		ImmatureTestInstructionContainerMessageMap[testInstructionContainerUuid] = immatureTestInstructionContainerMessage

	}

	// **** immatureTestInstructionContainerInformation **** **** immatureTestInstructionContainerInformation **** **** immatureTestInstructionContainerInformation ****
	sqlToExecute = ""
	sqlToExecute = sqlToExecute + "SELECT ITICI.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"ImmatureTestInstructionContainerMessage\" ITICI "
	sqlToExecute = sqlToExecute + "ORDER BY ITICI.\"DomainUuid\" ASC, ITICI.\"TestInstructionContainerUuid\" ASC,  ITICI.\"DropZoneUuid\" ASC, ITICI.\"TestInstructionAttributeUuid\" ASC; "

	// Query DB
	rows, err = fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "aa4b0e8e-3644-491d-be99-8c87ea9b9c23",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Get number of rows for 'immatureTestInstructionContainerInformation'
	immatureTestInstructionContainerInformationSQLCount = rows.CommandTag().RowsAffected()

	// Create map to store ImmatureTestInstructionContainerInformationMessages
	immatureTestInstructionContainerInformationMessagesMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage)

	// Temp variables used when extracting data
	var domainUuid, previousDomainUuid string
	var domainName string
	var testInstructionContainerUuid, previousTestInstructionContainerUuid string
	var testInstructionContainerName string
	var tempTestInstructionAttributeType string
	// First Row in TestData
	var firstRowInSQLRespons bool
	firstRowInSQLRespons = true

	var (
		availableDropZone, previousAvailableDropZone fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage
		availableDropZones                           []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage
	)

	var (
		dropZonePreSetTestInstructionAttribute, previousDropZonePreSetTestInstructionAttribute fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
		dropZonePreSetTestInstructionAttributes                                                []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	)

	var firstImmatureElementUuid string

	var dataStateChange uint8

	// Clear previous variables
	previousDomainUuid = ""
	previousTestInstructionContainerUuid = ""

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		availableDropZone = fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
		dropZonePreSetTestInstructionAttribute = fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}

		err := rows.Scan(

			// temp-data which is not stored in object
			&domainUuid,
			&domainName,
			&testInstructionContainerUuid,
			&testInstructionContainerName,

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
				"Id":           "525079b7-8484-4e61-a811-fa863a41ee2f",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert 'tempTestInstructionAttributeType' into gRPC-type
		dropZonePreSetTestInstructionAttribute.TestInstructionAttributeType = fenixTestCaseBuilderServerGrpcApi.TestInstructionAttributeTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestInstructionAttributeTypeEnum_value[tempTestInstructionAttributeType])

		// Handle the correct order of building together the full object
		dataStateChange = 0

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		dataStateChangeFound :=
			firstRowInSQLRespons == true &&
				domainUuid != previousDomainUuid &&
				testInstructionContainerUuid != previousTestInstructionContainerUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 1
		}

		// All UUIDs are changed and this is not the first row [dataStateChange=2]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid != previousDomainUuid &&
				testInstructionContainerUuid != previousTestInstructionContainerUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 2
		}

		// Only DropZonePreSetTestInstructionAttributeUuid is changed and this is not the first row [dataStateChange=3]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid == previousDomainUuid &&
				testInstructionContainerUuid == previousTestInstructionContainerUuid &&
				availableDropZone.DropZoneUuid == previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 3
		}

		// Only AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=4]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid == previousDomainUuid &&
				testInstructionContainerUuid == previousTestInstructionContainerUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid &&
				dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid != previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 4
		}

		// Only TestInstructionContainerUuid, AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=5]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				domainUuid == previousDomainUuid &&
				testInstructionContainerUuid != previousTestInstructionContainerUuid &&
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
		// Only TestInstructionContainerUuid, AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=5]
		case 2, 5:
			// New DropZone so add the previous DropZone-attributes to the DropZone-array
			newDropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{
				TestInstructionAttributeType: previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeType,
				TestInstructionAttributeUuid: previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid,
				TestInstructionAttributeName: previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeName,
				AttributeValueAsString:       previousDropZonePreSetTestInstructionAttribute.AttributeValueAsString,
				AttributeValueUuid:           previousDropZonePreSetTestInstructionAttribute.AttributeValueUuid,
			}
			dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, &newDropZonePreSetTestInstructionAttribute)

			// Add attributes to previousDropZone
			previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributes

			// Add previousAvailableDropZone to array of DropZone
			availableDropZones = append(availableDropZones, &previousAvailableDropZone)

			// Add the availableDropZones to the ImmatureTestInstructionInformationMessage-map
			immatureTestInstructionContainerInformation.AvailableDropZones = availableDropZones
			immatureTestInstructionContainerInformationMessagesMap[previousTestInstructionContainerUuid] = immatureTestInstructionContainerInformation

			// Create fresh versions of variables
			immatureTestInstructionContainerInformation = fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage{}
			availableDropZones = []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
			dropZonePreSetTestInstructionAttributes = []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}

		// Only DropZonePreSetTestInstructionAttributeUuid is changed and this is not the first row [dataStateChange=3]
		case 3:
			// Add the DropZone attribute to the array for attributes
			newDropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{
				TestInstructionAttributeType: previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeType,
				TestInstructionAttributeUuid: previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid,
				TestInstructionAttributeName: previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeName,
				AttributeValueAsString:       previousDropZonePreSetTestInstructionAttribute.AttributeValueAsString,
				AttributeValueUuid:           previousDropZonePreSetTestInstructionAttribute.AttributeValueUuid,
			}
			dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, &newDropZonePreSetTestInstructionAttribute)

		// Only AvailableDropZoneUuid and DropZonePreSetTestInstructionAttributeUuid are changed and this is not the first row [dataStateChange=4]
		case 4:
			// New DropZone so add the previous DropZone-attributes to the DropZone-array
			newDropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{
				TestInstructionAttributeType: previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeType,
				TestInstructionAttributeUuid: previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid,
				TestInstructionAttributeName: previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeName,
				AttributeValueAsString:       previousDropZonePreSetTestInstructionAttribute.AttributeValueAsString,
				AttributeValueUuid:           previousDropZonePreSetTestInstructionAttribute.AttributeValueUuid,
			}
			dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, &newDropZonePreSetTestInstructionAttribute)

			// Add attributes to previousDropZone
			previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributes

			// Add previousAvailableDropZone to array of DropZone
			availableDropZones = append(availableDropZones, &previousAvailableDropZone)

			// Create fresh versions of variables
			dropZonePreSetTestInstructionAttributes = []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{}

			// Something is wrong in the ordering of the testdata or the testdata itself
		default:
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                                     "0779886a-8280-42b6-9434-46ec1afd1d7f",
				"domainUuid":                             domainUuid,
				"previousDomainUuid":                     previousDomainUuid,
				"testInstructionContainerUuid":           testInstructionContainerUuid,
				"previousTestInstructionContainerUuid":   previousTestInstructionContainerUuid,
				"availableDropZone.DropZoneUuid":         availableDropZone.DropZoneUuid,
				"previousAvailableDropZone.DropZoneUuid": previousAvailableDropZone.DropZoneUuid,
				"dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid":         dropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid,
				"previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid": previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid,
			}).Fatal("Something is wrong in the ordering of the testdata or the testdata itself  --> Should bot happen")

		}

		// Move actual values into previous-variables
		previousDomainUuid = domainUuid
		previousTestInstructionContainerUuid = testInstructionContainerUuid
		previousAvailableDropZone = availableDropZone
		previousDropZonePreSetTestInstructionAttribute = dropZonePreSetTestInstructionAttribute

		// Set to not be the first row
		firstRowInSQLRespons = false

	}

	// Handle last row from database
	// Add the previous DropZone-attributes to the DropZone-array
	newDropZonePreSetTestInstructionAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{
		TestInstructionAttributeType: previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeType,
		TestInstructionAttributeUuid: previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeUuid,
		TestInstructionAttributeName: previousDropZonePreSetTestInstructionAttribute.TestInstructionAttributeName,
		AttributeValueAsString:       previousDropZonePreSetTestInstructionAttribute.AttributeValueAsString,
		AttributeValueUuid:           previousDropZonePreSetTestInstructionAttribute.AttributeValueUuid,
	}
	dropZonePreSetTestInstructionAttributes = append(dropZonePreSetTestInstructionAttributes, &newDropZonePreSetTestInstructionAttribute)

	// Add attributes to previousDropZone
	previousAvailableDropZone.DropZonePreSetTestInstructionAttributes = dropZonePreSetTestInstructionAttributes

	// Add previousAvailableDropZone to array of DropZone
	availableDropZones = append(availableDropZones, &previousAvailableDropZone)

	// Add the availableDropZones to the ImmatureTestInstructionContainerInformationMessage-map
	immatureTestInstructionContainerInformation.AvailableDropZones = availableDropZones
	immatureTestInstructionContainerInformationMessagesMap[previousTestInstructionContainerUuid] = immatureTestInstructionContainerInformation

	// Add 'basicTestInstructionContainerInformation' to map
	immatureTestInstructionContainerMessage, existsInMap := ImmatureTestInstructionContainerMessageMap[testInstructionContainerUuid]
	// testInstructionContainerUuid shouldn't exist in map. If so then there is a problem
	if existsInMap == false {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":                           "8630d2e6-261b-4dab-a499-71463346c5a3",
			"testInstructionContainerUuid": testInstructionContainerUuid,
		}).Fatal("TestInstructionUuid should exist in map. If not then there is a problem")
	}

	// Immature part to 'immatureTestInstructionContainerMessage'
	immatureTestInstructionContainerMessage.ImmatureTestInstructionContainerInformation = &immatureTestInstructionContainerInformation

	// Add 'firstImmatureElementUuid'
	immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel.FirstImmatureElementUuid = firstImmatureElementUuid

	// Store the information back in the map
	ImmatureTestInstructionContainerMessageMap[testInstructionContainerUuid] = immatureTestInstructionContainerMessage

	// ***************************************************************************************************

	// **** ImmatureElementModelMessage **** **** ImmatureElementModelMessage **** **** ImmatureElementModelMessage ****

	sqlToExecute = ""
	sqlToExecute = sqlToExecute + "SELECT IEM.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"BasicTestInstructionContainerInformation\" BTICI, "
	sqlToExecute = sqlToExecute + "\"" + usedDBSchema + "\".\"ImmatureElementModelMessage\" IEM "
	sqlToExecute = sqlToExecute + "WHERE BTICI.\"TestInstructionContainerUuid\" = IEM.\"ImmatureElementUuid\" "
	sqlToExecute = sqlToExecute + "ORDER BY BTICI.\"DomainUuid\" ASC, BTICI.\"TestInstructionContainerUuid\" ASC; "

	// Query DB
	rows, err = fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "c98209fd-150c-4e4c-bcce-303d66523213",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	// Get number of rows for 'immatureTestInstructionContainerInformation'
	immatureTestInstructionContainerInformationSQLCount = rows.CommandTag().RowsAffected()

	// Create map to store ImmatureTestInstructionInformationMessages
	//immatureTestInstructionContainerInformationMessagesMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage)

	// Temp variables used when extracting data
	var tempImmatureElementModelDomainUuid, previousTempImmatureDomainUuid string
	var tempImmatureElementModelDomainName string
	var tempTestCaseModelElementTypeAsString string
	//var previousOriginalElementUuid string
	//var testInstructionContainerUuid, previousTestInstructionContainerUuid string
	//var testInstructionContainerName string
	//var tempTestInstructionAttributeType string
	// First Row in TestData
	//var firstRowInSQLRespons bool
	firstRowInSQLRespons = true

	var (
	//availableDropZone, previousAvailableDropZone fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
	//availableDropZones                           []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
	)

	var (
	//dropZonePreSetTestInstructionAttribute, previousDropZonePreSetTestInstructionAttribute fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	//dropZonePreSetTestInstructionAttributes                                                []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
	)

	//var immatureElementModelMessage fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
	var immatureElementModelElement fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage
	var previousImmatureElementModelElement fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage
	var immatureElementModelElements []*fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage

	//var firstImmatureElementUuid string

	//var dataStateChange uint8

	// Clear previous variables
	previousDomainUuid = ""
	previousTestInstructionContainerUuid = ""

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		immatureElementModelElement = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage{}

		err = rows.Scan(

			// temp-data which is not stored in object
			&tempImmatureElementModelDomainUuid,
			&tempImmatureElementModelDomainName,

			// ImmatureElementModel

			&immatureElementModelElement.OriginalElementUuid,
			&immatureElementModelElement.OriginalElementName,
			&immatureElementModelElement.PreviousElementUuid,
			&immatureElementModelElement.NextElementUuid,
			&immatureElementModelElement.FirstChildElementUuid,
			&immatureElementModelElement.ParentElementUuid,
			&tempTestCaseModelElementTypeAsString,
			&immatureElementModelElement.CurrentElementModelElement,
		)

		if err != nil {
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":           "d4dcd3d8-ab65-46d2-b4a5-85d92481718d",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert 'tempTestCaseModelElementTypeAsString' into gRPC-type
		immatureElementModelElement.TestCaseModelElementType = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum_value[tempTestCaseModelElementTypeAsString])

		// Handle the correct order of building together the full object
		dataStateChange = 0

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		dataStateChangeFound :=
			firstRowInSQLRespons == true &&
				tempImmatureElementModelDomainUuid != previousTempImmatureDomainUuid &&
				immatureElementModelElement.OriginalElementUuid != previousImmatureElementModelElement.OriginalElementUuid
		if dataStateChangeFound == true {
			dataStateChange = 1
		}

		// All UUIDs are changed and this is not the first row [dataStateChange=2]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				tempImmatureElementModelDomainUuid != previousTempImmatureDomainUuid &&
				immatureElementModelElement.OriginalElementUuid != previousImmatureElementModelElement.OriginalElementUuid
		if dataStateChangeFound == true {
			dataStateChange = 2
		}

		// Only immatureElementModelElement.OriginalElementUuid is changed and this is not the first row [dataStateChange=3]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				tempImmatureElementModelDomainUuid == previousTempImmatureDomainUuid &&
				immatureElementModelElement.OriginalElementUuid != previousImmatureElementModelElement.OriginalElementUuid
		if dataStateChangeFound == true {
			dataStateChange = 3
		}

		// A new Element in the Element model, but it belongs to same 'OriginalElementUuid' as previous Element, and this is not the first row [dataStateChange=4]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				tempImmatureElementModelDomainUuid == previousTempImmatureDomainUuid &&
				immatureElementModelElement.OriginalElementUuid == previousImmatureElementModelElement.OriginalElementUuid
		if dataStateChangeFound == true {
			dataStateChange = 4
		}

		// Act on which 'dataStateChange' that was achieved
		switch dataStateChange {

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		case 1:

			newImmatureElementModelElements := []*fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage{}
			immatureElementModelElements = newImmatureElementModelElements

			// All UUIDs are changed and this is not the first row [dataStateChange=2]
		case 2:

			newImmatureElementModelElements := []*fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage{}
			immatureElementModelElements = newImmatureElementModelElements

			// New ElementModelElement so add the previous ElementModelElement to the ElementModelElements-array
			newElementModelToBeStored := fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage{
				OriginalElementUuid:        previousImmatureElementModelElement.OriginalElementUuid,
				OriginalElementName:        previousImmatureElementModelElement.OriginalElementName,
				MatureElementUuid:          previousImmatureElementModelElement.MatureElementUuid,
				PreviousElementUuid:        previousImmatureElementModelElement.PreviousElementUuid,
				NextElementUuid:            previousImmatureElementModelElement.NextElementUuid,
				FirstChildElementUuid:      previousImmatureElementModelElement.FirstChildElementUuid,
				ParentElementUuid:          previousImmatureElementModelElement.ParentElementUuid,
				TestCaseModelElementType:   previousImmatureElementModelElement.TestCaseModelElementType,
				CurrentElementModelElement: previousImmatureElementModelElement.CurrentElementModelElement,
			}
			immatureElementModelElements = append(immatureElementModelElements, &newElementModelToBeStored)

			// Add immatureElementModelElements to 'immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel' which can be found in map
			immatureTestInstructionContainerMessage, existsInMap = ImmatureTestInstructionContainerMessageMap[previousImmatureElementModelElement.OriginalElementUuid]
			// testInstructionContainerUuid should exist in map. If not so then there is a problem
			if existsInMap == false {
				fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
					"Id": "799ca2ef-6ded-4691-ae17-7c77a6a6f37e",
					"previousImmatureElementModelElement.OriginalElementUuid": previousImmatureElementModelElement.OriginalElementUuid,
				}).Fatal("TestInstructionContainerUuid should exist in map. If not then there is a problem")
			}

			//immatureElementModelMessage.TestCaseModelElements = immatureElementModelElements
			immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel.TestCaseModelElements = immatureElementModelElements
			ImmatureTestInstructionContainerMessageMap[previousImmatureElementModelElement.OriginalElementUuid] = immatureTestInstructionContainerMessage

			// Create fresh versions of variables
			//previousImmatureElementModelElement = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage
			//var immatureElementModelMessage = fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage{}

		// Only immatureElementModelElement.OriginalElementUuid is changed and this is not the first row [dataStateChange=3]
		case 3:
			// New ElementModelElement so add the previous ElementModelElement to the ElementModelElements-array
			newElementModelToBeStored := fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage{
				OriginalElementUuid:        previousImmatureElementModelElement.OriginalElementUuid,
				OriginalElementName:        previousImmatureElementModelElement.OriginalElementName,
				MatureElementUuid:          previousImmatureElementModelElement.MatureElementUuid,
				PreviousElementUuid:        previousImmatureElementModelElement.PreviousElementUuid,
				NextElementUuid:            previousImmatureElementModelElement.NextElementUuid,
				FirstChildElementUuid:      previousImmatureElementModelElement.FirstChildElementUuid,
				ParentElementUuid:          previousImmatureElementModelElement.ParentElementUuid,
				TestCaseModelElementType:   previousImmatureElementModelElement.TestCaseModelElementType,
				CurrentElementModelElement: previousImmatureElementModelElement.CurrentElementModelElement,
			}
			immatureElementModelElements = append(immatureElementModelElements, &newElementModelToBeStored)

			// A new Element in the Element model, but it belongs to same 'OriginalElementUuid' as previous Element, and this is not the first row [dataStateChange=4]
		case 4:

			// New ElementModelElement so add the previous ElementModelElement to the ElementModelElements-array
			newElementModelToBeStored := fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage{
				OriginalElementUuid:        previousImmatureElementModelElement.OriginalElementUuid,
				OriginalElementName:        previousImmatureElementModelElement.OriginalElementName,
				MatureElementUuid:          previousImmatureElementModelElement.MatureElementUuid,
				PreviousElementUuid:        previousImmatureElementModelElement.PreviousElementUuid,
				NextElementUuid:            previousImmatureElementModelElement.NextElementUuid,
				FirstChildElementUuid:      previousImmatureElementModelElement.FirstChildElementUuid,
				ParentElementUuid:          previousImmatureElementModelElement.ParentElementUuid,
				TestCaseModelElementType:   previousImmatureElementModelElement.TestCaseModelElementType,
				CurrentElementModelElement: previousImmatureElementModelElement.CurrentElementModelElement,
			}
			immatureElementModelElements = append(immatureElementModelElements, &newElementModelToBeStored)

			// Add immatureElementModelElements to 'immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel' which can be found in map
			immatureTestInstructionContainerMessage, existsInMap = ImmatureTestInstructionContainerMessageMap[previousImmatureElementModelElement.OriginalElementUuid]
			// testInstructionContainerUuid should exist in map. If not so then there is a problem
			if existsInMap == false {
				fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
					"Id": "799ca2ef-6ded-4691-ae17-7c77a6a6f37e",
					"previousImmatureElementModelElement.OriginalElementUuid": previousImmatureElementModelElement.OriginalElementUuid,
				}).Fatal("TestInstructionContainerUuid should exist in map. If not then there is a problem")
			}

			immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel.TestCaseModelElements = immatureElementModelElements
			ImmatureTestInstructionContainerMessageMap[previousImmatureElementModelElement.OriginalElementUuid] = immatureTestInstructionContainerMessage

			// Create fresh versions of variables
			//previousImmatureElementModelElement = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage
			//var immatureElementModelMessage = fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage{}

			// Something is wrong in the ordering of the testdata or the testdata itself
		default:
			fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
				"Id":                                   "24be5ad9-09b3-41a2-81e8-b4171dded878",
				"domainUuid":                           domainUuid,
				"previousDomainUuid":                   previousDomainUuid,
				"testInstructionContainerUuid":         testInstructionContainerUuid,
				"previousTestInstructionContainerUuid": previousTestInstructionContainerUuid,
			}).Fatal("Something is wrong in the ordering of the testdata or the testdata itself  --> Should bot happen")

		}

		// Move previous values to current
		previousImmatureElementModelElement = immatureElementModelElement
		previousTempImmatureDomainUuid = tempImmatureElementModelDomainUuid

		// Set to be not the first row
		firstRowInSQLRespons = false

	}
	// Handle last row from database

	// New ElementModelElement so add the previous ElementModelElement to the ElementModelElements-array
	newElementModelToBeStored := fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementMessage{
		OriginalElementUuid:        immatureElementModelElement.OriginalElementUuid,
		OriginalElementName:        immatureElementModelElement.OriginalElementName,
		MatureElementUuid:          immatureElementModelElement.MatureElementUuid,
		PreviousElementUuid:        immatureElementModelElement.PreviousElementUuid,
		NextElementUuid:            immatureElementModelElement.NextElementUuid,
		FirstChildElementUuid:      immatureElementModelElement.FirstChildElementUuid,
		ParentElementUuid:          immatureElementModelElement.ParentElementUuid,
		TestCaseModelElementType:   immatureElementModelElement.TestCaseModelElementType,
		CurrentElementModelElement: immatureElementModelElement.CurrentElementModelElement,
	}
	immatureElementModelElements = append(immatureElementModelElements, &newElementModelToBeStored)

	// Add immatureElementModelElements to 'immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel' which can be found in map
	immatureTestInstructionContainerMessage, existsInMap = ImmatureTestInstructionContainerMessageMap[immatureElementModelElement.OriginalElementUuid]
	// testInstructionContainerUuid shouldn exist in map. If not so then there is a problem
	if existsInMap == false {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "de167b26-f91a-4108-9ad2-4cc72b981d8a",
			"immatureElementModelElement.OriginalElementUuid": immatureElementModelElement.OriginalElementUuid,
		}).Fatal("TestInstructionContainerUuid should exist in map. If not then there is a problem")
	}

	immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel.TestCaseModelElements = immatureElementModelElements
	ImmatureTestInstructionContainerMessageMap[immatureElementModelElement.OriginalElementUuid] = immatureTestInstructionContainerMessage

	// Loop all ImmatureTestInstructionContainerMessage and create gRPC-response
	var allImmatureTestInstructionContainerMessage []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage

	for _, value := range ImmatureTestInstructionContainerMessageMap { // Order not specified
		allImmatureTestInstructionContainerMessage = append(allImmatureTestInstructionContainerMessage, &value)
	}

	cloudDBImmatureTestInstructionContainerItems = allImmatureTestInstructionContainerMessage

	fmt.Println(basicTestInstructionContainerInformationSQLCount)
	fmt.Println(immatureTestInstructionContainerInformationSQLCount)

	// No errors occurred
	return cloudDBImmatureTestInstructionContainerItems, nil

}
