package main

import (
	"github.com/bettersun/moist/proxy"
)

// 加载配置
func RunServer() {

	// 加载配置
	proxy.RunServer()
}

// 重新加载
func Reload() {

	var proxyInfo proxy.ProxyInfo

	var proxyUrls []proxy.ProxyUrl
	var proxyUrl proxy.ProxyUrl

	proxyUrl.Url = "/goodbye"
	proxyUrl.UseProxy = false
	proxyUrl.ResponseJson = "/json/goodbye.json"

	proxyUrls = append(proxyUrls, proxyUrl)

	proxyUrl.Url = "/hello"
	proxyUrl.UseProxy = false
	proxyUrl.ResponseJson = "/json/hello.json"

	proxyUrls = append(proxyUrls, proxyUrl)

	proxyInfo.TargetHost = "http://localhost:8002"
	proxyInfo.ProxyUrls = proxyUrls

	proxy.Reload(proxyInfo)
}

/// 关闭服务
func CloseServer() {

	// 关闭服务
	proxy.CloseServer()
}
