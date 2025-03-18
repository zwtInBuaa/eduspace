package logger

import (
	"EDU_TH_2_backend/gin/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.Logger

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.GetString("log.path"), // 日志文件路径，默认 os.TempDir()
		MaxSize:    10,                           // 每个日志文件保存10M，默认 100M
		MaxBackups: 30,                           // 保留30个备份，默认不限
		MaxAge:     7,                            // 保留7天，默认不限
		Compress:   true,                         // 是否压缩，默认不压缩
	}
	//return zapcore.AddSync(lumberJackLogger)
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger)) // 打印到控制台和文件
}

func getEncoder() zapcore.Encoder {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 大写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
	//return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogLevel() zapcore.Level {
	// debug 可以打印出 info debug warn
	// info  级别可以打印 warn info
	// warn  只能打印 warn
	// debug->info->warn->error
	var level zapcore.Level
	logLevel := config.GetString("log.level")
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	return level
}

func InitLogger() {
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(getLogLevel())
	core := zapcore.NewCore(
		getEncoder(),
		getLogWriter(),
		getLogLevel(),
	)

	needCaller := config.GetBool("log.needCaller")

	if needCaller {
		// 开启开发模式，堆栈跟踪（写出调用log语句位置）
		caller := zap.AddCaller()
		// 开启文件及行号
		development := zap.Development()
		// 设置初始化字段,如：添加一个服务器名称
		//filed := zap.Fields(zap.String("serviceName", "serviceName"))
		// 构造日志
		logger = zap.New(core, caller, development)
	} else {
		// 开启文件及行号
		development := zap.Development()
		// 设置初始化字段,如：添加一个服务器名称
		//filed := zap.Fields(zap.String("serviceName", "serviceName"))
		// 构造日志
		logger = zap.New(core, development)
	}

	//logger = zap.New(core, caller, development, filed)
	logger.Info("DefaultLogger init success")
}

func GetZapLogger() *zap.Logger {
	return logger
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	logger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
