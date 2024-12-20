# 在线词典

该小项目是字节跳动青训营GO入门课程中的实战项目<https://juejin.cn/course/bytetech/7140987981803814919/section/7141265019756675103>

## 需求分析

1. 执行程序时，在命令行传入单词

2. 根据在线词典给出该单词的音标和注释————利用第三方api进行查询

## 项目亮点

1. 使用http发送请求，解析json

2. 使用代码生成提升开发效率

## 实现思路

1. 抓包

可以通过浏览器开发者工具查看浏览器和服务器之间的请求响应，彩云小译的翻译请求接口为`https://api.interpreter.caiyunai.com/v1/dict`，这里要注意这个接口还有一个OPTIONS类型的请求，我们要选择的是POST类型的请求，从负载中可以看到该请求有两个json格式的参数，source表示待翻译的英文，trans_type表示要从哪种语言翻译为哪种语言，从响应中可以看到该请求响应的json格式数据

![post请求](../images/1731919360806.jpg)

![post负载](../images/fuzai.jpg)

![post响应](../images/xiangying.jpg)

2. 代码生成构建请求

由于请求非常复杂，用代码构建的话非常麻烦，可以使用代码生成方式来构建请求，首先需要右键请求复制为cURL(bash)，于是就得到下面这一段bash命令

```bash
curl 'https://api.interpreter.caiyunai.com/v1/dict' \
  -H 'accept: application/json, text/plain, */*' \
  -H 'accept-language: zh' \
  -H 'app-name: xiaoyi' \
  -H 'authorization: Bearer' \
  -H 'content-type: application/json;charset=UTF-8' \
  -H 'device-id: 5a1713039eefd90fca8064c503d00a26' \
  -H 'origin: https://fanyi.caiyunapp.com' \
  -H 'os-type: web' \
  -H 'os-version;' \
  -H 'priority: u=1, i' \
  -H 'referer: https://fanyi.caiyunapp.com/' \
  -H 'sec-ch-ua: "Microsoft Edge";v="131", "Chromium";v="131", "Not_A Brand";v="24"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Windows"' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: cross-site' \
  -H 'user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0' \
  -H 'x-authorization: token:qgemv4jr1y38jyq6vhvi' \
  --data-raw '{"trans_type":"en2zh","source":"good"}'
```

然后就可以到代码生成网站`https://curlconverter.com/go/`将上面的bash输入之后就可以得到使用go语言来编写得到的请求了，极大的减少了构建http请求的工作量

![代码生成](../images/daimashengcheng.jpg)

创建请求的时候第三个参数需要为一个流，因此需要使用`strings.NewReader(`{"trans_type":"en2zh","source":"good"}`)`将字符串转换成流，这是因为body可能是一个很大的字符串，如果直接使用字符串的话会导致非常大的内存开销，因此使用流来传输数据，这样就可以占用很少的内存，然后流式创建请求

响应的body同样是一个流，在go中，为了避免资源泄露，需要加一个defer来手动关闭这个流

> defer：会在函数结束之后，从下往上触发

```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	client := &http.Client{}
	var data = strings.NewReader(`{"trans_type":"en2zh","source":"good"}`)
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
    // 关闭流
	defer resp.Body.Close()
	// 读取响应
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
}
```

3. 生成request body

上面的请求中，body的数据是固定的，而我们需要用一个变量来作为body的输入，因此需要用json序列化

json序列化：需要构造一个结构体，使得其字段和json字段一一对应，然后直接调用json.Marshal即可

```go
type DictRequest struct {
	TransType string `json:"trans_type"`
	Source    string `json:"source"`
	UserID    string `json:"user_id"`
}

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
```

4. 解析response body

从响应中提取指定的字段，go是一门强类型的语言，在js或python等脚本语言中，这个body返回的是一个字典或者叫map的结构，可以直接用[]加点去取值，但是这不是go中的最佳实践，最常见的是写一个结构体，字段一一对应，然后反序列化到结构体中，但是通常返回的字段非常复杂，这种实现非常容易出错，因此可以继续使用代码生成的方法，使用在线网站`https://mholt.github.io/json-to-go/`即可自动生成对应的结构体

![json to go](../images/json2go.jpg)

```go
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
	// 解析response
	var dictResponse DictResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", dictResponse)
```

5. 打印结果

结构体很大，其中只有几个是我们所需要的，因此我们选择性的将结构体中的翻译和解释打印

```go
	fmt.Println(request.Source, "UK: ", dictResponse.Dictionary.Prons.En, "US: ", dictResponse.Dictionary.Prons.EnUs)
	for _, item := range dictResponse.Dictionary.Explanations {
		fmt.Println(item)
	}
```

6. 完善代码

将在线字典的功能独立成一个函数，然后在main函数中调用，并且将命令行参数传入

```go
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
example: simpleDict hello
		`)
		os.Exit(1)
	}
	word := os.Args[1]
	query(word)
}
```

## 总结

1. 利用代码生成的方式，借助在线网站工具，能够大大减少构建http请求的工作量

2. 学习了使用go发送http请求并解析json的基本流程，首先要抓包找到请求的url和请求头，然后使用代码生成工具生成请求，对于响应的json数据，需要通过代码生成的方式，借助在线网站将json数据转化成go的结构体，然后通过反序列化将json数据填充到结构体中

3. 学习了从命令行读取数据的方式，通过os.Args[]可以获取到命令行输入
