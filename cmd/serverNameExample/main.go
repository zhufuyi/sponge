package main

import (
	"context"
	"flag"
	"strconv"
	"time"

	"github.com/zhufuyi/sponge/configs"
	"github.com/zhufuyi/sponge/internal/config"
	"github.com/zhufuyi/sponge/internal/model"
	"github.com/zhufuyi/sponge/internal/server"

	"github.com/zhufuyi/sponge/pkg/app"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/nacoscli"
	"github.com/zhufuyi/sponge/pkg/tracer"

	"github.com/jinzhu/copier"

	// only grpc use start
	"fmt"

	"github.com/zhufuyi/sponge/pkg/registry"
	"github.com/zhufuyi/sponge/pkg/registry/etcd"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// only grpc use end
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
			panic(err)
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
// todo generate the code to register http and grpc services here
// delete the templates code start

func registerServers() []app.IServer {
	var servers []app.IServer

	// 创建http服务
	httpAddr := ":" + strconv.Itoa(config.Get().HTTP.Port)
	httpServer := server.NewHTTPServer(httpAddr,
		server.WithHTTPReadTimeout(time.Second*time.Duration(config.Get().HTTP.ReadTimeout)),
		server.WithHTTPWriteTimeout(time.Second*time.Duration(config.Get().HTTP.WriteTimeout)),
		server.WithHTTPIsProd(config.Get().App.Env == "prod"),
	)
	servers = append(servers, httpServer)

	// 创建grpc服务
	grpcAddr := ":" + strconv.Itoa(config.Get().Grpc.Port)
	grpcServer := server.NewGRPCServer(grpcAddr, grpcOptions()...)
	servers = append(servers, grpcServer)

	return servers
}

func grpcOptions() []server.GRPCOption {
	var opts []server.GRPCOption

	if config.Get().App.EnableRegistryDiscovery {
		iRegistry, instance := getETCDRegistry(
			config.Get().Etcd.Addrs,
			config.Get().App.Name,
			[]string{fmt.Sprintf("grpc://%s:%d", config.Get().App.Host, config.Get().Grpc.Port)},
		)
		opts = append(opts, server.WithRegistry(iRegistry, instance))
	}

	return opts
}

// 使用etcd实例化服务注册，consul和nacos也类似
func getETCDRegistry(etcdEndpoints []string, instanceName string, instanceEndpoints []string) (registry.Registry, *registry.ServiceInstance) {
	serviceInstance := registry.NewServiceInstance(instanceName, instanceEndpoints)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	})
	if err != nil {
		panic(err)
	}
	iRegistry := etcd.New(cli)

	return iRegistry, serviceInstance
}

// delete the templates code end

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
