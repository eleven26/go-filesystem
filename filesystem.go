package fs

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var Separator = fmt.Sprintf("%c", os.PathSeparator)

func Exists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else {
		// file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		return false, err
	}
}

func Get(path string) (b []byte, err error) {
	return os.ReadFile(path)
}

func GetString(path string) (content string, err error) {
	b, err := Get(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func Put(path string, content []byte) error {
	return ioutil.WriteFile(path, content, 0644)
}

func PutString(path string, content string) error {
	return Put(path, []byte(content))
}

func Append(path string, content []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_, err = f.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func Prepend(path string, content []byte) error {
	s, err := Get(path)
	if err != nil {
		return err
	}

	content = append(content, s...)
	err = Put(path, content)
	return err
}

func Chmod(path string, mode fs.FileMode) error {
	return os.Chmod(path, mode)
}

func Delete(paths ...string) error {
	for _, path := range paths {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func Move(from string, to string) error {
	return os.Rename(from, to)
}

func Copy(src, dst string) error {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	// Write data to dst
	err = ioutil.WriteFile(dst, data, 0644)
	return err
}

func Link(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

func Name(path string) string {
	return filepath.Base(path)
}

func Basename(path string) string {
	base := Name(path)

	return base[0:strings.Index(base, ".")]
}

func Dirname(path string) string {
	return filepath.Dir(path)
}

func Extension(path string) string {
	base := Name(path)

	if !strings.Contains(base, ".") {
		return ""
	}

	return base[strings.Index(base, ".")+1:]
}

func Size(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	size := fi.Size()
	return size, nil
}

func LastModified(path string) (t time.Time, err error) {
	fi, err := os.Stat(path)

	if err != nil {
		return
	}

	t = fi.ModTime()
	return
}

func IsDirectory(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fi.IsDir(), nil
}

func IsReadable(path string) (bool, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	defer file.Close()

	return err == nil, err
}

func IsWritable(path string) (bool, error) {
	file, err := os.OpenFile(path, os.O_WRONLY, 0666)
	defer file.Close()

	return err == nil, err
}

func IsFile(path string) (bool, error) {
	return Exists(path)
}

// Files List the files under the folder, excluding directories.
func Files(dir string) (files []string, err error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, f := range fileInfos {
		if !f.IsDir() {
			files = append(files, f.Name())
		}
	}

	return
}

func AllFiles(dir string) (files []string, err error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, f := range fileInfos {
		if f.IsDir() {
			allFiles, e := AllFiles(strings.Join([]string{dir, f.Name()}, Separator))
			if e != nil {
				err = e
				return
			}

			newFiles := make([]string, len(allFiles))
			for i, file := range allFiles {
				newFiles[i] = strings.Join([]string{f.Name(), file}, Separator)
			}

			files = append(files, newFiles...)
		} else {
			files = append(files, f.Name())
		}
	}

	return
}

func Directories(dir string) (dirs []string, err error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, f := range fileInfos {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}

	return
}

func MakeDirectory(path string, mode os.FileMode) error {
	return os.Mkdir(path, mode)
}

func MakeDirectories(path string, mode os.FileMode) error {
	return os.MkdirAll(path, mode)
}

func DeleteDirectory(path string) error {
	return os.RemoveAll(path)
}

func MoveDirectory(src, dst string) error {
	return os.Rename(src, dst)
}
