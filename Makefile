grpc_gen:
	protoc --go_out=./ping/proto --go_opt=paths=source_relative \
		--go-grpc_out=./ping/proto --go-grpc_opt=paths=source_relative \
		--proto_path=./ping/proto hello_ping.proto

	protoc --go_out=./pong/proto --go_opt=paths=source_relative \
		--go-grpc_out=./pong/proto --go-grpc_opt=paths=source_relative \
		--proto_path=./pong/proto hello_pong.proto