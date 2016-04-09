package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"os"
)

var floderName string = "槽边往事"

type blogInfo struct {
	Title   string
	Content string
	Cat     []string
	Strdate string
}

func main() {
	createFloder(floderName)
	times := 0
	blogGet(`http://www.caobian.info/?p=1037`, times)
}

func blogGet(urlstr string, times int) {

	// if times >= 3 {
	// 	return
	// }
	// proxy := func(_ *http.Request) (*url.URL, error) {
	// 	return url.Parse("http://127.0.0.1:1088")
	// }

	// transport := &http.Transport{Proxy: proxy}

	// client := &http.Client{Transport: transport}
	// resp, err := client.Get(urlstr)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	
	

	doc, err1 := goquery.NewDocument(urlstr)
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	//-------------------------抬头--------------------------------
	vartitle := doc.Find(".entry-title").Text()

	// fmt.Println(vartitle)
	//-------------------------内容--------------------------------
	varbody := doc.Find(".entry-content").Text()

	// fmt.Println(varbody)
	//-------------------------下一个URL---------------------------
	varNextURL, errNoNext := doc.Find(".nav-next a").Attr("href")
	if !errNoNext {
		fmt.Println("找不到链接")
		return
	}
	// fmt.Println(varNextURL)
	//------------------------分类---------------------------------
	varcatSel := doc.Find("a[rel='category']")
	strCat := make([]string, varcatSel.Length())

	varcatSel.Each(func(i int, s *goquery.Selection) {
		strCat[i] = s.Text()
	})

	// fmt.Println(strCat[0])
	// fmt.Println(strCat[1])
	//-----------------------时间---------------------------------
	tdate := doc.Find(".entry-date").Text()

	// fmt.Println(strdate)

	// _blogInfo := blogInfo{title: vartitle, content: varbody, cat: strCat, strdate: tdate}
	_blogInfo := blogInfo{Title: vartitle, Content: varbody, Cat: strCat, Strdate: tdate}
	blogJSON, errjson := json.Marshal(_blogInfo)
	if errjson != nil {
		fmt.Println("json解析错误")
	}
	// fmt.Println(string(blogJson))
	// file, _ := os.OpenFile(vartitle+".json", os.O_CREATE|os.O_WRONLY, 0)
	// defer file.Close()
	ioutil.WriteFile(floderName+"/"+vartitle+".json", blogJSON, 0777)

	times ++
	blogGet(varNextURL, times)

}

func createFloder(fName string) {
	err := os.Chdir(fName)
	if err != nil {
		os.Mkdir(fName, 0777)
	}
}

