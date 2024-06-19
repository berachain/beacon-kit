package benchmarks

import (
	"bytes"
	"testing"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
	"github.com/cosmos/cosmos-sdk/server"
)

// Benchmark function for phuslu logger.
func BenchmarkPhusluLoggerInfo(b *testing.B) {
	logger := phuslu.NewLogger[log.Logger]("Info", &bytes.Buffer{})
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info("This is an info message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for SDK logger Info.
func BenchmarkSDKLoggerInfo(b *testing.B) {
	serverCtx := server.NewDefaultContext()
	serverCtx.Viper.Set("log_level", "Info")
	logger, err := server.CreateSDKLogger(serverCtx, &bytes.Buffer{})
	if err != nil {
		b.Fatalf("failed to create cosmos logger: %v", err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info("This is an info message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for phuslu logger Warn.
func BenchmarkPhusluLoggerWarn(b *testing.B) {
	logger := phuslu.NewLogger[log.Logger]("Debug", &bytes.Buffer{})
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Warn("This is a warning message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for cosmos logger Warn.
func BenchmarkSDKLoggerWarn(b *testing.B) {
	serverCtx := server.NewDefaultContext()
	serverCtx.Viper.Set("log_level", "Debug")
	logger, err := server.CreateSDKLogger(serverCtx, &bytes.Buffer{})
	if err != nil {
		b.Fatalf("failed to create cosmos logger: %v", err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Warn("This is a warning message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for phuslu logger Error.
func BenchmarkPhusluLoggerError(b *testing.B) {
	logger := phuslu.NewLogger[log.Logger]("Error", &bytes.Buffer{})
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Error("This is an error message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for cosmos logger Error.
func BenchmarkSDKLoggerError(b *testing.B) {
	serverCtx := server.NewDefaultContext()
	serverCtx.Viper.Set("log_level", "Debug")
	logger, err := server.CreateSDKLogger(server.NewDefaultContext(),
		&bytes.Buffer{})
	if err != nil {
		b.Fatalf("failed to create cosmos logger: %v", err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Error("This is an error message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for phuslu logger Debug.
func BenchmarkPhusluLoggerDebug(b *testing.B) {
	logger := phuslu.NewLogger[log.Logger]("Debug", &bytes.Buffer{})
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Debug("This is a debug message", "key1", "value1", "key2", 2)
	}
}

func BenchmarkPhusluLoggerDebugSilent(b *testing.B) {
	logger := phuslu.NewLogger[log.Logger]("Info", &bytes.Buffer{})
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Debug("This is a debug message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for cosmos logger Debug.
func BenchmarkSDKLoggerDebug(b *testing.B) {
	serverCtx := server.NewDefaultContext()
	serverCtx.Viper.Set("log_level", "Debug")
	logger, err := server.CreateSDKLogger(serverCtx, &bytes.Buffer{})
	if err != nil {
		b.Fatalf("failed to create cosmos logger: %v", err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Debug("This is a debug message", "key1", "value1", "key2", 2)
	}
}

func BenchmarkSDKLoggerDebugSilent(b *testing.B) {
	serverCtx := server.NewDefaultContext()
	serverCtx.Viper.Set("log_level", "Info")
	logger, err := server.CreateSDKLogger(serverCtx, &bytes.Buffer{})
	if err != nil {
		b.Fatalf("failed to create cosmos logger: %v", err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Debug("This is a debug message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for phuslu logger With.
func BenchmarkPhusluLoggerWith(b *testing.B) {
	logger := phuslu.NewLogger[log.Logger]("Debug", &bytes.Buffer{})
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		newLogger := logger.With("contextKey", "contextValue")
		newLogger.Info("This is a contextual info message", "key1", "value1",
			"key2", 2)
	}
}

// Benchmark function for cosmos logger With.
func BenchmarkSDKLoggerWith(b *testing.B) {
	serverCtx := server.NewDefaultContext()
	serverCtx.Viper.Set("log_level", "Debug")
	logger, err := server.CreateSDKLogger(serverCtx, &bytes.Buffer{})
	if err != nil {
		b.Fatalf("failed to create cosmos logger: %v", err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		newLogger := logger.With("contextKey", "contextValue")
		newLogger.Info("This is a contextual info message", "key1", "value1",
			"key2", 2)
	}
}
