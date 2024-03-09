// Package benchmark is compression testing of rpc methods and generation of reported results.
package benchmark

import (
	"fmt"
	"os"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	"google.golang.org/protobuf/proto"

	"github.com/zhufuyi/sponge/pkg/gofile"
)

// Runner interface
type Runner interface {
	Run() error
}

// bench pressing parameters
type bench struct {
	rpcServerHost string // rpc server address

	protoFile     string        // proto file
	packageName   string        // proto file package name
	serviceName   string        // proto file service name
	methodName    string        // name of pressure test method
	methodRequest proto.Message // input parameters corresponding to the method

	total       uint     // number of requests
	importPaths []string // reliance on third party protobuf file locations
}

// New create a pressure test instance
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

// Run operational performance benchmarking
func (b *bench) Run() error {
	callMethod := fmt.Sprintf("%s.%s/%s", b.packageName, b.serviceName, b.methodName)
	fmt.Printf("benchmark '%s', total %d requests\n", callMethod, b.total)

	data, err := proto.Marshal(b.methodRequest)
	if err != nil {
		return err
	}

	report, err := runner.Run(
		callMethod, //  'package.Service/method' or 'package.Service.Method'
		b.rpcServerHost,
		runner.WithProtoFile(b.protoFile, b.importPaths),
		runner.WithBinaryData(data),
		runner.WithInsecure(true),
		runner.WithTotalRequests(b.total),
		// concurrent parameter
		//runner.WithConcurrencySchedule(runner.ScheduleLine),
		//runner.WithConcurrencyStep(5),  // add 5 workers per second
		//runner.WithConcurrencyStart(1), //
		//runner.WithConcurrencyEnd(20),  // maximum number of concurrent
	)
	if err != nil {
		return err
	}

	return b.saveReport(callMethod, report)
}

func (b *bench) saveReport(callMethod string, report *runner.Report) error {
	// specify the output path
	outputFile := fmt.Sprintf("%sreport_%s.html", os.TempDir()+gofile.GetPathDelimiter(), b.methodName)
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	rp := printer.ReportPrinter{
		Out:    file,
		Report: report,
	}

	fmt.Printf("benchmark '%s' finished, report file=%s\n", callMethod, outputFile)
	return rp.Print("html")
}
