package moist

import (
	"log"
	"runtime"
	"testing"
)

func TestExplorer(t *testing.T) {

	// var option ExplorerOption
	// option.RootPath = []string{CurrentDir()}
	// option.IncludeSubPath = true
	// option.IgnorePath = []string{`.git`, `.svn`}
	// option.IgnoreFile = []string{`.DS_Store`}
	// fmt.Println(option)

	// result := Explorer(option)
	// log.Println(result)

	var option ExplorerOption
	if runtime.GOOS == "darwin" {
		option.RootPath = []string{`/Users/sunjiashu/Documents/Develop/github.com/bettersun/xtool`}
	}
	if runtime.GOOS == "windows" {
		option.RootPath = []string{`E:\BS\Mac`}
	}

	option.IncludeSubPath = true
	option.IgnorePath = []string{`.git`, `.svn`}
	option.IgnoreFile = []string{`.DS_Store`}

	result := Explorer(option)

	m, err := StructToIfKeyMap(result)
	if err != nil {
		log.Println(err)
	}
	log.Println(m)

	OutFile("./explorer_tree.txt", result)
}
