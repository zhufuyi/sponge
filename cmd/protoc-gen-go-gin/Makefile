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
	  --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugin=handler \
	  --go-gin_opt=moduleName=yourModuleName --go-gin_opt=serverName=yourServerName --go-gin_opt=out=internal/handler \
	  --plugin=./protoc-gen-go-gin* \
	  api/v1/*.proto

service:
	@go build
	protoc --proto_path=. --proto_path=./third_party \
	  --go_out=. --go_opt=paths=source_relative \
	  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	  --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugin=service \
	  --go-gin_opt=moduleName=yourModuleName --go-gin_opt=serverName=yourServerName --go-gin_opt=out=internal/service \
	  --plugin=./protoc-gen-go-gin* \
	  api/v1/*.proto

clean:
	@rm -vrf api/v1/*.go
	@rm -vrf protoc-gen-go-gin*
	@rm -vrf internal/service
	@rm -vrf internal/handler
	@echo Done.
