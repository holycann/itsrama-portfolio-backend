package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/lmittmann/tint"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogLevel mendefinisikan tingkat log yang tersedia
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// LoggerConfig berisi konfigurasi untuk logger
type LoggerConfig struct {
	// Konfigurasi file log
	Path       string
	MaxSize    int  // Ukuran maksimal file log dalam MB
	MaxBackups int  // Jumlah maksimal file log cadangan
	MaxAge     int  // Umur maksimal file log dalam hari
	Compress   bool // Apakah file log lama akan dikompresi

	// Konfigurasi logging
	Level       LogLevel // Tingkat log minimal yang akan dicatat
	Development bool     // Mode pengembangan untuk log yang lebih detail
}

// Logger adalah wrapper untuk slog dengan fitur tambahan
type Logger struct {
	logger     *slog.Logger
	config     LoggerConfig
	mu         sync.Mutex
	fileWriter *lumberjack.Logger
}

// NewLogger membuat instance Logger baru dengan konfigurasi yang diberikan
func NewLogger(cfg LoggerConfig) *Logger {
	// Validasi dan atur default konfigurasi
	cfg = validateConfig(cfg)

	// Siapkan file writer
	fileWriter := &lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	// Gabungkan file writer dan stdout
	multiWriter := io.MultiWriter(os.Stdout, fileWriter)

	// Buat handler dengan tint untuk warna dan format
	handler := createHandler(multiWriter, cfg)

	// Buat logger
	slogLogger := slog.New(handler)

	return &Logger{
		logger:     slogLogger,
		config:     cfg,
		fileWriter: fileWriter,
	}
}

// validateConfig memvalidasi dan mengatur default konfigurasi logger
func validateConfig(cfg LoggerConfig) LoggerConfig {
	// Atur default jika tidak ditentukan
	if cfg.Path == "" {
		cfg.Path = filepath.Join("logs", "app.log")
	}

	// Buat direktori log jika belum ada
	if err := os.MkdirAll(filepath.Dir(cfg.Path), 0755); err != nil {
		fmt.Printf("Gagal membuat direktori log: %v\n", err)
	}

	// Default konfigurasi
	if cfg.MaxSize == 0 {
		cfg.MaxSize = 100 // 100 MB
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 3
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 30 // 30 hari
	}

	return cfg
}

// createHandler membuat handler log dengan opsi yang disesuaikan
func createHandler(writer io.Writer, cfg LoggerConfig) slog.Handler {
	// Tentukan level log
	var level slog.Level
	switch cfg.Level {
	case DebugLevel:
		level = slog.LevelDebug
	case InfoLevel:
		level = slog.LevelInfo
	case WarnLevel:
		level = slog.LevelWarn
	case ErrorLevel:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Opsi handler
	opts := &tint.Options{
		Level:      level,
		TimeFormat: "2006-01-02 15:04:05",
	}

	// Tambahkan opsi development
	if cfg.Development {
		opts.AddSource = true
	}

	return tint.NewHandler(writer, opts)
}

// withCallerInfo menambahkan informasi pemanggil ke log
func (l *Logger) withCallerInfo(args ...any) []any {
	if l.config.Development {
		// Dapatkan informasi pemanggil
		_, file, line, ok := runtime.Caller(2)
		if ok {
			callerInfo := fmt.Sprintf("%s:%d", shortenPath(file), line)
			args = append([]any{slog.String("caller", callerInfo)}, args...)
		}
	}
	return args
}

// shortenPath mempersingkat path file
func shortenPath(path string) string {
	// Ambil 2 direktori terakhir
	parts := strings.Split(path, string(os.PathSeparator))
	if len(parts) > 2 {
		return filepath.Join(parts[len(parts)-2], parts[len(parts)-1])
	}
	return path
}

// Debug mencatat pesan debug
func (l *Logger) Debug(msg string, args ...any) {
	if l.config.Level <= DebugLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.logger.Debug(msg, l.withCallerInfo(args...)...)
	}
}

// Info mencatat pesan informasi
func (l *Logger) Info(msg string, args ...any) {
	if l.config.Level <= InfoLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.logger.Info(msg, l.withCallerInfo(args...)...)
	}
}

// Warn mencatat pesan peringatan
func (l *Logger) Warn(msg string, args ...any) {
	if l.config.Level <= WarnLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.logger.Warn(msg, l.withCallerInfo(args...)...)
	}
}

// Error mencatat pesan kesalahan
func (l *Logger) Error(msg string, args ...any) {
	if l.config.Level <= ErrorLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.logger.Error(msg, l.withCallerInfo(args...)...)
	}
}

// Fatal mencatat pesan fatal dan keluar dari program
func (l *Logger) Fatal(msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Error(msg, l.withCallerInfo(args...)...)
	os.Exit(1)
}

// Rotate melakukan rotasi file log manual
func (l *Logger) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.fileWriter.Rotate()
}

// Close menutup dan membersihkan sumber daya logger
func (l *Logger) Close() error {
	return l.fileWriter.Close()
}

// Contoh penggunaan konfigurasi default
func DefaultLogger() *Logger {
	return NewLogger(LoggerConfig{
		Path:        filepath.Join("logs", "app.log"),
		Level:       InfoLevel,
		Development: false,
	})
}

// DevelopmentLogger membuat logger untuk mode pengembangan
func DevelopmentLogger() *Logger {
	return NewLogger(LoggerConfig{
		Path:        filepath.Join("logs", "dev.log"),
		Level:       DebugLevel,
		Development: true,
	})
}
