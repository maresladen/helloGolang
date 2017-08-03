package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type ebaoConfig struct {
	EUserName       string // 用户名
	EUserPwd        string // 密码
	EUserToken      string // 密钥
	EJavaRepoPath   string // java验证地址
	ENpmRepoPath    string // npm验证地址
	ENpmPath        string // npm下载地址
	EJavaSVNPAth    string // java仓库地址
	ERainbowSVNPath string // UI仓库地址
	ENpmRegistry    string // npm注册地址
	ECurlPath       string // curl下载地址
}

var (
	buildType, userName, userPwd, isSure, mavenToken string
	eJSON                                            ebaoConfig
)

func main() {

	// rebyte, _ := exec.Command("echo", "abbbcc").Output()
	// fmt.Println(string(rebyte))
	// var str3 = "\U0000006dif\U00000073\U00000063\U00000068\U00000069\U0000006e\U00000061\U00005f00\U00006e90\U00004e2d\U000056fd"
	// fmt.Println(string(str3))

	configSet()
	if eJSON.EUserName != "" {
		userName = eJSON.EUserName
	}
	if eJSON.EUserPwd != "" {
		userPwd = eJSON.EUserPwd
	}
	if eJSON.EUserToken != "" {
		mavenToken = eJSON.EUserToken
	}
	if userName == "" && userPwd == "" {
		//用户名的录入
		inputUserAndPwd()
	}
	// 选择后端还是前端
	backOrFront()
	// dotest()

}

//dotest 用于做测试
func dotest() {

	// strtemttt := `registry=http://repo.ebaotech.com/artifactory/api/npm/npm-all` + "\r\n"
	// strNpmrc := `_auth="YWxtLmNjaWMuZ3VmZW5nemdzOkFQNXpCbjVVdmc1ZlB3Szk2WUw3UzZld20yYQ=="always-auth=trueemail=asdfasdfas@asdfasdf.com`
	// tempindex1 := strings.Index(strNpmrc, "always-auth")
	// tempindex2 := strings.Index(strNpmrc, "email=")
	// strtemttt += strNpmrc[0:tempindex1] + "\r\n"
	// fmt.Println(string(strNpmrc[0:tempindex1]))
	// if tempindex2 >= 0 {
	// 	strAuth := string(strNpmrc[tempindex1:tempindex2])
	// 	fmt.Println(strAuth)
	// 	strtemttt += strAuth + "\r\n"
	// 	strEmail := string(strNpmrc[tempindex2:len(strNpmrc)])
	// 	fmt.Println(strEmail)
	// 	strtemttt += strEmail
	// } else {

	// 	strAuth := string(strNpmrc[tempindex1:len(strNpmrc)])
	// 	fmt.Println(strAuth)
	// 	strtemttt += strAuth + "\r\n"
	// }
	// fmt.Println("-----------------------------")
	// fmt.Println(strtemttt)

	//测试oracle监听文件拷贝问题

	// oraclePath := runCMD("where oci.dll")
	// fmt.Println(oraclePath)
	// var testScan string
	// fmt.Scanln(&testScan)
	// if testScan == "Y" {
	// 	fmt.Println("我们继续下一步")
	// }

	// oraCopyPath := strings.Replace(oraclePath, `BIN\oci.dll`, `network\admin\`, -1)
	// fmt.Println(oraCopyPath)
	// var testScan1 string
	// fmt.Scanln(&testScan1)
	// if testScan1 == "Y" {
	// 	fmt.Println("我们继续下一步")
	// }

	// result := runCMD("copy tnsnames.ora " + oraCopyPath + "tnsnames.ora")

	// fmt.Println(result)
	// var testScan2 string
	// fmt.Scanln(&testScan2)

	// if testScan1 == "Y" {
	// 	fmt.Println("结束了")
	// }
	//--------运行命令-------------------
	// cmd := exec.Command("/bin/zsh", `-c`, `where go`)

	// out, err := cmd.CombinedOutput()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(out))

	//-----------------测试批处理----------------------
	// excuPath := getCurrentPath()
	// strInstallnpmfile := `installNpm.bat`
	// strInstallYofile := `installYo.bat`
	// strInstallClifile := `installCli.bat`
	// strYoClifile := `yocli.bat`
	// strBatGBKfile := `batGBK.bat`
	// strInstallnpm := `echo 1`
	// // strInstallnpm := ` msiexec /i node.msi /qb `
	// saveFile(strInstallnpmfile, strInstallnpm)
	// strInstallYo := `echo 2`
	// // strInstallYo := `npm install –g yo `
	// saveFile(strInstallYofile, strInstallYo)
	// strInstallCli := `echo 3`
	// // strInstallCli := `npm install -g generator-rainbowui-cli `
	// saveFile(strInstallClifile, strInstallCli)
	// strYoCli := `echo 4`
	// // strYoCli := `yo rainbowui-cli`
	// saveFile(strYoClifile, strYoCli)

	// //这里是批处理文件，到时候运行   /qb是弹出窗口安装 /qn静默安装

	// strBatUTF := `@echo off
	// 	call ` + excuPath + `\` + strInstallnpmfile + `
	// 	call ` + excuPath + `\` + strInstallYofile + `
	// 	call ` + excuPath + `\` + strInstallClifile + `
	// 	mkdir  rainbowUI
	// 	cd rainbowUI
	// 	call ` + excuPath + `\` + strYoClifile + `
	// 	pause
	// 	exit`
	// byteGBK, err := Utf8ToGbk([]byte(strBatUTF))
	// if err != nil {
	// 	writelog(err, "编码转换失败")
	// }
	// strBatGBK := string(byteGBK)
	// saveFile(strBatGBKfile, strBatGBK)
}

