package gobash

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// linux default executor
var executor = "/bin/bash"

// SetExecutorPath 设置执行器
func SetExecutorPath(path string) {
	executor = path
}

// Exec 适合执行单条非阻塞命令，输出标准和错误日志，但日志输出不是实时，
// 注：如果执行命令永久阻塞，会造成协程泄露
func Exec(command string) ([]byte, error) {
	cmd := exec.Command(executor, "-c", command)

	stdout, stderr, err := getCmdReader(cmd)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	bytesErr, err := io.ReadAll(stderr)
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		if len(bytesErr) != 0 {
			return nil, errors.New(string(bytesErr))
		}
		return nil, err
	}

	return bytes, nil
}

// Run 执行命令，可以主动结束命令，执行结果实时返回在Result.StdOut中
func Run(ctx context.Context, command string) *Result {
	result := &Result{StdOut: make(chan string), Err: error(nil)}

	go func() {
		defer func() { close(result.StdOut) }() // 执行完毕，关闭通道

		cmd := exec.CommandContext(ctx, executor, "-c", command)
		handleExec(ctx, cmd, result)
	}()

	return result
}

// Result 执行命令的结果
type Result struct {
	StdOut chan string
	Err    error // 执行完毕命令后，如果为nil，执行命令成功
}

func handleExec(ctx context.Context, cmd *exec.Cmd, result *Result) {
	result.StdOut <- strings.Join(cmd.Args, " ") + "\n"

	stdout, stderr, err := getCmdReader(cmd)
	if err != nil {
		result.Err = err
		return
	}

	reader := bufio.NewReader(stdout)
	// 实时读取每行内容
	line := ""
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) { // 判断是否已经读取完毕
				break
			}
			result.Err = err
			break
		}
		select {
		case result.StdOut <- line:
		case <-ctx.Done():
			result.Err = fmt.Errorf("%v", ctx.Err())
			return
		}
	}

	// 捕获错误日志
	bytesErr, err := io.ReadAll(stderr)
	if err != nil {
		result.Err = err
		return
	}

	err = cmd.Wait()
	if err != nil {
		if len(bytesErr) != 0 {
			result.Err = errors.New(string(bytesErr))
			return
		}
		result.Err = err
	}
}

func getCmdReader(cmd *exec.Cmd) (io.ReadCloser, io.ReadCloser, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	return stdout, stderr, nil
}
