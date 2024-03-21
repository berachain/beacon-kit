package beacon

type GetRandaoResponse struct {
	ExecutionOptimistic bool    `json:"execution_optimistic"`
	Finalized           bool    `json:"finalized"`
	Data                *Randao `json:"data"`
}

type Randao struct {
	Randao string `json:"randao"`
}
