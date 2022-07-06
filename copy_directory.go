package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

func CopyDirectory(dir, dest string) error {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	if err := createIfNotExists(dest, 0755); err != nil {
		return err
	}

	for _, entry := range entries {
		src := filepath.Join(dir, entry.Name())
		dst := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(src)
		if err != nil {
			return err
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", src)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := createIfNotExists(dst, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(src, dst); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := copySymLink(src, dst); err != nil {
				return err
			}
		default:
			if err := Copy(src, dst); err != nil {
				return err
			}
		}

		if err := os.Lchown(dst, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		isSymlink := entry.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(dst, entry.Mode()); err != nil {
				return err
			}
		}
	}

	return nil
}

func createIfNotExists(dir string, perm os.FileMode) error {
	_, err := os.Stat(dir)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	if err := MakeDirectories(dir, perm); err != nil {
		return err
	}

	return nil
}

func copySymLink(src, dst string) error {
	link, err := os.Readlink(src)
	if err != nil {
		return err
	}

	return os.Symlink(link, dst)
}
