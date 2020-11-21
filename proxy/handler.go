package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bettersun/moist"
)

// 请求处理函数
func doHandle(w http.ResponseWriter, r *http.Request) {

	log.Println(urlInfo)

	for _, v := range urlNoProxy {
		if v.Url == r.URL.String() {

			// 响应用JSON文件
			jsonFile := moist.CurrentDir() + v.ResponseJson
			log.Println(jsonFile)

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
	}
}
