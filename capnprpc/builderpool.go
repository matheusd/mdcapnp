// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"sync"

	"github.com/colega/zeropool"
	types "matheusd.com/mdcapnp/capnprpc/types"
	"matheusd.com/mdcapnp/capnpser"
)

const (
	callMessageSizeOverhead    capnpser.WordCount = 6 // TODO: Measure this precisely.
	returnMessageSizeOverhead  capnpser.WordCount = 6 // TODO: Measure this precisely.
	singleCapReturnPayloadSize capnpser.WordCount = 3 // TODO: Measure this precisely.
)

type msgBuilderPool struct {
	p         *sync.Pool
	slicePool *zeropool.Pool[[]byte]
	alloc     capnpser.Allocator
	size      capnpser.WordCount
}

func (mbp *msgBuilderPool) put(v *capnpser.MessageBuilder) {
	mbp.p.Put(v)
}

func (mbp *msgBuilderPool) get() *capnpser.MessageBuilder {
	return mbp.p.Get().(*capnpser.MessageBuilder)
}

func (mbp *msgBuilderPool) newMb() any {
	serMb, err := capnpser.NewMessageBuilder(mbp.alloc)
	if err != nil {
		// SimpleSingleAlloc never errors here.
		panic(err)
	}
	return serMb
}

type MessageBuilderPool struct {
	// alloc capnpser.Allocator // Alloc strategy for outbound rpc messages.
	pools []*msgBuilderPool
	// p     *sync.Pool
}

func (mp *MessageBuilderPool) poolForSize(size capnpser.WordCount) *msgBuilderPool {
	// Find the smallest bucket that has can provide a message of at least
	// this size.
	for _, p := range mp.pools {
		if p.size >= size {
			return p
		}
	}
	return nil
}

func (mp *MessageBuilderPool) mustPoolForSize(size capnpser.WordCount) *msgBuilderPool {
	// No bucket found for this size. Go with the last bucket.
	res := mp.poolForSize(2 + size) // TODO: is this the right place to add the overhad?
	if res == nil {
		res = mp.pools[len(mp.pools)-1]
	}
	return res
}

func (mp *MessageBuilderPool) getRawMessageBuilder(sizeHint capnpser.WordCount) *capnpser.MessageBuilder {
	return mp.mustPoolForSize(sizeHint).get()
}

func (mp *MessageBuilderPool) getForPayloadSize(extraPayloadSize capnpser.WordCount) (outMsg, error) {
	// TODO: calculate the size hint.
	serMb := mp.getRawMessageBuilder(extraPayloadSize)
	mb, err := types.NewRootMessageBuilder(serMb)
	if err != nil {
		return outMsg{}, err
	}
	return outMsg{serMsg: serMb, mb: mb}, nil
}

func (mp *MessageBuilderPool) get() (outMsg, error) {
	return mp.getForPayloadSize(0)
}

func (mp *MessageBuilderPool) put(serMb *capnpser.MessageBuilder) {
	err := serMb.Reset()
	if err != nil {
		panic(err) // Simple allocator never errors on Reset().
	}
	msgCap := serMb.TotalCapacity()
	if msgCap == 0 {
		// Should not happen for msgs returned by mp.
		return
	}
	msgCapWords, ok := msgCap.StorageWordCount()
	if !ok {
		// Should not happen ever.
		panic("should never happen")
	}

	// Find the pool from where this was allocated based on the capacity.
	//
	// TODO: maybe consider if the allocator is the same too? This is all
	// bad.
	for _, p := range mp.pools {
		if p.size == msgCapWords {
			p.put(serMb)
			return
		}
	}
}

type messageBuilderPoolSerPoolAdapter struct {
	mbp *MessageBuilderPool
}

func (a messageBuilderPoolSerPoolAdapter) Get(size capnpser.WordCount) []byte {
	pool := a.mbp.poolForSize(size)
	if pool == nil {
		return make([]byte, size.ByteCount(), size.ByteCount()*2)
	}

	res := pool.slicePool.Get()
	return res[:size.ByteCount()]
}

func (a messageBuilderPoolSerPoolAdapter) Put(b []byte) {
	wc := capnpser.WordCount(cap(b) / capnpser.WordSize)
	pool := a.mbp.poolForSize(wc)
	if pool == nil {
		return
	}

	// Ensure wrongly sized buffers aren't put in the pool.
	if pool.size != wc {
		return
	}
	pool.slicePool.Put(b)
}

func NewMessageBuilderPool() *MessageBuilderPool {
	sizes := []capnpser.WordCount{16, 128, 512, 4096, 32768, 65536 /*, 262144*/}

	pools := make([]*msgBuilderPool, len(sizes))
	res := &MessageBuilderPool{pools: pools}
	poolAdapter := messageBuilderPoolSerPoolAdapter{mbp: res}
	for i, size := range sizes {
		alloc := capnpser.NewSingleSegmentPoolableAllocator(size, poolAdapter)
		mbp := &msgBuilderPool{alloc: alloc, size: size}
		mbp.p = &sync.Pool{New: mbp.newMb}
		sp := zeropool.New(func() []byte {
			return make([]byte, size.ByteCount())
		})
		mbp.slicePool = &sp
		pools[i] = mbp

		// Make this bucket hot??
		/*
			for range 0 {
				sp.Put(sp.Get())
			}
		*/
	}

	return res
}
