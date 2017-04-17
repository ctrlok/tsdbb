package log

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log, _ = zap.NewProduction()

func BenchmarkEmpty(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Debug("string")
	}
}

func BenchmarkKV(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Debug("string", zap.String("key", "value"))
	}
}

var ctx = context.WithValue(context.Background(), "key", "value")

func debugWContext(ctx context.Context) {
	fields := []zapcore.Field{}
	if n, ok := ctx.Value("key").(string); ok {
		fields = append(fields, zap.String("key", n))
	}
	log.Debug("string", fields...)
}
func BenchmarkKVWithContext(b *testing.B) {
	for n := 0; n < b.N; n++ {
		debugWContext(ctx)
	}
}

func funcString(ctx context.Context) string {
	if n, ok := ctx.Value("key").(string); ok {
		return n
	}
	return ""
}

func BenchmarkKVWithFuncString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		log.Debug("string", zap.String("key", funcString(ctx)))
	}
}

var debug = false

func debugWContextBool(ctx context.Context) {
	if !debug {
		return
	}
	fields := []zapcore.Field{}
	if n, ok := ctx.Value("key").(string); ok {
		fields = append(fields, zap.String("key", n))
	}
	log.Debug("string", fields...)
}

func BenchmarkKVWithContextBool(b *testing.B) {
	for n := 0; n < b.N; n++ {
		debugWContextBool(context.WithValue(context.Background(), "key", "value"))
	}
}
