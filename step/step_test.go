package step

import (
	"io/ioutil"
	"log"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/bettersun/moist"
)

func TestStep(t *testing.T) {

	cd, err := LoadCommentConfig("./config.yml")
	if err != nil {
		log.Printf("LoadRegExpConfig error! \n")
	}
	// log.Println(cd)
	// moist.OutFile("./comment.json", cd)

	cre := ToCommentRegExp(cd)
	// log.Println(cre)
	// moist.OutFile("./comment_reg.json", cre)

	file := []string{
		"./test/test.c",
		"./test/test.go",
		"./test/test.xxx",
		"./test/xxss.go",
		"./test/abcd.efg",
		"./test/test.sql"}

	// file = []string{"./test/test.c"}

	step := CountAll(file, &cre)
	// moist.OutFile("./step.json", step)

	ss := Summary(step)
	moist.OutFile("./step_01.json", ss)
}

func TestStepInputFile(t *testing.T) {

	/// 程序输入文件
	const INPUT_FILE = "./test/input.yml"

	// 输入文件
	input, err := ioutil.ReadFile(INPUT_FILE)
	if err != nil {
		log.Printf("文件读取错误 : %+v", err.Error())
	}

	// 取得输入值
	var file []string
	err = yaml.Unmarshal(input, &file)
	if err != nil {
		log.Printf("无法转换文件内容 : %+v", err.Error())
	}

	// 统计
	result := Step(file)
	moist.OutFile("./step_02.json", result)
}
