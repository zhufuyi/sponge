package commands

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/spf13/cobra"

	"github.com/zhufuyi/sponge/pkg/conf"
	"github.com/zhufuyi/sponge/pkg/gofile"
)

// GenGraphCommand generate graph command
func GenGraphCommand() *cobra.Command {
	var (
		isAll      bool
		projectDir string
		serverDir  []string
	)

	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Generate business architecture diagram for the project",
		Long:  "Generate business architecture diagram for the project.",
		Example: color.HiBlackString(`  # If there are multiple servers in a project, simply specify the project directory path to generate a diagram between servers 
  sponge graph --project-dir=/path/to/project

  # You can also specify multiple services to generate a business framework diagram
  sponge graph --server-dir=/path/to/server1 --server-dir=/path/to/server2

  # Includes database related servers
  sponge graph --project-dir=/path/to/project --all`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if projectDir == "" && len(serverDir) == 0 {
				return errors.New("no project directory or server directory specified\n\n" + cmd.Example)
			}

			// get yaml file from project and server directories
			yamlFiles := getYamlFiles(projectDir, serverDir)
			if len(yamlFiles) == 0 {
				fmt.Println("No yaml file found in project directory or server directories")
				return nil
			}

			outFile := "./business_architecture_diagram.svg"
			err := generateSvg(yamlFiles, isAll, outFile)
			if err != nil {
				return err
			}
			fmt.Printf("generated servers relationship diagram successfully, out = %s\n", color.HiCyanString("%s", outFile))
			return nil
		},
	}

	cmd.Flags().BoolVarP(&isAll, "all", "a", false, "does it include services such as databases")
	cmd.Flags().StringVarP(&projectDir, "project-dir", "p", "", "project directory")
	cmd.Flags().StringSliceVarP(&serverDir, "server-dir", "s", []string{}, "server directory, multiple parameters can be set")

	return cmd
}

func mergeYamlFile(files []string, yamlFiles []string) []string {
	for _, file := range files {
		if !strings.Contains(file, "_cc.") {
			yamlFiles = append(yamlFiles, file)
		}
	}
	return yamlFiles
}

func filterYamlFiles(configsDirs []string, yamlFiles []string) []string {
	for _, dir := range configsDirs {
		files, _ := gofile.ListFiles(dir, gofile.WithSuffix(".yaml"))
		yamlFiles = mergeYamlFile(files, yamlFiles)
		files, _ = gofile.ListFiles(dir, gofile.WithSuffix(".yml"))
		yamlFiles = mergeYamlFile(files, yamlFiles)
	}
	return yamlFiles
}

func getYamlFiles(projectDir string, serverDirs []string) []string {
	var yamlFiles []string

	configsDirs, _ := gofile.ListSubDirs(projectDir, "configs")
	yamlFiles = filterYamlFiles(configsDirs, yamlFiles)

	for _, dir := range serverDirs {
		configsDirs, _ = gofile.ListSubDirs(dir, "configs")
		yamlFiles = filterYamlFiles(configsDirs, yamlFiles)
	}
	return yamlFiles
}

// -------------------------------------------------------------------------------------

// Service represents a service in the project.
type Service struct {
	Name         string              `yaml:"name"`         // service name
	Type         string              `yaml:"type"`         // http, grpc, db, mq
	Dependencies map[string][]string `yaml:"dependencies"` // map[Type][]serviceName
}

// ProjectConfig represents the configuration of the project.
type ProjectConfig struct {
	Services []Service `yaml:"services"`
}

// NewProjectConfig creates a new ProjectConfig.
func NewProjectConfig() *ProjectConfig {
	return &ProjectConfig{}
}

// AddService adds a new service to the project configuration.
func (c *ProjectConfig) AddService(config *GenericConfig) {
	name := config.GetServiceName()
	dependencies, additionalServices := config.GetDependencies(name)

	service := Service{
		Name:         name,
		Type:         config.GetServiceType(),
		Dependencies: dependencies,
	}
	additionalServices = append(additionalServices, service)

	c.merge(additionalServices...)
}

