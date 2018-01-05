package main

import (
	"fmt"
	"time"
	"utils"
)

func main() {
	doMain()
	// doTest()
}

func doTest() {

	// fmt.Println("测试方法")

	// abc := utils.Md5File("hello")
	// fmt.Println(abc)

	// 处理String内容
	// stringEditor.DoStringEditor()
	// 有道翻译
	// youdaoTranslate.TranslateText()
}

//export doMain
func doMain() {

	// fmt.Println("主方法")

	// 上传服务器
	// fileUpload.FileUpload()

	// ES导入程序
	// ESImport.ImporterByText()

	//---------------gzip----------------

	// stream := utils.ReadAllBytes("/Users/BetaFun/Downloads/test.txt")
	// utils.DoGzip(stream)

	//---------------gzip----------------

	//大数字计算
	num := 108000
	tStart := time.Now()
	a, n := utils.MaxBigint(num)
	fmt.Println("字符串长度:", len(a))
	b := utils.Md5String(a)
	elapsed := time.Since(tStart)
	fmt.Println("总共耗时:", elapsed)
	fmt.Println("MD5:", b)
	fmt.Println("质数总数", n)

	//---------------zlib----------------
	// //read
	// var in bytes.Buffer
	// // bw := utils.ReadAllBytes("test.txt")
	// bw := []byte(a)
	// w := zlib.NewWriter(&in)
	// w.Write(bw)
	// w.Close()

	// utils.WriteBytes("hello", in.Bytes())

	// //write
	// var out bytes.Buffer
	// br := utils.ReadAllBytes("hello")
	// r, _ := zlib.NewReader(bytes.NewBuffer(br))
	// io.Copy(&out, r)
	// fmt.Println("---输出---")
	// fmt.Println(len(out.String()))

	//---------------zlib----------------

}

// //export sum
// func sum(a, b int) int {
// 	return a + b
// }

// "C"
// //export add
// func add(left, right *C.char) *C.char {
// 	// bytes对应ctypes的c_char_p类型,翻译成C类型就是 char *指针
// 	merge := C.GoString(left) + C.GoString(right)
// 	fmt.Println("---------------------go print---------------------")
// 	fmt.Println(merge)
// 	fmt.Println("---------------------go print---------------------")
// 	return C.CString(merge)
// }
