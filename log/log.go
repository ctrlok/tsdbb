package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type key int

const KeyOperation key = 0
const KeyUrl key = 1
const KeyClientNum key = 2
const KeyPlannedCount key = 3
const KeyBenchType key = 4
const KeyTSDBType key = 5
const KeyCount key = 6

var Logger *zap.Logger
var SLogger *zap.SugaredLogger
var DebugLevel = false

func ParseFields(ctx context.Context) []zapcore.Field {
	fields := []zapcore.Field{}
	if n, ok := ctx.Value(KeyOperation).(string); ok {
		fields = append(fields, zap.String("operation", n))
	}
	if n, ok := ctx.Value(KeyUrl).(string); ok {
		fields = append(fields, zap.String("dstURL", n))
	}
	if n, ok := ctx.Value(KeyClientNum).(int); ok {
		fields = append(fields, zap.Int("clientNumber", n))
	}
	if n, ok := ctx.Value(KeyPlannedCount).(int); ok {
		fields = append(fields, zap.Int("plannedMessagesToSend", n))
	}
	if n, ok := ctx.Value(KeyBenchType).(string); ok {
		fields = append(fields, zap.String("benchType", n))
	}
	if n, ok := ctx.Value(KeyTSDBType).(string); ok {
		fields = append(fields, zap.String("tsdbType", n))
	}
	if n, ok := ctx.Value(KeyCount).(int); ok {
		fields = append(fields, zap.Int("count", n))
	}
	return fields
}

type LogInfo struct{}

func (l LogInfo) Printf(s string, i ...interface{}) {
	SLogger.Infof(s, i)
}
