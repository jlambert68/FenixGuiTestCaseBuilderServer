package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"errors"
	"fmt"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
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
	var existsInTempTestCaseModelElementTypeMap bool

	for _, immatureElementModel := range supportedTestInstructionContainerInstance.ImmatureElementModel {

		// Handle 'TestCaseModelElementType'
		var tempTestCaseModelElementType fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum
		var tempInt32 int32
		tempInt32, existsInTempTestCaseModelElementTypeMap = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum_value[string(immatureElementModel.TestCaseModelElementType)]

		if existsInTempTestCaseModelElementTypeMap == false {
			// Shouldn't happen
			common_config.Logger.WithFields(logrus.Fields{
				"Id": "2012216a-9d78-4f66-84ea-3e8c9bd2723e",
				"tempImmatureTestInstructionInformation.TestInstructionAttributeType": immatureElementModel.TestCaseModelElementType,
				"tempImmatureTestInstructionInformation.DropZoneUUID":                 immatureElementModel.ImmatureElementUUID,
				"tempImmatureTestInstructionInformation.DropZoneName":                 immatureElementModel.ImmatureElementName,
			}).Error("Couldn't find TestCaseModelElementTypeEnum_value")

			err = errors.New(fmt.Sprintf("Couldn't find TestCaseModelElementTypeEnum_value for attribute in DropZone"+
				" in TestInstruction. DropZoneUUID=%s, DropZoneName=%s, TestInstructionAttributeType=%s",
				immatureElementModel.ImmatureElementUUID,
				immatureElementModel.ImmatureElementName,
				immatureElementModel.TestCaseModelElementType))

			return nil, err
		}

		tempTestCaseModelElementType = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum(tempInt32)

		var tempTestCaseModelElement *fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage
		tempTestCaseModelElement = &fenixTestCaseBuilderServerGrpcApi.ImmatureTestCaseModelElementMessage{
			OriginalElementUuid:      string(immatureElementModel.OriginalElementUUID),
			OriginalElementName:      string(immatureElementModel.ImmatureElementName),
			ImmatureElementUuid:      string(immatureElementModel.ImmatureElementUUID),
			PreviousElementUuid:      string(immatureElementModel.PreviousElementUUID),
			NextElementUuid:          string(immatureElementModel.NextElementUUID),
			FirstChildElementUuid:    string(immatureElementModel.FirstChildElementUUID),
			ParentElementUuid:        string(immatureElementModel.ParentElementUUID),
			TestCaseModelElementType: tempTestCaseModelElementType,
		}

		// Add TestCaseElementModelElement to list of elements
		tempTestCaseModelElements = append(tempTestCaseModelElements, tempTestCaseModelElement)

	}

	immatureTestInstructionContainerMessage = &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage{
		BasicTestInstructionContainerInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage{
			NonEditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_NonEditableBasicInformationMessage{
				DomainUuid:                       string(supportedTestInstructionContainerInstance.TestInstructionContainer.DomainUUID),
				DomainName:                       string(supportedTestInstructionContainerInstance.TestInstructionContainer.DomainName),
				TestInstructionContainerUuid:     string(supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.TestInstructionContainerUUID),
				TestInstructionContainerName:     string(supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.TestInstructionContainerName),
				TestInstructionContainerTypeUuid: string(supportedTestInstructionContainerInstance.TestInstructionContainer.TestInstructionContainerTypeUUID),
				TestInstructionContainerTypeName: string(supportedTestInstructionContainerInstance.TestInstructionContainer.TestInstructionContainerTypeName),
				Deprecated:                       supportedTestInstructionContainerInstance.TestInstructionContainer.Deprecated,
				MajorVersionNumber:               uint32(supportedTestInstructionContainerInstance.TestInstructionContainer.MajorVersionNumber),
				MinorVersionNumber:               uint32(supportedTestInstructionContainerInstance.TestInstructionContainer.MinorVersionNumber),
				UpdatedTimeStamp:                 timestamppb.New(tempUpdatedTimeStamp),
				TestInstructionContainerColor:    string(supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.TestInstructionContainerColor),
				TCRuleDeletion:                   string(supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.TCRuleDeletion),
				TCRuleSwap:                       string(supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.TCRuleSwap),
			},
			EditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_EditableBasicInformationMessage{
				TestInstructionContainerDescription:   supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.TestInstructionContainerDescription,
				TestInstructionContainerMouseOverText: supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.TestInstructionContainerMouseOverText,
			},
			InvisibleBasicInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionContainerInformationMessage_InvisibleBasicInformationMessage{
				Enabled: supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.Enabled,
			},
			EditableTestInstructionContainerAttributes: &fenixTestCaseBuilderServerGrpcApi.
				BasicTestInstructionContainerInformationMessage_EditableTestInstructionContainerAttributesMessage{
				TestInstructionContainerExecutionType: fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerExecutionTypeEnum(fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerExecutionTypeEnum_value[string(supportedTestInstructionContainerInstance.BasicTestInstructionContainerInformation.TestInstructionContainerExecutionType)])},
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
