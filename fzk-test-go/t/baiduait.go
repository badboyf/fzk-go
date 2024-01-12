package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	var host = "https://aip.baidubce.com/rpc/2.0/nlp/v1/news_summary"
	var accessToken = "24.20b17cefa3d05548cc832d9c2515193b.2592000.1654364445.282335-26122853"
	uri, err := url.Parse(host)
	if err != nil {
		fmt.Println(err)
	}
	query := uri.Query()
	query.Set("access_token", accessToken)
	query.Set("charset", "UTF-8")
	uri.RawQuery = query.Encode()

	var params = map[string]interface{}{}
	params["title"] = "今天天气怎么样"
	params["content"] = "今天天气怎么样"
	params["max_summary_len"] = 150
	sendBody, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
	}
	sendData := string(sendBody)
	client := &http.Client{}
	request, err := http.NewRequest("POST", uri.String(), strings.NewReader(sendData))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(result))
}