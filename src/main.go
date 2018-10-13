package main

import (
	// "fileUpload"
	"ES63Import"
	"fmt"
)

func main() {
	// doMain()
	doTest()
}

func doTest() {

	// fmt.Println(123)
	// // time.Sleep(1 * time.Second)
	// time.Sleep(100 * time.Millisecond)
	// fmt.Println(321)
	//otto test

	// do_test.TestJsVM()

	// do_test.HttpTest()
	// ESImport.GetWorkflowData()
	// ESImport.GetPolicyData()
	// ESImport.GetParentID()

	//----------------------------------------------------------------------------------------------------
	var doType string
	fmt.Println("请输入类型，1: 迁移索引  2: 从源索引下载  3: 注入到目标索引")
	fmt.Scanln(&doType)
	if doType == "1" {
		ES63Import.EsDataTrans()
	} else if doType == "2" {
		ES63Import.DownLoadData()
	} else if doType == "3" {
		ES63Import.UploadEsData()
	}

	//----------------------------------------------------------------------------------------------------

	// dbUtils.Test()

	// temp := html[strings.LastIndex(html, "/")+1 : len(html)]
	// fmt.Println(temp)

	// a := codeHunter.GetFileName("http://repo.ebaotech.com/artifactory/repo/com/ebao/unicorn/unicorn-api-gateway/4.3.0B10/unicorn-api-gateway-4.3.0B10-sources.jar")

	// codeHunter.NewFunc()

	// codeHunter.DownloadFiles(html)

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

	// youdaoTranslate.InitTrans()
	// var strTrans = `测试内容和加密`
	// url := youdaoTranslate.GenerateURL(strTrans)
	// youdaoTranslate.TranslateTextByBaidu(url)

	//------------------测试方法 导入ES数据，自动生成---------------------
	// ESImport.DoImportForTest()
	//------------------测试方法 导入ES数据，自动生成---------------------

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

	//-----------------------大数字计算-----------------------
	// num := 18000000
	// tStart := time.Now()
	// a, n := utils.MaxBigint(num)
	// fmt.Println("字符串长度:", len(a))
	// // b := utils.Md5String(a)
	// elapsed := time.Since(tStart)
	// fmt.Println("总共耗时:", elapsed)
	// // fmt.Println("MD5:", b)
	// fmt.Println("质数总数", n)
	// fmt.Println("最大数", a)
	//-----------------------大数字计算-----------------------

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

	// ESImport.ImporterByText()

	//todoCode Example-------------------------------

	// ch := make(chan int, 1)
	// y := 1
	// for x := 0; x < 3; x++ {
	// 	go test(x, y, add, ch)
	// }
	// for ele := range ch {
	// 	if len(ch) == 0 {
	// 		close(ch)
	// 	}
	// 	fmt.Println(ele)
	// }

	//todoCode Example-------------------------------
}

//-----------todo code-----------------
// var z = 0

// type Callback func(x, y int, c chan int)

// //提供一个接口，让外部去实现
// func test(x, y int, callback Callback, c chan int) {
// 	if x > 1 {
// 		callback(x, y, c)
// 	}
// }

// func add(x, y int, c chan int) {
// 	z = x + y
// 	c <- z
// }

//-----------todo code-----------------

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
