package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(filepath string, size, backups, age int, compress bool) *zap.SugaredLogger {
	now := time.Now()
	filename := fmt.Sprintf("%s/%04d-%02d-%02d.log", filepath, now.Year(), now.Month(), now.Day())
	hook := &lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    size,     // 最大尺寸, M
		MaxBackups: backups,  // 备份数
		MaxAge:     age,      // 存放天数
		Compress:   compress, // 是否压缩
	}
	defer hook.Close()

	enConfig := zap.NewProductionEncoderConfig() // 生成配置
	// 时间格式
	enConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(enConfig),                                            // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook)), // 打印到控制台和文件
		zapcore.InfoLevel,                                                              // 日志等级
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}
