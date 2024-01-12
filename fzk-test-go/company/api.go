package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)


func main() {

	// 测试
	projectId := "61a860b84e8ef87926e63a0f"
	appId := "61a9b53af32af13eec30c63f"
	// prod
	//projectId = "61b056caff5f093a472e87bf"
	//appId = "61b17c289e6cb947db56bae5"

	ids := []int{ 178 }
	for i := range ids {
		id := strconv.Itoa(ids[i])
		converApi(id, projectId, appId, true, true)

		fmt.Printf("\n\n\n")
	}

	//getOtherId(strconv.Itoa(178))
}

func converApi(id string, projectId string, appId string, onlyPrintExt bool, createApi bool) (map[string]interface{}) {
	otherId := getOtherId(id);
	if otherId == "" {
		fmt.Printf("%v otherId is null", id)
		return nil
	}
	api := cJuheScrawApi(id, otherId)
	api["projectId"] = projectId
	api["appId"] = appId

	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", "  ")
	if onlyPrintExt {
		jsonEncoder.Encode(api["ext"])
	} else {
		jsonEncoder.Encode(api)
	}

	fmt.Printf("%v\n%v\n", id, bf.String())
	ioutil.WriteFile(fmt.Sprintf("zfile/%v.json", id), bf.Bytes(), 0666)

	if createApi {
		doCreateApi(api, projectId)
	}

	return api
}

func doCreateApi(api map[string]interface{}, prjId string)  {
	client := &http.Client{}
	bodyByte, err := json.Marshal(api)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return
	}
	req, err := http.NewRequest("POST", "https://test-data-portal-api.csdn.net/portal/v1/interface/create", bytes.NewReader(bodyByte))
	if err != nil {
		fmt.Printf("error  %v \n", err)
		return
	}
	req.Header.Add("Cookie", "x_data_token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJjcmVkZW50aWFsIjoiZjZmZGZmZTQ4YzkwOGRlYjBmNGMzYmQzNmMwMzJlNzIiLCJ1c2VySWQiOiJhZG1pbiJ9.dAhjQLQ1YTDoJYN43eekN9K-3SY3ZkAKXz8UAy6OcgM;")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Data-Project-Id", prjId)
	req.Header.Add("Accept", "application/json, text/plain, */*")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error  %v \n", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error  %v \n", err)
		return
	}
	var m map[string]interface{}
	err = json.Unmarshal(resBody, &m)
	fmt.Printf("%v\n", m)
}

func getOtherId(id string) string {
	//idMap := map[string]string{}
	//idMap["103"] = "309"
	//idMap["248"] = "966"
	//idMap["249"] = "967"
	//idMap["54"] = "121"
	//idMap["208"] = "741"
	//idMap["339"] = "1189"
	//idMap["188"] = "678"
	//idMap["207"] = "724"
	//idMap["213"] = "734"
	//idMap["299"] = "1113"
	//idMap["43"] = "99"
	//idMap["33"] = "76"
	//idMap["99"] = "269"
	//idMap["498"] = "1381"
	//idMap["387"] = "1278"
	//idMap["386"] = "1277"
	//idMap["287"] = "1071"
	//idMap["395"] = "1286"
	//idMap["153"] = "492"
	//
	//rs := idMap[id]
	url := fmt.Sprintf("https://www.juhe.cn/docs/api/id/%v", id)
	response, err := http.Get(url); if err != nil {
		fmt.Printf("error %v \n", err)
		return ""
	}
	doc, err := goquery.NewDocumentFromReader(response.Body); if err != nil {
		fmt.Printf("error %v \n", err)
		return ""
	}

	aid := ""
	doc.Find("ul[class=api-sub-list] li").Each(func(i int, li *goquery.Selection) {
		if aid != "" {
			return
		}
		attr, exists := li.Find("a").Attr("href")
		if exists && strings.Contains(attr, "/docs/api/id/") {
			aid = attr[strings.LastIndex(attr, "/")+1:]
		}
	})

	return aid
}

