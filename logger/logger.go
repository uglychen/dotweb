package logger

import (
	"errors"
	"github.com/devfeel/dotweb/framework/file"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	LogLevel_Debug = "debug"
	LogLevel_Info  = "info"
	LogLevel_Warn  = "warn"
	LogLevel_Error = "error"
)

type AppLog interface {
	SetLogPath(logPath string)
	SetEnabledLog(enabledLog bool)
	Debug(log string, logTarget string)
	Info(log string, logTarget string)
	Warn(log string, logTarget string)
	Error(log string, logTarget string)
	Log(log string, logTarget string, logLevel string)
}

var (
	appLog         AppLog
	DefaultLogPath string
	EnabledLog     bool = true
)

func Logger() AppLog {
	return appLog
}

func SetLogger(logger AppLog) {
	appLog = logger
	logger.SetLogPath(DefaultLogPath)
	logger.SetEnabledLog(EnabledLog)
}

func SetLogPath(path string) {
	DefaultLogPath = path
	if appLog != nil {
		appLog.SetLogPath(path)
	}
}

func SetEnabledLog(isLog bool) {
	EnabledLog = isLog
	if appLog != nil {
		appLog.SetEnabledLog(isLog)
	}
}

func InitLog() {
	if DefaultLogPath == "" {
		DefaultLogPath = file.GetCurrentDirectory()
	}
	if appLog == nil {
		appLog = NewXLog()
	}

	SetLogPath(DefaultLogPath) //set default log path
	SetEnabledLog(EnabledLog)  //set default enabled log
}

//日志内容
// fileName 文件名字
// line 调用行号
// fullPath 文件全路径
// funcName 那个方法进行调用
type logContext struct {
	fileName string
	line     int
	fullPath string
	funcName string
}

//打印
// skip=0  runtime.Caller 的调用者.
// skip=1  runtime/proc.c 的 runtime.main
// skip=2  runtime/proc.c 的 runtime.goexit
//
//Go的普通程序的启动顺序:
//1.runtime.goexit 为真正的函数入口(并不是main.main)
//2.然后 runtime.goexit 调用 runtime.main 函数
//3.最终 runtime.main 调用用户编写的 main.main 函数
func callerInfo(skip int) (ctx *logContext, err error) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return nil, errors.New("error  during runtime.Callers")
	}

	funcInfo := runtime.FuncForPC(pc)
	if funcInfo == nil {
		return nil, errors.New("error during runtime.FuncForPC")
	}

	funcName := funcInfo.Name()
	if strings.HasPrefix(funcName, ".") {
		funcName = funcName[strings.Index(funcName, "."):]
	}

	ctx = &logContext{
		funcName: filepath.Base(funcName),
		line:     line,
		fullPath: file,
		fileName: filepath.Base(file),
	}

	return ctx, nil

}
