package yamltest

import (
	"log"
	"testing"

	"github.com/bettersun/moist/yaml"
)

// 配置
type Config struct {
	Name string `yaml:"name"`
}

func TestYaml_01(t *testing.T) {
	file := "config.yml"

	var config Config
	err := yaml.YamlFileToStruct(file, &config)
	if err != nil {
		log.Println(err)
	}

	log.Println(config)
	log.Println(config.Name)
}

func TestYaml_02(t *testing.T) {
	file := "config2.yml"

	var config []Config
	err := yaml.YamlFileToStruct(file, &config)
	if err != nil {
		log.Println(err)
	}

	log.Println(config)
	for _, v := range config {
		log.Println(v.Name)
	}
}

func TestYaml_03(t *testing.T) {

	file := "outyaml.yml"

	var config Config
	config.Name = "OutYamlTest<html>-:23"

	yaml.OutYaml(file, config)
}
