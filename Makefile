generate_grpc_code:
	protoc \
	--go_out=proto \
	--go_opt=paths=source_relative \
	--go-grpc_out=proto \
	--go-grpc_opt=paths=source_relative \
	playlist.proto