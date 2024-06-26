package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ListAllRepositoryApiUrls
// TestCase GUI use this gRPC-api to Load a full TestCase with all its data
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) ListAllRepositoryApiUrls(
	ctx context.Context,
	userIdentificationMessage *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (
	*fenixTestCaseBuilderServerGrpcApi.ListAllRepositoryApiUrlsResponseMessage, error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "406f85bd-828a-4cb8-81f5-b3badec79453",
	}).Debug("Incoming 'gRPC - ListAllRepositoryApiUrls'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "c117bff2-5966-4bc4-aa6b-24a90ccd1d7e",
	}).Debug("Outgoing 'gRPC - ListAllRepositoryApiUrls'")

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(userIdentificationMessage.UserIdOnComputer, userIdentificationMessage.GetProtoFileVersionUsedByClient())
	if returnMessage != nil {

		// Not correct proto-file version is used
		var responseMessage *fenixTestCaseBuilderServerGrpcApi.ListAllRepositoryApiUrlsResponseMessage
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListAllRepositoryApiUrlsResponseMessage{
			AckNackResponse:   returnMessage,
			RepositoryApiUrls: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Load list with TestCase Hashes from Database
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.ListAllRepositoryApiUrlsResponseMessage
	responseMessage = fenixCloudDBObject.PrepareLoadUsersTemplateRepositoryUrls(&userIdentificationMessage)

	return responseMessage, nil
}
