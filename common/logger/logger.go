package logger

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinish/lumberjack.v2"
)

var (
	logger     *Logger
	once       sync.Once
	LogDir     = "logs"
	LogLevel   = "info"
)

// Logger wraps standard logging with file rotation support
type Logger struct {
	mu      sync.Mutex
	writers []io.Writer
}

// Init initializes the global logger with file and stdout output
func Init(logDir string, logLevel string) {
	once.Do(func() {
		LogDir = logDir
		LogLevel = logLevel

		if err := os.MkdirAll(logDir, 0755); err != nil {
			panic("failed to create log directory: " + err.Error())
		}

		rotatingLog := &lumberjack.Logger{
			Filename:   filepath.Join(logDir, "new-api.log"),
			MaxSize:    50, // megabytes - reduced from 100 to keep disk usage lower
			MaxBackups: 5,  // keep fewer backups since this is a personal instance
			MaxAge:     14, // days - 2 weeks is enough for personal use
			Compress:   true,
		}

		logger = &Logger{
			writers: []io.Writer{os.Stdout, rotatingLog},
		}

		// Configure gin logging to use our writer
		gin.DefaultWriter = io.MultiWriter(logger.writers...)
		gin.DefaultErrorWriter = io.MultiWriter(logger.writers...)
	})
}

// write formats and writes a log entry to all configured writers
func (l *Logger) write(level, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006/01/02 15:04:05")
	entry := timestamp + " [" + level + "] " + msg + "\n"

	for _, w := range l.writers {
		_, _ = io.WriteString(w, entry)
	}
}

// SysLog logs a system-level informational message
func SysLog(msg string) {
	if logger == nil {
		Init(LogDir, LogLevel)
	}
	logger.write("INFO", msg)
}

// SysError logs a system-level error message
func SysError(msg string) {
	if logger == nil {
		Init(LogDir, LogLevel)
	}
	logger.write("ERROR", msg)
}

// SysWarn logs a system-level warning message
func SysWarn(msg string) {
	if logger == nil {
		Init(LogDir, LogLevel)
	}
	logger.write("WARN", msg)
}

// SysDebug logs a debug message (only when log level is debug)
func SysDebug(msg string) {
	if LogLevel != "debug" {
		return
	}
	if logger == nil {
		Init(LogDir, LogLevel)
	}
	logger.write("DEBUG", msg)
}
