package step

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/bettersun/moist"
)

/// 程序设置文件
const CONFIG_FILE = "./config.yml"

/// 统计代码文件列表的代码行数
func Step(file []string) StepSummary {

	// 加载注释配置
	commentConfig, err := LoadCommentConfig(CONFIG_FILE)
	if err != nil {
		log.Printf("LoadRegExpConfig error! \n")
	}

	// 转换为正则表达式
	commentRegExp := ToCommentRegExp(commentConfig)

	// 统计代码行数
	stepInfo := CountAll(file, &commentRegExp)

	// 汇总
	stepResult := Summary(stepInfo)

	return stepResult
}

/// 从注释定义文件读取注释定义
func LoadCommentConfig(commentDefineFile string) ([]CommentDefine, error) {

	// 读取代码注释定义文件
	b, err := ioutil.ReadFile(commentDefineFile)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// 取得程序设置项
	var cmtDef []CommentDefine
	err = yaml.Unmarshal(b, &cmtDef)
	if err != nil {
		log.Printf("Get the setting error! %v\n", err.Error())
	}

	return cmtDef, nil
}

// 统计多文件代码行数
func CountAll(file []string, mCommentRegExp *map[string]CommentRegExp) []StepInfo {

	var stepInfo []StepInfo
	for _, f := range file {
		err := Count(f, mCommentRegExp, &stepInfo)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}

	return stepInfo
}

// 代码行数统计结果
func Summary(stepInfo []StepInfo) StepSummary {

	// 代码行数统计信息
	var stepSummary StepSummary

	// 汇总
	stepSummary.FileCount = len(stepInfo)
	stepSummary.StepInfo = stepInfo

	// 合计
	for _, step := range stepInfo {
		// 总行数
		stepSummary.TotalStep = stepSummary.TotalStep + step.TotalStep
		// 空行总行数
		stepSummary.EmptyLineStep = stepSummary.EmptyLineStep + step.EmptyLineStep
		// 注释总行数
		stepSummary.CommentStep = stepSummary.CommentStep + step.CommentStep
		// 代码总行数
		stepSummary.SourceStep = stepSummary.SourceStep + step.SourceStep
		// 有效总行数(注释+代码)
		stepSummary.ValidStep = stepSummary.ValidStep + step.ValidStep

		// 无注标志定义文件
		if !step.CommentRuleDefined {
			stepSummary.FlatFile = append(stepSummary.FlatFile, step.File)
		}
		// 未统计文件
		if !step.Counted {
			stepSummary.UnCountedFile = append(stepSummary.UnCountedFile, step.File)
		}
	}

	return stepSummary
}

