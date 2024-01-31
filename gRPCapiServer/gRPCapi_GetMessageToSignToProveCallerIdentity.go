package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	uuidGenerator "github.com/google/uuid"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
)

// GetMessageToSignToProveCallerIdentity
// A Worker calls BuilderServer to receive a message that will be signed when later sending 'SupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage'
func (s *fenixTestCaseBuilderServerGrpcWorkerServicesServerStruct) GetMessageToSignToProveCallerIdentity(
	ctx context.Context,
	emptyParameter *fenixTestCaseBuilderServerGrpcApi.EmptyParameter) (
	getMessageToSignToProveCallerIdentityResponse *fenixTestCaseBuilderServerGrpcApi.GetMessageToSignToProveCallerIdentityResponse,
	err error) {

	var messageToSign string
	messageToSign = uuidGenerator.New().String()

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "a0381a07-7482-4ee6-acc6-8acc72112c66",
	}).Debug("Incoming 'gRPC - GetMessageToSignToProveCallerIdentity'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id":            "d1f8394e-c100-4e39-9a4f-3c633db1a4ea",
		"messageToSign": messageToSign,
	}).Debug("Outgoing 'gRPC - GetMessageToSignToProveCallerIdentity'")

	callingUser := "WorkerServer"

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(callingUser, emptyParameter.GetProtoFileVersionUsedByClient())
	if returnMessage != nil {

		// Not correct proto-file version is used
		getMessageToSignToProveCallerIdentityResponse = &fenixTestCaseBuilderServerGrpcApi.GetMessageToSignToProveCallerIdentityResponse{
			AckNack:       returnMessage,
			MessageToSign: "",
		}

		// Exiting
		return getMessageToSignToProveCallerIdentityResponse, nil
	}

	// Create response message containing message to sign
	getMessageToSignToProveCallerIdentityResponse = &fenixTestCaseBuilderServerGrpcApi.GetMessageToSignToProveCallerIdentityResponse{
		AckNack: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		},
		MessageToSign: messageToSign,
	}

	return getMessageToSignToProveCallerIdentityResponse, nil
}
