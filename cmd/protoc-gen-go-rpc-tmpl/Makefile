
proto:
	@go build
	protoc --proto_path=. --proto_path=./third_party \
	  --go_out=. --go_opt=paths=source_relative \
	  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	  --go-rpc-tmpl_out=. --go-rpc-tmpl_opt=paths=source_relative \
	  --go-rpc-tmpl_opt=moduleName=yourModuleName --go-rpc-tmpl_opt=serverName=yourServerName \
	  --plugin=./protoc-gen-go-rpc-tmpl* \
	  api/v1/*.proto

clean:
	@rm -vrf api/v1/*.go
	@rm -vrf protoc-gen-go-rpc-tmpl*
	@rm -vrf internal/service
	@rm -vrf internal/ecode
