// Package initial is the package that starts the service to initialize the service, including
// the initialization configuration, service configuration, connecting to the database, and
// resource release needed when shutting down the service.
package initial

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/jinzhu/copier"

	"github.com/go-dev-frame/sponge/pkg/conf"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/nacoscli"
	"github.com/go-dev-frame/sponge/pkg/stat"
	"github.com/go-dev-frame/sponge/pkg/tracer"

	"github.com/go-dev-frame/sponge/configs"
	"github.com/go-dev-frame/sponge/internal/config"
	"github.com/go-dev-frame/sponge/internal/database"
)

var (
	version            string
	configFile         string
	enableConfigCenter bool
)

// InitApp initial app configuration
func InitApp() {
	initConfig()
	cfg := config.Get()

	// initializing log
	_, err := logger.Init(
		logger.WithLevel(cfg.Logger.Level),
		logger.WithFormat(cfg.Logger.Format),
		logger.WithSave(
			cfg.Logger.IsSave,
			//logger.WithFileName(cfg.Logger.LogFileConfig.Filename),
			//logger.WithFileMaxSize(cfg.Logger.LogFileConfig.MaxSize),
			//logger.WithFileMaxBackups(cfg.Logger.LogFileConfig.MaxBackups),
			//logger.WithFileMaxAge(cfg.Logger.LogFileConfig.MaxAge),
			//logger.WithFileIsCompression(cfg.Logger.LogFileConfig.IsCompression),
		),
	)
	if err != nil {
		panic(err)
	}
	logger.Debug(config.Show())
	logger.Info("[logger] was initialized")

	// initializing tracing
	if cfg.App.EnableTrace {
		tracer.InitWithConfig(
			cfg.App.Name,
			cfg.App.Env,
			cfg.App.Version,
			cfg.Jaeger.AgentHost,
			strconv.Itoa(cfg.Jaeger.AgentPort),
			cfg.App.TracingSamplingRate,
		)
		logger.Info("[tracer] was initialized")
	}

	// initializing the print system and process resources
	if cfg.App.EnableStat {
		stat.Init(
			stat.WithLog(logger.Get()),
			stat.WithAlarm(), // invalid if it is windows, the default threshold for cpu and memory is 0.8, you can modify them
			stat.WithPrintField(logger.String("service_name", cfg.App.Name), logger.String("host", cfg.App.Host)),
		)
		logger.Info("[resource statistics] was initialized")
	}

	// initializing database
	database.InitDB()
	logger.Infof("[%s] was initialized", cfg.Database.Driver)
	database.InitCache(cfg.App.CacheType)
	if cfg.App.CacheType != "" {
		logger.Infof("[%s] was initialized", cfg.App.CacheType)
	}
}

func initConfig() {
	flag.StringVar(&version, "version", "", "service Version Number")
	flag.BoolVar(&enableConfigCenter, "enable-cc", false, "whether to get from the configuration center, "+
		"if true, the '-c' parameter indicates the configuration center")
	flag.StringVar(&configFile, "c", "", "configuration file")
	flag.Parse()

	if enableConfigCenter {
		getConfigFromNacos()
	} else {
		getConfigFromLocal()
	}

	if version != "" {
		config.Get().App.Version = version
	}
}

// get the configuration from the configuration center (first get the nacos configuration,
// then read the service configuration according to the nacos configuration center)
func getConfigFromNacos() {
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
	format, data, err := nacoscli.GetConfig(params)
	if err != nil {
		panic(fmt.Sprintf("connect to configuration center err, %v", err))
	}
	err = conf.ParseConfigData(data, format, appConfig)
	if err != nil {
		panic(fmt.Sprintf("parse configuration data err, %v", err))
	}
	if appConfig.App.Name == "" {
		panic("read the config from center error, config data is empty")
	}
	config.Set(appConfig)
}

// get configuration from local configuration file
func getConfigFromLocal() {
	if configFile == "" {
		configFile = configs.Path("serverNameExample.yml")
	}
	err := config.Init(configFile)
	if err != nil {
		panic("init config error: " + err.Error())
	}
}
