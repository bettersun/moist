package step

import (
	"regexp"
)

// 代码注释统计用正则表达式定义
type CommentDefine struct {
	// 扩展名(多个)
	FileExtension []string `yaml:"fileExtension"`
	// 单行正则表达式
	SingleLine []string `yaml:"singleLine"`
	// 多行正则表达式 开始
	MultiLineStart []string `yaml:"multiLineStart"`
	// 多行正则表达式 结束
	MultiLineEnd []string `yaml:"multiLineEnd"`
}

// 代码注释统计用正则表达式定义
type CommentRegExp struct {
	// 扩展名
	FileExtension string
	// 有单行注释
	HasSingleLineMark bool
	// 有多行注释
	HasMultiLineMark bool

	// 空正则表达式
	RegExEmptyLine *regexp.Regexp
	// 单行正则表达式
	RegExSingleLine []*regexp.Regexp
	// 写在单行的多行正则表达式 开始结束
	RegExSingleLineStartEnd []*regexp.Regexp
	// 多行正则表达式 开始
	RegExMultiLineStart []*regexp.Regexp
	// 多行正则表达式 结束
	RegExMultiLineEnd []*regexp.Regexp
}

// 代码行数信息
type StepInfo struct {
	CommentRuleDefined bool   `json:"commentRuleDefined"` // 存在注释标志定义 true:存在 false:不存在
	File               string `json:"file"`               // 文件名（全路径）
	FileName           string `json:"fileName"`           // 文件名
	TotalStep          int    `json:"totalStep"`          // 总行数
	EmptyLineStep      int    `json:"emptyLineStep"`      // 空行数
	CommentStep        int    `json:"commentStep"`        // 注释行数
	SourceStep         int    `json:"sourceStep"`         // 代码行数
	ValidStep          int    `json:"validStep"`          // 有效行数(注释+代码)
	ExInfo             string `json:"exInfo"`             // 扩展信息
	Counted            bool   `json:"counted"`            // 已统计标志 true:已统计 false:未统计
}

// 代码行数信息汇总
type StepSummary struct {
	StepInfo      []StepInfo `json:"stepInfo"`      // 代码行数统计结果
	FlatFile      []string   `json:"flatFile"`      // 无注标志定义文件
	UnCountedFile []string   `json:"unCountedFile"` // 未统计文件一览
	FileCount     int        `json:"fileCount"`     // 文件总数
	TotalStep     int        `json:"totalStep"`     // 总行数
	EmptyLineStep int        `json:"emptyLineStep"` // 空行总行数
	CommentStep   int        `json:"commentStep"`   // 注释总行数
	SourceStep    int        `json:"sourceStep"`    // 代码总行数
	ValidStep     int        `json:"validStep"`     // 有效总行数(注释+代码)
}
