package types

type SyncingData struct {
	HeadSlot     int64 `json:"head_slot"`
	SyncDistance int64 `json:"sync_distance"`
	IsSyncing    bool  `json:"is_syncing"`
	IsOptimistic bool  `json:"is_optimistic"`
	ELOffline    bool  `json:"el_offline"`
}

type GenericResponse struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                any  `json:"data"`
}

// NewResponse creates a new response with CometBFT's finality guarantees.
func NewResponse(data any) GenericResponse {
	return GenericResponse{
		// All data is finalized in CometBFT since we only return data for slots up to head
		Finalized: true,
		// Never optimistic since we only return finalized data
		ExecutionOptimistic: false,
		Data:                data,
	}
}
