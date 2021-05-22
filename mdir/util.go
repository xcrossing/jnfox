package mdir

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func PathOfName(name string, lens ...int) (string, error) {
	str := md5Str(name)
	dirs, err := splitStr(str, lens...)
	if err != nil {
		return "", err
	}
	path := filepath.Join(dirs...)
	return path, err
}

func md5Str(fileName string) string {
	hash := md5.Sum([]byte(fileName))
	return fmt.Sprintf("%x", hash)
}

func splitStr(str string, lens ...int) ([]string, error) {
	start := 0
	limit := len(str)
	slice := make([]string, 0, len(lens))

	for _, l := range lens {
		end := start + l
		if end > limit {
			return nil, errors.New("1")
		}
		subStr := str[start : start+l]
		slice = append(slice, subStr)
		start = end
	}

	return slice, nil
}

func mkdirs(force bool, dirs ...string) string {
	path := filepath.Join(dirs...)
	if force {
		os.MkdirAll(path, os.ModePerm)
	}
	return path
}

type baseNameNoExtSet map[string]void

var baseNameNoExtSetMember void

type fileList []*fileInfo
type fileInfo struct {
	oldPath       string
	baseName      string
	baseNameNoExt string
	newPath       string
}

func listFiles(dir string) (fileList, error) {
	npm := make(baseNameNoExtSet)
	list := make(fileList, 0, 0)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		baseName := filepath.Base(info.Name())
		baseNameNoExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))
		if _, exists := npm[baseNameNoExt]; exists {
			return errors.New("duplicate " + baseNameNoExt)
		}
		npm[baseNameNoExt] = baseNameNoExtSetMember
		list = append(list, &fileInfo{oldPath: path, baseName: baseName, baseNameNoExt: baseNameNoExt})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return list, nil
}
