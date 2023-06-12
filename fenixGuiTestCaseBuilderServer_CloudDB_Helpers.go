package main

import (
	"context"
	"errors"
	"fmt"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
)

// Load BasicInformation for TestCase to be able to populate the TestCaseExecution
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObjectStruct) getNexTestCaseVersion(testCaseUuid string) (nextTestCaseVersion uint32, err error) {

	//usedDBSchema := "FenixBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT MAX(TC.\"TestCaseVersion\") "
	sqlToExecute = sqlToExecute + "FROM \"FenixBuilder\".\"TestCases\" TC "
	sqlToExecute = sqlToExecute + "WHERE TC.\"TestCaseUuid\" = '" + testCaseUuid + "';"

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "78b049e0-ab58-4c5e-a7f8-1aa1416d8535",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

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
				fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
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
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":             "492e218c-9b5d-4b69-aad7-dff9137b4170",
			"Number of Rows": rows.CommandTag().RowsAffected(),
		}).Error("Expected 0 or 1 row")

		err := errors.New(fmt.Sprintf("Expected 0 or 1 row, but got %s rows ", rows.CommandTag().RowsAffected()))

		return 0, err
	}

	return currentTestCaseVersion + 1, err

}
