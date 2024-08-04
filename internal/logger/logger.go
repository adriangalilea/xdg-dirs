package logger

// Rationale:
// The tool is designed to be evaluated by the shell, not to modify its own environment.
// This means every single output from the binary should be silent unless `-debug` is specified,
// as any single output will be `eval`'d (executed). Therefore, we only output log.Export(exports)
// to stdout and log all other messages to ~/xdg.log unless in debug mode.

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	"sync"

	"xdg-dirs/internal/conf"
	"github.com/charmbracelet/lipgloss"
)

const maxLogFileSize = 10 * 1024 * 1024 // 10 MB

var (
	timestampStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	messageStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
)

type Logger struct {
	debug        bool
	logFile      *os.File
	logFilePath  string
	infoLogger   *log.Logger
	debugLogger  *log.Logger
	errorLogger  *log.Logger
	exportLogger *log.Logger
	mu           sync.RWMutex
	writers      struct {
		info   io.Writer
		debug  io.Writer
		error  io.Writer
		export io.Writer
	}
}

func NewLogger(debug bool, logFilePath string) *Logger {
	var logFile *os.File
	var err error

	if logFilePath == "" {
		logFilePath = conf.DefaultLogFilePath
	}

	// Ensure the directory exists
	err = os.MkdirAll(filepath.Dir(logFilePath), 0755)
	if err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	logger := &Logger{
		debug:        debug,
		logFile:      logFile,
		logFilePath:  logFilePath,
		exportLogger: log.New(os.Stdout, "", 0),
	}

	logger.updateWriters()

	return logger
}

func (l *Logger) updateWriters() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.debug {
		l.writers.info = io.MultiWriter(os.Stdout, l.logFile)
		l.writers.debug = io.MultiWriter(os.Stdout, l.logFile)
		l.writers.error = io.MultiWriter(os.Stderr, l.logFile)
	} else {
		l.writers.info = l.logFile
		l.writers.debug = l.logFile
		l.writers.error = l.logFile
	}
	l.writers.export = os.Stdout

	l.infoLogger = log.New(l.writers.info, "INFO: ", log.Ldate|log.Ltime)
	l.debugLogger = log.New(l.writers.debug, "DEBUG: ", log.Ldate|log.Ltime)
	l.errorLogger = log.New(l.writers.error, "ERROR: ", log.Ldate|log.Ltime)
}

func (l *Logger) rotateLogFile() error {
	l.mu.RLock()
	fileInfo, err := l.logFile.Stat()
	l.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to get log file info: %w", err)
	}

	if fileInfo.Size() < maxLogFileSize {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	backupPath := fmt.Sprintf("%s.%s", l.logFilePath, time.Now().Format("20060102T150405"))
	if err := os.Rename(l.logFilePath, backupPath); err != nil {
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	l.logFile.Close()
	l.logFile, err = os.OpenFile(l.logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}

	l.updateWriters()

	return nil
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if err := l.rotateLogFile(); err != nil {
		l.errorLogger.Printf("Failed to rotate log file: %v", err)
	}
	l.infoLogger.Printf(format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.debug {
		l.mu.RLock()
		defer l.mu.RUnlock()
		if err := l.rotateLogFile(); err != nil {
			l.errorLogger.Printf("Failed to rotate log file: %v", err)
		}
		message := messageStyle.Render(fmt.Sprintf(format, v...))
		l.debugLogger.Print(message)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if err := l.rotateLogFile(); err != nil {
		l.errorLogger.Printf("Failed to rotate log file: %v", err)
	}
	l.errorLogger.Printf(format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if err := l.rotateLogFile(); err != nil {
		l.errorLogger.Printf("Failed to rotate log file: %v", err)
	}
	l.errorLogger.Fatalf(format, v...)
}

func (l *Logger) Export(format string, v ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	l.exportLogger.Print(fmt.Sprintf(format, v...))
}
