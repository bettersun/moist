package moist

import (
	"fmt"
	"testing"
)

// func TestSearch(t *testing.T) {

// 	//	path := []string{`C:\Users\bs\go\src\github.com\bettersun\moist`}
// 	path := []string{`/Users/bs/Documents/Develop/go/src/github.com/bettersun/moist`}
// 	fileNamePart := []string{""}
// 	//	fileNamePart := []string{"search_result", "BsTeST", "aBc"}
// 	fileType := []string{".go", "java", ".txt", ".yml", ".json"}
// 	//	fileType := []string{""}

// 	var option SearchOption
// 	// 目标根目录
// 	option.RootPath = path
// 	// 查找子文件夹
// 	option.SearchSubPath = true
// 	// 区分大小写
// 	option.MatchCase = true
// 	// 查找模式： 包含
// 	option.Pattern = PATTERN_INCLUDE
// 	// 目标文件名关键字
// 	option.FileNamePart = fileNamePart
// 	// 目标文件类型
// 	option.FileType = fileType
// 	// 忽略目录(设置此项有助于提高查找效率)
// 	option.IgnorePath = []string{".git", ".svn"}
// 	// 忽略类型
// 	option.IgnoreType = []string{}
// 	// 忽略模式： 默认
// 	option.IgnorePattern = PATTERN_INCLUDE
// 	// 忽略文件名关键字
// 	option.IgnoreFileNamePart = []string{"search_result", "search_test", "BsTest", "abc", "eFeGoP"}
// 	// 获取文件详细信息
// 	option.ShowDetail = true

// 	result := SearchFile(option)
// 	//	fmt.Println(result)

// 	OutFile("./search_result.txt", result)
// }

// func TestSearch_02(t *testing.T) {

// 	path := []string{"E:/BS", "E:/Test"}
// 	fileType := []string{".xlsx", ".xlsm"}
// 	result := SearchByType(path, fileType, nil)
// 	fmt.Println(result)
// }

// func TestSearch_03(t *testing.T) {

// 	path := []string{"E:/BS", "E:/Test"}
// 	result := SearchXlsx(path, nil)
// 	fmt.Println(result)
// }

func TestSearch_04(t *testing.T) {

	var option SearchOption
	option.RootPath = []string{`E:/BS/Test`}
	option.SearchSubPath = true
	option.FileNamePart = []string{`1`}
	// option.MatchCase = matchCase
	option.Pattern = "1"
	// option.FileType = fType
	// option.IgnorePath = igPath
	// option.IgnoreFileNamePart = []string{``}
	// option.IgnorePattern = ignorePattern
	// option.IgnoreType = igFType

	result := Search(option)
	fmt.Println(result)
	OutFile("./search_result.txt", result)
}
