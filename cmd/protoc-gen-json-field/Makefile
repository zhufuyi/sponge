json:
	@go build
	protoc --proto_path=. --proto_path=./third_party \
	  --json-field_out=. --json-field_opt=paths=source_relative \
	  --plugin=./protoc-gen-json-field* \
	  api/v1/*.proto

clean:
	@rm -vrf protoc-gen-json-field*
	@rm -vrf api/v1/*.json
