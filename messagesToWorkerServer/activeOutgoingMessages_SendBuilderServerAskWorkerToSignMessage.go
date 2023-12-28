package messagesToWorkerServer

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"errors"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"time"
)

// SendGetMessageToSignToProveCallerIdentity
// Worker ask BuilderServer for a message to sign and use the signature to prove identity when sending 'SupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage'
func (messagesToWorkerServerObject *MessagesToWorkerServerObjectStruct) SendBuilderServerAskWorkerToSignMessage(
	messageToSignToProveIdentity string,
	workerAddressToDial string) (
	signMessageResponse *fenixExecutionWorkerGrpcApi.SignMessageResponse,
	err error) {

	messagesToWorkerServerObject.Logger.WithFields(logrus.Fields{
		"id": "fc751387-38c9-44b1-8252-28d0b5c098e0",
	}).Debug("Incoming 'SendBuilderServerAskWorkerToSignMessage'")

	defer messagesToWorkerServerObject.Logger.WithFields(logrus.Fields{
		"id": "b87471cc-5bb0-4074-8702-09cc5d754e20",
	}).Debug("Outgoing 'SendBuilderServerAskWorkerToSignMessage'")

	var ctx context.Context

	// Set up connection to BuilderServer, if that is not already done
	if messagesToWorkerServerObject.connectionToWorkerServerInitiated == false {
		err = messagesToWorkerServerObject.SetConnectionToFenixGuiBuilderServer(workerAddressToDial)
		if err != nil {
			return nil, err
		}
	}

	// Do gRPC-call
	//ctx := context.Background()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() {
		messagesToWorkerServerObject.Logger.WithFields(logrus.Fields{
			"ID": "d013b80c-23bf-42af-b51b-d71a2af45a4e",
		}).Debug("Running Defer Cancel function")
		cancel()
	}()

	// Only add access token when run on GCP
	if common_config.ExecutionLocationForBuilderServer == common_config.GCP {

		// Add Access token
		var returnMessageAckNack bool
		var returnMessageString string
		ctx, returnMessageAckNack, returnMessageString = messagesToWorkerServerObject.generateGCPAccessToken(ctx, workerAddressToDial)
		if returnMessageAckNack == false {
			var newError error
			newError = errors.New(returnMessageString)

			return nil, newError
		}

	}

	// Creates a new temporary client only to be used for this call
	var tempFenixWorkerServerGrpcClient fenixExecutionWorkerGrpcApi.FenixBuilderGprcServicesClient
	tempFenixWorkerServerGrpcClient = fenixExecutionWorkerGrpcApi.
		NewFenixBuilderGprcServicesClient(remoteFenixWorkerServerConnection)

	// Create request message
	var signMessageRequest *fenixExecutionWorkerGrpcApi.SignMessageRequest
	signMessageRequest = &fenixExecutionWorkerGrpcApi.SignMessageRequest{
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(getHighestExecutionWorkerProtoFileVersion()),
		MessageToBeSigned:            messageToSignToProveIdentity,
	}

	// Do gRPC-call
	signMessageResponse, err = tempFenixWorkerServerGrpcClient.BuilderServerAskWorkerToSignMessage(
		ctx, signMessageRequest)

	// Shouldn't happen
	if err != nil {
		messagesToWorkerServerObject.Logger.WithFields(logrus.Fields{
			"ID":                           "9a1324d6-496f-4ac4-9d98-c199bffd5b75",
			"error":                        err,
			"messageToSignToProveIdentity": messageToSignToProveIdentity,
		}).Error("Problem to do gRPC-call to WorkerServer for 'SendBuilderServerAskWorkerToSignMessage'")

		// Set that a new connection needs to be done next time
		messagesToWorkerServerObject.connectionToWorkerServerInitiated = false

		return nil, err
	}

	return signMessageResponse, err

}
