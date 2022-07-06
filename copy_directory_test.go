package fs

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCopyDirectory(t *testing.T) {
	os.Mkdir("test/tmp/a", os.ModePerm)
	PutString("test/tmp/a/b.txt", "abc")
	defer os.RemoveAll("test/tmp/a")

	err := CopyDirectory("test/tmp/a", "test/tmp/b")
	if err != nil {
		t.Fatalf("%#v\n", err)
	}
	defer os.RemoveAll("test/tmp/b")

	ok, err := Exists("test/tmp/b/b.txt")
	assert.Nil(t, err)
	assert.True(t, ok)

	s, err := GetString("test/tmp/b/b.txt")
	if err != nil {
		t.Fatalf("%#v\n", err)
	}
	assert.Equal(t, "abc", s)
}