func (c *ProjectConfig) merge(services ...Service) {
	for _, service := range services {
		for i, s := range c.Services {
			if s.Name == service.Name {
				// merge dependencies
				s.Name = service.Name
				s.Type = service.Type
				for k, v := range service.Dependencies {
					s.Dependencies[k] = append(s.Dependencies[k], v...)
				}
				c.Services[i] = s
				//continue
			}
		}
		c.Services = append(c.Services, service)
	}
}

func generateSvg(yamlFiles []string, isAll bool, outFile string) error {
	pc := NewProjectConfig()

	for _, file := range yamlFiles {
		config, err := ParseYaml(file, isAll)
		if err != nil {
			return fmt.Errorf("Failed to parse YAML: %v, file: %s", err, file)
		}
		pc.AddService(config)
	}

	// create Graphviz graph
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		return fmt.Errorf("Failed to create Graphviz graph: %v", err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			fmt.Printf("Failed to close Graphviz graph: %v", err)
		}
		_ = g.Close()
	}()

	setCustomGraphStyle(graph)
	edgeMap := map[string]*cgraph.Edge{}

	// add nodes and edges to graph
	for _, service := range pc.Services {
		node1, _ := graph.CreateNode(service.Name)
		setCustomNodeStyle(node1, getNodeColor(service.Name))
		for typeLabel, dependencies := range service.Dependencies {
			for _, serviceName := range dependencies {
				key := sortAndMergeFields(node1.Name(), serviceName, typeLabel)
				if e, ok := edgeMap[key]; ok {
					setCustomEdgeStyle(e, typeLabel, getEdgeColor(typeLabel), true)
					continue
				}

				node2, _ := graph.CreateNode(serviceName)
				setCustomNodeStyle(node2, getNodeColor(serviceName))
				e, _ := graph.CreateEdge(typeLabel, node1, node2)
				setCustomEdgeStyle(e, typeLabel, getEdgeColor(typeLabel), false)
				edgeMap[key] = e
			}
		}
	}

	return g.RenderFilename(graph, graphviz.SVG, outFile)
}

// Attribute style https://graphviz.gitlab.io/docs/attr-types/style/

func setCustomNodeStyle(n *cgraph.Node, nodeColor string) {
	//color := "#99BD25"
	n.SetStyle("rounded,filled")
	n.SetShape(cgraph.RectShape)
	n.SetWidth(1.2)
	n.SetHeight(0.66)
	n.SetColor(nodeColor)
	n.SetFillColor(nodeColor)
	n.SetFontColor("#ffffff")
	n.SetFontSize(22.0)
}

func setCustomEdgeStyle(e *cgraph.Edge, typeLabel string, edgeColor string, isBothArrow bool) {
	//color := "#909090"
	e.SetLabel(typeLabel)
	e.SetStyle("solid")   // solid, dashed, dotted, bold
	e.SetColor(edgeColor) // edge color
	e.SetArrowSize(0.6)
	e.SetFontColor(edgeColor) // label font color
	e.SetFontSize(10.0)
	e.SetConstraint(true)

	if isBothArrow {
		e.SetDir(cgraph.BothDir)
	}

	// rewriting the db edge style
	if typeLabel == "db" {
		e.SetLabel("")
		e.SetStyle("dashed")
		e.SetArrowSize(0.4)
	}
}

func setCustomGraphStyle(g *cgraph.Graph) {
	//g.SetBackgroundColor("#eeeeee")
	g.SetCenter(true)
}

func getEdgeColor(typeLabel string) string {
	switch typeLabel {
	case "mq":
		return "#16B69E"
	case "db":
		return "#3999C6"
	case "http":
		return "#955E42"
	default:
		return "#909090"
	}
}

