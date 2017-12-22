package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path"
	"time"
	"os"
	"strings"
)

type Logger struct{
	logger *zap.Logger

}

func (l *Logger) New() *zap.Logger{
	wd, _ := os.Getwd()
	filepath := path.Join(path.Dir(wd), "/logs/", strings.Replace(time.Now().Format(time.RFC3339), ":", "", -1) )
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath,
		//MaxSize:    500, // megabytes
		//MaxBackups: 3,
		//MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		w,
		zap.InfoLevel,
	)
	logger := zap.New(core)
	l.logger = logger
	return l.logger
}
