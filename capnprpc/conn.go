// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"sync"

	"github.com/rs/zerolog"
	types "matheusd.com/mdcapnp/capnprpc/types"
)

type msgBatch struct {
	isSingle bool
	single   message
	msgs     []message
}

type outMsg struct {
	msg              *message
	remainingInBatch int
	sentChan         chan struct{}
}

func (om *outMsg) wantSentAck() *outMsg {
	om.sentChan = make(chan struct{})
	return om
}

func (om *outMsg) waitSentAck(ctx context.Context) error {
	select {
	case <-om.sentChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func singleMsgBatch(msg *message) outMsg {
	return outMsg{msg: msg, remainingInBatch: 0}
}

type conn interface {
	send(context.Context, message, int) error
	receive(context.Context) (message, error) // Ok because message goes to stack.
	remoteName() string

	// TODO: Allow conn-owned buffer (io_uring)?
	// usesReceiverBuffer() bool
	// receiveMsg(context.Context) (*message, error)
}

type ConnectionAndProvisionId struct {
	connection *runningConn
	provision  provisionId
}

type connAndProvisionPromise struct {
	capId thirdPartyCapDescriptor
}

func (cpp *connAndProvisionPromise) isWaiting() bool {
	return cpp.capId.vineId > 0 // FIXME: not a great way to track this.
}

func (cpp *connAndProvisionPromise) Wait(ctx context.Context) (ConnectionAndProvisionId, error) {
	panic("todo")
}

func (cpp *connAndProvisionPromise) Fail(err error) {
	panic("todo")
}

func (cpp *connAndProvisionPromise) Fulfill(rc *runningConn, provId provisionId) {
	panic("todo")
}

var errConnDone = errors.New("conn is done")

// runningConn is a connection that is running to another vat.
type runningConn struct {
	// Design note: most of the fields are only meant to be accessed from
	// within a vat's runStep() call. They are not safe for concurrent
	// access from within client code.
	//
	// TODO: maybe convert the public runningConn into a handle instead of
	// pointer?

	c   conn
	vat *Vat
	log zerolog.Logger
	rid uint64 // A unique running id, set by the vat when creating this.

	boot bootstrapCap

	// bootExportId is the export id of the bootstrap cap offered by the vat
	// on this conn.
	bootExportId ExportId

	outQueue chan outMsg

	crb callReturnBuilder

	// TODO: question and export IDs are set by local vat, answer and import
	// ids are set by the remote vat. Split table type into two
	// (incoming/outgoing table) to protect from remote misuse and restrict
	// API.
	mu        sync.Mutex
	questions questionsTable // table[QuestionId, question]
	answers   table[AnswerId, answer]
	imports   table[ImportId, imprt]
	exports   table[ExportId, export]

	ctx    context.Context
	cancel func(error) // Closes runningConn.
}

func (rc *runningConn) String() string {
	return rc.c.remoteName()
}

func (rc *runningConn) queue(ctx context.Context, m outMsg) error {
	/*
		rc.log.Trace().
			Int("remInBatch", m.remainingInBatch).
			Str("which", m.msg.Which().String()).
			Msg("Queueing outgoing message")
	*/

	which := m.msg.Which()

	select {
	case <-ctx.Done():
		return context.Cause(ctx)

	case rc.outQueue <- m:
		rc.log.Trace().
			Int("remInBatch", m.remainingInBatch).
			Str("which", which.String()).
			Msg("Queued outgoing message")
		return nil

	default:
		// Assume rc.outQueue is properly buffered. If the default case
		// is triggered, it means the buffer is full and sending is too
		// slow.
		return errors.New("outbound queue is full")
	}
}

func (rc *runningConn) cleanupQuestionIdDueToUnref(qid QuestionId) {
	if rc == nil {
		return
	}

	// Early check to see if conn is still running.
	if rc.ctx.Err() != nil {
		return
	}

	err := rc.vat.queueFinish(rc.ctx, rc, qid)
	if err != nil {
		rc.log.Err(err).Int("qid", int(qid)).Msg("Error sending Finish")
	} else {
		rc.log.Debug().Int("qid", int(qid)).Msg("Sent Finish")
	}
}

func (rc *runningConn) inLoop(ctx context.Context) error {
	v := rc.vat
	for {
		msg, err := rc.c.receive(ctx)
		if err != nil {
			return err
		}

		if msg.rawSerMsg != nil {
			var rpcMsg types.Message
			err = rpcMsg.ReadFromRoot(msg.rawSerMsg)
			if err == nil {
				err = v.processInMessageAlt(ctx, rc, rpcMsg, msg.rawSerMsg)
			}
		} else {
			err = v.processInMessage(ctx, rc, msg)
		}

		// Process input msg.
		if err != nil {
			return err
		}
	}
}

func (rc *runningConn) outLoop(ctx context.Context) error {
	v := rc.vat
	c := rc.c
	for {
		select {
		case outMsg := <-rc.outQueue:
			// Debug.
			/*
				if mb.msg.rawSerMb != nil {
					msgRawData, _ := mb.msg.rawSerMb.Serialize()
					rc.log.Trace().
						Hex("msg", msgRawData).
						Msg("DEBUG MSG")
				}
			*/

			err := c.send(ctx, *outMsg.msg, outMsg.remainingInBatch)
			if err != nil {
				return err
			}
			if outMsg.sentChan != nil {
				close(outMsg.sentChan)
			}

			if outMsg.msg.IsFinish() {
				v.mp.put(outMsg.msg)
			} else if outMsg.msg.rawSerMb != nil {
				v.mbp.put(outMsg.msg.rawSerMb)
				v.mp.put(outMsg.msg) // Transitional code
			}

		case <-ctx.Done():
			return context.Cause(ctx)
		}
	}
}

func newRunningConn(c conn, v *Vat) *runningConn {
	log := v.log.With().Str("remote", c.remoteName()).Logger()

	rc := &runningConn{
		c:   c,
		vat: v,
		log: log,

		boot: bootstrapCap(newRootFutureCap(v)),

		outQueue:  make(chan outMsg, 60000), // TODO: Parametrize buffer size.
		questions: makeQuestionsTable(),     // makeTable[QuestionId, question](),
		answers:   makeTable[AnswerId, answer](),
		imports:   makeTable[ImportId, imprt](),
		exports:   makeTable[ExportId, export](),
	}

	rc.crb.rc = rc

	// Prepare boot message.
	rc.boot.step.value.Set(pipeStepStateBuilding, pipelineStepStateValue{conn: rc})

	return rc
}

type bootstrapCap callFuture

func (bc bootstrapCap) Wait(ctx context.Context) (capability, error) {
	return castCallResultOrErr[capability](waitResult(ctx, callFuture(bc)))
}

func castBootstrap(bc bootstrapCap) callFuture {
	return callFuture{step: bc.step}
}

func (rc *runningConn) Bootstrap() bootstrapCap {
	return rc.boot // Any calls fork the pipeline.
}

type twoConnLocker struct {
	rc1 *runningConn
	rc2 *runningConn
}

func makeTwoConnLocker(rc1, rc2 *runningConn) twoConnLocker {
	// Total ordering based on run id.
	if rc2.rid < rc1.rid {
		rc1, rc2 = rc2, rc1
	}
	return twoConnLocker{rc1: rc1, rc2: rc2}
}

func (tcl *twoConnLocker) lock() {
	tcl.rc1.mu.Lock()
	tcl.rc2.mu.Lock()
}

func (tcl *twoConnLocker) unlock() {
	tcl.rc2.mu.Unlock()
	tcl.rc1.mu.Unlock()
}
