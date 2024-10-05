package common_config

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// GenerateDatetimeTimeStampForDB
// Generate DataBaseTimeStamp, eg '2022-02-08 17:35:04.000000'
func GenerateDatetimeTimeStampForDB() (currentTimeStampAsString string) {

	timeStampLayOut := "2006-01-02 15:04:05.000000" //milliseconds
	currentTimeStamp := time.Now()
	currentTimeStampAsString = currentTimeStamp.Format(timeStampLayOut)

	return currentTimeStampAsString
}

// GenerateDatetimeFromTimeInputForDB
// Generate DataBaseTimeStamp, eg '2022-02-08 17:35:04.000000'
func GenerateDatetimeFromTimeInputForDB(currentTime time.Time) (currentTimeStampAsString string) {

	timeStampLayOut := "2006-01-02 15:04:05.000000" //milliseconds
	currentTimeStampAsString = currentTime.Format(timeStampLayOut)

	return currentTimeStampAsString
}

// ConvertGrpcTimeStampToStringForDB
// Convert a gRPC-timestamp into a string that can be used to store in the database
func ConvertGrpcTimeStampToStringForDB(grpcTimeStamp *timestamppb.Timestamp) (grpcTimeStampAsTimeStampAsString string) {
	grpcTimeStampAsTimeStamp := grpcTimeStamp.AsTime()

	timeStampLayOut := "2006-01-02 15:04:05.000000" //milliseconds

	grpcTimeStampAsTimeStampAsString = grpcTimeStampAsTimeStamp.Format(timeStampLayOut)

	return grpcTimeStampAsTimeStampAsString
}
