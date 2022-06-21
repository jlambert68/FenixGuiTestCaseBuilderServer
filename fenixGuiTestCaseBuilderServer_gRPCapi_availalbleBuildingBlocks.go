package main

import (
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

/*

  // ************************ Data Used to build Available structure, which is what the user can chose from ************************

  // *** Get data from DB ***

  // The TestCase Builder asks for all TestInstructions and Pre-defined TestInstructionContainer that the user can add to a TestCase
  rpc ListAllAvailableTestInstructionsAndTestContainers (UserIdentificationMessage) returns (AvailableTestInstructionsAndPreCreatedTestContainersResponseMessage) {
  }

  // The TestCase Builder asks for which TestInstructions and Pre-defined TestInstructionContainer that the user has pinned in the GUI
  rpc ListAllAvailablePinnedTestInstructionsAndTestContainers (UserIdentificationMessage) returns (AvailablePinnedTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage) {
  }

  // The TestCase Builder asks for all Bonds-elements that can be used in the TestCase-model
  rpc ListAllAvailableBonds (UserIdentificationMessage) returns (ImmatureBondsMessage) {
  }


  // *** Send data to DB ***

  // The TestCase Builder sends all TestInstructions and Pre-defined TestInstructionContainer that the user has pinned in the GUI by the user
  rpc SaveAllPinnedTestInstructionsAndTestContainers (SavePinnedTestInstructionsAndPreCreatedTestInstructionContainersMessage) returns (AckNackResponse) {
  }
*/

// GetTestInstructionsAndTestContainers - *********************************************************************
// The TestCase Builder asks for all TestInstructions and Pre-defined TestInstructionContainer that the user can add to a TestCase
func (s *fenixTestCaseBuilderServerGrpcServicesServer) ListAllAvailableTestInstructionsAndTestContainers(ctx context.Context, userIdentificationMessage *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (*fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestContainersResponseMessage, error) {

	// Define the response message
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestContainersResponseMessage

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "a55f9c82-1d74-44a5-8662-058b8bc9e48f",
	}).Debug("Incoming 'gRPC - ListAllAvailableTestInstructionsAndTestContainers'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "27fb45fe-3266-41aa-a6af-958513977e28",
	}).Debug("Outgoing 'gRPC - ListAllAvailableTestInstructionsAndTestContainers'")

	// Check if Client is using correct proto files version
	returnMessage := fenixGuiTestCaseBuilderServerObject.isClientUsingCorrectTestDataProtoFileVersion("666", userIdentificationMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestContainersResponseMessage{
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

	// Get users ImmatureTestInstruction-data from CloudDB
	err := fenixGuiTestCaseBuilderServerObject.loadClientsImmatureTestInstructionsFromCloudDB(userID, &cloudDBImmatureTestInstructionItems)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestContainersResponseMessage{
			ImmatureTestInstructions:          nil,
			ImmatureTestInstructionContainers: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   "Got some Error when retrieving ImmatureTestInstructions from database",
				ErrorCodes: []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Get users ImmatureTestInstructionContainer-data from CloudDB
	//err = fenixGuiTestCaseBuilderServerObject.loadClientsTestInstructionContainersFromCloudDB(userID, //&cloudDBImmatureTestInstructionContainerItems)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestContainersResponseMessage{
			ImmatureTestInstructions:          nil,
			ImmatureTestInstructionContainers: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   "Got some Error when retrieving ImmatureTestInstructionContainers from database",
				ErrorCodes: []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestContainersResponseMessage{
		ImmatureTestInstructions:          cloudDBImmatureTestInstructionItems,
		ImmatureTestInstructionContainers: cloudDBImmatureTestInstructionContainerItems,
		AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    true,
			Comments:   "",
			ErrorCodes: nil,
		},
	}

	return responseMessage, nil
}