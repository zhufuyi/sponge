package generate

const (
	dockerFileHTTPCode = `# add curl, used for http service checking, can be installed without it if deployed in k8s
RUN apk add curl

COPY configs/ /app/configs/
COPY serverNameExample /app/serverNameExample
RUN chmod +x /app/serverNameExample

# http port
EXPOSE 8080`

	dockerFileGrpcCode = `# add grpc_health_probe for health check of grpc services
COPY grpc_health_probe /bin/grpc_health_probe
RUN chmod +x /bin/grpc_health_probe

COPY configs/ /app/configs/
COPY serverNameExample /app/serverNameExample
RUN chmod +x /app/serverNameExample

# grpc and http port
EXPOSE 8282 8283`

	dockerFileBuildHTTPCode = `# compressing binary files
#cd /
#upx -9 serverNameExample


# building images with binary
FROM alpine:latest
MAINTAINER zhufuyi "g.zhufuyi@gmail.com"

# set the time zone to Shanghai
RUN apk add tzdata  \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# add curl, used for http service checking, can be installed without it if deployed in k8s
RUN apk add curl

COPY --from=build /serverNameExample /app/serverNameExample
COPY --from=build /go/src/serverNameExample/configs/serverNameExample.yml /app/configs/serverNameExample.yml

# http port
EXPOSE 8080`

	dockerFileBuildGrpcCode = `# install grpc-health-probe, for health check of grpc service
RUN go install github.com/grpc-ecosystem/grpc-health-probe@v0.4.12
RUN cd $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-health-probe@v0.4.12 \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "all=-s -w" -o /grpc_health_probe

# compressing binary files
#cd /
#upx -9 serverNameExample
#upx -9 grpc_health_probe


# building images with binary
FROM alpine:latest
MAINTAINER zhufuyi "g.zhufuyi@gmail.com"

# set the time zone to Shanghai
RUN apk add tzdata  \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# add grpc_health_probe for health check of grpc services
COPY --from=build /grpc_health_probe /bin/grpc_health_probe
COPY --from=build /serverNameExample /app/serverNameExample
COPY --from=build /go/src/serverNameExample/configs/serverNameExample.yml /app/configs/serverNameExample.yml

# grpc and http port
EXPOSE 8282 8283`

	imageBuildFileHTTPCode = `# compressing binary file
#cd ${DOCKERFILE_PATH}
#upx -9 ${serverName}
#cd -

echo "docker build -f ${DOCKERFILE} -t ${IMAGE_NAME_TAG} ${DOCKERFILE_PATH}"
docker build -f ${DOCKERFILE} -t ${IMAGE_NAME_TAG} ${DOCKERFILE_PATH}`

	imageBuildFileGrpcCode = `# install grpc-health-probe, for health check of grpc service
rootDockerFilePath=$(pwd)/${DOCKERFILE_PATH}
go install github.com/grpc-ecosystem/grpc-health-probe@v0.4.12
cd $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-health-probe@v0.4.12 \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "all=-s -w" -o "${rootDockerFilePath}/grpc_health_probe"
cd -

# compressing binary file
#cd ${DOCKERFILE_PATH}
#upx -9 ${serverName}
#upx -9 grpc_health_probe
#cd -

echo "docker build -f ${DOCKERFILE} -t ${IMAGE_NAME_TAG} ${DOCKERFILE_PATH}"
docker build -f ${DOCKERFILE} -t ${IMAGE_NAME_TAG} ${DOCKERFILE_PATH}

if [ -f "${DOCKERFILE_PATH}/grpc_health_probe" ]; then
    rm -f ${DOCKERFILE_PATH}/grpc_health_probe
fi`

	imageBuildLocalFileHTTPCode = `# compressing binary file
#cd ${DOCKERFILE_PATH}
#upx -9 ${serverName}
#cd -

mkdir -p ${DOCKERFILE_PATH}/configs && cp -f configs/${serverName}.yml ${DOCKERFILE_PATH}/configs/
echo "docker build -f ${DOCKERFILE} -t ${IMAGE_NAME}:latest ${DOCKERFILE_PATH}"
docker build -f ${DOCKERFILE} -t ${IMAGE_NAME}:latest ${DOCKERFILE_PATH}`

	imageBuildLocalFileGrpcCode = `# install grpc-health-probe, for health check of grpc service
rootDockerFilePath=$(pwd)/${DOCKERFILE_PATH}
go install github.com/grpc-ecosystem/grpc-health-probe@v0.4.12
cd $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-health-probe@v0.4.12 \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "all=-s -w" -o "${rootDockerFilePath}/grpc_health_probe"
cd -

# compressing binary file
#cd ${DOCKERFILE_PATH}
#upx -9 ${serverName}
#upx -9 grpc_health_probe
#cd -

mkdir -p ${DOCKERFILE_PATH}/configs && cp -f configs/${serverName}.yml ${DOCKERFILE_PATH}/configs/
echo "docker build -f ${DOCKERFILE} -t ${IMAGE_NAME}:latest ${DOCKERFILE_PATH}"
docker build -f ${DOCKERFILE} -t ${IMAGE_NAME}:latest ${DOCKERFILE_PATH}

if [ -f "${DOCKERFILE_PATH}/grpc_health_probe" ]; then
    rm -f ${DOCKERFILE_PATH}/grpc_health_probe
fi`

	dockerComposeFileHTTPCode = `    ports:
      - "8080:8080"   # http port
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]   # http health check, note: mirror must contain curl command`

	dockerComposeFileGrpcCode = `
    ports:
      - "8282:8282"   # grpc port
      - "8283:8283"   # grpc metrics or pprof port
    healthcheck:
      test: ["CMD", "grpc_health_probe", "-addr=localhost:8282"]    # grpc health check, note: the image must contain the grpc_health_probe command`

	k8sDeploymentFileHTTPCode = `
          ports:
            - name: http-port
              containerPort: 8080
          readinessProbe:
            httpGet:
              port: http-port
              path: /health
            initialDelaySeconds: 10
            timeoutSeconds: 2
            periodSeconds: 10
            successThreshold: 1
            failureThreshold: 3
          livenessProbe:
            httpGet:
              port: http-port
              path: /health`

	k8sDeploymentFileGrpcCode = `
          ports:
            - name: grpc-port
              containerPort: 8282
            - name: metrics-port
              containerPort: 8283
          readinessProbe:
            exec:
              command: ["/bin/grpc_health_probe", "-addr=:8282"]
            initialDelaySeconds: 10
            timeoutSeconds: 2
            periodSeconds: 10
            successThreshold: 1
            failureThreshold: 3
          livenessProbe:
            exec:
              command: ["/bin/grpc_health_probe", "-addr=:8282"]`

	k8sServiceFileHTTPCode = `  ports:
    - name: server-name-example-svc-http-port
      port: 8080
      targetPort: 8080`

	k8sServiceFileGrpcCode = `  ports:
    - name: server-name-example-svc-grpc-port
      port: 8282
      targetPort: 8282
    - name: server-name-example-svc-grpc-metrics-port
      port: 8283
      targetPort: 8283`

	configFileCode = `// code generated by https://github.com/zhufuyi/sponge

package config

import (
	"github.com/zhufuyi/sponge/pkg/conf"
)

var config *Config

func Init(configFile string, fs ...func()) error {
	config = &Config{}
	return conf.Parse(configFile, config, fs...)
}

func Show() string {
	return conf.Show(config)
}

func Get() *Config {
	if config == nil {
		panic("config is nil")
	}
	return config
}

func Set(conf *Config) {
	config = conf
}
`

	configFileCcCode = `// code generated by https://github.com/zhufuyi/sponge

package config

import (
	"github.com/zhufuyi/sponge/pkg/conf"
)

func NewCenter(configFile string) (*Center, error) {
	nacosConf := &Center{}
	err := conf.Parse(configFile, nacosConf)
	return nacosConf, err
}
`

	protoShellGRPCMark = `
  # generate files *_grpc_pb.go
  protoc --proto_path=. --proto_path=./third_party \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    $allProtoFiles

  checkResult $?`

	// for rpc and rpc-pb
	protoShellServiceTmplCode = `
  moduleName=$(cat docs/gen.info | head -1 | cut -d , -f 1)
  serverName=$(cat docs/gen.info | head -1 | cut -d , -f 2)
  # Generate 2 files, a logic code template file *.go (default save path in internal/service), a return error code template file *_rpc.go (default save path in internal/ecode)
  protoc --proto_path=. --proto_path=./third_party \
    --go-rpc-tmpl_out=. --go-rpc-tmpl_opt=paths=source_relative \
    --go-rpc-tmpl_opt=moduleName=${moduleName} --go-rpc-tmpl_opt=serverName=${serverName} \
    $specifiedProtoFiles

  checkResult $?

  sponge merge rpc-pb
  checkResult $?

  colorCyan='\e[1;36m'
  highBright='\e[1m'
  markEnd='\e[0m'

  echo ""
  echo -e "${highBright}Tip:${markEnd} execute the command ${colorCyan}make run${markEnd} and then test rpc method is in the file ${colorCyan}internal/service/xxx_client_test.go${markEnd}."
  echo ""`

	// for http-pb
	protoShellHandlerCode = `
  # generate the swagger document and merge all files into docs/apis.swagger.json
  protoc --proto_path=. --proto_path=./third_party \
    --openapiv2_out=. --openapiv2_opt=logtostderr=true --openapiv2_opt=allow_merge=true --openapiv2_opt=merge_file_name=docs/apis.json \
    $specifiedProtoFiles

  checkResult $?

  sponge web swagger --file=docs/apis.swagger.json
  checkResult $?

  moduleName=$(cat docs/gen.info | head -1 | cut -d , -f 1)
  serverName=$(cat docs/gen.info | head -1 | cut -d , -f 2)
  # A total of four files are generated, namely the registration route file _*router.pb.go (saved in the same directory as the protobuf file),
  # the injection route file *_router.go (saved by default in the path internal/routers), the logical code template file *.go (default path
  # is in internal/handler), return error code template file*_http.go (default path is in internal/ecode)
  protoc --proto_path=. --proto_path=./third_party \
    --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugin=handler \
    --go-gin_opt=moduleName=${moduleName} --go-gin_opt=serverName=${serverName} \
    $specifiedProtoFiles

  checkResult $?

  sponge merge http-pb
  checkResult $?

  colorCyan='\e[1;36m'
  highBright='\e[1m'
  markEnd='\e[0m'

  echo ""
  echo -e "${highBright}Tip:${markEnd} execute the command ${colorCyan}make run${markEnd} and then visit ${colorCyan}http://localhost:8080/apis/swagger/index.html${markEnd} in your browser. "
  echo ""`

	// for rpc-gw
	protoShellServiceCode = `
  # Generate the swagger document and merge all files into docs/apis.swagger.json
  protoc --proto_path=. --proto_path=./third_party \
    --openapiv2_out=. --openapiv2_opt=logtostderr=true --openapiv2_opt=allow_merge=true --openapiv2_opt=merge_file_name=docs/apis.json \
    $specifiedProtoFiles

  checkResult $?

  sponge micro swagger --file=docs/apis.swagger.json
  checkResult $?

  moduleName=$(cat docs/gen.info | head -1 | cut -d , -f 1)
  serverName=$(cat docs/gen.info | head -1 | cut -d , -f 2)
  # A total of 4 files are generated, namely the registration route file _*router.pb.go (saved in the same directory as the protobuf file),
  # the injection route file *_router.go (default save path in internal/routers), the logical code template file *.go (saved in
  # internal/service by default), return error code template file*_rpc.go (saved in internal/ecode by default)
  protoc --proto_path=. --proto_path=./third_party \
    --go-gin_out=. --go-gin_opt=paths=source_relative --go-gin_opt=plugin=service \
    --go-gin_opt=moduleName=${moduleName} --go-gin_opt=serverName=${serverName} \
    $specifiedProtoFiles

  checkResult $?

  sponge merge rpc-gw-pb
  checkResult $?

  colorCyan='\e[1;36m'
  highBright='\e[1m'
  markEnd='\e[0m'

  echo ""
  echo -e "${highBright}Tip:${markEnd} execute the command ${colorCyan}make run${markEnd} and then visit ${colorCyan}http://localhost:8080/apis/swagger/index.html${markEnd} in your browser."
  echo ""`

	httpServerConfigCode = `# http server settings
http:
  port: 8080            # listen port
  readTimeout: 3     # read timeout, unit(second)
  writeTimeout: 60  # write timeout, unit(second), if enableHTTPProfile is true, it needs to be greater than 60s, the default value for pprof to do profiling is 60s`

	rpcServerConfigCode = `# grpc service settings
grpc:
  port: 8282             # listen port
  httpPort: 8283       # profile and metrics ports
  readTimeout: 5      # read timeout, unit(second)
  writeTimeout: 5     # write timeout, unit(second)
  enableToken: false  # whether to enable server-side token authentication, default appID=grpc, appKey=123456
  # serverSecure parameter setting
  # if type="", it means no secure connection, no need to fill in any parameters
  # if type="one-way", it means server-side certification, only the fields "certFile" and "keyFile" should be filled in
  # if type="two-way", it means both client and server side certification, fill in all fields
  serverSecure:
    type: ""               # secures type, "", "one-way", "two-way"
    certFile: ""           # server side cert file, absolute path
    keyFile: ""           # server side key file, absolute path
    caFile: ""             # ca certificate file, valid only in "two-way", absolute path



# grpc client settings, support for setting up multiple rpc clients
grpcClient:
  - name: "your-rpc-server-name"   # rpc service name, used for service discovery
    host: "127.0.0.1"                    # rpc service address, used for direct connection
    port: 8282                              # rpc service port
    registryDiscoveryType: ""         # registration and discovery types: consul, etcd, nacos, if empty, connecting to server using host and port
    enableLoadBalance: false         # whether to turn on the load balancer
    # clientSecure parameter setting
    # if type="", it means no secure connection, no need to fill in any parameters
    # if type="one-way", it means server-side certification, only the fields "serverName" and "certFile" should be filled in
    # if type="two-way", it means both client and server side certification, fill in all fields
    clientSecure:
      type: ""              # secures type, "", "one-way", "two-way"
      serverName: ""   # server name, e.g. *.foo.com
      caFile: ""            # client side ca file, valid only in "two-way", absolute path
      certFile: ""          # client side cert file, absolute path, if secureType="one-way", fill in server side cert file here
      keyFile: ""          # client side key file, valid only in "two-way", absolute path
    clientToken:
      enable: false      # whether to enable token authentication
      appID: ""           # app id
      appKey: ""         # app key`

	rpcGwServerConfigCode = `# http server settings
http:
  port: 8080            # listen port
  readTimeout: 3     # read timeout, unit(second)
  writeTimeout: 60  # write timeout, unit(second), if enableHTTPProfile is true, it needs to be greater than 60s, the default value for pprof to do profiling is 60s


# grpc client settings, support for setting up multiple rpc clients
grpcClient:
  - name: "your-rpc-server-name"   # rpc service name, used for service discovery
    host: "127.0.0.1"                    # rpc service address, used for direct connection
    port: 8282                              # rpc service port
    registryDiscoveryType: ""         # registration and discovery types: consul, etcd, nacos, if empty, connecting to server using host and port
    enableLoadBalance: false         # whether to turn on the load balancer
    # clientSecure parameter setting
    # if type="", it means no secure connection, no need to fill in any parameters
    # if type="one-way", it means server-side certification, only the fields "serverName" and "certFile" should be filled in
    # if type="two-way", it means both client and server side certification, fill in all fields
    clientSecure:
      type: ""              # secures type, "", "one-way", "two-way"
      serverName: ""   # server name, e.g. *.foo.com
      caFile: ""            # client side ca file, valid only in "two-way", absolute path
      certFile: ""          # client side cert file, absolute path, if secureType="one-way", fill in server side cert file here
      keyFile: ""          # client side key file, valid only in "two-way", absolute path
    clientToken:
      enable: false      # whether to enable token authentication
      appID: ""           # app id
      appKey: ""         # app key`
)
