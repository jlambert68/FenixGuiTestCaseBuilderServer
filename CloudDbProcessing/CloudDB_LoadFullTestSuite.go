package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
	"strconv"
	"time"
)

// This date is used as delete date when the TestSuite is not deleted
const testSuiteNotDeletedDate = "2068-11-18"

// PrepareLoadFullTestSuite
// Load Full TestSuite from Database
func (fenixCloudDBObject *FenixCloudDBObjectStruct) PrepareLoadFullTestSuite(
	testSuiteUuidToLoad string,
	gCPAuthenticatedUser string) (
	responseMessage *fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse) {

	// Begin SQL Transaction
	txn, err := fenixSyncShared.DbPool.Begin(context.Background())

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "b9867b4a-dbf1-4d17-ad80-04265d1d4964",
			"error": err,
		}).Error("Problem to do 'DbPool.Begin'  in 'prepareSaveFullTestSuite'")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "Problem when Loading from database",
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
			AckNackResponse:   ackNackResponse,
			DetailedTestSuite: nil,
		}

		return responseMessage
	}

	defer txn.Commit(context.Background())

	// Load Domains that User has access to
	var domainAndAuthorizations []DomainAndAuthorizationsStruct
	domainAndAuthorizations, err = fenixCloudDBObject.PrepareLoadUsersDomains(gCPAuthenticatedUser)

	// If user doesn't have access to any domains then exit with warning in log
	if len(domainAndAuthorizations) == 0 {
		common_config.Logger.WithFields(logrus.Fields{
			"id":                   "4d865b21-b923-438c-bd63-250034583c88",
			"gCPAuthenticatedUser": gCPAuthenticatedUser,
		}).Warning("User doesn't have access to any domains")

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack: false,
			Comments: fmt.Sprintf("User %s doesn't have access to any domains",
				gCPAuthenticatedUser),
			ErrorCodes: nil,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
			AckNackResponse:   ackNackResponse,
			DetailedTestSuite: nil,
		}

		return responseMessage

	}

	// Load the TestSuite
	var fullTestSuiteMessage *fenixTestCaseBuilderServerGrpcApi.FullTestSuiteMessage
	fullTestSuiteMessage, err = fenixCloudDBObject.loadFullTestSuite(txn, testSuiteUuidToLoad, domainAndAuthorizations)

	// Error when retrieving TestSuite
	if err != nil {
		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "Problem when Loading from database",
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
			AckNackResponse:   ackNackResponse,
			DetailedTestSuite: nil,
		}

		return responseMessage
	}

	// Generate list with TestCasesUuid's in TestSuite
	var testCaseUuidsInTestSuite []string
	for _, tempTestCaseInTestSuite := range fullTestSuiteMessage.GetTestCasesInTestSuite().GetTestCasesInTestSuite() {
		testCaseUuidsInTestSuite = append(testCaseUuidsInTestSuite, tempTestCaseInTestSuite.GetTestCaseUuid())
	}

	// Load PreViews for all TestCases in TestSuite
	var tempTestCasesPreview []*fenixTestCaseBuilderServerGrpcApi.TestCasePreviewMessage
	tempTestCasesPreview, err = fenixCloudDBObject.loadTestCasesPreviewForTestSuite(txn, testCaseUuidsInTestSuite)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"id":    "2daadb5d-5c4e-47aa-b25d-c0d222b4fc0b",
			"error": err,
		}).Error("Got some problem when loading TestCasesPreView from database")

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "Problem when Loading from database",
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
			AckNackResponse:   ackNackResponse,
			DetailedTestSuite: nil,
		}

		return responseMessage
	}

	// Create TestSuite-preview and add to TestSuite-Object
	// Create the 'TestSuitePreview'
	var tempSelectedTestSuiteMetaDataValuesMap map[string]*fenixTestCaseBuilderServerGrpcApi.
		TestSuitePreviewStructureMessage_SelectedTestSuiteMetaDataValueMessage
	tempSelectedTestSuiteMetaDataValuesMap = make(map[string]*fenixTestCaseBuilderServerGrpcApi.
		TestSuitePreviewStructureMessage_SelectedTestSuiteMetaDataValueMessage)

	var tempTestSuitePreview *fenixTestCaseBuilderServerGrpcApi.TestSuitePreviewMessage
	tempTestSuitePreview = &fenixTestCaseBuilderServerGrpcApi.TestSuitePreviewMessage{
		TestSuitePreview: &fenixTestCaseBuilderServerGrpcApi.TestSuitePreviewStructureMessage{
			TestSuiteUuid:                 fullTestSuiteMessage.TestSuiteBasicInformation.GetTestSuiteUuid(),
			TestSuiteName:                 fullTestSuiteMessage.TestSuiteBasicInformation.GetTestSuiteName(),
			TestSuiteVersion:              strconv.Itoa(int(fullTestSuiteMessage.TestSuiteBasicInformation.GetTestSuiteVersion())),
			DomainUuidThatOwnTheTestSuite: fullTestSuiteMessage.TestSuiteBasicInformation.GetDomainUuid(),
			DomainNameThatOwnTheTestSuite: fullTestSuiteMessage.TestSuiteBasicInformation.GetDomainName(),
			TestSuiteDescription:          fullTestSuiteMessage.GetTestSuiteBasicInformation().GetTestSuiteDescription(),
			TestSuiteStructureObjects: &fenixTestCaseBuilderServerGrpcApi.
				TestSuitePreviewStructureMessage_TestSuiteStructureObjectMessage{TestCasePreViews: tempTestCasesPreview},
			LastSavedByUserOnComputer:          fullTestSuiteMessage.UpdatedByAndWhen.GetUserIdOnComputer(),
			LastSavedByUserGCPAuthorization:    fullTestSuiteMessage.UpdatedByAndWhen.GetGCPAuthenticatedUser(),
			LastSavedTimeStamp:                 fullTestSuiteMessage.UpdatedByAndWhen.GetUpdateTimeStamp(),
			SelectedTestSuiteMetaDataValuesMap: tempSelectedTestSuiteMetaDataValuesMap,
		},
		TestSuitePreviewHash: "",
	}

	// Generate 'SelectedTestSuiteMetaDataValuesMap'
	for _, tempMetaDataGroupMessagePtr := range fullTestSuiteMessage.TestSuiteMetaData.MetaDataGroupsMap {

		// Get the 'MetaDataGroupMessage' from ptr
		tempMetaDataGroupMessage := *tempMetaDataGroupMessagePtr

		// Loop 'MetaDataInGroupMap'
		for SelectedTestSuiteMetaDataValuesMap, tempMetaDataInGroupMessagePtr := range tempMetaDataGroupMessage.MetaDataInGroupMap {

			switch tempMetaDataInGroupMessagePtr.GetSelectType() {

			case fenixTestCaseBuilderServerGrpcApi.MetaDataSelectTypeEnum_MetaDataSelectType_SingleSelect:

				// Create 'SelectedTestSuiteMetaDataValueMessage' to be stored in Map
				var tempSelectedTestSuiteMetaDataValueMessage *fenixTestCaseBuilderServerGrpcApi.
					TestSuitePreviewStructureMessage_SelectedTestSuiteMetaDataValueMessage

				tempSelectedTestSuiteMetaDataValueMessage = &fenixTestCaseBuilderServerGrpcApi.
					TestSuitePreviewStructureMessage_SelectedTestSuiteMetaDataValueMessage{
					OwnerDomainUuid:   tempSelectedTestSuiteMetaDataValueMessage.GetOwnerDomainUuid(),
					OwnerDomainName:   tempSelectedTestSuiteMetaDataValueMessage.GetOwnerDomainName(),
					MetaDataGroupName: tempMetaDataInGroupMessagePtr.GetMetaDataGroupName(),
					MetaDataName:      tempMetaDataInGroupMessagePtr.GetMetaDataName(),
					MetaDataNameValue: tempMetaDataInGroupMessagePtr.GetSelectedMetaDataValueForSingleSelect(),
					SelectType:        tempSelectedTestSuiteMetaDataValueMessage.GetSelectType(),
					IsMandatory:       tempSelectedTestSuiteMetaDataValueMessage.GetIsMandatory(),
				}

				// Add value to map
				tempSelectedTestSuiteMetaDataValuesMap[SelectedTestSuiteMetaDataValuesMap] = tempSelectedTestSuiteMetaDataValueMessage

			case fenixTestCaseBuilderServerGrpcApi.MetaDataSelectTypeEnum_MetaDataSelectType_MultiSelect:

				// Loop all selected values and add to map
				for _, tempSelectedMetaDataValue := range tempMetaDataInGroupMessagePtr.GetSelectedMetaDataValuesForMultiSelect() {

					// Create 'SelectedTestSuiteMetaDataValueMessage' to be stored in Map
					var tempSelectedTestSuiteMetaDataValueMessage *fenixTestCaseBuilderServerGrpcApi.
						TestSuitePreviewStructureMessage_SelectedTestSuiteMetaDataValueMessage

					tempSelectedTestSuiteMetaDataValueMessage = &fenixTestCaseBuilderServerGrpcApi.
						TestSuitePreviewStructureMessage_SelectedTestSuiteMetaDataValueMessage{
						OwnerDomainUuid:   tempSelectedTestSuiteMetaDataValueMessage.GetOwnerDomainUuid(),
						OwnerDomainName:   tempSelectedTestSuiteMetaDataValueMessage.GetOwnerDomainName(),
						MetaDataGroupName: tempMetaDataInGroupMessagePtr.GetMetaDataGroupName(),
						MetaDataName:      tempMetaDataInGroupMessagePtr.GetMetaDataName(),
						MetaDataNameValue: tempSelectedMetaDataValue,
						SelectType:        tempSelectedTestSuiteMetaDataValueMessage.GetSelectType(),
						IsMandatory:       tempSelectedTestSuiteMetaDataValueMessage.GetIsMandatory(),
					}

					// Add value to map
					tempSelectedTestSuiteMetaDataValuesMap[SelectedTestSuiteMetaDataValuesMap] = tempSelectedTestSuiteMetaDataValueMessage

				}

			default:
				common_config.Logger.WithFields(logrus.Fields{
					"Id":                   "6972f543-5fba-414e-b9d2-079a374d0f48",
					"fullTestSuiteMessage": fullTestSuiteMessage,
					"tempMetaDataInGroupMessagePtr.GetSelectType()": tempMetaDataInGroupMessagePtr.GetSelectType(),
				}).Error("Unknown SelectType in 'MetaDataInGroupMap'")

				// Set Error codes to return message
				var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
				var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

				errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
				errorCodes = append(errorCodes, errorCode)

				// Create response message
				var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
				ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
					AckNack:    false,
					Comments:   "Problem when Loading from database",
					ErrorCodes: errorCodes,
					ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
						CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
				}

				responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
					AckNackResponse:   ackNackResponse,
					DetailedTestSuite: nil,
				}

				return responseMessage

			}
		}
	}

	// Add 'tempSelectedTestSuiteMetaDataValuesMap' to 'tempTestSuitePreview'
	tempTestSuitePreview.TestSuitePreview.SelectedTestSuiteMetaDataValuesMap = tempSelectedTestSuiteMetaDataValuesMap

	// Calculate 'TestSuitePreviewHash'
	var tempTestSuitePreviewHash string
	tempJson := protojson.Format(tempTestSuitePreview)
	tempTestSuitePreviewHash = common_config.HashSingleValue(tempJson)

	tempTestSuitePreview.TestSuitePreviewHash = tempTestSuitePreviewHash

	// Add TestCases Preview-object to TestSuite
	fullTestSuiteMessage.TestSuitePreview = tempTestSuitePreview

	/*
		// Add information about who first created the TestSuite
		err = fenixCloudDBObject.addTestSuiteCreator(txn, testSuiteUuidToLoad, fullTestSuiteMessage)

		// Error when retrieving TestSuite creator
		if err != nil {
			// Set Error codes to return message
			var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
			var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

			errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
			errorCodes = append(errorCodes, errorCode)

			// Create response message
			var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
			ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
				AckNack:    false,
				Comments:   "Problem when Loading from database",
				ErrorCodes: errorCodes,
				ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
					CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
			}

			responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
				AckNackResponse:   ackNackResponse,
				DetailedTestSuite: nil,
			}

			return responseMessage
		}
	*/

	// TestSuite
	if fullTestSuiteMessage == nil {

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_DATABASE_PROBLEM
		errorCodes = append(errorCodes, errorCode)

		// Create response message
		var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
		ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:    false,
			Comments:   "TestCase couldn't be found in Database or the user doesn't have access to the TestSuite",
			ErrorCodes: errorCodes,
			ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
				CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
		}

		responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
			AckNackResponse:   ackNackResponse,
			DetailedTestSuite: nil,
		}

		return responseMessage
	}

	// Create response message
	var ackNackResponse *fenixTestCaseBuilderServerGrpcApi.AckNackResponse
	ackNackResponse = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
		AckNack:    true,
		Comments:   "",
		ErrorCodes: nil,
		ProtoFileVersionUsedByClient: fenixTestCaseBuilderServerGrpcApi.
			CurrentFenixTestCaseBuilderProtoFileVersionEnum(common_config.GetHighestFenixGuiBuilderProtoFileVersion()),
	}

	responseMessage = &fenixTestCaseBuilderServerGrpcApi.GetDetailedTestSuiteResponse{
		AckNackResponse:   ackNackResponse,
		DetailedTestSuite: fullTestSuiteMessage,
	}

	return responseMessage
}

