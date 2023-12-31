package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ListAllAvailableBonds - *********************************************************************
// The TestCase Builder asks for all Bonds
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) ListAllAvailableBonds(
	ctx context.Context,
	userIdentificationMessage *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (
	*fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage,
	error) {

	// Define the response message
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "ffb95b37-9ab0-4933-a53c-b7676a12c8f2",
	}).Debug("Incoming 'gRPC - ListAllAvailableBonds'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "b46215a0-6620-428d-bee9-9fa4c5e4e98b",
	}).Debug("Outgoing 'gRPC - ListAllAvailableBonds'")

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(userIdentificationMessage.UserIdOnComputer, userIdentificationMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage{
			ImmatureBonds:   nil,
			AckNackResponse: returnMessage,
		}

		// Exiting
		return responseMessage, nil
	}

	// Define variables to store data from DB in
	var cloudDBAvailableBonds []*fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage_ImmatureBondMessage

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Get users ImmatureTestInstruction-data from CloudDB
	cloudDBAvailableBonds, err := fenixCloudDBObject.LoadAvailableBondsFromCloudDB()
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage{
			ImmatureBonds: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving Available Bonds from database",
				ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage{
		ImmatureBonds: cloudDBAvailableBonds,
		AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		},
	}

	return responseMessage, nil
}
