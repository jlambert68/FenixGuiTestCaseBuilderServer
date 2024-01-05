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

// ListAllImmatureTestInstructionAttributes - *********************************************************************
// The TestCase Builder asks for all TestInstructions Attributes that the user must add values to in TestCase
func (s *fenixTestCaseBuilderServerGrpcServicesServerStruct) ListAllImmatureTestInstructionAttributes(
	ctx context.Context,
	userIdentificationMessage *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (
	responseMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionAttributesMessage, err error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "4a40306f-a431-4498-914f-1a7d92e5b856",
	}).Debug("Incoming 'gRPC - ListAllImmatureTestInstructionAttributes'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "ae43d1a6-80a1-48de-a879-00184e154184",
	}).Debug("Outgoing 'gRPC - ListAllImmatureTestInstructionAttributes'")

	// Current user
	var gCPAuthenticatedUser string
	var userIdOnComputer string
	gCPAuthenticatedUser = userIdentificationMessage.GCPAuthenticatedUser
	userIdOnComputer = userIdentificationMessage.UserIdOnComputer

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(userIdOnComputer, userIdentificationMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionAttributesMessage{
			TestInstructionAttributesList: nil,
			AckNackResponse:               returnMessage,
		}

		// Exiting
		return responseMessage, nil
	}

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Load Domains that User has access to
	var domainAndAuthorizations []CloudDbProcessing.DomainAndAuthorizationsStruct
	domainAndAuthorizations, err = fenixCloudDBObject.PrepareLoadUsersDomains(gCPAuthenticatedUser)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "4cebdec1-86d2-46d5-9c2d-9bfd343c03ec",
			"error":                err,
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Got some problem when loading users domains from database")

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ImmatureTestInstructionAttributesMessage{
			TestInstructionAttributesList: nil,
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
			"id":                   "02a762dc-932b-4edd-a8b5-a6d0a53ba36b",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Warning("User doesn't have access to any domains")

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ImmatureTestInstructionAttributesMessage{
			TestInstructionAttributesList: nil,
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
				"id":                   "4b98acbc-31fb-49f6-9a75-3f7d16f50704",
				"error":                err,
				"gCPAuthenticatedUser": gCPAuthenticatedUser,
			}).Error("Got some problem when loading users published TestInstruction andTestInstructionContainers")

			responseMessage = &fenixTestCaseBuilderServerGrpcApi.
				ImmatureTestInstructionAttributesMessage{
				TestInstructionAttributesList: nil,
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

	// Define variables to store Attribute data from DB in
	var testInstructionAttributesList []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionAttributesMessage_TestInstructionAttributeMessage

	// Convert structured TestInstructionAttributes data into "raw" list of supported attributes
	testInstructionAttributesList, err = s.convertSupportedTestInstructionsAttributesIntoAttributesList(
		testInstructionsAndTestInstructionContainersFromGrpcBuilderMessages)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":  "cb607ccf-c9c9-4f2c-b788-a5b113fc885c",
			"err": err,
		}).Error("Couldn't convert Attributes belonging to TestInstructions into gRPC version to be sent to TesterGui")

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ImmatureTestInstructionAttributesMessage{
			TestInstructionAttributesList: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Couldn't convert Attributes belonging to TestInstructions into gRPC version to be sent to TesterGui",
				ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
			},
		}

		return responseMessage, err

	}

	// Create the response to caller
	responseMessage = &fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionAttributesMessage{
		TestInstructionAttributesList: testInstructionAttributesList,
		AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    true,
			Comments:   "",
			ErrorCodes: nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(
				common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		},
	}

	return responseMessage, nil
}