func getNodeColor(name string) string {
	name = strings.ToLower(name)
	switch name {
	case "kafka", "rabbitmq", "rocketmq", "nsq", "nats", "activemq", "pulsar", "zeromq":
		return "#16B69E"
	case "mysql", "redis", "mongodb", "postgresql", "sqlite", "oracle", "sqlserver", "cassandra", "influxdb", "elasticsearch", "clickhouse", "cockroachdb", "tidb":
		return "#3BD0FB"
	default:
		return "#99BD25"
	}
}

func sortAndMergeFields(field ...string) string {
	sort.Strings(field)
	return strings.Join(field, "-")
}

// --------------------------------- sponge yaml file --------------------------------------

// ParseYaml parses the YAML file and returns a GenericConfig, yaml file content example:
/*
app:
  name: "eshop-gw"

# http server settings example
http:
  port: 8080
  timeout: 0

# http client settings example
httpClient:
  - name: "flashSale"
    baseURL: "http://127.0.0.1:8080"
  - name: "product"
    baseURL: "http://127.0.0.1:8081"


# grpc server settings example
grpc:
  port: 8282
  httpPort: 8283

# grpc client settings example
grpcClient:
  - name: "user"
    host: "127.0.0.1"
    port: 18282
  - name: "order"
    host: "127.0.0.1"
    port: 28282


# db settings example
database:
  mysql:
    dsn: "root:123456@(192.168.3.37:3306)/eshop_order?parseTime=true&loc=Local&charset=utf8,utf8mb4"
  mongodb:
    dsn: "root:123456@192.168.3.37:27017/account?connectTimeoutMS=15000"

# or
redis:
  dsn: "default:123456@192.168.3.37:6379/0"

postgresql:
  dsn: "root:123456@192.168.3.37:5432/account?sslmode=disable"


# mq settings example
kafka:
  mode: "producer,consumer"
  brokers: ["192.168.3.37:9092"]

rabbitmq:
  mode: "consumer"
  host: "192.168.3.37"
  port: 5672
*/
func ParseYaml(configFile string, isAll bool) (*GenericConfig, error) {
	config := &GenericConfig{}
	err := conf.Parse(configFile, config)
	config.isAll = isAll
	return config, err
}

type GenericConfig struct {
	App        App          `yaml:"app" json:"app"`
	Grpc       Grpc         `yaml:"grpc" json:"grpc"`
	GrpcClient []GrpcClient `yaml:"grpcClient" json:"grpcClient"`
	HTTP       HTTP         `yaml:"http" json:"http"`
	HTTPClient []HTTPClient `yaml:"httpClient" json:"httpClient"`

	// database config, match 2 configuration modes, one is separate configuration, the other is unified configuration
	Database      Database               `yaml:"database" json:"database"`
	Mysql         map[string]interface{} `yaml:"mysql" json:"mysql"`
	Redis         map[string]interface{} `yaml:"redis" json:"redis"`
	Mongodb       map[string]interface{} `yaml:"mongodb" json:"mongodb"`
	Postgresql    map[string]interface{} `yaml:"postgresql" json:"postgresql"`
	Sqlite        map[string]interface{} `yaml:"sqlite" json:"sqlite"`
	Oracle        map[string]interface{} `yaml:"oracle" json:"oracle"`
	SQLServer     map[string]interface{} `yaml:"sqlserver" json:"sqlserver"`
	Cassandra     map[string]interface{} `yaml:"cassandra" json:"cassandra"`
	InfluxDB      map[string]interface{} `yaml:"influxdb" json:"influxdb"`
	Elasticsearch map[string]interface{} `yaml:"elasticsearch" json:"elasticsearch"`
	Clickhouse    map[string]interface{} `yaml:"clickhouse" json:"clickhouse"`
	Cockroachdb   map[string]interface{} `yaml:"cockroachdb" json:"cockroachdb"`
	Tidb          map[string]interface{} `yaml:"tidb" json:"tidb"`
	// add more database

	// mq config, match 2 configuration modes, one is separate configuration, the other is unified configuration
	MqClient []MqClient             `yaml:"mqClient" json:"mqClient"`
	Rabbitmq map[string]interface{} `yaml:"rabbitmq" json:"rabbitmq"`
	Kafka    map[string]interface{} `yaml:"kafka" json:"kafka"`
	Activemq map[string]interface{} `yaml:"activemq" json:"activemq"`
	Rocketmq map[string]interface{} `yaml:"rocketmq" json:"rocketmq"`
	Nats     map[string]interface{} `yaml:"nats" json:"nats"`
	Nsq      map[string]interface{} `yaml:"nsq" json:"nsq"`
	Asynq    map[string]interface{} `yaml:"asynq" json:"asynq"`
	Pulsar   map[string]interface{} `yaml:"pulsar" json:"pulsar"`
	Zeromq   map[string]interface{} `yaml:"zeromq" json:"zeromq"`
	// add more mq

	isAll bool
}