//Utf8ToGbk utf转为GBK格式
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func tempfunc(strCmd, dirPath string) string {
	cmd := exec.Command(strCmd)
	cmd.Dir = dirPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	return string(out)
}

//getZS 寻找数字以内的所有的质数
func getZS() int {
	result := 1
	slice1 := []int{2, 3}
	result *= 2
	result *= 3
	n := 5
	for n < 50500 {
		jumpSign := false
		for i := 0; i < len(slice1); i++ {
			if n%slice1[i] == 0 {
				jumpSign = true
				break

			}
		}
		if !jumpSign {
			result *= n
			slice1 = append(slice1, n)
		}
		n += 2
	}
	fmt.Println(len(slice1))
	return result
}

//读取配置文件
func configSet() {
	fi, err := os.Open("ebaoConfig.json")
	if err != nil {
		writelog(err, "get config json data wrong")
	} else {
		temp, _ := ioutil.ReadAll(fi)
		json.Unmarshal(temp, &eJSON)

	}
}

//后端和前端的判断
func backOrFront() {
	fmt.Println("请输入搭建环境的参数:S(后台Java环境),C(前端React环境)")
	fmt.Scanln(&buildType)
	if strings.ToUpper(buildType) == "S" {
		if mavenToken == "" {
			fmt.Println("您将搭建后端java环境，请登录以下网址获得eBao的maven私有库的权限密钥,并录入")
			fmt.Println(eJSON.EJavaRepoPath)
			fmt.Scanln(&mavenToken)
		}

		mavenNameToken := `
		<server>
		<id>eBaoTech</id>
		<username>` + userName + `</username>
		<password>` + mavenToken + `</password>
		</server>
		`
		betafunSign := "BETAFUNSIGN"

		mavenStr := readfile("setting.xxxml")
		saveMavenStr := strings.Replace(mavenStr, betafunSign, mavenNameToken, -1)
		fmt.Println(saveMavenStr)
		saveFile("maven-setting.xml", saveMavenStr)

		oraclePath := runCMD("where oci.dll")
		oraCopyPath := strings.Replace(oraclePath, `BIN\oci.dll`, `network\admin\`, -1)

		runCMD("copy tnsnames.ora " + oraCopyPath + "tnsnames.ora")
		// runCMD("mkdir ebaoJava")

		excuPath := getCurrentPath() // 当前程序所在的文件路径
		javaProjectbat := `@echo off` + "\r\n" + ` rem ` + "\r\n" + `c:` + "\r\n" + ` cd C:\Program Files\TortoiseSVN\bin` + "\r\n" + ` TortoiseProc.exe /command:checkout /url:` + eJSON.EJavaSVNPAth + ` /path:"` + excuPath + `\ebaoJava"` + ` /closeend:1 ` + "\r\n" + `exit`
		saveAndRunBat(javaProjectbat, "jsvnget.bat")
		// 读取setting文件，然后替换其中的用户和密钥，再保存
		// 查找oracle的监听文件目录,保存监听文件到这个路径
		// 通过bat下载源代码，这里的源代码路径是可配置的

	} else {
		fmt.Println("您将搭建前端React环境，请等待curl、nodejs环境的下载和安装")
		strExt := getFileExt(eJSON.ECurlPath, FileSuffix)
		curlFile := "curl" + strExt
		// strings. eJSON.ECurlPath
		fmt.Println("下载中......请等待")
		// 下载curl
		downLoadFile(eJSON.ECurlPath, curlFile)

		// 解压
		excuPath := getCurrentPath()             // 当前程序所在的文件路径
		zipCurlPath := excuPath + `\` + curlFile // curl文件名,带路径
		unpackCurlPath := excuPath + `\cfile`    // curl解压文件路径

		unzipCommandtemp := `@echo off` + "\r\n" + ` rem ` + "\r\n" + `c:` + "\r\n" + `cd C:\Program Files\7-Zip ` + "\r\n" + `7z x "` + zipCurlPath + `" -y -aos -o"` + unpackCurlPath + `"` + "\r\n" + ` exit ` // 解压命令行,也许不一定能过，需要分开做
		unzipCommand, err := Utf8ToGbk([]byte(unzipCommandtemp))
		if err != nil {
			writelog(err, "转换中文失败")
		}

		saveAndRunBat(string(unzipCommand), "unzip.bat")
		// 请求注册仓库的token并保存
		strDisk := string(unpackCurlPath[0]) + string(unpackCurlPath[1]) // 获取盘符，请允许我用这么low的方式

		tempfiles, _ := ioutil.ReadDir(unpackCurlPath)
		var curlRealPath string
		for _, tempfile := range tempfiles {
			if tempfile.IsDir() {
				if strings.Contains(tempfile.Name(), "curl") {

					curlRealPath = unpackCurlPath + `\` + tempfile.Name()
				}
			}
		}

		registerCommandtemp := `@echo off` + "\r\n" + "rem \r\n" + strDisk + "\r\n" + `cd ` + curlRealPath + "\r\n" + `curl -u ` + userName + `:` + userPwd + ` http://repo.ebaotech.com/artifactory/api/npm/auth` + "\r\n" + `exit` // 获取命令行

		registerCommand, err := Utf8ToGbk([]byte(registerCommandtemp))
		if err != nil {
			writelog(err, "转换GBK失败")
		}

		registerOutPut := saveAndRunBat(string(registerCommand), "register.bat")
		strNpmrc := `registry=http://repo.ebaotech.com/artifactory/api/npm/npm-all` + "\r\n"

		tempindex1 := strings.Index(registerOutPut, "always-auth")
		tempindex2 := strings.Index(registerOutPut, "email =")
		if tempindex2 <= 0 {
			tempindex2 = strings.Index(registerOutPut, "email=")
		}
		fmt.Println(string(registerOutPut[0:tempindex1]))
		strNpmrc += registerOutPut[0:tempindex1] + "\r\n"
		if tempindex2 >= 0 {
			strAuth := string(registerOutPut[tempindex1:tempindex2])
			fmt.Println(strAuth)
			strNpmrc += strAuth + "\r\n"
			strEmail := string(registerOutPut[tempindex2:len(registerOutPut)])
			fmt.Println(strEmail)
			strNpmrc += strEmail
		} else {

			strAuth := string(registerOutPut[tempindex1:len(registerOutPut)])
			fmt.Println(strAuth)
			strNpmrc += strAuth + "\r\n"
		}
		fmt.Println("----------------------------")
		fmt.Println(strNpmrc)

		saveFile(".npmrc", strNpmrc) //保存文件
		// 下载npm 下载地址需要可以设定
		downLoadFile(eJSON.ENpmPath, "node.msi")

		//这里把bat文件拆分，用于解决命令不执行问题

		strInstallnpmfile := `installNpm.bat`
		strRegistryNpmFactoryfile := `registryNpmFactory.bat`
		strInstallYofile := `installYo.bat`
		strInstallClifile := `installCli.bat`
		strYoClifile := `yocli.bat`
		strBatGBKfile := `batGBK.bat`

		strInstallnpmbyte, err := Utf8ToGbk([]byte(` msiexec /i node.msi /qb `))
		if err != nil {
			writelog(err, "installnpm转换GBK失败")
		}
		strInstallnpm := string(strInstallnpmbyte)
		saveFile(strInstallnpmfile, strInstallnpm)
		strRegistryNpmFactorybyte, err := Utf8ToGbk([]byte(` npm config set registry ` + eJSON.ENpmRegistry))
		if err != nil {
			writelog(err, "strRegistryNpmFactorybyte转换GBK失败")
		}
		strRegistryNpmFactory := string(strRegistryNpmFactorybyte)
		saveFile(strRegistryNpmFactoryfile, strRegistryNpmFactory)
		strInstallYobyte, err := Utf8ToGbk([]byte(`npm install -g yo `))
		if err != nil {
			writelog(err, "strInstallYobyte转换GBK失败")
		}
		strInstallYo := string(strInstallYobyte)
		saveFile(strInstallYofile, strInstallYo)
		strInstallClibyte, err := Utf8ToGbk([]byte(`npm install -g generator-rainbowui-cli `))
		if err != nil {
			writelog(err, "strInstallClibyte转换GBK失败")
		}
		strInstallCli := string(strInstallClibyte)
		saveFile(strInstallClifile, strInstallCli)
		strYoClibyte, err := Utf8ToGbk([]byte(`yo rainbowui-cli`))
		if err != nil {
			writelog(err, "strYoClibyte转换GBK失败")
		}
		strYoCli := string(strYoClibyte)
		saveFile(strYoClifile, strYoCli)

		//这里是批处理文件，到时候运行   /qb是弹出窗口安装 /qn静默安装

		installedNode := runCMD("where node")

		var strInstallNnnpm string
		if strings.Contains(installedNode, "node.exe") {
			strInstallNnnpm = ""
		} else {
			strInstallNnnpm = ` call ` + excuPath + `\` + strInstallnpmfile + "\r\n"
		}

		strBatUTF := `@echo off` + "\r\n" + strInstallNnnpm + `call ` + excuPath + `\` + strRegistryNpmFactoryfile + "\r\n" + `call ` + excuPath + `\` + strInstallYofile + "\r\n" + `call ` + excuPath + `\` + strInstallClifile + "\r\n" + `mkdir  rainbowUI` + "\r\n" + `cd rainbowUI` + "\r\n" + `call ` + excuPath + `\` + strYoClifile + "\r\n" + `pause` + "\r\n" + `exit`
		byteGBK, err := Utf8ToGbk([]byte(strBatUTF))
		if err != nil {
			writelog(err, "编码转换失败")
		}
		strBatGBK := string(byteGBK)

		saveAndRunBat(strBatGBK, strBatGBKfile)

		// 注册npm仓库
		// TODO: 选择SVN项目还是脚手架自建，进入目录做下载
		// -----------------------脚手架-----------------------
		// 安装脚手架 npm install –g yo
		// 安装rainbowUI npm install -g generator-rainbowui-cli
		// 进入目录 搭建rainbow   yo rainbowui-cli
		// -----------------------现有svn-----------------------
		// 这里的源代码路径可以配置

	}
}

func saveAndRunBat(command, fileName string) string {
	saveFile(fileName, command)

	rebyte, err := exec.Command("cmd.exe", "/c", fileName).Output()
	if err != nil {
		writelog(err, "运行批处理失败")
	}
	return string(rebyte)
}

func downLoadFile(url, fileName string) {
	res, err := http.Get(url)
	if err != nil {
		writelog(err, "获取文件错误!")
	}
	defer res.Body.Close()
	file, _ := os.Create(fileName)

	defer file.Close()
	io.Copy(file, res.Body)
}

func readfile(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		writelog(err, "读取文件失败")
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

//runCMD 弃用
func runCMD(strCmd string) (result string) {
	cmd := exec.Command("cmd.exe", `/c`, strCmd)

	out, err := cmd.CombinedOutput()
	if err != nil {
		writelog(err, "批处理执行失败")
	}
	result = string(out)
	return result
}

//用户名和密码的录入
func inputUserAndPwd() {
	index := 0

	for {
		if index == 0 {
			fmt.Println("请输入用户名")
			fmt.Scanln(&userName)
			fmt.Println("请输入密码")
			fmt.Scanln(&userPwd)
		} else {

			fmt.Println("请重新输入用户名")
			fmt.Scanln(&userName)
			fmt.Println("请重新输入密码")
			fmt.Scanln(&userPwd)
		}
		index++

		fmt.Printf("请确认您的用户名：%s,密码：%s。y/n \r\n", userName, userPwd)
		fmt.Scanln(&isSure)
		if strings.ToUpper(isSure) == "Y" {
			return
		}
	}

}

func saveFile(fileName, content string) {
	file, err := os.Create(fileName)
	if err != nil {
		writelog(err, "建立文件失败")
	}

	defer file.Close()

	file.WriteString(content)
}

func writelog(err error, strDefine string) {
	if checkFileIsExist("errlog") {
		file, _ := os.OpenFile("errlog", os.O_APPEND, 0666)
		defer file.Close()
		io.WriteString(file, err.Error())
	} else {
		file, _ := os.Create("errorlog")

		defer file.Close()

		file.WriteString(err.Error() + "  |  " + strDefine + "\r\n")
	}

	// panic(strDefine)

}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//FileExtType 通过文件名要获取的字符串类型
type FileExtType int

const (
	//FullNameWithSuffix 文件全名
	FullNameWithSuffix FileExtType = iota
	//FileSuffix 后缀名 有点
	FileSuffix
	//FileNameOnly 文件名
	FileNameOnly
)

func getFileExt(fullFilename string, resultType FileExtType) string {

	result := ""
	filenameWithSuffix := path.Base(fullFilename)
	if resultType == FullNameWithSuffix {
		result = filenameWithSuffix
	}
	fileSuffix := path.Ext(filenameWithSuffix)
	if resultType == FileSuffix {
		result = fileSuffix
	}

	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	if resultType == FileNameOnly {
		result = filenameOnly
	}

	return result

}

func getCurrentPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		writelog(err, "获取执行路径失败")
	}
	path, err := filepath.Abs(file)
	if err != nil {
		writelog(err, "获取执行路径失败")
	}
	result := filepath.Dir(path)
	return result
}
