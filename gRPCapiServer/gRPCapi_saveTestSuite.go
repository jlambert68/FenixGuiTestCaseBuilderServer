package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// SaveFullTestSuite
// TestSuite GUI use this gRPC-api to save a full TestSuite with all its data
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) SaveFullTestSuite(
	ctx context.Context,
	fullTestSuiteMessage *fenixTestCaseBuilderServerGrpcApi.SaveFullTestSuiteMessageRequest) (
	*fenixTestCaseBuilderServerGrpcApi.AckNackResponse, error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "d72d035d-ccc7-40ca-8a9e-14d7c1941907",
	}).Debug("Incoming 'gRPC - SaveFullTestSuite'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "a9188d8b-02df-455d-a5ca-f20d003b87eb",
	}).Debug("Outgoing 'gRPC - SaveFullTestSuite'")

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		fullTestSuiteMessage.UserIdentification.UserIdOnComputer, fullTestSuiteMessage.UserIdentification.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		// Exiting
		return returnMessage, nil
	}

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Save full TestCase to Cloud DB
	returnMessage = fenixCloudDBObject.PrepareSaveFullTestSuite(fullTestSuiteMessage)
	if returnMessage != nil {
		// Something went wrong when saving to database
		// Exiting
		return returnMessage, nil
	}

	return &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
	}, nil
}
