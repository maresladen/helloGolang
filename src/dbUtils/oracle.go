package dbUtils

import (
	ct "commonTools"
	"database/sql"
	"log"

	_ "github.com/mattn/go-oci8"
)

//Test 测试啊底码
func Test() {

	db, err := sql.Open("oci8", "ccic_rel_wf/ccic_rel_wfpwd@10.1.12.143:1521/ccic_rel")
	ct.CheckErr(err)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select 3.14, 'foo' from dual")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var f1 float64
		var f2 string
		rows.Scan(&f1, &f2)
		println(f1, f2) // 3.14 foo
	}
}
