router:
	@go build
	protoc --proto_path=. --proto_path=./third_party \
	  --go_out=. --go_opt=paths=source_relative \
	  --go-gin_out=. --go-gin_opt=paths=source_relative  \
	  --plugin=./protoc-gen-go-gin* \
	  api/v1/*.proto

handler:
	@go build
	protoc --proto_path=. --proto_path=./third_party \
	  --go_out=. --go_opt=paths=source_relative \
	  --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugins=handler \
	  --plugin=./protoc-gen-go-gin* \
	  api/v1/*.proto

service:
	@go build
	protoc --proto_path=. --proto_path=./third_party \
	  --go_out=. --go_opt=paths=source_relative \
	  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	  --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugins=service \
	  --plugin=./protoc-gen-go-gin* \
	  api/v1/*.proto

clean:
	@rm -vrf api/v1/*.go
	@rm -vrf protoc-gen-go-gin*
	@echo Done.
