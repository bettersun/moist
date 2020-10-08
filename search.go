package moist

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)


// 搜索模式(忽略模式)
//   0 ： 默认
//   1 ： 包含
//   2 ： 相等
//   3 ： 以开头
//   4 ： 以结尾
//   5 ： 正则表达式
const PatternDefault string = "0"
const PatternInclude string = "1"
const PatternEqual string = "2"
const PatternPrefix string = "3"
const PatternSuffix string = "4"
const PatternRegexp string = "5"

// 搜索选项
type SearchOption struct {
	RootPath           []string `yaml:"rootPath"`           // 目标根目录
	SearchSubPath      bool     `yaml:"searchSubPath"`      // 搜索子目录
	MatchCase          bool     `yaml:"matchCase"`          // 区分大小写
	FileNamePart       []string `yaml:"fileNamePart"`       // 目标文件名关键字
	FileType           []string `yaml:"fileType"`           // 目标文件类型
	Pattern            string   `yaml:"pattern"`            // 搜索模式
	IgnorePath         []string `yaml:"ignorePath"`         // 忽略目录（目录名完全相同，设置此项有助于提高搜索效率，常用忽略目录: .git, .svn）
	IgnoreFileNamePart []string `yaml:"ignoreFileNamePart"` // 忽略文件名关键字
	IgnoreType         []string `yaml:"ignoreType"`         // 忽略文件类型
	IgnorePattern      string   `yaml:"ignorePattern"`      // 忽略模式
	ShowDetail         bool     `yaml:"showDetail"`         // 是否显示文件详细信息标志
}

// 搜索结果
type SearchResult struct {
	File     []string   `json:"file"`     // 文件全路径
	FileInfo []FileInfo `json:"fileInfo"` // 文件信息
	Count    int        `json:"count"`    // 文件总数
	ErrInfo  []string   `json:"errInfo"`  // 错误信息
	ExInfo   []string   `json:"exInfo"`   // 其它信息
}

// 文件信息(自定义FileInfo, 区别于os.FileInfo)
type FileInfo struct {
	File     string      `json:"file"`     // 文件全路径
	FileName string      `json:"fileName"` // 文件名
	FileType string      `json:"fileType"` // 文件类型
	Size     int64       `json:"size"`     // 文件大小
	Mode     os.FileMode `json:"mode"`     // 文件模式
	ModTime  string      `json:"modTime"`  // 最终修改时间
	//	SysInfo  interface{} `json:"sysInfo"`  // 底层数据来源
}

// 搜索多个目录下的文件,返回文件全路径切片
// 添加参数判断
//   目标文件类型已指定, 忽略文件类型无须指定。后续处理中忽略文件类型作为空处理
// result : 文件全路径切片
// option : 搜索选项
func SearchFile(option SearchOption) []string {

	result := Search(option)

	return result.File
}

// 搜索多个目录下的文件
// 添加参数判断
//   目标文件类型已指定, 忽略文件类型无须指定。后续处理中忽略文件类型作为空处理
// result : 搜索结果
// option : 搜索选项
func Search(option SearchOption) (result SearchResult) {

	//	log.Printf("Option: %+v\n", option)

	// 多个目录搜索
	for _, p := range option.RootPath {
		if strings.TrimSpace(p) == "" {
			info := "空目录已忽略。"
			log.Println(info)
			result.ExInfo = append(result.ExInfo, info)
		}
	}

	if !IsAllEmpty(option.FileType) && !IsAllEmpty(option.IgnoreType) {
		info := "目标文件类型已指定, 忽略文件类型指定无效。后续处理中忽略文件类型作为空处理。"
		log.Println(info)
		option.IgnoreType = []string{}
		result.ExInfo = append(result.ExInfo, info)
	}

	// 不区分大小写的场合，比较用属性(目标文件名关键字、忽略文件名关键字)需转换为小写。
	if !option.MatchCase {
		for i, _ := range option.FileNamePart {
			option.FileNamePart[i] = strings.ToLower(option.FileNamePart[i])
		}
		for i, _ := range option.IgnoreFileNamePart {
			option.IgnoreFileNamePart[i] = strings.ToLower(option.IgnoreFileNamePart[i])
		}
	}

	// 目标文件名关键字和忽略文件名关键字有相同元素的场合，忽略文件名关键字的相同元素无效。
	for _, nmPart := range option.FileNamePart {
		if nmPart != "" {
			for j, igNmPart := range option.IgnoreFileNamePart {
				if nmPart == igNmPart {
					info := "目标文件名关键字和忽略文件名关键字重复，重复的忽略文件名关键字已重置为空。> " + igNmPart
					log.Println(info)
					option.IgnoreFileNamePart[j] = ""
					result.ExInfo = append(result.ExInfo, info)
				}
			}
		}
	}

	// 移除空元素
	option.RootPath = RemoveEmpty(option.RootPath)
	option.FileNamePart = RemoveEmpty(option.FileNamePart)
	option.IgnoreFileNamePart = RemoveEmpty(option.IgnoreFileNamePart)

	// 多个目录搜索
	for _, p := range option.RootPath {
		searchRecursive(&result, p, &option)
	}

	if len(result.ErrInfo) != 0 {
		log.Println("搜索过程中有错误发生。错误信息如下：")
		log.Println(result.ErrInfo)
	}

	return result
}

