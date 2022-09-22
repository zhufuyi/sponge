package main

import (
	"context"
	"flag"

	"strconv"
	"time"

	"github.com/zhufuyi/sponge/config"
	"github.com/zhufuyi/sponge/internal/model"
	"github.com/zhufuyi/sponge/internal/server"
	"github.com/zhufuyi/sponge/pkg/app"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/tracer"

	// grpc import start
	"fmt"

	"github.com/zhufuyi/sponge/pkg/registry"
	"github.com/zhufuyi/sponge/pkg/registry/etcd"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// grpc import end
)

var (
	version    string
	configFile string
)

// @title sponge api docs
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
	// 初始化配置文件，必须优先执行，后面的初始化需要依赖配置
	func() {
		flag.StringVar(&configFile, "c", "", "配置文件")
		flag.StringVar(&version, "version", "", "服务版本号")
		flag.Parse()
		if configFile == "" {
			configFile = config.Path("conf.yml") // 默认配置文件config/conf.yml
		}
		err := config.Init(configFile)
		if err != nil {
			panic("init config error: " + err.Error())
		}
		if version != "" {
			config.Get().App.Version = version
		}
		//config.Show()
	}()

	// 执行初始化日志
	func() {
		_, err := logger.Init(
			logger.WithLevel(config.Get().Logger.Level),
			logger.WithFormat(config.Get().Logger.Format),
			logger.WithSave(config.Get().Logger.IsSave,
				logger.WithFileName(config.Get().Logger.LogFileConfig.Filename),
				logger.WithFileMaxSize(config.Get().Logger.LogFileConfig.MaxSize),
				logger.WithFileMaxBackups(config.Get().Logger.LogFileConfig.MaxBackups),
				logger.WithFileMaxAge(config.Get().Logger.LogFileConfig.MaxAge),
				logger.WithFileIsCompression(config.Get().Logger.LogFileConfig.IsCompression),
			),
		)
		if err != nil {
			panic("init logger error: " + err.Error())
		}
	}()

	var inits []app.Init

	// 初始化数据库
	inits = append(inits, func() {
		model.InitMysql()
		model.InitRedis()
	})

	if config.Get().App.EnableTracing { // 根据配置是否开启链路跟踪
		inits = append(inits, func() {
			// 初始化链路跟踪
			exporter, err := tracer.NewJaegerAgentExporter(config.Get().Jaeger.AgentHost, config.Get().Jaeger.AgentPort)
			if err != nil {
				panic("init trace error:" + err.Error())
			}
			resource := tracer.NewResource(
				tracer.WithServiceName(config.Get().App.Name),
				tracer.WithEnvironment(config.Get().App.Env),
				tracer.WithServiceVersion(config.Get().App.Version),
			)

			tracer.Init(exporter, resource, config.Get().Jaeger.SamplingRate) // 如果SamplingRate=0.5表示只采样50%
		})
	}

	return inits
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
