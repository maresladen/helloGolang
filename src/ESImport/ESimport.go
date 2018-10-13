package ESImport

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	strings "strings"
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
	SliptSign         string       `json:"SliptSign"`
	TrimSign          string       `json:"TrimSign"`
	IdField           string       `json:"IdField"`
}

type colmunProp struct {
	ColumnName      string
	ColumnIndex     int
	IsAllTextColumn bool
	FieldType       string
	Before          string
	After           string
}

type chanLogType struct {
	time  string
	count int
}

var cConfig columnConfig
var keyMap map[int]colmunProp
var splitFileNameNum = 1
var importAllCount = 0

//ImporterByText 通过Config文件做导入操作
func ImporterByText() {
	doImportByText()
}

func TestTest() {
	ctx1 := context.Background()

	// Create a client
	client, err := elastic.NewClient(elastic.SetURL("http://10.1.14.178:9206"))
	if err != nil {
		// Handle error
		panic(err)
	}
	templength := 1000
	chDoIndexArr := make(chan int64, templength)

	for i := 0; i < templength; i++ {
		go dotestRequest(client, &ctx1, chDoIndexArr)
	}

	for i := 0; i < templength; i++ {
		temptime := <-chDoIndexArr

		fmt.Println(temptime)
	}

}

func dotestRequest(client *elastic.Client, ctx1 *context.Context, chData chan int64) {

	oldNowTime := time.Now()
	tQuery := elastic.NewTermQuery("TaskName", "NewProposal")
	searchResult, err := client.Search().
		Index("workflow_v3"). // search in index "twitter"
		Type("Workflow").
		Query(tQuery).     // specify the query
		From(0).Size(500). // take documents 0-9
		Pretty(true).      // pretty print requestand response JSON
		Do(*ctx1)          // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	count := searchResult.TotalHits()
	nowTime := time.Now()
	fmt.Println(nowTime.Sub(oldNowTime))

	chData <- count

}

func doImportByJSON() {

	//解析json内容,然后做导入
}

func doImportByText() {

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

	fmt.Println("执行完毕 ", time.Now().Format("2006-01-02 15:04:05"))

}

func readConfig() {
	fi, err := os.Open("./config.json")
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
	//=============重构版===============

	chLog := make(chan chanLogType, 1000)
	var textArrs []string
	indexedCount := 0
	chDoIndexArr := make(chan []string, cConfig.ChannelSize*3)
	chControlGoRoutine := make(chan int, cConfig.ChannelSize)
	for {

		importAllCount++
		a, _, c := br.ReadLine()
		if c == io.EOF {
			if len(textArrs) > 0 {
				indexedCount++
				chDoIndexArr <- textArrs
				chControlGoRoutine <- 1
				go doWorkNew(chDoIndexArr, chControlGoRoutine, client, &ctx, chLog, importAllCount, cConfig.SliptSign, cConfig.TrimSign)
			}
			break
		} else {
			if importAllCount%cConfig.SliptFileRowCount == 0 && importAllCount != 0 {

				textArrs = append(textArrs, strings.TrimSpace(string(a)))
				chDoIndexArr <- textArrs
				textArrs = nil

				chControlGoRoutine <- 1
				textArrs = make([]string, cConfig.SliptFileRowCount)
				indexedCount++
				go doWorkNew(chDoIndexArr, chControlGoRoutine, client, &ctx, chLog, importAllCount, cConfig.SliptSign, cConfig.TrimSign)
			} else {

				textArrs = append(textArrs, strings.TrimSpace(string(a)))
			}
		}
	}

	for index := 0; index < indexedCount; index++ {
		log := <-chLog
		fmt.Println("导入了第", log.count, "*** 时间节点为", log.time)
	}

}

func doWorkNew(chData chan []string, chControlGoRoutine chan int, client *elastic.Client, ctx *context.Context, ch chan chanLogType, importCount int, sliptSign string, trimSign string) {
	strArr := <-chData
	bulkService := client.Bulk()
	for _, text := range strArr {
		textArr := strings.Split(text, sliptSign)
		if len(text) == 0 {
			continue
		}

		//m:字段对应列表内容
		m := make(map[string]string)

		for index, str := range textArr {
			prop, ok := keyMap[index+1]
			if ok {
				m[prop.ColumnName] = strings.Trim(str, trimSign)
			}

		}

		jsonStr, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		docID := m[cConfig.IdField]

		req := elastic.NewBulkIndexRequest().Index(cConfig.EsIndex).Type(cConfig.EsType).Doc(string(jsonStr)).Id(docID)
		bulkService.Add(req)
	}

	rep, err := bulkService.Do(*ctx)

	<-chControlGoRoutine
	if err != nil {
		temp := &chanLogType{
			count: importCount,
			time:  "call error",
		}
		ch <- *temp
	} else {
		if !rep.Errors {
			temp := &chanLogType{
				count: importCount,
				time:  time.Now().Format("2006-01-02 15:04:05"),
			}
			ch <- *temp
		} else {
			temp := &chanLogType{
				count: importCount,
				time:  "do index error",
			}
			ch <- *temp
		}
	}
}

