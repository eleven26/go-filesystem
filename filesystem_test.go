package fs

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"os"
	"testing"
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
