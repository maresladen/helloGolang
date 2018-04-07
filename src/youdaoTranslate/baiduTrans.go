package youdaoTranslate

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var appid = "20171116000095862"
var secretKey = "hng4WYKMsidixy1q3B1j"

// var myurl = "http://api.fanyi.baidu.com/api/trans/vip/translate"
// var fromLang = "zh"
// var toLang = "en"

type tranResultModel struct {
	From        string        `json:"from"`
	To          string        `json:"to"`
	TransResult []transResult `json:"trans_result"`
}

type transResult struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type configJSON struct {
	SourceFileName  string `json:"SourceFileName"`
	FromLang        string `json:"FromLang"`
	ToLang          string `json:"ToLang"`
	ProcessCount    int    `json:"ProcessCount"`
	SinglePostCount int    `json:"SinglePostCount"`
	Appid           string `json:"Appid"`
	SecretKey       string `json:"SecretKey"`
	TransURL        string `json:"TransURL"`
}

var cConfig configJSON

//InitTrans 初始化
func InitTrans() {
	readConfig()

	if cConfig.Appid != "" {
		appid = cConfig.Appid
	}
	if cConfig.SecretKey != "" {
		secretKey = cConfig.SecretKey
	}
}

func readConfig() {
	fi, err := os.Open("baiduTrans.json")
	if err != nil {
		writelog(err, "获取config文件失败，请确定启动同级目录下存在baiduTrans.json文件")
	} else {
		temp, _ := ioutil.ReadAll(fi)
		json.Unmarshal(temp, &cConfig)
	}
}

func generateTransText() {

}

func generateURL(strTrans string) string {
	salt := rand.Intn(32768) + 32768
	sign := appid + strTrans + strconv.Itoa(salt) + secretKey
	h := md5.New()
	io.WriteString(h, sign)

	signMD5 := hex.EncodeToString(h.Sum(nil))
	transText, err := url.Parse(strTrans)
	if err != nil {
		println("wrong")
	}

	return cConfig.TransURL + "?appid=" + appid + "&q=" + transText.String() + "&from=" + cConfig.FromLang + "&to=" + cConfig.ToLang + "&salt=" + strconv.Itoa(salt) + "&sign=" + signMD5
}

//TranslateTextByBaidu 百度翻译方法
func TranslateTextByBaidu(url string) {
	client := http.Client{
		Timeout: time.Duration(time.Second * 5),
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("网络访问出错")
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("读取回传内容出错")
	}

	var jData tranResultModel
	json.Unmarshal(data, &jData)

	transResult := jData.TransResult
	for _, tempTransResult := range transResult {
		fmt.Println(tempTransResult.Src)
		fmt.Println("---")
		fmt.Println(tempTransResult.Dst)
	}
}
