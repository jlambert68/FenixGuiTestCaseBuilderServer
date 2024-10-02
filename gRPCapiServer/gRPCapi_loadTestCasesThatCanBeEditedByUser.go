package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ListTestCasesThatCanBeEdited
// TestCase GUI use this gRPC-api to List all TestCases that can be edited by the user
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) ListTestCasesThatCanBeEdited(
	ctx context.Context,
	listTestCasesRequestMessage *fenixTestCaseBuilderServerGrpcApi.ListTestCasesRequestMessage) (
	*fenixTestCaseBuilderServerGrpcApi.ListTestCasesThatCanBeEditedResponseMessage,
	error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "2558d47c-f9a5-489f-8ead-9662dfc3cb17",
	}).Debug("Incoming 'gRPC - ListTestCasesThatCanBeEdited'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "3f970753-89d6-4126-85a5-9bc3987e3586",
	}).Debug("Outgoing 'gRPC - ListTestCasesThatCanBeEdited'")

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(listTestCasesRequestMessage.UserIdOnComputer, listTestCasesRequestMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {

		// Not correct proto-file version is used
		var responseMessage *fenixTestCaseBuilderServerGrpcApi.ListTestCasesThatCanBeEditedResponseMessage
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestCasesThatCanBeEditedResponseMessage{
			AckNackResponse:                returnMessage,
			TestCasesThatCanBeEditedByUser: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Load Full TestCase from Database
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.ListTestCasesThatCanBeEditedResponseMessage
	responseMessage = fenixCloudDBObject.PrepareListTestCasesThatCanBeEdited(
		listTestCasesRequestMessage.GetGCPAuthenticatedUser())

	return responseMessage, nil
}
