// Package benchmark is compression testing of rpc methods and generation of reported results.
package benchmark

import (
	"fmt"
	"os"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	"google.golang.org/protobuf/proto"
)

type Option = runner.Option

// Runner interface
type Runner interface {
	Run() error
}

// bench pressing parameters
type bench struct {
	rpcServerHost string // rpc server address

	protoFile              string        // proto file
	packageName            string        // proto file package name
	serviceName            string        // proto file service name
	methodName             string        // name of pressure test method
	methodRequest          proto.Message // input parameters corresponding to the method
	dependentProtoFilePath []string      // dependent proto file path

	total uint // number of requests

	options []runner.Option
}

// New create a pressure test instance
//
// invalid parameter total if the option runner.WithRunDuration is set
func New(host string, protoFile string, methodName string, req proto.Message, dependentProtoFilePath []string, total int, options ...runner.Option) (Runner, error) {
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
		protoFile:              protoFile,
		packageName:            packageName,
		serviceName:            serviceName,
		methodName:             mName,
		methodRequest:          req,
		rpcServerHost:          host,
		total:                  uint(total),
		dependentProtoFilePath: dependentProtoFilePath,
		options:                options,
	}, nil
}

// Run operational performance benchmarking
func (b *bench) Run() error {
	callMethod := fmt.Sprintf("%s.%s/%s", b.packageName, b.serviceName, b.methodName)

	data, err := proto.Marshal(b.methodRequest)
	if err != nil {
		return err
	}

	opts := []runner.Option{
		runner.WithTotalRequests(b.total),
		runner.WithProtoFile(b.protoFile, b.dependentProtoFilePath),
		runner.WithBinaryData(data),
		runner.WithInsecure(true),
		// more parameter settings https://github.com/bojand/ghz/blob/master/runner/options.go#L41
		// example settings: https://github.com/bojand/ghz/blob/master/runner/options_test.go#L79
	}
	opts = append(opts, b.options...)

	report, err := runner.Run(callMethod, b.rpcServerHost, opts...)
	if err != nil {
		return err
	}

	return b.saveReport(callMethod, report)
}

func (b *bench) saveReport(callMethod string, report *runner.Report) error {
	// specify the output path
	outDir := os.TempDir() + string(os.PathSeparator) + "sponge_grpc_benchmark"
	_ = os.MkdirAll(outDir, 0777)
	outputFile := fmt.Sprintf("%sreport_%s.html", outDir+string(os.PathSeparator), b.methodName)
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

	fmt.Printf("\nperformance testing of api '%s' is now complete, copy the report file path to your browser to view,\nreport file: %s\n\n", callMethod, outputFile)
	return rp.Print("html")
}
