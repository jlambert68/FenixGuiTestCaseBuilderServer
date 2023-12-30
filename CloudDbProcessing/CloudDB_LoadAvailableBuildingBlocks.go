package CloudDbProcessing

import (
	"FenixGuiTestCaseBuilderServer/common_config"
	"context"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	fenixSyncShared "github.com/jlambert68/FenixSyncShared"
	"github.com/sirupsen/logrus"
	"time"
)

// Load TestInstructions for Client
func (fenixCloudDBObject *FenixCloudDBObjectStruct) LoadClientsImmatureTestInstructionsFromCloudDB(
	gCPAuthenticatedUser string) (
	cloudDBImmatureTestInstructionItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage, err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "273dceef-7982-4e7d-98db-c132342e530b",
	}).Debug("Entering: loadClientsImmatureTestInstructionsFromCloudDB()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "0894a272-3e91-407a-b5a4-1b70f8e00e6b",
		}).Debug("Exiting: loadClientsImmatureTestInstructionsFromCloudDB()")
	}()

	immatureTestInstructionMessageMap := make(map[string]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage)

	err = fenixCloudDBObject.processTestInstructionsBasicTestInstructionInformation(immatureTestInstructionMessageMap)
	if err != nil {
		return nil, err
	}

	err = fenixCloudDBObject.processTestInstructionsImmatureTestInstructionInformation(immatureTestInstructionMessageMap)
	if err != nil {
		return nil, err
	}

	err = fenixCloudDBObject.processTestInstructionsImmatureElementModel(immatureTestInstructionMessageMap)
	if err != nil {
		return nil, err
	}

	// Loop all ImmatureTestInstructionMessage and create gRPC-response
	var allImmatureTestInstructionMessage []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionMessage

	for _, value := range immatureTestInstructionMessageMap { // Order not specified
		allImmatureTestInstructionMessage = append(allImmatureTestInstructionMessage, value)
	}

	cloudDBImmatureTestInstructionItems = allImmatureTestInstructionMessage

	// No errors occurred
	return cloudDBImmatureTestInstructionItems, nil

}

// Load TestInstructionContainers for Client
func (fenixCloudDBObject *FenixCloudDBObjectStruct) LoadClientsImmatureTestInstructionContainersFromCloudDB(userID string) (cloudDBImmatureTestInstructionContainerItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage, err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "68b965ea-234c-425b-b525-1f8b7154850b",
	}).Debug("Entering: loadClientsImmatureTestInstructionContainersFromCloudDB()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "12021bfa-154f-48f2-bd8c-0809e1877fd4",
		}).Debug("Exiting: loadClientsImmatureTestInstructionContainersFromCloudDB()")
	}()

	immatureTestInstructionContainerMessageMap := make(map[string]*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage)

	err = fenixCloudDBObject.processTestInstructionContainersBasicTestInstructionContainerInformation(immatureTestInstructionContainerMessageMap)
	if err != nil {
		return nil, err
	}

	err = fenixCloudDBObject.processTestInstructionContainersImmatureTestInstructionContainerInformation(immatureTestInstructionContainerMessageMap)
	if err != nil {
		return nil, err
	}

	err = fenixCloudDBObject.processTestInstructionContainersImmatureElementModel(immatureTestInstructionContainerMessageMap)
	if err != nil {
		return nil, err
	}

	// Loop all ImmatureTestInstructionContainerMessage and create gRPC-response
	var allImmatureTestInstructionContainerMessage []*fenixTestCaseBuilderServerGrpcApi.ImmatureTestInstructionContainerMessage

	for _, value := range immatureTestInstructionContainerMessageMap { // Order not specified
		allImmatureTestInstructionContainerMessage = append(allImmatureTestInstructionContainerMessage, value)
	}

	cloudDBImmatureTestInstructionContainerItems = allImmatureTestInstructionContainerMessage

	// No errors occurred
	return cloudDBImmatureTestInstructionContainerItems, nil

}

