# Proxy

一个简单的代理工具

## ui

使用fyne开发的简易画面。
config.yml和proxy.yml需要放在ui目录下。

## config.yml

```yaml
proxyPort: 9012
```

## proxy.yml

```yaml
targetHost: http://localhost:8002
baseUrl: /bettersun
proxyUrls:
 -
  useProxy: true
  url: /hello
  responseJson: /json/hello.json
 -
  useProxy: false
  url: /goodbye
  responseJson: /json/goodbye.json
 -
  useProxy: false
  url: /
  responseJson: /json/welcome.json

```