// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package benchmarks_test

import (
	"bytes"
	"testing"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/log/pkg/phuslu"
	"github.com/cosmos/cosmos-sdk/server"
)

/* -------------------------------------------------------------------------- */
/*                                   Info                                     */
/* -------------------------------------------------------------------------- */

// Benchmark function for phuslu logger with pretty style.
func BenchmarkPhusluLoggerPrettyInfo(b *testing.B) {
	logger := newPhusluLogger().WithConfig(configWithPretty("info"))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info("This is an info message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for phuslu logger with JSON style.
func BenchmarkPhusluLoggerJSONInfo(b *testing.B) {
	logger := newPhusluLogger().WithConfig(configWithJSON("info"))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info("This is an info message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for SDK logger Info.
func BenchmarkSDKLoggerInfo(b *testing.B) {
	logger := newSDKLoggerWithLevel(b, "Info")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Info("This is an info message", "key1", "value1", "key2", 2)
	}
}

/* -------------------------------------------------------------------------- */
/*                                   Warn                                     */
/* -------------------------------------------------------------------------- */

// Benchmark function for phuslu logger Warn.
func BenchmarkPhusluLoggerPrettyWarn(b *testing.B) {
	logger := newPhusluLogger().WithConfig(configWithPretty("warn"))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Warn("This is a warning message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for phuslu logger with JSON style.
func BenchmarkPhusluLoggerJSONWarn(b *testing.B) {
	logger := newPhusluLogger().WithConfig(configWithJSON("warn"))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Warn("This is a warning message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for cosmos logger Warn.
func BenchmarkSDKLoggerWarn(b *testing.B) {
	logger := newSDKLoggerWithLevel(b, "Warn")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Warn("This is a warning message", "key1", "value1", "key2", 2)
	}
}

/* -------------------------------------------------------------------------- */
/*                                   Error                                    */
/* -------------------------------------------------------------------------- */

// Benchmark function for phuslu logger Error.
func BenchmarkPhusluLoggerPrettyError(b *testing.B) {
	logger := newPhusluLogger().WithConfig(configWithPretty("error"))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Error("This is an error message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for phuslu logger with JSON style.
func BenchmarkPhusluLoggerJSONError(b *testing.B) {
	logger := newPhusluLogger().WithConfig(configWithJSON("error"))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Error("This is an error message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for cosmos logger Error.
func BenchmarkSDKLoggerError(b *testing.B) {
	logger := newSDKLoggerWithLevel(b, "Error")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Error("This is an error message", "key1", "value1", "key2", 2)
	}
}

/* -------------------------------------------------------------------------- */
/*                                   Debug                                    */
/* -------------------------------------------------------------------------- */

// Benchmark function for phuslu logger Debug.
func BenchmarkPhusluLoggerPrettyDebug(b *testing.B) {
	logger := newPhusluLogger().WithConfig(configWithPretty("debug"))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Debug("This is a debug message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for phuslu logger with JSON style.
func BenchmarkPhusluLoggerJSONDebug(b *testing.B) {
	logger := newPhusluLogger().WithConfig(configWithJSON("debug"))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Debug("This is a debug message", "key1", "value1", "key2", 2)
	}
}

func BenchmarkPhusluLoggerPrettyDebugSilent(b *testing.B) {
	logger := newPhusluLogger().WithConfig(configWithPretty("info"))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Debug("This is a debug message", "key1", "value1", "key2", 2)
	}
}

func BenchmarkPhusluLoggerJSONDebugSilent(b *testing.B) {
	logger := newPhusluLogger().WithConfig(configWithJSON("info"))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Debug("This is a debug message", "key1", "value1", "key2", 2)
	}
}

// Benchmark function for cosmos logger Debug.
func BenchmarkSDKLoggerDebug(b *testing.B) {
	logger := newSDKLoggerWithLevel(b, "Debug")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Debug("This is a debug message", "key1", "value1", "key2", 2)
	}
}

func BenchmarkSDKLoggerDebugSilent(b *testing.B) {
	logger := newSDKLoggerWithLevel(b, "Info")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.Debug("This is a debug message", "key1", "value1", "key2", 2)
	}
}

/* -------------------------------------------------------------------------- */
/*                                   With                                     */
/* -------------------------------------------------------------------------- */

// Benchmark function for phuslu logger With.
func BenchmarkPhusluLoggerPrettyWith(b *testing.B) {
	logger := newPhusluLogger().WithConfig(configWithPretty("info"))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		newLogger := logger.With("contextKey", "contextValue")
		newLogger.Info("This is a contextual info message", "key1", "value1",
			"key2", 2)
	}
}

// Benchmark function for phuslu logger With JSON style.
func BenchmarkPhusluLoggerJSONWith(b *testing.B) {
	logger := newPhusluLogger().WithConfig(configWithJSON("info"))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		newLogger := logger.With("contextKey", "contextValue")
		newLogger.Info("This is a contextual info message", "key1", "value1",
			"key2", 2)
	}
}

// Benchmark function for cosmos logger With.
func BenchmarkSDKLoggerWith(b *testing.B) {
	logger := newSDKLoggerWithLevel(b, "Debug")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		newLogger := logger.With("contextKey", "contextValue")
		newLogger.Info("This is a contextual info message", "key1", "value1",
			"key2", 2)
	}
}

/* -------------------------------------------------------------------------- */
/*                                   Helpers                                  */
/* -------------------------------------------------------------------------- */

// setup func to create a new cosmos logger with the given log level.
func newSDKLoggerWithLevel(b *testing.B, level string) log.Logger {
	b.Helper()
	serverCtx := server.NewDefaultContext()
	serverCtx.Viper.Set("log_level", level)
	logger, err := server.CreateSDKLogger(serverCtx, &bytes.Buffer{})
	if err != nil {
		b.Fatalf("failed to create cosmos logger: %v", err)
	}
	return logger
}

// setup func to create a new phuslu logger with the given log level.
func newPhusluLogger() *phuslu.Logger[log.Logger] {
	cfg := phuslu.DefaultConfig() // dummy config
	l := phuslu.NewLogger[log.Logger](
		&bytes.Buffer{}, &cfg)
	return l
}

// setup func to create a phuslu logger config with pretty style.
func configWithPretty(level string) phuslu.Config {
	cfg := phuslu.DefaultConfig()
	cfg.LogLevel = level
	cfg.Style = phuslu.StylePretty
	return cfg
}

// setup func to create a phuslu logger config with JSON style.
func configWithJSON(level string) phuslu.Config {
	cfg := phuslu.DefaultConfig()
	cfg.LogLevel = level
	cfg.Style = phuslu.StyleJSON
	return cfg
}
