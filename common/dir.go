package common

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

func GetServiceName() (string, error) {
	cur, err := GetCurPath()
	if err != nil {
		return "", err
	}

	return filepath.Base(cur), nil
}

func GetCurPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return dir, nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(path string) error {
	exist, err := PathExists(path)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	return os.MkdirAll(path, os.ModePerm)
}

func WriteToFile(buffer *bytes.Buffer, outputFile string, force bool) error {
	exists, _ := PathExists(outputFile)
	if exists && !force {
		var override string
		fmt.Fprint(os.Stdout, "文件("+outputFile+")已存在，是否需要覆盖(Y/y,默认不覆盖)? ")
		fmt.Scanln(&override)

		if override != "Y" && override != "y" {
			fmt.Fprintln(os.Stdout, "略过"+outputFile)
			return nil
		}
	}

	f, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(buffer.String())

	return nil
}
