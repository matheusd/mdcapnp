package capnprpc

import (
	"fmt"
	"sync"

	"matheusd.com/mdcapnp/capnpser"
	"matheusd.com/mdcapnp/internal/pools"
)

type inMsgPool struct {
	slices *pools.SlicesPool
	msgs   *sync.Pool
}

func (ap *inMsgPool) new() *InMsg {
	m := ap.msgs.Get().(*InMsg)
	m.buf = ap.slices.Get(0)
	return m
}

func (ap *inMsgPool) put(m *InMsg) {
	if m.pool != ap {
		// *Really* shouldn't happen.
		panic(fmt.Sprintf("trying to return InMsg to pool %p when "+
			"InMsg was from pool %p", ap, m.pool))
	}
	clear(m.buf)
	arena := m.Msg.Arena()
	arena.ReadLimiter().Init(0)
	arena.Reset(nil)
	ap.slices.Put(m.buf[:cap(m.buf)])
	m.buf = nil
	ap.msgs.Put(m)
}

func (ap *inMsgPool) replaceBuf(old []byte, wantWordLen capnpser.WordCount) []byte {
	wantLen := int(wantWordLen.ByteCount())
	new := ap.slices.Get(wantLen)
	copy(new, old)
	ap.slices.Put(old)
	return new[:wantLen]
}

func newArenaPool(slicePoolSizes ...int) *inMsgPool {
	p := &inMsgPool{
		slices: pools.NewSlicesPool(slicePoolSizes...),
		msgs:   &sync.Pool{},
	}

	p.msgs.New = func() any {
		return &InMsg{
			Msg:  capnpser.MakeMsg(&capnpser.Arena{}),
			pool: p,
		}
	}
	return p

}
