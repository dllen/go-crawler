package logger

import (
	"fmt"
	"go.uber.org/zap"
)

func log(log string, level string) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	if level == "info" {
		sugar.Infof(log)
	}
	if level == "debug" {
		sugar.Debugf(log)
	}
	if level == "warn" {
		sugar.Warnf(log)
	}
	if level == "error" {
		sugar.Errorf(log)
	}

}

func Info(v ...interface{}) {
	l := fmt.Sprintf("%v", v)
	log(l, "info")
}

func Debug(v ...interface{}) {
	l := fmt.Sprintf("%v", v)
	log(l, "debug")
}

func Warn(v ...interface{}) {
	l := fmt.Sprintf("%v", v)
	log(l, "warn")
}

func Error(v ...interface{}) {
	l := fmt.Sprintf("%v", v)
	log(l, "error")
}
