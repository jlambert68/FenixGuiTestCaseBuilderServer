package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"strings"
)

// ConnectorPublishTemplateRepositoryConnectionParameters
// Connector publish Template Repository Connection Parameters
func (s *fenixTestCaseBuilderServerGrpcWorkerServicesServerStruct) ConnectorPublishTemplateRepositoryConnectionParameters(
	ctx context.Context,
	allTemplateRepositoryConnectionParameters *fenixTestCaseBuilderServerGrpcApi.
		AllTemplateRepositoryConnectionParameters) (
	returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse,
	err error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "56eabbe6-8be7-42c0-86ae-78869e952a90",
		"allTemplateRepositoryConnectionParameters": allTemplateRepositoryConnectionParameters,
	}).Debug("Incoming 'gRPCWorker- ConnectorPublishTemplateRepositoryConnectionParameters'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "c0f0711a-ada9-480e-afda-4d6785049bfa",
	}).Debug("Outgoing 'gRPCWorker - ConnectorPublishTemplateRepositoryConnectionParameters'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	returnMessage = common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userId, allTemplateRepositoryConnectionParameters.GetClientSystemIdentification().ProtoFileVersionUsedByClient)
	if returnMessage != nil {

		fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
			"id":            "ed80eaca-72f7-431c-ad89-8ed565e2fc01",
			"returnMessage": returnMessage,
			"allTemplateRepositoryConnectionParameters": allTemplateRepositoryConnectionParameters,
		}).Debug("Not correct proto-file version")

		// Exiting
		return returnMessage, nil
	}

	// Extract the Hashes that are bases as for the message that was signed
	// ReCreate the  message
	var reCreatedMessageHashThatWasSigned string
	var hashesToHash []string

	// Loop all Template-data and convert into json
	for _, tempTemplate := range allTemplateRepositoryConnectionParameters.GetAllTemplateRepositories() {
		var tempTemplateAsJson string
		tempTemplateAsJson = protojson.Format(tempTemplate)

		// Remove spaces in json
		tempTemplateAsJson = strings.ReplaceAll(tempTemplateAsJson, " ", "")

		// Append to slice to be hashed
		hashesToHash = append(hashesToHash, tempTemplateAsJson)

	}

	// Create a hash of the slice
	reCreatedMessageHashThatWasSigned = fenixSyncShared.HashValues(hashesToHash, true)

	// Save ConnectorPublishTemplateRepositoryConnectionParameters in CloudDB
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	err = fenixCloudDBObject.PrepareSavePublishedTemplateRepositoryConnectionParameters(
		allTemplateRepositoryConnectionParameters.GetClientSystemIdentification().GetDomainUuid(),
		allTemplateRepositoryConnectionParameters.GetAllTemplateRepositories(),
		allTemplateRepositoryConnectionParameters.GetMessageSignatureData(),
		reCreatedMessageHashThatWasSigned)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "70484037-5e3e-43f2-9213-9ff52c3ccbea",
			"error": err,
		}).Error("Couldn't save supported Template Repository Connection Parameters in CloudDB")

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
