// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package suite

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/services"
)

const (
	// DefaultLogLinesToFetch is the default number of log lines to fetch.
	DefaultLogLinesToFetch = 500

	// DefaultLogCollectionTimeout is the default timeout for collecting logs.
	DefaultLogCollectionTimeout = 10 * time.Second

	// LogDumpDir is the base output directory for dumped service logs.
	LogDumpDir = "e2e-logs"

	// LogDumpNumLines is the maximum number of log lines to fetch per service.
	LogDumpNumLines = 100_000

	// LogDumpTimeout is the timeout for fetching logs from a single service.
	LogDumpTimeout = 60 * time.Second
)

// GetServiceLogs fetches logs from a service by name.
// Returns a slice of log lines from the service.
func (s *KurtosisE2ESuite) GetServiceLogs(serviceName string) ([]string, error) {
	return s.GetServiceLogsWithOptions(serviceName, DefaultLogLinesToFetch, DefaultLogCollectionTimeout)
}

// GetServiceLogsWithOptions fetches logs with custom options.
func (s *KurtosisE2ESuite) GetServiceLogsWithOptions(
	serviceName string,
	numLines int,
	timeout time.Duration,
) ([]string, error) {
	// Get service context to obtain UUID
	sCtx, err := s.enclave.GetServiceContext(serviceName)
	if err != nil {
		return nil, err
	}
	serviceUUID := sCtx.GetServiceUUID()

	// Prepare service UUID map for log query
	serviceUUIDs := map[services.ServiceUUID]bool{
		serviceUUID: true,
	}

	// Get enclave identifier
	enclaveID := s.enclave.GetEnclaveUuid()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()

	// Fetch logs (not following, return all available, limited to numLines)
	//#nosec G115 // numLines is always positive and bounded by DefaultLogLinesToFetch
	logsChan, cancelLogs, err := s.kCtx.GetServiceLogs(
		ctx,
		string(enclaveID),
		serviceUUIDs,
		false, // shouldFollowLogs
		true,  // shouldReturnAllLogs
		uint32(numLines),
		nil, // no filter
	)
	if err != nil {
		return nil, err
	}
	defer cancelLogs()

	// Collect logs from channel using type inference
	var logs []string
	for {
		select {
		case content, ok := <-logsChan:
			if !ok {
				return logs, nil
			}
			// GetServiceLogsByServiceUuids is an exported method
			serviceLogs := content.GetServiceLogsByServiceUuids()
			if logLines, exists := serviceLogs[serviceUUID]; exists {
				for _, logLine := range logLines {
					logs = append(logs, logLine.GetContent())
				}
			}
		case <-ctx.Done():
			return logs, nil
		}
	}
}

// DumpAllServiceLogs fetches logs from all services in the enclave and writes
// them to disk. Output directory is e2e-logs/<SuiteName>/ with one file per
// service. Best-effort: logs warnings but does not fail teardown.
func (s *KurtosisE2ESuite) DumpAllServiceLogs() {
	svcMap, err := s.enclave.GetServices()
	if err != nil {
		s.logger.Warn("Failed to enumerate services for log dump", "error", err)
		return
	}

	// Build output directory: e2e-logs/<TestSuiteName>/
	outDir := filepath.Join(LogDumpDir, s.T().Name())

	// Remove and recreate so only the latest run's logs persist.
	if err = os.RemoveAll(outDir); err != nil {
		s.logger.Warn("Failed to remove old log dump directory", "dir", outDir, "error", err)
	}
	if err = os.MkdirAll(outDir, 0o750); err != nil {
		s.logger.Warn("Failed to create log dump directory", "dir", outDir, "error", err)
		return
	}

	for svcName := range svcMap {
		logs, fetchErr := s.GetServiceLogsWithOptions(
			string(svcName), LogDumpNumLines, LogDumpTimeout,
		)
		if fetchErr != nil {
			s.logger.Warn("Failed to fetch logs for service",
				"service", svcName, "error", fetchErr)
			continue
		}

		filePath := filepath.Join(outDir, string(svcName)+".log")
		content := strings.Join(logs, "\n")
		if writeErr := os.WriteFile(filePath, []byte(content), 0o600); writeErr != nil {
			s.logger.Warn("Failed to write log file",
				"service", svcName, "path", filePath, "error", writeErr)
		}
	}

	s.logger.Info("Service logs dumped", "dir", outDir, "services", len(svcMap))
}

// ContainsLogMessage checks if any log line contains the target message.
func ContainsLogMessage(logs []string, target string) bool {
	for _, log := range logs {
		if strings.Contains(log, target) {
			return true
		}
	}
	return false
}
