package main

import (
	"fmt"
    "os/exec"
    "encoding/json"
    "strings"
    "os"
    "io/ioutil"
)

type Argjson struct{
    Uin string
    Key string
    Biz string
    PassTicket string
    FormMsgID string
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

func main() {
    ar,result := configSet() 
    // fmt.Println(result)
    if !result{
        fmt.Println("configset fun error")
    }else{
        // fmt.Println(ar.Uin)
        // fmt.Println(ar.Key)
        // fmt.Println(ar.Biz)
        // fmt.Println(ar.PassTicket)
        // fmt.Println(ar.FormMsgID)
        
        jsonContent := doExec(&ar)
        getJSONData(jsonContent,&ar)
        if ar.FormMsgID != ""{
            fmt.Println("we do for here")
        }
    }

    
}  

func configSet()  (ar Argjson ,result bool){
    result =true
	fi, err := os.Open("config.json")
	if err != nil {
		fmt.Println("no config file")
        result =false
	} else {
		temp, _ := ioutil.ReadAll(fi)
        err = json.Unmarshal(temp,&ar)
        if(err != nil){
            fmt.Println("config explain error")
            result =false
        }
	}
    return ar,result
}

func doExec(ar *Argjson) []byte{
    jsonHTML:=`https://mp.weixin.qq.com/mp/getmasssendmsg?__biz=`+ ar.Biz+`&uin=`+ar.Uin+`&key=`+ar.Key+`&f=json&frommsgid=`+ar.FormMsgID+`&count=10&uin=`+ar.Uin+`&key=`+ar.Key+`&pass_ticket=`+ar.PassTicket+`&wxtoken=&x5=0`
    fmt.Println(jsonHTML)
    cmd := exec.Command("curl",jsonHTML)
    out, err := cmd.CombinedOutput()  
    if err != nil {  
        fmt.Println(err)  
    }  
    fmt.Println(string(out))
    return out
}

func getJSONData(doc []byte,ar *Argjson)  {
     var txtJSON Wxjson
	err := json.Unmarshal(doc, &txtJSON)
	if err != nil {
		fmt.Println("cuole")
	}

	var contentJSON General_msg_json
	err = json.Unmarshal([]byte(txtJSON.General_msg_list), &contentJSON)
	if err != nil {
		fmt.Println("cuole2")
	}
    
    ar.FormMsgID =""
    
    for index, jst := range contentJSON.List {
		// fmt.Println(Jlist)
		if jst.Comm_msg_info.Type == 49 {
			htmlurl := strings.Replace(jst.App_msg_ext_info.Content_url, `amp;`, ``, -1)
                       
			fmt.Println(htmlurl)
		}
        if index == len(contentJSON.List) -1{
            ar.FormMsgID = string(jst.Comm_msg_info.Id)
        }
	}
}