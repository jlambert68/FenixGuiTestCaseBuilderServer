package main

import (
	fenixGuiTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// AreYouAlive - *********************************************************************
//Anyone can check if Fenix TestCase Builder server is alive with this service
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) AreYouAlive(ctx context.Context, emptyParameter *fenixGuiTestCaseBuilderServerGrpcApi.EmptyParameter) (*fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse, error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "1ff67695-9a8b-4821-811d-0ab8d33c4d8b",
	}).Debug("Incoming 'gRPC - AreYouAlive'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "9c7f0c3d-7e9f-4c91-934e-8d7a22926d84",
	}).Debug("Outgoing 'gRPC - AreYouAlive'")

	return &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{AckNack: true, Comments: "I'am alive, from Client"}, nil
}

// GetTestInstructionsAndTestContainers - *********************************************************************
// The TestCase Builder asks for all TestInstructions and Pre-defined TestInstructionContainer that the user can add to a TestCase
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) GetTestInstructionsAndTestContainers(ctx context.Context, userIdentificationMessage *fenixGuiTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (*fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage, error) {

	// Define the response message
	var responseMessage *fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "a55f9c82-1d74-44a5-8662-058b8bc9e48f",
	}).Debug("Incoming 'gRPC - GetTestInstructionsAndTestContainers'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "27fb45fe-3266-41aa-a6af-958513977e28",
	}).Debug("Outgoing 'gRPC - GetTestInstructionsAndTestContainers'")

	// Create the response
	responseMessage = &fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage{
		DomainUuid:                       "",
		DomainName:                       "",
		SystemUuid:                       "",
		SystemName:                       "",
		TestInstructionMessages:          nil,
		TestInstructionContainerMessages: nil,
		AckNackResponse: &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "",
			ErrorCodes: nil,
		},
	}

	return responseMessage, nil
}

// GetPinnedTestInstructionsAndTestContainers - *********************************************************************
// The TestCase Builder asks for which TestInstructions and Pre-defined TestInstructionContainer that the user has pinned in the GUI
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) GetPinnedTestInstructionsAndTestContainers(ctx context.Context, userIdentificationMessage *fenixGuiTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (*fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage, error) {

	// Define the response message
	var responseMessage *fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "cffc25f0-b0e6-407a-942a-71fc74f831ac",
	}).Debug("Incoming 'gRPC - GetPinnedTestInstructionsAndTestContainers'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "61e2c28d-b091-442a-b7f8-d2502d9547cf",
	}).Debug("Outgoing 'gRPC - GetPinnedTestInstructionsAndTestContainers'")

	// Create the response
	responseMessage = &fenixGuiTestCaseBuilderServerGrpcApi.TestInstructionsAndTestContainersMessage{
		DomainUuid:                       "",
		DomainName:                       "",
		SystemUuid:                       "",
		SystemName:                       "",
		TestInstructionMessages:          nil,
		TestInstructionContainerMessages: nil,
		AckNackResponse: &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "",
			ErrorCodes: nil,
		},
	}

	return responseMessage, nil
}

// SendTestInstructionsAndTestContainers - *********************************************************************
// The TestCase Builder sends all TestInstructions and Pre-defined TestInstructionContainer that the user has pinned in the GUI
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) SendTestInstructionsAndTestContainers(ctx context.Context, pinnedTestInstructionsAndTestContainersMessage *fenixGuiTestCaseBuilderServerGrpcApi.PinnedTestInstructionsAndTestContainersMessage) (*fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse, error) {

	// Define the response message
	var responseMessage *fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "ff642667-2cbd-4f23-91eb-a6f8e76d9177",
	}).Debug("Incoming 'gRPC - SendTestInstructionsAndTestContainers'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "2c24b079-1e0b-46e9-ad1f-d47e8ff0d3b4",
	}).Debug("Outgoing 'gRPC - SendTestInstructionsAndTestContainers'")

	// Create the response
	responseMessage = &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:    false,
		Comments:   "",
		ErrorCodes: nil,
	}

	return responseMessage, nil

}
