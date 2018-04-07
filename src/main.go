package main

import (
	"fmt"
	"time"
	"utils"
	"youdaoTranslate"
)

func main() {
	doMain()
	// doTest()
}

func doTest() {

	// path, _ := utils.GetCurrentPath()
	// temp, _ := codeHunter.WalkDir(path, ".txt")
	// temp, _ := codeHunter.WalkDir(path)

	// for _, zipFile := range temp {
	// 	codeHunter.Base64T(zipFile)
	// }

	// for _, txtFile := range temp {
	// 	codeHunter.BaseBack(txtFile)
	// }

	// codeHunter.BaseBack("temp.txt")
	// fmt.Println("测试方法")

	// abc := utils.Md5File("hello")
	// fmt.Println(abc)

	// 处理String内容
	// stringEditor.DoStringEditor()
	// 有道翻译
	// youdaoTranslate.TranslateText()

	// -------------------------------------------------------------------------------
	// 百度翻译

	var strTrans = `如何在一次请求中翻译多个单词或者多段文本` + "\n" + `为什么我的请求总是返回错误码54001`
	youdaoTranslate.TranslateTextByBaidu(strTrans)
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
	num := 180000
	tStart := time.Now()
	a, n := utils.MaxBigint(num)
	fmt.Println("字符串长度:", len(a))
	b := utils.Md5String(a)
	elapsed := time.Since(tStart)
	fmt.Println("总共耗时:", elapsed)
	fmt.Println("MD5:", b)
	fmt.Println("质数总数", n)
	fmt.Println("数字", a)

	// c := utils.DecimalToAny(9999999999, 76)
	// fmt.Println("76进制", c)

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