type HTTPClient struct {
	BaseURL string `yaml:"baseURL" json:"baseURL"`
	Name    string `yaml:"name" json:"name"`
}

type Grpc struct {
	HTTPPort int `yaml:"httpPort" json:"httpPort"`
	Port     int `yaml:"port" json:"port"`
}

type GrpcClient struct {
	Host string `yaml:"host" json:"host"`
	Name string `yaml:"name" json:"name"`
	Port int    `yaml:"port" json:"port"`
}

type MqClient struct {
	Name string `yaml:"name" json:"name"` // rabbitmq, kafka, activemq, rocketmq, nats, nsq, redis, asynq, pulsar, zeromq, etc.
	Mode string `yaml:"mode" json:"mode"` // producer, consumer.

	Rabbitmq map[string]interface{} `yaml:"rabbitmq" json:"rabbitmq"`
	Kafka    map[string]interface{} `yaml:"kafka" json:"kafka"`
	Activemq map[string]interface{} `yaml:"activemq" json:"activemq"`
	Rocketmq map[string]interface{} `yaml:"rocketmq" json:"rocketmq"`
	Nats     map[string]interface{} `yaml:"nats" json:"nats"`
	Nsq      map[string]interface{} `yaml:"nsq" json:"nsq"`
	Redis    map[string]interface{} `yaml:"redis" json:"redis"`
	Asynq    map[string]interface{} `yaml:"asynq" json:"asynq"`
	Pulsar   map[string]interface{} `yaml:"pulsar" json:"pulsar"`
	Zeromq   map[string]interface{} `yaml:"zeromq" json:"zeromq"`
	// add more mq
}

type App struct {
	Name                  string `yaml:"name" json:"name"`
	EnableTrace           bool   `yaml:"enableTrace" json:"enableTrace"`
	RegistryDiscoveryType string `yaml:"registryDiscoveryType" json:"registryDiscoveryType"`
	CacheType             string `yaml:"cacheType" json:"cacheType"`
}

type HTTP struct {
	Port    int `yaml:"port" json:"port"`
	Timeout int `yaml:"timeout" json:"timeout"`
}

type Database struct {
	Driver        string `yaml:"driver" json:"driver"`
	Mysql         DbAddr `yaml:"mysql" json:"mysql"`
	Redis         DbAddr `yaml:"redis" json:"redis"`
	Mongodb       DbAddr `yaml:"mongodb" json:"mongodb"`
	Postgresql    DbAddr `yaml:"postgresql" json:"postgresql"`
	Sqlite        DbAddr `yaml:"sqlite" json:"sqlite"`
	Oracle        DbAddr `yaml:"oracle" json:"oracle"`
	SQLServer     DbAddr `yaml:"sqlserver" json:"sqlserver"`
	Cassandra     DbAddr `yaml:"cassandra" json:"cassandra"`
	InfluxDB      DbAddr `yaml:"influxdb" json:"influxdb"`
	Elasticsearch DbAddr `yaml:"elasticsearch" json:"elasticsearch"`
	Clickhouse    DbAddr `yaml:"clickhouse" json:"clickhouse"`
	Cockroachdb   DbAddr `yaml:"cockroachdb" json:"cockroachdb"`
	Tidb          DbAddr `yaml:"tidb" json:"tidb"`
	// add more database
}

