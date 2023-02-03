package grtmp

import (
	"io"
	"log"
	"os"
)

type LevelOutputLogger interface {
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
}

// Logger
// 基于log标准库设计的轻量日志工具，在标准库基础上仅添加了简单的日志级别功能
type Logger struct {
	l     *log.Logger
	level LoggerLevel
}

type LoggerLevel int

const (
	Silent LoggerLevel = iota
	Error
	Warn
	Info
)

const (
	errPrefix  = " [Error] "
	warnPrefix = " [Warn] "
	infoPrefix = " [Info] "
)

var logger Logger

func init() {
	logger.l = log.New(os.Stdout, "[grtmp] ", log.LstdFlags|log.Lmicroseconds)
	SetLogLevel(Info)
}

// SetLogOutput
// 设定rtmp日志输出
func SetLogOutput(writer io.Writer) {
	logger.l.SetOutput(writer)
}

// SetLogLevel
// 设定rtmp日志级别，低于设定级别的日志会被打印
func SetLogLevel(level LoggerLevel) {
	logger.level = level
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.level >= Error {
		logger.l.Printf(errPrefix+format, v...)
	}
}

func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level >= Warn {
		logger.l.Printf(warnPrefix+format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level >= Info {
		logger.l.Printf(infoPrefix+format, v...)
	}
}
