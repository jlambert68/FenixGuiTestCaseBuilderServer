package messagesToWorkerServer

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"crypto/tls"
	fenixExecutionWorkerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixExecutionServer/fenixExecutionWorkerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/idtoken"
	grpcMetadata "google.golang.org/grpc/metadata"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"golang.org/x/net/context"
)

// ********************************************************************************************************************

// SetConnectionToFenixGuiBuilderServer - Set upp connection and Dial to FenixExecutionServer
func (messagesToWorkerServerObject *MessagesToWorkerServerObjectStruct) SetConnectionToFenixGuiBuilderServer(
	workerServerAddressToDial string) (err error) {

	// slice with sleep time, in milliseconds, between each attempt to Dial to Server
	var sleepTimeBetweenDialAttempts []int
	sleepTimeBetweenDialAttempts = []int{100, 100, 200, 200, 300, 300, 500, 500, 600, 1000} // Total: 3.6 seconds

	var opts []grpc.DialOption

	// Do multiple attempts to do connection to Execution Server
	var numberOfDialAttempts int
	var dialAttemptCounter int
	numberOfDialAttempts = len(sleepTimeBetweenDialAttempts)
	dialAttemptCounter = 0

	for {

		dialAttemptCounter = dialAttemptCounter + 1

		//When running on GCP then use credential otherwise not
		if common_config.ExecutionLocationForBuilderServer == common_config.GCP {
			creds := credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			})

			opts = []grpc.DialOption{
				grpc.WithTransportCredentials(creds),
			}
		}

		// Set up connection to Fenix Execution Server
		// When run on GCP, use credentials
		if common_config.ExecutionLocationForBuilderServer == common_config.GCP {
			// Run on GCP
			remoteFenixWorkerServerConnection, err = grpc.Dial(workerServerAddressToDial, opts...)

			if err != nil {
				messagesToWorkerServerObject.Logger.WithFields(logrus.Fields{
					"ID":                        "1cb282c4-4864-42b6-b943-89d4ed5b5300",
					"workerServerAddressToDial": workerServerAddressToDial,
					"error message":             err,
					"dialAttemptCounter":        dialAttemptCounter,
				}).Error("Couldn't dial WorkerServer")

			} else {
				messagesToWorkerServerObject.Logger.WithFields(logrus.Fields{
					"ID":                        "7f946da0-c377-4cda-bbe5-8abb612f75b1",
					"workerServerAddressToDial": workerServerAddressToDial,
					"dialAttemptCounter":        dialAttemptCounter,
				}).Debug("Success in dialing WorkerServer")
			}

		} else {
			// Run Local
			remoteFenixWorkerServerConnection, err = grpc.Dial(common_config.LocalFenixWorkerServerAddressToDial, grpc.WithInsecure())
			if err != nil {
				messagesToWorkerServerObject.Logger.WithFields(logrus.Fields{
					"ID": "85757210-ab12-4e81-9f1c-5efc438a9cd1",
					"common_config.LocalFenixWorkerServerAddressToDial": common_config.LocalFenixWorkerServerAddressToDial,
					"error message":      err,
					"dialAttemptCounter": dialAttemptCounter,
				}).Error("Couldn't dial WorkerServer")

			} else {
				messagesToWorkerServerObject.Logger.WithFields(logrus.Fields{
					"ID": "f9efb46d-d896-4049-8cd5-aed832fee861",
					"common_config.LocalFenixWorkerServerAddressToDial": common_config.LocalFenixWorkerServerAddressToDial,
					"dialAttemptCounter": dialAttemptCounter,
				}).Debug("Success in dialing WorkerServer")
			}
		}

		// Only return the error after last attempt
		if dialAttemptCounter >= numberOfDialAttempts {
			return err
		}

		// uccess in dialing WorkerServer
		if err == nil {
			// Creates a new gRPC-Client
			fenixWorkerServerGrpcWorkerServicesClient = fenixExecutionWorkerGrpcApi.
				NewFenixBuilderGprcServicesClient(remoteFenixWorkerServerConnection)

			return err

		}

		// Sleep for some time before retrying to connect
		time.Sleep(time.Millisecond * time.Duration(sleepTimeBetweenDialAttempts[dialAttemptCounter-1]))
	}
}

// Generate Google access token. Used when running in GCP
func (messagesToWorkerServerObject *MessagesToWorkerServerObjectStruct) generateGCPAccessToken(
	ctx context.Context,
	workerServerAddressToDial string) (
	appendedCtx context.Context,
	returnAckNack bool,
	returnMessage string) {

	// Only create the token if there is none, or it has expired
	if messagesToWorkerServerObject.gcpAccessToken == nil || messagesToWorkerServerObject.gcpAccessToken.Expiry.Before(time.Now()) {

		// Create an identity token.
		// With a global TokenSource tokens would be reused and auto-refreshed at need.
		// A given TokenSource is specific to the audience.
		tokenSource, err := idtoken.NewTokenSource(ctx, "https://"+workerServerAddressToDial)
		if err != nil {
			messagesToWorkerServerObject.Logger.WithFields(logrus.Fields{
				"ID":  "1c9e7d77-ae3f-403a-9bbd-1574b5858e88",
				"err": err,
			}).Error("Couldn't generate access token")

			return nil, false, "Couldn't generate access token"
		}

		token, err := tokenSource.Token()
		if err != nil {
			messagesToWorkerServerObject.Logger.WithFields(logrus.Fields{
				"ID":  "b5f6e3f1-a93c-4df9-9cda-708acaffa1f9",
				"err": err,
			}).Error("Problem getting the token")

			return nil, false, "Problem getting the token"
		} else {
			messagesToWorkerServerObject.Logger.WithFields(logrus.Fields{
				"ID": "378b6abd-23d1-4048-965a-699fa36f0f50",
				//"token": token,
			}).Debug("Got Bearer Token")
		}

		messagesToWorkerServerObject.gcpAccessToken = token

	}

	messagesToWorkerServerObject.Logger.WithFields(logrus.Fields{
		"ID": "50366c49-36f5-4b0a-93fb-a60349b134fb",
	}).Debug("Will use Bearer Token")

	// Add token to GrpcServer Request.
	appendedCtx = grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+messagesToWorkerServerObject.gcpAccessToken.AccessToken)

	return appendedCtx, true, ""

}

// ********************************************************************************************************************
// Get the highest ClientProtoFileVersionEnumeration for Execution Worker
func getHighestExecutionWorkerProtoFileVersion() int32 {

	// Check if there already is a 'highestclientProtoFileVersion' saved, if so use that one
	if highestExecutionWorkerProtoFileVersion != -1 {
		return highestExecutionWorkerProtoFileVersion
	}

	// Find the highest value for proto-file version
	var maxValue int32
	maxValue = 0

	for _, v := range fenixExecutionWorkerGrpcApi.CurrentFenixExecutionWorkerProtoFileVersionEnum_value {
		if v > maxValue {
			maxValue = v
		}
	}

	highestExecutionWorkerProtoFileVersion = maxValue

	return highestExecutionWorkerProtoFileVersion
}
