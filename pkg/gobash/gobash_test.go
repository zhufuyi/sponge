package gobash

import (
	"context"
	"fmt"
	"reflect"
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
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*500)
		result := Run(ctx, cmd)
		for v := range result.StdOut { // Real-time output of logs and error messages
			t.Logf(v)
		}
		if result.Err != nil {
			t.Logf("exec command failed, %v", result.Err)
		}
	}
}
func TestRunC(t *testing.T) {
	cmds := map[string][]string{
		"ping": []string{"www.baidu.com"},
		"pwd":  []string{},
		"go":   []string{"env"},
	}

	for cmd, args := range cmds {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		result := RunC(ctx, cmd, args...)
		for v := range result.StdOut { // Real-time output of logs and error messages
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

func TestExecC(t *testing.T) {
	cmds := map[string][]string{
		"pwd":    []string{},
		"go":     []string{"env"},
		"sponge": []string{"-h"},
	}

	for cmd, args := range cmds {
		out, err := ExecC(cmd, args...)
		if err != nil {
			t.Logf("exec command[%s] failed, %v\n", cmd, err)
			continue
		}
		t.Logf("%s\n", out)
	}
}

func reflectObject(obj interface{}) {
	if obj == nil {
		return
	}

	to := reflect.TypeOf(obj)
	vo := reflect.ValueOf(obj)

	// if not a struct and then exit
	if to.Kind() == reflect.Ptr {
		if vo.IsNil() {
			return
		}
		if to.Elem().Kind() != reflect.Struct {
			fmt.Println("unsupported type: " + to.Kind().String())
			return
		}
		// pointer to struct
		to = to.Elem()
		vo = vo.Elem()
	} else if to.Kind() != reflect.Struct {
		fmt.Println("unsupported type: " + to.Kind().String())
		return
	}

	handleStruct(to, vo)
}

func handleStruct(to reflect.Type, vo reflect.Value) {
	num := to.NumField()
	for i := 0; i < num; i++ {
		field := to.Field(i)
		fieldName := field.Name
		fieldValue := vo.FieldByName(fieldName)

		if fieldValue.Kind() == reflect.Struct {
			subTO := field.Type
			subVO := fieldValue
			fmt.Println()
			fmt.Printf("[field %d] name=%s, type=%s, typeName=%s\n", i+1, fieldName, fieldValue.Kind().String(), field.Type.String())
			handleStruct(subTO, subVO)
			fmt.Println()
		} else if fieldValue.Kind() == reflect.Ptr {
			subTO := field.Type.Elem()
			subVO := fieldValue.Elem()
			fmt.Println()
			fmt.Printf("[field %d] name=%s, type=%s, typeName=%s\n", i+1, fieldName, fieldValue.Kind().String(), field.Type.String())
			handleStruct(subTO, subVO)
			fmt.Println()
		} else {
			fmt.Printf("[field %d] name=%s, type=%s, typeName=%s, value=%v\n", i+1, fieldName, fieldValue.Kind().String(), field.Type.String(), fieldValue)
		}
	}
	handleStructMethods(to)
}

func handleStructMethods(to reflect.Type) {
	num := to.NumMethod()
	for i := 0; i < num; i++ {
		method := to.Method(i)
		fmt.Printf("[method %d] name=%s\n", i+1, method.Name)
		for j := 0; j < method.Type.NumIn(); j++ {
			fmt.Printf("[method parameter %d] type=%s\n", j+1, method.Type.In(j).Name())
		}
		for k := 0; k < method.Type.NumOut(); k++ {
			fmt.Printf("[method return %d] type=%s\n", k+1, method.Type.Out(k).Name())
		}
	}
}

type secStruct struct {
	Cnt []int64
	Sec string
}

type ThiStruct struct {
	Cnt []int64
	Sec string
}

func (s secStruct) Receive(msg string) {
	fmt.Println(msg)
}

type myStruct struct {
	Num        int    `json:"num_json" orm:"column:num_orm"`
	Desc       string `json:"desc_json" orm:"column:desc_orm"`
	Child      secStruct
	Child2     *secStruct
	IInterface interface{}
}

func (s myStruct) Send(msg string) (int, error) {
	fmt.Println(msg)
	return 10, nil
}

func TestExecC2(t *testing.T) {

	s := myStruct{
		Num:        100,
		Desc:       "这是描述",
		Child:      secStruct{[]int64{1, 2, 3}, "子字段"},
		Child2:     &secStruct{[]int64{4, 5, 6}, "子字段2"},
		IInterface: nil,
	}

	reflectObject(&s)
}
