package log

import (
	"testing"

	"go.uber.org/zap"
)

func init() {
	Log, _ = zap.NewProduction()
	SLog = Log.Sugar()
}

func BenchmarkZapSugar(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SLog.Debug("string", "key", "value")
	}
}

func BenchmarkZapDefault(b *testing.B) {
	str := zap.String("key", "value")
	for i := 0; i < b.N; i++ {
		Log.Debug("string", str)
	}
}

func BenchmarkZapDefault2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Log.Debug("string", zap.String("key", "value"))
	}
}

func BenchmarkZapEntry(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := Log.Check(zap.DebugLevel, "message")
		if l.Level == zap.DebugLevel {
			Log.Debug("string", zap.String("key", "value"))
		}
	}

}
