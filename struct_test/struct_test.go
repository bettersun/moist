package structtest

import (
	"log"
	"testing"

	"github.com/bettersun/moist"
)

type User struct {
	Name string `json:"name"`
	Age  int8   `json:"age"`
	Info Info   `json:"info"`
}

type Info struct {
	Tel  string `json:"tel"`
	Mail string `json:"mail"`
}

func TestStruct_01(t *testing.T) {
	file := "user.json"

	var usr User
	err := moist.JsonFileToStruct(file, &usr)
	if err != nil {
		log.Println(err)
	}

	log.Println(usr)
	log.Println(usr.Name)
	log.Println(usr.Age)
	log.Println(usr.Info.Tel)
	log.Println(usr.Info.Mail)
}

func TestStruct_02(t *testing.T) {
	m := make(map[interface{}]interface{})
	m["name"] = "张三"
	m["age"] = 26

	info := make(map[interface{}]interface{})
	info["tel"] = "12356789"
	info["mail"] = "bettersun@163.com"
	m["info"] = info

	var usr User

	mTemp, err := moist.ToStringKeyMap(m)
	if err != nil {
		log.Println(err)
	}

	err = moist.MapToStruct(mTemp, &usr)
	if err != nil {
		log.Println(err)
	}

	log.Println(usr)
	log.Println(usr.Name)
	log.Println(usr.Age)
	log.Println(usr.Info.Tel)
	log.Println(usr.Info.Mail)
}
