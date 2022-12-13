// Package sponge is a microservice framework, a tool for quickly creating microservice code.
// sponge has a rich generating code commands, a total of 12 different functional code,
// these functional code can be combined into a complete service (similar to artificially broken
// sponge cells can be automatically reorganized into a new sponge ). Microservice code features
// include logging, service registration and discovery, registry, rate limit, circuit breaker, trace,
// metrics monitoring, pprof performance analysis, statistics, caching, CICD. The code uses a decoupled
// layered structure and it's easy to add or replace functional code. As an efficiency-enhancing tool, commonly
// repeated code is basically generated automatically and only business logic code needs to be populated
// based on the generated template code examples.
//
// https://raw.githubusercontent.com/zhufuyi/sponge/main/assets/README.md
//
//	sponge -h
//	sponge a microservice framework, a tool for quickly creating microservice code.
//
//	Usage:
//	sponge [command]
//
//	Available Commands:
//	completion  Generate the autocompletion script for the specified shell
//	config         Generate go config codes from yaml file
//	help           Help about any command
//	init            Initialize sponge
//	micro        Generate proto, model, dao, service, rpc, rpc-gw, rpc-cli codes
//	tools         Managing sponge dependency tools
//	update      Update sponge to the latest version
//	web          Generate model, dao, handler, http codes
package sponge
