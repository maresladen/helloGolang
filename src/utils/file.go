package utils

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//GetCurrentPath 获取当前执行路径
func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".
			`)
	}
	return string(path[0 : i+1]), nil
}

//ReadAllBytes 读文整个文件
func ReadAllBytes(path string) []byte {
	fi, err := os.Open(path)
	if err != nil {
		println(err)
	}
	result, err := ioutil.ReadAll(fi)

	if err != nil {
		println(err)
	}
	return result
}

//WriteFile 写文件
func WriteFile(strFileName, strDefine string) {
	if CheckFileIsExist(strFileName) {
		file, _ := os.OpenFile(strFileName, os.O_APPEND, 0666)
		defer file.Close()
		io.WriteString(file, strDefine)
	} else {
		file, _ := os.Create(strFileName)

		defer file.Close()

		file.WriteString(strDefine)
	}
}

//WriteBytes 写文件
func WriteBytes(strFileName string, b []byte) {
	if CheckFileIsExist(strFileName) {
		file, _ := os.OpenFile(strFileName, os.O_APPEND, 0666)
		defer file.Close()
		file.Write(b)
	} else {
		file, _ := os.Create(strFileName)

		defer file.Close()

		file.Write(b)
	}
}

//CheckFileIsExist 检查是否存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
