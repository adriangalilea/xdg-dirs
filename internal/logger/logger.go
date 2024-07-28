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
)

const maxLogFileSize = 10 * 1024 * 1024 // 10 MB

type Logger struct {
	debug        bool
	logFile      *os.File
	logFilePath  string
	infoLogger   *log.Logger
	debugLogger  *log.Logger
	errorLogger  *log.Logger
	exportLogger *log.Logger
	mu           sync.Mutex
}

func NewLogger(debug bool, logFilePath string) *Logger {
	var logFile *os.File
	var err error

	if logFilePath == "" {
		logFilePath = filepath.Join(os.TempDir(), "xdg-user-dirs-update.log")
	}

	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	var infoWriter, debugWriter, errorWriter io.Writer
	if debug {
		infoWriter = io.MultiWriter(os.Stdout, logFile)
		debugWriter = io.MultiWriter(os.Stdout, logFile)
		errorWriter = io.MultiWriter(os.Stderr, logFile)
	} else {
		infoWriter = logFile
		debugWriter = logFile
		errorWriter = logFile
	}

	logger := &Logger{
		debug:        debug,
		logFile:      logFile,
		logFilePath:  logFilePath,
		infoLogger:   log.New(infoWriter, "INFO: ", log.Ldate|log.Ltime),
		debugLogger:  log.New(debugWriter, "DEBUG: ", log.Ldate|log.Ltime),
		errorLogger:  log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime),
		exportLogger: log.New(os.Stdout, "", 0),
	}

	return logger
}

func (l *Logger) rotateLogFile() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	fileInfo, err := l.logFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get log file info: %w", err)
	}

	if fileInfo.Size() < maxLogFileSize {
		return nil
	}

	backupPath := fmt.Sprintf("%s.%s", l.logFilePath, time.Now().Format("20060102T150405"))
	if err := os.Rename(l.logFilePath, backupPath); err != nil {
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	l.logFile.Close()
	l.logFile, err = os.OpenFile(l.logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}

	l.infoLogger.SetOutput(io.MultiWriter(os.Stdout, l.logFile))
	l.debugLogger.SetOutput(io.MultiWriter(os.Stdout, l.logFile))
	l.errorLogger.SetOutput(io.MultiWriter(os.Stderr, l.logFile))

	return nil
}

func (l *Logger) Info(format string, v ...interface{}) {
	if err := l.rotateLogFile(); err != nil {
		l.mu.Lock()
		l.errorLogger.Printf("Failed to rotate log file: %v", err)
		l.mu.Unlock()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.infoLogger.Printf(format, v...)
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.debug {
		if err := l.rotateLogFile(); err != nil {
			l.mu.Lock()
			l.errorLogger.Printf("Failed to rotate log file: %v", err)
			l.mu.Unlock()
		}
		l.mu.Lock()
		defer l.mu.Unlock()
		l.debugLogger.Printf(format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if err := l.rotateLogFile(); err != nil {
		l.mu.Lock()
		l.errorLogger.Printf("Failed to rotate log file: %v", err)
		l.mu.Unlock()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.errorLogger.Printf(format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	if err := l.rotateLogFile(); err != nil {
		l.mu.Lock()
		l.errorLogger.Printf("Failed to rotate log file: %v", err)
		l.mu.Unlock()
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.errorLogger.Fatalf(format, v...)
}

func (l *Logger) Export(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.exportLogger.Print(fmt.Sprintf(format, v...))
}
