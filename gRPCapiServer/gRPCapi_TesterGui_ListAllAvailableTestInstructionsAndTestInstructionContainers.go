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
	responseMessage *fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage,
	err error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "a55f9c82-1d74-44a5-8662-058b8bc9e48f",
	}).Debug("Incoming 'gRPC - ListAllAvailableTestInstructionsAndTestContainers'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "27fb45fe-3266-41aa-a6af-958513977e28",
	}).Debug("Outgoing 'gRPC - ListAllAvailableTestInstructionsAndTestContainers'")

	// Current user
	var gCPAuthenticatedUser string
	var userIdOnComputer string
	gCPAuthenticatedUser = userIdentificationMessage.GCPAuthenticatedUser
	userIdOnComputer = userIdentificationMessage.UserIdOnComputer

	// Check if Client is using correct proto files version
	var returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
	returnMessage = common_config.IsClientUsingCorrectTestDataProtoFileVersion(
		userIdOnComputer, userIdentificationMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.
			AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
			ImmatureTestInstructions:          nil,
			ImmatureTestInstructionContainers: nil,
			AckNackResponse:                   returnMessage,
		}

		// Exiting
		return responseMessage, nil
	}

	// Define variables to store data from DB in
	var cloudDBImmatureTestInstructionItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage
	var cloudDBExecutionDomainsThatCanReceiveDirectTargetedTestInstructions []*fenixTestCaseBuilderServerGrpcApi.
		ExecutionDomainsThatCanReceiveDirectTargetedTestInstructionsMessage
	var cloudDBImmatureTestInstructionContainerItems []*fenixTestCaseBuilderServerGrpcApi.
		ImmatureTestInstructionContainerMessage

	// Temporary map to hinder duplicates of ExecutionDomains
	var executionDomainsMap map[string]*fenixTestCaseBuilderServerGrpcApi.
		ExecutionDomainsThatCanReceiveDirectTargetedTestInstructionsMessage
	executionDomainsMap = make(map[string]*fenixTestCaseBuilderServerGrpcApi.
		ExecutionDomainsThatCanReceiveDirectTargetedTestInstructionsMessage)

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
			DomainsThatCanOwnTheTestCase:                                 nil,
			ImmatureTestInstructions:                                     nil,
			ImmatureTestInstructionContainers:                            nil,
			ExecutionDomainsThatCanReceiveDirectTargetedTestInstructions: nil,
			ImmatureTestInstructionAttributes:                            nil,
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
				DomainsThatCanOwnTheTestCase:                                 nil,
				ImmatureTestInstructions:                                     nil,
				ImmatureTestInstructionContainers:                            nil,
				ExecutionDomainsThatCanReceiveDirectTargetedTestInstructions: nil,
				ImmatureTestInstructionAttributes:                            nil,
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
			var tempExecutionDomainThatCanReceiveDirectTargetedTestInstructions *fenixTestCaseBuilderServerGrpcApi.ExecutionDomainsThatCanReceiveDirectTargetedTestInstructionsMessage
			tempImmatureTestInstructionMessage,
				tempExecutionDomainThatCanReceiveDirectTargetedTestInstructions,
				err = s.convertSupportedTestInstructionsIntoABBResultTI(
				tempTestInstructions.TestInstructionVersions[0].TestInstructionInstance,
				tempTestInstructions.TestInstructionVersions[0].ResponseVariablesMapStructure)

			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "9a3be6ab-0486-47df-a8ef-9f1d30503817",
					"err": err,
					"tempTestInstructions.TestInstructionVersions[0].TestInstructionInstance": tempTestInstructions.TestInstructionVersions[0].TestInstructionInstance,
				}).Error("Couldn't convert TestInstruction into gRPC version to be sent to TesterGui")

				responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
					DomainsThatCanOwnTheTestCase:                                 nil,
					ImmatureTestInstructions:                                     nil,
					ImmatureTestInstructionContainers:                            nil,
					ExecutionDomainsThatCanReceiveDirectTargetedTestInstructions: nil,
					ImmatureTestInstructionAttributes:                            nil,
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
			cloudDBImmatureTestInstructionItems = append(
				cloudDBImmatureTestInstructionItems, tempImmatureTestInstructionMessage)

			executionDomainsMap[tempExecutionDomainThatCanReceiveDirectTargetedTestInstructions.ExecutionDomainUuid] =
				tempExecutionDomainThatCanReceiveDirectTargetedTestInstructions

		}

		// Loop all TestInstructionContainers from 'this domain'
		for _, tempTestInstructionContainer := range testInstructionsAndTestInstructionContainersFromGrpcBuilderMessage.
			TestInstructionContainers.TestInstructionContainersMap {

			// Convert TestInstructionContainer. Slice position '0' is always the latest one so use that
			var tempImmatureTestInstructionContainerMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage
			tempImmatureTestInstructionContainerMessage, err = s.convertSupportedTestInstructionContainersIntoABBResultTIC(
				tempTestInstructionContainer.TestInstructionContainerVersions[0].TestInstructionContainerInstance)

			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":  "730c9ffd-f0fb-40d6-b5c9-335e82556520",
					"err": err,
					"tempTestInstructionContainer.TestInstructionContainerVersions[0]." +
						"TestInstructionContainerInstance": tempTestInstructionContainer.
						TestInstructionContainerVersions[0].TestInstructionContainerInstance,
				}).Error("Couldn't convert TestInstructionContainer into gRPC version to be sent to TesterGui")

				responseMessage = &fenixTestCaseBuilderServerGrpcApi.
					AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
					DomainsThatCanOwnTheTestCase:                                 nil,
					ImmatureTestInstructions:                                     nil,
					ImmatureTestInstructionContainers:                            nil,
					ExecutionDomainsThatCanReceiveDirectTargetedTestInstructions: nil,
					ImmatureTestInstructionAttributes:                            nil,
					AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
						AckNack:  false,
						Comments: "Couldn't convert TestInstructionContainer into gRPC version to be sent to TesterGui",
						ErrorCodes: []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{
							fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
						ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
							CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
								GetHighestFenixGuiBuilderProtoFileVersion()),
					},
				}

				return responseMessage, err
			}

			// Append to list of TestInstructions to be sent to TesterGui
			cloudDBImmatureTestInstructionContainerItems = append(
				cloudDBImmatureTestInstructionContainerItems, tempImmatureTestInstructionContainerMessage)

		}
	}

	// Extract User's Domains that can own a TestCase by looping Domains and check which one that can own a TestCase
	var domainsThatCanOwnTheTestCase []*fenixTestCaseBuilderServerGrpcApi.DomainsThatCanOwnTheTestCaseMessage

	for _, domainAndAuthorization := range domainAndAuthorizations {
		if domainAndAuthorization.CanBuildAndSaveTestCaseOwnedByThisDomain > 0 {

			// When value is set then the Domain can own a TestCase
			var tempDomainsThatCanOwnTheTestCase *fenixTestCaseBuilderServerGrpcApi.DomainsThatCanOwnTheTestCaseMessage
			tempDomainsThatCanOwnTheTestCase = &fenixTestCaseBuilderServerGrpcApi.DomainsThatCanOwnTheTestCaseMessage{
				DomainUuid: domainAndAuthorization.DomainUuid,
				DomainName: domainAndAuthorization.DomainName,
			}

			// Add to lists of Domains that can own a TestCase
			domainsThatCanOwnTheTestCase = append(domainsThatCanOwnTheTestCase, tempDomainsThatCanOwnTheTestCase)
		}
	}

	// Convert map with ExecutionDomains to slice of ExecutionDomains
	var availableExecutionDomains []string
	for _, executionDomainThatCanReceiveDirectTargetedTestInstructions := range executionDomainsMap {

		// If ExecutionDomainUUid is a ZeroUuid then skip that ExecutionDomain
		// The reason is that ZeroUuid is used to indicate that a TestInstruction can have a dynamic ExecutionDomain and are set in TesterGui
		//if executionDomainThatCanReceiveDirectTargetedTestInstructions.ExecutionDomainUuid != common_config.ZeroUuid {

		// Add to slice with alla information about ExecutionDomains
		cloudDBExecutionDomainsThatCanReceiveDirectTargetedTestInstructions = append(
			cloudDBExecutionDomainsThatCanReceiveDirectTargetedTestInstructions,
			executionDomainThatCanReceiveDirectTargetedTestInstructions)

		// Add to simple slice of ExecutionDomains
		availableExecutionDomains = append(availableExecutionDomains, executionDomainThatCanReceiveDirectTargetedTestInstructions.GetNameUsedInGui())

		//}

	}

	// Define variables to store Attribute data from DB in
	var testInstructionAttributesList []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionAttributesMessage_TestInstructionAttributeMessage

	// Convert structured TestInstructionAttributes data into "raw" list of supported attributes
	testInstructionAttributesList, err = s.convertSupportedTestInstructionsAttributesIntoAttributesList(
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages, availableExecutionDomains)

	// Create the response to caller
	responseMessage = &fenixTestCaseBuilderServerGrpcApi.
		AvailableTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
		DomainsThatCanOwnTheTestCase:                                 domainsThatCanOwnTheTestCase,
		ImmatureTestInstructions:                                     cloudDBImmatureTestInstructionItems,
		ImmatureTestInstructionContainers:                            cloudDBImmatureTestInstructionContainerItems,
		ExecutionDomainsThatCanReceiveDirectTargetedTestInstructions: cloudDBExecutionDomainsThatCanReceiveDirectTargetedTestInstructions,
		ImmatureTestInstructionAttributes: &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionAttributesMessage{
			TestInstructionAttributesList: testInstructionAttributesList},
		AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    true,
			Comments:   "",
			ErrorCodes: nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		},
	}

	return responseMessage, nil
}
