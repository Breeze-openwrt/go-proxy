package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log 是全局 SugaredLogger 实例，支持简单的 Printf 风格调用
	Log *zap.SugaredLogger
	// AtomicLevel 用于动态调整日志级别
	AtomicLevel zap.AtomicLevel
)

// Init 初始化日志系统
// vCount: 命令行参数 -v 的个数
// output: 日志输出文件路径
// configLevel: 配置文件中的默认级别
func Init(vCount int, output string, configLevel string) error {
	AtomicLevel = zap.NewAtomicLevel()

	// 1. 决定日志级别 (优先级：命令行 > 配置文件)
	var level zapcore.Level
	switch {
	case vCount >= 2:
		// 我们将 -vv 定义为超越 Debug 的极度详细模式（在 Zap 中依然对应 Debug，但我们可以根据这个标记启用更多追踪）
		level = zap.DebugLevel
	case vCount == 1:
		level = zap.DebugLevel
	default:
		// 解析配置文件中的级别
		switch strings.ToLower(configLevel) {
		case "debug":
			level = zap.DebugLevel
		case "warn":
			level = zap.WarnLevel
		case "error":
			level = zap.ErrorLevel
		default:
			level = zap.InfoLevel
		}
	}
	AtomicLevel.SetLevel(level)

	// 2. 检查目录安全性 (如果指定了输出文件)
	var writeSyncer zapcore.WriteSyncer
	if output != "" && output != "stdout" && output != "stderr" {
		dir := filepath.Dir(output)
		if dir != "." {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				return fmt.Errorf("logging directory does not exist: %s", dir)
			}
		}

		// 3. 配置滚动 (Max 20MB)
		lumberjackLogger := &lumberjack.Logger{
			Filename:   output,
			MaxSize:    20, // megabytes
			MaxBackups: 5,
			MaxAge:     30,   // days
			Compress:   true, // 压缩旧文件
		}
		writeSyncer = zapcore.AddSync(lumberjackLogger)
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 4. 配置 Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 终端带有颜色，文件建议用普通级别
	
	// 如果输出到文件，取消颜色
	if output != "" && output != "stdout" && output != "stderr" {
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		writeSyncer,
		AtomicLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Log = logger.Sugar()

	return nil
}

// 导出便捷方法
func Info(template string, args ...interface{})  { Log.Infof(template, args...) }
func Debug(template string, args ...interface{}) { Log.Debugf(template, args...) }
func Warn(template string, args ...interface{})  { Log.Warnf(template, args...) }
func Error(template string, args ...interface{}) { Log.Errorf(template, args...) }
func Fatal(template string, args ...interface{}) { Log.Fatalf(template, args...) }
