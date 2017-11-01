package stringEditor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

//DoStringEditor 字符串处理
func DoStringEditor() {

	fi, err := os.Open("stringeditor.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)

	index := 0
	strSign := ""
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		index++
		ChangeText(index, string(a), &strSign)
	}
}

//ChangeText 测试使用
func ChangeText(index int, text string, strSign *string) {
	// text = "validationGroup||string||Defines a set of elements for the validation specified, between groups can be separated by a comma."
	temp := strings.Split(text, "||")
	tempSign := strings.Trim(temp[3], " ")
	if *strSign != tempSign {
		*strSign = tempSign
		fmt.Println()
		index = 1
	}
	result := ""
	result += `{id:` + strconv.Itoa(index) + `, Name:"` + temp[0] + `", ` + `Type:"` + temp[1] + `",` + `Description:"i18n.` + temp[3] + temp[0] + "Description},"

	fmt.Println(result)
}
