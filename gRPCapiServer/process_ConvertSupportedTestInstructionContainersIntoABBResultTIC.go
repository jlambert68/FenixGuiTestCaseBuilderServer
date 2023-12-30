package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TypeAndStructs"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// Convert message from SupportedTestInstructionContainers into ABBResultTI used for sending to TesterGui
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) convertSupportedTestInstructionContainersIntoABBResultTIC(
	supportedTestInstructionContainerInstance *TestInstructionAndTestInstuctionContainerTypes.TestInstructionContainerStruct) (
	immatureTestInstructionContainerMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage,
	err error) {

	// Convert UpdatedTimeStamp into time-variable
	var timeStampLayoutForParser string
	timeStampLayoutForParser, err = common_config.GenerateTimeStampParserLayout(
		string(supportedTestInstructionContainerInstance.TestInstructionContainer.UpdatedTimeStamp))

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":  "4b596e63-4e7d-474b-8f8f-7a5769314b4a",
			"err": err,
			"supportedTestInstructionContainerInstance.TestInstructionContainer.UpdatedTimeStamp": supportedTestInstructionContainerInstance.TestInstructionContainer.UpdatedTimeStamp,
		}).Error("Couldn't generate parser layout from TimeStamp")

		return nil, err
	}

	var tempUpdatedTimeStamp time.Time
	tempUpdatedTimeStamp, err = time.Parse(timeStampLayoutForParser,
		string(supportedTestInstructionContainerInstance.TestInstructionContainer.UpdatedTimeStamp))

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":  "a43dbf2b-a22f-4acf-abfa-c3d73f25bd16",
			"err": err,
			"supportedTestInstructionContainerInstance.TestInstructionContainer.UpdatedTimeStamp": supportedTestInstructionContainerInstance.TestInstructionContainer.UpdatedTimeStamp,
		}).Error("Couldn't parse TimeStamp in Broadcast-message")

		return nil, err
	}

	// Create 'AvailableDropZones'
	var tempAvailableDropZonesMap map[TypeAndStructs.DropZoneUUIDType]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage
	tempAvailableDropZonesMap = make(map[TypeAndStructs.DropZoneUUIDType]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage)
	var existsInTempAvailableDropZonesMap bool
	var testInstructionContainerAttributeTypeEnumExistsInMap bool

	var tempAvailableDropZonesToBeSentToTesterGui []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage

	for _, tempImmatureTestInstructionContainerInformation := range supportedTestInstructionContainerInstance.ImmatureTestInstructionContainer {

		var tempAvailableDropZoneToBeSentToTesterGui *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage

		// Check of the dropzone already exists within the map
		tempAvailableDropZoneToBeSentToTesterGui, existsInTempAvailableDropZonesMap = tempAvailableDropZonesMap[tempImmatureTestInstructionContainerInformation.DropZoneUUID]
		if existsInTempAvailableDropZonesMap == true {
			// Already Exist
		} else {
			// Create the DropZone and add it to the Map
			tempAvailableDropZoneToBeSentToTesterGui = &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{
				DropZoneUuid:                            string(tempImmatureTestInstructionContainerInformation.DropZoneUUID),
				DropZoneName:                            string(tempImmatureTestInstructionContainerInformation.DropZoneName),
				DropZoneDescription:                     tempImmatureTestInstructionContainerInformation.DropZoneDescription,
				DropZoneMouseOver:                       tempImmatureTestInstructionContainerInformation.DropZoneMouseOver,
				DropZoneColor:                           string(tempImmatureTestInstructionContainerInformation.DropZoneColor),
				DropZonePreSetTestInstructionAttributes: nil,
			}

			tempAvailableDropZonesMap[tempImmatureTestInstructionContainerInformation.DropZoneUUID] = tempAvailableDropZoneToBeSentToTesterGui
		}

		// Handle 'TestInstructionContainerAttributeType'
		var tempTestInstructionAttributeType fenixTestCaseBuilderServerGrpcApi.TestInstructionAttributeTypeEnum
		var tempInt32 int32
		tempInt32, testInstructionContainerAttributeTypeEnumExistsInMap = fenixTestCaseBuilderServerGrpcApi.TestInstructionAttributeTypeEnum_value[string(tempImmatureTestInstructionContainerInformation.TestInstructionAttributeType)]

		if testInstructionContainerAttributeTypeEnumExistsInMap == false {
			// Shouldn't happen
			common_config.Logger.WithFields(logrus.Fields{
				"Id": "b785a710-7076-4829-bd03-74bd2f4e1aff",
				"tempImmatureTestInstructionContainerInformation.TestInstructionContainerAttributeType": tempImmatureTestInstructionContainerInformation.TestInstructionAttributeType,
				"tempImmatureTestInstructionContainerInformation.DropZoneUUID":                          tempImmatureTestInstructionContainerInformation.DropZoneUUID,
				"tempImmatureTestInstructionContainerInformation.DropZoneName":                          tempImmatureTestInstructionContainerInformation.DropZoneName,
			}).Error("Couldn't find TestInstructionContainerAttributeTypeEnum_value")

			return nil, err
		}

		tempTestInstructionAttributeType = fenixTestCaseBuilderServerGrpcApi.TestInstructionAttributeTypeEnum(tempInt32)

		// Add the Attribute to slice of attributes
		var tempDropZonePreSetTestInstructionAttribute *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
		tempDropZonePreSetTestInstructionAttribute = &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{
			TestInstructionAttributeType: tempTestInstructionAttributeType,
			TestInstructionAttributeUuid: string(tempImmatureTestInstructionContainerInformation.TestInstructionAttributeUUID),
			TestInstructionAttributeName: string(tempImmatureTestInstructionContainerInformation.TestInstructionAttributeName),
			AttributeValueAsString:       string(tempImmatureTestInstructionContainerInformation.AttributeValueAsString),
			AttributeValueUuid:           string(tempImmatureTestInstructionContainerInformation.AttributeValueUUID),
			AttributeActionCommand:       fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage_AttributeActionCommandEnum(tempImmatureTestInstructionContainerInformation.AttributeActionCommand),
		}

		tempAvailableDropZoneToBeSentToTesterGui.DropZonePreSetTestInstructionAttributes = append(
			tempAvailableDropZoneToBeSentToTesterGui.DropZonePreSetTestInstructionAttributes,
			tempDropZonePreSetTestInstructionAttribute)
	}

	// Convert DropZone-Map into slice to be used in gRPC-message to TesterGui
	for _, tempAvailableDropZone := range tempAvailableDropZonesMap {
		tempAvailableDropZonesToBeSentToTesterGui = append(tempAvailableDropZonesToBeSentToTesterGui, tempAvailableDropZone)
	}

	// Create 'TestCaseModelElements'
	var tempTestCaseModelElements []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage

	for _, tempImmatureElementModel := range supportedTestInstructionContainerInstance.ImmatureElementModel {

		var tempTestCaseModelElement *fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage
		tempTestCaseModelElement = &fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{
			OriginalElementUuid:      "",
			OriginalElementName:      "",
			ImmatureElementUuid:      "",
			PreviousElementUuid:      "",
			NextElementUuid:          "",
			FirstChildElementUuid:    "",
			ParentElementUuid:        "",
			TestCaseModelElementType: 0,
		}
	}

	immatureTestInstructionContainerMessage = &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage{
		BasicTestInstructionContainerInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage{
			NonEditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_NonEditableBasicInformationMessage{
				DomainUuid:                           string(supportedTestInstructionContainerInstance.TestInstructionContainer.DomainUUID),
				DomainName:                           string(supportedTestInstructionContainerInstance.TestInstructionContainer.DomainName),
				TestInstructionContainerOrignalUuid:  string(supportedTestInstructionContainerInstance.TestInstructionContainer.TestInstructionContainerUUID),
				TestInstructionContainerOriginalName: string(supportedTestInstructionContainerInstance.TestInstructionContainer.TestInstructionContainerName),
				TestInstructionContainerTypeUuid:     string(supportedTestInstructionContainerInstance.TestInstructionContainer.TestInstructionContainerTypeUUID),
				TestInstructionContainerTypeName:     string(supportedTestInstructionContainerInstance.TestInstructionContainer.TestInstructionContainerTypeName),
				Deprecated:                           supportedTestInstructionContainerInstance.TestInstructionContainer.Deprecated,
				MajorVersionNumber:                   uint32(supportedTestInstructionContainerInstance.TestInstructionContainer.MajorVersionNumber),
				MinorVersionNumber:                   uint32(supportedTestInstructionContainerInstance.TestInstructionContainer.MinorVersionNumber),
				UpdatedTimeStamp:                     timestamppb.New(tempUpdatedTimeStamp),
				TestInstructionContainerColor:        string(supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.TestInstructionContainerColor),
				TCRuleDeletion:                       string(supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.TCRuleDeletion),
				TCRuleSwap:                           string(supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.TCRuleSwap),
			},
			EditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableBasicInformationMessage{
				TestInstructionContainerDescription:   supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.TestInstructionContainerDescription,
				TestInstructionContainerMouseOverText: supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.TestInstructionContainerMouseOverText,
			},
			InvisibleBasicInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_InvisibleBasicInformationMessage{
				Enabled: supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.Enabled,
			},
		},
		ImmatureTestInstructionContainerInformation: &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage{
			AvailableDropZones: tempAvailableDropZonesToBeSentToTesterGui,
		},
		ImmatureSubTestCaseModel: &fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage{
			FirstImmatureElementUuid: string(supportedTestInstructionContainerInstance.ImmatureElementModel[0].TopImmatureElementUUID),
			TestCaseModelElements:    tempTestCaseModelElements,
		},
	}

	return immatureTestInstructionContainerMessage, err
}

