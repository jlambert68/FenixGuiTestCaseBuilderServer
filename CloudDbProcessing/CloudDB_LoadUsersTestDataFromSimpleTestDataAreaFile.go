package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

// Do initial preparations to be able to load all domains for a specific user
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareLoadUsersTestDataFromSimpleTestDataAreaFile(
	gCPAuthenticatedUser string) (
	listAllTestDataForTestDataAreasResponseMessage *fenixTestCaseBuilderServerGrpcApi.ListAllTestDataForTestDataAreasResponseMessage) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "63b0e8ce-2fbf-4353-9b2e-858360db4a26",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'PrepareLoadUsersTestDataFromSimpleTestDataAreaFile'")

		listAllTestDataForTestDataAreasResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListAllTestDataForTestDataAreasResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   err.Error(),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			TestDataFromSimpleTestDataAreaFiles: nil,
		}

		return listAllTestDataForTestDataAreasResponseMessage

	}

	defer txn.Commit(context.Background())

	// Load Domains that User has access to
	var domainAndAuthorizations []DomainAndAuthorizationsStruct
	domainAndAuthorizations, err = fenixCloudDBObject.PrepareLoadUsersDomains(gCPAuthenticatedUser)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "8ddfe224-6702-40a5-90b4-048bf80c6987",
			"error":                err,
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Got some problem when loading users domains from database")

		listAllTestDataForTestDataAreasResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListAllTestDataForTestDataAreasResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   err.Error(),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			TestDataFromSimpleTestDataAreaFiles: nil,
		}

		return listAllTestDataForTestDataAreasResponseMessage

	}

	// If user doesn't have access to any domains then exit with warning in log
	if len(domainAndAuthorizations) == 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "a6b88092-a04e-4f3c-811d-4320ee09d0d1",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Warning("User doesn't have access to any domains")

		listAllTestDataForTestDataAreasResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListAllTestDataForTestDataAreasResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   fmt.Sprintf("User %s doesn't have access to any domains", gCPAuthenticatedUser),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			TestDataFromSimpleTestDataAreaFiles: nil,
		}

		return listAllTestDataForTestDataAreasResponseMessage

	}

	// Extract User's Domains that can own a TestCase by looping Domains and check which one that can own a TestCase
	var domainsThatCanOwnTheTestCase []*fenixTestCaseBuilderServerGrpcApi.DomainsThatCanOwnTheTestCaseMessage

	for _, domainAndAuthorization := range domainAndAuthorizations {
		if domainAndAuthorization.CanBuildAndSaveTestCaseOwnedByThisDomain > 0 {

			// When value is set then the Domain can own a TestCase
			var tempDomainsThatCanOwnTheTestCase *fenixTestCaseBuilderServerGrpcApi.DomainsThatCanOwnTheTestCaseMessage
			tempDomainsThatCanOwnTheTestCase = &fenixTestCaseBuilderServerGrpcApi.DomainsThatCanOwnTheTestCaseMessage{
				DomainUuid: domainAndAuthorization.DomainUuid,
				DomainName: domainAndAuthorization.DomainName,
			}

			// Add to lists of Domains that can own a TestCase
			domainsThatCanOwnTheTestCase = append(domainsThatCanOwnTheTestCase, tempDomainsThatCanOwnTheTestCase)
		}
	}

	// Check if there are any Domains that can own TestCases
	if len(domainsThatCanOwnTheTestCase) == 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "8b42795c-a124-459e-8f1b-3dda63bc2c22",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Warning("User doesn't have access to any domains that can own TestCase. This will automatically give that no 'simple' TestData can be loaded")

		listAllTestDataForTestDataAreasResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListAllTestDataForTestDataAreasResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   fmt.Sprintf("User %s doesn't have access to any domains that can own TestCase and therefor no 'simple' TestData can be loaded", gCPAuthenticatedUser),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			TestDataFromSimpleTestDataAreaFiles: nil,
		}

		return listAllTestDataForTestDataAreasResponseMessage

	}

	// Load all TestData for the Domains
	var testDataFromSimpleTestDataAreaFiles []*fenixTestCaseBuilderServerGrpcApi.TestDataFromOneSimpleTestDataAreaFileMessage
	testDataFromSimpleTestDataAreaFiles, err = fenixCloudDBObject.loadUsersTestDataFromSimpleTestDataAreaFile(
		txn,
		domainsThatCanOwnTheTestCase)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "43fc33a6-16d8-445a-95c5-04fa2df22546",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Got some error when loading TestData from Database")

		listAllTestDataForTestDataAreasResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListAllTestDataForTestDataAreasResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   fmt.Sprintf("Got some error when loading TestData from Database for User '%s'", gCPAuthenticatedUser),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			TestDataFromSimpleTestDataAreaFiles: nil,
		}

		return listAllTestDataForTestDataAreasResponseMessage
	}

	// Create ResponseMessage
	listAllTestDataForTestDataAreasResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
		ListAllTestDataForTestDataAreasResponseMessage{
		AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    true,
			Comments:   "",
			ErrorCodes: nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
					GetHighestFenixGuiBuilderProtoFileVersion()),
		},
		TestDataFromSimpleTestDataAreaFiles: testDataFromSimpleTestDataAreaFiles,
	}

	return listAllTestDataForTestDataAreasResponseMessage

}

