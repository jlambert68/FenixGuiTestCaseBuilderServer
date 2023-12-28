package messagesToWorkerServer

import (
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
)

type MessagesToWorkerServerObjectStruct struct {
	Logger                            *logrus.Logger
	gcpAccessToken                    *oauth2.Token
	connectionToWorkerServerInitiated bool
}

// Variables used for contacting Fenix Worker Server
var (
	remoteFenixWorkerServerConnection         *grpc.ClientConn
	fenixWorkerServerGrpcWorkerServicesClient fenixExecutionWorkerGrpcApi.FenixBuilderGprcServicesClient
)

var highestExecutionWorkerProtoFileVersion int32 = -1