type DbAddr struct {
	Dsn    string `yaml:"dsn" json:"dsn"`
	DBFile string `yaml:"dbFile" json:"dbFile"` // sqlite
}

// GetServiceName returns the name of the service.
func (c *GenericConfig) GetServiceName() string {
	return c.App.Name
}

// GetServiceType returns the type of the service.
func (c *GenericConfig) GetServiceType() string {
	var serviceType string
	if c.HTTP.Port > 0 {
		serviceType = "http"
	} else if c.Grpc.Port > 0 {
		serviceType = "grpc"
	}
	return serviceType
}

// GetDependencies returns the dependencies of the service.
func (c *GenericConfig) GetDependencies(serviceName string) (map[string][]string, []Service) {
	mapDependencies := make(map[string][]string)

	var httpDependencies []string
	for _, client := range c.HTTPClient {
		if client.Name != "" {
			httpDependencies = append(httpDependencies, client.Name)
		}
	}
	if len(httpDependencies) > 0 {
		mapDependencies["http"] = httpDependencies
	}

	var grpcDependencies []string
	for _, client := range c.GrpcClient {
		if client.Name != "" && client.Name != c.App.Name && client.Name != "your_grpc_service_name" && client.Name != "your-rpc-server-name" {
			grpcDependencies = append(grpcDependencies, client.Name)
		}
	}
	if len(grpcDependencies) > 0 {
		mapDependencies["grpc"] = grpcDependencies
	}

	if c.isAll {
		dbDependencies := c.getDBDependencies()
		if len(dbDependencies) > 0 {
			mapDependencies["db"] = dbDependencies
		}
	}

	mqDependencies, services := c.getMqDependencies(serviceName)
	if len(mqDependencies) > 0 {
		mapDependencies["mq"] = mqDependencies
	}

	return mapDependencies, services
}

// nolint
func (c *GenericConfig) getDBDependencies() []string {
	db := c.Database
	var dbDependencies []string

	if db.Mysql.Dsn != "" || len(c.Mysql) > 0 {
		dbDependencies = append(dbDependencies, "mysql")
	}
	if db.Redis.Dsn != "" || len(c.Redis) > 0 {
		dbDependencies = append(dbDependencies, "redis")
	}
	if db.Mongodb.Dsn != "" || len(c.Mongodb) > 0 {
		dbDependencies = append(dbDependencies, "mongodb")
	}
	if db.Postgresql.Dsn != "" || len(c.Postgresql) > 0 {
		dbDependencies = append(dbDependencies, "postgresql")
	}
	if db.Sqlite.DBFile != "" || db.Sqlite.Dsn != "" || len(c.Sqlite) > 0 {
		dbDependencies = append(dbDependencies, "sqlite")
	}
	if db.Oracle.Dsn != "" || len(c.Oracle) > 0 {
		dbDependencies = append(dbDependencies, "oracle")
	}
	if db.SQLServer.Dsn != "" || len(c.SQLServer) > 0 {
		dbDependencies = append(dbDependencies, "sqlserver")
	}
	if db.Cassandra.Dsn != "" || len(c.Cassandra) > 0 {
		dbDependencies = append(dbDependencies, "cassandra")
	}
	if db.InfluxDB.Dsn != "" || len(c.InfluxDB) > 0 {
		dbDependencies = append(dbDependencies, "influxdb")
	}
	if db.Elasticsearch.Dsn != "" || len(c.Elasticsearch) > 0 {
		dbDependencies = append(dbDependencies, "elasticsearch")
	}
	if db.Clickhouse.Dsn != "" || len(c.Clickhouse) > 0 {
		dbDependencies = append(dbDependencies, "clickhouse")
	}
	if db.Cockroachdb.Dsn != "" || len(c.Cockroachdb) > 0 {
		dbDependencies = append(dbDependencies, "cockroachdb")
	}
	if db.Tidb.Dsn != "" || len(c.Tidb) > 0 {
		dbDependencies = append(dbDependencies, "tidb")
	}

	return dbDependencies
}

