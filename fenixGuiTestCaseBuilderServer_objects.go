package main

import (
	fenixGuiTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"net"
	//	ecpb "github.com/jlambert68/FenixGrpcApi/Client/fenixGuiTestCaseBuilderServerGrpcApi/echo/go_grpc_api"
)

type fenixGuiTestCaseBuilderServerObjectStruct struct {
	logger         *logrus.Logger
	gcpAccessToken *oauth2.Token
}

// Variable holding everything together
var fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct

// gRPC variables
var (
	registerfenixGuiTestCaseBuilderServerServer *grpc.Server
	lis                                         net.Listener
)

// gRPC Server used for register clients Name, Ip and Por and Clients Test Enviroments and Clients Test Commandst
type FenixGuiTestCaseBuilderGrpcServicesServer struct {
	fenixGuiTestCaseBuilderServerGrpcApi.UnimplementedFenixTestCaseBuilderServerGrpcServicesServer
}

//TODO FIXA DENNA PATH, HMMM borde köra i DB framöver
// For now hardcoded MerklePath
//var merkleFilterPath string = //"AccountEnvironment/ClientJuristictionCountryCode/MarketSubType/MarketName/" //SecurityType/"

var testFile_1 = "data/FenixRawTestdata_14rows_211216.csv"

var testFile_2 = "data/FenixRawTestdata_14rows_211216_change.csv"

var testFileSelection bool = true

var testFile = testFile_2

var highestFenixProtoFileVersion int32 = -1
var highestClientProtoFileVersion int32 = -1

// Echo gRPC-server
/*
type ecServer struct {
	echo.UnimplementedEchoServer
}


*/
