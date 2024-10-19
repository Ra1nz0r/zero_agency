package logger

import (
	"path/filepath"
	"strings"
	"time"

	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapService interface {
	Debug(fields ...interface{})
	Info(fields ...interface{})
	Error(fields ...interface{})
	Fatal(fields ...interface{})
}

type ZapStorage struct {
	*zap.Logger
}

var Zap ZapService = &ZapStorage{zap.NewNop()}

// Initialize инициализирует логгер.
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return fmt.Errorf("parse atomic level error: %w", err)
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level = lvl
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	config.DisableStacktrace = true

	// Кастомный энкодер для вызова (caller)
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeCaller = zapcore.CallerEncoder(func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		// Форматируем вызов с добавлением стабильных отступов
		formattedCaller := formatCallerWithPadding(caller.File, caller.Line, 20)
		enc.AppendString(formattedCaller)
	})

	logger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return fmt.Errorf("logger build error: %w", err)
	}

	Zap = &ZapStorage{logger}

	return nil
}

// Debug логирует сообщения уровня DEBUG.
func (z *ZapStorage) Debug(fields ...interface{}) {
	z.Logger.Sugar().Debug(fields...)
}

// Info логирует сообщения уровня INFO.
func (z *ZapStorage) Info(fields ...interface{}) {
	z.Logger.Sugar().Info(fields...)
}

// Error логирует сообщения уровня ERROR.
func (z *ZapStorage) Error(fields ...interface{}) {
	z.Logger.Sugar().Error(fields...)
}

// Fatal логирует сообщения уровня FATAL.
func (z *ZapStorage) Fatal(fields ...interface{}) {
	z.Logger.Sugar().Fatal(fields...)
}

// formatCallerWithPadding форматирует имя файла и номер строки с фиксированной длиной
func formatCallerWithPadding(file string, line int, width int) string {
	// Получаем только имя файла без пути
	fileName := filepath.Base(file)
	callerInfo := fmt.Sprintf("%s:%d", fileName, line)

	if len(callerInfo) > width {
		// Если строка слишком длинная, обрезаем её с начала
		return "..." + callerInfo[len(callerInfo)-width+3:]
	}
	// Добавляем пробелы, если строка короче
	return callerInfo + strings.Repeat(" ", width-len(callerInfo))
}
