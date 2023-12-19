package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/CloudDbProcessing"
	"FenixGuiTestCaseBuilderServer/common_config"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

/*

  // ************************ Data Used to build Available structure, which is what the user can chose from ************************

  // *** Get data from DB ***

  // The TestCase Builder asks for all TestInstructions and Pre-defined TestInstructionContainer that the user can add to a TestCase
  rpc ListAllAvailableTestInstructionsAndTestContainers (UserIdentificationMessage) returns (AvailableTestInstructionsAndPreCreatedTestContainersResponseMessage) {
  }

  // The TestCase Builder asks for which TestInstructions and Pre-defined TestInstructionContainer that the user has pinned in the GUI
  rpc ListAllAvailablePinnedTestInstructionsAndTestContainers (UserIdentificationMessage) returns (AvailablePinnedTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage) {
  }

  // The TestCase Builder asks for all Bonds-elements that can be used in the TestCase-model
  rpc ListAllAvailableBonds (UserIdentificationMessage) returns (ImmatureBondsMessage) {
  }


  // *** Send data to DB ***

  // The TestCase Builder sends all TestInstructions and Pre-defined TestInstructionContainer that the user has pinned in the GUI by the user
  rpc SaveAllPinnedTestInstructionsAndTestContainers (SavePinnedTestInstructionsAndPreCreatedTestInstructionContainersMessage) returns (AckNackResponse) {
  }
*/

// ListAllAvailableTestInstructionsAndTestInstructionContainers - *********************************************************************
// The TestCase Builder asks for all TestInstructions and Pre-defined TestInstructionContainer that the user can add to a TestCase
func (s *fenixTestCaseBuilderServerGrpcServicesServer) ListAllAvailableTestInstructionsAndTestInstructionContainers(
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

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion("666", userIdentificationMessage.ProtoFileVersionUsedByClient)
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

	// Current user
	userID := userIdentificationMessage.UserId

	// Define variables to store data from DB in
	var cloudDBImmatureTestInstructionItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage
	var cloudDBImmatureTestInstructionContainerItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

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

// ListAllAvailablePinnedTestInstructionsAndTestInstructionContainers - *********************************************************************
// The TestCase Builder asks for all Pinned TestInstructions and Pinned Pre-defined TestInstructionContainer that the user can add to a TestCase
func (s *fenixTestCaseBuilderServerGrpcServicesServer) ListAllAvailablePinnedTestInstructionsAndTestInstructionContainers(ctx context.Context, userIdentificationMessage *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (*fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage, error) {

	// Define the response message
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "5a72e9c7-602e-4a16-a551-961f96fac457",
	}).Debug("Incoming 'gRPC - ListAllAvailablePinnedTestInstructionsAndTestInstructionContainers'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "28a7d2e7-ebdc-4e98-a5e9-08491f1ff181",
	}).Debug("Outgoing 'gRPC - ListAllAvailablePinnedTestInstructionsAndTestInstructionContainers'")

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion("666", userIdentificationMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
			AvailablePinnedTestInstructions:                    nil,
			AvailablePinnedPreCreatedTestInstructionContainers: nil,
			AckNackResponse: returnMessage,
		}

		// Exiting
		return responseMessage, nil
	}

	// Current user
	userID := userIdentificationMessage.UserId

	// Define variables to store data from DB in
	var cloudDBPinnedTestInstructionMessages []*fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionMessage
	var cloudDBPinnedPreCreatedTestInstructionContainerMessages []*fenixTestCaseBuilderServerGrpcApi.AvailablePinnedPreCreatedTestInstructionContainerMessage

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Get users PinnedTestInstruction-data from CloudDB
	cloudDBPinnedTestInstructionMessages, err := fenixCloudDBObject.LoadClientsPinnedTestInstructionsFromCloudDB(userID)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
			AvailablePinnedTestInstructions:                    nil,
			AvailablePinnedPreCreatedTestInstructionContainers: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving PinnedTestInstructions from database",
				ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Get users PinnedPreCreatedTestInstructionContainer-data from CloudDB
	cloudDBPinnedPreCreatedTestInstructionContainerMessages, err = fenixCloudDBObject.LoadClientsPinnedTestInstructionContainersFromCloudDB(userID)
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
			AvailablePinnedTestInstructions:                    nil,
			AvailablePinnedPreCreatedTestInstructionContainers: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving PinnedPreCreatedTestInstructionContainers from database",
				ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionsAndPreCreatedTestInstructionContainersResponseMessage{
		AvailablePinnedTestInstructions:                    cloudDBPinnedTestInstructionMessages,
		AvailablePinnedPreCreatedTestInstructionContainers: cloudDBPinnedPreCreatedTestInstructionContainerMessages,
		AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		},
	}

	return responseMessage, nil
}

