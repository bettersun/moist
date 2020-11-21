package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// 配置
type ProxyConfig struct {
	ProxyPort  string `yaml:"proxyPort"`
	TargetHost string `yaml:"targetHost"`
}

// 代理的URL信息
type ProxyUrlInfo struct {
	Url          string `yaml:"url"`
	UseProxy     bool   `yaml:"useProxy"`
	ResponseJson string `yaml:"responseJson"`
}

const configFile = "proxy_config.yml"
const urlFile = "proxy_url.yml"

/// 读取配置
func LoadConfig(file string) (ProxyConfig, error) {

	var config ProxyConfig

	// 读取文件
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Print(err)
		return config, err
	}

	// 转换成Struct
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Printf("Get the setting error! %v\n", err.Error())
	}

	return config, nil
}

/// 读取URL信息
func LoadUrlInfo(file string) ([]ProxyUrlInfo, error) {

	var urlInfo []ProxyUrlInfo

	// 读取文件
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Print(err)
		return urlInfo, err
	}

	// 转换成Struct
	err = yaml.Unmarshal(b, &urlInfo)
	if err != nil {
		log.Printf("Get the setting error! %v\n", err.Error())
	}

	return urlInfo, nil
}
