package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/bettersun/moist"
)

// 全局变量： 配置
var config ProxyConfig

// 全局变量： 代理信息
var proxyInfo ProxyInfo

// 全局变量：HTTP服务
var server http.Server

var isRunning bool

// 启动服务(命令行调用)
func RunServerCommand() {

	RunServer()

	// 命令行执行时，需要下面这段代码
	// 从UI执行时，需要注释掉下面这段代码
	ch := make(chan string)
	<-ch
}

// 启动服务(界面调用)
func RunServer() {

	var err error

	// 加载配置
	config, err = LoadConfig(CONFIG_FILE)
	if err != nil {
		log.Print(err)
	}
	// log.Println(config)

	// 读取代理的URL信息
	proxyInfo, err = LoadProxyInfo(PROXY_FILE)
	if err != nil {
		log.Print(err)
	}
	// log.Println(proxyInfo)

	// 监听
	port := ":" + config.ProxyPort
	server = http.Server{
		Addr: port,
	}

	if isRunning {
		log.Println("The Server is running")
		return
	}

	// 响应函数
	for _, v := range proxyInfo.ProxyUrls {
		http.HandleFunc(v.Url, DoHandle)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome Home")
	})
	http.HandleFunc("/bettersun", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome bettersun")
	})

	// TODO 需要确认端口是否被占用
	go server.ListenAndServe()
	isRunning = true
	log.Printf("Listen and serve [%v]", port)
}

// 重新加载运行时代理信息
func Reload(p ProxyInfo) {
	proxyInfo = p
	log.Println("The proxy has been updated.")
	// log.Println(proxyInfo)

	// 加载后保存代理信息
	SaveProxy(p)
}

// 关闭服务
func CloseServer() {
	go server.Close()
	log.Println("Server closed")
}

// 请求处理函数
func DoHandle(w http.ResponseWriter, r *http.Request) {
	// 输出请求到文件
	OutRequest(r)

	for _, v := range proxyInfo.ProxyUrls {
		// log.Printf("URL : %v\n", v.Url)
		// log.Printf("Request URL : %v\n", r.URL.String())
		if v.Url == r.URL.String() {

			if v.UseProxy {
				doProxy(w, r, &proxyInfo)
			} else {
				doHandle(w, r, v)
			}
		}
	}
}

// 转发请求函数
// 参考：https://www.cnblogs.com/boxker/p/11046342.html
func doProxy(w http.ResponseWriter, r *http.Request, p *ProxyInfo) {
	fmt.Println("url: ", r.URL)

	// 创建一个HttpClient用于转发请求
	cli := &http.Client{}

	// 读取请求的Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Print("ioutil.ReadAll(r.Body) ", err.Error())
	}

	// 转发的URL
	reqUrl := p.TargetHost + r.URL.String()
	// log.Printf("转发URL: %v", reqUrl)

	// 创建转发用的请求
	req, err := http.NewRequest(r.Method, reqUrl, strings.NewReader(string(body)))
	if err != nil {
		log.Print("http.NewRequest ", err.Error())
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

	// 输出响应体到文件
	OutResponse(r, res)

	// 复制转发的响应Body到响应Body
	io.Copy(w, res.Body)
}

func doHandle(w http.ResponseWriter, r *http.Request, proxyUrl ProxyUrl) {

	// 响应用JSON文件
	jsonFile := moist.CurrentDir() + proxyUrl.ResponseJson
	// log.Printf("响应JSON：%v", jsonFile)

	// 响应文件不存在
	if !moist.IsExist(jsonFile) {
		msg := fmt.Sprintf("JSON文件不存在 %v", jsonFile)
		log.Println(msg)

		m := make(map[string]interface{}, 0)
		m["msg"] = msg
		w.WriteHeader(200)

		msgStream, err := json.Marshal(m)
		if err != nil {
			log.Println(err)
		}
		w.Write(msgStream)
		return
	}

	// JSON文件转换成Map
	obj, err := moist.JsonFileToMap(jsonFile)
	if err != nil {
		log.Println(err)
	}

	// 转换成字节
	objStream, err := json.Marshal(obj)

	w.WriteHeader(200)
	w.Write(objStream)
}

// 输出请求到文件
func OutRequest(r *http.Request) {

	url := r.URL.String()
	// 输出文件名
	filename := strings.ReplaceAll(url, "/", "_")
	filename = "request/" + filename + "_request.json"

	// 读取请求体
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("ioutil.ReadAll() ", err.Error())
	}
	body := string(data)

	m := make(map[string]interface{})
	m["header"] = r.Header
	m["body"] = body

	// 输出请求到文件
	err = moist.OutJson(filename, m)
	if err != nil {
		log.Print(err)
	}
}

// 输出响应到文件
func OutResponse(r *http.Request, res *http.Response) {

	// 输出文件名
	url := r.URL.String()
	filename := strings.ReplaceAll(url, "/", "_")
	filename = "response/" + filename + "_body.json"

	// 读取响应体
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("ioutil.ReadAll() ", err.Error())
	}
	resBody := string(data)

	// 输出响应体到文件
	err = moist.OutJson(filename, resBody)
	if err != nil {
		log.Print(err)
	}
	// // 输出响应头到文件
	// moist.OutJson("real_response_header.json", res.Header)
}
