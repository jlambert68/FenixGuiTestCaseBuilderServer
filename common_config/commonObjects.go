package common_config

import "github.com/sirupsen/logrus"

// ***********************************************************************************************************
// The following variables receives their values from environment variables

// Where is the client running
var ExecutionLocationForClient ExecutionLocationTypeType

// Where is the Fenix TestDataSync server running
// LocationForFenixTestDataServer
var ExecutionLocationForFenixTestDataServer ExecutionLocationTypeType

// Definitions for where client and Fenix Server is running
type ExecutionLocationTypeType int

// Constants used for where stuff is running
const (
	LocalhostNoDocker ExecutionLocationTypeType = iota
	LocalhostDocker
	GCP
)

// FenixGuiBuilderServer
var LocationForFenixGuiBuilderServerTypeMapping = map[ExecutionLocationTypeType]string{
	LocalhostNoDocker: "LOCALHOST_NODOCKER",
	LocalhostDocker:   "LOCALHOST_DOCKER",
	GCP:               "GCP",
}

// Address to Fenix TestData Server & Client, will have their values from Environment variables at startup
var (
	//FenixGuiServerAddress string // TODO remove, but is referenced by code that is not removed yet
	FenixGuiServerPort int
)

// LogAllSQLs - Environment variable for deciding if SQLs should be logged
var LogAllSQLs bool

// ***********************************************************************************************************

var highestFenixProtoFileVersion int32 = -1

var Logger *logrus.Logger
