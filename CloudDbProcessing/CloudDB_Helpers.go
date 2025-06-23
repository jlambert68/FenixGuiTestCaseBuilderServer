package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	"errors"
	"fmt"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"time"
)

// Load BasicInformation for TestCase to be able to populate the TestCaseExecution
func (fenixCloudDBObject *FenixCloudDBObjectStruct) getNexTestCaseVersion(testCaseUuid string) (nextTestCaseVersion uint32, err error) {

	//usedDBSchema := "FenixBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT MAX(TC.\"TestCaseVersion\") "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestCases\" TC "
	sqlToExecute = sqlToExecute + "WHERE TC.\"TestCaseUuid\" = '" + testCaseUuid + "';"

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "78b049e0-ab58-4c5e-a7f8-1aa1416d8535",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		databaseStatistics := fenixSyncShared.DbPool.Stat()

		common_config.Logger.WithFields(logrus.Fields{
			"Id":                 "92dce2b7-a9c1-40f2-b46e-618c91c83d65",
			"databaseStatistics": databaseStatistics,
		}).Error("DBPool Statistics")

		return 0, err
	}

	// Initiate a new variable to store the data
	var currentTestCaseVersion uint32
	var numberOfRows int64

	// Extract data from DB result set
	for rows.Next() {

		// Only process if there are row in response
		if rows.RawValues()[0] != nil {

			numberOfRows = int64(len(rows.RawValues()))

			err := rows.Scan(
				&currentTestCaseVersion,
			)

			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":           "0d94d43c-1989-4915-a22b-3d1f7ae4fe37",
					"Error":        err,
					"sqlToExecute": sqlToExecute,
				}).Error("Something went wrong when processing result from database")

				return 0, err
			}

		} else {
			numberOfRows = 0
			break

		}
	}

	// Verify that there are a maximum of 1 row
	switch numberOfRows {
	case 0:
		currentTestCaseVersion = 0

	case 1:

	default:
		common_config.Logger.WithFields(logrus.Fields{
			"Id":             "492e218c-9b5d-4b69-aad7-dff9137b4170",
			"Number of Rows": rows.CommandTag().RowsAffected(),
		}).Error("Expected 0 or 1 row")

		err := errors.New(fmt.Sprintf("Expected 0 or 1 row, but got %s rows ", rows.CommandTag().RowsAffected()))

		return 0, err
	}

	return currentTestCaseVersion + 1, err

}

// Get the Next version of the TestSuite
func (fenixCloudDBObject *FenixCloudDBObjectStruct) getNexTestSuiteVersion(
	testSuiteUuid string) (
	nextTestSuiteVersion uint32,
	err error) {

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT MAX(TS.\"TestSuiteVersion\") "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestSuites\" TS "
	sqlToExecute = sqlToExecute + "WHERE TS.\"TestSuitesUuid\" = '" + testSuiteUuid + "';"

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "f7d010a8-37fa-45ca-a52c-5931c3c68c39",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		databaseStatistics := fenixSyncShared.DbPool.Stat()

		common_config.Logger.WithFields(logrus.Fields{
			"Id":                 "988f38bd-01a1-4c50-a667-76215f3e8bcb",
			"databaseStatistics": databaseStatistics,
		}).Error("DBPool Statistics")

		return 0, err
	}

	// Initiate a new variable to store the data
	var currentTestSuiteVersion uint32
	var numberOfRows int64

	// Extract data from DB result set
	for rows.Next() {

		// Only process if there are row in response
		if rows.RawValues()[0] != nil {

			numberOfRows = int64(len(rows.RawValues()))

			err := rows.Scan(
				&currentTestSuiteVersion,
			)

			if err != nil {
				common_config.Logger.WithFields(logrus.Fields{
					"Id":           "0d94d43c-1989-4915-a22b-3d1f7ae4fe37",
					"Error":        err,
					"sqlToExecute": sqlToExecute,
				}).Error("Something went wrong when processing result from database")

				return 0, err
			}

		} else {
			numberOfRows = 0
			break

		}
	}

	// Verify that there are a maximum of 1 row
	switch numberOfRows {
	case 0:
		currentTestSuiteVersion = 0

	case 1:

	default:
		common_config.Logger.WithFields(logrus.Fields{
			"Id":             "d3677990-318d-41fd-9844-1ae6740c3d05",
			"Number of Rows": rows.CommandTag().RowsAffected(),
		}).Error("Expected 0 or 1 row")

		err := errors.New(fmt.Sprintf("Expected 0 or 1 row, but got %s rows ", rows.CommandTag().RowsAffected()))

		return 0, err
	}

	return currentTestSuiteVersion + 1, err

}

// GenerateDatetimeFromTimeInputForDB
// Generate DataBaseTimeStamp, eg '2022-02-08 17:35:04.000000'
func GenerateDatetimeFromTimeInputForDB(currentTime time.Time) (currentTimeStampAsString string) {

	timeStampLayOut := "2006-01-02 15:04:05.000000 -0700" //milliseconds
	currentTimeStampAsString = currentTime.Format(timeStampLayOut)

	return currentTimeStampAsString
}
