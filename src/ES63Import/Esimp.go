package ES63Import

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	elastic "github.com/olivere/elastic"
)

var esmConfig mconfig
var exitsign string

type mconfig struct {
	HTTPType      string `json:"HttpType"`
	Host          string `json:"Host"`
	Auth          string `json:"Auth"`
	ESAddr        string `json:"ESAddr"`
	WorkflowIndex string `json:"WorkflowIndex"`
	PolicyIndex   string `json:"PolicyIndex"`
}

// type transferContent struct {
// 	SrcURL    string `json:"srcURL"`
// 	TarURL    string `json:"tarURL"`
// 	SrcIndex  string `json:"srcIndex"`
// 	TarIndex  string `json:"tarIndex"`
// 	SrcType   string `json:"srcType"`
// 	TarType   string `json:"tarType"`
// 	JoinField string `json:"joinField"`
// 	TarModule string `json:"tarModule"`
// }

type transferContent struct {
	SrcURL  string `json:"srcURL"`
	TarURL  string `json:"tarURL"`
	Indexes []struct {
		SrcIndex  string   `json:"srcIndex"`
		TarIndex  string   `json:"tarIndex"`
		SrcType   string   `json:"srcType"`
		TarType   string   `json:"tarType"`
		IDKey     string   `json:"idKey,omitempty"`
		ParentKey string   `json:"parentKey,omitempty"`
		TarModule string   `json:"tarModule,omitempty"`
		JoinField string   `json:"joinField,omitempty"`
		Columns   []string `json:"columns,omitempty"`
	} `json:"indexes"`
}

var tConfig transferContent

func readConfigTrans() {
	fi, err := os.Open("./configTrans.json")
	if err != nil {
		// writelog(err, "get config json data wrong")
		fmt.Println("get config json data wrong")
	} else {
		temp, err := ioutil.ReadAll(fi)
		if err != nil {
			fmt.Println(err)
			fmt.Scanln(&exitsign)
		}
		json.Unmarshal(temp, &tConfig)
	}
}

//下载
func downloadFiles(content string, indexName string, index int) {

	// CreateFloder("download")
	file, _ := os.Create(strconv.Itoa(index) + "." + indexName)

	defer file.Close()
	io.WriteString(file, content)

}

//CreateFloder 建立文件夹
func CreateFloder(fName string) {
	err := os.Chdir(fName)
	if err != nil {
		os.Mkdir(fName, 0777)
	}
}

//Contain 判断集合是否包含内容
func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

func mergeText(extendName string) {
	outFileName := extendName + "_merge_result" + ".txt"

	outFile, openErr := os.OpenFile(outFileName, os.O_CREATE|os.O_WRONLY, 0755)
	if openErr != nil {
		fmt.Printf("Can not open file %s", outFileName)
	}
	bWriter := bufio.NewWriter(outFile)
	filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		fmt.Println("Processing:", path)
		//这里是文件过滤器，表示我仅仅处理txt文件
		if strings.HasSuffix(path, extendName) {
			fp, fpOpenErr := os.Open(path)
			if fpOpenErr != nil {
				fmt.Printf("Can not open file %v", fpOpenErr)
				return fpOpenErr
			}
			bReader := bufio.NewReader(fp)
			for {
				buffer := make([]byte, 1024)
				readCount, readErr := bReader.Read(buffer)
				if readErr == io.EOF {
					break
				} else {
					bWriter.Write(buffer[:readCount])
				}
			}
		}
		return err
	})
	bWriter.Flush()
}