// Load the TestSuite with all its data
func (fenixCloudDBObject *FenixCloudDBObjectStruct) loadFullTestSuite(
	dbTransaction pgx.Tx,
	testSuiteUuidToLoad string,
	domainAndAuthorizations []DomainAndAuthorizationsStruct) (
	fullTestSuiteMessage *fenixTestCaseBuilderServerGrpcApi.FullTestSuiteMessage, err error) {

	// Generate a Domains list and Calculate the Authorization requirements
	var tempCalculatedDomainAndAuthorizations DomainAndAuthorizationsStruct
	var domainList []string
	for _, domainAndAuthorization := range domainAndAuthorizations {
		// Add to DomainList
		domainList = append(domainList, domainAndAuthorization.DomainUuid)

		// Calculate the Authorization requirements for...
		// TestCaseAuthorizationLevelOwnedByDomain
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteOwnedByThisDomain =
			tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteOwnedByThisDomain +
				domainAndAuthorization.CanListAndViewTestCaseOrTestSuiteOwnedByThisDomain

		// TestCaseAuthorizationLevelHavingTiAndTicWithDomain
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain =
			tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain +
				domainAndAuthorization.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain
	}

	// Convert Values into string for TestSuiteAuthorizationLevelOwnedByDomain
	var tempCanListAndViewTestSuiteOwnedByThisDomainAsString string
	tempCanListAndViewTestSuiteOwnedByThisDomainAsString = strconv.FormatInt(
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain, 10)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":    "04ce498d-d440-4eaf-8f0a-8ade191a47ba",
			"Error": err,
			"tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain": tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain,
		}).Error("Couldn't convert into string representation")

		return nil, err
	}

	// Convert Values into string for TestSuiteAuthorizationLevelOwnedByDomain
	var tempCanListAndViewTestSuiteHavingTIandTICfromThisDomainAsString string
	tempCanListAndViewTestSuiteHavingTIandTICfromThisDomainAsString = strconv.FormatInt(
		tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain, 10)

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":    "9f094b06-34e4-46d0-bb25-1f2bbedb1d2f",
			"Error": err,
			"tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain": tempCalculatedDomainAndAuthorizations.CanListAndViewTestCaseOrTestSuiteHavingTIandTICFromThisDomain,
		}).Error("Couldn't convert into string representation")

		return nil, err
	}

	/*
		sqlToExecute = sqlToExecute + 	"SELECT TS.\"DomainUuid\", TS.\"DomainName\", TS.\"TestSuiteUuid\", TS.\"TestSuiteName\", " +
		sqlToExecute = sqlToExecute + 				"TS.\"TestSuiteVersion\", TS.\"TestSuiteHash\", " +
		sqlToExecute = sqlToExecute + 			"TS.\"InsertTimeStamp\", TS.\"InsertedByUserIdOnComputer\", TS.\"InsertedByGCPAuthenticatedUser\", " +
		sqlToExecute = sqlToExecute + 			"TS.\"TestCasesInTestSuite\", TS.\"TestSuitePreview\", TS.\"TestSuiteMetaData\" "

		WITH uniquecounters AS (
		    SELECT Distinct ON ("TestSuiteUuid")  "UniqueCounter"
		    FROM "FenixBuilder"."TestSuites"
		    WHERE "TestSuiteUuid" =  '3e903bc2-57bc-4083-bc0d-184ddb6522d4'


		    ORDER BY "TestSuiteUuid", "UniqueCounter" DESC
		)
		SELECT t.*
		FROM "FenixBuilder"."TestSuites" t
		WHERE t."UniqueCounter" IN (SELECT * FROM uniquecounters) AND
		      t."TestSuiteIsDeleted" = false;

	*/

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "WITH uniquecounters AS ( "
	sqlToExecute = sqlToExecute + "SELECT Distinct ON (\"TestSuiteUuid\")  \"UniqueCounter\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestSuites\" "
	sqlToExecute = sqlToExecute + fmt.Sprintf("WHERE \"TestSuiteUuid\" = '%s' ", testSuiteUuidToLoad)
	sqlToExecute = sqlToExecute + "ORDER BY \"TestSuiteUuid\", \"UniqueCounter\" DESC "
	sqlToExecute = sqlToExecute + ") "

	sqlToExecute = sqlToExecute + "SELECT TS.\"DomainUuid\", da.\"domain_name\", TS.\"TestSuiteUuid\", TS.\"TestSuiteName\", " +
		"TS.\"TestSuiteVersion\", TS.\"TestSuiteDescription\", TS.\"TestSuiteExecutionEnvironment\", TS.\"TestSuiteHash\", TS.\"DeleteTimestamp\", " +
		"TS.\"InsertTimeStamp\", TS.\"InsertedByUserIdOnComputer\", TS.\"InsertedByGCPAuthenticatedUser\", " +
		"TS.\"TestCasesInTestSuite\", TS.\"TestSuitePreview\", TS.\"TestSuiteMetaData\", TS.\"TestSuiteTestData\", " +
		"TS.\"TestSuiteType\", TS.\"TestSuiteTypeName\", TS.\"TestSuiteImplementedFunctions\" "

	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestSuites\" TS, \"FenixDomainAdministration\".\"domains\" da  "
	sqlToExecute = sqlToExecute + fmt.Sprintf("WHERE TS.\"TestSuiteUuid\" = '%s' ", testSuiteUuidToLoad)
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "(TS.\"CanListAndViewTestSuiteAuthorizationLevelOwnedByDomain\" & " +
		tempCanListAndViewTestSuiteOwnedByThisDomainAsString + ") "
	sqlToExecute = sqlToExecute + "= TS.\"CanListAndViewTestSuiteAuthorizationLevelOwnedByDomain\" "
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "(TS.\"CanListAndViewTestSuiteAuthorizationLevelHavingTiAndTicWith\" & " +
		tempCanListAndViewTestSuiteHavingTIandTICfromThisDomainAsString + ")"
	sqlToExecute = sqlToExecute + "= TS.\"CanListAndViewTestSuiteAuthorizationLevelHavingTiAndTicWith\" "
	sqlToExecute = sqlToExecute + "AND "
	sqlToExecute = sqlToExecute + "TS.\"UniqueCounter\" IN (SELECT * FROM uniquecounters) AND "
	sqlToExecute = sqlToExecute + "TS.\"DeleteTimestamp\" > now() AND "
	sqlToExecute = sqlToExecute + "TS.\"DomainUuid\" = da.\"domain_uuid\" "
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "3b331ac4-7944-4037-ba1a-aa8cb4d38074",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'loadFullTestSuite'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "c66b4f71-8b00-482f-82d2-a8f1e622f08b",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return nil, err
	}

	var (
		tempDomainUuid                          string
		tempDomainName                          string
		tempTestSuiteUuid                       string
		tempTestSuiteName                       string
		tempTestSuiteDescription                string
		tempTestSuiteExecutionEnvironment       string
		tempTestSuiteVersion                    int
		tempTestSuiteHash                       string
		tempDeleteTimestampAsString             string
		tempUpdatedByAndWhenAsString            string
		tempInsertedByUserIdOnComputer          string
		tempInsertedByGCPAuthenticatedUser      string
		tempTestCasesInTestSuiteAsJson          string
		tempTestSuitePreviewAsJson              string
		tempTestSuiteMetaDataAsJson             string
		tempTestSuiteTestDataAsJson             string
		tempTestSuiteType                       int
		tempTestSuiteTypeName                   string
		tempTestSuiteImplementedFunctionsAsJson string

		tempDeleteTimestampAsTimeStamp  time.Time
		tempUpdatedByAndWhenAsTimeStamp time.Time

		tempTestCasesInTestSuiteAsByteArray          []byte
		tempTestSuitePreviewAsByteArray              []byte
		tempTestSuiteMetaDataAsByteArray             []byte
		tempTestSuiteTestDataAsByteArray             []byte
		tempTestSuiteImplementedFunctionsAsByteArray []byte

		tempTestCasesInTestSuiteAsGrpc          fenixTestCaseBuilderServerGrpcApi.TestCasesInTestSuiteMessage
		tempTestSuitePreviewAsGrpc              fenixTestCaseBuilderServerGrpcApi.TestSuitePreviewMessage
		tempTestSuiteMetaDataAsGrpc             fenixTestCaseBuilderServerGrpcApi.UserSpecifiedTestSuiteMetaDataMessage
		tempTestSuiteTestDataAsGrpc             fenixTestCaseBuilderServerGrpcApi.UsersChosenTestDataForTestSuiteMessage
		tempTestSuiteImplementedFunctionsAsGrpc map[int32]bool
	)

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&tempDomainUuid,
			&tempDomainName,
			&tempTestSuiteUuid,
			&tempTestSuiteName,

			&tempTestSuiteVersion,
			&tempTestSuiteDescription,
			&tempTestSuiteExecutionEnvironment,
			&tempTestSuiteHash,
			&tempDeleteTimestampAsTimeStamp,

			&tempUpdatedByAndWhenAsTimeStamp,
			&tempInsertedByUserIdOnComputer,
			&tempInsertedByGCPAuthenticatedUser,

			&tempTestCasesInTestSuiteAsJson,
			&tempTestSuitePreviewAsJson,
			&tempTestSuiteMetaDataAsJson,
			&tempTestSuiteTestDataAsJson,

			&tempTestSuiteType,
			&tempTestSuiteTypeName,
			&tempTestSuiteImplementedFunctionsAsJson,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "89c7db75-5ff1-4bf3-a1bb-5470520af682",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// Format Dates
		// Format The Delete Date into a string
		tempDeleteTimestampAsString = tempDeleteTimestampAsTimeStamp.Format("2006-01-02")

		// When the TestCase is not deleted then it uses a Delete date far away in the future. If so then clear the Date sent to TesterGui

		if tempDeleteTimestampAsString == testSuiteNotDeletedDate {
			tempDeleteTimestampAsString = ""
		}

		// Format Insert date
		tempUpdatedByAndWhenAsString = tempUpdatedByAndWhenAsTimeStamp.String()

		// Convert json-strings into byte-arrays
		tempTestCasesInTestSuiteAsByteArray = []byte(tempTestCasesInTestSuiteAsJson)
		tempTestSuitePreviewAsByteArray = []byte(tempTestSuitePreviewAsJson)
		tempTestSuiteMetaDataAsByteArray = []byte(tempTestSuiteMetaDataAsJson)
		tempTestSuiteTestDataAsByteArray = []byte(tempTestSuiteTestDataAsJson)
		tempTestSuiteImplementedFunctionsAsByteArray = []byte(tempTestSuiteImplementedFunctionsAsJson)

		// Convert json-byte-arrays into proto-messages
		err = protojson.Unmarshal(tempTestCasesInTestSuiteAsByteArray, &tempTestCasesInTestSuiteAsGrpc)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "2e73987c-1693-4673-b366-0b56efcfbc09",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestCasesInTestSuiteAsByteArray' into proto-message")

			return nil, err
		}

		err = protojson.Unmarshal(tempTestSuitePreviewAsByteArray, &tempTestSuitePreviewAsGrpc)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "eafb06bb-1257-4c49-b57c-99edb3baa81d",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestSuitePreviewAsByteArray' into proto-message")

			return nil, err
		}

		err = protojson.Unmarshal(tempTestSuiteMetaDataAsByteArray, &tempTestSuiteMetaDataAsGrpc)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "4643c3d5-87bf-4a5f-a171-b01a9d3c0d4b",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestSuiteMetaDataAsByteArray' into proto-message")

			return nil, err
		}

		err = protojson.Unmarshal(tempTestSuiteTestDataAsByteArray, &tempTestSuiteTestDataAsGrpc)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "ef74050a-d52b-4282-a520-32f9ca1ceecf",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestSuiteTestDataAsByteArray' into proto-message")

			return nil, err
		}

		tempTestSuiteImplementedFunctionsAsGrpc = make(map[int32]bool)
		err = json.Unmarshal(tempTestSuiteImplementedFunctionsAsByteArray, &tempTestSuiteImplementedFunctionsAsGrpc)
		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":    "2a737639-2fbc-470b-80bc-c09972c6768d",
				"Error": err,
			}).Error("Something went wrong when converting 'tempTestSuiteImplementedFunctionsAsByteArray' into proto-message")

			return nil, err
		}

		// Add the different parts into full TestSuite-message
		fullTestSuiteMessage = &fenixTestCaseBuilderServerGrpcApi.FullTestSuiteMessage{
			TestSuiteBasicInformation: &fenixTestCaseBuilderServerGrpcApi.TestSuiteBasicInformationMessage{
				DomainUuid:                    tempDomainUuid,
				DomainName:                    tempDomainName,
				TestSuiteUuid:                 tempTestSuiteUuid,
				TestSuiteVersion:              uint32(tempTestSuiteVersion),
				TestSuiteName:                 tempTestSuiteName,
				TestSuiteDescription:          tempTestSuiteDescription,
				TestSuiteExecutionEnvironment: tempTestSuiteExecutionEnvironment,
			},
			TestSuiteTestData:    &tempTestSuiteTestDataAsGrpc,
			TestSuitePreview:     &tempTestSuitePreviewAsGrpc,
			TestSuiteMetaData:    &tempTestSuiteMetaDataAsGrpc,
			TestCasesInTestSuite: &tempTestCasesInTestSuiteAsGrpc,
			DeletedDate:          tempDeleteTimestampAsString,
			UpdatedByAndWhen: &fenixTestCaseBuilderServerGrpcApi.UpdatedByAndWhenMessage{
				UserIdOnComputer:     tempInsertedByUserIdOnComputer,
				GCPAuthenticatedUser: tempInsertedByGCPAuthenticatedUser,
				UpdateTimeStamp:      tempUpdatedByAndWhenAsString,
			},
			TestSuiteType: &fenixTestCaseBuilderServerGrpcApi.TestSuiteTypeMessage{
				TestSuiteType:     fenixTestCaseBuilderServerGrpcApi.TestSuiteTypeEnum(tempTestSuiteType),
				TestSuiteTypeName: tempTestSuiteTypeName,
			},
			TestSuiteImplementedFunctionsMap: tempTestSuiteImplementedFunctionsAsGrpc,
			MessageHash:                      tempTestSuiteHash,
		}

	}

	return fullTestSuiteMessage, err

}

