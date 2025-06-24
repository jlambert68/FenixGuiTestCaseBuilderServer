package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// GetDetailedTestSuite
// TestSuite GUI use this gRPC-api to Load a full TestSuite with all its data
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) GetDetailedTestSuite(
	ctx context.Context,
	getTestSuiteRequestMessage *fenixTestCaseBuilderServerGrpcApi.GetTestSuiteRequestMessage) (
	*fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse, error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "6776db6f-1e7e-421b-b78c-6bb8b8946ad1",
	}).Debug("Incoming 'gRPC - GetDetailedTestSuite'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "3b5064b3-8943-46a1-b7f8-af0f21ae2f58",
	}).Debug("Outgoing 'gRPC - GetDetailedTestSuite'")

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		getTestSuiteRequestMessage.UserIdOnComputer, getTestSuiteRequestMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {

		// Not correct proto-file version is used
		var responseMessage *fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
			AckNackResponse:   returnMessage,
			DetailedTestSuite: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Load Full TestSuite from Database
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse
	responseMessage = fenixCloudDBObject.PrepareLoadFullTestSuite(
		getTestSuiteRequestMessage.GetTestSuiteUuid(),
		getTestSuiteRequestMessage.GetGCPAuthenticatedUser())

	return responseMessage, nil
}
