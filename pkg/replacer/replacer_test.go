package replacer

import (
	"embed"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
			subFiles := []string{"testDir/foo.txt"}
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
			r.SetSubDirsAndFiles(subDirs, subFiles...)
			r.SetIgnoreSubDirs(ignoreDirs...)
			r.SetIgnoreSubFiles(ignoreFiles...)
			r.SetReplacementFields(fields)
			_ = r.SetOutputDir(fmt.Sprintf("%s/replacer_test/%s_%s",
				os.TempDir(), tt.name, time.Now().Format("150405")))
			_, err := r.ReadFile("replace.txt")
			assert.NoError(t, err)
			err = r.SaveFiles()
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
	r.SetIgnoreSubFiles()
	r.SetSubDirsAndFiles(nil)
	err = r.SetOutputDir("/tmp/yourServerName")
	assert.NoError(t, err)
	path := r.GetSourcePath()
	assert.NotEmpty(t, path)

	r = &replacerInfo{}
	err = r.SaveFiles()
	assert.NoError(t, err)
}
