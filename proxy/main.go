package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"strings"
)

// 全局变量： 配置
var config ProxyConfig

// 全局变量： 代理的URL信息
var urlInfo []ProxyUrlInfo

// 全局变量： 不使用代理的URL信息
var urlNoProxy []ProxyUrlInfo

func main() {

	var err error

	// 加载配置
	config, err = LoadConfig(configFile)
	if err != nil {
		log.Print(err)
	}
	// log.Println(config)

	// 读取¥代理的URL信息
	urlInfo, err = LoadUrlInfo(urlFile)
	if err != nil {
		log.Print(err)
	}
	// log.Println(urlInfo)

	for _, v := range urlInfo {
		if !v.UseProxy {
			// 不使用代理的URL信息
			urlNoProxy = append(urlNoProxy, v)
		}
	}

	// 响应函数
	for _, v := range urlInfo {
		if v.UseProxy {
			// 使用代理
			http.HandleFunc(v.Url, doProxy)
		} else {
			// 不使用代理
			http.HandleFunc(v.Url, doHandle)
		}
	}

	// 监听
	port := ":" + config.ProxyPort
	log.Printf("Listen and serve [%v]", port)
	http.ListenAndServe(port, nil)
}

// 转发请求函数
// 参考：https://www.cnblogs.com/boxker/p/11046342.html
func doProxy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("url: ", r.URL)

	// 创建一个HttpClient用于转发请求
	cli := &http.Client{}

	// 读取请求的Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Print("ioutil.ReadAll(r.Body) ", err.Error())
	}

	// 转发的URL
	reqUrl := config.TargetHost + r.URL.String()

	// 创建转发用的请求
	req, err := http.NewRequest(r.Method, reqUrl, strings.NewReader(string(body)))
	if err != nil {
		fmt.Print("http.NewRequest ", err.Error())
		return
	}

	// 转发请求的表头
	for k, v := range r.Header {
		req.Header.Set(k, v[0])
	}

	// 发起请求
	res, err := cli.Do(req)
	if err != nil {
		fmt.Print("cli.Do(req) ", err.Error())
		return
	}
	defer res.Body.Close()

	// 响应Header
	for k, v := range res.Header {
		w.Header().Set(k, v[0])
	}

	// 复制转发的响应Body到响应Body
	io.Copy(w, res.Body)
}