//EsDataTrans 数据迁移
func EsDataTrans() {

	readConfigTrans()

	ctx := context.Background()

	clientSrc, err := elastic.NewClient(elastic.SetURL(tConfig.SrcURL))
	clientTar, err := elastic.NewClient(elastic.SetURL(tConfig.TarURL))
	if err != nil {
		fmt.Println(err)
		fmt.Scanln(&exitsign)
		panic(err)
	}

	for _, indexConfig := range tConfig.Indexes {

		svc := clientSrc.Scroll(indexConfig.SrcIndex).Type(indexConfig.SrcType).Size(5000)
		bulkService := clientTar.Bulk()
		for {
			res, err := svc.Do(ctx)
			if err == io.EOF {
				fmt.Println(indexConfig.SrcIndex+"-"+indexConfig.SrcType, " 内容跑完")
				break
			}
			if err != nil {
				fmt.Println("上来就跪了", err.Error())
				fmt.Scanln(&exitsign)
				panic(err)
			}

			for _, hit := range res.Hits.Hits {
				j, err := json.Marshal(&hit.Source)
				if err != nil {
					fmt.Println("转换json失败")

					fmt.Scanln(&exitsign)
					panic(err)
				}

				jsonStr := string(j)
				var dat map[string]interface{}
				var fData map[string]interface{}
				fData = make(map[string]interface{})
				err = json.Unmarshal(j, &dat)
				if err == nil {
					if indexConfig.JoinField != "" {

						var temp map[string]interface{}
						temp = make(map[string]interface{})
						temp["name"] = indexConfig.TarModule

						if hit.Parent != "" {
							temp["parent"] = hit.Parent
						} else {
							if indexConfig.ParentKey != "" {
								temp["parent"] = "transVirtual"
							}
						}

						fData[indexConfig.JoinField] = temp
					}
					if indexConfig.Columns != nil {

						for _, v := range indexConfig.Columns {
							if val, ok := dat[v]; ok {
								fData[v] = val
							}
						}

					}
				} else {
					fmt.Println(err)
					fmt.Scanln(&exitsign)
				}

				jsonStrTemp, err := json.Marshal(fData)
				if err != nil {
					fmt.Println("map transError", err)
					fmt.Scanln(&exitsign)
				}
				jsonStr = string(jsonStrTemp)

				req := elastic.NewBulkIndexRequest().Index(indexConfig.TarIndex).Type(indexConfig.TarType).Doc(jsonStr).Id(hit.Id)
				if hit.Parent != "" {
					req.Routing(hit.Parent)
				} else {
					if indexConfig.JoinField != "" {
						if indexConfig.ParentKey != "" {
							req.Routing("transVirtual")
						}
					}
				}
				bulkService.Add(req)

			}

			rep, err := bulkService.Do(ctx)
			if err != nil {
				fmt.Println(err)
				fmt.Scanln(&exitsign)
			} else if rep.Errors {
				for _, item := range rep.Items {
					for _, v := range item {
						fmt.Println(v.Error.Reason)
					}
				}

				fmt.Println("rep err")
				fmt.Scanln(&exitsign)
			}
		}

	}

	// if err != nil {
	// 	panic(err)
	// }
	// for _, hit := range searchResult.Hits.Hits {
	// 	j, err := json.Marshal(&hit.Source)
	// 	if err != nil {
	// 		fmt.Println("转换json失败")

	// 		fmt.Scanln(&exitsign)
	// 		panic(err)
	// 	}
	// 	var dat map[string]interface{}
	// 	if err := json.Unmarshal(j, &dat); err == nil {
	// 		fmt.Println(dat)
	// 	}
	// }

	fmt.Println("恭喜你")
	fmt.Scanln(&exitsign)
}

