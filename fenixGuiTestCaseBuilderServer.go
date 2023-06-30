package main

import (
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

// Used for only process cleanup once
var cleanupProcessed = false

func cleanup() {

	if cleanupProcessed == false {

		cleanupProcessed = true

		// Cleanup before close down application
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{}).Info("Clean up and shut down servers")

		// Stop Backend gRPC Server
		fenixGuiTestCaseBuilderServerObject.StopGrpcServer()

		// Close Database Connection
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "587cc9b8-88eb-422c-b419-53fa4c51ebce",
		}).Info("Closing Database connection")

		fenixSyncShared.DbPool.Close()

	}
}

func fenixGuiTestCaseBuilderServerMain() {

	// Connect to CloudDB
	fenixSyncShared.ConnectToDB()

	// Set up BackendObject
	fenixGuiTestCaseBuilderServerObject = &fenixGuiTestCaseBuilderServerObjectStruct{}

	// Init logger
	fenixGuiTestCaseBuilderServerObject.InitLogger("")

	// Clean up when leaving. Is placed after logger because shutdown logs information
	defer cleanup()

	// Start Backend gRPC-server
	fenixGuiTestCaseBuilderServerObject.InitGrpcServer()

}
