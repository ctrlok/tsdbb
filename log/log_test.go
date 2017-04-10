package log

import (
	"testing"

	"go.uber.org/zap"
)

func init() {
	Logger, _ = zap.NewProduction()
	SugaredLogger = Logger.Sugar()
}

func BenchmarkZapSugar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SugaredLogger.Debug("string", "key", "value")
	}
}

func BenchmarkZapDefault(b *testing.B) {
	str := zap.String("key", "value")
	for i := 0; i < b.N; i++ {
		Logger.Debug("string", str)
	}
}

func BenchmarkZapDefault2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Logger.Debug("string", zap.String("key", "value"))
	}
}

func BenchmarkZapEntry(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := Logger.Check(zap.DebugLevel, "message")
		if l.Level == zap.DebugLevel {
			Logger.Debug("string", zap.String("key", "value"))
		}
	}

}
