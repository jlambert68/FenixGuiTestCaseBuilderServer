package gRPCapiServer

import (
	fenixGuiTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"net"
	//	ecpb "github.com/jlambert68/FenixGrpcApi/Client/fenixGuiTestCaseBuilderServerGrpcApi/echo/go_grpc_api"
)

type fenixGuiTestCaseBuilderServerObjectStruct struct {
	Logger         *logrus.Logger
	gcpAccessToken *oauth2.Token
}

// Variable holding everything together
var fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct

// gRPC variables
var (
	registerFenixTestCaseBuilderServerGrpcServicesServer *grpc.Server // registerFenixTestCaseBuilderServerGrpcServicesServer *grpc.Server
	lis                                                  net.Listener
)

// gRPC Server used for register clients Name, Ip and Por and Clients Test Enviroments and Clients Test Commandst
type fenixTestCaseBuilderServerGrpcServicesServerStruct struct {
	fenixGuiTestCaseBuilderServerGrpcApi.UnimplementedFenixTestCaseBuilderServerGrpcServicesServer
}

// gRPC Server used for register clients Name, Ip and Por and Clients Test Enviroments and Clients Test Commandst
type fenixTestCaseBuilderServerGrpcWorkerServicesServerStruct struct {
	fenixGuiTestCaseBuilderServerGrpcApi.UnimplementedFenixTestCaseBuilderServerGrpcWorkerServicesServer
}
