package codeSpider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"utils"

	"github.com/PuerkitoBio/goquery"
)

type spiderConfig struct {
	RootPath       string   `json:"rootPath"`
	FileExtendName []string `json:"fileExtendName"`
	FolderVersion  []string `json:"folderVersion"`
}

var c spiderConfig

//EntryPoint 入口方法
func EntryPoint() {
	//读取文件配置文件
	fi, err := os.Open("spider.json")
	if err != nil {
		utils.WriteFile("error", err.Error()+"get config json data wrong")
	} else {
		temp, _ := ioutil.ReadAll(fi)
		json.Unmarshal(temp, &c)
	}
	rootPath := c.RootPath

	lookupPath(rootPath)
}

func lookupPath(url string) {
	folderV := c.FolderVersion
	for f := range folderV {
		fmt.Print(f)
	}
}

func downLoadFile(url, name string) {

}

func getURLDoc(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		utils.WriteFile("error", err.Error()+"get url doc wrong")
	}
	//查找所有的A标签,进行遍历
	aTag := doc.Find("a")
	aText := aTag.First().Text()
	if strings.HasSuffix(aText, `/`) {
		//文件夹逻辑,判断是否是要找的文件夹，是则透视，否则跳过
	} else {
		href, _ := aTag.Attr("href")
		fmt.Println(href)
		//文件逻辑，是不是要找的文件，是则加入下载通道，否则跳过
	}
}
