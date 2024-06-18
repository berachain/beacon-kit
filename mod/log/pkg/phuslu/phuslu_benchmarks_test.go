package phuslu_test

import (
	"io"
	"testing"

	"github.com/berachain/beacon-kit/mod/cli/pkg/builder"
	"github.com/cosmos/cosmos-sdk/server"
)

// Benchmark function for phuslu logger
func BenchmarkPhusluLoggerInfo(b *testing.B) {
	logger, _ := builder.CreatePhusluLogger(server.NewDefaultContext(), io.Discard)
	for n := 0; n < b.N; n++ {
		logger.Info("This is an info message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for SDK logger Info
func BenchmarkSDKLoggerInfo(b *testing.B) {
	logger, _ := server.CreateSDKLogger(server.NewDefaultContext(), io.Discard)
	for n := 0; n < b.N; n++ {
		logger.Info("This is an info message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for phuslu logger Warn
func BenchmarkPhusluLoggerWarn(b *testing.B) {
	serverCtx := server.NewDefaultContext()
	logger, err := builder.CreatePhusluLogger(serverCtx, io.Discard)
	if err != nil {
		b.Fatalf("failed to create phuslu logger: %v", err)
	}
	for n := 0; n < b.N; n++ {
		logger.Warn("This is a warning message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for cosmos logger Warn
func BenchmarkSDKLoggerWarn(b *testing.B) {
	serverCtx := server.NewDefaultContext()
	logger, err := server.CreateSDKLogger(serverCtx, io.Discard)
	if err != nil {
		b.Fatalf("failed to create cosmos logger: %v", err)
	}
	for n := 0; n < b.N; n++ {
		logger.Warn("This is a warning message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for phuslu logger Error
func BenchmarkPhusluLoggerError(b *testing.B) {
	serverCtx := server.NewDefaultContext()
	logger, err := builder.CreatePhusluLogger(serverCtx, io.Discard)
	if err != nil {
		b.Fatalf("failed to create phuslu logger: %v", err)
	}
	for n := 0; n < b.N; n++ {
		logger.Error("This is an error message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for cosmos logger Error
func BenchmarkSDKLoggerError(b *testing.B) {
	logger, err := server.CreateSDKLogger(server.NewDefaultContext(), io.Discard)
	if err != nil {
		b.Fatalf("failed to create cosmos logger: %v", err)
	}
	for n := 0; n < b.N; n++ {
		logger.Error("This is an error message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for phuslu logger Debug
func BenchmarkPhusluLoggerDebug(b *testing.B) {
	serverCtx := server.NewDefaultContext()
	logger, err := builder.CreatePhusluLogger(serverCtx, io.Discard)
	if err != nil {
		b.Fatalf("failed to create phuslu logger: %v", err)
	}
	for n := 0; n < b.N; n++ {
		logger.Debug("This is a debug message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for cosmos logger Debug
func BenchmarkSDKLoggerDebug(b *testing.B) {
	serverCtx := server.NewDefaultContext()
	logger, err := server.CreateSDKLogger(serverCtx, io.Discard)
	if err != nil {
		b.Fatalf("failed to create cosmos logger: %v", err)
	}
	for n := 0; n < b.N; n++ {
		logger.Debug("This is a debug message", "key1", "value1", "key2", 2)
	}
}

// // Benchmark function for phuslu logger With
// func BenchmarkPhusluLoggerWith(b *testing.B) {
// 	serverCtx := server.NewDefaultContext()
// 	logger, err := builder.CreatePhusluLogger(serverCtx, io.Discard)
// 	if err != nil {
// 		b.Fatalf("failed to create phuslu logger: %v", err)
// 	}
// 	for n := 0; n < b.N; n++ {
// 		newLogger := logger.With("contextKey", "contextValue")
// 		newLogger.Info("This is a contextual info message", "anotherKey", "anotherValue")
// 	}
// }

// // Benchmark function for cosmos logger With
// func BenchmarkSDKLoggerWith(b *testing.B) {
// 	serverCtx := server.NewDefaultContext()
// 	logger, err := server.CreateSDKLogger(serverCtx, io.Discard)
// 	if err != nil {
// 		b.Fatalf("failed to create cosmos logger: %v", err)
// 	}
// 	for n := 0; n < b.N; n++ {
// 		newLogger := logger.With("contextKey", "contextValue")
// 		newLogger.Info("This is a contextual info message", "anotherKey", "anotherValue")
// 	}
// }

func main() {
}
