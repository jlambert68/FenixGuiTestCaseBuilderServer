package main

import (
	fenixGuiTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// *********************************************************************
//Fenix client can check if Fenix Testdata sync server is alive with this service
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) AreYouAlive(ctx context.Context, emptyParameter *fenixGuiTestCaseBuilderServerGrpcApi.EmptyParameter) (*fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse, error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "1ff67695-9a8b-4821-811d-0ab8d33c4d8b",
	}).Debug("Incoming 'AreYouAlive'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "9c7f0c3d-7e9f-4c91-934e-8d7a22926d84",
	}).Debug("Outgoing 'AreYouAlive'")

	return &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{AckNack: true, Comments: "I'am alive, from Client"}, nil
}

// *********************************************************************
// Fenix client can register itself with the Fenix Testdata sync server
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) SendMerkleHash(ctx context.Context, merkleHashMessage *fenixGuiTestCaseBuilderServerGrpcApi.EmptyParameter) (*fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse, error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "a55f9c82-1d74-44a5-8662-058b8bc9e48f",
	}).Debug("Incoming 'SendMerkleHash'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "27fb45fe-3266-41aa-a6af-958513977e28",
	}).Debug("Outgoing 'SendMerkleHash'")

	// Send MerkleHash to Fenix after sending return message back to caller
	fenixGuiTestCaseBuilderServerObject.SendMerkleHash()

	return &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{AckNack: true, Comments: ""}, nil
}

// *********************************************************************
// Fenix client can send TestData MerkleTree to Fenix Testdata sync server with this service
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) SendMerkleTree(ctx context.Context, merkleTreeMessage *fenixGuiTestCaseBuilderServerGrpcApi.EmptyParameter) (*fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse, error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "cffc25f0-b0e6-407a-942a-71fc74f831ac",
	}).Debug("Incoming 'SendMerkleTree'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "61e2c28d-b091-442a-b7f8-d2502d9547cf",
	}).Debug("Outgoing 'SendMerkleTree'")

	// Send MerkleTree to Fenix after sending return message back to caller
	defer fenixGuiTestCaseBuilderServerObject.SendMerkleTree()

	return &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{AckNack: true, Comments: ""}, nil
}

// *********************************************************************
// Fenix client can send TestDataHeaderHash to Fenix Testdata sync server with this service
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) SendTestDataHeaderHash(ctx context.Context, testDataHeaderMessage *fenixGuiTestCaseBuilderServerGrpcApi.EmptyParameter) (*fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse, error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "ff642667-2cbd-4f23-91eb-a6f8e76d9177",
	}).Debug("Incoming 'SendTestDataHeaderHash'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "2c24b079-1e0b-46e9-ad1f-d47e8ff0d3b4",
	}).Debug("Outgoing 'SendTestDataHeaderHash'")

	// Send TestDataHeaderHash to Fenix after sending return message back to caller
	defer fenixGuiTestCaseBuilderServerObject.SendTestDataHeaderHash()

	return &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{AckNack: true, Comments: ""}, nil
}

// *********************************************************************
// Fenix client can send TestDataHeaders to Fenix Testdata sync server with this service
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) SendTestDataHeaders(ctx context.Context, testDataHeaderMessage *fenixGuiTestCaseBuilderServerGrpcApi.EmptyParameter) (*fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse, error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "aee48999-12ad-4bb7-bc8a-96b62a8eeedf",
	}).Debug("Incoming 'SendTestDataHeaders'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "ca0b58a8-6d56-4392-8751-45906670e86b",
	}).Debug("Outgoing 'SendTestDataHeaders'")

	// Send TestDataHeaders to Fenix after sending return message back to caller
	defer fenixGuiTestCaseBuilderServerObject.SendTestDataHeaders()

	return &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{AckNack: true, Comments: ""}, nil
}

// *********************************************************************
// Fenix client can send TestData rows to Fenix Testdata sync server with this service
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) SendTestDataRows(ctx context.Context, merklePathsMessage *fenixGuiTestCaseBuilderServerGrpcApi.MerklePathsMessage) (*fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse, error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "2b1c8752-eb84-4c15-b8a7-22e2464e5168",
	}).Debug("Incoming 'SendTestDataRows'")

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id":                 "7c7cd700-953f-4e31-9ca8-e1a4262c62b8",
		"merklePathsMessage": merklePathsMessage,
	}).Debug("Requested TestData")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "755e8b4f-f184-4277-ad41-e041714c2ca8",
	}).Debug("Outgoing 'SendTestDataRows'")

	// Send requested TestDataRows to Fenix after sending return message back to caller
	defer fenixGuiTestCaseBuilderServerObject.SendTestDataRows(merklePathsMessage.MerklePath)

	return &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{AckNack: true, Comments: ""}, nil
}

// *********************************************************************
// Fenix client can send All TestData rows to Fenix Testdata sync server with this service
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) SendAllTestDataRows(ctx context.Context, emptyParameter *fenixGuiTestCaseBuilderServerGrpcApi.EmptyParameter) (*fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse, error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "7708888f-edb0-4b87-97b7-cb2ce3b93d4a",
	}).Debug("Incoming 'SendTestDataRows'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "7bc8a6bd-8d8e-4244-98bf-cd5ca686d3f2",
	}).Debug("Outgoing 'SendTestDataRows'")

	// Send all TestDataRows to Fenix after sending return message back to caller
	defer fenixGuiTestCaseBuilderServerObject.SendTestDataRows([]string{})

	return &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{AckNack: true, Comments: ""}, nil
}

// Fenix client can register itself with the Fenix Testdata sync server
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) RegisterTestDataClient(ctx context.Context, testDataClientInformationMessage *fenixGuiTestCaseBuilderServerGrpcApi.EmptyParameter) (*fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse, error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "5133b80b-6f3a-4562-9e62-1b3ceb169cc1",
	}).Debug("Incoming 'RegisterTestDataClient'")

	defer fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"id": "316dcd7e-2229-4a82-b15b-0f808c2dd8aa",
	}).Debug("Outgoing 'RegisterTestDataClient'")

	// Send Client registration to Fenix after sending return message back to caller
	defer fenixGuiTestCaseBuilderServerObject.SendMerkleHash()

	return &fenixGuiTestCaseBuilderServerGrpcApi.AckNackResponse{AckNack: true, Comments: ""}, nil
}

/*
func (s *FenixGuiTestCaseBuilderGrpcServicesServer) mustEmbedUnimplementedFenixClientTestDataGrpcServicesServer() {
	//TODO implement me
	panic("implement me")
}


*/
