package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// FIXME: глобальная переменная. Подумать что можно с ней сделать.
var Zlog *zap.SugaredLogger

const (
	colorCyan = "\033[36m"
)

func colorTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	coloredTime := fmt.Sprintf("%s%s", colorCyan, t.Format("2006-01-02 15:04:05"))
	enc.AppendString(coloredTime)
}

func InitLogger() *zap.Logger {

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		EncodeTime:       colorTimeEncoder,
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: " | ",
	}

	// Создаём core с цветной консолью
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)
	core := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)

	// Создаём логгер
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	Zlog = logger.Sugar()
	return logger
}

func Error(msg string, args ...any) {
	checkLoggerInitialized()
	Zlog.Errorln(msg, args)
}

func Info(msg string, args ...any) {
	checkLoggerInitialized()
	Zlog.Infoln(msg, args)
}

func Debug(msg string, args ...any) {
	checkLoggerInitialized()
	Zlog.Debugln(msg, args)
}

func Warn(msg string, args ...any) {
	checkLoggerInitialized()
	Zlog.Warnln(msg, args)
}

func Fatal(msg string, args ...any) {
	checkLoggerInitialized()
	Zlog.Fatalln(msg, args)
}

func checkLoggerInitialized() {
	if Zlog == nil {
		panic("logger not initialized")
	}
}
