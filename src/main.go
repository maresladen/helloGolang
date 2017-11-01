package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"stringEditor"
	strings "strings"
	"sync"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"
)

type columnConfig struct {
	FileName          string       `json:"FileName"`
	Cloumns           []colmunProp `json:"Cloumns"`
	SliptFileRowCount int          `json:"SliptFileRowCount"`
	ChannelSize       int          `json:"ChannelSize"`
	SubmitURL         string       `json:"SubmitURL"`
	EsIndex           string       `json:"EsIndex"`
	EsType            string       `json:"EsType"`
	LoginUser         string       `json:"LoginUser"`
	LoginPWD          string       `json:"LoginPWD"`
}

type colmunProp struct {
	ColumnName      string
	ColumnIndex     int
	IsAllTextColumn bool
	FieldType       string
	Before          string
	After           string
}

var cConfig columnConfig
var keyMap map[int]colmunProp
var splitFileNameNum = 1
var importAllCount = 0
var wg sync.WaitGroup
var wgPost sync.WaitGroup
var wgPostChild sync.WaitGroup

var wgChild sync.WaitGroup

func main() {
	doMain()
	// doTest()
}

func doTest() {

	stringEditor.DoStringEditor()
	// youdaoTranslate.TranslateText()
}

func doMain() {

	fmt.Println("let's go! ", time.Now().Format("2006-01-02 15:04:05"))
	//读取配置
	readConfig()
	//初始化一个hash表，用于后续处理字符串
	keyMap = make(map[int]colmunProp)
	for _, prop := range cConfig.Cloumns {
		keyMap[prop.ColumnIndex] = prop
	}

	//调用逐行读取文件方法
	readFileByLine()

	//阻塞等待所有组内成员都执行完毕退栈
	wgChild.Wait()
	wg.Wait()
	wgChild.Wait()
	wgPost.Wait()
	fmt.Println("执行完毕 ", time.Now().Format("2006-01-02 15:04:05"))
}

func readConfig() {
	fi, err := os.Open("config.json")
	if err != nil {
		writelog(err, "get config json data wrong")
	} else {
		temp, _ := ioutil.ReadAll(fi)
		json.Unmarshal(temp, &cConfig)
	}
}

func readFileByLine() {

	fi, err := os.Open(cConfig.FileName)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)

	//----------------------
	ctx := context.Background()

	// // 连接es
	// fmt.Println(cConfig.LoginUser)
	// fmt.Println(cConfig.LoginPWD)

	client, err := elastic.NewClient(elastic.SetURL(cConfig.SubmitURL))

	if cConfig.LoginUser != "" {
		client, err = elastic.NewClient(
			elastic.SetURL(cConfig.SubmitURL),
			elastic.SetBasicAuth(cConfig.LoginUser, cConfig.LoginPWD))
	}
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// ----------------------
	chRequest := make(chan *elastic.BulkService, cConfig.ChannelSize)
	chText := make(chan []string, cConfig.ChannelSize)
	chString := make(chan string, cConfig.SliptFileRowCount*3)
	readEnd := false
	// var regArr []*elastic.BulkIndexRequest

	// go processStringNew(string(a), client, chRequest, &regArr, readEnd)
	// for index := 0; index < len(chRequest); index++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		bulkRequest := <-chRequest
	// 		num := bulkRequest.NumberOfActions()
	// 		bulkRequest.Do(ctx)
	// 		fmt.Println("本地导入", num, "条数据,已经读取<", importAllCount, ">条数据")
	// 		defer wg.Done()
	// 	}()
	// }
	wg.Add(1)
	go func(wg *sync.WaitGroup, wgChild *sync.WaitGroup) {
		for b := range chText {

			bulkService := client.Bulk()
			for _, text := range b {
				// 协程里面套协程，惊艳

				// go processStringNew(c, client, tempService)
				wgChild.Add(1)
				go func(wgChild *sync.WaitGroup, text *string, bulkService *elastic.BulkService) {
					//这里我加了一个同步方法，这个警告是不是可以避免
					textArr := strings.Split(*text, `","`)
					allText := ``
					//m:字段对应列表内容
					m := make(map[string]string)

					for index, str := range textArr {
						prop, ok := keyMap[index+1]
						if ok {

							if prop.IsAllTextColumn {
								tempAll := strings.Trim(str, `"`)
								if prop.Before != "" {
									tempAll = prop.Before + tempAll
								}
								if prop.After != "" {
									tempAll = tempAll + prop.After
								}
								allText += tempAll + `,`
							}
							m[prop.ColumnName] = strings.Trim(str, `"`)
						}

					}

					m["_alltext"] = allText

					jsonStr, err := json.Marshal(m)
					if err != nil {
						panic(err)
					}

					req := elastic.NewBulkIndexRequest().Index(cConfig.EsIndex).Type(cConfig.EsType).Doc(string(jsonStr))
					bulkService.Add(req)
					wgChild.Done()
				}(wgChild, &text, bulkService)
			}
			wgChild.Wait()
			chRequest <- bulkService
		}
		wg.Done()
		close(chRequest)
	}(&wg, &wgChild)

	wgPost.Add(1)
	go func(wgPost *sync.WaitGroup, wgPostChild *sync.WaitGroup, ctxAddr *context.Context) {

		for temp := range chRequest {
			wgPostChild.Add(1)
			go func(bulkServiceTemp *elastic.BulkService, wgPostChild *sync.WaitGroup, ctxAddrChild *context.Context) {
				bulkServiceTemp.Do(*ctxAddrChild)
				wgPostChild.Done()
			}(temp, wgPostChild, ctxAddr)
		}
		wgPost.Done()
	}(&wgPost, &wgPostChild, &ctx)

	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			var chTextEle []string
			for textIndex := 0; textIndex < importAllCount%cConfig.SliptFileRowCount; textIndex++ {
				temp := <-chString
				chTextEle = append(chTextEle, temp)
			}
			chText <- chTextEle
			close(chText)
			fmt.Println("text 关闭了 ", time.Now().Format("2006-01-02 15:04:05"))
			readEnd = true
		}
		if readEnd {
			break
		} else {
			if importAllCount%cConfig.SliptFileRowCount == 0 && importAllCount != 0 {
				var chTextEle []string
				for textIndex := 0; textIndex < cConfig.SliptFileRowCount; textIndex++ {
					temp := <-chString
					chTextEle = append(chTextEle, temp)

				}
				chText <- chTextEle
			}
		}

		chString <- string(a)
		importAllCount++

	}

}

