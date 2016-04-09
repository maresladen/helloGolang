package main

import (
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "path/filepath"
    "os"
)

type blogInfo struct {
	Title   string
	Content string
	Cat     []string
	Strdate string
}

func main() {
   cbwsList := getFilelist("/home/cooldan/下载/槽边往事")
    for indexs,fname := range cbwsList{
        fmt.Println(indexs ,"begin")
        dosql(fname)
    }
        fmt.Println("success")
}

func getFilelist(path string) []string{

    //    result := make([]string,0)
       var result []string
        err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
                if ( f == nil ) {return err}
                if f.IsDir() {return nil}
                result = append(result,path)
                return nil
        })
        if err != nil {
                fmt.Printf("filepath.Walk() returned %v\n", err)
        }
        return result
}


func dosql(fpath string) {
    jsbyte,_:= ioutil.ReadFile(fpath)
    var _blog blogInfo
    err := json.Unmarshal(jsbyte,&_blog)
    if err != nil{
        fmt.Println(err)
    }
    insertValues(&_blog)

}

func insertValues(st *blogInfo)  {
    db, _ := sql.Open("mysql", "root:miao0308@/blogTable?charset=utf8")
    stmt, _ := db.Prepare("INSERT blogTable SET blogTitle=?,blogContent=?,blogCat=?,blogDate=?")
   var strCat string
    for _,str := range st.Cat{
        strCat += str+","
    }
    stmt.Exec(st.Title, st.Content, strCat,st.Strdate)
    db.Close()
}

// func dbManager() {
    
    
//     fmt.Println("sadfjl")
    // db, err := sql.Open("mysql", "root:miao0308@/blogTable?charset=utf8")
    // checkErr(err)

    // //插入数据
    // stmt, err := db.Prepare("INSERT userinfo SET username=?,departname=?,created=?")
    // checkErr(err)

    // res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
    // checkErr(err)

    // id, err := res.LastInsertId()
    // checkErr(err)

    // fmt.Println(id)
    // //更新数据
    // stmt, err = db.Prepare("update userinfo set username=? where uid=?")
    // checkErr(err)

    // res, err = stmt.Exec("astaxieupdate", id)
    // checkErr(err)

    // affect, err := res.RowsAffected()
    // checkErr(err)

    // fmt.Println(affect)

    // // 查询数据
    // rows, err := db.Query("SELECT * FROM userinfo")
    // checkErr(err)

    // for rows.Next() {
    //     var uid int
    //     var username string
    //     var department string
    //     var created string
    //     err = rows.Scan(&uid, &username, &department, &created)
    //     checkErr(err)
    //     fmt.Println(uid)
    //     fmt.Println(username)
    //     fmt.Println(department)
    //     fmt.Println(created)
    // }

    // //删除数据
    // stmt, err = db.Prepare("delete from userinfo where uid=?")
    // checkErr(err)

    // res, err = stmt.Exec(id)
    // checkErr(err)

    // affect, err = res.RowsAffected()
    // checkErr(err)

    // fmt.Println(affect)

    // db.Close()

// }

// func checkErr(err error) {
//     if err != nil {
//         panic(err)
//     }
// }