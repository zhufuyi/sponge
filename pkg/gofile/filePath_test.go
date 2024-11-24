package gofile

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsExists(t *testing.T) {
	ok := IsExists("/tmp/tmp/tmp")
	assert.Equal(t, false, ok)
	ok = IsExists("README.md")
	assert.Equal(t, true, ok)
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

	t.Run("no filepath abs", func(t *testing.T) {
		files, err := ListFiles(dir, WithSuffix(".go"), WithNoAbsolutePath())
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

	name = GetFileSuffixName("./README.md")
	assert.Equal(t, ".md", name)

	name = GetDir("gofile/README.md")
	assert.Equal(t, "gofile", name)

	name = GetSuffixDir("gofile/")
	assert.Equal(t, "gofile", name)

	name = GetFileDir("gofile/README.md")
	assert.Equal(t, "gofile/", name)

	err := CreateDir("./")
	assert.Nil(t, err)
	notfoundDir := os.TempDir() + "/notfoundDir"
	err = CreateDir(notfoundDir)
	assert.Nil(t, err)
	time.Sleep(time.Millisecond * 100)
	_ = os.RemoveAll(notfoundDir)

	name = GetFilenameWithoutSuffix("./README.md")
	assert.Equal(t, "README", name)
}

func TestGetPathDelimiter(t *testing.T) {
	d := GetPathDelimiter()
	t.Log(d)
}

func TestNotMatch(t *testing.T) {
	fn := matchPrefix("")
	assert.Equal(t, false, fn("."))

	fn = matchContain("")
	assert.Equal(t, false, fn("."))

	fn = matchSuffix("")
	assert.NotNil(t, fn)
}

func TestIsWindows(t *testing.T) {
	t.Log(IsWindows())
}

func TestErrorPath(t *testing.T) {
	dir := "/notfound"

	_, err := ListFiles(dir)
	assert.Error(t, err)

	_, err = ListDirsAndFiles(dir)
	assert.Error(t, err)

	err = walkDirWithFilter(dir, nil, nil)
	assert.Error(t, err)

	err = walkDir(dir, nil)
	assert.Error(t, err)

	err = walkDir2(dir, nil, nil)
	assert.Error(t, err)
}

func TestFuzzyMatchFiles(t *testing.T) {
	files := FuzzyMatchFiles("./README.md")
	assert.Equal(t, 1, len(files))

	files = FuzzyMatchFiles("./*_test.go")
	assert.Equal(t, 2, len(files))
}

func TestJoin(t *testing.T) {
	elements := []string{"a/b", "c"}
	path := Join(elements...)
	t.Log(path)

	elements = []string{"a\\b", "c"}
	path = Join(elements...)
	t.Log(path)
}

func TestListDirs(t *testing.T) {
	dir := ".."
	dirs, err := ListDirs(dir)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(FilterDirs(dirs, WithSuffix(".txt")))
	t.Log(FilterDirs(dirs, WithPrefix("query")))
	t.Log(FilterDirs(dirs, WithContain("auth")))
}

func TestListSubDirs(t *testing.T) {
	dir := ".."
	dirs, err := ListSubDirs(dir, "gin")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dirs)
}
