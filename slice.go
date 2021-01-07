package moist

import (
	"bytes"
	"strings"
)

// 判断目标字符串是否是在切片中
func IsInSlice(slice []string, s string) (isIn bool) {

	if len(slice) == 0 {
		return false
	}

	isIn = false
	for _, f := range slice {

		if f == s {
			isIn = true
			break
		}
	}

	return isIn
}

// 判断目标字符串部分是否是在切片中
func IsPartInSlice(slice []string, s string) (isIn bool) {

	if len(slice) == 0 {
		return false
	}

	isIn = false
	for _, part := range slice {
		if strings.Contains(s, part) {
			isIn = true
			break
		}
	}

	return isIn
}

// 判断目标字符串是否包含在切片的某个元素中
func Contains(slice []string, s string) (result bool) {

	if len(slice) == 0 {
		return false
	}

	result = false
	for _, v := range slice {
		if strings.Contains(v, s) {
			result = true
			break
		}
	}

	return result
}

// 判断目标字符串部分是否包含切片中所有项目
func ContainsAll(slice []string, s string) (isContains bool) {

	if len(slice) == 0 {
		return false
	}

	isContains = true
	for _, part := range slice {

		if !strings.Contains(s, part) {
			isContains = false
			break
		}
	}

	return isContains
}

// 目标字符串如果在切片中,返回下标
func InSlice(slice []string, s string) int {

	idx := -1
	if len(slice) == 0 {
		return idx
	}

	for i, f := range slice {

		if f == s {
			idx = i
			break
		}
	}

	return idx
}

// 判断目标字符串的开头是否含有切片中指定的字符串
func IsInPrefix(slice []string, s string) (isIn bool) {

	isIn = false
	for _, f := range slice {

		tmp := strings.ToLower(s)
		f := strings.ToLower(f)
		if strings.TrimSpace(f) != "" && strings.HasPrefix(tmp, f) {
			isIn = true
			break
		}
	}

	return isIn
}

// 判断目标字符串的末尾是否含有切片中指定的字符串
func IsInSuffix(slice []string, s string) (isIn bool) {

	isIn = false
	for _, f := range slice {

		tmp := strings.ToLower(s)
		f := strings.ToLower(f)
		if strings.TrimSpace(f) != "" && strings.HasSuffix(tmp, f) {
			isIn = true
			break
		}
	}

	return isIn
}

// 判断切片各元素是否是空字符串或空格
func IsAllEmpty(slice []string) (isEmpty bool) {

	if len(slice) == 0 {
		return true
	}

	isEmpty = true
	for _, s := range slice {

		if strings.TrimSpace(s) != "" {
			isEmpty = false
			break
		}
	}

	return isEmpty
}

// 字符串切片转换拼接成一个字符切片
//  各字符串间添加换行符
func StringSliceToByte(s []string) []byte {
	var buffer bytes.Buffer

	l := len(s)
	for i, v := range s {
		buffer.WriteString(v)

		if i < l-1 {
			buffer.WriteString("\n")
		}
	}

	return buffer.Bytes()
}

// 删除字符串切片中的空字符串
// 参考：
//   https://blog.csdn.net/fwhezfwhez/article/details/79931415
func RemoveEmpty(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}

	for i, v := range slice {
		if v == "" {
			slice = append(slice[:i], slice[i+1:]...)
			return RemoveEmpty(slice)
			break
		}
	}

	return slice
}

// 去除重复
func Duplicate(slice []string) []string {

	var result []string
	tmpMap := make(map[string]interface{})
	for _, v := range slice {
		if _, ok := tmpMap[v]; !ok {
			result = append(result, v)
			tmpMap[v] = struct{}{}
		}
	}

	return result
}
