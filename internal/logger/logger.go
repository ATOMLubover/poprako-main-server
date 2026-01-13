// Package logger provides logging functionalities for the application.
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() *zap.Logger {
	var lvl zapcore.Level

	switch os.Getenv("GO_LOG") {
	case "debug":
		lvl = zapcore.DebugLevel
	case "info":
		lvl = zapcore.InfoLevel
	case "warn":
		lvl = zapcore.WarnLevel
	case "error":
		lvl = zapcore.ErrorLevel
	default:
		lvl = zapcore.InfoLevel
	}

	atmLvl := zap.NewAtomicLevelAt(lvl)

	encCfg := zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		MessageKey:    "M",
		StacktraceKey: "S",
		LineEnding:    zapcore.DefaultLineEnding,
		// NOTICE: show colors in console output.
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	cnslEnc := zapcore.NewConsoleEncoder(encCfg)
	wrt := zapcore.Lock(os.Stdout)

	core := zapcore.NewCore(cnslEnc, wrt, atmLvl)

	lgr := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	lgr.Info("Logger initialized", zap.String("level", atmLvl.String()))

	return lgr
}
