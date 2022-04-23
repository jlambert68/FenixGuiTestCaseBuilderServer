package main

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"strconv"

	//"flag"
	"fmt"
	"log"
	"os"
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

func main() {
	//time.Sleep(15 * time.Second)
	fenixGuiTestCaseBuilderServerMain()
}

func init() {
	//executionLocationForClient := flag.String("startupType", "0", "The application should be started with one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
	//flag.Parse()

	var err error

	// Get Environment variable to tell how this program was started
	var executionLocationForClient = mustGetenv("ExecutionLocationForClient")

	switch executionLocationForClient {
	case "LOCALHOST_NODOCKER":
		common_config.ExecutionLocationForClient = common_config.LocalhostNoDocker

	case "LOCALHOST_DOCKER":
		common_config.ExecutionLocationForClient = common_config.LocalhostDocker

	case "GCP":
		common_config.ExecutionLocationForClient = common_config.GCP

	default:
		fmt.Println("Unknown Execution location for Client: " + executionLocationForClient + ". Expected one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
		os.Exit(0)

	}

	// Get Environment variable to tell where Fenix TestData Sync Server is started
	var executionLocationForFenixTestDataServer = mustGetenv("ExecutionLocationForFenixTestDataServer")

	switch executionLocationForFenixTestDataServer {
	case "LOCALHOST_NODOCKER":
		common_config.ExecutionLocationForFenixTestDataServer = common_config.LocalhostNoDocker

	case "LOCALHOST_DOCKER":
		common_config.ExecutionLocationForFenixTestDataServer = common_config.LocalhostDocker

	case "GCP":
		common_config.ExecutionLocationForFenixTestDataServer = common_config.GCP

	default:
		fmt.Println("Unknown Execution location for Fenix TestData Syn Server: " + executionLocationForFenixTestDataServer + ". Expected one of the following: LOCALHOST_NODOCKER, LOCALHOST_DOCKER, GCP")
		os.Exit(0)

	}

	// Extract all other Environment variables
	// Address to Fenix TestData Sync server
	common_config.FenixTestDataSyncServerAddress = mustGetenv("FenixTestDataSyncServerAddress")

	// Port for Fenix TestData Sync server
	common_config.FenixTestDataSyncServerPort, err = strconv.Atoi(mustGetenv("FenixTestDataSyncServerPort"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'FenixTestDataSyncServerPort' to an integer, error: ", err)
		os.Exit(0)

	}

	// Address to Client TestData Sync server
	common_config.ClientTestDataSyncServerAddress = mustGetenv("ClientTestDataSyncServerAddress")

	// Port for Client TestData Sync server
	common_config.ClientTestDataSyncServerPort, err = strconv.Atoi(mustGetenv("ClientTestDataSyncServerPort"))
	if err != nil {
		fmt.Println("Couldn't convert environment variable 'ClientTestDataSyncServerPort' to an integer, error: ", err)
		os.Exit(0)

	}

	// Create the Dial up string to Fenix TestData SyncServer
	fenixGuiTestCaseBuilderServer_address_to_dial = common_config.FenixTestDataSyncServerAddress + ":" + strconv.Itoa(common_config.FenixTestDataSyncServerPort)

}
