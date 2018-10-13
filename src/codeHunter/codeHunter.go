package codeHunter

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"utils"

	"github.com/PuerkitoBio/goquery"
)

//NewFunc hehe
func NewFunc() {
	readEnd := false
	fi, err := os.Open("sourceJ")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)

	chControlGoRoutine := make(chan int, 10)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			readEnd = true
		}
		if readEnd {
			break
		} else {
			chControlGoRoutine <- 1
			go DownloadFiles(string(a), chControlGoRoutine)
		}

	}
	fmt.Println("asdf")
}

type matchJSON struct {
	URLSuffixMatch  []string `json:"urlSuffixMatch"`
	FileSuffixMatch []string `json:"fileSuffixMatch"`
}

var c matchJSON

//DownloadFiles 下载
func DownloadFiles(url string, chControlGoRoutine chan int) {

	res, err := http.Get(url)
	if err != nil {
		<-chControlGoRoutine
		fmt.Println("有问题---"+url, res.StatusCode, res.Status)
	} else {
		if res.StatusCode != 200 {
			<-chControlGoRoutine
			fmt.Println("有问题---"+url, res.StatusCode, res.Status)
		} else {
			defer res.Body.Close()
			name := GetFileName(url)
			file, _ := os.Create(name)

			defer file.Close()
			io.Copy(file, res.Body)

			<-chControlGoRoutine
		}
	}

}

//GetDocumentByHTML 测试内容
func GetDocumentByHTML(url string) {
	fileMap := make(map[string]string)
	configReader()

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Status Code Error : %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		attrValue, exist := s.Attr("href")
		if exist {
			if hasuffixMatch(attrValue, 1) {
				//todo: 这是链接
				doLoopURL(url+attrValue, attrValue, fileMap)

			} else if hasuffixMatch(attrValue, 2) {
				//todo: 这里需要判断列表匹配的内容
				//repo.ebaotech.com/artifactory/repo
				fileMap[attrValue] = url + attrValue
			}
		}
	})

	fileTxt := ""
	for _, value := range fileMap {
		fileTxt += value + "\n"
	}

	file, _ := os.Create("pomFile.txt")

	defer file.Close()
	file.WriteString(fileTxt)

}

func GetFileName(url string) string {
	resultName := url
	temp := strings.LastIndex(url, "/") + 1
	resultName = url[temp:len(url)]
	return resultName
}

func doLoopURL(url string, folderName string, fileMap map[string]string) {

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Status Code Error : %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		attrValue, exist := s.Attr("href")
		if exist {
			if hasuffixMatch(attrValue, 1) {
				doLoopURL(url+attrValue, attrValue, fileMap)

			} else if hasuffixMatch(attrValue, 2) {
				//todo: 这里需要判断列表匹配的内容
				//repo.ebaotech.com/artifactory/repo
				fileMap[attrValue] = url + attrValue

			}
		}
	})
}

func configReader() {
	//读取文件配置文件
	fi, err := os.Open("config.json")
	if err != nil {
		utils.WriteFile("error.log", err.Error()+"get config json data wrong")
	} else {
		temp, _ := ioutil.ReadAll(fi)
		json.Unmarshal(temp, &c)
	}

}

func hasuffixMatch(str string, matchType int) bool {
	if matchType == 1 {
		if str == "../" {
			return false
		}
		for _, se := range c.URLSuffixMatch {
			if strings.HasSuffix(str, se) {
				if strings.HasPrefix(str, "unicorn") || strings.HasPrefix(str, "ccic") || strings.HasPrefix(str, "AP26") || strings.HasPrefix(str, "platform") || strings.HasPrefix(str, "starters") || strings.HasPrefix(str, "4.3") || strings.HasPrefix(str, "1.0") {
					return true
				}
			}
		}
	} else {
		fmt.Println(str)
		for _, s := range c.FileSuffixMatch {
			a := strings.HasSuffix(str, s)
			fmt.Println(a)
			if a {
				return true
			}
		}
	}
	return false
}

// 要做的事情，需要判断什么时候建立文件夹，理论上只需要在记录一层，只要判断到找到需要下载了，那就找到上一层的URL，建立文件夹

func getMatchFile() {

}

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
