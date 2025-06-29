package server

import (
	"fmt"

	"github.com/Phillezi/common/utils/or"
	"go.uber.org/zap"
)

func getString(m map[string]interface{}, key string, logger ...*zap.Logger) string {
	l := or.Or(logger...)
	if val, exists := m[key]; exists {
		if s, ok := val.(string); ok {
			return s
		}
		if l != nil {
			l.Warn("Unexpected type for config value",
				zap.String("key", key),
				zap.String("expected", "string"),
				zap.String("actual", fmt.Sprintf("%T", val)),
			)
		}
	}
	return ""
}

func getT[T any](m map[string]interface{}, key string, logger ...*zap.Logger) T {
	var zero T
	l := or.Or(logger...)
	if val, exists := m[key]; exists {
		if s, ok := val.(T); ok {
			return s
		}
		if l != nil {
			l.Warn("Unexpected type for config value",
				zap.String("key", key),
				zap.String("expected",
					fmt.Sprintf("%T", zero)),
				zap.String("actual", fmt.Sprintf("%T", val)),
			)
		}
	}
	return zero
}