// Add who created the TestSuite and when it was created
func (fenixCloudDBObject *FenixCloudDBObjectStruct) addTestSuiteCreator(
	dbTransaction pgx.Tx,
	testSuiteUuidToLoad string,
	fullTestSuiteMessage *fenixTestCaseBuilderServerGrpcApi.FullTestSuiteMessage) (
	err error) {

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT \"InsertTimeStamp\", \"InsertedByUserIdOnComputer\", \"InsertedByGCPAuthenticatedUser\" "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestSuites\" "
	sqlToExecute = sqlToExecute + fmt.Sprintf("WHERE \"TestSuiteUuid\" = '%s' ", testSuiteUuidToLoad)
	sqlToExecute = sqlToExecute + "ORDER BY \"UniqueCounter\" ASC "
	sqlToExecute = sqlToExecute + "LIMIT 1 "
	sqlToExecute = sqlToExecute + "; "

	// Log SQL to be executed if Environment variable is true
	if common_config.LogAllSQLs == true {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "e195cf4d-d9fc-444b-bfe8-65c4a666bfd5",
			"sqlToExecute": sqlToExecute,
		}).Debug("SQL to be executed within 'addTestSuiteCreator'")
	}

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := dbTransaction.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "95382e00-59e0-4082-be16-fe56bee992d0",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	var (
		tempInsertTimeStampAsString    string
		InsertedByUserIdOnComputer     string
		InsertedByGCPAuthenticatedUser string

		tempInsertTimeStampAsTimeStamp time.Time

		numberOfRows int
	)

	// Extract data from DB result set
	for rows.Next() {

		err = rows.Scan(
			&tempInsertTimeStampAsTimeStamp,
			&InsertedByUserIdOnComputer,
			&InsertedByGCPAuthenticatedUser,
		)

		if err != nil {

			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "1f5d50c8-d35a-4d87-be33-80f7e64f7a93",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return err
		}

		// Format Dates
		// Format Insert date
		tempInsertTimeStampAsString = tempInsertTimeStampAsTimeStamp.String()

		// increase number of rows
		numberOfRows = numberOfRows + 1
	}

	if numberOfRows != 1 {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "d1987560-3ea1-4081-baba-c9e3fc699553",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
			"numberOfRows": numberOfRows,
		}).Error("Number of rows expected to be one, but was not.")

		err = errors.New("number of rows expected to be one, but was not")

		return err

	}

	// Update Created information
	fullTestSuiteMessage.UpdatedByAndWhen.CreatedByComputerLogin = InsertedByUserIdOnComputer
	fullTestSuiteMessage.UpdatedByAndWhen.CreatedByGcpLogin = InsertedByGCPAuthenticatedUser
	fullTestSuiteMessage.UpdatedByAndWhen.CreatedDate = tempInsertTimeStampAsString

	return err

}
