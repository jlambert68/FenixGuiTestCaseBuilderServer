package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/shared_code"
	"github.com/sirupsen/logrus"
)

// PublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers
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
	}).Debug("Incoming 'gRPCWorker- PublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "97688fc1-7010-4820-9d6f-26ffde24504e",
	}).Debug("Outgoing 'gRPCWorker - PublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")

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
			"in 'PublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")
	}

	// Verify recreated Hashes from gRPC-Worker-message
	var errorSliceWorker []error
	errorSliceWorker = shared_code.VerifyTestInstructionAndTestInstructionContainerAndUsersMessageHashesAndDomain(
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.ConnectorsDomain.ConnectorsDomainUUID,
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)
	if errorSliceWorker != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"ID":               "98ad28e7-3bff-405d-8d61-e94886db5f08",
			"errorSliceWorker": errorSliceWorker,
		}).Error("Problem when recreated Hashes from gRPC-Worker-message " +
			"in 'PublishSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers'")

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

	// Save SupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers in CloudDB
	// Save the TestCase
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	err = fenixCloudDBObject.PrepareSaveSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "70484037-5e3e-43f2-9213-9ff52c3ccbea",
			"error": err,
		}).Error("Couldn't save supported TestInstructions, TestInstructionContainers and Allowed Users in CloudDB")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     err.Error(),
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		return returnMessage, nil
	}

	return returnMessage, nil

}
