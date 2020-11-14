package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	HttpProxy("/", 8080)
}

// 目标Host
func TargetHost() string {
	return "http://localhost:8002"
}

// 转发请求
// 参考：https://www.cnblogs.com/boxker/p/11046342.html
func HttpProxy(pattern string, port int) {
	http.HandleFunc(pattern, doProxy)
	strPort := strconv.Itoa(port)
	fmt.Print("listenning on :", " ", pattern, " ", strPort, "\n")
	err := http.ListenAndServe(":"+strPort, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// 转发请求
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
	reqUrl := TargetHost() + r.URL.String()

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
