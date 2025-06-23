package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"time"
)

// Do initial preparations to be able to load all domains for a specific user
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareLoadUsersTestCaseMetaData(
	gCPAuthenticatedUser string) (
	testCaseAndTestSuiteMetaDataResponseMessage *fenixTestCaseBuilderServerGrpcApi.ListTestCaseAndTestSuiteMetaDataResponseMessage) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "f7b71d48-2991-4402-800a-1cd88260c646",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'PrepareLoadUsersTestCaseMetaData'")

		testCaseAndTestSuiteMetaDataResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListTestCaseAndTestSuiteMetaDataResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   err.Error(),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			TestCaseAndTestSuiteMetaDataForDomains: nil,
		}

		return testCaseAndTestSuiteMetaDataResponseMessage

	}

	defer txn.Commit(context.Background())

	// Load Domains that User has access to
	var domainAndAuthorizations []DomainAndAuthorizationsStruct
	domainAndAuthorizations, err = fenixCloudDBObject.PrepareLoadUsersDomains(gCPAuthenticatedUser)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "05a41aae-8f43-44f5-9678-79cc7264da68",
			"error":                err,
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Got some problem when loading users domains from database")

		testCaseAndTestSuiteMetaDataResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListTestCaseAndTestSuiteMetaDataResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   err.Error(),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			TestCaseAndTestSuiteMetaDataForDomains: nil,
		}

		return testCaseAndTestSuiteMetaDataResponseMessage

	}

	// If user doesn't have access to any domains then exit with warning in log
	if len(domainAndAuthorizations) == 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "97ce0a8b-1647-4b69-9295-76b367e83052",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Warning("User doesn't have access to any domains")

		testCaseAndTestSuiteMetaDataResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListTestCaseAndTestSuiteMetaDataResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   fmt.Sprintf("User %s doesn't have access to any domains", gCPAuthenticatedUser),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			TestCaseAndTestSuiteMetaDataForDomains: nil,
		}

		return testCaseAndTestSuiteMetaDataResponseMessage

	}

	// Extract User's Domains that can own a TestCase by looping Domains and check which one that can own a TestCase
	var domainsThatCanOwnTheTestCase []*fenixTestCaseBuilderServerGrpcApi.DomainsThatCanOwnTheTestCaseMessage

	for _, domainAndAuthorization := range domainAndAuthorizations {
		if domainAndAuthorization.CanBuildAndSaveTestCaseOrTestSuiteOwnedByThisDomain > 0 {

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
			"id":                   "499cfd12-3dbe-4756-9567-a431ec158953",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Warning("User doesn't have access to any domains that can own TestCase. This will automatically give that no TestCaseMetaData can be loaded")

		testCaseAndTestSuiteMetaDataResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListTestCaseAndTestSuiteMetaDataResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   fmt.Sprintf("User %s doesn't have access to any domains that can own TestCase and therefor no TestCaseMetaData can be loaded", gCPAuthenticatedUser),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			TestCaseAndTestSuiteMetaDataForDomains: nil,
		}

		return testCaseAndTestSuiteMetaDataResponseMessage

	}

	// Load all parameters to be able to construct the TemplateApiUrls
	var testCaseMetaDataForDomains []*fenixTestCaseBuilderServerGrpcApi.TestCaseAndTestSuiteMetaDataForOneDomainMessage
	testCaseMetaDataForDomains, err = fenixCloudDBObject.loadUsersTestCaseAndTestSuiteMetaDataParameters(
		txn,
		domainsThatCanOwnTheTestCase)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "5c403e57-b581-478a-a25f-7b2230a93931",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Got some error when loading Template-url-parameters from Database")

		testCaseAndTestSuiteMetaDataResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListTestCaseAndTestSuiteMetaDataResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   fmt.Sprintf("Got some error when loading TestCaseMetaData-arameters from Database for User '%s'", gCPAuthenticatedUser),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			TestCaseAndTestSuiteMetaDataForDomains: nil,
		}

		return testCaseAndTestSuiteMetaDataResponseMessage
	}

	// Check if any template-parameters was found in database
	if len(testCaseMetaDataForDomains) == 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "dee2a5af-bf82-445f-8dbf-549bf07f8f68",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Warning("Didn't find any Template-parameter for user, which shouldn't happen")

		testCaseAndTestSuiteMetaDataResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListTestCaseAndTestSuiteMetaDataResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   fmt.Sprintf("Didn't find any TestCaseMetaData-parameter for user '%s', which shouldn't happen", gCPAuthenticatedUser),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			TestCaseAndTestSuiteMetaDataForDomains: nil,
		}

		return testCaseAndTestSuiteMetaDataResponseMessage
	}

	// Create ResponseMessage
	testCaseAndTestSuiteMetaDataResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
		ListTestCaseAndTestSuiteMetaDataResponseMessage{
		AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    true,
			Comments:   "",
			ErrorCodes: nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
					GetHighestFenixGuiBuilderProtoFileVersion()),
		},
		TestCaseAndTestSuiteMetaDataForDomains: testCaseMetaDataForDomains,
	}

	return testCaseAndTestSuiteMetaDataResponseMessage

}

