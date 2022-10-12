// Package sponge is a go microservices framework, a tool for quickly creating complete microservices code for http or grpc.
// Generate `config`, `ecode`, `model`, `dao`, `handler`, `router`, `http`, `proto`, `service`, `grpc` code from the SQL DDL,
// which can be combined into full services(similar to how a broken sponge cell automatically reorganises itself into a new sponge).
//
// combined with the [sponge](https://github.com/zhufuyi/sponge@sponge) tool to generate framework code。
//
//	sponge -h
//	sponge management tools
//
//	Usage:
//	sponge [command]
//
//	Available Commands:
//	completion  Generate the autocompletion script for the specified shell
//	config      Generate go config code
//	dao         Generate dao code
//	grpc        Generate grpc server code
//	handler     Generate handler code
//	help        Help about any command
//	http        Generate http code
//	model       Generate model code
//	proto       Generate protobuf code
//	service     Generate grpc service code
package sponge
