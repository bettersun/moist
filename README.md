# moist
a go library



## YAML转struct / JSON转struct

### struct定义
```go
type Config struct {
	Name string `yaml:"name"`
}
```

### YAML转struct

#### YAML文件
config.yml
```yaml
name: bettersun
```

#### 调用YAML转struct

导入包
```go
import (
	yml "github.com/bettersun/moist/yaml"
)

```

转换代码
```go
	file := "config.yml"

	var config Config
	result, err := yml.YamlFileToStruct(file, &config)
	if err != nil {
		log.Println(err)
	}

	log.Println(result)
	log.Println(result.(*Config).Name)
```

#### YAML文件内容为数组

config.yml
```yaml
- name: bettersun
- name: better
- name: sun
```

转换代码
```go
	file := "config.yml"

	var config Config
	result, err := yml.YamlFileToStruct(file, &config)
	if err != nil {
		log.Println(err)
	}

	log.Println(result.(*[]Config))
	log.Println(*(result.(*[]Config)))
	for _, v := range *(result.(*[]Config)) {
		log.Println(v.Name)
	}
```

### JSON转struct

#### JSON文件
config.json

```json
{"name":"bettersun"}
```

#### 调用JSON转struct

导入包
```go
import (
	yml "github.com/bettersun/moist"
)

```

转换代码
```go
	file := "config.json"

	var config Config
	result, err := moist.JsonFileToStruct(file, &config)
	if err != nil {
		log.Println(err)
	}

	log.Println(result)
	log.Println(result.(*Config).Name)
```