func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) processTestInstructionContainersBasicTestInstructionContainerInformation(immatureTestInstructionContainerMessageMap map[string]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage) (err error) {

	var (
	//	basicTestInstructionContainerInformation            fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage
	//basicTestInstructionContainerInformationSQLCount    int64
	//immatureTestInstructionContainerInformation fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage
	//immatureTestInstructionContainerInformationSQLCount int64
	//immatureSubTestCaseModel                   fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
	//immatureSubTestCaseModelSQLCount           int64
	)

	usedDBSchema := "FenixBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	// **** BasicTestInstructionContainerInformation **** **** BasicTestInstructionContainerInformation **** **** BasicTestInstructionContainerInformation ****
	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT BTI.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"BasicTestInstructionContainerInformation\" BTI "
	sqlToExecute = sqlToExecute + "ORDER BY BTI.\"DomainUuid\" ASC,  BTI.\"TestInstructionContainerTypeUuid\" ASC, BTI.\"TestInstructionContainerUuid\" ASC; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "b944c506-4ded-4f5e-98c4-06f272d16e1a",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Variables to used when extract data from result set
	//var basicTestInstructionContainerInformation fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerMessage
	var tempTimeStamp time.Time
	//var tempTestInstructionContainerExecutionType string

	// Get number of rows for 'basicTestInstructionContainerInformation'
	//basicTestInstructionContainerInformationSQLCount = rows.CommandTag().RowsAffected()
	var (
		nonEditableInformation    fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_NonEditableBasicInformationMessage
		editableInformation       fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableBasicInformationMessage
		invisibleBasicInformation fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_InvisibleBasicInformationMessage
		//editableTestInstructionContainerAttribute fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableTestInstructionContainerAttributesMessage
		//immatureElementModelMessage                        fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
		//immatureTestInstructionContainerInformationMessage fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage
	)

	// Extract data from DB result set
	for rows.Next() {

		// Initiate a new variable to store the data
		nonEditableInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_NonEditableBasicInformationMessage{}
		editableInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableBasicInformationMessage{}
		invisibleBasicInformation = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_InvisibleBasicInformationMessage{}
		//editableTestInstructionContainerAttribute = fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableTestInstructionContainerAttributesMessage{}

		err := rows.Scan(
			// NonEditableInformation
			&nonEditableInformation.DomainUuid,
			&nonEditableInformation.DomainName,
			&nonEditableInformation.TestInstructionContainerOrignalUuid,
			&nonEditableInformation.TestInstructionContainerOriginalName,
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
			//&tempTestInstructionContainerExecutionType,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "7d082f7c-f987-44e7-97b7-c3c1652955c3",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// Convert TimeStamp into proto-format for TimeStamp
		nonEditableInformation.UpdatedTimeStamp = timestamppb.New(tempTimeStamp)

		// Convert 'tempTestInstructionContainerExecutionType' gRPC-type
		//editableTestInstructionContainerAttribute.TestInstructionContainerExecutionType = fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerExecutionTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerExecutionTypeEnum_value[tempTestInstructionContainerExecutionType])

		// Add 'basicTestInstructionContainerInformation' to map
		testInstructionContainerUuid := nonEditableInformation.TestInstructionContainerOrignalUuid

		_, existsInMap := immatureTestInstructionContainerMessageMap[testInstructionContainerUuid]
		// testInstructionContainerUuid shouldn't exist in map. If so then there is a problem
		if existsInMap == true {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                           "58cd4928-e4b5-4faf-9724-047c1cbc82a1",
				"testInstructionContainerUuid": testInstructionContainerUuid,
				"sqlToExecute":                 sqlToExecute,
			}).Fatal("TestInstructionContainerUuid shouldn't exist in map. If so then there is a problem")

		}

		// Create 'basicTestInstructionContainerInformation' of the parts
		basicTestInstructionContainerInformation := fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage{
			NonEditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_NonEditableBasicInformationMessage{
				DomainUuid:                           nonEditableInformation.DomainUuid,
				DomainName:                           nonEditableInformation.DomainName,
				TestInstructionContainerOrignalUuid:  nonEditableInformation.TestInstructionContainerOrignalUuid,
				TestInstructionContainerOriginalName: nonEditableInformation.TestInstructionContainerOriginalName,
				TestInstructionContainerTypeUuid:     nonEditableInformation.TestInstructionContainerTypeUuid,
				TestInstructionContainerTypeName:     nonEditableInformation.TestInstructionContainerTypeName,
				Deprecated:                           nonEditableInformation.Deprecated,
				MajorVersionNumber:                   nonEditableInformation.MajorVersionNumber,
				MinorVersionNumber:                   nonEditableInformation.MinorVersionNumber,
				UpdatedTimeStamp:                     nonEditableInformation.UpdatedTimeStamp,
				TestInstructionContainerColor:        nonEditableInformation.TestInstructionContainerColor,
				TCRuleDeletion:                       nonEditableInformation.TCRuleDeletion,
				TCRuleSwap:                           nonEditableInformation.TCRuleSwap,
			},
			EditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableBasicInformationMessage{
				TestInstructionContainerDescription:   editableInformation.TestInstructionContainerDescription,
				TestInstructionContainerMouseOverText: editableInformation.TestInstructionContainerMouseOverText,
			},
			InvisibleBasicInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_InvisibleBasicInformationMessage{
				Enabled: invisibleBasicInformation.Enabled},
			//EditableTestInstructionContainerAttributes: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableTestInstructionContainerAttributesMessage{
			//	TestInstructionContainerExecutionType: editableTestInstructionContainerAttribute.TestInstructionContainerExecutionType},
		}

		immatureTestInstructionContainerInformationMessage := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage{}
		immatureElementModelMessage := fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage{}

		// Create 'immatureTestInstructionContainerMessage' and add 'BasicTestInstructionContainerInformation' and a small part of 'ImmatureSubTestCaseModel'
		newImmatureTestInstructionContainerMessage := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage{
			BasicTestInstructionContainerInformation:    &basicTestInstructionContainerInformation,
			ImmatureTestInstructionContainerInformation: &immatureTestInstructionContainerInformationMessage,
			ImmatureSubTestCaseModel:                    &immatureElementModelMessage}

		// Save immatureTestInstructionContainerMessage in map
		immatureTestInstructionContainerMessageMap[testInstructionContainerUuid] = &newImmatureTestInstructionContainerMessage

	}
	return nil
}

