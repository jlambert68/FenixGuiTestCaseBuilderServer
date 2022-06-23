package main

import (
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) processTestInstructionContainersBasicTestInstructionContainerInformation(immatureTestInstructionContainerMessageMap *map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage) (err error) {

	var (
		//	basicTestInstructionContainerInformation            fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage
		//basicTestInstructionContainerInformationSQLCount    int64
		immatureTestInstructionContainerInformation fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage
		//immatureTestInstructionContainerInformationSQLCount int64
		//immatureSubTestCaseModel                   fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
		//immatureSubTestCaseModelSQLCount           int64
	)

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

		return err
	}

	// Variables to used when extract data from result set
	//var basicTestInstructionContainerInformation fenixTestCaseBuilderServerGrpcApi.TestInstructionMessage
	var tempTimeStamp time.Time
	var tempTestInstructionContainerExecutionType string

	// Get number of rows for 'basicTestInstructionContainerInformation'
	//basicTestInstructionContainerInformationSQLCount = rows.CommandTag().RowsAffected()
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

			return err
		}

		// Convert TimeStamp into proto-format for TimeStamp
		nonEditableInformation.UpdatedTimeStamp = timestamppb.New(tempTimeStamp)

		// Convert 'tempTestInstructionContainerExecutionType' gRPC-type
		editableTestInstructionContainerAttribute.TestInstructionContainerExecutionType = fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerExecutionTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerExecutionTypeEnum_value[tempTestInstructionContainerExecutionType])

		// Add 'basicTestInstructionContainerInformation' to map
		testInstructionContainerUuid := nonEditableInformation.TestInstructionContainerUuid
		x := immatureTestInstructionContainerMessageMap[testInstructionContainerUuid]
		immatureTestInstructionContainerMessage, existsInMap := &immatureTestInstructionContainerMessageMap[testInstructionContainerUuid]
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

		return err
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

}

// **** ImmatureElementModelMessage **** **** ImmatureElementModelMessage **** **** ImmatureElementModelMessage ****
func processTestInstructionContainersImmatureElementModel() (err error) {

	usedDBSchema := "FenixGuiBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT IEM.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"BasicTestInstructionContainerInformation\" BTICI, "
	sqlToExecute = sqlToExecute + "\"" + usedDBSchema + "\".\"ImmatureElementModelMessage\" IEM "
	sqlToExecute = sqlToExecute + "WHERE BTICI.\"TestInstructionContainerUuid\" = IEM.\"ImmatureElementUuid\" "
	sqlToExecute = sqlToExecute + "ORDER BY BTICI.\"DomainUuid\" ASC, BTICI.\"TestInstructionContainerUuid\" ASC; "

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "c98209fd-150c-4e4c-bcce-303d66523213",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Get number of rows for 'immatureTestInstructionContainerInformation'
	//immatureTestInstructionContainerInformationSQLCount = rows.CommandTag().RowsAffected()

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
	firstRowInSQLRespons := true

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
	previousDomainUuid := ""
	previousTestInstructionContainerUuid := ""

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

			return err
		}

		// Convert 'tempTestCaseModelElementTypeAsString' into gRPC-type
		immatureElementModelElement.TestCaseModelElementType = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum_value[tempTestCaseModelElementTypeAsString])

		// Handle the correct order of building together the full object
		dataStateChange := 0

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
			immatureTestInstructionContainerMessage, existsInMap := ImmatureTestInstructionContainerMessageMap[previousImmatureElementModelElement.OriginalElementUuid]
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

}
