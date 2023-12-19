package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"encoding/json"
	"fmt"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/shared_code"
	"github.com/sirupsen/logrus"
	"os"
)

// ConnectorReportCompleteTestInstructionExecutionResult
// When a TestInstruction has been fully executed the Execution Connector use this to inform the results of the execution result to the Worker
func (s *fenixTestCaseBuilderServerGrpcServicesServer) ConnectorPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
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

		var byteSlice []byte
		var byteSliceAsString string
		// Convert TestInstructionVersion to byte-string and then Hash message
		byteSlice, err = json.Marshal(errorSliceWorker)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"ID":               "1f484750-f756-4107-8d5a-7c92b132dc69",
				"errorSliceWorker": errorSliceWorker,
				"err":              err,
			}).Error("Problem when converting into byteSlice")

			// Create return message
			returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Problem when converting into byteSlice",
				ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{},
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
			}

			return returnMessage, nil
		}

		// Convert byteSlice into string
		byteSliceAsString = string(byteSlice)

		// Create return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     byteSliceAsString,
			ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{},
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage, nil
	}

	// Create gRPC-message towards GuiBuilderServer for 'SupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'

	// First
	// Convert back supported TestInstructions, TestInstructionContainers and Allowed Users message from a gRPC-Worker version of the message and check correctness of Hashes
	var testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage *TestInstructionAndTestInstuctionContainerTypes.TestInstructionsAndTestInstructionsContainersStruct
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage, err = shared_code.
		GenerateStandardFromGrpcWorkerMessageForTestInstructionsAndUsers(
			supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage)

	if err != nil {
		// Create return message
		returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     err.Error(),
			ErrorCodes:                   []fenixExecutionWorkerGrpcApi.ErrorCodesEnum{},
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		return returnMessage, nil
	}

	// Second
	// Verify recreated Hashes from gRPC-Builder-message
	var errorSliceBuilder []error
	errorSliceBuilder = shared_code.VerifyTestInstructionAndTestInstructionContainerAndUsersMessageHashes(
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

	// If there are error then loop and concatenate error message to be sent to user
	if errorSliceBuilder != nil {
		var errToReturn string
		for _, errFromBuilder := range errorSliceBuilder {
			if len(errToReturn) == 0 {
				errToReturn = errFromBuilder.Error()
			} else {
				errToReturn = errToReturn + " - " + errFromBuilder.Error()
			}

		}

		// Create return message
		returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     errToReturn,
			ErrorCodes:                   []fenixExecutionWorkerGrpcApi.ErrorCodesEnum{},
			ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
		}

		return returnMessage, nil
	}

	// Third
	// Convert supported TestInstructions, TestInstructionContainers and Allowed Users message into a gRPC-Builder version of the message
	var supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcBuilderMessage *fenixTestCaseBuilderServerGrpcApi.SupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersMessage
	supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcBuilderMessage, err = shared_code.
		GenerateTestInstructionAndTestInstructionContainerAndUserGrpcBuilderMessage(
			supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcWorkerMessage.GetConnectorDomain().GetConnectorsDomainUUID(),
			testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	succeededToSend, responseMessage := fenixGuiBuilderObject.
		SendPublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsersToFenixGuiBuilderServer(
			supportedTestInstructionsAndTestInstructionContainersAndAllowedUsersGrpcBuilderMessage)

	if succeededToSend == false {
		fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
			"id":              "532dff93-5786-4350-96a2-ddf977ee5ec5",
			"responseMessage": responseMessage,
		}).Error("Got some error when sending 'CompleteTestInstructionExecutionResultToFenixExecutionServer'")
	}

	// Create Error Codes
	var errorCodes []fenixExecutionWorkerGrpcApi.ErrorCodesEnum

	// Generate response
	returnMessage = &fenixExecutionWorkerGrpcApi.AckNackResponse{
		AckNack:                      succeededToSend,
		Comments:                     fmt.Sprintf("Messagage from ExecutionServer: '%s'", responseMessage),
		ErrorCodes:                   errorCodes,
		ProtoFileVersionUsedByClient: fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum(common_config.GetHighestExecutionWorkerProtoFileVersion()),
	}

	return returnMessage, nil

}
