// nolint
// code generated from config file

package config

import "github.com/zhufuyi/sponge/pkg/conf"

type Config = GenerateName

var config *Config

// Init parsing configuration files to struct, including yaml, toml, json, etc.
func Init(configFile string, fs ...func()) error {
	config = &Config{}
	return conf.Parse(configFile, config, fs...)
}

func Show() {
	conf.Show(config)
}

func Get() *Config {
	if config == nil {
		panic("config is nil")
	}
	return config
}

type GenerateName struct {
	App         App         `yaml:"app" json:"app"`
	Etcd        Etcd        `yaml:"etcd" json:"etcd"`
	Grpc        Grpc        `yaml:"grpc" json:"grpc"`
	HTTP        HTTP        `yaml:"http" json:"http"`
	Jaeger      Jaeger      `yaml:"jaeger" json:"jaeger"`
	Logger      Logger      `yaml:"logger" json:"logger"`
	Metrics     Metrics     `yaml:"metrics" json:"metrics"`
	Mysql       Mysql       `yaml:"mysql" json:"mysql"`
	RateLimiter RateLimiter `yaml:"rateLimiter" json:"rateLimiter"`
	Redis       Redis       `yaml:"redis" json:"redis"`
}

type Redis struct {
	Addr         string `yaml:"addr" json:"addr"`
	DB           int    `yaml:"dB" json:"dB"`
	DialTimeout  int    `yaml:"dialTimeout" json:"dialTimeout"`
	Dsn          string `yaml:"dsn" json:"dsn"`
	MinIdleConn  int    `yaml:"minIdleConn" json:"minIdleConn"`
	Password     string `yaml:"password" json:"password"`
	PoolSize     int    `yaml:"poolSize" json:"poolSize"`
	PoolTimeout  int    `yaml:"poolTimeout" json:"poolTimeout"`
	ReadTimeout  int    `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout" json:"writeTimeout"`
}

type Etcd struct {
	Addrs []string `yaml:"addrs" json:"addrs"`
}

type Jaeger struct {
	AgentHost    string  `yaml:"agentHost" json:"agentHost"`
	AgentPort    string  `yaml:"agentPort" json:"agentPort"`
	SamplingRate float64 `yaml:"samplingRate" json:"samplingRate"`
}

type Mysql struct {
	ConnMaxLifetime int    `yaml:"connMaxLifetime" json:"connMaxLifetime"`
	Dsn             string `yaml:"dsn" json:"dsn"`
	EnableLog       bool   `yaml:"enableLog" json:"enableLog"`
	MaxIdleConns    int    `yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns" json:"maxOpenConns"`
	SlowThreshold   int    `yaml:"slowThreshold" json:"slowThreshold"`
}

type RateLimiter struct {
	Dimension string `yaml:"dimension" json:"dimension"`
	MaxLimit  int    `yaml:"maxLimit" json:"maxLimit"`
	QPSLimit  int    `yaml:"qpsLimit" json:"qpsLimit"`
}

type App struct {
	EnableRegistryDiscovery bool   `yaml:"enableRegistryDiscovery" json:"enableRegistryDiscovery"`
	EnableLimit             bool   `yaml:"enableLimit" json:"enableLimit"`
	EnableMetrics           bool   `yaml:"enableMetrics" json:"enableMetrics"`
	EnableProfile           bool   `yaml:"enableProfile" json:"enableProfile"`
	EnableTracing           bool   `yaml:"enableTracing" json:"enableTracing"`
	Env                     string `yaml:"env" json:"env"`
	HostIP                  string `yaml:"hostIP" json:"hostIP"`
	Name                    string `yaml:"name" json:"name"`
	Version                 string `yaml:"version" json:"version"`
}

type LogFileConfig struct {
	Filename      string `yaml:"filename" json:"filename"`
	IsCompression bool   `yaml:"isCompression" json:"isCompression"`
	MaxAge        int    `yaml:"maxAge" json:"maxAge"`
	MaxBackups    int    `yaml:"maxBackups" json:"maxBackups"`
	MaxSize       int    `yaml:"maxSize" json:"maxSize"`
}

type Logger struct {
	Format        string        `yaml:"format" json:"format"`
	IsSave        bool          `yaml:"isSave" json:"isSave"`
	Level         string        `yaml:"level" json:"level"`
	LogFileConfig LogFileConfig `yaml:"logFileConfig" json:"logFileConfig"`
}

type HTTP struct {
	Port         int    `yaml:"port" json:"port"`
	ReadTimeout  int    `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout" json:"writeTimeout"`
	ServiceName  string `yaml:"serviceName" json:"serviceName"`
}

type Grpc struct {
	Port         int    `yaml:"port" json:"port"`
	ReadTimeout  int    `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout" json:"writeTimeout"`
	ServiceName  string `yaml:"serviceName" json:"serviceName"`
}

type Metrics struct {
	Port int `yaml:"port" json:"port"`
}
