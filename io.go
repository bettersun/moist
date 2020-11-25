package moist

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 获取程序运行路径
// 参考：
//  https://blog.csdn.net/sufu1065/article/details/80116627
func CurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// 引用： https://www.cnblogs.com/wangqishu/p/5147107.html
// 判断文件或目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// 在文件全路径中截取文件名
func FileName(filePath string, separator string) string {

	if separator == "" {
		separator = string(os.PathSeparator)
	}

	pathIndex := strings.LastIndex(filePath, separator)

	var fName string
	if pathIndex == -1 {
		log.Printf("There is no separator in the file path. [%v]\n", separator)
	} else {
		path := strings.Split(filePath, separator)
		fName = path[len(path)-1]
	}

	return fName
}

// 在文件全路径中截取路径
func FilePath(filePath string, separator string) string {

	if separator == "" {
		separator = string(os.PathSeparator)
	}

	pathIndex := strings.LastIndex(filePath, separator)

	var fPath []string
	if pathIndex == -1 {
		log.Printf("There is no separator in the file path. [%v]\n", separator)
	} else {
		path := strings.Split(filePath, separator)
		fPath = path[0 : len(path)-1]
	}

	return strings.Join(fPath, separator)
}

// 复制多个文件到指定目录
func Copy(list []string, destPath string) {

	separator := string(os.PathSeparator)

	log.Printf("separator: %v\n", separator)
	log.Printf("destPath Before: %v\n", destPath)
	if !strings.HasSuffix(destPath, separator) {
		destPath = destPath + separator
	}
	log.Printf("destPath After: %v\n", destPath)

	for i, srcFile := range list {

		if !IsExist(srcFile) {

			log.Printf("Error: File does not exits. %v\n", srcFile)
			list[i] = ""
		}
	}

	for _, srcFile := range list {

		if srcFile != "" {

			log.Printf("srcFile: %v\n", srcFile)
			fNm := FileName(srcFile, "")
			log.Printf("FileName: %v\n", fNm)
			destFile := destPath + fNm
			log.Printf("destFile: %v\n", destFile)

			copy(srcFile, destFile)
		}
	}
}

// 引用： https://www.jianshu.com/p/6cc1938260ba
func copy(src string, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// 参考： https://www.golangtc.com/t/53749f5d320b521bba000007
// 多个文件合并成一个文件
func Merge(srcFile []string, outFileName string) {

	outFile, openErr := os.OpenFile(outFileName, os.O_CREATE|os.O_WRONLY, 0600)
	if openErr != nil {
		log.Printf("Can not create file %s\n", outFileName)
	}

	bWriter := bufio.NewWriter(outFile)

	for _, f := range srcFile {
		fp, fpOpenErr := os.Open(f)
		if fpOpenErr != nil {
			log.Printf("Can not open file %v\n", fpOpenErr)
			continue
		}

		bReader := bufio.NewReader(fp)
		for {
			buffer := make([]byte, 1024)
			readCount, readErr := bReader.Read(buffer)
			if readErr == io.EOF {
				break
			} else {
				bWriter.Write(buffer[:readCount])
			}
		}
	}

	bWriter.Flush()
}

// 将字符串数组写入指定文件
//  各字符串间添加换行符
func WriteFile(filename string, strarr []string) error {

	tmp := StringSliceToByte(strarr)
	err := ioutil.WriteFile(filename, tmp, 0644)

	if err != nil {
		log.Printf("WriteToFile Error ： %v\n", err)
		return err
	}

	return nil
}

// 输出到控制台
func OutTerminal(obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		log.Println(err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, S_EMPTY, S_INDENT)

	out.WriteTo(os.Stdout)
}

// 输出到Json文件
func OutJson(file string, obj interface{}) error {
	// HTML转义的处理 START
	// b, err := json.Marshal(obj)
	// if err != nil {
	// 	log.Println(err)
	// }

	// var out bytes.Buffer
	// err = json.Indent(&out, b, S_EMPTY, S_INDENT)
	// HTML转义的处理 END

	// HTML原文输出
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent(S_EMPTY, S_INDENT)
	jsonEncoder.Encode(obj)

	err := ioutil.WriteFile(file, bf.Bytes(), os.ModePerm) // 覆盖所有Unix权限位（用于通过&获取类型位）
	if err != nil {
		log.Println(err)
	}

	return err
}

// 创建目录树和文件
func CreateTree(rootPath string, file []string, isCreateFile bool) error {

	if len(file) == 0 {
		file = []string{""}
	}

	for _, s := range file {

		// TODO ../
		if strings.HasPrefix(s, "../") {
			s = strings.Replace(s, "../", "/", 1)
		}
		if strings.HasPrefix(s, "./") {
			s = strings.Replace(s, "./", "/", 1)
		}

		if strings.HasPrefix(s, "///") {
			s = strings.Replace(s, "///", "", 1)
		}
		if strings.HasPrefix(s, "//") {
			s = strings.Replace(s, "//", "", 1)
		}
		if strings.HasPrefix(s, "/") {
			s = strings.Replace(s, "/", "", 1)
		}

		if rootPath != "" && !strings.HasSuffix(rootPath, "/") {
			rootPath = rootPath + "/"
		}

		p := rootPath + FilePath(s, S_SLASH)

		//递归创建文件夹
		err := os.MkdirAll(p, os.ModePerm)
		if err != nil {
			log.Println(err)
			continue
		}

		if isCreateFile {

			s = rootPath + s
			// 创建文件
			_, err = os.Create(s)
			if err != nil {
				log.Println(err)
			}

		}
	}

	return nil
}

// // 创建目录树
// func CreatePathTree(rootPath string, subPath []string) error {

// 	for _, s := range subPath {

// 		log.Println("0:" + s)
// 		if strings.HasPrefix(s, "../") {
// 			s = strings.Replace(s, "../", "/", 1)
// 		}
// 		log.Println("1:" + s)
// 		if strings.HasPrefix(s, "./") {
// 			s = strings.Replace(s, "./", "/", 1)
// 		}

// 		log.Println(s)
// 		p := rootPath + FilePath(s, S_SLASH)
// 		log.Println(p)

// 		//递归创建文件夹
// 		err := os.MkdirAll(p, os.ModePerm)
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 	}

// 	return nil
// }

/// 读取文件内容
/// 本方法适用于小文件
func ReadFile(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Read file failed %v\n", err)
		return "", err
	}
	return string(bytes), nil
}

// 单文件复制
func CopyFile(src string, dst string) error {
	p := FilePath(dst, S_SLASH)

	//递归创建文件夹
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		log.Println(err)
		return err
	}

	// 创建文件
	_, err = os.Create(dst)
	if err != nil {
		log.Println(err)
		return err
	}

	copy(src, dst)

	return nil
}
