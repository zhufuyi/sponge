package main

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/internal/model"
	"github.com/zhufuyi/sponge/internal/server"

	"github.com/zhufuyi/sponge/pkg/app"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/nacoscli"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/consul"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/etcd"
	"github.com/zhufuyi/sponge/pkg/tracer"

	"github.com/jinzhu/copier"
)

var (
	version            string
	configFile         string
	enableConfigCenter bool
)

// @title serverNameExample api docs
// @description http server api docs
// @schemes http https
// @version v0.0.0
// @host localhost:8080
func main() {
	inits := registerInits()
	servers := registerServers()
	closes := registerCloses(servers)

	s := app.New(inits, servers, closes)
	s.Run()
}

// -------------------------------- 注册app初始化 ---------------------------------
func registerInits() []app.Init {
	// 初始化配置
	initConfig()

	// 初始化日志
	_, _ = logger.Init(
		logger.WithLevel(config.Get().Logger.Level),
		logger.WithFormat(config.Get().Logger.Format),
		logger.WithSave(config.Get().Logger.IsSave),
	)

	var inits []app.Init

	// 初始化数据库
	inits = append(inits, func() {
		model.InitMysql()
		model.InitRedis()
	})

	// 初始化链路跟踪
	if config.Get().App.EnableTracing {
		inits = append(inits, func() {
			tracer.InitWithConfig(config.Get().App.Name, config.Get().App.Env, config.Get().App.Version,
				config.Get().Jaeger.AgentHost, config.Get().Jaeger.AgentPort, config.Get().Jaeger.SamplingRate)
		})
	}

	return inits
}

// 初始化配置
func initConfig() {
	flag.StringVar(&version, "version", "", "service Version Number")
	flag.BoolVar(&enableConfigCenter, "enable-cc", false, "whether to get from the configuration center, "+
		"if true, the '-c' parameter indicates the configuration center")
	flag.StringVar(&configFile, "c", "", "configuration file")
	flag.Parse()

	if enableConfigCenter {
		// 从配置中心获取配置(先获取nacos配置，再根据nacos配置中心读取服务配置)
		if configFile == "" {
			configFile = configs.Path("serverNameExample_cc.yml")
		}
		nacosConfig, err := config.NewCenter(configFile)
		if err != nil {
			panic(err)
		}
		appConfig := &config.Config{}
		params := &nacoscli.Params{}
		_ = copier.Copy(params, &nacosConfig.Nacos)
		err = nacoscli.Init(appConfig, params)
		if err != nil {
			panic(fmt.Sprintf("connect to configuration center err, %v", err))
		}
		if appConfig.App.Name == "" {
			panic("read the config from center error, config data is empty")
		}
		config.Set(appConfig)
	} else {
		// 从本地配置文件获取配置
		if configFile == "" {
			configFile = configs.Path("serverNameExample.yml")
		}
		err := config.Init(configFile)
		if err != nil {
			panic("init config error: " + err.Error())
		}
	}

	if version != "" {
		config.Get().App.Version = version
	}
	//fmt.Println(config.Show())
}

// -------------------------------- 注册app服务 ---------------------------------
func registerServers() []app.IServer {
	var servers []app.IServer

	// todo generate the code to register http and grpc services here
	// delete the templates code start
	// 创建http服务
	httpAddr := ":" + strconv.Itoa(config.Get().HTTP.Port)
	httpRegistry, httpInstance := registryService("http", config.Get().App.Host, config.Get().HTTP.Port)
	httpServer := server.NewHTTPServer(httpAddr,
		server.WithHTTPReadTimeout(time.Second*time.Duration(config.Get().HTTP.ReadTimeout)),
		server.WithHTTPWriteTimeout(time.Second*time.Duration(config.Get().HTTP.WriteTimeout)),
		server.WithHTTPRegistry(httpRegistry, httpInstance),
		server.WithHTTPIsProd(config.Get().App.Env == "prod"),
	)
	servers = append(servers, httpServer)

	// 创建grpc服务
	grpcAddr := ":" + strconv.Itoa(config.Get().Grpc.Port)
	grpcRegistry, grpcInstance := registryService("grpc", config.Get().App.Host, config.Get().Grpc.Port)
	grpcServer := server.NewGRPCServer(grpcAddr,
		server.WithGrpcReadTimeout(time.Duration(config.Get().Grpc.ReadTimeout)*time.Second),
		server.WithGrpcWriteTimeout(time.Duration(config.Get().Grpc.WriteTimeout)*time.Second),
		server.WithGrpcRegistry(grpcRegistry, grpcInstance),
	)
	servers = append(servers, grpcServer)
	// delete the templates code end

	return servers
}

func registryService(scheme string, host string, port int) (registry.Registry, *registry.ServiceInstance) {
	instanceEndpoint := fmt.Sprintf("%s://%s:%d", scheme, host, port)

	if config.Get().App.EnableRegistryDiscovery {
		// 使用consul注册服务
		if config.Get().App.RegistryDiscoveryType == "consul" {
			iRegistry, instance, err := consul.NewRegistry(
				config.Get().Consul.Addr,
				config.Get().App.Name+"_"+scheme+"_"+host,
				config.Get().App.Name,
				[]string{instanceEndpoint},
			)
			if err != nil {
				panic(err)
			}
			return iRegistry, instance
		}

		// 使用etcd注册服务
		if config.Get().App.RegistryDiscoveryType == "etcd" {
			iRegistry, instance, err := etcd.NewRegistry(
				config.Get().Etcd.Addrs,
				config.Get().App.Name+"_"+scheme+"_"+host,
				config.Get().App.Name,
				[]string{instanceEndpoint},
			)

			if err != nil {
				panic(err)
			}
			return iRegistry, instance
		}
	}

	return nil, nil
}

// -------------------------- 注册app需要释放的资源  -------------------------------------------

func registerCloses(servers []app.IServer) []app.Close {
	var closes []app.Close

	// 关闭服务
	for _, s := range servers {
		closes = append(closes, s.Stop)
	}

	// 关闭mysql
	closes = append(closes, func() error {
		return model.CloseMysql()
	})

	// 关闭redis
	closes = append(closes, func() error {
		return model.CloseRedis()
	})

	// 关闭trace
	if config.Get().App.EnableTracing {
		closes = append(closes, func() error {
			ctx, _ := context.WithTimeout(context.Background(), 2*time.Second) //nolint
			return tracer.Close(ctx)
		})
	}

	return closes
}
