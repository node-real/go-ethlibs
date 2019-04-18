package node

import (
	"context"

	"github.com/INFURA/go-ethlibs/eth"
	"github.com/INFURA/go-ethlibs/jsonrpc"
)

type Requester interface {
	// Request method can be used to send JSONRPC requests and receive JSONRPC responses
	Request(ctx context.Context, r *jsonrpc.Request) (*jsonrpc.RawResponse, error)
}

type Subscriber interface {
	// Subscribe method can be used to subscribe via eth_subscribe
	Subscribe(ctx context.Context, r *jsonrpc.Request) (Subscription, error)
}

// Client represents a connection to an ethereum node
type Client interface {
	Requester
	Subscriber

	// URL returns the backend URL we are connected to
	URL() string

	// BlockNumber returns the current block number at head
	BlockNumber(ctx context.Context) (uint64, error)

	// BlockByNumber can be used to get a block by its number
	BlockByNumber(ctx context.Context, number uint64, full bool) (*eth.Block, error)

	// BlockByHash can be used to get a block by its hash
	BlockByHash(ctx context.Context, hash string, full bool) (*eth.Block, error)

	// eth_getTransactionByHash can be used to get transaction by its hash
	TransactionByHash(ctx context.Context, hash string) (*eth.Transaction, error)

	// NewHeads subscription
	NewHeads(ctx context.Context) (Subscription, error)

	// NewPendingTransactions subscriptions
	NewPendingTransaction(ctx context.Context) (Subscription, error)

	// TransactionReceipt for a particular transaction
	TransactionReceipt(ctx context.Context, hash string) (*eth.TransactionReceipt, error)

	// GetLogs
	GetLogs(ctx context.Context, filter eth.LogFilter) ([]eth.Log, error)

	SupportsSubscriptions() bool
}

type Subscription interface {
	Response() *jsonrpc.RawResponse
	ID() string
	Ch() chan *jsonrpc.Notification
	Unsubscribe(ctx context.Context) error
	Done() <-chan struct{}
	Err() error
}