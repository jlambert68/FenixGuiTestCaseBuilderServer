package main

// All TestDataRowItems in CloudDB
var cloudDBTestDataRowItems []cloudDBTestDataRowItemCurrentStruct

type cloudDBTestDataRowItemCurrentStruct struct {
	clientUuid            string
	rowHash               string
	testdataValueAsString string
	leafNodeName          string
	leafNodePath          string
	leafNodeHash          string
	valueColumnOrder      int
	valueRowOrder         int
	updatedTimeStamp      string
}

// Client Info in CloudDB
var cloudDBClientInfo cloudDBClientInfoStruct

type cloudDBClientInfoStruct struct {
	clientUuid                    string
	merklehash                    string
	labels_hash                   string
	meklehash_created_timestamp   string
	labels_hash_created_timestamp string
	meklehash_sent_timestamp      string
	labels_hash_sent_timestamp    string
	full_merkle_filter_path       string
	domain_name                   string
	domain_uuid                   string
}

// Exposed TestData for Client
var cloudDBExposedTestDataRowItems []cloudDBExposedTestDataRowItemsStruct

type cloudDBExposedTestDataRowItemsStruct struct {
	rowHash                string
	testdataValueAsString  string
	merkleTreeLeafNodeName string
	valueColumnOrder       int
	valueRowOrder          int
	updatedTimeStamp       string
}
