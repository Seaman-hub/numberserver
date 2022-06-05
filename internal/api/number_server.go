package api

import (
	"golang.org/x/net/context"

	"github.com/Seaman-hub/numberserver/api/ns"
)

var (
	numberPool *IDPool
)

// NumberServerAPI defines the number-server API.
type NumberServerAPI struct {
}

// NewNumberServerAPI returns a new NumberServerAPI.
func NewNumberServerAPI(initValue uint) *NumberServerAPI {
	numberPool = NewIDPool(initValue)
	numberPool.Fillhole()
	return &NumberServerAPI{}
}

// GetSequenceNum returns a sequence number.
func (n *NumberServerAPI) GetSequenceNum(ctx context.Context, req *ns.GetSequenceNumRequest) (*ns.GetSequenceNumResponse, error) {
	mid := (uint32)(numberPool.Acquire())
	return &ns.GetSequenceNumResponse{
		Number: mid,
	}, nil
}

// PutSequenceNum puts a sequence number back to pool.
func (n *NumberServerAPI) PutSequenceNum(ctx context.Context, req *ns.PutSequenceNumRequest) (*ns.PutSequenceNumResponse, error) {
	numberPool.Release((uint)(req.Number))
	return &ns.PutSequenceNumResponse{}, nil
}
