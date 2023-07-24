package main

import (
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

/*

  // ************************************************  TestCase Builder ************************************************

  // *** Send data to server ***

  // Save full TestCase in DB
  rpc SaveFullTestCase(FullTestCaseMessage) returns (AckNackResponse) {
  }


  // Save a Basic TestCase info in DB
  rpc SaveTestCase(TestCaseBasicInformationMessage) returns (AckNackResponse) {
  }

  // Save all TestInstructions from the TestCase
  rpc SaveAllTestCaseTestInstructions(SaveAllTestInstructionsForSpecificTestCaseRequestMessage) returns (AckNackResponse) {
  }

  // Save all TestInstructionContainers from the TestCase
  rpc SaveAllTestCaseTestInstructionContainers(SaveAllTestInstructionContainersForSpecificTestCaseRequestMessage) returns (AckNackResponse) {
  }

*/

// SaveFullTestCase
// TestCase GUI use this gRPC-api to save a full TestCase with all its data
func (s *fenixTestCaseBuilderServerGrpcServicesServer) SaveFullTestCase(ctx context.Context, fullTestCaseMessage *fenixTestCaseBuilderServerGrpcApi.FullTestCaseMessage) (*fenixTestCaseBuilderServerGrpcApi.AckNackResponse, error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "d5168677-cf4f-4c22-81b5-235f1c34b079",
	}).Debug("Incoming 'gRPC - SaveFullTestCase'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "3670e241-49d1-4931-b729-c95f00199f66",
	}).Debug("Outgoing 'gRPC - SaveFullTestCase'")

	// Check if Client is using correct proto files version
	returnMessage := fenixGuiTestCaseBuilderServerObject.isClientUsingCorrectTestDataProtoFileVersion(fullTestCaseMessage.TestCaseBasicInformation.UserIdentification.UserId, fullTestCaseMessage.TestCaseBasicInformation.UserIdentification.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		// Exiting
		return returnMessage, nil
	}

	// Save full TestCase to Cloud DB
	returnMessage = fenixGuiTestCaseBuilderServerObject.prepareSaveFullTestCase(fullTestCaseMessage)
	if returnMessage != nil {
		// Something went wrong when saving to database
		// Exiting
		return returnMessage, nil
	}

	return &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(fenixGuiTestCaseBuilderServerObject.getHighestFenixTestDataProtoFileVersion()),
	}, nil
}
