package utils

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// 打印参数
func Info(v ...interface{}) {
	format := strings.Repeat("%v ", len(v))
	msg := fmt.Sprintf(format, v...)
	_info(2, msg, true)
}

// 打印参数
func InfoNoLine(v ...interface{}) {
	format := strings.Repeat("%v ", len(v))
	msg := fmt.Sprintf(format, v...)
	_info(2, msg, false)
}

// 根据格式打印参数
func Infof(format string, v ...interface{}) {
	// format := strings.Repeat("%v ", len(v))
	msg := fmt.Sprintf(format, v...)
	_info(2, msg, true)
}

// // 检查错误,如果err不为空则打印返回true
// func IsErr(err error) bool {
// 	if err != nil {
// 		msg := fmt.Sprintf("错误信息: %v ", err)
// 		_info(2, msg)
// 		return true
// 	}
// 	return false
// }

// 检查错误,如果err不为空则打印返回true
// errInfo : 附加错误信息
func IsErr(err error) bool {
	if err != nil {
		format := "错误信息:"
		msg := fmt.Sprintf(format+" %v ", err)
		_info(2, msg, true)
		return true
	}
	return false
}

// // 检查错误,如果err不为空则打印返回true
// // errInfo : 附加错误信息
// func IsErr(err error, errInfo ...interface{}) bool {
// 	if err != nil {
// 		format := "错误信息:"
// 		if len(errInfo) > 0 {
// 			format += strings.Repeat(" %v >>>", len(errInfo))
// 		}
// 		errInfo = append(errInfo, err)
// 		msg := fmt.Sprintf(format+" %v ", errInfo...)
// 		_info(2, msg)
// 		return true
// 	}
// 	return false
// }

// 检查错误,如果err不为空则打印并退出
func ExitWithErr(err error) bool {
	if err != nil {
		msg := fmt.Sprintf("错误信息: %v ", err)
		_info(2, msg, true)
		syscall.Exit(557)
		return true
	}
	return false
}

func _info(step int, msg string, isNewLine bool) {
	now := time.Now().Format("2006-01-02 15:04:05")
	_, filePath, line, ok := runtime.Caller(step)
	if !ok {
		filePath = "????"
		line = 0
	}
	//获取文件名
	_, file := path.Split(filePath)
	if isNewLine {
		fmt.Println(now + " " + "[" + file + ":" + strconv.Itoa(line) + "]" + msg)
	} else {
		fmt.Print(now + " " + "[" + file + ":" + strconv.Itoa(line) + "]" + msg)
	}
}

// defer utils.GetFuncRunTime(time.Now())
// 获取运行的方法及运行时间
func GetFuncRunTime(t time.Time) (funcName, runTime string) {
	msg := ""
	now := time.Now().Format("2006-01-02 15:04:05")
	pc, filePath, line, ok := runtime.Caller(1)
	if !ok {
		filePath = "????"
		line = 0
		msg = "获取执行时间信息异常"
	}
	//获取文件名
	_, file := path.Split(filePath)
	if msg == "" {
		msg = "方法名: " + runtime.FuncForPC(pc).Name() + "\t" + "运行时间: " + time.Since(t).String()
	}
	fmt.Println(now + " " + "[" + file + ":" + strconv.Itoa(line) + "]" + msg)
	return
}
