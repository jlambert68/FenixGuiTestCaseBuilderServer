package main

import (
	"context"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"time"
)

// ****************************************************************************************************************
// Load data from CloudDB into memory structures, to speed up stuff
//
// All TestDataRowItems in CloudDB
func (fenixGuiTestCaseBuilderServerObject *fenixGuiTestCaseBuilderServerObject_struct) loadAllTestDataRowItemsForClientFromCloudDB(testDataRowItems *[]cloudDBExposedTestDataRowItemsStruct) (err error) {

	fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
		"Id": "61b8b021-9568-463e-b867-ac1ddb10584d",
	}).Debug("Entering: loadAllTestDataRowItemsForClientFromCloudDB()")

	defer func() {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id": "78a97c41-a098-4122-88d2-01ed4b6c4844",
		}).Debug("Exiting: loadAllTestDataRowItemsForClientFromCloudDB()")
	}()

	/* Example
	SELECT *
	FROM "FenixTestDataSyncClient"."CurrentExposedTestDataForClient"

	    row_hash                  varchar   not null,
	    testdata_value_as_string  varchar   not null,
	    value_column_order        integer   not null,
	    value_row_order           integer   not null,
	    updated_timestamp         timestamp not null,
	    merkletree_leaf_node_name varchar   not null,

	*/

	usedDBSchema := fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT * "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"CurrentExposedTestDataForClient\";"

	// Query DB
	rows, err := fenixSyncShared.DbPool.Query(context.Background(), sqlToExecute)

	if err != nil {
		fenixGuiTestCaseBuilderServerObject.logger.WithFields(logrus.Fields{
			"Id":           "2f130d7e-f8aa-466f-b29d-0fb63608c1a6",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when executing SQL")

		return err
	}

	// Variables to used when extract data from result set
	var cloudDBExposedTestDataRowItem cloudDBExposedTestDataRowItemsStruct
	var tempTimeStamp time.Time
	var timeStampAsString string
	timeStampLayOut := fenixSyncShared.TimeStampLayOut //"2006-01-02 15:04:05.000000" //milliseconds

	// Extract data from DB result set
	for rows.Next() {
		err := rows.Scan(
			&cloudDBExposedTestDataRowItem.rowHash,
			&cloudDBExposedTestDataRowItem.testdataValueAsString,
			&cloudDBExposedTestDataRowItem.valueColumnOrder,
			&cloudDBExposedTestDataRowItem.valueRowOrder,
			&tempTimeStamp,
			&cloudDBExposedTestDataRowItem.merkleTreeLeafNodeName,
		)

		if err != nil {
			return err
		}

		// Convert timestamp into string representation and add to  extracted data
		timeStampAsString = tempTimeStamp.Format(timeStampLayOut)
		cloudDBExposedTestDataRowItem.updatedTimeStamp = timeStampAsString

		// Add values to the object that is pointed to by variable in function
		*testDataRowItems = append(*testDataRowItems, cloudDBExposedTestDataRowItem)

	}

	// No errors occurred
	return nil

}