// 递归搜索目录及子目录下的文件
// result : 搜索结果
// path : 单个目录
// option : 搜索选项
func searchRecursive(result *SearchResult, path string, option *SearchOption) {
	sub, err := ioutil.ReadDir(path)
	if err != nil {
		info := "目录不存在，或打开错误。"
		result.ErrInfo = append(result.ErrInfo, path+" : "+info)
		return
	}

	separator := string(filepath.Separator)

	for _, f := range sub {
		var tmp string
		if strings.HasSuffix(path, separator) {
			tmp = string(path + f.Name())
		} else {
			tmp = string(path + separator + f.Name())
		}

		nameCompare := f.Name()
		// 不区分大小写的场合，都按照小写比较
		if !option.MatchCase {
			nameCompare = strings.ToLower(f.Name())
		}

		// 不忽略子目录、且是目录的场合,递归搜索。
		if option.SearchSubPath && f.IsDir() {
			// 过滤被忽略的目录（目录名完全相同）
			if !IsInSlice(option.IgnorePath, f.Name()) {
				searchRecursive(result, tmp, option)
			}
		} else if !f.IsDir() {
			// 搜索模式匹配
			var isMatchPattern bool
			if !IsAllEmpty(option.FileNamePart) {
				if nameCompare != S_EMPTY {
					// 相等
					if option.Pattern == PatternEqual {
						for _, nmPart := range option.FileNamePart {
							if nameCompare == nmPart {
								isMatchPattern = true
								break
							}
						}
					}
					// 包含
					if option.Pattern == PatternInclude {
						for _, nmPart := range option.FileNamePart {
							if strings.Index(nameCompare, nmPart) >= 0 {
								isMatchPattern = true
								break
							}
						}
					}
					// 以开头
					if option.Pattern == PatternPrefix {
						for _, nmPart := range option.FileNamePart {
							if strings.HasPrefix(nameCompare, nmPart) {
								isMatchPattern = true
								break
							}
						}
					}
					// 以结尾
					if option.Pattern == PatternSuffix {
						for _, nmPart := range option.FileNamePart {
							if strings.HasSuffix(nameCompare, nmPart) {
								isMatchPattern = true
								break
							}
						}
					}
					// 正则表达式
					if option.Pattern == PatternRegexp {
						for _, nmPart := range option.FileNamePart {
							// 正则表达式
							regExp := regexp.MustCompile(nmPart)
							if err != nil {
								errRegExp := "正则表达式转换错误。" + err.Error()
								log.Println(errRegExp)
								result.ErrInfo = append(result.ErrInfo, errRegExp)
								continue
							}
							// 匹配
							if regExp.MatchString(nameCompare) {
								isMatchPattern = true
								break
							}
						}
					}
				}
			}

			// 忽略模式匹配
			var isIgnorePattern bool
			if nameCompare != S_EMPTY && !IsAllEmpty(option.IgnoreFileNamePart) {
				// 相等
				if option.IgnorePattern == PatternEqual {
					for _, nmPart := range option.IgnoreFileNamePart {
						if nameCompare == nmPart {
							isIgnorePattern = true
							break
						}
					}
				}
				// 包含
				if option.IgnorePattern == PatternInclude {
					for _, nmPart := range option.IgnoreFileNamePart {
						if strings.Index(nameCompare, nmPart) >= 0 {
							isIgnorePattern = true
							break
						}
					}
				}
				// 以开头
				if option.IgnorePattern == PatternPrefix {
					for _, nmPart := range option.IgnoreFileNamePart {
						if strings.HasPrefix(nameCompare, nmPart) {
							isIgnorePattern = true
							break
						}
					}
				}
				// 以结尾
				if option.IgnorePattern == PatternSuffix {
					for _, nmPart := range option.IgnoreFileNamePart {
						if strings.HasSuffix(nameCompare, nmPart) {
							isIgnorePattern = true
							break
						}
					}
				}
				// 正则表达式
				if option.IgnorePattern == PatternRegexp {
					for _, nmPart := range option.IgnoreFileNamePart {
						// 正则表达式
						regExp := regexp.MustCompile(nmPart)
						if err != nil {
							errRegExp := "正则表达式转换错误。" + err.Error()
							log.Println(errRegExp)
							result.ErrInfo = append(result.ErrInfo, errRegExp)
							continue
						}
						// 匹配
						if regExp.MatchString(nameCompare) {
							isIgnorePattern = true
							break
						}
					}
				}
			}

			var isTarget bool
			// 非忽略对象、且、部分文件名未指定或者已指定并已符合要求后， 再进一步判断。
			if !isIgnorePattern && (IsAllEmpty(option.FileNamePart) || isMatchPattern) {
				// 目标文件类型被指定
				if !IsAllEmpty(option.FileType) {

					// 属于目标文件类型
					if IsInSuffix(option.FileType, nameCompare) {
						isTarget = true
					}
				} else { // 目标文件类型为空

					// 忽略文件类型被指定
					if !IsAllEmpty(option.IgnoreType) {

						// 不属于忽略文件类型
						if !IsInSuffix(option.IgnoreType, nameCompare) {
							isTarget = true
						}
					} else { // 忽略文件类型为空
						isTarget = true
					}
				}
			}

			// 搜索结果
			if isTarget {

				tmp = strings.Replace(tmp, separator, S_SLASH, -1)

				// 文件全路径
				result.File = append(result.File, tmp)
				// 文件详细信息
				if option.ShowDetail {
					var fInfo FileInfo
					// 文件全路径
					fInfo.File = tmp
					// 文件名
					fInfo.FileName = f.Name()
					// 文件类型
					if strings.Index(fInfo.FileName, S_DOT) >= 0 {
						types := strings.Split(fInfo.FileName, S_DOT)
						fInfo.FileType = types[len(types)-1]
					}
					// 文件大小
					fInfo.Size = f.Size()
					// 文件模式
					fInfo.Mode = f.Mode()
					// 最终修改时间
					fInfo.ModTime = f.ModTime().Format(DateTimeFormatYmdHmsSlash)
					//					fInfo.SysInfo = f.Sys()
					result.FileInfo = append(result.FileInfo, fInfo)
				}

				// 文件总数
				result.Count = result.Count + 1
			}
		}
	}
}


