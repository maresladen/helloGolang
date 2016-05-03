package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
    "strconv"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

type Argjson struct {
	Uin        string
	Key        string
	Biz        string
	PassTicket string
	FormMsgID  string
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

var msgID string
var existSQLID int64

func main() {

	
	db, err := sql.Open("mysql", "root:miao0308@/blogDB?charset=utf8")
	if err !=nil{
		fmt.Println(err)
	}
	
	rows, err1 := db.Query("select max(articleid) aa from wxTable")
		if err1 !=nil{
		fmt.Println(err)
	}
    for rows.Next(){
		rows.Scan(&existSQLID)
	}

	ar, result := configSet()
	if !result {
		fmt.Println("explian config error")
	}

	
    msgID = ar.FormMsgID

	for i := 0; i < 100; i++ {
		if msgID == "" {
			break
		}
        // fmt.Println(msgID)
		data := HTMLCreateGetJSON(&ar)
        if data == nil{
            break
        }
		msgID = getJSONData(data, &ar)
        if msgID == ""{
            break
        } 
	}

}

func configSet() (ar Argjson, result bool) {
	result = true
	fi, err := os.Open("config.json")
	if err != nil {
		fmt.Println("no config file")
		result = false
	} else {
		temp, _ := ioutil.ReadAll(fi)
		err = json.Unmarshal(temp, &ar)
		if err != nil {
			fmt.Println("config explain error")
			result = false
		}
	}
	return ar, result
}

//HTMLCreateGetJSON 获取json数据
func HTMLCreateGetJSON(ar *Argjson) ([]byte) {

	jsonHTML := `https://mp.weixin.qq.com/mp/getmasssendmsg?__biz=` + ar.Biz + `&uin=` + ar.Uin + `&key=` + ar.Key + `&f=json&frommsgid=` + msgID + `&count=10&uin=` + ar.Uin + `&key=` + ar.Key + `&pass_ticket=` + ar.PassTicket + `&wxtoken=&x5=0`

	// fmt.Println(jsonHTML)
    
	doc, err := goquery.NewDocument(jsonHTML)
	if err != nil {
		fmt.Println(err)
        
	}
	varbody := doc.Find("body").Text()

	return []byte(varbody)
}

func getJSONData(doc []byte, ar *Argjson) (msgIDArg string){
	var txtJSON Wxjson
	err := json.Unmarshal(doc, &txtJSON)
	if err != nil {
		fmt.Println("cuole")
        return ""
	}

	var contentJSON General_msg_json
	err = json.Unmarshal([]byte(txtJSON.General_msg_list), &contentJSON)
	if err != nil {
		fmt.Println("cuole2")
        return ""
	}

	// ar.FormMsgID = ""

	// fmt.Println("start")
	for index, jst := range contentJSON.List {
        
        
        // fmt.Println(string(len(contentJSON.List)))
        msgIDArg =""
		if jst.Comm_msg_info.Type == 49 {
			htmlurl := strings.Replace(jst.App_msg_ext_info.Content_url, `amp;`, ``, -1)

			// fmt.Println(htmlurl)
			
			if jst.Comm_msg_info.Id == existSQLID{
				break
			}
			
			jst.App_msg_ext_info.Content_url = htmlurl
			insertValues(jst.App_msg_ext_info,jst.Comm_msg_info.Id)

		}
		if index == len(contentJSON.List)-1 {
			msgIDArg =strconv.FormatInt(jst.Comm_msg_info.Id,10)
		}
        
        
	}
    return msgIDArg
}

func insertValues(st *App_msg_ext_json,pageID int64) {
	db, _ := sql.Open("mysql", "root:miao0308@/blogTable?charset=utf8")
	stmt, _ := db.Prepare("INSERT wxTable SET articleid=?,title=?,url=?")

	stmt.Exec(pageID, st.Title, st.Content_url)
	db.Close()
}