// 统计单文件代码行数
func Count(file string, mCommentRegExp *map[string]CommentRegExp, stepInfo *[]StepInfo) error {

	// 取得代码文件对应的注释标志定义
	var cmtRegExp CommentRegExp
	var isDefined bool = false
	for k, v := range *mCommentRegExp {
		if moist.HasSuffixIgnoreCase(file, k) {
			cmtRegExp = v
			isDefined = true
			break
		}
	}

	var info StepInfo
	if isDefined {
		// 存在注释标志定义
		info.CommentRuleDefined = true
	} else {
		// 无该类型代码文件注释标志定义
		// 默认注释正则表达式
		defaultRegEx(&cmtRegExp)

		info.File = file
		info.FileName = moist.FileName(file, moist.S_SLASH)
		// BLogger.Warn(Message["CMN_W0003"], file)
		log.Printf("[%v] 无该类型代码文件注释标志定义, 不统计注释行数。", file)
		info.ExInfo = "无注释标志定义, 不统计注释行数。"
		// 无注释标志定义
		info.CommentRuleDefined = false
	}

	var totalStep int     // 总行数
	var emptyLineStep int // 空行数
	var commentStep int   // 注释
	var sourceStep int    // 代码行数

	var isMultiLineComment bool // 多行注释统计中标志
	var matchIndex int          // 多行注释结束 正则表达式 下标

	var isMatch bool

	// 打开文件
	f, err := os.Open(file)
	if err != nil {
		info.File = file
		info.FileName = moist.FileName(file, moist.S_SLASH)
		info.ExInfo = err.Error()
		info.Counted = false

		*stepInfo = append(*stepInfo, info)
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	// 行读取
	for {
		line, err := reader.ReadString('\n')

		// 读取出错,但非文件结尾
		if err != nil && err != io.EOF {
			log.Printf("Read Line Error: %v\n", err)
			break
		}

		// 总行数
		totalStep++
		if !isMultiLineComment {

			if cmtRegExp.RegExEmptyLine.MatchString(line) { //空行
				emptyLineStep++
			} else if _, isMatchSingle := MatchIn(line, cmtRegExp.RegExSingleLine); cmtRegExp.HasSingleLineMark && isMatchSingle { // 单行注释
				commentStep++
			} else if _, isMatchSingleStartEnd := MatchIn(line, cmtRegExp.RegExSingleLineStartEnd); cmtRegExp.HasMultiLineMark && isMatchSingleStartEnd { // 多行注释标志开始结束的单行代码
				commentStep++
			} else if matchIndex, isMatch = MatchIn(line, cmtRegExp.RegExMultiLineStart); cmtRegExp.HasMultiLineMark && isMatch { // 多行注释 开始
				commentStep++
				isMultiLineComment = true
			}
		} else if cmtRegExp.HasMultiLineMark { // 有多行注释

			if isMultiLineComment {
				if cmtRegExp.RegExEmptyLine.MatchString(line) { //多行注释里的空行
					emptyLineStep++
				} else {
					//多行注释
					commentStep++

					//多行注释结束
					if cmtRegExp.RegExMultiLineEnd[matchIndex].MatchString(line) { //多行注释 结束
						isMultiLineComment = false
					}
				}
			}
		}

		// 文件结尾
		if err == io.EOF {
			break
		}
	}

	// 代码行数
	sourceStep = totalStep - commentStep - emptyLineStep

	info.File = file
	info.FileName = moist.FileName(file, moist.S_SLASH)
	info.TotalStep = totalStep
	info.EmptyLineStep = emptyLineStep
	info.SourceStep = sourceStep
	info.CommentStep = commentStep
	// 有效行数（注释 + 代码）
	info.ValidStep = totalStep - emptyLineStep
	info.Counted = true

	*stepInfo = append(*stepInfo, info)
	return nil
}

/// 注释定义转注释正则表达式
func ToCommentRegExp(cmtDef []CommentDefine) map[string]CommentRegExp {

	// 相同注释定义按照文件类型整理到Map
	var mCmtRegex map[string]CommentRegExp
	mCmtRegex = make(map[string]CommentRegExp)
	for _, v := range cmtDef {

		for _, ext := range v.FileExtension {
			_, ok := mCmtRegex[ext]
			if !ok {
				mCmtRegex[ext] = ToRegExp(ext, v)
			}
		}
	}

	return mCmtRegex
}

// 转换为正则表达式
func ToRegExp(ext string, def CommentDefine) CommentRegExp {

	var cmtRegExp CommentRegExp

	lenSingle := len(def.SingleLine)
	lenMuilti := len(def.MultiLineStart)

	// 分配切片空间
	cmtRegExp.RegExSingleLine = make([]*regexp.Regexp, lenSingle)
	cmtRegExp.RegExSingleLineStartEnd = make([]*regexp.Regexp, lenMuilti)
	cmtRegExp.RegExMultiLineStart = make([]*regexp.Regexp, lenMuilti)
	cmtRegExp.RegExMultiLineEnd = make([]*regexp.Regexp, lenMuilti)

	// 扩展名
	cmtRegExp.FileExtension = ext

	if !moist.IsAllEmpty(def.SingleLine) {
		// 有单行注释
		cmtRegExp.HasSingleLineMark = true
	}
	if !moist.IsAllEmpty(def.MultiLineStart) {
		// 有多行注释
		cmtRegExp.HasMultiLineMark = true
	}

	// 空行正则表达式
	cmtRegExp.RegExEmptyLine = regexp.MustCompile(`^[\s]*$`)

	// 单行正则表达式
	for i, v := range def.SingleLine {

		cmtRegExp.RegExSingleLine[i] =
			regexp.MustCompile(`^[\s]*` + v + `.*`)
	}

	// 单行正则表达式 注释开始符和注释结束符在同一行
	for i, _ := range def.MultiLineStart {
		start := Escape(def.MultiLineStart[i])
		end := Escape(def.MultiLineEnd[i])

		cmtRegExp.RegExSingleLineStartEnd[i] =
			regexp.MustCompile(`^[\s]*(` + start + `).*(` + end + `)[\s]*$`)
	}

	// 多行正则表达式 开始
	for i, v := range def.MultiLineStart {
		v = Escape(v)
		cmtRegExp.RegExMultiLineStart[i] =
			regexp.MustCompile(`^[\s]*(` + v + `).*`)
	}

	// 多行正则表达式 结束
	for i, v := range def.MultiLineEnd {
		v = Escape(v)
		cmtRegExp.RegExMultiLineEnd[i] =
			regexp.MustCompile(`.*(` + v + `)[\s]*$`)
	}

	return cmtRegExp
}

// 多个标志的正则表达式确认
func MatchIn(line string, regexList []*regexp.Regexp) (matchIndex int, isMatch bool) {

	isMatch = false
	for i, v := range regexList {

		if v.MatchString(line) {

			// log.Printf("=========================================")
			// log.Printf(line)
			// log.Printf(v.String())

			isMatch = true
			matchIndex = i
			break
		}
	}

	return matchIndex, isMatch
}

// 默认正则表达式（不统计注释）
func defaultRegEx(cmtRegExp *CommentRegExp) {

	// 有单行注释标志
	cmtRegExp.HasSingleLineMark = true
	// 有多行注释标志
	cmtRegExp.HasMultiLineMark = true
	// 空行正则表达式(字符串使用锐音符(`)不需要转义))
	cmtRegExp.RegExEmptyLine = regexp.MustCompile(`^[\s]*$`)
	// 单行正则表达式(.包括换行符，不需要在最后加上$)
	cmtRegExp.RegExSingleLine = nil
	// 单行正则表达式 开始结束
	cmtRegExp.RegExSingleLineStartEnd = nil
	// 多行正则表达式 开始
	cmtRegExp.RegExMultiLineStart = nil
	// 多行正则表达式 结束
	cmtRegExp.RegExMultiLineEnd = nil
}

// 转义
func Escape(sRegExp string) string {
	// return strings.ReplaceAll(strings.ReplaceAll(resexp, `*`, `\*`), `/`, `/`)

	v := sRegExp
	v = strings.ReplaceAll(v, `*`, `\*`)
	v = strings.ReplaceAll(v, `/`, `\/`)

	return v
}
