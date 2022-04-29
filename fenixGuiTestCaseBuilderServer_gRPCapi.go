package main

import (
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// AreYouAlive - *********************************************************************
//Anyone can check if Fenix TestCase Builder server is alive with this service
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) AreYouAlive(ctx context.Context, emptyParameter *fenixTestCaseBuilderServerGrpcApi.EmptyParameter) (*fenixTestCaseBuilderServerGrpcApi.AckNackResponse, error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "1ff67695-9a8b-4821-811d-0ab8d33c4d8b",
	}).Debug("Incoming 'gRPC - AreYouAlive'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "9c7f0c3d-7e9f-4c91-934e-8d7a22926d84",
	}).Debug("Outgoing 'gRPC - AreYouAlive'")

	return &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{AckNack: true, Comments: "I'am alive, from Client"}, nil
}

// GetTestInstructionsAndTestContainers - *********************************************************************
// The TestCase Builder asks for all TestInstructions and Pre-defined TestInstructionContainer that the user can add to a TestCase
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) GetTestInstructionsAndTestContainers(ctx context.Context, userIdentificationMessage *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (*fenixTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage, error) {

	// Define the response message
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "a55f9c82-1d74-44a5-8662-058b8bc9e48f",
	}).Debug("Incoming 'gRPC - GetTestInstructionsAndTestContainers'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "27fb45fe-3266-41aa-a6af-958513977e28",
	}).Debug("Outgoing 'gRPC - GetTestInstructionsAndTestContainers'")

	// Check if Client is using correct proto files version
	returnMessage := fenixGuiTestCaseBuilderServerObject.isClientUsingCorrectTestDataProtoFileVersion("666", userIdentificationMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage{
			TestInstructionMessages:          nil,
			TestInstructionContainerMessages: nil,
			AckNackResponse:                  returnMessage,
		}

		// Exiting
		return responseMessage, nil
	}

	// Current user
	userID := userIdentificationMessage.UserId

	// Define variables to store data from DB in
	var cloudDBTestInstructionItems []*fenixTestCaseBuilderServerGrpcApi.TestInstructionMessage
	var cloudDBTestInstructionContainerItems []*fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerMessage

	// Get users TestInstruction-data from CloudDB
	err := fenixGuiTestCaseBuilderServerObject.loadClientsTestInstructionsFromCloudDB(userID, &cloudDBTestInstructionItems)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage{
			TestInstructionMessages:          nil,
			TestInstructionContainerMessages: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   "Got some Error when retrieving TestInstructions from database",
				ErrorCodes: []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_UNSPECIFIED},
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Get users TestInstructionContainer-data from CloudDB
	err = fenixGuiTestCaseBuilderServerObject.loadClientsTestInstructionContainersFromCloudDB(userID, &cloudDBTestInstructionContainerItems)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage{
			TestInstructionMessages:          nil,
			TestInstructionContainerMessages: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   "Got some Error when retrieving TestInstructionContainers from database",
				ErrorCodes: []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_UNSPECIFIED},
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage{
		TestInstructionMessages:          cloudDBTestInstructionItems,
		TestInstructionContainerMessages: cloudDBTestInstructionContainerItems,
		AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    true,
			Comments:   "",
			ErrorCodes: nil,
		},
	}

	return responseMessage, nil
}

// GetPinnedTestInstructionsAndTestContainers - *********************************************************************
// The TestCase Builder asks for which TestInstructions and Pre-defined TestInstructionContainer that the user has pinned in the GUI
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) GetPinnedTestInstructionsAndTestContainers(ctx context.Context, userIdentificationMessage *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (*fenixTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage, error) {

	// Define the response message
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "cffc25f0-b0e6-407a-942a-71fc74f831ac",
	}).Debug("Incoming 'gRPC - GetPinnedTestInstructionsAndTestContainers'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "61e2c28d-b091-442a-b7f8-d2502d9547cf",
	}).Debug("Outgoing 'gRPC - GetPinnedTestInstructionsAndTestContainers'")

	// Check if Client is using correct proto files version
	returnMessage := fenixGuiTestCaseBuilderServerObject.isClientUsingCorrectTestDataProtoFileVersion("666", userIdentificationMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage{
			TestInstructionMessages:          nil,
			TestInstructionContainerMessages: nil,
			AckNackResponse:                  returnMessage,
		}

		// Exiting
		return responseMessage, nil
	}

	// Current user
	userID := userIdentificationMessage.UserId

	// Define variables to store data from DB in
	var cloudDBTestInstructionItems []*fenixTestCaseBuilderServerGrpcApi.TestInstructionMessage
	var cloudDBTestInstructionContainerItems []*fenixTestCaseBuilderServerGrpcApi.TestInstructionContainerMessage

	// Get users TestInstruction-data from CloudDB
	err := fenixGuiTestCaseBuilderServerObject.loadClientsTestInstructionsFromCloudDB(userID, &cloudDBTestInstructionItems)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage{
			TestInstructionMessages:          nil,
			TestInstructionContainerMessages: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   "Got some Error when retrieving TestInstructions from database",
				ErrorCodes: []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_UNSPECIFIED},
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Get users TestInstructionContainer-data from CloudDB
	err = fenixGuiTestCaseBuilderServerObject.loadClientsTestInstructionContainersFromCloudDB(userID, &cloudDBTestInstructionContainerItems)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage{
			TestInstructionMessages:          nil,
			TestInstructionContainerMessages: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   "Got some Error when retrieving TestInstructionContainers from database",
				ErrorCodes: []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_UNSPECIFIED},
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage{
		TestInstructionMessages:          cloudDBTestInstructionItems,
		TestInstructionContainerMessages: cloudDBTestInstructionContainerItems,
		AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    true,
			Comments:   "",
			ErrorCodes: nil,
		},
	}

	return responseMessage, nil
}
