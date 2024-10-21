package types

import (
	"context"

	"github.com/berachain/beacon-kit/testing/e2e/suite/types/proof"
	"github.com/berachain/offchain-sdk/log"
	"github.com/hashicorp/go-retryablehttp"
)

// Client is a client for the Beacon node API.
type Client struct {
	// Proof API Client
	*proof.Client

	// Config for Beacon API HTTP calls.
	cfg *Config

	// HTTP client that handles retries with a default retry policy.
	httpClient *retryablehttp.Client

	// The logger to handle logs
	logger log.Logger
}

// NewClient creates a client for the Beacon node API.
func NewClient(cfg *Config, logger log.Logger) (*Client, error) {
	// TODO: setup and start the beacon node.

	// Ensure the given configuration is valid.
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	// Setup and configure the retryable HTTP client.
	httpClient := retryablehttp.NewClient()
	httpClient.HTTPClient.Timeout = cfg.HttpTimeout
	httpClient.Logger = logger
	httpClient.RetryMax = cfg.MaxRetries

	// Setup the proof client.
	proofs := proof.NewClient(httpClient, cfg.ApiURL)

	return &Client{
		Client:     proofs,
		cfg:        cfg,
		httpClient: httpClient,
		logger:     logger,
	}, nil
}

// Shutdown gracefully shuts down the Beacon client.
func (c *Client) Shutdown(context.Context) error {
	// TODO: shutdown the beacon node.

	c.httpClient.HTTPClient.CloseIdleConnections()

	return nil
}