// SavePinnedTestInstructionsAndTestContainers - *********************************************************************
// The TestCase Builder sends all TestInstructions and Pre-defined TestInstructionContainer that the user has pinned in the GUI
func (s *fenixTestCaseBuilderServerGrpcServicesServer) SaveAllPinnedTestInstructionsAndTestInstructionContainers(ctx context.Context, pinnedTestInstructionsAndTestContainersMessage *fenixTestCaseBuilderServerGrpcApi.SavePinnedTestInstructionsAndPreCreatedTestInstructionContainersMessage) (*fenixTestCaseBuilderServerGrpcApi.AckNackResponse, error) {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "a93fb1bd-1a5b-4417-80c3-082d34267c06",
	}).Debug("Incoming 'gRPC - SavePinnedTestInstructionsAndTestContainers'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "981ad10a-2bfb-4a39-9b4d-35cac0d7481a",
	}).Debug("Outgoing 'gRPC - SavePinnedTestInstructionsAndTestContainers'")

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(pinnedTestInstructionsAndTestContainersMessage.UserId, pinnedTestInstructionsAndTestContainersMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		// Exiting
		return returnMessage, nil
	}

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Save Pinned TestInstructions and pre-created TestInstructionContainers to Cloud DB
	returnMessage = fenixCloudDBObject.PrepareSavePinnedTestInstructionsAndPinnedTestInstructionContainersToCloudDB(pinnedTestInstructionsAndTestContainersMessage)
	if returnMessage != nil {
		// Something went wrong when saving to database
		// Exiting
		return returnMessage, nil
	}

	return &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:                      true,
		Comments:                     "",
		ErrorCodes:                   nil,
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
	}, nil
}

// ListAllAvailableTestInstructionsAndTestInstructionContainers - *********************************************************************
// The TestCase Builder asks for all TestInstructions and Pre-defined TestInstructionContainer that the user can add to a TestCase
func (s *fenixTestCaseBuilderServerGrpcServicesServer) ListAllAvailableBonds(ctx context.Context, userIdentificationMessage *fenixTestCaseBuilderServerGrpcApi.UserIdentificationMessage) (*fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage, error) {

	// Define the response message
	var responseMessage *fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "ffb95b37-9ab0-4933-a53c-b7676a12c8f2",
	}).Debug("Incoming 'gRPC - ListAllAvailableBonds'")

	defer fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"id": "b46215a0-6620-428d-bee9-9fa4c5e4e98b",
	}).Debug("Outgoing 'gRPC - ListAllAvailableBonds'")

	// Check if Client is using correct proto files version
	returnMessage := common_config.IsClientUsingCorrectTestDataProtoFileVersion(userIdentificationMessage.UserId, userIdentificationMessage.ProtoFileVersionUsedByClient)
	if returnMessage != nil {
		// Not correct proto-file version is used
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage{
			ImmatureBonds:   nil,
			AckNackResponse: returnMessage,
		}

		// Exiting
		return responseMessage, nil
	}

	// Define variables to store data from DB in
	var cloudDBAvailableBonds []*fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage_ImmatureBondMessage

	// Initiate object forCloudDB-processing
	var fenixCloudDBObject *CloudDbProcessing.FenixCloudDBObjectStruct
	fenixCloudDBObject = &CloudDbProcessing.FenixCloudDBObjectStruct{}

	// Get users ImmatureTestInstruction-data from CloudDB
	cloudDBAvailableBonds, err := fenixCloudDBObject.LoadAvailableBondsFromCloudDB()
	if err != nil {
		// Something went wrong so return an error to caller
		responseMessage = &fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage{
			ImmatureBonds: nil,
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:                      false,
				Comments:                     "Got some Error when retrieving Available Bonds from database",
				ErrorCodes:                   []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum{fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM},
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
			},
		}

		// Exiting
		return responseMessage, nil
	}

	// Create the response to caller
	responseMessage = &fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage{
		ImmatureBonds: cloudDBAvailableBonds,
		AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      true,
			Comments:                     "",
			ErrorCodes:                   nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		},
	}

	return responseMessage, nil
}
