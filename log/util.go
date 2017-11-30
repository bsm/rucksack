package log

import (
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func parseFields(s string) map[string]interface{} {
	if s == "" {
		return nil
	}

	pairs := strings.Split(s, ",")
	fields := make(map[string]interface{}, len(pairs))

	for _, pair := range pairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 || parts[0] == "" {
			continue
		}

		var v interface{} = parts[1]
		if n, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
			v = n
		}
		fields[parts[0]] = v
	}
	return fields
}

func buildLogger(name, level string, fields map[string]interface{}) (*zap.Logger, error) {
	var config zap.Config

	// Select config
	if name == "" && len(fields) == 0 {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	config.InitialFields = fields

	// Parse level
	if level != "" {
		if v := zap.NewAtomicLevel(); v.UnmarshalText([]byte(level)) == nil {
			config.Level = v
		}
	}

	// Build logger
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	// Set name
	if name != "" {
		logger = logger.Named(name)
	}
	return logger, nil
}
