package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// DeleteTestCaseAtThisDate
// elete-command that updates database that TestCase will be deleted per a certain date
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) DeleteTestCaseAtThisDate(
	ctx context.Context,
	deleteTestCaseAtThisDateRequest *fenixTestCaseBuilderServerGrpcApi.DeleteTestCaseAtThisDateRequest) (
	ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse, err error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "478e5933-6fb9-4f61-82b1-067371f8d8fc",
	}).Debug("Incoming 'gRPC - DeleteTestCaseAtThisDate'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "b2b71f4d-4da6-42d3-b66b-bc548f873c65",
	}).Debug("Outgoing 'gRPC - DeleteTestCaseAtThisDate'")

	// Check if Client is using correct proto files version
	ackNackResponse = common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		deleteTestCaseAtThisDateRequest.UserIdentification.UserIdOnComputer,
		deleteTestCaseAtThisDateRequest.UserIdentification.GetProtoFileVersionUsedByClient())
	if ackNackResponse != nil {

		// Exiting
		return ackNackResponse, nil
	}

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Try to set DeleteDate for TestCase
	err = fenixCloudDBObject.PrepareDeleteTestCaseAtThisDate(deleteTestCaseAtThisDateRequest)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                              "c89fc4ab-1455-40e0-a1b4-bead902343eb",
			"error":                           err,
			"deleteTestCaseAtThisDateRequest": deleteTestCaseAtThisDateRequest,
		}).Error("Couldn't update Delete date for TestCase")

		// Create Return message
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     err.Error(),
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return ackNackResponse, nil
	}

	ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
	}

	return ackNackResponse, nil
}
