package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

// ConnectorPublishSupportedMetaData
// Connector publish supported TestCaseMetaData-parameters
func (s *fenixTestCaseBuilderServerGrpcWorkerServicesServerStruct) ConnectorPublishSupportedMetaData(
	ctx context.Context,
	supportedTestCaseMetaData *fenixTestCaseBuilderServerGrpcApi.
		SupportedTestCaseMetaData) (
	returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse,
	err error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id":                        "76d02500-60b0-41bf-b7c8-37ba314b4f3d",
		"supportedTestCaseMetaData": supportedTestCaseMetaData,
	}).Debug("Incoming 'gRPCWorker- ConnectorPublishSupportedMetaData'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "0d1ccbe3-2122-499b-ab2c-7c97f03b4a91",
	}).Debug("Outgoing 'gRPCWorker - ConnectorPublishSupportedMetaData'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	returnMessage = common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userId, supportedTestCaseMetaData.GetClientSystemIdentification().ProtoFileVersionUsedByClient)
	if returnMessage != nil {

		fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
			"id":                        "ed80eaca-72f7-431c-ad89-8ed565e2fc01",
			"returnMessage":             returnMessage,
			"supportedTestCaseMetaData": supportedTestCaseMetaData,
		}).Debug("Not correct proto-file version")

		// Exiting
		return returnMessage, nil
	}

	// Extract the Hashes that are bases as for the message that was signed
	// ReCreate the  message
	var reCreatedMessageHashThatWasSigned string

	// Create a hash of the slice
	reCreatedMessageHashThatWasSigned = fenixSyncShared.HashSingleValue(supportedTestCaseMetaData.GetSupportedMetaDataAsJson())

	// Save ConnectorPublishTemplateRepositoryConnectionParameters in CloudDB
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	err = fenixCloudDBObject.PrepareSavePublishedSupportedTestCaseMetaDataParameters(
		supportedTestCaseMetaData.GetClientSystemIdentification().GetDomainUuid(),
		supportedTestCaseMetaData.GetSupportedMetaDataAsJson(),
		supportedTestCaseMetaData.GetMessageSignatureData(),
		reCreatedMessageHashThatWasSigned)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "4ce83265-cf63-4547-a566-b3ae9eb10026",
			"error": err,
		}).Error("Couldn't save supported TestCaseMetaData-parameters in CloudDB")

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
