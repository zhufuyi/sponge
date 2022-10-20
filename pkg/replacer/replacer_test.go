package replacer

import (
	"embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

//go:embed testDir
var fs embed.FS

func TestNewWithFS(t *testing.T) {
	type args struct {
		fn func() Replacer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "New",
			args: args{
				fn: func() Replacer {
					replacer, err := New("testDir")
					if err != nil {
						panic(err)
					}
					return replacer
				},
			},
			wantErr: false,
		},

		{
			name: "NewFS",
			args: args{
				fn: func() Replacer {
					replacer, err := NewFS("testDir", fs)
					if err != nil {
						panic(err)
					}
					return replacer
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.args.fn()

			subDirs := []string{"testDir/replace"}
			ignoreDirs := []string{"testDir/ignore"}
			ignoreFiles := []string{"test.txt"}
			fields := []Field{
				{
					Old: "1234",
					New: "....",
				},
				{
					Old:             "abcdef",
					New:             "hello_",
					IsCaseSensitive: true,
				},
			}
			r.SetSubDirs(subDirs...)         // 只处理指定子目录，为空时表示指定全部文件
			r.SetIgnoreFiles(ignoreDirs...)  // 忽略替换目录
			r.SetIgnoreFiles(ignoreFiles...) // 忽略替换文件
			r.SetReplacementFields(fields)   // 设置替换文本
			_ = r.SetOutputDir(fmt.Sprintf("%s/replacer_test/%s_%s",
				os.TempDir(), tt.name, time.Now().Format("150405"))) // 设置输出目录和名称
			_, err := r.ReadFile("replace.txt")
			assert.NoError(t, err)
			err = r.SaveFiles() // 保存替换后文件
			if (err != nil) != tt.wantErr {
				t.Logf("SaveFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("save files successfully, out = %s", r.GetOutputDir())
		})
	}
}

func TestReplacerError(t *testing.T) {
	_, err := New("/notfound")
	assert.Error(t, err)
	_, err = NewFS("/notfound", embed.FS{})
	assert.Error(t, err)

	r, err := New("testDir")
	assert.NoError(t, err)
	r.SetIgnoreFiles()
	r.SetSubDirs()
	err = r.SetOutputDir("/tmp/yourServerName")
	assert.NoError(t, err)
	path := r.GetSourcePath()
	assert.NotEmpty(t, path)

	r = &replacerInfo{}
	err = r.SaveFiles()
	assert.NoError(t, err)
}