func (c *GenericConfig) getMqDependencies(serviceName string) ([]string, []Service) {
	var mqDependencies []string
	var services []Service

	for _, client := range c.MqClient {
		mqMode := strings.ReplaceAll(client.Mode, " ", "")
		mqDependencies, services = addMqDependencies(mqMode, client.Name, serviceName, mqDependencies, services)
	}

	if len(c.Rabbitmq) > 0 {
		mqDependencies, services = addMqDependencies(getMqMode(c.Rabbitmq), "rabbitmq", serviceName, mqDependencies, services)
	}
	if len(c.Kafka) > 0 {
		mqDependencies, services = addMqDependencies(getMqMode(c.Kafka), "kafka", serviceName, mqDependencies, services)
	}
	if len(c.Activemq) > 0 {
		mqDependencies, services = addMqDependencies(getMqMode(c.Activemq), "activemq", serviceName, mqDependencies, services)
	}
	if len(c.Rocketmq) > 0 {
		mqDependencies, services = addMqDependencies(getMqMode(c.Rocketmq), "rocketmq", serviceName, mqDependencies, services)
	}
	if len(c.Nats) > 0 {
		mqDependencies, services = addMqDependencies(getMqMode(c.Nats), "nats", serviceName, mqDependencies, services)
	}
	if len(c.Nsq) > 0 {
		mqDependencies, services = addMqDependencies(getMqMode(c.Nsq), "nsq", serviceName, mqDependencies, services)
	}
	if len(c.Asynq) > 0 {
		mqDependencies, services = addMqDependencies(getMqMode(c.Asynq), "asynq", serviceName, mqDependencies, services)
	}
	if len(c.Pulsar) > 0 {
		mqDependencies, services = addMqDependencies(getMqMode(c.Pulsar), "pulsar", serviceName, mqDependencies, services)
	}
	if len(c.Zeromq) > 0 {
		mqDependencies, services = addMqDependencies(getMqMode(c.Zeromq), "zeromq", serviceName, mqDependencies, services)
	}

	return mqDependencies, services
}

func getMqMode(mq map[string]interface{}) string {
	mqMode := "producer"
	if mode, ok := mq["mode"]; ok {
		if v, ok2 := mode.(string); ok2 {
			mqMode = v
		}
	}
	return mqMode
}

func addMqDependencies(mqMode string, mqName string, serviceName string, mqDependencies []string, services []Service) ([]string, []Service) {
	isOnlyConsumer, s := checkMqMode(mqMode, mqName, serviceName)
	if !isOnlyConsumer {
		mqDependencies = append(mqDependencies, mqName)
	}
	if s.Name != "" {
		services = append(services, s)
	}
	return mqDependencies, services
}

func checkMqMode(mqMode string, mqName string, serviceName string) (bool, Service) {
	isOnlyConsumer := false
	switch mqMode {
	case "consumer", "producer-consumer", "producer,consumer", "producer|consumer", "producerconsumer":
		if mqMode == "consumer" {
			isOnlyConsumer = true
		}
		s := Service{
			Name: mqName,
			Type: "mq",
			Dependencies: map[string][]string{
				"mq": {serviceName},
			},
		}
		return isOnlyConsumer, s
	}
	return isOnlyConsumer, Service{}
}
