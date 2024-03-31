package provider

import "context"

const (
	finalized_key = "fc_finalized"
)

func (p *Provider) GetTrustedEth1Hash() []byte {
	hash, err := p.QueryWithProof(context.Background(), finalized_key, p.latestBlockHeight)
	if err != nil {
		panic(err)
	}
	return hash.Bytes()
}
