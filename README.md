# moist
a go library

## YAML转struct / JSON转struct

struct定义
```go
type Config struct {
	Name string `yaml:"name"`
}
```

config.yml

```yaml
name: bettersun
```

调用YAML转struct
```go
	file := "config.yml"

	var config Config
	result, err := YamlToStruct(file, &config)
	if err != nil {
		log.Println(err)
	}

	log.Println(result)
	log.Println(result.(*Config).Name)
```

config.json

```json
{"name":"bettersun"}
```

调用JSON转struct
```go
	file := "config.json"

	var config Config
	result, err := JsonToStruct(file, &config)
	if err != nil {
		log.Println(err)
	}

	log.Println(result)
	log.Println(result.(*Config).Name)
```