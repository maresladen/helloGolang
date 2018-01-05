package utils

import (
	"io"
	"io/ioutil"
	"os"
)

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
	if checkFileIsExist(strFileName) {
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
	if checkFileIsExist(strFileName) {
		file, _ := os.OpenFile(strFileName, os.O_APPEND, 0666)
		defer file.Close()
		file.Write(b)
	} else {
		file, _ := os.Create(strFileName)

		defer file.Close()

		file.Write(b)
	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
