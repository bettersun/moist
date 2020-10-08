package moist

import (
	"errors"
	"strconv"
	"unicode"
)

// 行（数字）列（字母）组合单元格
func GetRange(c string, r int) string {
	return c + strconv.Itoa(r)
}

// 从单元格获取行（数字）列（字母）
func GetColumnRow(cell string) (column string, row int, err error) {

	r := []rune(cell)

	iLetter := -1
	iDigit := -1
	fisrtDigit := -1
	for i, c := range r {
		if unicode.IsLetter(c) {
			iLetter = i
		}

		if unicode.IsDigit(c) {
			iDigit = i
		}

		if fisrtDigit == -1 && unicode.IsDigit(c) {
			fisrtDigit = i
		}
	}

	if iLetter == -1 {
		err = errors.New("非正确单元格: 不包含字母")
		return column, row, err
	}
	if iDigit == -1 {
		err = errors.New("非正确单元格: 不包含数字")
		return column, row, err
	}
	if iDigit <= iLetter {
		err = errors.New("非正确单元格: 数字在字母前面")
		return column, row, err
	}

	// 列
	column = cell[0:fisrtDigit]
	// 行转换成整数
	row, err = strconv.Atoi(cell[fisrtDigit:])
	if err != nil {
		return column, row, err
	}

	return column, row, nil
}
