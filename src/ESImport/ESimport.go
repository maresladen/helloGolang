package ESImport

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
	// doTest()
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

	//阻塞等待所有组内成员都执行完毕退栈
	fmt.Println("执行完毕 ", time.Now().Format("2006-01-02 15:04:05"))
}

func readConfig() {
	fi, err := os.Open("../ESImport/config.json")
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
	chText := make(chan chanLogType, 1000)
	chControlGoRoutine := make(chan int, cConfig.ChannelSize)
	chString := make(chan string, cConfig.SliptFileRowCount*3)
	readEnd := false
	doworkCount := 0
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			var chTextEle []string
			for textIndex := 0; textIndex < importAllCount%cConfig.SliptFileRowCount; textIndex++ {
				temp := <-chString
				chTextEle = append(chTextEle, temp)
			}
			fmt.Println("读取了", importAllCount, "条数据，text 关闭了 ", time.Now().Format("2006-01-02 15:04:05"))
			readEnd = true
			doworkCount++
			go doWork(&chTextEle, client, &ctx, chText, chControlGoRoutine, importAllCount)
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

				doworkCount++
				go doWork(&chTextEle, client, &ctx, chText, chControlGoRoutine, importAllCount)
			}
		}

		chString <- string(a)
		importAllCount++

	}

	for index := 0; index < doworkCount; index++ {
		// <-chText
		temp := <-chText
		fmt.Println("导入了第", temp.count, "*** 时间节点为", temp.time)

	}

}

func doWork(b *[]string, client *elastic.Client, ctx *context.Context, ch chan chanLogType, chControlGoRoutine chan int, importCount int) {
	bulkService := client.Bulk()
	for _, text := range *b {

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
	}
	chControlGoRoutine <- 1
	rep, err := bulkService.Do(*ctx)
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
			time:  "has error",
		}
		ch <- *temp
		<-chControlGoRoutine
	}
	if err != nil {
		temp := &chanLogType{
			count: importCount,
			time:  "http has error",
		}
		<-chControlGoRoutine
		ch <- *temp
		fmt.Println("------------------error----------------------")
		fmt.Println("|     ", err)
		fmt.Println("------------------error----------------------")
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
