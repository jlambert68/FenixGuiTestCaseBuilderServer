package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
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

		//log.Println("Close DB_session: %v", DB_session)
		//DB_session.Close()
	}
}

func fenixGuiTestCaseBuilderServerMain() {

	// Connect to CloudDB
	fenixSyncShared.ConnectToDB()

	// Set up BackendObject
	fenixGuiTestCaseBuilderServerObject = &fenixGuiTestCaseBuilderServerObject_struct{
		fenixGuiTestCaseBuilderServer_TestDataClientUuid: fenixSyncShared.MustGetEnvironmentVariable("TestDataClientUuid"),
		fenixGuiTestCaseBuilderServer_DomainUuid:         fenixSyncShared.MustGetEnvironmentVariable("TestDomainUuid"),
		fenixGuiTestCaseBuilderServer_DomainName:         fenixSyncShared.MustGetEnvironmentVariable("TestDomainName"),
		merkleFilterPath:                                 fenixSyncShared.MustGetEnvironmentVariable("MerkleFilterPath"), //TODO Remove all references to HARDCODED merkleFilterPath
	}

	// Init logger
	fenixGuiTestCaseBuilderServerObject.InitLogger("")

	// TODO Endast för Test
	fenixGuiTestCaseBuilderServerObject.loadAllTestDataRowItemsForClientFromCloudDB(&cloudDBExposedTestDataRowItems)

	// Clean up when leaving. Is placed after logger because shutdown logs information
	defer cleanup()

	// TODO remove only for testing gRPC connection between Cloud Run containers at SEB-GCP
	/*
		go func() {
			// Sleep for 60 second
			fmt.Println("Sleep for 60 seconds")
			time.Sleep(60 * time.Second)

			// Printed after sleep is over
			fmt.Println("Sleep Over.....")

			fmt.Println("Try to do gRPC-call to Server")
			serverStatus, serverMessage := fenixGuiTestCaseBuilderServerObject.SendAreYouAliveToFenixTestDataServer()
			fmt.Println("serverStatus", serverStatus)
			fmt.Println("serverMessage", serverMessage)

		}()
	*/

	// TODO remove only for testing https connection between Cloud Run containers at SEB-network
	go func() {
		addr := flag.String("addr", ":4000", "HTTPS network address")
		certFile := flag.String("certfile", "server.crt", "certificate PEM file")
		keyFile := flag.String("keyfile", "server.key", "key PEM file")
		flag.Parse()

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/" {
				http.NotFound(w, req)
				return
			}
			fmt.Fprintf(w, "Proudly served, from GCP, with Go and HTTPS!")
		})

		srv := &http.Server{
			Addr:    *addr,
			Handler: mux,
			TLSConfig: &tls.Config{
				MinVersion:               tls.VersionTLS13,
				PreferServerCipherSuites: true,
			},
		}

		log.Printf("Starting server on %s", *addr)
		err := srv.ListenAndServeTLS(*certFile, *keyFile)
		log.Fatal(err)
	}()

	// Start Backend gRPC-server
	fenixGuiTestCaseBuilderServerObject.InitGrpcServer()

}
