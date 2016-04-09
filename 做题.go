package main

import (
    "fmt"
    "net/http"
	"io/ioutil"
)
func main() {
//    b:=`pvwdgazxubqfsnrhocitlkeymj,." `
//    a:=`abcdefghijklmnopqrstuvwxyz,." `
//    rsa := []rune(a) 
//    rsb := []rune(b)
//     ml := make(map[string]string)
//    for i:=0;i<len(a);i++{
//        ml[string(rsa[i])] =string(rsb[i])

//    }
   
//    c := `wxgcg txgcg ui p ixgff, txgcg ui p epm. I gyhgwt mrl lig txg ixgff wrsspnd tr irfkg txui hcrvfgs, nre, hfgpig tcm liunz txg crt13 ra "ixgff" tr gntgc ngyt fgkgf.`
//    rsc := []rune(c)
//    d :=""
//    for _,v := range rsc{
//        d += ml[string(v)];
//    }

//     fmt.Println(d)

// tempstr :=`http://fun.coolshell.cn/n/`
s := `http://fun.coolshell.cn/n/32722.html`

 fmt.Println(getAddress(s))



}

func getAddress(s string) (ttt string) {
    res, err := http.Get(s)
	if err != nil {
		fmt.Println("get错误")
	}
	body, err := ioutil.ReadAll(res.Body) //转换byte数组
	if err != nil {
		fmt.Println("read错误")
	}
	defer res.Body.Close()
	//io.Copy(os.Stdout, res.Body)//写到输出流，
	bodystr := string(body)
    fmt.Println(bodystr)
    if len(bodystr) <=5{
       ttt = getAddress(`http://fun.coolshell.cn/n/`+bodystr+".html")
    }
    return ttt
} 