// Load Pinned TestInstructions for Client
func (fenixCloudDBObject *FenixCloudDBObjectStruct) LoadClientsPinnedTestInstructionsFromCloudDB(userID string) (availablePinnedPreCreatedTestInstructionContainerMessage []*fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionMessage, err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "9901525f-a271-4f4f-a798-fea7fdf29dfb",
	}).Debug("Entering: loadClientsPinnedTestInstructionsFromCloudDB()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "0f7be73c-4065-4d4a-ae02-40a4d93fc2a3",
		}).Debug("Exiting: loadClientsPinnedTestInstructionsFromCloudDB()")
	}()

	/*
		SELECT PTITIC.*
		FROM "FenixBuilder"."PinnedTestInstructionsAndPreCreatedTestInstructionContainers" PTITIC
		WHERE PTITIC."PinnedType" = 1;
	*/

	usedDBSchema := "FenixBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT PTITIC.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"PinnedTestInstructionsAndPreCreatedTestInstructionContainers\" PTITIC "
	sqlToExecute = sqlToExecute + "WHERE PTITIC.\"PinnedType\" = 1 AND " // 1 = TestInstructions
	sqlToExecute = sqlToExecute + "PTITIC.\"UserId\" = '" + userID + "' "
	sqlToExecute = sqlToExecute + "ORDER BY PTITIC.\"PinnedUuid\" ASC; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "6be84b22-613f-4d93-afe8-e8ee22826e7b",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return nil, err
	}

	var tempUserId string
	var tempPinnedType int
	var tempTimeStamp time.Time
	var availablePinnedTestInstructionMessages []fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionMessage

	// Extract data from DB result set
	for rows.Next() {

		// Initiate new fresh variable
		newAvailablePinnedTestInstructionMessage := fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionMessage{}

		err = rows.Scan(

			&tempUserId,
			&newAvailablePinnedTestInstructionMessage.TestInstructionUuid,
			&newAvailablePinnedTestInstructionMessage.TestInstructionName,
			&tempPinnedType,
			&tempTimeStamp,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "e1d695b7-ec8a-4692-9e9a-416869923e82",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// append PinnedTestInstruction to array
		availablePinnedTestInstructionMessages = append(availablePinnedTestInstructionMessages, newAvailablePinnedTestInstructionMessage)

	}

	// Convert to pointer-array that fits gRPC api
	var availablePinnedTestInstructionToSendOvergRPC []*fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionMessage
	for _, tempAvailablePinnedTestInstructionMessage := range availablePinnedTestInstructionMessages {
		newAvailablePinnedTestInstructionMessage := fenixTestCaseBuilderServerGrpcApi.AvailablePinnedTestInstructionMessage{}
		newAvailablePinnedTestInstructionMessage = tempAvailablePinnedTestInstructionMessage
		availablePinnedTestInstructionToSendOvergRPC = append(availablePinnedTestInstructionToSendOvergRPC, &newAvailablePinnedTestInstructionMessage)
	}

	return availablePinnedTestInstructionToSendOvergRPC, err
}

// Load Pinned TestInstructionContainers for Client
func (fenixCloudDBObject *FenixCloudDBObjectStruct) LoadClientsPinnedTestInstructionContainersFromCloudDB(userID string) (availablePinnedPreCreatedTestInstructionContainerContainerMessage []*fenixTestCaseBuilderServerGrpcApi.AvailablePinnedPreCreatedTestInstructionContainerMessage, err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "c2decb25-9f53-44c0-be49-88ac5c9cde5d",
	}).Debug("Entering: loadClientsPinnedTestInstructionContainersFromCloudDB()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "a9863d2a-4f59-4eef-a939-117bcddea3c4",
		}).Debug("Exiting: loadClientsPinnedTestInstructionContainersFromCloudDB()")
	}()

	usedDBSchema := "FenixBuilder" // TODO should this env variable be used? fenixSyncShared.GetDBSchemaName()

	sqlToExecute := ""
	sqlToExecute = sqlToExecute + "SELECT PTITIC.* "
	sqlToExecute = sqlToExecute + "FROM \"" + usedDBSchema + "\".\"PinnedTestInstructionsAndPreCreatedTestInstructionContainers\" PTITIC "
	sqlToExecute = sqlToExecute + "WHERE PTITIC.\"PinnedType\" = 2 AND " // 2 = TestInstructionContainers
	sqlToExecute = sqlToExecute + "PTITIC.\"UserId\" = '" + userID + "' "
	sqlToExecute = sqlToExecute + "ORDER BY PTITIC.\"PinnedUuid\" ASC; "

	// Query DB
	var ctx context.Context
	ctx, timeOutCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer timeOutCancel()

	rows, err := fenixSyncShared.DbPool.Query(ctx, sqlToExecute)
	defer rows.Close()

	if err != nil {
		common_config.Logger.WithFields(logrus.Fields{
			"Id":           "6be84b22-613f-4d93-afe8-e8ee22826e7b",
			"Error":        err,
			"sqlToExecute": sqlToExecute,
		}).Error("Something went wrong when processing result from database")

		return nil, err
	}

	var tempUserId string
	var tempPinnedType int
	var tempTimeStamp time.Time
	var availablePinnedTestInstructionContainerMessages []fenixTestCaseBuilderServerGrpcApi.AvailablePinnedPreCreatedTestInstructionContainerMessage

	// Extract data from DB result set
	for rows.Next() {

		// Initiate new fresh variable
		newAvailablePinnedTestInstructionContainerMessage := fenixTestCaseBuilderServerGrpcApi.AvailablePinnedPreCreatedTestInstructionContainerMessage{}

		err = rows.Scan(

			&tempUserId,
			&newAvailablePinnedTestInstructionContainerMessage.TestInstructionContainerUuid,
			&newAvailablePinnedTestInstructionContainerMessage.TestInstructionContainerName,
			&tempPinnedType,
			&tempTimeStamp,
		)

		if err != nil {
			common_config.Logger.WithFields(logrus.Fields{
				"Id":           "e1d695b7-ec8a-4692-9e9a-416869923e82",
				"Error":        err,
				"sqlToExecute": sqlToExecute,
			}).Error("Something went wrong when processing result from database")

			return nil, err
		}

		// append PinnedTestInstructionContainer to array
		availablePinnedTestInstructionContainerMessages = append(availablePinnedTestInstructionContainerMessages, newAvailablePinnedTestInstructionContainerMessage)

	}

	// Convert to pointer-array that fits gRPC api
	var availablePinnedTestInstructionContainerToSendOvergRPC []*fenixTestCaseBuilderServerGrpcApi.AvailablePinnedPreCreatedTestInstructionContainerMessage
	for _, tempAvailablePinnedTestInstructionContainerMessage := range availablePinnedTestInstructionContainerMessages {
		newAvailablePinnedTestInstructionContainerMessage := fenixTestCaseBuilderServerGrpcApi.AvailablePinnedPreCreatedTestInstructionContainerMessage{}
		newAvailablePinnedTestInstructionContainerMessage = tempAvailablePinnedTestInstructionContainerMessage
		availablePinnedTestInstructionContainerToSendOvergRPC = append(availablePinnedTestInstructionContainerToSendOvergRPC, &newAvailablePinnedTestInstructionContainerMessage)
	}

	return availablePinnedTestInstructionContainerToSendOvergRPC, err
}

