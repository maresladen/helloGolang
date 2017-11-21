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

	fmt.Println("测试方法")
	// 处理String内容
	// stringEditor.DoStringEditor()
	// 有道翻译
	// youdaoTranslate.TranslateText()
}

func doMain() {

	// fmt.Println("主方法")

	// 上传服务器
	// fileUpload.FileUpload()

	// ES导入程序
	// ESImport.Importer()

	tStart := time.Now()
	a, n := utils.MaxBigint(1000000)
	b := utils.Md5fun(a)
	elapsed := time.Since(tStart)
	fmt.Println("总共耗时:", elapsed)
	fmt.Println(b)
	fmt.Println(n)

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
