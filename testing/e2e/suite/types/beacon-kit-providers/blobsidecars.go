package beacon_kit_providers

import (
	"context"
	"fmt"
	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/deneb"
)

func (s *http.Service) BlobSidecars(ctx context.Context,
	opts *api.BlobSidecarsOpts,
) (
	*api.Response[[]*deneb.BlobSidecar],
	error,
) {
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}
	if opts == nil {
		return nil, client.ErrNoOptions
	}
	if opts.Block == "" {
		return nil, errors.Join(errors.New("no block specified"), client.ErrInvalidOptions)
	}

	endpoint := fmt.Sprintf("/eth/v1/beacon/blob_sidecars/%s", opts.Block)
	httpResponse, err := s.get(ctx, endpoint, "", &opts.Common, true)
	if err != nil {
		return nil, err
	}

	var response *api.Response[[]*deneb.BlobSidecar]
	switch httpResponse.contentType {
	case ContentTypeSSZ:
		response, err = s.blobSidecarsFromSSZ(httpResponse)
	case ContentTypeJSON:
		response, err = s.blobSidecarsFromJSON(httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}
	if err != nil {
		return nil, err
	}

	return response, nil
}
