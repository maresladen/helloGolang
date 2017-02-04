package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

type mconfig struct {
	ProjectPath string
	DotnetPort  int
}

// gitconfig.json的json格式和内容
// {
// 	"ProjectPath":"/Users/BetaFun/CodeWork/DotNet/asptest",
// 	"DotnetPort": 5000
// }

var m mconfig

func main() {
	configSet()
	diffout := runCMD(m.ProjectPath, "git", "diff", "develop", "origin/develop", "--stat")
	if len(diffout) == 0 {
		return
	}
	pullout := runCMD(m.ProjectPath, "git", "pull", "origin")
	if len(pullout) == 0 {
		writelog(nil, "pull error")
		return
	}

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
	file, err := os.OpenFile("shcmd.sh", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0775)
	defer file.Close()
	if err != nil {
		writelog(err, "create script file err")
	} else {

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
		_, err = io.WriteString(file, scmd)
		if err != nil {
			writelog(err, "write script file err")
		}
		cmd := exec.Command("/bin/sh", "./shcmd.sh")
		cmd.Dir, _ = os.Getwd()
		err = cmd.Run()
		if err != nil {
			writelog(err, "run commod err")
		}
	}
}

func writelog(err error, strDefine string) {
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