// 快捷搜索
// 搜索多个目录下的指定文件类型, 返回文件全路径切片
//
// path : 搜索目录
// fileType : 搜索文件类型
// ignorePath : 忽略目录(设置此项有助于提高查找效率)，例：[]string{".git", ".svn"}
func SearchByType(path []string, fileType []string, ignorePath []string) []string {

	var option SearchOption
	// 目标根目录
	option.RootPath = path
	// 查找子文件夹
	option.SearchSubPath = true
	// 目标文件类型
	option.FileType = fileType
	// 忽略目录
	option.IgnorePath = ignorePath

	result := SearchFile(option)

	return result
}

// 快捷搜索
// 以目标文件名关键字搜索多个目录下的指定文件类型, 返回文件全路径切片
//
//
// path : 搜索目录
// fileNamePart : 目标文件名关键字
// fileType : 搜索文件类型
// ignorePath : 忽略目录(设置此项有助于提高查找效率)，例：[]string{".git", ".svn"}
func SearchByNamePartType(path []string, fileNamePart []string, fileType []string, ignorePath []string) []string {

	var option SearchOption
	// 目标根目录
	option.RootPath = path
	// 查找子文件夹
	option.SearchSubPath = true
	// 区分大小写
	option.MatchCase = false
	// 查找模式： 包含
	option.Pattern = PatternInclude
	// 目标文件名关键字
	option.FileNamePart = fileNamePart
	// 目标文件类型
	option.FileType = fileType
	// 忽略目录
	option.IgnorePath = ignorePath
	// 忽略模式： 默认
	option.IgnorePattern = PatternDefault

	result := SearchFile(option)

	return result
}

// 快捷搜索XLSX(XLSM)文件
//
// path : 搜索目录
// ignorePath : 忽略目录(设置此项有助于提高查找效率)，例：[]string{".git", ".svn"}
func SearchXlsx(path []string, ignorePath []string) []string {

	var option SearchOption
	// 目标根目录
	option.RootPath = path
	// 查找子文件夹
	option.SearchSubPath = true
	// 区分大小写
	option.MatchCase = true
	// 查找模式： 包含
	option.Pattern = PatternInclude
	// 目标文件类型
	option.FileType = []string{"xlsx", "xlsm"}
	// 忽略目录
	option.IgnorePath = ignorePath
	// 忽略模式： 以开头
	option.IgnorePattern = PatternPrefix
	// 忽略文件名部分
	option.IgnoreFileNamePart = []string{"~$"}

	result := SearchFile(option)

	return result
}
