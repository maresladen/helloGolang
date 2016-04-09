package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
    "github.com/PuerkitoBio/goquery"
    "github.com/robertkrimen/otto"
)



var (
    strFromMsgID string
    strUin string
    strKey string
    strBiz string
    strPassTicket string
    sharehtmlAddress string
)



func main() {

	// data := readfile("test.json")
	// var txtJSON Wxjson
	// err := json.Unmarshal(data, &txtJSON)
	// if err != nil {
	// 	fmt.Println("cuole")
	// }

	// var contentJSON General_msg_json
	// err = json.Unmarshal([]byte(txtJSON.General_msg_list), &contentJSON)
	// if err != nil {
	// 	fmt.Println("cuole2")
	// }

	// for _, jst := range contentJSON.List {
	// 	// fmt.Println(Jlist)
	// 	if jst.Comm_msg_info.Type == 49 {
	// 		htmlurl := strings.Replace(jst.App_msg_ext_info.Content_url, `amp;`, ``, -1)
	// 		fmt.Println(htmlurl)
	// 	}

	// }
    
    configSet();
    if sharehtmlAddress != ""{
            getWXDataHtml(sharehtmlAddress);
    }
    
    doForJsonData()
    
	// dd := strings.Replace(cc, `amp;`, ``, -1)

}

type Wxjson struct {
	Ret              int64
	Errmsg           string
	General_msg_list string
	Bizuin_code      string
	Uin_code         string
	Key              string
	Is_friend        int64
	Is_continue      int64
	Count            int64
}

type General_msg_json struct {
	List []*Jlist
}

type Jlist struct {
	Comm_msg_info      *Comm_msg_json
	App_msg_ext_info   *App_msg_ext_json
	Image_msg_ext_info *Image_msg_ext_Json
}

type Comm_msg_json struct {
	Id       int64
	Type     int64
	Datetime int64
	Fakeid   string
	Status   int64
	Content  string
}

type Image_msg_ext_Json struct {
	Length  int64
	Fileid  int64
	Mediaid int64
}

type App_msg_ext_json struct {
	Title                   string
	Digest                  string
	Content                 string
	Fileid                  int64
	Content_url             string
	Source_url              string
	Cover                   string
	Subtype                 int64
	Is_multi                int64
	Multi_app_msg_item_list []string
	Author                  string
}


func configSet()  {
	fi, err := os.Open("config.txt")
	if err != nil {
		fmt.Println("呃...")
	} else {
		temp, _ := ioutil.ReadAll(fi)
        sharehtmlAddress =string(temp)	
        fmt.Println(sharehtmlAddress)
	}
}

func getWXDataHtml(urlstr string) {
	//通过GOquery获取内容，并取得第一个msgid，然后调用循环获取jsondata方法
    doc, err := goquery.NewDocument(urlstr)
	if err != nil {
		fmt.Println(err)
		return
	}
    // fmt.Println(doc.Html()) 
    
    vm := otto.New()
    vm.Run(doc)
    uin,_ := vm.Get("uin")
    key,_:=vm.Get("key")
    biz,_:=vm.Get("biz")
    passTicket,_:=vm.Get("pass_ticket")
    formMsgID,_:=vm.Get("frmMsgId")
    
    strUin,_ = uin.ToString()
    strKey,_ = key.ToString()
    strBiz,_ = biz.ToString()
    strPassTicket,_ = passTicket.ToString()
    strFromMsgID,_ = formMsgID.ToString()
    
    
    fmt.Println(strUin)
    fmt.Println(strKey)
    fmt.Println(strBiz)
    fmt.Println(strPassTicket)
    fmt.Println(strFromMsgID)
    // wxjsonAdd:=`https://mp.weixin.qq.com/mp/getmasssendmsg?__biz=`+ strBiz+`&uin=`+strUin+`&key=`+strKey+`&f=json&frommsgid=`+strFromMsgID+`&count=10&uin=`+strUin+`&key=`+strKey+`&pass_ticket=`+strPassTicket+`&wxtoken=&x5=0`
    
    
    
    // doForJsonData()
}

func doForJsonData() {
    jsonHtml :=`https://mp.weixin.qq.com/mp/getmasssendmsg?__biz=`+ strBiz+`&uin=`+strUin+`&key=`+strKey+`&f=json&frommsgid=`+strFromMsgID+`&count=10&uin=`+strUin+`&key=`+strKey+`&pass_ticket=`+strPassTicket+`&wxtoken=&x5=0`
	doc, err := goquery.NewDocument(jsonHtml)
	if err != nil {
		fmt.Println(err)
		return
	}
    
    
    var txtJSON Wxjson
	err = json.Unmarshal([]byte(doc.Text()), &txtJSON)
	if err != nil {
		fmt.Println("cuole")
	}

	var contentJSON General_msg_json
	err = json.Unmarshal([]byte(txtJSON.General_msg_list), &contentJSON)
	if err != nil {
		fmt.Println("cuole2")
	}
    
    for index, jst := range contentJSON.List {
		// fmt.Println(Jlist)
		if jst.Comm_msg_info.Type == 49 {
			htmlurl := strings.Replace(jst.App_msg_ext_info.Content_url, `amp;`, ``, -1)
                       
			fmt.Println(htmlurl)
		}
        if index == len(contentJSON.List) -1{
            strFromMsgID = string(jst.Comm_msg_info.Id)
        }

	}
    
}

func readfile(path string) []byte {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, _ := ioutil.ReadAll(fi)
	return fd
}