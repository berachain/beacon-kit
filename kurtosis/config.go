package kurtosis

import (
	"encoding/json"
)

// E2ETestConfig defines the configuration for end-to-end tests, including any additional services and participants involved.
type E2ETestConfig struct {
	AdditionalServices []interface{} `json:"additional_services"` // AdditionalServices specifies any extra services that should be included in the test environment.
	Participants       []Participant `json:"participants"`        // Participants lists the configurations for each participant in the test.
}

// Participant holds the configuration for a single participant in the test, including client images and types.
type Participant struct {
	// ClClientImage specifies the Docker image to use for the consensus layer client.
	ClClientImage string `json:"cl_client_image"`
	// ClClientType denotes the type of consensus layer client (e.g., beaconkit).
	ClClientType string `json:"cl_client_type"`
	// ElClientType denotes the type of execution layer client (e.g., reth).
	ElClientType string `json:"el_client_type"`
}

// DefaultE2ETestConfig provides a default configuration for end-to-end tests, pre-populating with a standard set of participants and no additional services.
func DefaultE2ETestConfig() *E2ETestConfig {
	return &E2ETestConfig{
		AdditionalServices: []interface{}{}, // Initialize with no additional services.
		Participants: []Participant{ // Pre-populate with a default set of participants.
			{
				ElClientType:  "reth",
				ClClientImage: "beacond:kurtosis-local",
				ClClientType:  "beaconkit",
			},
			{
				ElClientType:  "reth",
				ClClientImage: "beacond:kurtosis-local",
				ClClientType:  "beaconkit",
			},
			{
				ElClientType:  "reth",
				ClClientImage: "beacond:kurtosis-local",
				ClClientType:  "beaconkit",
			},
			{
				ElClientType:  "reth",
				ClClientImage: "beacond:kurtosis-local",
				ClClientType:  "beaconkit",
			},
		},
	}
}

func (e *E2ETestConfig) MustMarshalJSON() []byte {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	return jsonBytes
}
