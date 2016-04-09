package main

import (
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "fmt"
    //"time"
)

func main() {
    db, err := sql.Open("mysql", "root:miao0308@/blogDB?charset=utf8")
    
    checkErr(err)

    // 插入数据
    // stmt, err := db.Prepare("INSERT blogTable SET blogTitle=?,blogContent=?,blogCat=?,blogDate=?")
    // checkErr(err)

    // _, err = stmt.Exec("中文", "dd", "dd","dd")
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


    // // //删除数据
    // stmt, err := db.Prepare("delete from blogTable where id=?")
    // checkErr(err)

    // res, err := stmt.Exec(2)
    // checkErr(err)

    // affect, err := res.RowsAffected()
    // checkErr(err)

    // fmt.Println(affect)


    // //查询数据
    rows, err := db.Query("SELECT * FROM  blogTable")
    checkErr(err)
    // ccc ,_ := rows.Columns()

    // fmt.Println(ccc)
    for rows.Next() {
        var id int
        var blogTitle string
        var blogContent string
        var blogCat string
        var blogDate string
        err = rows.Scan(&id,&blogTitle,&blogContent,&blogCat,&blogDate)
        ccc ,_ := rows.Columns()
        fmt.Println(ccc)
        // checkErr(err)
        
        fmt.Println(ccc)
        fmt.Println(id)
        fmt.Println(blogTitle)
        fmt.Println(blogContent)
        fmt.Println(blogCat)
        fmt.Println(blogDate)
        
        
    }



    db.Close()

}

func checkErr(err error) {
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
}