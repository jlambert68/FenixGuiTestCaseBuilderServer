package common_config

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	fenixTestCaseBuilderServerGrpcApi "github.com/jlambert68/FenixGrpcApi/FenixTestCaseBuilderServer/fenixTestCaseBuilderServerGrpcApi/go_grpc_api"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

// Hash a single value
func HashSingleValue(valueToHash string) (hashValue string) {

	hash := sha256.New()
	hash.Write([]byte(valueToHash))
	hashValue = hex.EncodeToString(hash.Sum(nil))

	return hashValue

}

// GenerateDatetimeTimeStampForDB
// Generate DataBaseTimeStamp, eg '2022-02-08 17:35:04.000000'
func GenerateDatetimeTimeStampForDB() (currentTimeStampAsString string) {

	timeStampLayOut := "2006-01-02 15:04:05.000000 -0700" //milliseconds
	currentTimeStamp := time.Now().UTC()
	currentTimeStampAsString = currentTimeStamp.Format(timeStampLayOut)

	return currentTimeStampAsString
}

// GenerateDatetimeFromTimeInputForDB
// Generate DataBaseTimeStamp, eg '2022-02-08 17:35:04.000000'
func GenerateDatetimeFromTimeInputForDB(currentTime time.Time) (currentTimeStampAsString string) {

	timeStampLayOut := "2006-01-02 15:04:05.000000 -0700" //milliseconds
	currentTimeStampAsString = currentTime.Format(timeStampLayOut)

	return currentTimeStampAsString
}

// ********************************************************************************************************************
// Check if Calling Client is using correct proto-file version
func IsClientUsingCorrectTestDataProtoFileVersion(
	callingClientUuid string,
	usedProtoFileVersion fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum) (
	returnMessage *fenixTestCaseBuilderServerGrpcApi.AckNackResponse) {

	var clientUseCorrectProtoFileVersion bool
	var protoFileExpected fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum
	var protoFileUsed fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum

	protoFileUsed = usedProtoFileVersion
	protoFileExpected = fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum(GetHighestFenixGuiBuilderProtoFileVersion())

	// Check if correct proto files is used
	if protoFileExpected == protoFileUsed {
		clientUseCorrectProtoFileVersion = true
	} else {
		clientUseCorrectProtoFileVersion = false
	}

	// Check if Client is using correct proto files version
	if clientUseCorrectProtoFileVersion == false {
		// Not correct proto-file version is used

		// Set Error codes to return message
		var errorCodes []fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum
		var errorCode fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum

		errorCode = fenixTestCaseBuilderServerGrpcApi.ErrorCodesEnum_ERROR_WRONG_PROTO_FILE_VERSION
		errorCodes = append(errorCodes, errorCode)

		// Create Return message
		returnMessage = &fenixTestCaseBuilderServerGrpcApi.AckNackResponse{
			AckNack:                      false,
			Comments:                     "Wrong proto file used. Expected: '" + protoFileExpected.String() + "', but got: '" + protoFileUsed.String() + "'",
			ErrorCodes:                   errorCodes,
			ProtoFileVersionUsedByClient: protoFileExpected,
		}

		Logger.WithFields(logrus.Fields{
			"id": "513dd8fb-a0bb-4738-9a0b-b7eaf7bb8adb",
		}).Debug("Wrong proto file used. Expected: '" + protoFileExpected.String() + "', but got: '" + protoFileUsed.String() + "' for Client: " + callingClientUuid)

		return returnMessage

	} else {
		return nil
	}

}

// ********************************************************************************************************************
// Get the highest FenixProtoFileVersionEnumeration
func GetHighestFenixGuiBuilderProtoFileVersion() int32 {

	// Check if there already is a 'highestFenixProtoFileVersion' saved, if so use that one
	if highestFenixProtoFileVersion != -1 {
		return highestFenixProtoFileVersion
	}

	// Find the highest value for proto-file version
	var maxValue int32
	maxValue = 0

	for _, v := range fenixTestCaseBuilderServerGrpcApi.CurrentFenixTestCaseBuilderProtoFileVersionEnum_value {
		if v > maxValue {
			maxValue = v
		}
	}

	highestFenixProtoFileVersion = maxValue

	return highestFenixProtoFileVersion
}

// Extracts 'ParserLayout' from the TimeStamp(as string)
func GenerateTimeStampParserLayout(timeStampAsString string) (parserLayout string, err error) {
	// "2006-01-02 15:04:05.999999999 -0700 MST"

	var timeStampParts []string
	var timeParts []string
	var numberOfDecimals int

	// Split TimeStamp into separate parts
	timeStampParts = strings.Split(timeStampAsString, " ")

	// Validate that first part is a date with the following form '2006-01-02'
	if len(timeStampParts[0]) != 10 {

		Logger.WithFields(logrus.Fields{
			"Id":                "ffbf0682-ebc7-4e27-8ad1-0e5005fbc364",
			"timeStampAsString": timeStampAsString,
			"timeStampParts[0]": timeStampParts[0],
		}).Error("Date part has not the correct form, '2006-01-02'")

		err = errors.New(fmt.Sprintf("Date part, '%s' has not the correct form, '2006-01-02'", timeStampParts[0]))

		return "", err

	}

	// Add Date to Parser Layout
	parserLayout = "2006-01-02"

	// Add Time to Parser Layout
	parserLayout = parserLayout + " 15:04:05."

	// Split time into time and decimals
	timeParts = strings.Split(timeStampParts[1], ".")

	// Get number of decimals
	if len(timeParts) > 1 {
		numberOfDecimals = len(timeParts[1])

		// Add Decimals to Parser Layout
		parserLayout = parserLayout + strings.Repeat("9", numberOfDecimals)

	} else {
		numberOfDecimals = 0

		// remove added '.' decimal divider
		parserLayout = parserLayout[0 : len(parserLayout)-1]
	}

	// Add time zone, part 1, if that information exists
	if len(timeStampParts) > 2 {
		parserLayout = parserLayout + " -0700"
	}

	// Add time zone, part 2, if that information exists
	if len(timeStampParts) > 3 {
		parserLayout = parserLayout + " MST"
	}

	return parserLayout, err
}
