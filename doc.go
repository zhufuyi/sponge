// Package sponge is a microservices framework for quickly creating http or grpc code.
// Generate codes `config`, `ecode`, `model`, `dao`, `handler`, `router`, `http`, `proto`, `service`, `grpc` from the SQL DDL,
// these codes can be combined into complete services (similar to how a broken sponge cell can automatically reorganize into a new sponge).
//
//	sponge -h
//	sponge management tools
//
//	Usage:
//	sponge [command]
//
//	Available Commands:
//	completion  Generate the autocompletion script for the specified shell
//	config      Generate go config code from yaml file
//	dao         Generate dao code
//	grpc        Generate grpc server code
//	handler     Generate handler code
//	help        Help about any command
//	http        Generate http code
//	model       Generate model code
//	proto       Generate protobuf code
//	service     Generate grpc service code
//	update		Update sponge to the latest version
package sponge
