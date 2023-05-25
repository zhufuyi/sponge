package jy2struct

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvert(t *testing.T) {
	type args struct {
		args *Args
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "json to struct from data",
			args: args{args: &Args{
				Data:   `{"name":"foo","age":11}`,
				Format: "json",
			}},
			wantErr: false,
		},
		{
			name: "yaml to struct from data",
			args: args{args: &Args{
				Data: `name: "foo"
age: 10`,
				Format: "yaml",
			}},
			wantErr: false,
		},
		{
			name: "json to struct from file",
			args: args{args: &Args{
				InputFile: "test.json",
				Format:    "json",
				SubStruct: true,
				Tags:      "gorm",
			}},
			wantErr: false,
		},
		{
			name: "yaml to struct from file",
			args: args{args: &Args{
				InputFile: "test.yaml",
				Format:    "yaml",
				SubStruct: true,
			}},
			wantErr: false,
		},
		{
			name: "json to slice from data",
			args: args{args: &Args{
				Data:   `[{"name":"foo","age":11},{"name":"foo2","age":22}]`,
				Format: "json",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}

	// test Convert error
	arg := &Args{Format: "unknown"}
	_, err := Convert(arg)
	assert.Error(t, err)
	arg = &Args{Format: "yaml", InputFile: "notfound.yaml"}
	_, err = Convert(arg)
	assert.Error(t, err)
}
