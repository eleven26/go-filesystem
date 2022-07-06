package fs

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExists(t *testing.T) {
	ok, err := Exists("filesystem.go")
	assert.True(t, ok)
	assert.Nil(t, err)

	ok, err = Exists("foo.txt")
	assert.False(t, ok)
	assert.True(t, errors.Is(err, os.ErrNotExist))
	assert.True(t, errors.Is(err, fs.ErrNotExist))
}

func TestGet(t *testing.T) {
	bs, err := Get("filesystem1.go")
	assert.Empty(t, bs)
	assert.True(t, errors.Is(err, os.ErrNotExist))
	assert.True(t, errors.Is(err, fs.ErrNotExist))

	bs, _ = Get("test/foo.txt")
	assert.Equal(t, bs, []byte("abc"))
}

func TestGetString(t *testing.T) {
	s, _ := GetString("test/foo.txt")
	assert.Equal(t, s, "abc")
}

func TestPut(t *testing.T) {
	path := "test/tmp/foo.txt"
	err := Put(path, []byte("abc"))
	defer func(name string) {
		_ = os.Remove(name)
	}(path)

	assert.Equal(t, err, nil)

	s, _ := GetString(path)
	assert.Equal(t, s, "abc")
}

func TestPutString(t *testing.T) {
	path := "test/tmp/foo.txt"
	err := PutString(path, "abc")
	defer func(name string) {
		_ = os.Remove(name)
	}(path)

	assert.Equal(t, err, nil)

	s, _ := GetString(path)
	assert.Equal(t, s, "abc")
}

func TestAppend(t *testing.T) {
	path := "test/tmp/foo.txt"
	PutString(path, "abc")
	defer func(name string) {
		_ = os.Remove(name)
	}(path)

	Append(path, []byte("aaa"))
	s, _ := GetString(path)
	assert.Equal(t, s, "abcaaa")
}

func TestPrepend(t *testing.T) {
	path := "test/tmp/foo.txt"
	PutString(path, "abc")
	defer func(name string) {
		_ = os.Remove(name)
	}(path)

	Prepend(path, []byte("aaa"))
	s, _ := GetString(path)
	assert.Equal(t, s, "aaaabc")
}

func TestChmod(t *testing.T) {
	path := "test/tmp/foo.txt"
	PutString(path, "abc")
	defer func(name string) {
		_ = os.Remove(name)
	}(path)

	Chmod(path, 0600)
	f, _ := os.Stat(path)
	assert.Equal(t, "600", fmt.Sprintf("%o", f.Mode().Perm()))
}

func TestDelete(t *testing.T) {
	path1 := "test/tmp/foo1.txt"
	path2 := "test/tmp/foo2.txt"
	PutString(path1, "abc")
	PutString(path2, "def")

	Delete(path1, path2)

	exists, err := Exists(path1)
	assert.False(t, exists)
	assert.True(t, errors.Is(err, os.ErrNotExist))

	exists, err = Exists(path2)
	assert.False(t, exists)
	assert.True(t, errors.Is(err, os.ErrNotExist))
}

func TestMove(t *testing.T) {
	path1 := "test/tmp/foo1.txt"
	path2 := "test/tmp/foo2.txt"
	defer func(name string) {
		_ = os.Remove(name)
	}(path2)

	PutString(path1, "abc")

	Move(path1, path2)

	exists, err := Exists(path1)
	assert.False(t, exists)
	assert.True(t, errors.Is(err, os.ErrNotExist))

	exists, err = Exists(path2)
	assert.True(t, exists)
	assert.Nil(t, err)
}

func TestCopy(t *testing.T) {
	path1 := "test/tmp/foo1.txt"
	path2 := "test/tmp/foo2.txt"
	defer func(name string) {
		_ = os.Remove(name)
	}(path1)
	defer func(name string) {
		_ = os.Remove(name)
	}(path2)

	PutString(path1, "abc")
	Copy(path1, path2)

	exists, err := Exists(path1)
	assert.True(t, exists)
	assert.Nil(t, err)

	exists, err = Exists(path2)
	assert.True(t, exists)
	assert.Nil(t, err)
}

func TestLink(t *testing.T) {
	path1 := "test/tmp/foo1.txt"
	path2 := "test/tmp/foo2.txt"
	defer func(name string) {
		_ = os.Remove(name)
	}(path1)
	defer func(name string) {
		_ = os.Remove(name)
	}(path2)

	PutString(path1, "abc")
	p1, _ := filepath.Abs(path1)
	p2, _ := filepath.Abs(path2)
	Link(p1, p2)

	s, _ := GetString(path2)
	assert.Equal(t, s, "abc")
}

func TestName(t *testing.T) {
	path := "/a/b/c.txt"
	assert.Equal(t, Name(path), "c.txt")
}

func TestBasename(t *testing.T) {
	path := "/a/b/c.txt"
	assert.Equal(t, Basename(path), "c")
}

func TestDirname(t *testing.T) {
	path := "/a/b/c.txt"
	assert.Equal(t, Dirname(path), "/a/b")
}

func TestExtension(t *testing.T) {
	path := "/a/b/c.txt"
	assert.Equal(t, Extension(path), "txt")

	assert.Equal(t, Extension("a"), "")
}

func TestSize(t *testing.T) {
	path := "test/foo.txt"

	s, _ := Size(path)
	assert.Equal(t, int64(3), s)
}

func TestLastModified(t *testing.T) {
	path := "test/foo.txt"

	ti, _ := LastModified(path)
	assert.IsType(t, time.Time{}, ti)
}

func TestIsDirectory(t *testing.T) {
	path := "test/foo.txt"
	b, _ := IsDirectory(path)
	assert.False(t, b)

	path = "test"
	b, _ = IsDirectory(path)
	assert.True(t, b)
}

func TestIsReadable(t *testing.T) {
	path := "test/foo.txt"
	b, _ := IsReadable(path)
	assert.True(t, b)
}

func TestIsWritable(t *testing.T) {
	path := "test/foo.txt"
	b, _ := IsWritable(path)
	assert.True(t, b)
}

func TestIsFile(t *testing.T) {
	path := "test/foo.txt"
	b, _ := IsFile(path)
	assert.True(t, b)
}

func TestFiles(t *testing.T) {
	path := "test"
	files, _ := Files(path)
	assert.Equal(t, 1, len(files))
	assert.Equal(t, "foo.txt", files[0])
}
