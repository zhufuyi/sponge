package gobash

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func init() {
	if runtime.GOOS == "windows" {
		SetExecutorPath("D:\\Program Files\\cmder\\vendor\\git-for-windows\\bin\\bash.exe")
	}
}

func TestRun(t *testing.T) {
	cmds := []string{
		"for i in $(seq 1 3); do  exit 1; done",
		"notFoundCommand",
		"pwd",
		"for i in $(seq 1 5); do echo 'test cmd' $i;sleep 0.2; done",
	}

	for _, cmd := range cmds {
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500) // 超时控制
		result := Run(ctx, cmd)                                                   // 执行
		for v := range result.StdOut {                                            // 实时输出日志和错误信息
			t.Logf(v)
		}
		if result.Err != nil {
			t.Logf("exec command failed, %v", result.Err)
		}
	}
}

func TestExec(t *testing.T) {
	cmds := []string{
		"for i in $(seq 1 3); do  exit 1; done",
		"notFoundCommand",
		"pwd",
		"for i in $(seq 1 3); do echo 'test cmd' $i;sleep 0.2; done",
	}

	for _, cmd := range cmds {
		out, err := Exec(cmd)
		if err != nil {
			t.Logf("exec command[%s] failed, %v\n", cmd, err)
			continue
		}
		t.Logf("%s\n", out)
	}
}
