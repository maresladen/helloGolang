package stringEditor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

//DoStringEditor 字符串处理
func DoStringEditor() {

	fi, err := os.Open("./sourceFile/stringeditor.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)

	file1 := ""
	file2 := ""
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		tempfile1, tempfile2 := ChangeText(string(a))
		file1 += tempfile1 + "\n"
		file2 += tempfile2 + "\n"
	}

	saveFile("i18n.txt", file1)
	saveFile("dec.txt", file2)
}

//ChangeText 测试使用
func ChangeText(text string) (string, string) {
	// text = "validationGroup||string||Defines a set of elements for the validation specified, between groups can be separated by a comma."
	temp := strings.Split(text, "|")
	// tempSign := strings.Trim(temp[3], " ")
	// if *strSign != tempSign {
	// 	*strSign = tempSign
	// 	fmt.Println()
	// }
	result := ""
	result1 := ""
	tempName := ""
	tempName = strings.Title(temp[0]) + "Description"
	result += temp[0] + `: {` + `Type: "` + temp[1] + `",DefaultValue: "",` + `Description: i18n.` + tempName + `},`

	result1 += tempName + `:"` + temp[2] + `",`
	// fmt.Println(result)
	return result, result1
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
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
