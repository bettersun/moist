package moist

import "strings"

// 判断目标字符串的末尾是否含有指定的字符串,忽略大小写
func HasSuffixIgnoreCase(s string, f string) (has bool) {

	has = false
	tmp := strings.ToLower(s)
	f = strings.ToLower(f)
	if strings.TrimSpace(f) != "" && strings.HasSuffix(tmp, f) {
		has = true
	}

	return has
}
