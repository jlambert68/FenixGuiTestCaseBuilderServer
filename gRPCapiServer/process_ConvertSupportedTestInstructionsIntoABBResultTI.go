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

// Convert message from SupportedTestInstructions into ABBResultTI used for sending to TesterGui
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) convertSupportedTestInstructionsIntoABBResultTI(
	supportedTestInstructionInstance *TestInstructionAndTestInstuctionContainerTypes.TestInstructionStruct,
	responseVariablesMapStructure *TestInstructionAndTestInstuctionContainerTypes.ResponseVariablesMapStructureStruct) (
	immatureTestInstructionMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage,
	executionDomainThatCanReceiveDirectTargetedTestInstructions *fenixTestCaseBuilderServerGrpcApi.ExecutionDomainsThatCanReceiveDirectTargetedTestInstructionsMessage,
	err error) {

	// Convert UpdatedTimeStamp into time-variable
	var timeStampLayoutForParser string
	timeStampLayoutForParser, err = common_config.GenerateTimeStampParserLayout(
		string(supportedTestInstructionInstance.TestInstruction.UpdatedTimeStamp))

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":  "189fa929-84c6-4563-80ce-06f67caf3923",
			"err": err,
			"supportedTestInstructionInstance.TestInstruction.UpdatedTimeStamp": supportedTestInstructionInstance.TestInstruction.UpdatedTimeStamp,
		}).Error("Couldn't generate parser layout from TimeStamp")

		return nil, nil, err
	}

	var tempUpdatedTimeStamp time.Time
	tempUpdatedTimeStamp, err = time.Parse(timeStampLayoutForParser,
		string(supportedTestInstructionInstance.TestInstruction.UpdatedTimeStamp))

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":  "a43dbf2b-a22f-4acf-abfa-c3d73f25bd16",
			"err": err,
			"supportedTestInstructionInstance.TestInstruction.UpdatedTimeStamp": supportedTestInstructionInstance.TestInstruction.UpdatedTimeStamp,
		}).Error("Couldn't parse TimeStamp in Broadcast-message")

		return nil, nil, err
	}

	// Create 'AvailableDropZones'
	var tempAvailableDropZonesMap map[TypeAndStructs.DropZoneUUIDType]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage
	tempAvailableDropZonesMap = make(map[TypeAndStructs.DropZoneUUIDType]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage)
	var existsInTempAvailableDropZonesMap bool
	var testInstructionAttributeTypeEnumExistsInMap bool

	var tempAvailableDropZonesToBeSentToTesterGui []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage

	for _, tempImmatureTestInstructionInformation := range supportedTestInstructionInstance.ImmatureTestInstructionInformation {

		var tempAvailableDropZoneToBeSentToTesterGui *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage

		// Check of the dropzone already exists within the map
		tempAvailableDropZoneToBeSentToTesterGui, existsInTempAvailableDropZonesMap = tempAvailableDropZonesMap[tempImmatureTestInstructionInformation.DropZoneUUID]
		if existsInTempAvailableDropZonesMap == true {
			// Already Exist
		} else {
			// Create the DropZone and add it to the Map
			tempAvailableDropZoneToBeSentToTesterGui = &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage{
				DropZoneUuid:                            string(tempImmatureTestInstructionInformation.DropZoneUUID),
				DropZoneName:                            string(tempImmatureTestInstructionInformation.DropZoneName),
				DropZoneDescription:                     tempImmatureTestInstructionInformation.DropZoneDescription,
				DropZoneMouseOver:                       tempImmatureTestInstructionInformation.DropZoneMouseOver,
				DropZoneColor:                           string(tempImmatureTestInstructionInformation.DropZoneColor),
				DropZonePreSetTestInstructionAttributes: nil,
			}

			tempAvailableDropZonesMap[tempImmatureTestInstructionInformation.DropZoneUUID] = tempAvailableDropZoneToBeSentToTesterGui
		}

		// Handle 'TestInstructionAttributeType'
		var tempTestInstructionAttributeType fenixTestCaseBuilderServerGrpcApi.TestInstructionAttributeTypeEnum
		var tempInt32 int32
		tempInt32, testInstructionAttributeTypeEnumExistsInMap = fenixTestCaseBuilderServerGrpcApi.TestInstructionAttributeTypeEnum_value[string(tempImmatureTestInstructionInformation.TestInstructionAttributeType)]

		if testInstructionAttributeTypeEnumExistsInMap == false {
			// Shouldn't happen
			common_config.Logger.WithFields(logrus.Fields{
				"Id": "b785a710-7076-4829-bd03-74bd2f4e1aff",
				"tempImmatureTestInstructionInformation.TestInstructionAttributeType": tempImmatureTestInstructionInformation.TestInstructionAttributeType,
				"tempImmatureTestInstructionInformation.DropZoneUUID":                 tempImmatureTestInstructionInformation.DropZoneUUID,
				"tempImmatureTestInstructionInformation.DropZoneName":                 tempImmatureTestInstructionInformation.DropZoneName,
			}).Error("Couldn't find TestInstructionAttributeTypeEnum_value")

			err = errors.New(fmt.Sprintf("Couldn't find TestInstructionAttributeTypeEnum_value for attribute"+
				" in TestInstruction. DropZoneUUID=%s, DropZoneName=%s, TestInstructionAttributeType=%s",
				tempImmatureTestInstructionInformation.DropZoneUUID,
				tempImmatureTestInstructionInformation.DropZoneName,
				tempImmatureTestInstructionInformation.TestInstructionAttributeType))

			return nil, nil, err
		}

		tempTestInstructionAttributeType = fenixTestCaseBuilderServerGrpcApi.TestInstructionAttributeTypeEnum(tempInt32)

		// Add the Attribute to slice of attributes
		var tempDropZonePreSetTestInstructionAttribute *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage
		tempDropZonePreSetTestInstructionAttribute = &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage{
			TestInstructionAttributeType: tempTestInstructionAttributeType,
			TestInstructionAttributeUuid: string(tempImmatureTestInstructionInformation.TestInstructionAttributeUUID),
			TestInstructionAttributeName: string(tempImmatureTestInstructionInformation.TestInstructionAttributeName),
			AttributeValueAsString:       string(tempImmatureTestInstructionInformation.AttributeValueAsString),
			AttributeValueUuid:           string(tempImmatureTestInstructionInformation.AttributeValueUUID),
			AttributeActionCommand: fenixTestCaseBuilderServerGrpcApi.
				ImmatureTestInstructionInformationMessage_AvailableDropZoneMessage_DropZonePreSetTestInstructionAttributeMessage_AttributeActionCommandEnum(
					tempImmatureTestInstructionInformation.AttributeActionCommand),
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

	for _, immatureElementModel := range supportedTestInstructionInstance.ImmatureElementModel {

		// Handle 'TestCaseModelElementType'
		var tempTestCaseModelElementType fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum
		var tempInt32 int32
		tempInt32, existsInTempTestCaseModelElementTypeMap = fenixTestCaseBuilderServerGrpcApi.TestCaseModelElementTypeEnum_value[string(immatureElementModel.TestCaseModelElementType)]

		if existsInTempTestCaseModelElementTypeMap == false {
			// Shouldn't happen
			common_config.Logger.WithFields(logrus.Fields{
				"Id": "9202c5de-3c0e-486b-8805-491d79bbbc89",
				"tempImmatureTestInstructionInformation.TestInstructionAttributeType": immatureElementModel.TestCaseModelElementType,
				"tempImmatureTestInstructionInformation.DropZoneUUID":                 immatureElementModel.ImmatureElementUUID,
				"tempImmatureTestInstructionInformation.DropZoneName":                 immatureElementModel.ImmatureElementName,
			}).Error("Couldn't find TestCaseModelElementTypeEnum_value")

			err = errors.New(fmt.Sprintf("Couldn't find TestCaseModelElementTypeEnum_value for attribute in DropZone"+
				" in TestInstruction. DropZoneUUID=%s, DropZoneName=%s, TestInstructionAttributeType=%s",
				immatureElementModel.ImmatureElementUUID,
				immatureElementModel.ImmatureElementName,
				immatureElementModel.TestCaseModelElementType))

			return nil, nil, err
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

	// Add ResponseVariables
	var responseVariablesMapForGrpc map[string]*fenixTestCaseBuilderServerGrpcApi.ResponseVariableMessage
	responseVariablesMapForGrpc = make(map[string]*fenixTestCaseBuilderServerGrpcApi.ResponseVariableMessage)

	// Loop Response Variable Map and convert
	if responseVariablesMapStructure != nil {
		for responseVariableUuid, tempResponseVariables := range responseVariablesMapStructure.ResponseVariablesMap {

			var responseVariableMessageForGrpc *fenixTestCaseBuilderServerGrpcApi.ResponseVariableMessage
			responseVariableMessageForGrpc = &fenixTestCaseBuilderServerGrpcApi.ResponseVariableMessage{
				ResponseVariableUuid:        string(tempResponseVariables.ResponseVariable.ResponseVariableUuid),
				ResponseVariableName:        string(tempResponseVariables.ResponseVariable.ResponseVariableName),
				ResponseVariableDescription: string(tempResponseVariables.ResponseVariable.ResponseVariableDescription),
				ResponseVariableIsMandatory: bool(tempResponseVariables.ResponseVariable.ResponseVariableIsMandatory),
				ResponseVariableTypeUuid:    string(tempResponseVariables.ResponseVariable.ResponseVariableTypeUuid),
				ResponseVariableTypeName:    string(tempResponseVariables.ResponseVariable.ResponseVariableTypeName),
			}

			// Add to gRPC-map
			responseVariablesMapForGrpc[string(responseVariableUuid)] = responseVariableMessageForGrpc
		}
	}

	// Create the Immature TestInstruction
	immatureTestInstructionMessage = &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage{
		BasicTestInstructionInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage{
			NonEditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_NonEditableBasicInformationMessage{
				DomainUuid:                  string(supportedTestInstructionInstance.TestInstruction.DomainUUID),
				DomainName:                  string(supportedTestInstructionInstance.TestInstruction.DomainName),
				ExecutionDomainUuid:         string(supportedTestInstructionInstance.TestInstruction.ExecutionDomainUUID),
				ExecutionDomainName:         string(supportedTestInstructionInstance.TestInstruction.ExecutionDomainName),
				TestInstructionOriginalUuid: string(supportedTestInstructionInstance.TestInstruction.TestInstructionUUID),
				TestInstructionOriginalName: string(supportedTestInstructionInstance.TestInstruction.TestInstructionName),
				TestInstructionTypeUuid:     string(supportedTestInstructionInstance.TestInstruction.TestInstructionTypeUUID),
				TestInstructionTypeName:     string(supportedTestInstructionInstance.TestInstruction.TestInstructionTypeName),
				Deprecated:                  supportedTestInstructionInstance.TestInstruction.Deprecated,
				MajorVersionNumber:          uint32(supportedTestInstructionInstance.TestInstruction.MajorVersionNumber),
				MinorVersionNumber:          uint32(supportedTestInstructionInstance.TestInstruction.MinorVersionNumber),
				UpdatedTimeStamp:            timestamppb.New(tempUpdatedTimeStamp),
				TestInstructionColor:        string(supportedTestInstructionInstance.BasicTestInstructionInformation.TestInstructionColor),
				TCRuleDeletion:              string(supportedTestInstructionInstance.BasicTestInstructionInformation.TCRuleDeletion),
				TCRuleSwap:                  string(supportedTestInstructionInstance.BasicTestInstructionInformation.TCRuleSwap),
			},
			EditableInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_EditableBasicInformationMessage{
				TestInstructionDescription:   supportedTestInstructionInstance.BasicTestInstructionInformation.TestInstructionDescription,
				TestInstructionMouseOverText: supportedTestInstructionInstance.BasicTestInstructionInformation.TestInstructionMouseOverText,
			},
			InvisibleBasicInformation: &fenixTestCaseBuilderServerGrpcApi.BasicTestInstructionInformationMessage_InvisibleBasicInformationMessage{
				Enabled: supportedTestInstructionInstance.BasicTestInstructionInformation.Enabled,
			},
		},
		ImmatureTestInstructionInformation: &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionInformationMessage{
			AvailableDropZones: tempAvailableDropZonesToBeSentToTesterGui,
		},
		ImmatureSubTestCaseModel: &fenixTestCaseBuilderServerGrpcApi.ImmatureElementModelMessage{
			FirstImmatureElementUuid: string(supportedTestInstructionInstance.ImmatureElementModel[0].TopImmatureElementUUID),
			TestCaseModelElements:    tempTestCaseModelElements,
		},
		ResponseVariablesMapStructure: &fenixTestCaseBuilderServerGrpcApi.ImmatureResponseVariablesMapStructureMessage{
			ResponseVariablesMap: responseVariablesMapForGrpc},
	}

	// Create the GuiName for the ExecutionDomain
	var nameUsedInGui string
	nameUsedInGui = fmt.Sprintf("%s/%s [%s/%s]",
		string(supportedTestInstructionInstance.TestInstruction.DomainName),
		string(supportedTestInstructionInstance.TestInstruction.ExecutionDomainUUID),
		string(supportedTestInstructionInstance.TestInstruction.DomainUUID)[0:6],
		string(supportedTestInstructionInstance.TestInstruction.ExecutionDomainUUID)[0:6])

	// Create the ExecutionDomain information
	executionDomainThatCanReceiveDirectTargetedTestInstructions = &fenixTestCaseBuilderServerGrpcApi.
		ExecutionDomainsThatCanReceiveDirectTargetedTestInstructionsMessage{
		DomainUuid:          string(supportedTestInstructionInstance.TestInstruction.DomainUUID),
		DomainName:          string(supportedTestInstructionInstance.TestInstruction.DomainName),
		ExecutionDomainUuid: string(supportedTestInstructionInstance.TestInstruction.ExecutionDomainUUID),
		ExecutionDomainName: string(supportedTestInstructionInstance.TestInstruction.ExecutionDomainName),
		NameUsedInGui:       nameUsedInGui,
	}

	return immatureTestInstructionMessage, executionDomainThatCanReceiveDirectTargetedTestInstructions, err
}
