protoc -I ../internal/transport/grpc/proto/ \
 trade.proto \
 --go-grpc_out=../internal/transport --go_out=../internal/transport