func cJuheScrawApi(id string, other string) (map[string]interface{}) {
	name, mark := cJuheApiInfo("https://www.juhe.cn/docs/api/packages/" + id)
	htmlUrl := fmt.Sprintf("https://www.juhe.cn/docs/api/id/%v/aid/%v/%v", id, other, id)
	url, method, params, resParams, respCode := sJuheUrl(htmlUrl)

	apiPath := strings.Replace(url, "http://v.juhe.cn", "", -1)
	apiPath = strings.Replace(apiPath, "https://v.juhe.cn", "", -1)
	apiPath = strings.Replace(apiPath, ".php", "", -1)
	index := strings.LastIndex(apiPath, "/")
	index = strings.LastIndex(apiPath[:index], "/")
	apiPath = apiPath[index:]
	apiPath = strings.ToLower(apiPath)

	apiUrl := url + "?key=${token}"
	ext := map[string]interface{}{}
	extParams := []interface{}{}
	for i := range params {
		item := map[string]string{}
		if params[i][0] == "key" {
			continue
		} else if params[i][0] == "dtype" {
			apiUrl += "&dtype=json"
			continue
		}
		item["name"] = params[i][0]
		item["source"] = "QUERY"
		if params[i][2] == "int" {
			item["paramType"] = "NUMBER"
		} else {
			item["paramType"] = "STRING"
		}
		if params[i][1] == "是" {
			item["required"] = "true"
		} else {
			item["required"] = "false"
		}
		item["desc"] = params[i][3]
		extParams = append(extParams, item)
		apiUrl += "&" + fmt.Sprintf("%v=${query.%v}", params[i][0], params[i][0])
	}

	tokenParam := map[string]string{}
	tokenParam["name"] = "token"
	tokenParam["source"] = "QUERY"
	tokenParam["paramType"] = "STRING"
	tokenParam["required"] = "true"
	tokenParam["desc"] = "query参数，接口token，在我的api列表获取"
	extParams = append(extParams, tokenParam)

	extProxyParam := map[string]interface{}{}
	extProxyParam["output"] = "{\"code\": <#if result.error_code==0>200<#else>${result.error_code}</#if>,\"msg\": \"${result.reason}\",\"data\": ${JSONObject.toJSONString(result.result)}}"
	extProxyParam["outputType"] = "JSON"
	extProxyParam["goodsImg"] = fmt.Sprintf("https://juheimg.oss-cn-hangzhou.aliyuncs.com/www/api/a%v.png", id)
	extProxyParam["chargeType"] = "PAY"
	extProxyParam["goodsDesc"] = mark
	extProxyParam["remoteToken"] = getRemoteToken(id)
	extProxyParam["outputContentDisposition"] = nil
	extProxyParam["outputContentType"] = nil

	remoteParam := map[string]interface{}{}
	remoteParam["url"] = apiUrl
	remoteParam["connectTimeout"] = "3000"
	remoteParam["readTimeout"] = "6000"
	remoteParam["header"] = map[string]string{}
	remoteParam["body"] = nil
	remoteParam["dataType"] = "JSON"
	remoteParam["resultType"] = "JSON"
	if strings.Contains(method, "get") {
		remoteParam["method"] = "GET"
	} else {
		remoteParam["method"] = "POST"
	}
	extProxyParam["remoteParam"] = remoteParam


	exampleOutput := []interface{}{}
	for i := range resParams {
		item := map[string]string{}
		key := resParams[i][0]
		if key == "error_code" {
			key = "code"
		} else if key == "reason" {
			key = "msg"
		} else if key == "result" {
			key = "data"
		} else if key == "" {
			continue
		}
		item["name"] = key
		if resParams[i][1] == "int" {
			item["paramType"] = "NUMBER"
		} else {
			item["paramType"] = "STRING"
		}
		item["desc"] = resParams[i][2]
		exampleOutput = append(exampleOutput, item)
	}
	extProxyParam["exampleOutput"] = exampleOutput

	respCode = strings.Replace(respCode, "\"error_code\": 0", "\"code\": 200", -1)
	respCode = strings.Replace(respCode, "返回值为0才计费", "", -1)
	respCode = strings.Replace(respCode, "返回0时计费一次", "", -1)
	respCode = strings.Replace(respCode, "\"error_code\":", "\"code\":", -1)
	respCode = strings.Replace(respCode, "\"reason\":", "\"msg\":", -1)
	respCode = strings.Replace(respCode, "\"result\":", "\"data\":", -1)
	if id == "54" {
		respCode = "{\"code\":200,\"msg\":\"短信发送成功\",\"data\":{\"count\":1,\"fee\":1,\"sid\":\"23d6bc4913614919a823271d820662af\"}}"
	}
	var exampleOutputCodeMap map[string]interface{}
	err := json.Unmarshal([]byte(respCode), &exampleOutputCodeMap)
	if err != nil {
		extProxyParam["exampleOutputCode"] = respCode
	} else {
		for s := range exampleOutputCodeMap {
			realKey := s
			v := exampleOutputCodeMap[s]
			if s == "error_code" {
				realKey = "code"
				v = 200
			} else if s == "reason" {
				realKey = "msg"
			} else if s == "result" {
				realKey = "data"
			} else if s == "resultcode" {
				continue
			}
			delete(exampleOutputCodeMap, s)
			exampleOutputCodeMap[realKey] = v
		}
		extProxyParam["exampleOutputCode"] = exampleOutputCodeMap
	}

	ext["params"] = extParams
	ext["proxyParam"] = extProxyParam
	ext["needPersonalCertificate"] = "true"
	ext["needCompanyCertificate"] = "true"

	rs := map[string]interface{}{}
	rs["title"] = name
	rs["desc"] = mark
	rs["path"] = apiPath
	rs["methods"] = []string{"GET", "POST"}
	rs["typeCode"] = 2
	rs["ext"] = ext

	return rs
}