func doWork(b *[]string, client *elastic.Client, ctx *context.Context, ch chan chanLogType, chControlGoRoutine chan int, importCount int, sliptSign string, trimSign string) {
	bulkService := client.Bulk()
	unimportItems := ""
	for _, text := range *b {
		unimportItems += text + "\n"
		textArr := strings.Split(text, sliptSign)
		//m:字段对应列表内容
		m := make(map[string]string)

		for index, str := range textArr {
			prop, ok := keyMap[index+1]
			if ok {
				m[prop.ColumnName] = strings.Trim(str, trimSign)
			}

		}

		jsonStr, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		req := elastic.NewBulkIndexRequest().Index(cConfig.EsIndex).Type(cConfig.EsType).Doc(string(jsonStr))
		bulkService.Add(req)
	}
	chControlGoRoutine <- 1
	rep, err := bulkService.Do(*ctx)
	if err != nil {
		temp := &chanLogType{
			count: importCount,
			time:  "call error",
		}
		<-chControlGoRoutine
		ch <- *temp
	} else {
		if !rep.Errors {
			temp := &chanLogType{
				count: importCount,
				time:  time.Now().Format("2006-01-02 15:04:05"),
			}
			ch <- *temp
			<-chControlGoRoutine
		} else {
			temp := &chanLogType{
				count: importCount,
				time:  "do index error",
			}
			ch <- *temp
			<-chControlGoRoutine
		}
	}
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

//------------------------guf for test---------------------------------

//DoImportForTest  测试用方法
func DoImportForTest() {

	fmt.Println("let's go123! ", time.Now().Format("2006-01-02 15:04:05"))
	//读取配置
	readConfig()
	//初始化一个hash表，用于后续处理字符串
	keyMap = make(map[int]colmunProp)
	for _, prop := range cConfig.Cloumns {
		keyMap[prop.ColumnIndex] = prop
	}

	//调用逐行读取文件方法
	readFileByLineForTest()

	fmt.Println("执行完毕 ", time.Now().Format("2006-01-02 15:04:05"))

}

func readFileByLineForTest() {

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
	chText := make(chan chanLogType, 1000)
	chControlGoRoutine := make(chan int, cConfig.ChannelSize)
	chString := make(chan string, cConfig.SliptFileRowCount*3)
	doworkCount := 0
	for i := 5000000; i <= 6000000; i++ {
		if importAllCount%cConfig.SliptFileRowCount == 0 && importAllCount != 0 {
			var chTextEle []string
			for textIndex := 0; textIndex < cConfig.SliptFileRowCount; textIndex++ {
				temp := <-chString
				chTextEle = append(chTextEle, temp)

			}

			doworkCount++
			go doWorkForTest(&chTextEle, client, &ctx, chText, chControlGoRoutine, importAllCount)
		}

		// chString <- string(`"` + strconv.Itoa(i) + `","` + strconv.Itoa(i) + `","` + strconv.Itoa(i) + `"`)
		chString <- string(`"` + strconv.Itoa(i) + `","` + strconv.Itoa(i) + `千万2号","` + strconv.Itoa(i) + `千万2号"`)

		importAllCount++
	}

	for index := 0; index < doworkCount; index++ {
		// <-chText
		temp := <-chText
		fmt.Println("导入了第", temp.count, "*** 时间节点为", temp.time)
	}

}

func doWorkForTest(b *[]string, client *elastic.Client, ctx *context.Context, ch chan chanLogType, chControlGoRoutine chan int, importCount int) {
	bulkService := client.Bulk()
	unimportItems := ""
	id := ""
	for _, text := range *b {
		unimportItems += text + "\n"
		textArr := strings.Split(text, `","`)
		//m:字段对应列表内容
		m := make(map[string]string)

		id = strings.Trim(textArr[0], `"`)
		for index, str := range textArr {
			prop, ok := keyMap[index+1]
			if ok {

				temp := strings.Trim(str, `"`)
				if prop.Before != "" {
					temp = prop.Before + temp
				}
				if prop.After != "" {
					temp = temp + prop.After
				}
				m[prop.ColumnName] = temp
			}

		}

		jsonStr, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		req := elastic.NewBulkIndexRequest().Index(cConfig.EsIndex).Type(cConfig.EsType).Doc(string(jsonStr)).Parent(id) //.Parent("-97")

		bulkService.Add(req)
	}
	chControlGoRoutine <- 1
	rep, err := bulkService.Do(*ctx)
	if !rep.Errors {
		temp := &chanLogType{
			count: importCount,
			time:  time.Now().Format("2006-07-02 15:04:05"),
		}
		ch <- *temp
		<-chControlGoRoutine
	} else {
		temp := &chanLogType{
			count: importCount,
			time:  "has error",
		}
		ch <- *temp
		<-chControlGoRoutine
	}
	if err != nil {
		temp := &chanLogType{
			count: importCount,
			time:  "has error",
		}
		<-chControlGoRoutine
		ch <- *temp
	}
}
