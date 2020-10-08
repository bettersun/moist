package moist


import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/// 首字母转小写
func LowerFirst(s string) string {
	if s != `` {
		s = strings.Replace(s, string(s[0]), strings.ToLower(string(s[0])), 1)
	}

	return s
}

/// 首字母转大写
func UperFirst(s string) string {
	if s != `` {
		s = strings.Replace(s, string(s[0]), strings.ToUpper(string(s[0])), 1)
	}

	return s
}

// 短线连接字符串转驼峰(首字母转大写)
//  无换行符
// 例：
//   test_id  => TestId
//   test_place_name  => TestPlaceName
//
// isUpFirst : true: 首字母转大写 false:首字母转小写
// convOther : true: 非短线连接转换成小写 false: 非短线连接不转换
func HyphenToCamelLine(s string, isUpFirst bool, convOther bool) string {

	if convOther {
		// 非短线连接转换成小写
		s = strings.ToLower(s)
	}

	regSnake := regexp.MustCompile(`-.`)
	// 短线连接字符串转驼峰
	s = regSnake.ReplaceAllStringFunc(s, hyphenToCamelSingle)
	if isUpFirst {

		// 首字母转大写
		s = strings.Replace(s, string(s[0]), strings.ToUpper(string(s[0])), 1)
	} else {

		// 首字母转小写
		s = strings.Replace(s, string(s[0]), strings.ToLower(string(s[0])), 1)
	}

	return s
}

// 下划线连接字符串转驼峰，正则替换用的函数
func hyphenToCamelSingle(s string) string {
	return strings.ToUpper(strings.Replace(s, "-", "", -1))
}

// 下划线连接字符串转驼峰
//  含有换行符
//
// isUpFirst : true: 首字母转大写 false:首字母转小写
// convOther : true: 非下划线连接转换成小写 false: 非下划线连接不转换
func SnakeToCamel(s string, isUpFirst bool, convOther bool) string {
	arr := strings.Split(s, "\n")

	for i, v := range arr {
		if len(v) > 0 {
			arr[i] = SnakeToCamelLine(v, isUpFirst, convOther)
		}
	}

	return strings.Join(arr, "\n")
}

// 下划线连接字符串转驼峰(首字母转大写)
//  无换行符
// 例：
//   test_id  => TestId
//   test_place_name  => TestPlaceName
//
// isUpFirst : true: 首字母转大写 false:首字母转小写
// convOther : true: 非下划线连接转换成小写 false: 非下划线连接不转换
func SnakeToCamelLine(s string, isUpFirst bool, convOther bool) string {

	if convOther {
		// 非下划线连接转换成小写
		s = strings.ToLower(s)
	}

	regSnake := regexp.MustCompile(`_.`)
	// 下划线连接字符串转驼峰
	s = regSnake.ReplaceAllStringFunc(s, snakeToCamelSingle)
	if isUpFirst {

		// 首字母转大写
		s = strings.Replace(s, string(s[0]), strings.ToUpper(string(s[0])), 1)
	} else {

		// 首字母转小写
		s = strings.Replace(s, string(s[0]), strings.ToLower(string(s[0])), 1)
	}

	return s
}

// 下划线连接字符串转驼峰，正则替换用的函数
func snakeToCamelSingle(s string) string {
	return strings.ToUpper(strings.Replace(s, "_", "", -1))
}

// 驼峰转下划线
//  含有换行符
func CamelToSnake(s string) string {

	arr := strings.Split(s, "\n")

	for i, v := range arr {
		arr[i] = CamelToSnakeLine(v)
	}

	return strings.Join(arr, "\n")
}

// 驼峰转下划线
//  无换行符
func CamelToSnakeLine(s string) string {

	arr := strings.Split(s, "")

	for i, c := range arr {

		low := strings.ToLower(c)

		// 转换成小写后，不等于原字符，则为大写
		// 首字母前不加下划线
		if low != c && i == 0 {
			arr[i] = low
		} else if low != c && i > 0 {
			arr[i] = "_" + low
		}
	}

	return strings.Join(arr, "")
}

func ToDecimal(value float64, c int) float64 {

	sfmt := "%." + strconv.Itoa(c) + "f"

	value, err := strconv.ParseFloat(fmt.Sprintf(sfmt, value), 64)
	if err != nil {
		fmt.Println(err)
	}

	return value
}
