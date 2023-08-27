SRC_DIR=api/proto
DST_DIR=api/go/strategopb

protoc \
  -I=$SRC_DIR \
  --go_out=$DST_DIR \
  --go_opt=paths=source_relative \
  --go-grpc_out=$DST_DIR \
  --go-grpc_opt=paths=source_relative \
  $(find $SRC_DIR -type f -name '*.proto')
