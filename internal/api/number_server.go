package api

import (
	"github.com/Seaman-hub/numberserver/api/ns"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
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
	log.WithFields(log.Fields{
		"mid": mid,
	}).Info("allocted")
	return &ns.GetSequenceNumResponse{
		Number: mid,
	}, nil
}

// PutSequenceNum puts a sequence number back to pool.
func (n *NumberServerAPI) PutSequenceNum(ctx context.Context, req *ns.PutSequenceNumRequest) (*ns.PutSequenceNumResponse, error) {
	numberPool.Release((uint)(req.Number))
	log.WithFields(log.Fields{
		"number": req.Number,
	}).Info("released")
	return &ns.PutSequenceNumResponse{}, nil
}