// Load TestInstructions for Client
func (fenixCloudDBObject *FenixCloudDBObjectStruct) LoadAvailableBondsFromCloudDB() (cloudDBAvailableBondsItems []*fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage_ImmatureBondMessage, err error) {

	common_config.Logger.WithFields(logrus.Fields{
		"Id": "4b7058fe-c46d-4ab8-8612-895c8e1102a1",
	}).Debug("Entering: loadAvailableBondsFromCloudDB()")

	defer func() {
		common_config.Logger.WithFields(logrus.Fields{
			"Id": "665cecf6-a2cc-4c3a-bcb7-9ac1170bd8d3",
		}).Debug("Exiting: loadAvailableBondsFromCloudDB()")
	}()

	//availableBondsAttributes := []fenixTestCaseBuilderServerGrpcApi.BasicBondInformationMessage_VisibleBondAttributesMessage

	availableBondsAttributes, err := fenixCloudDBObject.processVisibleBondAttributesInformation()
	if err != nil {
		return nil, err
	}

	// Loop all Bonds-messages and create gRPC-response
	for _, visibleBondAttributesMessage := range availableBondsAttributes {

		// Deep copy of values
		tempBondAttributesMessage := fenixTestCaseBuilderServerGrpcApi.BasicBondInformationMessage_VisibleBondAttributesMessage{
			BondUuid:                 visibleBondAttributesMessage.BondUuid,
			BondName:                 visibleBondAttributesMessage.BondName,
			BondDescription:          visibleBondAttributesMessage.BondDescription,
			BondMouseOverText:        visibleBondAttributesMessage.BondMouseOverText,
			Deprecated:               visibleBondAttributesMessage.Deprecated,
			Enabled:                  visibleBondAttributesMessage.Enabled,
			Visible:                  visibleBondAttributesMessage.Visible,
			BondColor:                visibleBondAttributesMessage.BondColor,
			CanBeDeleted:             visibleBondAttributesMessage.CanBeDeleted,
			CanBeSwappedOut:          visibleBondAttributesMessage.CanBeSwappedOut,
			UpdatedTimeStamp:         visibleBondAttributesMessage.UpdatedTimeStamp,
			TestCaseModelElementType: visibleBondAttributesMessage.TestCaseModelElementType,
			ShowBondAttributes:       visibleBondAttributesMessage.ShowBondAttributes,
			TCRuleDeletion:           visibleBondAttributesMessage.TCRuleDeletion,
			TCRuleSwap:               visibleBondAttributesMessage.TCRuleSwap,
		}

		basicBondInformationMessage := fenixTestCaseBuilderServerGrpcApi.BasicBondInformationMessage{
			VisibleBondAttributes: &tempBondAttributesMessage}

		immatureBondsMessage_ImmatureBondMessage_B0_BOND := fenixTestCaseBuilderServerGrpcApi.ImmatureBondsMessage_ImmatureBondMessage{
			BasicBondInformation: &basicBondInformationMessage}

		cloudDBAvailableBondsItems = append(cloudDBAvailableBondsItems, &immatureBondsMessage_ImmatureBondMessage_B0_BOND)
	}

	// No errors occurred
	return cloudDBAvailableBondsItems, nil

}
