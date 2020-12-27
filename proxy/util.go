package proxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/bettersun/moist"
)

///
func SaveProxy(p ProxyInfo) {

	bkFile := fmt.Sprintf("%v/backup/proxy%v.yml", moist.CurrentDir(), moist.NowYmdHms())
	// log.Println(bkFile)
	moist.CopyFile(PROXY_FILE, bkFile)

	// 备份成功后覆盖当前yml文件
	if moist.IsExist(bkFile) {

		d, err := yaml.Marshal(&p)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		// log.Println(d)

		err = ioutil.WriteFile(PROXY_FILE, d, os.ModePerm) // 覆盖所有Unix权限位（用于通过&获取类型位）
		if err != nil {
			log.Println(err)
		}
	}

}
