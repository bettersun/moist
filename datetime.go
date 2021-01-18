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
	return time.Now().Format(DTFmtYmdHms)
}

func NowMdHms() string {

	//
	return time.Now().Format(DTFmtMdHms)
}

func NowYmdHmsHyphen() string {

	//
	return time.Now().Format(DTFmtYmdHmsHyphen)
}

func NowYmdHmsSlash() string {

	//
	return time.Now().Format(DTFmtYmdHmsSlash)
}

func TodayYmd() string {

	//
	return time.Now().Format(DFmtYmd)
}

func TodayYmdHyphen() string {

	//
	return time.Now().Format(DFmtYmdHyphen)
}

func TodayYmdSlash() string {

	//
	return time.Now().Format(DFmtYmdSlash)
}
