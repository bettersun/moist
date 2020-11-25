package try

import (
	"log"
	"testing"

	"github.com/bettersun/moist"
	"github.com/bettersun/moist/yaml"
)

func TestTry_01(t *testing.T) {
	file := "config.yml"

	var config Config
	result, err := yaml.YamlFileToStruct(file, &config)
	if err != nil {
		log.Println(err)
	}

	log.Println(result)
	log.Println(result.(*Config).Name)
}

func TestTry_02(t *testing.T) {
	file := "config.json"

	var config Config
	result, err := moist.JsonFileToStruct(file, &config)
	if err != nil {
		log.Println(err)
	}

	log.Println(result)
	log.Println(result.(*Config).Name)
}

func TestTry_03(t *testing.T) {

	file := "outyaml.yml"

	var config Config
	config.Name = "OutYamlTest<html>-:23"

	yaml.OutYaml(file, config)
}
