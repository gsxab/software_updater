package logs

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

type LogLevel int8

const (
	AllLevels LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	NoLevels
)

var (
	logLevel = InfoLevel
	errorLog = log.New(os.Stdout, "ERROR", log.Ltime|log.Llongfile)
	warnLog  = log.New(os.Stdout, "WARN", log.Ltime|log.Llongfile)
	infoLog  = log.New(os.Stdout, "INFO", log.Ltime|log.Llongfile)
	debugLog = log.New(os.Stdout, "DEBUG", log.Ltime|log.Llongfile)
)

func SetLevel(level LogLevel) {
	logLevel = level
}

func format(kvs ...any) string {
	sb := &strings.Builder{}
	//_, _ = fmt.Fprintf(sb, "Module=%s Function=%s Message=%s", module, function, message)
	for i := 0; i < len(kvs)-1; i += 2 {
		_, _ = fmt.Fprintf(sb, " %s=%s", kvs[i], kvs[i+1])
	}
	return sb.String()
}

func makeKvLogFunc(level LogLevel, logger *log.Logger) func(ctx context.Context, kvs ...any) {
	return func(ctx context.Context, kvs ...any) {
		if level < logLevel {
			return
		}
		str := format(kvs...)
		_ = logger.Output(3, str)
	}
}

func makeBasicLogFunc(f func(context.Context, ...any)) func(ctx context.Context, kvs ...any) {
	return func(ctx context.Context, kvs ...any) {
		f(ctx, kvs...)
	}
}

func makeErrLogFunc(f func(context.Context, ...any)) func(ctx context.Context, err error, kvs ...any) {
	return func(ctx context.Context, err error, kvs ...any) {
		kvs = append(kvs, "err", err.Error())
		f(ctx, kvs...)
	}
}

//goland:noinspection GoUnusedGlobalVariable
var (
	errorLogger = makeKvLogFunc(ErrorLevel, errorLog)
	warnLogger  = makeKvLogFunc(WarnLevel, warnLog)
	infoLogger  = makeKvLogFunc(InfoLevel, infoLog)
	debugLogger = makeKvLogFunc(DebugLevel, debugLog)

	Error = makeBasicLogFunc(errorLogger)
	Warn  = makeBasicLogFunc(warnLogger)
	Info  = makeBasicLogFunc(infoLogger)
	Debug = makeBasicLogFunc(debugLogger)

	ErrorE = makeErrLogFunc(errorLogger)
	WarnE  = makeErrLogFunc(warnLogger)
	InfoE  = makeErrLogFunc(infoLogger)
	DebugE = makeErrLogFunc(debugLogger)
)
