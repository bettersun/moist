package yaml

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// 输出到Yaml文件
func OutYaml(file string, s interface{}) error {
	b, err := yaml.Marshal(&s)
	if err != nil {
		log.Println(err)
		return err
	}
	// log.Println(b)

	err = ioutil.WriteFile(file, b, os.ModePerm) // 覆盖所有Unix权限位（用于通过&获取类型位）
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

/// YML文件转struct
///  file: yaml文件
///  s   : 定义的struct的地址(调用处需要加&)
///
/// 使用例：参照README.md
func YamlFileToStruct(file string, s interface{}) (interface{}, error) {

	// 读取文件
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Print(err)
		return s, err
	}

	// 转换成Struct
	err = yaml.Unmarshal(b, s)
	if err != nil {
		log.Printf("Get the setting error! %v\n", err.Error())
	}

	return s, nil
}
