package pools

import (
	"github.com/colega/zeropool"
)

type sizedSlicePool struct {
	size int
	pool *zeropool.Pool[[]byte]
}

type SlicesPool struct {
	pools []sizedSlicePool
}

func (sp *SlicesPool) poolForSize(minSize int) *sizedSlicePool {
	// Find the smallest bucket that has can provide a message of at least
	// this size.
	for i := range sp.pools {
		if sp.pools[i].size >= minSize {
			return &sp.pools[i]
		}
	}
	return nil
}

func (sp *SlicesPool) Get(minCap int) []byte {
	p := sp.poolForSize(minCap)
	if p == nil {
		return make([]byte, minCap)
	}

	res := p.pool.Get()
	if res == nil {
		res = make([]byte, p.size, p.size)
	}
	return res
}

func (sp *SlicesPool) Put(b []byte) {
	clear(b)
	p := sp.poolForSize(cap(b))
	if p != nil && p.size == cap(b) { // Only store when cap matches exactly.
		p.pool.Put(b[:cap(b)])
	}
}

var defaultSlicesPoolSizes = []int{16 * 8, 128 * 8, 512 * 8, 4096 * 8, 32768 * 8, 65536 * 8}

func NewSlicesPool(sizes ...int) *SlicesPool {
	if sizes == nil {
		sizes = defaultSlicesPoolSizes
	}

	res := &SlicesPool{
		pools: make([]sizedSlicePool, len(sizes)),
	}
	for i, size := range sizes {
		res.pools[i].size = size
		res.pools[i].pool = &zeropool.Pool[[]byte]{}
	}
	return res
}
