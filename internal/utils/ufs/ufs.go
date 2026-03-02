package ufs

import (
	"io"
	"os"
	"path"
	"slices"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}

	return true
}

func WalkListDir(dir string, exclude []string) []string {
	paths := []string{}

	entries, err := os.ReadDir(dir)

	if slices.Contains(exclude, dir) {
		return paths
	}

	if err != nil {
		println("Failed to read template directory.")
		os.Exit(1)
	}

	for _, entry := range entries {
		current := path.Join(dir, entry.Name())

		if slices.Contains(exclude, current) || slices.Contains(exclude, entry.Name()) {
			continue
		}

		if entry.IsDir() {
			childrens := WalkListDir(current, exclude)
			paths = append(paths, childrens...)
		} else {
			paths = append(paths, current)
		}
	}

	return paths
}

func Copy(src string, dst string) (bool, error) {
	if !Exists(src) {
		return false, os.ErrNotExist
	}

	if Exists(dst) {
		return false, os.ErrExist
	}

	dir := path.Dir(dst)

	if !Exists(dir) {
		err := MkDir(dir, true)

		if err != nil {

			return false, err
		}
	}

	in, err := os.Open(src)

	if err != nil {
		return false, err
	}

	defer in.Close()

	out, err := os.Create(dst)

	if err != nil {
		return false, err
	}

	defer out.Close()

	_, err = io.Copy(in, out)

	if err != nil {
		return false, err
	}

	err = out.Sync()

	if err != nil {
		return false, err
	}

	return true, nil
}

func CopyAll(src string, dst string) (bool, error) {
	paths := WalkListDir(src, []string{})
	for _, p := range paths {
		target := path.Join(dst, strings.TrimPrefix(p, src))

		_, err := Copy(p, target)

		if err != nil {
			println(err.Error())

			return false, err
		}
	}

	return true, nil
}

func MkDir(dir string, all bool) error {
	if Exists(dir) {
		return os.ErrExist
	}

	if all {
		return os.MkdirAll(dir, 0755)
	}

	return os.Mkdir(dir, 0755)
}

func ListDir(dir string) []string {
	entries, err := os.ReadDir(dir)

	if err != nil {
		println("Failed to read template directory.")
		os.Exit(1)
	}

	paths := []string{}

	for _, entry := range entries {
		paths = append(paths, path.Join(dir, entry.Name()))
	}

	return paths
}

func WriteFileByByte(path string, content []byte) error {
	err := os.WriteFile(path, content, 0644)

	if err != nil {
		return err
	}

	return nil
}

func ReadFileByByte(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func WriteFile(path string, content string) error {
	return WriteFileByByte(path, []byte(content))
}

func ReadFile(path string) (string, error) {
	bytes, err := ReadFileByByte(path)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func Rename(oldPath string, newPath string) error {
	return os.Rename(oldPath, newPath)
}