func cJuheApiInfo(url string) (string, string) {
	response, err := http.Get(url); if err != nil {
		fmt.Printf("error %v \n", err)
		return "", ""
	}
	resBody, err := ioutil.ReadAll(response.Body)
	var m map[string]interface{}
	err = json.Unmarshal(resBody, &m)
	//fmt.Printf("%v", m)
	result := m["result"].(map[string]interface{})
	apiData := result["apiData"].(map[string]interface{})
	name := apiData["name"].(string)
	marks := apiData["marks"].(string)
	return name, marks
}

func sJuheUrl(url string) (string, string, [][]string, [][]string, string) {
	response, err := http.Get(url); if err != nil {
		fmt.Printf("error %v \n", err)
		return "", "", nil, nil, ""
	}
	resBody, err := ioutil.ReadAll(response.Body)
	var m map[string]interface{}
	err = json.Unmarshal(resBody, &m)
	result := m["result"].(map[string]interface{})
	html := result["html"].(string)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html)); if err != nil {
		fmt.Printf("error %v \n", err)
		return "", "", nil, nil, ""
	}
	apiUrl := doc.Find("div[class=simpleline]:nth-child(1)").Text()
	apiUrl = apiUrl[15:]
	method := doc.Find("div[class=simpleline]:nth-child(3)").Text()

	var p [][]string
	var resp [][]string
	var respCode string
	doc.Find("div[class=simpleTable]").Each(func(i int, table *goquery.Selection) {
		if strings.Contains(table.Text(), "请求参数说明") { // 参数
			table.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
				key := tr.Find("td:nth-child(2)").Text()
				key = strings.TrimSpace(key)
				if key == "" {
					return
				}
				require := tr.Find("td:nth-child(3)").Text()
				typ := tr.Find("td:nth-child(4)").Text()
				desc := tr.Find("td:nth-child(5)").Text()
				p = append(p, []string{key, require, typ, desc})
			})
		}
		if strings.Contains(table.Text(), "返回参数说明") { // 返回参数说明
			table.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
				key := tr.Find("td:nth-child(2)").Text()
				typ := tr.Find("td:nth-child(3)").Text()
				desc := tr.Find("td:nth-child(4)").Text()
				resp = append(resp, []string{key, typ, desc})
			})
		}
		if strings.Contains(table.Text(), "JSON返回示例") { // JSON返回示例
			table.Find("code").Each(func(i int, tr *goquery.Selection) {
				content := tr.Text()
				index := strings.Index(content, "{")
				respCode = content[index:]
			})
		}

	})

	//fmt.Printf("%v \n%v \n%v \n", apiUrl, method, p)
	return apiUrl, method, p, resp, respCode
}

func getRemoteToken(id string) string {
	content, err := ioutil.ReadFile("zfile/token.json");
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return ""
	}
	var a []map[string]interface{}
	err = json.Unmarshal(content, &a)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return ""
	}
	idInt, _ := strconv.Atoi(id)
	for i := range a {
		s := fmt.Sprintf("%v", a[i]["did"])
		idItem, _ := strconv.Atoi(s)
		if idItem == idInt {
			return a[i]["appkey"].(string)
		}
	}
	return ""
}