package moist

import (
	"time"
)

// 现在日期（现在时刻）
//   DATE_TIME_Format
//   DATE_TIME_Format_YMD
func Now(format string) string {

	//
	return time.Now().Format(format)
}

func NowYmdHms() string {

	//
	return time.Now().Format(DateTimeFormatYmdHms)
}

func NowYmdHmsHyphen() string {

	//
	return time.Now().Format(DateTimeFormatYmdHmsHyphen)
}

func NowYmdHmsSlash() string {

	//
	return time.Now().Format(DateTimeFormatYmdHmsSlash)
}

func TodayYmd() string {

	//
	return time.Now().Format(DateFormatYmd)
}

func TodayYmdHyphen() string {

	//
	return time.Now().Format(DateFormatYmdHyphen)
}

func TodayYmdSlash() string {

	//
	return time.Now().Format(DateFormatYmdSlash)
}
