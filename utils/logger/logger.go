package logger

import (
	"fmt"
	"github.com/muchlist/berita_acara/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

var (
	log logger
)

type logger struct {
	log *zap.Logger
}

type loggerInterface interface {
	Printf(format string, v ...interface{})
}

func Init() {
	logConfig := zap.Config{
		OutputPaths: []string{getOutput()},
		Level:       zap.NewAtomicLevelAt(getLevel()),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "lvl",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error
	if log.log, err = logConfig.Build(); err != nil {
		panic(err)
	}
}

// GetLogger digunakan untuk keperluan setLogger di dependency lain,
// misalnya ada di database logger, dia memerlukan logger yang menerapkan interface
// tertentu
func GetLogger() logger {
	return log
}

func getLevel() zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(configs.Config.LOGLEVEL)) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func getOutput() string {
	output := strings.TrimSpace(configs.Config.LOGOUTPUT)
	if output == "" {
		return "stdout"
	}
	return output
}

// Printf diperlukan untuk setLogger di ElasticSearch
func (l logger) Printf(format string, v ...interface{}) {
	if len(v) == 0 {
		Info(format)
	} else {
		Info(fmt.Sprintf(format, v...))
	}
}

func (l logger) Print(format string, v ...interface{}) {
	if len(v) == 0 {
		Info(format)
	} else {
		Info(fmt.Sprintf(format, v...))
	}
}

func Info(msg string, tags ...zap.Field) {
	log.log.Info(msg, tags...)
	_ = log.log.Sync()
}

func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	log.log.Error(msg, tags...)
	_ = log.log.Sync()
}
