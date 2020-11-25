package main

import (
	"github.com/bettersun/moist/proxy"
)

func runServer() {

	// 加载配置
	proxy.RunServer()
}

func reload() {

	var proxyInfo proxy.ProxyInfo

	var proxyUrls []proxy.ProxyUrl
	var proxyUrlInfo proxy.ProxyUrl

	proxyUrlInfo.Url = "/goodbye"
	proxyUrlInfo.UseProxy = false
	proxyUrlInfo.ResponseJson = "/json/goodbye.json"

	proxyUrls = append(proxyUrls, proxyUrlInfo)

	proxyUrlInfo.Url = "/hello"
	proxyUrlInfo.UseProxy = false
	proxyUrlInfo.ResponseJson = "/json/hello.json"

	proxyUrls = append(proxyUrls, proxyUrlInfo)

	proxyInfo.TargetHost = "http://localhost:8002"
	proxyInfo.BaseUrl = "/bettersun"
	proxyInfo.ProxyUrls = proxyUrls

	proxy.Reload(proxyInfo)
}

func closeServer() {

	// 加载配置
	proxy.CloseServer()
}