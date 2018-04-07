package codeHunter

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// func EnterPoint() {

// }

//Base64T 测试代码
func Base64T(name string) {
	ff, _ := ioutil.ReadFile(name) //我还是喜欢用这个快速读文件
	sEnc := base64.StdEncoding.EncodeToString(ff)
	ioutil.WriteFile(name+".txt", []byte(sEnc), 0666) //直接写入到文件就ok完活了。
}

//BaseBack 还原
func BaseBack(name string) {
	ff, _ := ioutil.ReadFile(name)
	ddd, _ := base64.StdEncoding.DecodeString(string(ff)) //成图片文件并把文件写入到buffer
	newName := strings.Replace(name, ".txt", "", 1)
	ioutil.WriteFile(newName, ddd, 0666) //buffer输出到jpg文件中（不做处理，直接写到文件）
}

//WalkDir 获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		//if err != nil { //忽略错误
		// return err
		//}

		if fi.IsDir() { // 忽略目录
			return nil
		}
		if suffix != "" {
			if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
				files = append(files, filename)
			}
		} else {
			files = append(files, filename)
		}

		return nil
	})

	return files, err
}
