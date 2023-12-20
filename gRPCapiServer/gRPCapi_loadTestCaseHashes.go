package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// GetDetailedTestCase
// TestCase GUI use this gRPC-api to Load a full TestCase with all its data
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) GetTestCasesHashes(ctx context.Context, testCasesHashRequest *fenixTestCaseBuilderServerGrpcApi.TestCasesHashRequest) (*fenixTestCaseBuilderServerGrpcApi.TestCasesHashResponse, error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "04aaa37e-1786-47e9-b4da-551624c2ee7d",
	}).Debug("Incoming 'gRPC - GetTestCaseHashes'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "c9f65151-f237-4a58-8495-7ee5cb5bd656",
	}).Debug("Outgoing 'gRPC - GetTestCaseHashes'")

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(testCasesHashRequest.UserIdentification.UserId, testCasesHashRequest.UserIdentification.ProtoFileVersionUsedByClient)
	if returnMessage != nil {

		// Not correct proto-file version is used
		var responseMessage *fenixTestCaseBuilderServerGrpcApi.TestCasesHashResponse
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.TestCasesHashResponse{
			AckNack:         returnMessage,
			TestCasesHashes: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Load list with TestCase Hashes from Database
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.TestCasesHashResponse
	responseMessage = fenixCloudDBObject.PrepareLoadTestCaseHashes(&testCasesHashRequest.TestCaseUuids)

	return responseMessage, nil
}
