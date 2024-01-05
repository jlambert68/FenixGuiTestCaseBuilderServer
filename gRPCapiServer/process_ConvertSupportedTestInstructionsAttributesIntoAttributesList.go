package gRPCapiServer

import (
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
)

// Convert attributes from SupportedTestInstructions into "raw" attributes list used for sending to TesterGui
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) convertSupportedTestInstructionsAttributesIntoAttributesList(
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages []*TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct) (
	testInstructionAttributesList []*fenixTestCaseBuilderServerGrpcApi.
		ImmatureTestInstructionAttributesMessage_TestInstructionAttributeMessage,
	err error) {

	// Loop Domains
	for _, tempTestInstructionsAndTestInstructionContainersFromGrpcBuilderMessage := range testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages {

		// Loop TestInstructions within a Domain
		for _, tempTestInstruction := range tempTestInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.TestInstructions.TestInstructionsMap {

			// Pick first TestInstructionsInstance, always the latest, and loop attributes
			for _, tempTestInstructionAttribute := range tempTestInstruction.TestInstructionVersions[0].TestInstructionInstance.TestInstructionAttribute {

				// Create gRPC-version of attribute message
				var testInstructionAttribute *fenixTestCaseBuilderServerGrpcApi.
					ImmatureTestInstructionAttributesMessage_TestInstructionAttributeMessage
				testInstructionAttribute = &fenixTestCaseBuilderServerGrpcApi.
					ImmatureTestInstructionAttributesMessage_TestInstructionAttributeMessage{
					DomainUuid:                                    string(tempTestInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID),
					DomainName:                                    string(tempTestInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainName),
					TestInstructionUuid:                           string(tempTestInstruction.TestInstructionVersions[0].TestInstructionInstance.TestInstruction.TestInstructionUUID),
					TestInstructionName:                           string(tempTestInstruction.TestInstructionVersions[0].TestInstructionInstance.TestInstruction.TestInstructionName),
					TestInstructionAttributeUuid:                  string(tempTestInstructionAttribute.TestInstructionAttributeUUID),
					TestInstructionAttributeName:                  string(tempTestInstructionAttribute.TestInstructionAttributeName),
					TestInstructionAttributeTypeUuid:              string(tempTestInstructionAttribute.TestInstructionAttributeTypeUUID),
					TestInstructionAttributeTypeName:              string(tempTestInstructionAttribute.TestInstructionAttributeTypeName),
					TestInstructionAttributeDescription:           tempTestInstructionAttribute.TestInstructionAttributeDescription,
					TestInstructionAttributeMouseOver:             tempTestInstructionAttribute.TestInstructionAttributeMouseOver,
					TestInstructionAttributeVisible:               tempTestInstructionAttribute.TestInstructionAttributeVisible,
					TestInstructionAttributeEnable:                tempTestInstructionAttribute.TestInstructionAttributeEnabled,
					TestInstructionAttributeMandatory:             tempTestInstructionAttribute.TestInstructionAttributeMandatory,
					TestInstructionAttributeVisibleInTestCaseArea: tempTestInstructionAttribute.TestInstructionAttributeVisibleInTestCaseArea,
					TestInstructionAttributeIsDeprecated:          tempTestInstructionAttribute.TestInstructionAttributeIsDeprecated,
					TestInstructionAttributeValueAsString:         string(tempTestInstructionAttribute.TestInstructionAttributeValueAsString),
					TestInstructionAttributeValueUuid:             string(tempTestInstructionAttribute.TestInstructionAttributeValueUUID),
					TestInstructionAttributeUIType:                string(tempTestInstructionAttribute.TestInstructionAttributeType),
				}

				// Add Attribute to list of Attributes
				testInstructionAttributesList = append(testInstructionAttributesList, testInstructionAttribute)

			}
		}
	}

	return testInstructionAttributesList, err

}
