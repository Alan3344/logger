package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// 自定义日志输出
type Logger struct {
}

type ConsoleColor struct {
	Gray    func(interface{}) string
	Red     func(interface{}) string
	Green   func(interface{}) string
	Yellow  func(interface{}) string
	Blue    func(interface{}) string
	Magenta func(interface{}) string
	Cyan    func(interface{}) string
	White   func(interface{}) string
	Default func(interface{}) string
}

var cc = func(consoleFgColorValue int, msg interface{}) string {
	return fmt.Sprintf("\033[%dm%v\033[0m", consoleFgColorValue, msg)
}

var CC ConsoleColor = ConsoleColor{
	Gray: func(msg interface{}) string {
		return cc(30, msg)
	},
	Red: func(msg interface{}) string {
		return cc(31, msg)
	},
	Green: func(msg interface{}) string {
		return cc(32, msg)
	},
	Yellow: func(msg interface{}) string {
		return cc(33, msg)
	},
	Blue: func(msg interface{}) string {
		return cc(34, msg)
	},
	Magenta: func(msg interface{}) string {
		return cc(35, msg)
	},
	Cyan: func(msg interface{}) string {
		return cc(36, msg)
	},
	White: func(msg interface{}) string {
		return cc(37, msg)
	},
	Default: func(msg interface{}) string {
		return cc(0, msg)
	},
}

/**
 * 打印日志
 * isError 是否是错误日志 打印错误，错误文本以红色着色
 * deep 调用深度
 * messsages 要打印的信息 任意类型
 */
func Fprint(anyColor func(interface{}) string, deep int, args ...interface{}) {
	if len(args) == 0 {
		return
	}
	var msgs []string = []string{}
	stamp := "[" + CC.Green(time.Now().Format(time.TimeOnly)) + "]"

	pc, file, line, ok := runtime.Caller(deep)
	funcName := runtime.FuncForPC(pc).Name()
	funcName = filepath.Ext(funcName)
	funcName = funcName[1:]
	if !ok {
		fmt.Println(strings.Join(msgs, " "))
		return
	}
	if funcName == "" || funcName == "0" || funcName == "1" {
		funcName = "init"
	}

	dir, err := os.Getwd()
	if err != nil {
		msgs = append(msgs, CC.Red(strings.Join([]string{funcName, err.Error()}, " ")))
		fmt.Println(msgs)
	} else {
		basename := filepath.Base(dir)
		pths := strings.SplitN(file, basename, 2)
		if len(pths) > 2 {
			file = pths[len(pths)-2] + basename + pths[len(pths)-1]
		} else {
			file = pths[len(pths)-1]
		}
		// msgs = append(msgs, fmt.Sprintf("\033[36m%11s:%-4d\033[0m \033[33m%-9s\033[0m", file, line, funcName))
		pathline := fmt.Sprintf("%s:%s", CC.Cyan(file), CC.Blue(line))
		msgs = append(msgs, fmt.Sprintf("%-35s", pathline)+" "+CC.Yellow(funcName))
		// msgs = append(msgs, fmt.Sprintf("\033[36m%11s:%-4d\033[0m", file, line))
	}

	// 用户传入的要打印的参数
	for _, msg := range args {
		msgs = append(msgs, anyColor(fmt.Sprintf("%v", msg)))
	}

	msgs = append([]string{stamp}, msgs...)
	ret := fmt.Sprint(strings.Join(msgs, " "))
	fmt.Println(ret)
}

func (_l *Logger) Error(args ...interface{}) {
	Fprint(CC.Red, 2, args)
}

func (_l *Logger) Info(args ...interface{}) {
	Fprint(CC.Default, 2, args)
}

func (_l *Logger) High(args ...interface{}) {
	Fprint(CC.Yellow, 2, args)
}
