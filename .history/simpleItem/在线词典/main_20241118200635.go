package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type DictRequest struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

type DictResponse struct {
	Rc   int `json:"rc"`
	Wiki struct {
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string   `json:"explanations"`
		Synonym      []string   `json:"synonym"`
		Antonym      []string   `json:"antonym"`
		WqxExample   [][]string `json:"wqx_example"`
		Entry        string     `json:"entry"`
		Type         string     `json:"type"`
		Related      []any      `json:"related"`
		Source       string     `json:"source"`
	} `json:"dictionary"`
}

func main() {
	client := &http.Client{}
	// var data = strings.NewReader(`{"trans_type":"en2zh","source":"good"}`)
	request := DictRequest{
		TransType: "en2zh",
		Source:    "good",
	}
	// 序列化结构体，变成一个buf数组
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)

	// 创建请求
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	// 设置请求头
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh")
	req.Header.Set("app-name", "xiaoyi")
	req.Header.Set("authorization", "Bearer")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("device-id", "5a1713039eefd90fca8064c503d00a26")
	req.Header.Set("origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("os-type", "web")
	req.Header.Set("os-version", "")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("sec-ch-ua", `"Microsoft Edge";v="131", "Chromium";v="131", "Not_A Brand";v="24"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0")
	req.Header.Set("x-authorization", "token:qgemv4jr1y38jyq6vhvi")
	// 发起请求
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// 读取响应存储到内存中
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// 打印出来既是响应的json数据
	// fmt.Printf("%s\n", bodyText)
	// 解析response
	var dictResponse DictResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%#v\n", dictResponse)

	fmt.Println(request.Source, "UK: ", dictResponse.Dictionary.Prons.En, "US: ", dictResponse.Dictionary.Prons.EnUs)
	for _, item := range dictResponse.Dictionary.Explanations {
		fmt.Println(item)
	}
}
