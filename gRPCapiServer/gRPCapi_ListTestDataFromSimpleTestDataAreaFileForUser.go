package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ListAllTestDataForTestDataAreas
// TesterGui use this gRPC-api to get TestData from 'simple' files
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) ListAllTestDataForTestDataAreas(
	ctx context.Context,
	userIdentificationMessage *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (
	*fenixTestCaseBuilderServerGrpcApi.
		ListAllTestDataForTestDataAreasResponseMessage,
	error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "ecd28511-e892-4afb-b6d3-d2366ac76c2f",
	}).Debug("Incoming 'gRPC - ListTestDataFromSimpleTestDataAreaFile'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "f47fac8f-aa97-47a8-b7ce-2ff652b5feb3",
	}).Debug("Outgoing 'gRPC - ListTestDataFromSimpleTestDataAreaFile'")

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userIdentificationMessage.UserIdOnComputer, userIdentificationMessage.GetProtoFileVersionUsedByClient())
	if returnMessage != nil {

		// Not correct proto-file version is used
		var responseMessage *fenixTestCaseBuilderServerGrpcApi.ListAllTestDataForTestDataAreasResponseMessage
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListAllTestDataForTestDataAreasResponseMessage{
			AckNackResponse:                     returnMessage,
			TestDataFromSimpleTestDataAreaFiles: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Load 'simple' TestData
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.ListAllTestDataForTestDataAreasResponseMessage
	responseMessage = fenixCloudDBObject.PrepareLoadUsersTestDataFromSimpleTestDataAreaFile(
		userIdentificationMessage.GCPAuthenticatedUser)

	return responseMessage, nil
}
