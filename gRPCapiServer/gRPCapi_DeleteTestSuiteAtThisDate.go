package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// DeleteTestSuiteAtThisDate
// Delete-command that updates database that TestSuite will be deleted per a certain date
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) DeleteTestSuiteAtThisDate(
	ctx context.Context,
	deleteTestSuiteAtThisDateRequest *fenixTestCaseBuilderServerGrpcApi.DeleteTestSuiteAtThisDateRequest) (
	ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse, err error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "467146bf-b329-4300-975e-633f33226221",
	}).Debug("Incoming 'gRPC - DeleteTestSuiteAtThisDate'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "ac44672b-efcb-4f34-9bfc-6467dbc8145f",
	}).Debug("Outgoing 'gRPC - DeleteTestSuiteAtThisDate'")

	// Check if Client is using correct proto files version
	ackNackResponse = common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		deleteTestSuiteAtThisDateRequest.UserIdentification.UserIdOnComputer,
		deleteTestSuiteAtThisDateRequest.UserIdentification.GetProtoFileVersionUsedByClient())
	if ackNackResponse != nil {

		// Exiting
		return ackNackResponse, nil
	}

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Try to set DeleteDate for TestCase
	err = fenixCloudDBObject.PrepareDeleteTestSuiteAtThisDate(deleteTestSuiteAtThisDateRequest)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                               "26bdd238-2dcd-4ed5-9240-b77a317a1d7c",
			"error":                            err,
			"deleteTestSuiteAtThisDateRequest": deleteTestSuiteAtThisDateRequest,
		}).Error("Couldn't update Delete date for TestSuite")

		// Create Return message
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   err.Error(),
			ErrorCodes: nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return ackNackResponse, nil
	}

	ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:    true,
		Comments:   "",
		ErrorCodes: nil,
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
			CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
	}

	return ackNackResponse, nil
}
