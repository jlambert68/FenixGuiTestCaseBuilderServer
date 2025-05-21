package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ListTestCaseMetaData
// TesterGui use this gRPC-api to get all TestCaseMetaData for user to use when building TestCases
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) ListTestCaseMetaData(
	ctx context.Context,
	userIdentificationMessage *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (
	*fenixTestCaseBuilderServerGrpcApi.ListTestCaseMetaDataResponseMessage, error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "741ccbff-4db4-4a80-b147-51a196714c91",
	}).Debug("Incoming 'gRPC - ListTestCaseMetaData'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "1f056731-f612-4d91-95fa-c181e6dc0f4a",
	}).Debug("Outgoing 'gRPC - ListTestCaseMetaData'")

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(userIdentificationMessage.UserIdOnComputer, userIdentificationMessage.GetProtoFileVersionUsedByClient())
	if returnMessage != nil {

		// Not correct proto-file version is used
		var responseMessage *fenixTestCaseBuilderServerGrpcApi.ListTestCaseMetaDataResponseMessage
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestCaseMetaDataResponseMessage{
			AckNackResponse:            returnMessage,
			TestCaseMetaDataForDomains: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Load list with TestCase Hashes from Database
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.ListTestCaseMetaDataResponseMessage
	responseMessage = fenixCloudDBObject.PrepareLoadUsersTestCaseMetaData(userIdentificationMessage.GCPAuthenticatedUser)

	return responseMessage, nil
}
