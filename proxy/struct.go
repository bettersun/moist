package proxy

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// 配置
type ProxyConfig struct {
	ProxyPort string `yaml:"proxyPort"`
}

// 代理信息
type ProxyInfo struct {
	TargetHost string     `yaml:"targetHost"`
	ProxyUrls  []ProxyUrl `yaml:"proxyUrls"`
}

// 代理URL
type ProxyUrl struct {
	Url          string `yaml:"url"`
	UseProxy     bool   `yaml:"useProxy"`
	ResponseJson string `yaml:"responseJson"`
}

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
func LoadProxyInfo(file string) (ProxyInfo, error) {

	var info ProxyInfo

	// 读取文件
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Print(err)
		return info, err
	}

	// 转换成Struct
	err = yaml.Unmarshal(b, &info)
	if err != nil {
		log.Printf("Get the setting error! %v\n", err.Error())
	}

	return info, nil
}
