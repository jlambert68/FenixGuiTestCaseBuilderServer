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

const templateApiUrlBase string = "%s/%s/%s/contents%s"

// Do initial preparations to be able to load all domains for a specific user
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareLoadUsersTemplateRepositoryUrls(
	gCPAuthenticatedUser string) (
	allRepositoryApiUrlsResponseMessage *fenixTestCaseBuilderServerGrpcApi.ListAllRepositoryApiUrlsResponseMessage,
	err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "f7b71d48-2991-4402-800a-1cd88260c646",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'PrepareLoadUsersTemplateRepositoryUrls'")

		return nil, err

	}

	defer txn.Commit(context.Background())

	// Load Domains that User has access to
	var domainAndAuthorizations []DomainAndAuthorizationsStruct
	domainAndAuthorizations, err = fenixCloudDBObject.PrepareLoadUsersDomains(gCPAuthenticatedUser)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "8acf400f-8cd0-4e3b-b020-931be4dc3ad4",
			"error":                err,
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Got some problem when loading users domains from database")

		allRepositoryApiUrlsResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListAllRepositoryApiUrlsResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   err.Error(),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			RepositoryApiUrls: nil,
		}

		return allRepositoryApiUrlsResponseMessage, err

	}

	// If user doesn't have access to any domains then exit with warning in log
	if len(domainAndAuthorizations) == 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "8de8e0b7-ab8e-4fd5-9821-760113934cc3",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Warning("User doesn't have access to any domains")

		allRepositoryApiUrlsResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListAllRepositoryApiUrlsResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   fmt.Sprintf("User %s doesn't have access to any domains", gCPAuthenticatedUser),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			RepositoryApiUrls: nil,
		}

		return allRepositoryApiUrlsResponseMessage, err

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
			"id":                   "70775a53-38ce-41f4-a08e-38b37c0a8ddb",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Warning("User doesn't have access to any domains that can own TestCase. This will automatically give that no templates can be loaded")

		allRepositoryApiUrlsResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListAllRepositoryApiUrlsResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   fmt.Sprintf("User %s doesn't have access to any domains that can own TestCase and therefor no tempaltes can be loaded", gCPAuthenticatedUser),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			RepositoryApiUrls: nil,
		}

		return allRepositoryApiUrlsResponseMessage, err

	}

	// Load all parameters to be able to construct the TemplateApiUrls
	var templateRepositoryConnectionParameters []*fenixTestCaseBuilderServerGrpcApi.TemplateRepositoryConnectionParameters
	templateRepositoryConnectionParameters, err = fenixCloudDBObject.loadUsersTemplateRepositoryUrlParameters(
		txn,
		domainsThatCanOwnTheTestCase)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "5c403e57-b581-478a-a25f-7b2230a93931",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Got some error when loading Template-url-parameters from Database")

		allRepositoryApiUrlsResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListAllRepositoryApiUrlsResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   fmt.Sprintf("Got some error when loading Template-url-parameters from Database for User '%s'", gCPAuthenticatedUser),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			RepositoryApiUrls: nil,
		}

		return allRepositoryApiUrlsResponseMessage, err
	}

	// Check if any template-parameters was found in database
	if len(templateRepositoryConnectionParameters) == 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "96f9a028-621e-462e-8863-cbffa32fe571",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Warning("Didn't find any Template-parameter for user, which shouldn't happen")

		allRepositoryApiUrlsResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
			ListAllRepositoryApiUrlsResponseMessage{
			AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   fmt.Sprintf("Didn't find any Template-parameter for user '%s', which shouldn't happen", gCPAuthenticatedUser),
				ErrorCodes: nil,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
						GetHighestFenixGuiBuilderProtoFileVersion()),
			},
			RepositoryApiUrls: nil,
		}

		return allRepositoryApiUrlsResponseMessage, err
	}

	// Build Template-Urls to be sent back to TesterGui
	var repositoryApiUrlResponseMessages []*fenixTestCaseBuilderServerGrpcApi.RepositoryApiUrlResponseMessage
	var repositoryApiFullUrl string

	for _, templateRepositoryConnectionParameter := range templateRepositoryConnectionParameters {

		// Create the full 'RepositoryApiFullUrl'
		repositoryApiFullUrl = fmt.Sprintf(
			templateApiUrlBase,
			templateRepositoryConnectionParameter.GetRepositoryApiUrl(),
			templateRepositoryConnectionParameter.GetRepositoryOwner(),
			templateRepositoryConnectionParameter.GetRepositoryName(),
			templateRepositoryConnectionParameter.GetRepositoryPath())

		// Create one repositoryApiUrlResponseMessage and to slice of alla Urls
		var repositoryApiUrlResponseMessage *fenixTestCaseBuilderServerGrpcApi.RepositoryApiUrlResponseMessage
		repositoryApiUrlResponseMessage = &fenixTestCaseBuilderServerGrpcApi.RepositoryApiUrlResponseMessage{
			RepositoryApiUrlName: templateRepositoryConnectionParameter.GetRepositoryApiUrlName(),
			RepositoryApiFullUrl: repositoryApiFullUrl,
			GitHubApiKey:         templateRepositoryConnectionParameter.GetGitHubApiKey(),
		}

		repositoryApiUrlResponseMessages = append(repositoryApiUrlResponseMessages, repositoryApiUrlResponseMessage)

	}

	// Create ResponseMessage
	allRepositoryApiUrlsResponseMessage = &fenixTestCaseBuilderServerGrpcApi.
		ListAllRepositoryApiUrlsResponseMessage{
		AckNackResponse: &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    true,
			Comments:   "",
			ErrorCodes: nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.
					GetHighestFenixGuiBuilderProtoFileVersion()),
		},
		RepositoryApiUrls: repositoryApiUrlResponseMessages,
	}

	return allRepositoryApiUrlsResponseMessage, err

}

