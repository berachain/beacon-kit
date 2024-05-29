package version

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/log"
)

// defaultReportingInterval is the default interval at which the version is reported.
const defaultReportingInterval = 30 * time.Second

// ReportingService is a service that periodically logs the running chain version.
type ReportingService struct {
	logger  log.Logger[any]
	version string
	ticker  *time.Ticker
}

// NewReportingService creates a new VersionReporterService.
func NewReportingService(
	logger log.Logger[any],
	version string,
) *ReportingService {
	return &ReportingService{
		logger:  logger,
		version: version,
		ticker: time.NewTicker(
			defaultReportingInterval),
	}
}

// Name returns the name of the service.
func (v *ReportingService) Name() string {
	return "ReportingService"
}

// Start begins the periodic logging of the chain version.
func (v *ReportingService) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				v.ticker.Stop()
				return
			case <-v.ticker.C:
				v.logger.Info("this node is running", "version", v.version)
			}
		}
	}()
	return nil
}

// Status returns nil if the service is healthy.
func (v *ReportingService) Status() error {
	return nil
}

// WaitForHealthy waits for all registered services to be healthy.
func (v *ReportingService) WaitForHealthy(ctx context.Context) {}