//DownLoadData 下载
func DownLoadData() {

	readConfigTrans()

	ctx := context.Background()

	clientSrc, err := elastic.NewClient(elastic.SetURL(tConfig.SrcURL))
	if err != nil {
		fmt.Println(err)
		fmt.Scanln(&exitsign)
		panic(err)
	}

	for _, indexConfig := range tConfig.Indexes {

		outFileName := indexConfig.SrcIndex + "_merge_result" + ".txt"

		outFile, openErr := os.OpenFile(outFileName, os.O_CREATE|os.O_WRONLY, 0755)
		if openErr != nil {
			fmt.Printf("Can not open file %s", outFileName)
		}
		bWriter := bufio.NewWriter(outFile)
		fmt.Println(indexConfig.SrcIndex, "开始")
		svc := clientSrc.Scroll(indexConfig.SrcIndex).Type(indexConfig.SrcType).Size(5000)
		i := 0
		for {
			res, err := svc.Do(ctx)
			if err == io.EOF {
				fmt.Println(indexConfig.SrcIndex + "跑完了")
				break
			}
			if err != nil {
				fmt.Println("上来就跪了", err.Error())
				fmt.Scanln(&exitsign)
				panic(err)
			}
			saveDoc := ""

			for _, hit := range res.Hits.Hits {
				j, err := json.Marshal(&hit.Source)
				if err != nil {
					fmt.Println("转换json失败")

					fmt.Scanln(&exitsign)
					panic(err)
				}

				jsonStr := string(j)

				if indexConfig.JoinField != "" {
					var dat map[string]interface{}
					err := json.Unmarshal(j, &dat)
					if err == nil {
						var temp map[string]string
						temp = make(map[string]string)
						temp["name"] = indexConfig.TarModule
						if hit.Parent != "" {
							temp["parent"] = hit.Parent
						} else {
							if indexConfig.ParentKey != "" {
								temp["parent"] = "transVirtual"
							}
						}
						dat[indexConfig.JoinField] = temp

					} else {
						fmt.Println(err)
						fmt.Scanln(&exitsign)
					}

					jsonStrTemp, _ := json.Marshal(dat)
					jsonStr = string(jsonStrTemp)

					saveDoc += jsonStr + "\n"

				} else {
					saveDoc += jsonStr + "\n"
				}

			}
			i++
			// downloadFiles(saveDoc, indexConfig.SrcIndex, i)
			bWriter.Write([]byte(saveDoc))

			fmt.Println("开始下载")
			fmt.Println((i)*5000, "条")

		}
		bWriter.Flush()

		// mergeText(indexConfig.SrcIndex)
	}

	// if err != nil {
	// 	panic(err)
	// }
	// for _, hit := range searchResult.Hits.Hits {
	// 	j, err := json.Marshal(&hit.Source)
	// 	if err != nil {
	// 		fmt.Println("转换json失败")

	// 		fmt.Scanln(&exitsign)
	// 		panic(err)
	// 	}
	// 	var dat map[string]interface{}
	// 	if err := json.Unmarshal(j, &dat); err == nil {
	// 		fmt.Println(dat)
	// 	}
	// }

	fmt.Println("恭喜你")
	fmt.Scanln(&exitsign)
}

//UploadEsData 上传完成
func UploadEsData() {

	readConfigTrans()
	var isExit string
	fmt.Println("请检查要注入数据的URL地址和index,导入地址以第一个配置项为准")
	fmt.Println("导入URL: ", tConfig.TarURL)
	fmt.Println("导入index: ", tConfig.Indexes[0].TarIndex)
	fmt.Println("导入Type: ", tConfig.Indexes[0].TarType)
	fmt.Println("输入N退出，其他任意键继续")
	fmt.Scanln(&isExit)
	if isExit == "n" || isExit == "N" {
		return
	}
	var fileName string
	fmt.Println("请输入要上传的文件名")
	fmt.Scanln(&fileName)

	fi, err := os.Open(fileName)
	if err != nil {
		fmt.Println("打开文件失败", err)
		fmt.Scanln(&isExit)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)

	//----------------------
	ctx := context.Background()

	client, err := elastic.NewClient(elastic.SetURL(tConfig.TarURL))

	bulkService := client.Bulk()
	iCount := 0
	for {

		a, _, c := br.ReadLine()
		if c == io.EOF {
			if iCount%5000 != 0 {
				rep, err := bulkService.Do(ctx)
				if err != nil {
					fmt.Println(err)
					fmt.Scanln(&exitsign)
				} else if rep.Errors {
					for _, item := range rep.Items {
						for _, v := range item {
							fmt.Println(v.Error.Reason)
						}
					}

					fmt.Println("rep err")
					fmt.Scanln(&exitsign)
				}
			}
			break
		} else {

			iCount++
			jsonStr := string(a)
			var dat map[string]interface{}
			err = json.Unmarshal(a, &dat)
			tempID := dat[tConfig.Indexes[0].IDKey]

			req := elastic.NewBulkIndexRequest().Index(tConfig.Indexes[0].TarIndex).Type(tConfig.Indexes[0].TarType).Doc(jsonStr).Id(tempID.(string))
			if tConfig.Indexes[0].ParentKey != "" {
				req.Routing(dat[tConfig.Indexes[0].ParentKey].(string))
			}
			bulkService.Add(req)

			if iCount%5000 == 0 {
				rep, err := bulkService.Do(ctx)
				if err != nil {
					fmt.Println(err)
					fmt.Scanln(&exitsign)
				} else if rep.Errors {
					for _, item := range rep.Items {
						for _, v := range item {
							fmt.Println(v.Error.Reason)
						}
					}

					fmt.Println("rep err")
					fmt.Scanln(&exitsign)
				}
				bulkService = client.Bulk()
			}
		}
	}

	fmt.Println("导入完成")
}
