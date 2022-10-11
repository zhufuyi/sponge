package benchmark

import (
	"fmt"
	"os"

	"github.com/zhufuyi/sponge/pkg/gofile"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	"github.com/golang/protobuf/proto"
)

type Runner interface {
	Run() error
}

// bench 压测参数
type bench struct {
	rpcServerHost string // rpc server端地址

	protoFile     string        // proto file
	packageName   string        // proto file的package
	serviceName   string        // proto file的Service
	methodName    string        // 当前压测方法名称
	methodRequest proto.Message // 压测输入参数，也就也方法参数

	total       uint     // 请求数量
	importPaths []string // 依赖第三方protobuf文件位置
}

// New 创建一个压测实例
func New(host string, protoFile string, methodName string, req proto.Message, total uint, importPaths ...string) (Runner, error) {
	data, err := os.ReadFile(protoFile)
	if err != nil {
		return nil, err
	}

	packageName := getName(data, packagePattern)
	if packageName == "" {
		return nil, fmt.Errorf("not found package name in protobuf file %s", protoFile)
	}

	serviceName := getName(data, servicePattern)
	if serviceName == "" {
		return nil, fmt.Errorf("not found service name in protobuf file %s", protoFile)
	}

	methodNames := getMethodNames(data, methodPattern)
	mName := matchName(methodNames, methodName)
	if mName == "" {
		return nil, fmt.Errorf("not found name %s in protobuf file %s", methodName, protoFile)
	}

	return &bench{
		protoFile:     protoFile,
		packageName:   packageName,
		serviceName:   serviceName,
		methodName:    mName,
		methodRequest: req,
		rpcServerHost: host,
		total:         total,
		importPaths:   importPaths,
	}, nil
}

func (b *bench) Run() error {
	callMethod := fmt.Sprintf("%s.%s/%s", b.packageName, b.serviceName, b.methodName)
	fmt.Printf("benchmark '%s', total %d requests\n", callMethod, b.total)

	buf := proto.Buffer{}
	err := buf.EncodeMessage(b.methodRequest)
	if err != nil {
		return err
	}

	report, err := runner.Run(
		callMethod, //  'package.Service/method' or 'package.Service.Method'
		b.rpcServerHost,
		runner.WithProtoFile(b.protoFile, b.importPaths),
		runner.WithBinaryData(buf.Bytes()),
		runner.WithInsecure(true),
		runner.WithTotalRequests(b.total),
		// 并发参数
		//runner.WithConcurrencySchedule(runner.ScheduleLine),
		//runner.WithConcurrencyStep(5),  // 每秒增加5个worker
		//runner.WithConcurrencyStart(1), //
		//runner.WithConcurrencyEnd(20),  // 最大并发数
	)
	if err != nil {
		return err
	}

	// 指定输出路径
	outputFile := fmt.Sprintf("%sreport_%s.html", gofile.GetRunPath()+gofile.GetPathDelimiter(), b.methodName)
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	rp := printer.ReportPrinter{
		Out:    file,
		Report: report,
	}

	fmt.Printf("benchmark '%s' finished, report file=%s\n", callMethod, outputFile)
	return rp.Print("html")
}
