package gRPCapiServer

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	fenixGuiTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
)

// InitGrpcServer - Set up and start Backend gRPC-server
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) InitGrpcServer() {

	var err error

	// Find first non allocated port from defined start port
	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"Id": "054bc0ef-93bb-4b75-8630-74e3823f71da",
	}).Info("Backend Server tries to start")

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"Id":                                 "ca3593b1-466b-4536-be91-5e038de178f4",
		"common_config.FenixGuiServerPort: ": common_config.FenixGuiServerPort,
	}).Info("Start listening on:")
	lis, err = net.Listen("tcp", ":"+strconv.Itoa(common_config.FenixGuiServerPort))

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
			"Id":    "ad7815b3-63e8-4ab1-9d4a-987d9bd94c76",
			"err: ": err,
		}).Error("failed to listen:")
	} else {
		fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
			"Id":                                 "ba070b9b-5d57-4c0a-ab4c-a76247a50fd3",
			"common_config.FenixGuiServerPort: ": common_config.FenixGuiServerPort,
		}).Info("Success in listening on port:")

	}

	// Creates a new RegisterWorkerServer gRPC server
	//go func() {
	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"Id": "b0ccffb5-4367-464c-a3bc-460cafed16cb",
	}).Info("Starting Backend gRPC Server")

	registerFenixTestCaseBuilderServerGrpcServicesServer = grpc.NewServer()
	fenixGuiTestCaseBuilderServerGrpcApi.RegisterFenixTestCaseBuilderServerGrpcServicesServer(
		registerFenixTestCaseBuilderServerGrpcServicesServer,
		&fenixTestCaseBuilderServerGrpcServicesServerStruct{})
	fenixGuiTestCaseBuilderServerGrpcApi.RegisterFenixTestCaseBuilderServerGrpcWorkerServicesServer(
		registerFenixTestCaseBuilderServerGrpcServicesServer,
		&fenixTestCaseBuilderServerGrpcWorkerServicesServerStruct{})

	// Register RouteGuide on the same server.
	reflection.Register(registerFenixTestCaseBuilderServerGrpcServicesServer)

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"Id":                                 "e843ece9-b707-4c60-b1d8-14464305e68f",
		"common_config.FenixGuiServerPort: ": common_config.FenixGuiServerPort,
	}).Info("registerFenixTestCaseBuilderServerGrpcServicesServer for " +
		"TestInstruction Backend Server started")
	registerFenixTestCaseBuilderServerGrpcServicesServer.Serve(lis)
	//}()

}

// StopGrpcServer - Stop Backend gRPC-server
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) StopGrpcServer() {

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{}).Info("Gracefull stop for: registerTaxiHardwareStreamServer")
	registerFenixTestCaseBuilderServerGrpcServicesServer.GracefulStop()

	fenixGuiTestCaseBuilderServerObject.Logger.WithFields(logrus.Fields{
		"common_config.FenixGuiServerPort: ": common_config.FenixGuiServerPort,
	}).Info("Close net.Listing")
	_ = lis.Close()

}
