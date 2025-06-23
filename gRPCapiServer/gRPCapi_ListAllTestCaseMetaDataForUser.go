package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ListTestCaseAndTestSuiteMetaData
// TesterGui use this gRPC-api to get all TestCaseMetaData and TestSuiteMeta for user to use when building TestCases and TestSuites
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) ListTestCaseAndTestSuiteMetaData(
	ctx context.Context,
	userIdentificationMessage *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (
	*fenixTestCaseBuilderServerGrpcApi.ListTestCaseAndTestSuiteMetaDataResponseMessage, error) {

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
		var responseMessage *fenixTestCaseBuilderServerGrpcApi.ListTestCaseAndTestSuiteMetaDataResponseMessage
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ListTestCaseAndTestSuiteMetaDataResponseMessage{
			AckNackResponse:                        returnMessage,
			TestCaseAndTestSuiteMetaDataForDomains: nil,
		}

		// Exiting
		return responseMessage, nil
	}

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Load list with TestCase Hashes from Database
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.ListTestCaseAndTestSuiteMetaDataResponseMessage
	responseMessage = fenixCloudDBObject.PrepareLoadUsersTestCaseMetaData(userIdentificationMessage.GCPAuthenticatedUser)

	return responseMessage, nil
}
