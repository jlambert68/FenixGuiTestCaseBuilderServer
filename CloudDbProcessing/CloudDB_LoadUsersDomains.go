package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"github.com/jackc/pgx/v4"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"time"
)

// Do initial preparations to be able to load all domains for a specific user
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareLoadUsersDomains(
	gCPAuthenticatedUser string) (
	domainAndAuthorizations []DomainAndAuthorizationsStruct,
	err error) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "de61ddc1-d9ff-4f5a-b2bf-5dcfbe7ac619",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'prepareLoadUsersDomains'")

		return nil, err

	}

	defer txn.Commit(context.Background())

	// Load all domains for a specific user
	domainAndAuthorizations, err = fenixCloudDBObject.loadUsersDomains(txn, gCPAuthenticatedUser)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "2cd14dd8-0d44-4fc4-ae72-adaaaf37bacc",
			"error":                err,
			"GCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Couldn't load all Users Domains from CloudDB")

		return nil, err
	}

	return domainAndAuthorizations, err
}

// Used for holding a Users domain and the Authorizations for that Domain
type DomainAndAuthorizationsStruct struct {
	GCPAuthenticatedUser                                       string
	DomainUuid                                                 string
	DomainName                                                 string
	CanListAndViewTestCaseOwnedByThisDomain                    int64
	CanBuildAndSaveTestCaseOwnedByThisDomain                   int64
	CanListAndViewTestCaseHavingTIandTICFromThisDomain         int64
	CanListAndViewTestCaseHavingTIandTICFromThisDomainExtended int64
	CanBuildAndSaveTestCaseHavingTIandTICFromThisDomain        int64
}

// Load all domains for a specific user
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadUsersDomains(
	dbTransaction pgx.Tx,
	gCPAuthenticatedUser string) (
	domainAndAuthorizations []DomainAndAuthorizationsStruct,
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
