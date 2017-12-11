package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func buildLogger(name, level string, enableStackTrace bool, fields map[string]interface{}) (*zap.Logger, error) {
	var config zap.Config

	// Select config
	if name == "" && len(fields) == 0 {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	config.InitialFields = fields
	config.DisableStacktrace = !enableStackTrace

	// Parse level
	if level != "" {
		if v := zap.NewAtomicLevel(); v.UnmarshalText([]byte(level)) == nil {
			config.Level = v
		}
	}

	// Build logger
	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	// Set name
	if name != "" {
		logger = logger.Named(name)
	}
	return logger, nil
}
