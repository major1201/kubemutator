package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	initLog("stdout", "stderr", zapcore.DebugLevel)
}

func initLog(stdout, stderr string, level zapcore.Level) {
	logger, _ := zap.Config{
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{stdout},
		ErrorOutputPaths: []string{stderr},
	}.Build()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
}
