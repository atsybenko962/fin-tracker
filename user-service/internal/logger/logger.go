package logger

import (
	"github.com/atsybenko962/task-platform/user-service/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	Debug  = "debug"
	Info   = "info"
	Empty  = ""
	Warn   = "warn"
	Error  = "error"
	Dpanic = "dpanic"
	Panic  = "panic"
	Fatal  = "fatal"
)

func NewLogger(conf config.Config, sync zapcore.WriteSyncer) *zap.Logger {
	levels := map[string]zapcore.Level{
		Debug:  zapcore.DebugLevel,
		Info:   zapcore.InfoLevel,
		Empty:  zapcore.InfoLevel,
		Warn:   zapcore.WarnLevel,
		Error:  zapcore.ErrorLevel,
		Dpanic: zapcore.DPanicLevel,
		Panic:  zapcore.PanicLevel,
		Fatal:  zapcore.FatalLevel,
	}

	zapConf := zap.NewProductionConfig()
	zapConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConf.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	atom := zap.NewAtomicLevelAt(levels[conf.LogLevel])
	zapConf.Level = atom
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapConf.EncoderConfig),
		sync,
		atom,
	)
	logger := zap.New(core)
	return logger.Named(conf.AppName)

}
