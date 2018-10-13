package ESImport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

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

type transferContent struct {
	SrcURL    string `json:"srcURL"`
	TarURL    string `json:"tarURL"`
	SrcIndex  string `json:"srcIndex"`
	TarIndex  string `json:"tarIndex"`
	SrcType   string `json:"srcType"`
	TarType   string `json:"tarType"`
	JoinField string `json:"joinField"`
	TarModule string `json:"tarModule"`
}

var tConfig transferContent

func readConfigTrans() {
	fi, err := os.Open("./configTrans.json")
	if err != nil {
		writelog(err, "get config json data wrong")
	} else {
		temp, err := ioutil.ReadAll(fi)
		if err != nil {
			fmt.Println(err)
			fmt.Scanln(&exitsign)
		}
		json.Unmarshal(temp, &tConfig)
	}
}

//GetParentID 获取父id
func GetParentID() string {

	readConfigTrans()

	ctx := context.Background()

	clientSrc, err := elastic.NewClient(elastic.SetURL(tConfig.SrcURL))
	clientTar, err := elastic.NewClient(elastic.SetURL(tConfig.TarURL))
	if err != nil {
		fmt.Println(err)
		fmt.Scanln(&exitsign)
		panic(err)
	}

	svc := clientSrc.Scroll(tConfig.SrcIndex).Type(tConfig.SrcType).Size(5000)
	bulkService := clientTar.Bulk()
	i := 0
	for {
		res, err := svc.Do(ctx)
		if err == io.EOF {
			fmt.Println("跑完结束了")
			fmt.Scanln(&exitsign)
			break
		}
		if err != nil {
			fmt.Println("上来就跪了")
			fmt.Scanln(&exitsign)
			panic(err)
		}
		i++
		fmt.Println(i)
		if i > 600 {
			fmt.Println("超过范围结束的")
			fmt.Scanln(&exitsign)
			break
		}
		for _, hit := range res.Hits.Hits {
			j, err := json.Marshal(&hit.Source)
			if err != nil {
				fmt.Println("转换json失败")

				fmt.Scanln(&exitsign)
				panic(err)
			}

			jsonStr := string(j)

			if tConfig.JoinField != "" {
				var dat map[string]interface{}
				err := json.Unmarshal(j, &dat)
				if err == nil {
					var temp map[string]string
					temp = make(map[string]string)
					temp["name"] = tConfig.TarModule
					if hit.Parent != "" {
						temp["parent"] = hit.Parent
					}
					dat[tConfig.JoinField] = temp

				} else {
					fmt.Println(err)
					fmt.Scanln(&exitsign)
				}

				jsonStrTemp, _ := json.Marshal(dat)
				jsonStr = string(jsonStrTemp)

				fmt.Println(jsonStr)
				// fmt.Scanln(&exitsign)
			}

			req := elastic.NewBulkIndexRequest().Index(tConfig.TarIndex).Doc(jsonStr).Id(hit.Id)
			if hit.Parent != "" {
				req.Routing(hit.Parent)
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

	fmt.Scanln(&exitsign)
	return "恭喜你"
}

// func GetPolicyData() {
// 	readtempConfig()

// 	// url := esmConfig.HTTPType + "://" + esmConfig.Host + "/ah-proposal/pa/policy/v1/" + pid
// 	// m2 := make(map[string]string)
// 	// domapData(m2)
// 	// m := make(map[string]interface{})
// 	// getRequestData2(url, esmConfig.Host, esmConfig.Auth, m, m2)
// 	// fmt.Println("---------------")
// 	// fmt.Println(m["ProposalNo"])
// 	// fmt.Scanln(&exitsign)

// 	fi, err := os.Open("./pol.csv")
// 	if err != nil {
// 		fmt.Printf("Error: %s\n", err)
// 		return
// 	}
// 	defer fi.Close()

// 	br := bufio.NewReader(fi)
// 	//----------------------
// 	ctx := context.Background()
// 	client, err := elastic.NewClient(elastic.SetURL(esmConfig.ESAddr))
// 	if err != nil {
// 		// Handle error
// 		fmt.Println(err)
// 		fmt.Println("获取client连接失败")
// 		fmt.Scanln(&exitsign)
// 		panic(err)
// 	}

// 	m2 := make(map[string]string)
// 	domapData(m2)

// 	for {
// 		a, _, c := br.ReadLine()
// 		if c == io.EOF {
// 			break
// 		} else {
// 			pid := strings.TrimSpace(string(a))

// 			policyId := "Policy_" + pid

// 			time.Sleep(60 * time.Millisecond)
// 			url := esmConfig.HTTPType + "://" + esmConfig.Host + "/gi-proposal/pa/policy/v1/" + pid
// 			fmt.Println(url)

// 			m := make(map[string]interface{})
// 			getRequestData2(url, esmConfig.Host, esmConfig.Auth, m, m2)

// 			m["PlanCode"] = "buchang3"
// 			m["entity_type"] = "Policy"
// 			m["index_time"] = "2018-09-03"
// 			jsonStr, _ := json.Marshal(m)
// 			fmt.Println("---------------------")
// 			fmt.Println(pid)

// 			get1, err := client.Get().
// 				Index(esmConfig.PolicyIndex).
// 				Type("Policy").
// 				Id(policyId).
// 				Do(ctx)
// 			if err != nil || !get1.Found {
// 				// Handle error
// 				put1, err := client.Index().Index(esmConfig.PolicyIndex).Type("Policy").Id(policyId).BodyString(string(jsonStr)).Do(ctx)
// 				if err != nil {
// 					fmt.Println("do index faild", err.Error())
// 					fmt.Scanln(&exitsign)
// 				} else {
// 					fmt.Println(put1)
// 				}
// 			}
// 		}
// 	}
// }

// //GetWorkflowData testFunction
// func GetWorkflowData() {

// 	readtempConfig()
// 	// fmt.Println("请输入ES地址")
// 	// fmt.Scanln(&esAddress)

// 	// fmt.Println("请输入workflow的index名称")
// 	// fmt.Scanln(&wfIndexName)

// 	// fmt.Println("请输入policy的index名称")
// 	// fmt.Scanln(&policyIndexName)

// 	// fmt.Println("请输入host")
// 	// fmt.Scanln(&hostStr)

// 	// fmt.Println("auth内容")
// 	// fmt.Scanln(&authStr)

// 	ctx := context.Background()

// 	// Create a client
// 	client, err := elastic.NewClient(elastic.SetURL(esmConfig.ESAddr))
// 	if err != nil {
// 		// Handle error
// 		fmt.Println(err)
// 		fmt.Println("获取client连接失败")
// 		fmt.Scanln(&exitsign)
// 		panic(err)
// 	}

// 	// tempQuery := elastic.NewMatchQuery("PolicyId", 5038335199549)
// 	formatTimeStr := "2018-09-15 00:00:00"

// 	formatTime, err := time.Parse("2006-01-02 15:04:05", formatTimeStr)

// 	matchQuery := elastic.NewRangeQuery("index_time").Lte(formatTime)

// 	workflowResult, err := client.Search().Index(esmConfig.WorkflowIndex).Type("Workflow").Query(matchQuery).Sort("index_time", false).From(0).Size(6000).Pretty(true).Do(ctx)

// 	if err != nil {
// 		// Handle error
// 		fmt.Println(err)
// 		fmt.Println("请求失败")

// 		fmt.Scanln(&exitsign)
// 		panic(err)
// 	}

// 	for _, hit := range workflowResult.Hits.Hits {
// 		j, err := json.Marshal(&hit.Source)
// 		if err != nil {
// 			fmt.Println("转换json失败")

// 			fmt.Scanln(&exitsign)
// 			panic(err)
// 		}
// 		var dat map[string]interface{}
// 		m := make(map[string]interface{})
// 		if err := json.Unmarshal(j, &dat); err == nil {
// 			fmt.Println("==============重新设置map内容=======================")
// 			doMapAction(dat, m)

// 			// eid := fmt.Sprintf("%f", m["BusinessEntityId"])
// 			m["id"] = m["BusinessEntityId"]

// 			poId := getPID(m["PolicyId"])
// 			if poId == "" {
// 				poId = getPID(m["BusinessEntityId"])
// 			}

// 			policyId := "Policy_" + poId

// 			fmt.Println(policyId)

// 			m["entity_id"] = policyId
// 			m["entity_type"] = "Policy"
// 			m["Status"] = m["PolicyStatus"]
// 			fmt.Println("==============请求保单对象=======================")

// 			time.Sleep(60 * time.Millisecond)
// 			url := esmConfig.HTTPType + "://" + esmConfig.Host + "/gi-proposal/pa/policy/v1/" + poId
// 			fmt.Println(url)

// 			delete(m, "BusinessEntityId")

// 			getRequestData(url, esmConfig.Host, esmConfig.Auth, m)

// 			m["PlanCode"] = "buchang3"
// 			// "entity_id":PolicyObj.policy_id
// 			jsonStr, _ := json.Marshal(m)
// 			fmt.Println("---------------------")
// 			fmt.Println(policyId)

// 			get1, err := client.Get().
// 				Index(esmConfig.PolicyIndex).
// 				Type("Policy").
// 				Id(policyId).
// 				Do(ctx)
// 			if err != nil || !get1.Found {
// 				// Handle error
// 				put1, err := client.Index().Index(esmConfig.PolicyIndex).Type("Policy").Id(policyId).BodyString(string(jsonStr)).Do(ctx)
// 				if err != nil {
// 					fmt.Println("do index faild", err.Error())
// 					fmt.Scanln(&exitsign)
// 				} else {
// 					fmt.Println(put1)
// 				}
// 			} else {
// 				// put1, err := client.Index()
// 				put1, err := client.Update().Index(esmConfig.PolicyIndex).Type("Policy").Id(policyId).Doc(m).Do(ctx)
// 				if err != nil {
// 					fmt.Println("update index faild", err.Error())
// 					fmt.Scanln(&exitsign)
// 				} else {
// 					fmt.Println(put1)
// 				}
// 			}

// 		}

// 	}
// 	// tempSearchHits := workflowResult.Hits
// 	// var tempMap map[string]interface{}
// 	// err = json.Unmarshal(tempSearchHits.Hits[0].Source, &tempMap)
// 	// if err != nil {
// 	// 	fmt.Println(string(tempMap))
// 	// } else {
// 	// 	fmt.Println("no")
// 	// }

// 	// for _, hit := range tempSearchHits.Hits {
// 	// 	fmt.Println(hit.Id)
// 	// 	fmt.Println(len(hit.Fields))
// 	// }

// 	// TotalHits is another convenience function that works even when something goes wrong.

// 	fmt.Println("exit?")
// 	fmt.Scanln(&exitsign)
// 	// fmt.Scanf("%s %s", &firstName, &lastName)

// }
