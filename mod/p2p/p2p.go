package p2p

import "context"

type PublisherReceiver[InPubT, OutPubT, InReceiverT, OutReceiverT any] interface {
	Publisher[InPubT, OutPubT]
	Receiver[InReceiverT, OutReceiverT]
}

type Publisher[InT, OutT any] interface {
	Publish(ctx context.Context, data InT) (OutT, error)
}

type Receiver[InT, OutT any] interface {
	Request(ctx context.Context, ref InT, out OutT) error
}
