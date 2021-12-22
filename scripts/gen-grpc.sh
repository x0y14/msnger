protoc --go_out=../pkg/protobuf --go_opt=paths=source_relative \
  --go-grpc_out=../pkg/protobuf --go-grpc_opt=paths=source_relative \
  --go-grpc_opt=require_unimplemented_servers=true \
  -I ../api ../api/op.proto ../api/talk.proto ../api/types.proto
