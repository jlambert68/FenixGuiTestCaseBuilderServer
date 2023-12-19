package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

/*

  // ************************************************  TestCase Builder ************************************************

  // *** Get data from DB ***

  // The TestCase Builder asks for a list of TestCase, with some basic information. Messages will be streamed in chunks of e.g. 100 TestCases per chunk
  rpc ListAllTestCases(ListTestCasesRequestMessage) returns (stream ListTestCasesResponseMessage) {
  }

  // The TestCase Builder asks for one complete TestCase
  rpc GetDetailedTestCase(GetTestCaseRequestMessage) returns (FullTestCaseMessage) {
  }

  // List all TestInstructions in the TestCase
  rpc ListAllTestCaseTestInstructions(ListAllTestInstructionsForSpecificTestCaseRequestMessage) returns (MatureTestInstructionsMessage) {
  }

  // List all TestInstructionContainers in the TestCase
  rpc ListAllTestCaseTestInstructionContainers(ListAllTestInstructionContainersForSpecificTestCaseRequestMessage) returns (MatureTestInstructionContainerMessage) {
  }

*/

// GetDetailedTestCase
// TestCase GUI use this gRPC-api to Load a full TestCase with all its data
func (s *fenixTestCaseBuilderServerGrpcServicesServer) GetDetailedTestCase(ctx context.Context, getTestCaseRequestMessage *fenixTestCaseBuilderServerGrpcApi.GetTestCaseRequestMessage) (*fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse, error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "76c53d30-a25e-43e3-87ef-523181a0d949",
	}).Debug("Incoming 'gRPC - GetDetailedTestCase'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "d6b330e4-c313-4a7f-be93-044d1ab363e0",
	}).Debug("Outgoing 'gRPC - GetDetailedTestCase'")

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(getTestCaseRequestMessage.UserId, getTestCaseRequestMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {

		// Not correct proto-file version is used
		var responseMessage *fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse{
			AckNackResponse:  returnMessage,
			DetailedTestCase: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Load Full TestCase from Database
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.GetDetailedTestCaseResponse
	responseMessage = fenixCloudDBObject.PrepareLoadFullTestCase(getTestCaseRequestMessage.TestCaseUuid)

	return responseMessage, nil
}
