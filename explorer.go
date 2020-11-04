package moist

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

// 遍历选项
type ExplorerOption struct {
	RootPath       []string `yaml:"rootPath"`       // 目标根目录
	IncludeSubPath bool     `yaml:"includeSubPath"` // 搜索子目录
	IgnorePath     []string `yaml:"ignorePath"`     // 忽略目录
	IgnoreFile     []string `yaml:"ignoreFile"`     // 忽略文件
}

// 树节点
type TreeNode struct {
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	Children []*TreeNode `json:"children"`
	IsDir    bool        `json:"isDir"`
}

// 遍历多个目录下的文件
// option : 遍历选项
// tree : 遍历结果
func Explorer(option ExplorerOption) (tree TreeNode) {

	//	log.Printf("Option: %+v\n", option)

	// 多个目录搜索
	for _, p := range option.RootPath {
		if strings.TrimSpace(p) == "" {
			info := "空目录已忽略。"
			log.Println(info)
		}
	}

	// 移除空元素
	option.RootPath = RemoveEmpty(option.RootPath)

	separator := string(filepath.Separator)

	// 多个目录搜索
	for _, p := range option.RootPath {
		var child TreeNode

		// 目录名
		p = strings.ReplaceAll(p, separator, S_SLASH)
		names := strings.Split(p, S_SLASH)
		var name string
		leng := len(names)
		if leng > 0 {
			if names[leng-1] == "" {
				if leng > 1 {
					name = names[leng-2]
				}
			} else {
				name = names[leng-1]
			}
		}
		child.Name = name

		// 目录路径
		child.Path = p
		// 默认传过来的值是目录
		child.IsDir = true

		// 遍历
		explorerRecursive(&child, &option)
		tree.Children = append(tree.Children, &child)
	}

	return tree
}

// 递归遍历子目录下的文件
// child : 子目录
// option : 遍历选项
func explorerRecursive(child *TreeNode, option *ExplorerOption) {

	path := child.Path
	sub, err := ioutil.ReadDir(path)
	if err != nil {
		info := "目录不存在，或打开错误。"
		log.Println(info)
		log.Println(err)
		return
	}

	// separator := string(filepath.Separator)
	separator := S_SLASH

	for _, f := range sub {

		var grandChild TreeNode

		var tmp string
		if strings.HasSuffix(path, separator) {
			tmp = string(path + f.Name())
		} else {
			tmp = string(path + separator + f.Name())
		}

		// 子目录
		grandChild.Path = tmp

		if f.IsDir() {
			grandChild.IsDir = true
		} else {
			// 叶子节点(文件)
			grandChild.IsDir = false
		}

		// 不忽略子目录、且是目录的场合,递归遍历
		if option.IncludeSubPath && f.IsDir() {

			if IsInSlice(option.IgnorePath, f.Name()) {
				continue
			}

			explorerRecursive(&grandChild, option)
		}

		// 忽略文件不处理
		if !f.IsDir() && IsInSlice(option.IgnoreFile, f.Name()) {
			continue
		}

		grandChild.Name = f.Name()
		child.Children = append(child.Children, &grandChild)
	}
}
