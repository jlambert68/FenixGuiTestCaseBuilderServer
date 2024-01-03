package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"fmt"
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
	gCPAuthenticatedUser = userIdentificationMessage.GCPAuthenticatedUser

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
	var domainAndAuthorizations []CloudDbProcessing.DomainAndAuthorizationsStruct
	domainAndAuthorizations, err = fenixCloudDBObject.PrepareLoadUsersDomains(gCPAuthenticatedUser)

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

	// If user doesn't have access to any domains then exit with warning in log
	if len(domainAndAuthorizations) == 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "e9e86616-7484-4c6f-b7dd-96f223a24cb3",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Warning("User doesn't have access to any domains")

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.
			AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
			ImmatureTestInstructions:          nil,
			ImmatureTestInstructionContainers: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   fmt.Sprintf("User %s doesn't have access to any domains", gCPAuthenticatedUser),
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
		PrepareLoadDomainsSupportedTestInstructionsAndTestInstructionContainersAndAllowedUsers(domainAndAuthorizations)

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
		for _, tempTestInstructions := range testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.
			TestInstructions.TestInstructionsMap {

			// Convert TestInstruction. Slice position '0' is always the latest one so use that
			var tempImmatureTestInstructionMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage
			tempImmatureTestInstructionMessage, err = s.convertSupportedTestInstructionsIntoABBResultTI(tempTestInstructions.TestInstructionVersions[0].TestInstructionInstance)

			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "9a3be6ab-0486-47df-a8ef-9f1d30503817",
					"err": err,
					"tempTestInstructions.TestInstructionVersions[0].TestInstructionInstance": tempTestInstructions.TestInstructionVersions[0].TestInstructionInstance,
				}).Error("Couldn't convert TestInstruction into gRPC version to be sent to TesterGui")

				responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
					ImmatureTestInstructions:          nil,
					ImmatureTestInstructionContainers: nil,
					AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
						AckNack:                      false,
						Comments:                     "Couldn't convert TestInstruction into gRPC version to be sent to TesterGui",
						ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
						ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
					},
				}

				return responseMessage, err

			}

			// Append to list of TestInstructions to be sent to TesterGui
			cloudDBImmatureTestInstructionItems = append(cloudDBImmatureTestInstructionItems, tempImmatureTestInstructionMessage)

		}

		// Loop all TestInstructionContainers from 'this domain'
		for _, tempTestInstructionContainer := range testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.
			TestInstructionContainers.TestInstructionContainersMap {

			// Convert TestInstructionContainer. Slice position '0' is always the latest one so use that
			var tempImmatureTestInstructionContainerMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage
			tempImmatureTestInstructionContainerMessage, err = s.convertSupportedTestInstructionContainersIntoABBResultTIC(tempTestInstructionContainer.TestInstructionContainerVersions[0].TestInstructionContainerInstance)

			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "730c9ffd-f0fb-40d6-b5c9-335e82556520",
					"err": err,
					"tempTestInstructionContainer.TestInstructionContainerVersions[0].TestInstructionContainerInstance": tempTestInstructionContainer.TestInstructionContainerVersions[0].TestInstructionContainerInstance,
				}).Error("Couldn't convert TestInstructionContainer into gRPC version to be sent to TesterGui")

				responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
					ImmatureTestInstructions:          nil,
					ImmatureTestInstructionContainers: nil,
					AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
						AckNack:                      false,
						Comments:                     "Couldn't convert TestInstructionContainer into gRPC version to be sent to TesterGui",
						ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
						ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
					},
				}

				return responseMessage, err
			}

			// Append to list of TestInstructions to be sent to TesterGui
			cloudDBImmatureTestInstructionContainerItems = append(cloudDBImmatureTestInstructionContainerItems, tempImmatureTestInstructionContainerMessage)

		}
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
