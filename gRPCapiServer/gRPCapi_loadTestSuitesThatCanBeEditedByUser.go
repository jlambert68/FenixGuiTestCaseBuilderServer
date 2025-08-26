package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ListTestSuitesThatCanBeEdited
// TestCase GUI use this gRPC-api to List all TestSuites that can be edited by the user
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) ListTestSuitesThatCanBeEdited(
	ctx context.Context,
	listTestSuitesRequestMessage *fenixTestCaseBuilderServerGrpcApi.ListTestSuitesRequestMessage) (
	*fenixTestCaseBuilderServerGrpcApi.ListTestSuitesResponseMessage,
	error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "8994f142-eb2c-4cd5-83c2-6eca86079131",
	}).Debug("Incoming 'gRPC - ListTestSuitesThatCanBeEdited'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "50252fec-2ded-4dff-8dfd-0353a439e852",
	}).Debug("Outgoing 'gRPC - ListTestSuitesThatCanBeEdited'")

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(listTestSuitesRequestMessage.UserIdOnComputer, listTestSuitesRequestMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {

		// Not correct proto-file version is used
		var responseMessage *fenixTestCaseBuilderServerGrpcApi.ListTestSuitesResponseMessage
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestSuitesResponseMessage{

			AckNackResponse:           returnMessage,
			BasicTestSuiteInformation: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// List TestSuites from Database
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.ListTestSuitesResponseMessage
	responseMessage = fenixCloudDBObject.PrepareListTestSuitesThatCanBeEdited(
		listTestSuitesRequestMessage.GetGCPAuthenticatedUser(),
		listTestSuitesRequestMessage.TestSuiteUpdatedMinTimeStamp.AsTime(),
		listTestSuitesRequestMessage.TestSuiteExecutionUpdatedMinTimeStamp.AsTime())

	return responseMessage, nil
}
