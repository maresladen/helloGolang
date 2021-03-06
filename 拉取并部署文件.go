//说明：使用的config文件为gitconfig.json
//使用的sh命令文件为shcmd.sh
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type mconfig struct {
	ProjectPath string
	DotnetPort  int
	Branch      string
}

// gitconfig.json的json格式和内容
// {
// 	"ProjectPath":"/Users/BetaFun/CodeWork/DotNet/asptest",
// 	"DotnetPort": 5000
// }

var m mconfig

func main() {
	fmt.Println("开始配置读取中...")
	configSet()
	fmt.Println("开始文件拉取...")
	//拉取远端代码
	runCMD(m.ProjectPath, "git", "fetch", "origin")

	fmt.Println("开始文件对比中...")
	diffout := runCMD(m.ProjectPath, "git", "diff", "origin/"+m.Branch, "--stat")
	if len(diffout) == 0 {
		fmt.Println("工程没有差异,程序退出")
		return
	}
	fmt.Println("开始文件合并...")
	mergeout := runCMD(m.ProjectPath, "git", "merge", "origin/"+m.Branch)
	if len(mergeout) == 0 {
		writelog(nil, "merge error")
		return
	}
	fmt.Println("开始生成脚本,网站服务将停止...")

	writeScript(strconv.Itoa(m.DotnetPort), m.ProjectPath)
}

func runCMD(runPath, name string, args ...string) []byte {
	cmd := exec.Command(name, args...)
	cmd.Dir = runPath

	out, err := cmd.CombinedOutput()

	if err != nil {
		writelog(err, string(out))
	}
	return out
}

//读取配置文件
func configSet() {

	fi, err := os.Open("gitconfig.json")
	if err != nil {
		writelog(err, "get config json data wrong")
	} else {
		temp, _ := ioutil.ReadAll(fi)
		json.Unmarshal(temp, &m)
	}
}

func writeScript(port, runPath string) {
	//原本这个方法是直接执行,后来改为写sh文件来执行,现在改回原状
	//完全没有好处,直接执行和调用文件执行最后的结果是一样的,还增加开销
	//唯一的好处是知道了覆盖文本文件的命令....
	// file, err := os.OpenFile("shcmd.sh", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0775)
	// defer file.Close()
	// if err != nil {
	// 	writelog(err, "create script file err")
	// } else {

	scmd := `#!/bin/sh
useport=` + port + `
testData=""
for PID in $(lsof -i:$useport |awk '{print $2}'); do
if [ $PID != "PID" ]; then
    if [ -z $PID ]; then
        break;
        else
        testData=$PID
        kill $PID
    fi
fi
done
cd ` + runPath + `
dotnet run`
	// _, err = io.WriteString(file, scmd)
	// if err != nil {
	// 	writelog(err, "write script file err")
	// }
	// cmd := exec.Command("/bin/sh", "./shcmd.sh")
	cmd := exec.Command("/bin/sh", "-c", scmd)
	cmd.Dir, _ = os.Getwd()
	fmt.Println("开始执行脚本,网站将重新启用")
	err := cmd.Run()
	if err != nil {
		if strings.Contains(err.Error(), "exit status 143") {
			fmt.Println("进程被关闭,请检查服务是否被重新开启")
		} else {

			writelog(err, "run commod err")
		}
	}
	// }
}

func writelog(err error, strDefine string) {
	fmt.Println("发生错误,请查看errlog文件")
	if checkFileIsExist("errlog") {
		file, _ := os.OpenFile("errlog", os.O_APPEND, 0666)
		defer file.Close()
		io.WriteString(file, err.Error())
	} else {
		file, _ := os.Create("errorlog")

		defer file.Close()

		file.WriteString(err.Error() + "  |  " + strDefine + "\n\r")
	}

}
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//Json格式的配置文件
// {
// 	"ProjectPath":"/Users/BetaFun/CodeWork/DotNet/asptest",
// 	"DotnetPort": 5000,
// 	"branch" :"develop"
// }

//sh文件内容
// #!/bin/sh
// useport=5000
// testData=""
// for PID in $(lsof -i:$useport |awk '{print $2}'); do
// if [ $PID != "PID" ]; then
//     if [ -z $PID ]; then
//         break;
//         else
//         testData=$PID
//         kill $PID
//     fi
// fi
// done
// cd /Users/BetaFun/CodeWork/DotNet/asptest
// dotnet run
