package main

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"fmt"
	"log"
	"os"
	"strconv"
)

// mustGetEnv is a helper function for getting environment variables.
// Displays a warning if the environment variable is not set.
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

func init() {
	//executionLocation := flag.String("startupType", "0", "The application should be started with one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
	//flag.Parse()

	var err error

	// Get Environment variable to tell how this program was started
	var executionLocation = mustGetenv("ExecutionLocation")

	switch executionLocation {
	case "LOCALHOST_NODOCKER":
		common_config.ExecutionLocationForClient = common_config.LocalhostNoDocker

	case "LOCALHOST_DOCKER":
		common_config.ExecutionLocationForClient = common_config.LocalhostDocker

	case "GCP":
		common_config.ExecutionLocationForClient = common_config.GCP

	default:
		fmt.Println("Unknown Execution location for FenixGuiServer: " + executionLocation + ". Expected one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
		os.Exit(0)

	}

	/*
		// Address to GuiBuilderServer - Not needed
		common_config.FenixGuiServerAddress = mustGetenv("FenixGuiBuilderServerAddress")
	*/

	// Port for GuiBuilderServer
	common_config.FenixGuiServerPort, err = strconv.Atoi(mustGetenv("FenixGuiBuilderServerPort"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'FenixGuiBuilderServerPort' to an integer, error: ", err)
		os.Exit(0)

	}

	// Max number of DB-connection from Pool. Not stored because it is re-read when connecting the DB-pool
	_ = mustGetenv("DB_POOL_MAX_CONNECTIONS")

	_, err = strconv.Atoi(mustGetenv("DB_POOL_MAX_CONNECTIONS"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'DB_POOL_MAX_CONNECTIONS' to an integer, error: ", err)
		os.Exit(0)

	}
	// Should all SQL-queries be logged before executed
	var tempBoolAsString string
	var tempBool bool
	tempBoolAsString = mustGetenv("LogAllSQLs")
	tempBool, err = strconv.ParseBool(tempBoolAsString)
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'LogAllSQLs' to a boolean, error: ", err)
		os.Exit(0)
	}
	common_config.LogAllSQLs = tempBool

}
