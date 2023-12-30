package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/jlambert68/FenixTestInstructionsAdminShared/TestInstructionAndTestInstuctionContainerTypes"
	"github.com/sirupsen/logrus"
)

// ListAllAvailableTestInstructionsAndTestInstructionContainers - *********************************************************************
// The TestCase Builder asks for all TestInstructions and Pre-defined TestInstructionContainer that the user can add to a TestCase
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) ListAllAvailableTestInstructionsAndTestInstructionContainers(
	ctx context.Context, userIdentificationMessage *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (
	*fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage, error) {

	// Define the response message
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "a55f9c82-1d74-44a5-8662-058b8bc9e48f",
	}).Debug("Incoming 'gRPC - ListAllAvailableTestInstructionsAndTestContainers'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "27fb45fe-3266-41aa-a6af-958513977e28",
	}).Debug("Outgoing 'gRPC - ListAllAvailableTestInstructionsAndTestContainers'")

	var err error

	// Current user
	var gCPAuthenticatedUser string
	gCPAuthenticatedUser = userIdentificationMessage.UserId

	// Check if Client is using correct proto files version
	var returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
	returnMessage = common_config.IsClientUsingCorrectTestDataProtoFileVersion(gCPAuthenticatedUser, userIdentificationMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
			ImmatureTestInstructions:          nil,
			ImmatureTestInstructionContainers: nil,
			AckNackResponse:                   returnMessage,
		}

		// Exiting
		return responseMessage, nil
	}

	// Define variables to store data from DB in
	var cloudDBImmatureTestInstructionItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage
	var cloudDBImmatureTestInstructionContainerItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Load Domains that User has access to
	var domainUuidList []string
	domainUuidList, err = fenixCloudDBObject.PrepareLoadUsersDomains(gCPAuthenticatedUser)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "a9cdb129-0a67-4bec-adfc-ae8cba283c98",
			"error":                err,
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Got some problem when loading users domains from database")

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.
			AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
			ImmatureTestInstructions:          nil,
			ImmatureTestInstructionContainers: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   err.Error(),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
		}

		return responseMessage, err

	}

	// Load all Supported TestInstructions, TestInstructionContainers belonging to all domain in DomainList
	var testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages []*TestInstructionAndTestInstuctionContainerTypes.
		TestInstructionsAndTestInstructionsContainersStruct
	testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages, err = fenixCloudDBObject.
		PrepareLoadDomainsSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(domainUuidList)

	if err != nil {
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"id":                   "2ea14e8a-f0b6-423b-8f2b-79049fb444a8",
				"error":                err,
				"gCPAuthenticatedUser": gCPAuthenticatedUser,
			}).Error("Got some problem when loading users published TestInstruction andTestInstructionContainers")

			responseMessage = &fenixTestCaseBuilderServerGrpcApi.
				AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
				ImmatureTestInstructions:          nil,
				ImmatureTestInstructionContainers: nil,
				AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
					AckNack:    false,
					Comments:   err.Error(),
					ErrorCodes: nil,
					ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
						CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
							GetHighestFenixGuiBuilderProtoFileVersion()),
				},
			}

			return responseMessage, err

		}
	}

	// 	Loop all received messages from database and extract users ImmatureTestInstruction-data and ImmatureTestInstructionContainer-data
	for _, testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage := range testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages {

		// Loop all TestInstructions from 'this domain'
		for _, tempTestInstruction := range testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.
			TestInstructions.TestInstructionsMap {

			tempTestInstruction.TestInstructionVersions[0].TestInstructionInstance.
		}
	}

	// Get users ImmatureTestInstruction-data from CloudDB
	cloudDBImmatureTestInstructionItems, err := fenixCloudDBObject.LoadClientsImmatureTestInstructionsFromCloudDB(userID)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
			ImmatureTestInstructions:          nil,
			ImmatureTestInstructionContainers: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving ImmatureTestInstructions from database",
				ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Get users ImmatureTestInstructionContainer-data from CloudDB
	cloudDBImmatureTestInstructionContainerItems, err = fenixCloudDBObject.LoadClientsImmatureTestInstructionContainersFromCloudDB(userID)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
			ImmatureTestInstructions:          nil,
			ImmatureTestInstructionContainers: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving ImmatureTestInstructionContainers from database",
				ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
		ImmatureTestInstructions:          cloudDBImmatureTestInstructionItems,
		ImmatureTestInstructionContainers: cloudDBImmatureTestInstructionContainerItems,
		AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		},
	}

	return responseMessage, nil
}
