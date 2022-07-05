package fs

import (
	"io/fs"
	"io/ioutil"
	"os"
)

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