// **** immatureTestInstructionContainerInformation **** **** immatureTestInstructionContainerInformation **** **** immatureTestInstructionContainerInformation ****
func (fenixCloudDBObject *FenixCloudDBObjectStruct) processTestInstructionContainersImmatureTestInstructionContainerInformation(immatureTestInstructionContainerMessageMap map[string]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage) (err error) {

	usedDBSchema := "FenixBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT ITII.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"ImmatureTestInstructionContainerInformation\" ITII "
	sqlToExecute = sqlToExecute + "ORDER BY ITII.\"DomainUuid\" ASC, ITII.\"TestInstructionContainerUuid\" ASC,  ITII.\"DropZoneUuid\" ASC, ITII.\"TestInstructionContainerAttributeUuid\" ASC; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "aa4b0e8e-3644-491d-be99-8c87ea9b9c23",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Get number of rows for 'immatureTestInstructionContainerInformation'
	//immatureTestInstructionContainerInformationSQLCount = rows.CommandTag().RowsAffected()

	// Create map to store ImmatureTestInstructionContainerInformationMessages
	//immatureTestInstructionContainerInformationMessagesMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage)

	// Temp variables used when extracting data
	var domainUuid, previousDomainUuid string
	var domainName string
	var testInstructionContainerUuid, previousTestInstructionContainerUuid string
	var testInstructionContainerName string
	var tempTestInstructionContainerAttributeType string

	// First Row in TestData
	var firstRowInSQLRespons bool
	firstRowInSQLRespons = true

	var (
		availableDropZone, previousAvailableDropZone fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage
		availableDropZones                           []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage
	)

	var (
		dropZonePreSetTestInstructionContainerAttribute, previousDropZonePreSetTestInstructionContainerAttribute fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage
		dropZonePreSetTestInstructionContainerAttributes                                                         []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage
	)

	var firstImmatureElementUuid string

	var dataStateChange uint8

	// Clear previous variables
	previousDomainUuid = ""
	previousTestInstructionContainerUuid = ""

	// Initiate a new variable to store the data
	newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
	availableDropZone = newAvailableDropZone

	newDropZonePreSetTestInstructionContainerAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage{}
	dropZonePreSetTestInstructionContainerAttribute = newDropZonePreSetTestInstructionContainerAttribute

	// Extract data from DB result set
	for rows.Next() {

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
			&tempTestInstructionContainerAttributeType,
			&dropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid,
			&dropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeName,
			&dropZonePreSetTestInstructionContainerAttribute.AttributeValueAsString,
			&dropZonePreSetTestInstructionContainerAttribute.AttributeValueUuid,

			// Reference to first element in element-model
			&firstImmatureElementUuid,

			// Attribute Action Command controls have to use the attribute
			&dropZonePreSetTestInstructionContainerAttribute.AttributeActionCommand,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "e514dbca-530d-490e-9fb7-58eaa114a721",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// Convert 'tempTestInstructionContainerAttributeType' into gRPC-type
		dropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeType = fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerAttributeTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerAttributeTypeEnum_value[tempTestInstructionContainerAttributeType])

		// Handle the correct order of building together the full object
		dataStateChange = 0

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		dataStateChangeFound :=
			firstRowInSQLRespons == true &&
				//domainUuid != previousDomainUuid &&
				testInstructionContainerUuid != previousTestInstructionContainerUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid
		//dropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid != previousDropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 1
		}

		// All UUIDs are changed and this is not the first row [dataStateChange=2]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				//domainUuid != previousDomainUuid &&
				testInstructionContainerUuid != previousTestInstructionContainerUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid
		//dropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid != previousDropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 2
		}

		// Only DropZonePreSetTestInstructionContainerAttributeUuid is changed and this is not the first row [dataStateChange=3]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				//domainUuid == previousDomainUuid &&
				testInstructionContainerUuid == previousTestInstructionContainerUuid &&
				availableDropZone.DropZoneUuid == previousAvailableDropZone.DropZoneUuid
		//dropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid != previousDropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 3
		}

		// Only AvailableDropZoneUuid and DropZonePreSetTestInstructionContainerAttributeUuid are changed and this is not the first row [dataStateChange=4]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				//domainUuid == previousDomainUuid &&
				testInstructionContainerUuid == previousTestInstructionContainerUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid
		//dropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid != previousDropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 4
		}

		// Only TestInstructionContainerUuid, AvailableDropZoneUuid and DropZonePreSetTestInstructionContainerAttributeUuid are changed and this is not the first row [dataStateChange=5]
		dataStateChangeFound =
			firstRowInSQLRespons == false &&
				//domainUuid == previousDomainUuid &&
				testInstructionContainerUuid != previousTestInstructionContainerUuid &&
				availableDropZone.DropZoneUuid != previousAvailableDropZone.DropZoneUuid
		//dropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid != previousDropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid
		if dataStateChangeFound == true {
			dataStateChange = 5
		}

		// Act on which 'dataStateChange' that was achieved
		switch dataStateChange {

		// All UUIDs are changed and this is the first row [dataStateChange=1]
		case 1:
			newDropZonePreSetTestInstructionContainerAttributes := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage{}
			dropZonePreSetTestInstructionContainerAttributes = newDropZonePreSetTestInstructionContainerAttributes

			newAvailableDropZones := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
			availableDropZones = newAvailableDropZones

		// All UUIDs are changed and this is not the first row [dataStateChange=2]
		// Only TestInstructionContainerUuid, AvailableDropZoneUuid and DropZonePreSetTestInstructionContainerAttributeUuid are changed and this is not the first row [dataStateChange=5]
		case 2, 5:
			// New DropZone so add the previous DropZone-attributes to the DropZone-array
			dropZonePreSetTestInstructionContainerAttributes = append(dropZonePreSetTestInstructionContainerAttributes, previousDropZonePreSetTestInstructionContainerAttribute)

			// Convert to pointer object instead before storing in map
			var dropZonePreSetTestInstructionContainerAttributesToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage
			for _, tempDropZonePreSetTestInstructionContainerAttributeToStore := range dropZonePreSetTestInstructionContainerAttributes {
				newAdropZonePreSetTestInstructionContainerAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage{}
				newAdropZonePreSetTestInstructionContainerAttribute = tempDropZonePreSetTestInstructionContainerAttributeToStore
				dropZonePreSetTestInstructionContainerAttributesToStore = append(dropZonePreSetTestInstructionContainerAttributesToStore, &newAdropZonePreSetTestInstructionContainerAttribute)
			}

			// Add attributes to previousDropZone
			previousAvailableDropZone.DropZonePreSetTestInstructionContainerAttributes = dropZonePreSetTestInstructionContainerAttributesToStore

			// Add previousAvailableDropZone to array of DropZone
			availableDropZones = append(availableDropZones, previousAvailableDropZone)

			// Add the availableDropZones to the ImmatureTestInstructionContainerInformationMessage-map
			immatureTestInstructionContainerMessage, existsInMap := immatureTestInstructionContainerMessageMap[previousTestInstructionContainerUuid]
			// testInstructionContainerUuid shouldn't exist in map. If so then there is a problem
			if existsInMap == false {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":                           "9fd1b07e-c87a-4583-869b-b3ed28b44616",
					"testInstructionContainerUuid": testInstructionContainerUuid,
					"sqlToExecute":                 sqlToExecute,
				}).Fatal("TestInstructionContainerUuid should exist in map. If not so then there is a problem")
			}

			// Convert to pointer object instead before storing in map
			var availableDropZoneMessageToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage
			for _, tempAvailableDropZones := range availableDropZones {
				newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
				newAvailableDropZone = tempAvailableDropZones
				availableDropZoneMessageToStore = append(availableDropZoneMessageToStore, &newAvailableDropZone)
			}

			immatureTestInstructionContainerMessage.ImmatureTestInstructionContainerInformation.AvailableDropZones = availableDropZoneMessageToStore
			immatureTestInstructionContainerMessageMap[previousTestInstructionContainerUuid] = immatureTestInstructionContainerMessage

			// Create fresh versions of variables
			//newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
			//availableDropZone = newAvailableDropZone

			newAailableDropZones := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
			availableDropZones = newAailableDropZones

			//newDropZonePreSetTestInstructionContainerAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage{}
			//dropZonePreSetTestInstructionContainerAttribute = newDropZonePreSetTestInstructionContainerAttribute

			newDropZonePreSetTestInstructionContainerAttributes := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage{}
			dropZonePreSetTestInstructionContainerAttributes = newDropZonePreSetTestInstructionContainerAttributes

		// Only DropZonePreSetTestInstructionContainerAttributeUuid is changed and this is not the first row [dataStateChange=3]
		case 3:
			// Add the DropZone attribute to the array for attributes
			dropZonePreSetTestInstructionContainerAttributes = append(dropZonePreSetTestInstructionContainerAttributes, previousDropZonePreSetTestInstructionContainerAttribute)

		// Only AvailableDropZoneUuid and DropZonePreSetTestInstructionContainerAttributeUuid are changed and this is not the first row [dataStateChange=4]
		case 4:
			// New DropZone so add the previous DropZone-attributes to the DropZone-array
			dropZonePreSetTestInstructionContainerAttributes = append(dropZonePreSetTestInstructionContainerAttributes, previousDropZonePreSetTestInstructionContainerAttribute)

			// Convert to pointer object instead before storing in map
			var dropZonePreSetTestInstructionContainerAttributesToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage
			for _, tempDropZonePreSetTestInstructionContainerAttributeToStore := range dropZonePreSetTestInstructionContainerAttributes {
				newAdropZonePreSetTestInstructionContainerAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage{}
				newAdropZonePreSetTestInstructionContainerAttribute = tempDropZonePreSetTestInstructionContainerAttributeToStore
				dropZonePreSetTestInstructionContainerAttributesToStore = append(dropZonePreSetTestInstructionContainerAttributesToStore, &newAdropZonePreSetTestInstructionContainerAttribute)
			}

			// Add attributes to previousDropZone
			previousAvailableDropZone.DropZonePreSetTestInstructionContainerAttributes = dropZonePreSetTestInstructionContainerAttributesToStore

			// Add previousAvailableDropZone to array of DropZone
			availableDropZones = append(availableDropZones, previousAvailableDropZone)

			// Clear DropZone-attributes-array
			newDropZonePreSetTestInstructionContainerAttributes := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage{}
			dropZonePreSetTestInstructionContainerAttributes = newDropZonePreSetTestInstructionContainerAttributes

			// Something is wrong in the ordering of the testdata or the testdata itself
		default:
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                                     "352075ba-32f8-4374-86d9-a936ec91b179",
				"domainUuid":                             domainUuid,
				"previousDomainUuid":                     previousDomainUuid,
				"testInstructionContainerUuid":           testInstructionContainerUuid,
				"previousTestInstructionContainerUuid":   previousTestInstructionContainerUuid,
				"availableDropZone.DropZoneUuid":         availableDropZone.DropZoneUuid,
				"previousAvailableDropZone.DropZoneUuid": previousAvailableDropZone.DropZoneUuid,
				"dropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid":         dropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid,
				"previousDropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid": previousDropZonePreSetTestInstructionContainerAttribute.TestInstructionContainerAttributeUuid,
			}).Fatal("Something is wrong in the ordering of the testdata or the testdata itself  --> Should bot happen")

		}

		// Move actual values into previous-variables
		previousDomainUuid = domainUuid
		previousTestInstructionContainerUuid = testInstructionContainerUuid
		previousAvailableDropZone = availableDropZone
		previousDropZonePreSetTestInstructionContainerAttribute = dropZonePreSetTestInstructionContainerAttribute

		// Create fresh versions of variables
		newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
		availableDropZone = newAvailableDropZone

		newDropZonePreSetTestInstructionContainerAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage{}
		dropZonePreSetTestInstructionContainerAttribute = newDropZonePreSetTestInstructionContainerAttribute

		//newDropZonePreSetTestInstructionContainerAttributes := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage{}
		//dropZonePreSetTestInstructionContainerAttributes = newDropZonePreSetTestInstructionContainerAttributes

		// Set to not be the first row
		firstRowInSQLRespons = false

	}

	// Handle last row from database
	// Add the previous DropZone-attributes to the DropZone-array
	dropZonePreSetTestInstructionContainerAttributes = append(dropZonePreSetTestInstructionContainerAttributes, previousDropZonePreSetTestInstructionContainerAttribute)

	// Convert to pointer object instead before storing in map
	var dropZonePreSetTestInstructionContainerAttributesToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage
	for _, tempDropZonePreSetTestInstructionContainerAttributeToStore := range dropZonePreSetTestInstructionContainerAttributes {
		newAdropZonePreSetTestInstructionContainerAttribute := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage{}
		newAdropZonePreSetTestInstructionContainerAttribute = tempDropZonePreSetTestInstructionContainerAttributeToStore
		dropZonePreSetTestInstructionContainerAttributesToStore = append(dropZonePreSetTestInstructionContainerAttributesToStore, &newAdropZonePreSetTestInstructionContainerAttribute)
	}

	// Add attributes to previousDropZone
	previousAvailableDropZone.DropZonePreSetTestInstructionContainerAttributes = dropZonePreSetTestInstructionContainerAttributesToStore

	// Add previousAvailableDropZone to array of DropZone
	availableDropZones = append(availableDropZones, previousAvailableDropZone)

	// Add 'basicTestInstructionContainerInformation' to map
	immatureTestInstructionContainerMessage, existsInMap := immatureTestInstructionContainerMessageMap[testInstructionContainerUuid]
	if existsInMap == false {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":                           "0f59327f-84a9-47bd-bfe2-337c3402ab0c",
			"testInstructionContainerUuid": testInstructionContainerUuid,
		}).Fatal("TestInstructionContainerUuid should exist in map. If not then there is a problem")
	}

	// Convert to pointer object instead before storing in map
	var availableDropZoneMessageToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage
	for _, tempAvailableDropZones := range availableDropZones {
		newAvailableDropZone := fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage{}
		newAvailableDropZone = tempAvailableDropZones
		availableDropZoneMessageToStore = append(availableDropZoneMessageToStore, &newAvailableDropZone)
	}

	// Store the result back in the map
	immatureTestInstructionContainerMessage.ImmatureTestInstructionContainerInformation.AvailableDropZones = availableDropZoneMessageToStore
	immatureTestInstructionContainerMessageMap[testInstructionContainerUuid] = immatureTestInstructionContainerMessage

	return err
}

