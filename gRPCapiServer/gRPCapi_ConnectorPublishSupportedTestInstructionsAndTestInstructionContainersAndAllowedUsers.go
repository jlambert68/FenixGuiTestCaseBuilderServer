package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"fmt"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/shared_code"
	"github.com/sirupsen/logrus"
)

// ConnectorReportCompleteTestInstructionExecutionResult
// When a TestInstruction has been fully executed the Execution Connector use this to inform the results of the execution result to the Worker
func (s *fenixTestCaseBuilderServerGrpcWorkerServicesServerStruct) PublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
	ctx context.Context,
	supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage *fenixTestCaseBuilderServerGrpcApi.
		SupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage) (
	returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse,
	err error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "66ed33eb-a92c-4231-9f37-e04d44d48dfa",
		"supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage": supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage,
	}).Debug("Incoming 'gRPCWorker- ConnectorPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "97688fc1-7010-4820-9d6f-26ffde24504e",
	}).Debug("Outgoing 'gRPCWorker - ConnectorPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")

	// Calling system
	userId := "Execution Connector"

	// Check if Client is using correct proto files version
	returnMessage = common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userId, supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage.GetClientSystemIdentification().ProtoFileVersionUsedByClient)
	if returnMessage != nil {

		fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
			"id":            "95b92ac2-6597-4498-ab57-bb48530e7dfc",
			"returnMessage": returnMessage,
			"supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage": supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage,
		}).Debug("Not correct proto-file version")

		// Exiting
		return returnMessage, nil
	}

	// Convert back supported TestInstructions, TestInstructionContainers and Allowed Users message from a gRPC-Worker version of the message and check correctness of Hashes
	var testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.TestInstructionsAndTestInstructionsContainersStruct
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage, err = shared_code.
		GenerateStandardFromGrpcBuilderMessageForTestInstructionsAndUsers(
			supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"ID":    "a5b86695-a1d4-43b4-9b73-8d5a29871269",
			"error": err,
		}).Fatalln("Problem when Convert back supported TestInstructions, TestInstructionContainers and " +
			"Allowed Users message from a gRPC-Builder version of the message and check correctness of Hashes " +
			"in 'ConnectorPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")
	}

	// Verify recreated Hashes from gRPC-Worker-message
	var errorSliceWorker []error
	errorSliceWorker = shared_code.VerifyTestInstructionAndTestInstructionContainerAndUsersMessageHashes(
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)
	if errorSliceWorker != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"ID":               "98ad28e7-3bff-405d-8d61-e94886db5f08",
			"errorSliceWorker": errorSliceWorker,
		}).Error("Problem when recreated Hashes from gRPC-Worker-message " +
			"in 'ConnectorPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")

		// Loop error messages and concatenate into one string
		var errorMessageBackToConnector string
		for _, errorFromWorker := range errorSliceWorker {
			if len(errorMessageBackToConnector) == 0 {
				errorMessageBackToConnector = errorFromWorker.Error()
			} else {
				errorMessageBackToConnector = errorMessageBackToConnector + "; " + errorFromWorker.Error()
			}
		}

		// Create return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     errorMessageBackToConnector,
			ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{},
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage, nil
	}

	// Save Published TestInstructions, TestInstructionContainers and Allowed Users to CloudDB

	// Generate response when succeed to save to Database
	returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{},
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
	}

	fmt.Println("Hej hej")

	return returnMessage, nil

}
