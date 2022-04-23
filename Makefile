
# cat -e -t -v Makefile

#.DEFAULT_GOAL := buildAndRun
# Migration steps for new Golang gRPC installation
# https://stackoverflow.com/questions/60578892/protoc-gen-go-grpc-program-not-found-or-is-not-executable


compileProto:
	@echo "Compile proto file..."
# Old
#	protoc  -I ../grpc_api/fenixTestDataSyncServerGrpcApi --go_out=plugins=grpc:../grpc_api/fenixTestDataSyncServerGrpcApi ../grpc_api/fenixTestDataSyncServerGrpcApi/fenixTestDataSyncServerGrpcApi.proto
# New
 # generate the messages
# protoc --go_out="$GO_GEN_PATH" -I "$dependecies" "$proto"
	cd ../grpc_api/fenixClientTestDataSyncServerGrpcApi && protoc --go_out=. fenixClientTestDataSyncServerGrpcApi.proto

# generate the services
# protoc --go-grpc_out="$GO_GEN_PATH" -I "$dependecies" "$proto"
	cd ../grpc_api/fenixClientTestDataSyncServerGrpcApi && protoc --go-grpc_out=. fenixClientTestDataSyncServerGrpcApi.proto

#	sed -i 's/import "common_grpc_api.proto";/\/\/import "common_grpc_api.proto";/' ..\/grpc_api\/fenixClientTestDataSyncServerGrpcApi\/fenixClientTestDataSyncServerGrpcApi.proto
#	sed -i 's/\/\/import "..\/common_grpc_api\/common_grpc_api.proto";/import "..\/common_grpc_api\/common_grpc_api.proto";/' ..\/grpc_api\/fenixClientTestDataSyncServerGrpcApi\/fenixClientTestDataSyncServerGrpcApi.proto

#protoc version: 3.17.3
#protoc-gen-go: 1.26 (latest)
#protoc-gen-go-grpc: 1.1 (latest)
#golang: 1.17