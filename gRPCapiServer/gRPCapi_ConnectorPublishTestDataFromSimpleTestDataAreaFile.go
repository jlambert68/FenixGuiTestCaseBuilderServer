package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// ConnectorPublishTestDataFromSimpleTestDataAreaFile
// AllTemplateRepositoryConnectionParameters Connector Publish TestData From a Simple TestData-file for one TestData-area
func (s *fenixTestCaseBuilderServerGrpcWorkerServicesServerStruct) ConnectorPublishTestDataFromSimpleTestDataAreaFile(
	ctx context.Context,
	testDataFromSimpleTestDataAreaFileMessage *fenixTestCaseBuilderServerGrpcApi.
		TestDataFromSimpleTestDataAreaFileMessage) (
	returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse,
	err error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "10debd4d-3d93-4e85-8e24-59b9f8487234",
		"allTemplateRepositoryConnectionParameters": testDataFromSimpleTestDataAreaFileMessage,
	}).Debug("Incoming 'gRPCWorker- ConnectorPublishTestDataFromSimpleTestDataAreaFile'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "6440f5d3-79c9-4a20-84b8-7a5df5fc7f47",
	}).Debug("Outgoing 'gRPCWorker - ConnectorPublishTestDataFromSimpleTestDataAreaFile'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	returnMessage = common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userId, testDataFromSimpleTestDataAreaFileMessage.GetClientSystemIdentification().ProtoFileVersionUsedByClient)
	if returnMessage != nil {

		fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
			"id":            "cbc600d5-57b8-48b0-83e0-a7c61f70801d",
			"returnMessage": returnMessage,
			"testDataFromSimpleTestDataAreaFileMessage": testDataFromSimpleTestDataAreaFileMessage,
		}).Debug("Not correct proto-file version")

		// Exiting
		return returnMessage, nil
	}

	// Save PublishTestDataFromSimpleTestDataAreaFile
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	err = fenixCloudDBObject.PrepareSavePublishedTestDataFromSimpleTestDataAreaFile(
		testDataFromSimpleTestDataAreaFileMessage.GetClientSystemIdentification().GetDomainUuid(),
		testDataFromSimpleTestDataAreaFileMessage.GetTestDataFromSimpleTestDataAreaFiles(),
		testDataFromSimpleTestDataAreaFileMessage.GetSignedMessageByWorkerServiceAccount())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "c3cdb8aa-fd57-4280-84e6-1af75c5c6da1",
			"error": err,
		}).Error("Couldn't save supported TestData from 'simple' file in CloudDB")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   err.Error(),
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage, nil
	}

	// Generate response when succeed to save to Database
	returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:    true,
		Comments:   "",
		ErrorCodes: []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{},
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
			CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
	}

	return returnMessage, nil

}
