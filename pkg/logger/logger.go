package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"practice/pkg/conf"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger     AppLogger
	LoggerOnce sync.Once
)

type AppLogger struct {
	logger logr.Logger
}

func (l AppLogger) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l AppLogger) Info(str string) {
	l.logger.Info(str)
}

func (l AppLogger) Warning(str string) {
	l.WithName("warning").Info(str)
}

func (l AppLogger) Warningf(format string, args ...interface{}) {
	l.WithName("warning").Info(fmt.Sprintf(format, args...))
}

func (l AppLogger) Error(err error, str string) {
	l.logger.Error(err, str)
}
func (l AppLogger) Errorf(err error, format string, args ...interface{}) {
	l.logger.Error(err, fmt.Sprintf(format, args...))
}

func (l AppLogger) FatalError(err error, message string) {
	l.logger.Error(err, message)
	log.Fatal(err)
}

func (l AppLogger) WithName(name string) AppLogger {
	l.logger = l.logger.WithName(name)
	return l
}

//Logger return a logr logger with name
func Logger(name string) AppLogger {
	return getLogger().WithName(name)
}

func NewZapLogger() *zap.Logger {
	cfg := conf.GetConfig()
	logDir := cfg.LogConf.Path
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}
	//LogFile := filepath.Join(logDir, fmt.Sprintf("%s-%s", "operator", time.Now().Format("20060102-150405")))
	LogFile := filepath.Join(logDir, "App.log")

	hook := lumberjack.Logger{
		Filename:   LogFile,                // 日志文件路径
		MaxSize:    cfg.LogConf.MaxSize,    // 每个日志文件保存的大小 单位:M
		MaxAge:     cfg.LogConf.MaxAge,     // 文件最多保存多少天
		MaxBackups: cfg.LogConf.MaxBackups, // 日志文件最多保存多少个备份
		Compress:   cfg.LogConf.Compress,   // 是否压缩
	}
	zapConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "file",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	logLeve := zap.InfoLevel
	atomicLevel := zap.NewAtomicLevel()
	writes = append(writes, zapcore.AddSync(os.Stdout))
	if os.Getenv("TDC_OPERATOR_DEV_MODE") == "ON" {
		logLeve = zap.DebugLevel
	}
	atomicLevel.SetLevel(logLeve)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapConfig),
		zapcore.NewMultiWriteSyncer(writes...),
		//zapcore.AddSync(&hook),
		atomicLevel,
	)

	// 开启堆栈跟踪
	caller := zap.AddCaller()

	// 设置初始化字段
	field := zap.Fields(zap.String("tdc", "App"))
	option := zap.Development()
	option2 := zap.AddCallerSkip(1)
	zapLogger := zap.New(core, caller, option, field, option2)
	return zapLogger
}

func InitLogger() {
	LoggerOnce.Do(func() {
		logger = AppLogger{zapr.NewLogger(NewZapLogger())}
	})
}

func getLogger() AppLogger {
	return logger
}

func NewLogWriter(logger AppLogger, fd string) *logWriter {
	return &logWriter{
		buf:    new(bytes.Buffer),
		logger: logger.logger.WithValues("fd", fd),
	}
}

type logWriter struct {
	buf    *bytes.Buffer
	logger logr.Logger
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	l.logger.Info(strings.TrimSpace(string(p)))
	return l.buf.Write(p)
}

func (l *logWriter) String() string {
	return l.buf.String()
}