// Load 'simple' TestData from Database
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadUsersTestDataFromSimpleTestDataAreaFile(
	dbTransaction pgx.Tx,
	domainsToGetTemplateApiUrlFor []*fenixTestCaseBuilderServerGrpcApi.DomainsThatCanOwnTheTestCaseMessage) (
	testDataFromSimpleTestDataAreaFileMessages []*fenixTestCaseBuilderServerGrpcApi.TestDataFromOneSimpleTestDataAreaFileMessage,
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id":                            "0fa60f59-2822-4731-9146-a3d6deb5b067",
		"domainsToGetTemplateApiUrlFor": domainsToGetTemplateApiUrlFor,
	}).Debug("Entering: loadUsersTestDataFromSimpleTestDataAreaFile")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "382220d7-2d79-45bf-9ccc-5297b7066e9d",
		}).Debug("Exiting: loadUsersTestDataFromSimpleTestDataAreaFile")
	}()

	// Generate SQLINArray containing DomainUuids
	var sQLINArray string
	var domainSlice []string

	// Loop Domains and add to dataSlice
	for _, domainsTemplateApiUrl := range domainsToGetTemplateApiUrlFor {
		domainSlice = append(domainSlice, domainsTemplateApiUrl.GetDomainUuid())
	}

	// create the IN-array...('sdada', 'adadadf')
	sQLINArray = fenixCloudDBObject.generateSQLINArray(domainSlice)

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT \"TestDataDomainUuid\", \"TestDataDomainName\", \"TestDataDomainTemplateName\"," +
		" \"TestDataAreaUuid\", \"TestDataAreaName\", " +
		"\"TestDataFileSha256Hash\", \"ImportantDataInFileSha256Hash\", \"InsertedTimeStamp\", " +
		"\"TestDataFromOneSimpleTestDataAreaFileFullMessage\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestDataFromSimpleTestDataAreaFile\" "
	sqlToExecute = sqlToExecute + "WHERE \"TestDataDomainUuid\" IN " + sQLINArray + ""
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "4e86cd4b-6a13-4565-8919-6e25dbb63b78",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadUsersTestDataFromSimpleTestDataAreaFile'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "5461d4ab-461c-40a4-8b34-8eca12c26066",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return nil, err
	}

	var insertedTimeStampAsTimeStamp time.Time
	var tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsString string
	var tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsStringAsByteArray []byte

	// Extract data from DB result set
	for rows.Next() {

		var oneTestDataFromOneSimpleTestDataAreaFileMessage fenixTestCaseBuilderServerGrpcApi.TestDataFromOneSimpleTestDataAreaFileMessage
		var oneTestDataFromOneSimpleTestDataAreaFileFullMessage fenixTestCaseBuilderServerGrpcApi.TestDataFromOneSimpleTestDataAreaFileMessage

		err = rows.Scan(
			&oneTestDataFromOneSimpleTestDataAreaFileMessage.TestDataDomainUuid,
			&oneTestDataFromOneSimpleTestDataAreaFileMessage.TestDataDomainName,
			&oneTestDataFromOneSimpleTestDataAreaFileMessage.TestDataDomainTemplateName,
			&oneTestDataFromOneSimpleTestDataAreaFileMessage.TestDataAreaUuid,
			&oneTestDataFromOneSimpleTestDataAreaFileMessage.TestDataAreaName,
			&oneTestDataFromOneSimpleTestDataAreaFileMessage.TestDataFileSha256Hash,
			&oneTestDataFromOneSimpleTestDataAreaFileMessage.ImportantDataInFileSha256Hash,
			&insertedTimeStampAsTimeStamp,
			&tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsString,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "b5124d87-9d58-4844-873d-73a3284f4389",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert json-string into byte-arrays
		tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsStringAsByteArray = []byte(tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsString)

		// Convert json-byte-array into proto-messages
		err = protojson.Unmarshal(tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsStringAsByteArray, &oneTestDataFromOneSimpleTestDataAreaFileFullMessage)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "89bfcfac-d89d-40ae-8d47-0ddafc127138",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestDataFromOneSimpleTestDataAreaFileFullMessageAsStringAsByteArray' into proto-message")

			return nil, err
		}

		// Add TemplateRepositoryConnectionParameters to list
		testDataFromSimpleTestDataAreaFileMessages = append(testDataFromSimpleTestDataAreaFileMessages, &oneTestDataFromOneSimpleTestDataAreaFileFullMessage)

	}

	return testDataFromSimpleTestDataAreaFileMessages, err
}
