package logger

import (
	"context"
	"log"
	"os"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// GormLogger wraps gorm.io/gorm/logger.Interface
type GormLogger struct {
	LogLevel gormlogger.LogLevel
}

// GetGormLogger returns a GORM logger instance
func GetGormLogger(logLevel gormlogger.LogLevel) gormlogger.Interface {
	return &GormLogger{
		LogLevel: logLevel,
	}
}

// LogMode sets the log level
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info prints info
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		InfoLogger.Printf(msg, data...)
	}
}

// Warn prints warn messages
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		InfoLogger.Printf(msg, data...)
	}
}

// Error prints error messages
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		ErrorLogger.Printf(msg, data...)
	}
}

// Trace prints trace messages
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		ErrorLogger.Printf("[%.3fms] [rows:%v] %s; %s", float64(elapsed.Nanoseconds())/1e6, rows, sql, err)
	} else if l.LogLevel >= gormlogger.Info {
		InfoLogger.Printf("[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}
