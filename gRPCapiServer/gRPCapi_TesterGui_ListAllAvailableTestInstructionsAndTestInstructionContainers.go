package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ListAllAvailableTestInstructionsAndTestInstructionContainers - *********************************************************************
// The TestCase Builder asks for all TestInstructions and Pre-defined TestInstructionContainer that the user can add to a TestCase
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) ListAllAvailableTestInstructionsAndTestInstructionContainers(
	ctx context.Context, userIdentificationMessage *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (
	*fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage, error) {

	// Define the response message
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "a55f9c82-1d74-44a5-8662-058b8bc9e48f",
	}).Debug("Incoming 'gRPC - ListAllAvailableTestInstructionsAndTestContainers'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "27fb45fe-3266-41aa-a6af-958513977e28",
	}).Debug("Outgoing 'gRPC - ListAllAvailableTestInstructionsAndTestContainers'")

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion("666", userIdentificationMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
			ImmatureTestInstructions:          nil,
			ImmatureTestInstructionContainers: nil,
			AckNackResponse:                   returnMessage,
		}

		// Exiting
		return responseMessage, nil
	}

	// Current user
	userID := userIdentificationMessage.UserId

	// Define variables to store data from DB in
	var cloudDBImmatureTestInstructionItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage
	var cloudDBImmatureTestInstructionContainerItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Get users ImmatureTestInstruction-data from CloudDB
	cloudDBImmatureTestInstructionItems, err := fenixCloudDBObject.LoadClientsImmatureTestInstructionsFromCloudDB(userID)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
			ImmatureTestInstructions:          nil,
			ImmatureTestInstructionContainers: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving ImmatureTestInstructions from database",
				ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Get users ImmatureTestInstructionContainer-data from CloudDB
	cloudDBImmatureTestInstructionContainerItems, err = fenixCloudDBObject.LoadClientsImmatureTestInstructionContainersFromCloudDB(userID)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
			ImmatureTestInstructions:          nil,
			ImmatureTestInstructionContainers: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving ImmatureTestInstructionContainers from database",
				ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
		ImmatureTestInstructions:          cloudDBImmatureTestInstructionItems,
		ImmatureTestInstructionContainers: cloudDBImmatureTestInstructionContainerItems,
		AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		},
	}

	return responseMessage, nil
}