// func doBuldRequest(bulkRequest *elastic.BulkService, ctx context.Context) {
// 	bulkResponse, err := bulkRequest.Do(ctx)
// 	if err != nil {
// 		writelog(err, "doBulk faild")
// 	}
// 	indexed := bulkResponse.Indexed()
// 	fmt.Println("导入了", len(indexed), "条数据,现在导入的是第<", importAllCount, ">条数据")

// }

func processStringNew(text string, client *elastic.Client, bulkService *elastic.BulkService) {

	textArr := strings.Split(text, `","`)
	allText := ``
	//m:字段对应列表内容
	m := make(map[string]string)

	for index, str := range textArr {
		prop, ok := keyMap[index+1]
		if ok {

			if prop.IsAllTextColumn {
				tempAll := strings.Trim(str, `"`)
				if prop.Before != "" {
					tempAll = prop.Before + tempAll
				}
				if prop.After != "" {
					tempAll = tempAll + prop.After
				}
				allText += tempAll + `,`
			}
			m[prop.ColumnName] = strings.Trim(str, `"`)
		}

	}

	m["_alltext"] = allText

	jsonStr, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	req := elastic.NewBulkIndexRequest().Index(cConfig.EsIndex).Type(cConfig.EsType).Doc(string(jsonStr))
	bulkService.Add(req)
	// *regArr = append(*regArr, req)
	// if len(*regArr) == cConfig.SliptFileRowCount && !readEnd {
	// 	tempService := client.Bulk()
	// 	for _, tempReq := range *regArr {
	// 		tempService.Add(tempReq)
	// 	}
	// 	*regArr = make([]*elastic.BulkIndexRequest, 0)
	// 	ch <- tempService
	// }
	// if readEnd {
	// 	fmt.Println("已经读取完毕")
	// 	tempService := client.Bulk()
	// 	for _, tempReq := range *regArr {
	// 		tempService.Add(tempReq)
	// 	}
	// 	ch <- tempService
	// 	close(ch)
	// }
}

func processString(text string) string {

	textArr := strings.Split(text, `","`)
	allText := ``
	//m:字段对应列表内容
	m := make(map[string]string)

	for index, str := range textArr {
		prop, ok := keyMap[index+1]
		if ok {

			if prop.IsAllTextColumn {
				tempAll := strings.Trim(str, `"`)
				if prop.Before != "" {
					tempAll = prop.Before + tempAll
				}
				if prop.After != "" {
					tempAll = tempAll + prop.After
				}
				allText += tempAll + `,`
			}
			m[prop.ColumnName] = strings.Trim(str, `"`)
		}

	}

	m["_alltext"] = allText

	jsonStr, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(jsonStr)
}

//---------------------------------------------
//以后这里的内容需要集成到util包中去

//建立文件
func createFloder(fName string) {
	err := os.Chdir(fName)
	if err != nil {
		os.Mkdir(fName, 0777)
	}
}

func writeFile(str string) {
	file, _ := os.Create("text.txt")

	defer file.Close()

	file.WriteString(str)
}

//写log
func writelog(err error, strDefine string) {
	if checkFileIsExist("errlog") {
		file, _ := os.OpenFile("errlog", os.O_APPEND, 0666)
		defer file.Close()
		io.WriteString(file, err.Error())
	} else {
		file, _ := os.Create("errorlog")

		defer file.Close()

		file.WriteString(err.Error() + "  |  " + strDefine + "\n\r")
	}

}

//判断文件是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//---------------做请求----------------
//---------------弃用-----------------------
func httpDo(postData string) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", cConfig.SubmitURL, strings.NewReader(postData))
	if err != nil {
		writelog(err, "建立请求失败")
	}

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Cookie", "name=anny")
	if cConfig.LoginUser != "" {
		req.SetBasicAuth(cConfig.LoginUser, cConfig.LoginPWD)
	}

	resp, err := client.Do(req)
	// _, err = client.Do(req)
	if err != nil {
		writelog(err, "执行提交失败")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

//---------------弃用-----------------------
func saveOrPostSplitFile(writeResult string) {

	//这里要根据有没有填写URL地址，进行提交操作的判断

	fileName := "result/result" + strconv.Itoa(splitFileNameNum)
	file, err := os.Create(fileName)
	if err != nil {
		writelog(err, "建立文件失败")
	}

	defer file.Close()

	file.WriteString(writeResult)

	splitFileNameNum++
}
