package youdaoTranslate

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type i18 struct {
	zh string
}

var youdaoAPI = `http://fanyi.youdao.com/openapi.do?keyfrom=go-translator&key=307165215&type=data&doctype=json&version=1.1&q=`

var replaceSign = strings.NewReplacer(",", "", ":", "", "：", "", "，", "", "(", "", ")", "", "!", "", "?", "", ".", "", "！", "", "？", "", "。", "")

//Web 内容
type Web struct {
	Value []string `json:"value"`
	Key   string   `json:"key"`
}

//Basic 不知道
type Basic struct {
	Phonetic string   `json:"phonetic"`
	Explains []string `json:"explains"`
}

//Translation 不知道
type Translation struct {
	Translation []string `json:"translation"`
	Basic       Basic    `json:"basic"`
	Query       string   `json:"query"`
	ErrorCode   float64  `json:"errorCode"`
	Web         []Web    `json:"web"`
}

func tanslatefun(text string) string {

	client := http.Client{
		Timeout: time.Duration(time.Second * 5),
	}
	urlGet := youdaoAPI + text

	resp, err := client.Get(urlGet)

	if err != nil {
		fmt.Println("出错啦：网络不稳定啊少年，-1s")
		return ""
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println("-----------")
	fmt.Println(string(data))

	if err != nil {
		fmt.Println("出错啦：有道翻译好像出问题了，-1s")
	}

	var j Translation
	err = json.Unmarshal(data, &j)

	if err != nil {
		fmt.Println("出错啦：难道有道已经停用了，-1s")
	}

	if code := j.ErrorCode; code > 0 {
		//errorCode：
		//　0 - 正常
		//　20 - 要翻译的文本过长
		//　30 - 无法进行有效的翻译
		//　40 - 不支持的语言类型
		//　50 - 无效的key
		//　60 - 无词典结果，仅在获取词典结果生效
		switch code {
		case 20:
			fmt.Println("出错啦：要翻译的文本过长")
		case 30:
			fmt.Println("出错啦：无法进行有效的翻译")
		case 40:
			fmt.Println("出错啦：不支持的语言类型")
		case 50:
			fmt.Println("出错啦：无效的key")
		case 60:
			fmt.Println("出错啦：无词典结果，仅在获取词典结果生效")
		}

		return ""
	}

	fmt.Printf("翻译：\t%s\n", strings.Join(j.Translation[:], " / "))
	return strings.Join(j.Translation[:], "")
}

func httpGet() {
	resp, err := http.Get("http://www.01happy.com/demo/accept.php?id=1")
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

//TranslateText 翻译
func TranslateText() {
	// httpGet()
	// text := "红色的苹果"
	// enc := mahonia.NewEncoder("utf8")
	//converts a  string from UTF-8 to gbk encoding.
	// fmt.Println(enc.ConvertString(text))

	// tanslatefun(enc.ConvertString(text))
	// tanslatefun(url.QueryEscape(text))

	// urlGet := youdaoAPI + text
	// encodeurl := url.QueryEscape(urlGet)
	// fmt.Println(encodeurl)
	// decodeurl, err := url.QueryUnescape(encodeurl)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(decodeurl)

	textCN, textEng := configSet()

	saveFile("中文.txt", textCN)
	saveFile("英文.txt", textEng)

	// temptext := getCode("cardiovascular")
	// fmt.Println(temptext)

}

func dojob(hans string) (string, string) {

	cnText := ""
	engText := ""
	//这里套个循环
	// 默认
	// a := pinyin.NewArgs()
	// abc := pinyin.Pinyin(hans, a)
	// var result string
	// for _, b := range abc {
	// 	for _, c := range b {
	// 		singleResult := strings.ToUpper(string(c[0]))
	// 		result += singleResult
	// 	}
	// }
	english := tanslatefun(url.QueryEscape(hans))
	code := getCode(english)
	cnText += `` + code + `:'` + hans + `',`
	engText += `` + code + `:'` + english + `',`

	return cnText, engText
}

func getCode(text string) string {
	fmt.Println(text)

	text = replaceSign.Replace(text)
	codeList := strings.Fields(text)
	result := ""
	if len(codeList) > 1 {
		for index, eleText := range codeList {
			if index == 0 {
				result += strings.ToLower(eleText)
			} else {
				fir := strings.ToUpper(string(eleText[0]))
				result += fir
				oth := strings.ToLower(string(eleText[1:len(eleText)]))
				result += oth
				// result +=strings.ToUpper(string(eleText[0])) + strings.ToLower(string(eleText[1,len(eleText)]))
			}
		}
	} else {
		result = strings.ToLower(text)
	}
	return result
}

//读取配置文件
func configSet() (string, string) {

	fi, err := os.Open("zh.txt")
	if err != nil {
		writelog(err, "get config json data wrong")
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	cntext := ""
	engtext := ""
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		tempCN, tempEng := dojob(string(a))
		cntext += tempCN + "\n"
		engtext += tempEng + "\n"
	}
	return cntext, engtext

}

func saveFile(fileName, content string) {
	file, err := os.Create(fileName)
	if err != nil {
		writelog(err, "建立文件失败")
	}

	defer file.Close()

	file.WriteString(content)
}

func writelog(err error, strDefine string) {
	if checkFileIsExist("errlog") {
		file, _ := os.OpenFile("errlog", os.O_APPEND, 0666)
		defer file.Close()
		io.WriteString(file, err.Error())
	} else {
		file, _ := os.Create("errorlog")

		defer file.Close()

		file.WriteString(err.Error() + "  |  " + strDefine + "\r\n")
	}
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
