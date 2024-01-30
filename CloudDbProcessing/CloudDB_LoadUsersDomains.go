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

	// Concatenate Users specific Domains and Domains open for every one to use
	domainAndAuthorizations, err = fenixCloudDBObject.concatenateUsersDomainsAndDomainOpenToEveryOneToUse(
		txn, gCPAuthenticatedUser)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "7fb5b108-b74d-485d-8eb4-401bc77a2ee4",
			"error":                err,
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Got problem extracting users Domains")

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

// Concatenate Users specific Domains and Domains open for every one to use
func (fenixCloudDBObject *FenixCloudDBObjectStruct) concatenateUsersDomainsAndDomainOpenToEveryOneToUse(
	dbTransaction pgx.Tx,
	gCPAuthenticatedUser string) (
	domainAndAuthorizations []DomainAndAuthorizationsStruct,
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "687e4d77-3ff3-4fc6-acdc-da39b6b05bc0",
	}).Debug("Entering: concatenateUsersDomainsAndDomainOpenToEveryOneToUse()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":                      "1dd6aa13-1b6a-454e-8a9b-4ab4f1280f95",
			"domainAndAuthorizations": domainAndAuthorizations,
		}).Debug("Exiting: concatenateUsersDomainsAndDomainOpenToEveryOneToUse()")
	}()

	// Load all domains open for every one to use in some way
	var domainsOpenForEveryOneToUse []DomainAndAuthorizationsStruct
	domainsOpenForEveryOneToUse, err = fenixCloudDBObject.loadDomainsOpenForEveryOneToUse(dbTransaction, gCPAuthenticatedUser)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "c3d160ea-9122-46fc-a483-0afa54ba45d2",
			"error":                err,
			"GCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Couldn't load all Domains open for every one to use from CloudDB")

		return nil, err
	}

	// Load all domains for a specific user
	domainAndAuthorizations, err = fenixCloudDBObject.loadUsersDomains(dbTransaction, gCPAuthenticatedUser)
	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "b3a80162-e78e-48c2-a7a1-796a4d2df9f2",
			"error":                err,
			"GCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Error("Couldn't load all Users Domains from CloudDB")

		return nil, err
	}

	// Concatenate Domains and Authorizations
	var domainMap map[string]DomainAndAuthorizationsStruct
	domainMap = make(map[string]DomainAndAuthorizationsStruct)

	// Loop 'domainAndAuthorizations' and add to Map
	for _, termpDomainAndAuthorization := range domainAndAuthorizations {
		domainMap[termpDomainAndAuthorization.DomainUuid] = termpDomainAndAuthorization
	}

	// Loop 'domainsOpenForEveryOneToUse' and add to Map if they don't already exist, if so then replace certain values
	var existsInMap bool
	var termpDomainAndAuthorization DomainAndAuthorizationsStruct
	for _, tempdomainOpenForEveryOneToUse := range domainsOpenForEveryOneToUse {

		termpDomainAndAuthorization, existsInMap = domainMap[tempdomainOpenForEveryOneToUse.DomainUuid]
		if existsInMap == false {
			// Add to Map
			domainMap[tempdomainOpenForEveryOneToUse.DomainUuid] = tempdomainOpenForEveryOneToUse

		} else {
			// Replace values
			termpDomainAndAuthorization.CanBuildAndSaveTestCaseHavingTIandTICFromThisDomain =
				tempdomainOpenForEveryOneToUse.CanBuildAndSaveTestCaseHavingTIandTICFromThisDomain
			termpDomainAndAuthorization.CanListAndViewTestCaseHavingTIandTICFromThisDomain =
				tempdomainOpenForEveryOneToUse.CanListAndViewTestCaseHavingTIandTICFromThisDomain
		}
	}

	// Clear and rebuild 'domainAndAuthorizations'
	domainAndAuthorizations = nil
	for _, tempDomainAndAuthorizations := range domainMap {
		domainAndAuthorizations = append(domainAndAuthorizations, tempDomainAndAuthorizations)
	}

	return domainAndAuthorizations, err

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

// Load all domains open for every one to use in some way
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadDomainsOpenForEveryOneToUse(
	dbTransaction pgx.Tx,
	gCPAuthenticatedUser string) (
	domainsOpenForEveryOneToUse []DomainAndAuthorizationsStruct,
	err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "a7634a49-d32e-4b15-b2ca-68e86f7a6983",
	}).Debug("Entering: loadDomainsOpenForEveryOneToUse()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":                          "2c8b80dc-1058-4dd2-bb50-277b5952e731",
			"domainsOpenForEveryOneToUse": domainsOpenForEveryOneToUse,
		}).Debug("Exiting: loadDomainsOpenForEveryOneToUse()")
	}()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT dom.domain_uuid, dom.domain_name, " +
		"dom.\"AllUsersCanListAndViewTestCaseHavingTIandTICFromThisDomain\", " +
		"dom.\"AllUsersCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain\", dbpn.bitNumberValue "
	sqlToExecute = sqlToExecute + "FROM \"FenixDomainAdministration\".\"domains\" dom, " +
		"\"FenixDomainAdministration\".\"domainbitpositionenum\" dbpn "
	sqlToExecute = sqlToExecute + "WHERE (\"AllUsersCanListAndViewTestCaseHavingTIandTICFromThisDomain\" = true OR "
	sqlToExecute = sqlToExecute + "\"AllUsersCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain\" = true) AND "
	sqlToExecute = sqlToExecute + "dom.\"bitnumbername\" = dbpn.\"bitnumbername\" "
	sqlToExecute = sqlToExecute + ";"

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "9ac11b92-7536-4234-af11-7ee112f620d9",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadDomainsOpenForEveryOneToUse'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "54e8c6d1-490b-4576-bdd6-cb453013f21d",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return nil, err
	}

	var tempCanListAndViewTestCaseHavingTIandTICFromThisDomain bool
	var tempCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain bool
	var bitNumberValue int64

	// Extract data from DB result set
	for rows.Next() {

		var tempDomainAndAuthorizations DomainAndAuthorizationsStruct

		err = rows.Scan(
			&tempDomainAndAuthorizations.DomainUuid,
			&tempDomainAndAuthorizations.DomainName,
			&tempCanListAndViewTestCaseHavingTIandTICFromThisDomain,
			&tempCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain,
			&bitNumberValue,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "ede844c1-1f90-479e-8032-c0365f5bfb97",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Convert bool to int64 for 'tempCanListAndViewTestCaseHavingTIandTICFromThisDomain'
		if tempCanListAndViewTestCaseHavingTIandTICFromThisDomain == true {
			tempDomainAndAuthorizations.CanListAndViewTestCaseHavingTIandTICFromThisDomain = bitNumberValue
		}

		// Convert bool to int64 for 'tempCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain'
		if tempCanBuildAndSaveTestCaseHavingTIandTICFromThisDomain == true {
			tempDomainAndAuthorizations.CanBuildAndSaveTestCaseHavingTIandTICFromThisDomain = bitNumberValue
		}

		// Add user to the row-data
		tempDomainAndAuthorizations.GCPAuthenticatedUser = gCPAuthenticatedUser

		// Append DomainUuid to list of Domains
		domainsOpenForEveryOneToUse = append(domainsOpenForEveryOneToUse, tempDomainAndAuthorizations)

	}

	return domainsOpenForEveryOneToUse, err
}