// **** ImmatureElementModelMessage **** **** ImmatureElementModelMessage **** **** ImmatureElementModelMessage ****
func (fenixCloudDBObject *FenixCloudDBObjectStruct) processTestInstructionContainersImmatureElementModel(immatureTestInstructionContainerMessageMap map[string]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage) (err error) {

	usedDBSchema := "FenixBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT IEM.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"BasicTestInstructionContainerInformation\" BTII, "
	sqlToExecute = sqlToExecute + "\"" + usedDBSchema + "\".\"ImmatureElementModelMessage\" IEM "
	sqlToExecute = sqlToExecute + "WHERE BTII.\"TestInstructionContainerUuid\" = IEM.\"TopImmatureElementUuid\" "
	sqlToExecute = sqlToExecute + "ORDER BY IEM.\"DomainUuid\" ASC, IEM.\"TopImmatureElementUuid\" ASC, IEM.\"IsTopElement\" DESC; " //, IEM.\"CurrentElementModelElement\" ASC; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "c98209fd-150c-4e4c-bcce-303d66523213",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Get number of rows for 'immatureTestInstructionContainerInformation'
	//immatureTestInstructionContainerInformationSQLCount = rows.CommandTag().RowsAffected()

	// Create map to store ImmatureTestInstructionContainerInformationMessages
	//immatureTestInstructionContainerInformationMessagesMap := make(map[string]fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage)

	// Temp variables used when extracting data
	var tempImmatureElementModelDomainUuid string
	var tempImmatureElementModelDomainName string
	var tempTestCaseModelElementTypeAsString string
	var tempIsTopElement bool
	var tempTopElementUuid string
	var previousTempTopElementUuid string

	//var previousOriginalElementUuid string
	//var testInstructionContainerUuid, previousTestInstructionContainerUuid string
	//var testInstructionContainerName string
	//var tempTestInstructionContainerAttributeType string
	// First Row in TestData
	//var firstRowInSQLRespons bool
	firstRowInSQLRespons := true

	var (
	//availableDropZone, previousAvailableDropZone fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage
	//availableDropZones                           []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage
	)

	var (
	//dropZonePreSetTestInstructionContainerAttribute, previousDropZonePreSetTestInstructionContainerAttribute fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage
	//dropZonePreSetTestInstructionContainerAttributes                                                []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionContainerAttributeMessage
	)

	//var immatureElementModelMessage fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage
	var immatureElementModelElement fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage
	var previousImmatureElementModelElement fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage
	var immatureElementModelElements []fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage

	//var firstImmatureElementUuid string

	//var dataStateChange uint8

	// Clear previous variables
	//previousDomainUuid := ""
	//previousTestInstructionContainerUuid := ""

	// Initiate a new variable to store the data
	newImmatureElementModelElement := fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}
	immatureElementModelElement = newImmatureElementModelElement

	previousImmatureElementModelElement = fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}

	// Extract data from DB result set
	for rows.Next() {

		// Initiate new fresh variable
		newImmatureElementModelElement := fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}
		immatureElementModelElement = newImmatureElementModelElement

		err = rows.Scan(

			// temp-data which is not stored in object
			&tempImmatureElementModelDomainUuid,
			&tempImmatureElementModelDomainName,

			// ImmatureElementModel

			&immatureElementModelElement.ImmatureElementUuid,
			&immatureElementModelElement.OriginalElementName,
			&immatureElementModelElement.PreviousElementUuid,
			&immatureElementModelElement.NextElementUuid,
			&immatureElementModelElement.FirstChildElementUuid,
			&immatureElementModelElement.ParentElementUuid,
			&tempTestCaseModelElementTypeAsString,
			&immatureElementModelElement.OriginalElementUuid,
			&tempTopElementUuid,
			&tempIsTopElement,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "7a937579-bb0a-44d4-850f-4cbdd5fff3a5",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// Convert 'tempTestCaseModelElementTypeAsString' into gRPC-type
		immatureElementModelElement.TestCaseModelElementType = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum_value[tempTestCaseModelElementTypeAsString])

		// Handle the correct order of building together the full object
		dataStateChange := 0

		// This is the first row, and it is flagged as Top-element [dataStateChange=1]
		dataStateChangeFound :=
			firstRowInSQLRespons == true &&
				tempIsTopElement == true

		if dataStateChangeFound == true {
			dataStateChange = 1
		}

		// This is not the first row, and it is flagged as Top-element [dataStateChange=2]
		dataStateChangeFound =
			dataStateChange == 0 &&
				firstRowInSQLRespons == false &&
				tempIsTopElement == true

		if dataStateChangeFound == true {
			dataStateChange = 2
		}

		//  This is not the first row, and it is not flagged as Top-element [dataStateChange=3]
		dataStateChangeFound =
			dataStateChange == 0 &&
				firstRowInSQLRespons == false &&
				tempIsTopElement == false

		if dataStateChangeFound == true {
			dataStateChange = 3
		}

		// Act on which 'dataStateChange' that was achieved
		switch dataStateChange {

		// All UUIDs are changed and this is the first row [dataStateChange=1]

		case 1:

			newImmatureElementModelElements := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}
			immatureElementModelElements = newImmatureElementModelElements

			// All UUIDs are changed and this is not the first row [dataStateChange=2]
			// A new Element model Element and this is not the first row [dataStateChange=4]

		case 2, 4:

			immatureElementModelElements = append(immatureElementModelElements, previousImmatureElementModelElement)

			// Add immatureElementModelElements to 'immatureTestInstructionContainerMessage' which can be found in map
			var immatureTestInstructionContainerMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage
			var existsInMap bool
			immatureTestInstructionContainerMessage, existsInMap = immatureTestInstructionContainerMessageMap[previousTempTopElementUuid]
			if existsInMap == false {
				common_config.Logger.WithFields(logrus.Fields{
					"Id": "ef98b5ca-17d5-4bf8-8af4-a1a954736a47",
					"previousImmatureElementModelElement.ImmatureElementUuid": previousImmatureElementModelElement.ImmatureElementUuid,
				}).Fatal("ImmatureElementUuid should exist in map. If not then there is a problem")
			}

			// Convert to pointer object instead before storing in map
			var immatureElementModelElementsToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage
			for _, tempImmatureElementModelElement := range immatureElementModelElements {
				newImmatureElementModelElement := fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}
				newImmatureElementModelElement = tempImmatureElementModelElement
				immatureElementModelElementsToStore = append(immatureElementModelElementsToStore, &newImmatureElementModelElement)
			}

			immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel.TestCaseModelElements = immatureElementModelElementsToStore
			immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel.FirstImmatureElementUuid = previousTempTopElementUuid
			immatureTestInstructionContainerMessageMap[previousTempTopElementUuid] = immatureTestInstructionContainerMessage

			// Create fresh versions of variables
			newIimmatureElementModelElements := []fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}
			immatureElementModelElements = newIimmatureElementModelElements

			// A new Element model Element , but it belongs to same 'ImmatureElementUuid' as previous Element, and this is not the first row [dataStateChange=3]
		case 3:

			immatureElementModelElements = append(immatureElementModelElements, previousImmatureElementModelElement)

			// Something is wrong in the ordering of the testdata or the testdata itself
		default:
			common_config.Logger.WithFields(logrus.Fields{
				"Id":                                  "24be5ad9-09b3-41a2-81e8-b4171dded878",
				"immatureElementModelElement":         immatureElementModelElements,
				"previousImmatureElementModelElement": previousImmatureElementModelElement,
			}).Fatal("Something is wrong in the ordering of the testdata or the testdata itself  --> Should bot happen")

		}

		// Move previous values to current
		previousImmatureElementModelElement = immatureElementModelElement
		previousTempTopElementUuid = tempTopElementUuid

		// Set to be not the first row
		firstRowInSQLRespons = false

	}
	// Handle last row from database

	// New ElementModelElement so add the previous ElementModelElement to the ElementModelElements-array
	immatureElementModelElements = append(immatureElementModelElements, immatureElementModelElement)

	// Add immatureElementModelElements to 'immatureTestInstructionContainerMessage' which can be found in map
	var immatureTestInstructionContainerMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage
	var existsInMap bool
	immatureTestInstructionContainerMessage, existsInMap = immatureTestInstructionContainerMessageMap[tempTopElementUuid]
	if existsInMap == false {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "a1744497-782f-4e82-bec0-ae0205c6573f",
			"immatureElementModelElement.ImmatureElementUuid": immatureElementModelElement.ImmatureElementUuid,
		}).Fatal("ImmatureElementUuid should exist in map. If not then there is a problem")
	}

	// Convert to pointer object instead before storing in map
	var immatureElementModelElementsToStore []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage
	for _, tempImmatureElementModelElement := range immatureElementModelElements {
		newImmatureElementModelElement := fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{}
		newImmatureElementModelElement = tempImmatureElementModelElement
		immatureElementModelElementsToStore = append(immatureElementModelElementsToStore, &newImmatureElementModelElement)
	}

	immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel.TestCaseModelElements = immatureElementModelElementsToStore
	immatureTestInstructionContainerMessage.ImmatureSubTestCaseModel.FirstImmatureElementUuid = tempTopElementUuid
	immatureTestInstructionContainerMessageMap[tempTopElementUuid] = immatureTestInstructionContainerMessage

	return nil

}
