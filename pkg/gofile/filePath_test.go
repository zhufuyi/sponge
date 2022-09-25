package gofile

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestIsExists(t *testing.T) {
	ok := IsExists("/tmp/test")
	if !ok {
		t.Log("not exists")
	}
}

func TestGetRunPath(t *testing.T) {
	t.Log(GetRunPath())
}

func TestListFiles(t *testing.T) {
	dir := "."

	t.Run("all files", func(t *testing.T) {
		files, err := ListFiles(dir)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(strings.Join(files, "\n"))
	})

	t.Run("prefix name", func(t *testing.T) {
		files, err := ListFiles(dir, WithPrefix("READ"))
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(strings.Join(files, "\n"))
	})

	t.Run("suffix name", func(t *testing.T) {
		files, err := ListFiles(dir, WithSuffix(".go"))
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(strings.Join(files, "\n"))
	})

	t.Run("contain name", func(t *testing.T) {
		files, err := ListFiles(dir, WithContain("file"))
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(strings.Join(files, "\n"))
	})
}

func TestListDirsAndFiles(t *testing.T) {
	df, err := ListDirsAndFiles(".")
	if err != nil {
		t.Fatal(err)
	}
	for dir, files := range df {
		t.Log(dir, strings.Join(files, "\n"))
	}
}

func TestGetFilename(t *testing.T) {
	name := GetFilename("./README.md")
	assert.Equal(t, "README.md", name)
}

func TestGetPathDelimiter(t *testing.T) {
	d := GetPathDelimiter()
	t.Log(d)
}
