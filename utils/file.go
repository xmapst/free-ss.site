package utils

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
)

func GetSize(f multipart.File) (int, error) {
	context, err := ioutil.ReadAll(f)
	return len(context), err
}

func GetExt(fileName string) string {
	return path.Ext(fileName)
}

func CheckExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

func CheckPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func IsNotExistMKDir(src string) error {
	if notExist := CheckExist(src); notExist == true {
		if err := MKDir(src); err != nil {
			return err
		}
	}

	return nil
}

func MKDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func MustOpen(fileName, filePath string) (*os.File, error) {
	src, err := checkSrc(filePath)
	if err != nil {
		return nil, err
	}

	f, err := Open(filepath.Join(src, fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}

func checkSrc(filePath string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("os.Getwd err : %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		return "", fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = IsNotExistMKDir(src)
	if err != nil {
		return "", fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}
	return src, nil
}

func MustWrite(fileName, filePath string) (*os.File, error) {
	src, err := checkSrc(filePath)
	if err != nil {
		return nil, err
	}
	f, err := Open(src+fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}

	return f, nil
}