// Load all TestCaseMetaData- and TestSuite-parameters for all Domains that the User has access to
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadUsersTestCaseAndTestSuiteMetaDataParameters(
	dbTransaction pgx.Tx,
	domainsToTestCaseMetaDataFor []*fenixTestCaseBuilderServerGrpcApi.DomainsThatCanOwnTheTestCaseMessage) (
	testCaseAndTestSuiteMetaDataForDomains []*fenixTestCaseBuilderServerGrpcApi.TestCaseAndTestSuiteMetaDataForOneDomainMessage,
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id":                           "b482cff4-72a3-4ac3-a695-271edfb1b5d1",
		"domainsToTestCaseMetaDataFor": domainsToTestCaseMetaDataFor,
	}).Debug("Entering: loadUsersTestCaseAndTestSuiteMetaDataParameters")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "ec1b3cf2-d59b-43c8-9a1a-4a480f86883c",
		}).Debug("Exiting: loadUsersTestCaseAndTestSuiteMetaDataParameters")
	}()

	// Generate SQLINArray containing DomainUuids
	var sQLINArray string
	var domainSlice []string

	// Loop Domains and add to dataSlice
	for _, domainTestCaseMetaData := range domainsToTestCaseMetaDataFor {
		domainSlice = append(domainSlice, domainTestCaseMetaData.GetDomainUuid())
	}

	// create the IN-array...('sdada', 'adadadf')
	sQLINArray = fenixCloudDBObject.generateSQLINArray(domainSlice)

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT \"DomainUuid\", \"DomainName\", \"SupportedTestCaseMetaData\", " +
		"\"SupportedTestSuiteMetaData\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"SupportedTestCaseAndTestSuiteMetaData\" "
	sqlToExecute = sqlToExecute + "WHERE \"DomainUuid\" IN " + sQLINArray + ""
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "8c65ac46-2fc6-4e8b-9d7a-e617d55fbdf7",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadUsersTestCaseAndTestSuiteMetaDataParameters'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "80178a89-910e-48db-aa8f-898899135d22",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return nil, err
	}

	// Extract data from DB result set
	for rows.Next() {

		var testCaseAndTestSuiteMetaDataForDomain *fenixTestCaseBuilderServerGrpcApi.TestCaseAndTestSuiteMetaDataForOneDomainMessage
		testCaseAndTestSuiteMetaDataForDomain = &fenixTestCaseBuilderServerGrpcApi.TestCaseAndTestSuiteMetaDataForOneDomainMessage{}

		err = rows.Scan(

			&testCaseAndTestSuiteMetaDataForDomain.DomainUuid,
			&testCaseAndTestSuiteMetaDataForDomain.DomainName,
			&testCaseAndTestSuiteMetaDataForDomain.TestCaseMetaDataAsJson,
			&testCaseAndTestSuiteMetaDataForDomain.TestSuiteMetaDataAsJson,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "3d00038f-c027-4f01-8f83-9b6ee04fe450",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Add 'testCaseAndTestSuiteMetaDataForDomain' to list
		testCaseAndTestSuiteMetaDataForDomains = append(testCaseAndTestSuiteMetaDataForDomains, testCaseAndTestSuiteMetaDataForDomain)

	}

	return testCaseAndTestSuiteMetaDataForDomains, err
}
