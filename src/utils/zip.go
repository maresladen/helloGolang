package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

//DoGzip 压缩和解压
func DoGzip(stream []byte) {

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()

	w.Write(stream)

	w.Flush()
	fmt.Println("gzip size:", len(b.Bytes()))
	fmt.Println(b)

	r, _ := gzip.NewReader(&b)
	defer r.Close()
	undatas, _ := ioutil.ReadAll(r)
	fmt.Println("ungzip size:", len(undatas))
}

//Md5String 字符串转换为MD5
func Md5String(s string) string {
	signByte := []byte(s)
	return Md5ByteArr(signByte)
}

//Md5ByteArr 字符数组转换为MD5
func Md5ByteArr(b []byte) string {
	hash := md5.New()
	hash.Write(b)
	return hex.EncodeToString(hash.Sum(nil))
}

//Md5File 通过文件路径生成md5
func Md5File(path string) string {
	file, err := os.Open(path)
	if err == nil {
		hash := md5.New()
		io.Copy(hash, file)
		return hex.EncodeToString(hash.Sum(nil))
	}
	return ""
}