// Load all parameters to be able to construct the TemplateApiUrls
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadUsersTemplateRepositoryUrlParameters(
	dbTransaction pgx.Tx,
	domainsToGetTemplateApiUrlFor []*fenixTestCaseBuilderServerGrpcApi.DomainsThatCanOwnTheTestCaseMessage) (
	templateRepositoryConnectionParameters []*fenixTestCaseBuilderServerGrpcApi.TemplateRepositoryConnectionParameters,
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "8cbb9338-bb1e-45a0-b01b-ac7cb28fc52a",
	}).Debug("Entering: loadUsersDomains()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":                      "860b9b49-4d94-4f03-bf84-5b77f095ac7b",
			"domainAndAuthorizations": domainAndAuthorizations,
		}).Debug("Exiting: loadUsersDomains()")
	}()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT domainuuid, domainname, canlistandviewtestcaseownedbythisdomain, " +
		"canbuildandsavetestcaseownedbythisdomain, canlistandviewtestcasehavingtiandticfromthisdomain, " +
		"canlistandviewtestcasehavingtiandticfromthisdomainextended, canbuildandsavetestcasehavingtiandticfromthisdomain "
	sqlToExecute = sqlToExecute + "FROM \"FenixDomainAdministration\".\"allowedusers\" "
	sqlToExecute = sqlToExecute + "WHERE \"gcpauthenticateduser\" = '" + gCPAuthenticatedUser + "'"
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "d72d1a9c-079c-442f-bc0e-f95b557fd443",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadUsersDomains'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "9177c263-04d0-411d-8ac2-148279038fb3",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return nil, err
	}

	// Extract data from DB result set
	for rows.Next() {

		var tempDomainAndAuthorizations DomainAndAuthorizationsStruct

		err = rows.Scan(
			&tempDomainAndAuthorizations.DomainUuid,
			&tempDomainAndAuthorizations.DomainName,
			&tempDomainAndAuthorizations.CanListAndViewTestCaseOwnedByThisDomain,
			&tempDomainAndAuthorizations.CanBuildAndSaveTestCaseOwnedByThisDomain,
			&tempDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain,
			&tempDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomainExtended,
			&tempDomainAndAuthorizations.CanBuildAndSaveTestCaseHavingTIandTICFromThisDomain,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "b5ae17c4-fce1-4627-abd0-a6450cb17dd7",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Add user to the row-data
		tempDomainAndAuthorizations.GCPAuthenticatedUser = gCPAuthenticatedUser

		// Append DomainUuid to list of Domains
		domainAndAuthorizations = append(domainAndAuthorizations, tempDomainAndAuthorizations)

	}

	return domainAndAuthorizations, err
}
