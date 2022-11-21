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

// SetExecutorPath setting the executor
func SetExecutorPath(path string) {
	executor = path
}

// Exec suitable for executing a single non-blocking command, outputting standard and error logs, but the log output is not real time,
// Note: If the execution of a command blocks permanently, it can cause a concurrent leak.
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

// Run execute the command, you can actively end the command, the execution results are returned in real time in Result.StdOut
func Run(ctx context.Context, command string) *Result {
	result := &Result{StdOut: make(chan string), Err: error(nil)}

	go func() {
		defer func() { close(result.StdOut) }() // execution complete, channel closed

		cmd := exec.CommandContext(ctx, executor, "-c", command)
		handleExec(ctx, cmd, result)
	}()

	return result
}

// Result of the execution of the command
type Result struct {
	StdOut chan string
	Err    error // If nil after the command is executed, the command is executed successfully
}

func handleExec(ctx context.Context, cmd *exec.Cmd, result *Result) {
	result.StdOut <- strings.Join(cmd.Args, " ") + "\n"

	stdout, stderr, err := getCmdReader(cmd)
	if err != nil {
		result.Err = err
		return
	}

	reader := bufio.NewReader(stdout)
	// reads each line in real time
	line := ""
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) { // determine if it has been read
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

	// capture error logs
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
