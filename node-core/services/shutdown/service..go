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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package shutdown

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
)

// Service is a service that manages any startup or shutdown tasks that need to be done.
type Service struct {
	// logger is used for logging messages in the service.
	logger log.Logger

	// pidFile is the path to the pid file we use to detect unsafe shutdowns.
	pidFile string
}

func NewService(logger log.Logger, pidFile string) *Service {
	return &Service{logger: logger, pidFile: pidFile}
}

// Name returns the name of the service.
func (*Service) Name() string {
	return "shutdown"
}

func (s *Service) Start(_ context.Context) error {
	// if the pid file already exists it means we didn't gracefully shutdown.
	if _, err := os.Stat(s.pidFile); err == nil {
		s.printUnsafeShutDownDetected()
	}

	// create a new pid file and write our process id to it. If it already existed before
	// it will be truncated and overwritten with our new pid. The timestamp of the file
	// will be updated to the current time so we can compute the time it was running
	f, err := os.Create(s.pidFile)
	if err != nil {
		return errors.Wrap(err, "failed to create pidfile")
	}
	defer f.Close()
	_, err = f.WriteString(strconv.Itoa(os.Getpid()))
	if err != nil {
		return errors.Wrap(err, "failed to write pid to pidfile")
	}

	return nil
}

func (s *Service) Stop() error {
	s.logger.Info("Stopping shutdown service")

	err := os.Remove(s.pidFile)
	if err != nil {
		return errors.Wrap(err, "failed to remove pidfile")
	}

	return nil
}

func (s *Service) printUnsafeShutDownDetected() {
	// collect info about the previous pid instance and when it was started
	pid := "unknown"
	if bytes, readErr := os.ReadFile(s.pidFile); readErr == nil {
		pid = string(bytes)
	}
	started := "unknown"
	if fi, statErr := os.Stat(s.pidFile); statErr == nil {
		started = fi.ModTime().String()
	}

	s.logger.Warn(fmt.Sprintf(`

	+==========================================================================+
	+ ⚠️ Detected an unsafe shutdown!
	+ 🧩 Previously running:
	+      PID: %s
	+      Started: %s
	+==========================================================================+

`,
		pid, started))
}
