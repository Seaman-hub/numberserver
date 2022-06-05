package api

import (
	"sync"
)

type IDPool struct {
	mtx      sync.Mutex
	released []uint
	max_id   uint
}

// NewIDPool creates and initializes an IDPool.
func NewIDPool(initValue uint) *IDPool {
	return &IDPool{
		max_id: initValue,
	}
}

// Fill recycled number hole
func (p *IDPool) Fillhole() {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	// TBD to fill release from storage
}

func (p *IDPool) Acquire() (id uint) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if len(p.released) > 0 {
		id = p.released[len(p.released)-1]
		p.released = p.released[:len(p.released)-1]
		return id
	}
	id = p.max_id
	p.max_id++
	return id
}

func (p *IDPool) Release(id uint) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.released = append(p.released, id)
}
