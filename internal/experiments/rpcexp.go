// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package experiments

import (
	"context"
	"errors"
	"slices"

	"github.com/sourcegraph/conc/pool"
)

type msgBuilder struct{} // Alias to a serializer MessageBuilder

type message struct {
	// serialized msg?
}

type callable struct {
	// promise || local-callable || remote-capability
	// pipelinable
}

type callParamsBuilder func(*msgBuilder) error

type pipelineStep struct {
	conn          *runningConn
	interfaceId   uint64
	methodId      uint16
	argsBuilder   func(*msgBuilder) error // Builds an rpc.Call struct
	paramsBuilder callParamsBuilder       // Builds the Params field of an rpc.Call struct
}

type pipeline struct {
	steps []pipelineStep
}

//go:noinline
func (pipe *pipeline) addStep(iid uint64, mid uint16, pb callParamsBuilder) int {
	pipe.steps = append(pipe.steps, pipelineStep{
		interfaceId:   iid,
		methodId:      mid,
		paramsBuilder: pb,
	})
	return len(pipe.steps) - 1
}

type capability interface{}

type futureCap = struct {
	pipe      *pipeline
	stepIndex int
}

func remoteCall(obj futureCap, iid uint64, mid uint16, pb callParamsBuilder) futureCap {
	obj.pipe.steps = append(obj.pipe.steps, pipelineStep{
		interfaceId:   iid,
		methodId:      mid,
		paramsBuilder: pb,
	})
	return futureCap{obj.pipe, len(obj.pipe.steps) - 1}
}

func waitResult[T any](ctx context.Context, cap futureCap) (T, error) {
	// Run cap.pipe
	panic("boo")
}

type conn struct {
	in  chan message
	out chan message
}

type runningConn struct {
	*conn
	ctx    context.Context
	cancel func() // Closes runningConn.
}

func (rc *runningConn) queueOut(m message) {
	select {
	// TODO: what if out is blocked? Impose max limit before
	// effecting disconnection?
	case rc.conn.out <- m:
	case <-rc.ctx.Done():
	}
}

type bootstrapCap struct {
	_bootstrapCap struct{}
	fc            futureCap
}

func (rc *runningConn) bootstrap() bootstrapCap {
	panic("???")
}

type inMsg struct {
	conn *conn
	msg  message
}

type outMsg struct {
	conn *conn
	msg  message
}

type vat struct {
	conns []*conn

	newConn  chan *conn
	connDone chan *conn

	inMsg  chan inMsg
	outMsg chan outMsg
}

func (v *vat) RunConn(c *conn) *runningConn {
	panic("boo")
}

func (v *vat) runConn(g *pool.ContextPool, ctx context.Context, c *conn) *runningConn {
	connCtx, connCancel := context.WithCancel(ctx)
	rc := &runningConn{
		conn:   c,
		ctx:    connCtx,
		cancel: connCancel,
	}

	connG := pool.New().WithContext(connCtx).WithCancelOnError().WithFirstError()
	connG.Go(func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case msg := <-c.in:
				v.inMsg <- inMsg{conn: c, msg: msg}
			}
		}
	})

	g.Go(func(ctx context.Context) error {
		err := connG.Wait()
		// Consider canceled a graceful conn closure.
		if errors.Is(err, context.Canceled) {
			err = nil
		}
		return err
	})

	return rc
}

// vatRunState is the running state of the vat.
//
// This MUST NOT be accessed directly or stored outside of Run() and runStep(),
// under penalty of race conditions.
//
// Individual fields _may_ be accessed on other goroutines, as long as they have
// been properly captured by a closure.
type vatRunState struct {
	ctx   context.Context
	g     *pool.ContextPool
	conns []*runningConn
}

func (s *vatRunState) delConn(c *conn) {
	s.conns = slices.DeleteFunc(s.conns, func(rc *runningConn) bool {
		return rc.conn == c
	})
}

func (s *vatRunState) findConn(c *conn) *runningConn {
	i := slices.IndexFunc(s.conns, func(rc *runningConn) bool { return rc.conn == c })
	if i < 0 {
		return nil
	}
	return s.conns[i]
}

func (v *vat) runStep(rs *vatRunState) error {
	select {
	case nc := <-v.newConn:
		rc := v.runConn(rs.g, rs.ctx, nc)
		rs.conns = append(rs.conns, rc)

	case oc := <-v.connDone:
		rs.delConn(oc)

	case m := <-v.inMsg:
		// Process input msg.
		_ = m

	case m := <-v.outMsg:
		// Queue outgoing msg.
		rc := rs.findConn(m.conn)
		if rc != nil {
			rc.queueOut(m.msg)
		} // TODO: else alert caller conn does not exist anymore?

	case <-rs.ctx.Done():
		return rs.ctx.Err()
	}

	return nil
}

func (v *vat) Run(ctx context.Context) (err error) {
	rs := &vatRunState{
		g: pool.New().WithContext(ctx).WithCancelOnError().WithFirstError(),
	}
	rs.g.Go(func(ctx context.Context) error {
		var err error
		for err == nil {
			err = v.runStep(rs)
		}
		return err
	})

	return rs.g.Wait()
}
