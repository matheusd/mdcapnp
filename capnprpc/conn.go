// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/rs/zerolog"
	types "matheusd.com/mdcapnp/capnprpc/types"
	"matheusd.com/mdcapnp/capnpser"
)

type outMsg struct {
	serMsg   *capnpser.MessageBuilder
	mb       types.MessageBuilder
	sentChan chan struct{}
}

func (om *outMsg) ackSent() {
	if om.sentChan != nil {
		close(om.sentChan)
	}
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

type OutMsg struct {
	Msg *capnpser.MessageBuilder
}

type InMsg struct {
	Msg capnpser.Message
}

type conn interface {
	send(context.Context, OutMsg) error
	receive(context.Context) (InMsg, error) // Ok because message goes to stack.
	close() error
	remoteName() string

	// TODO: Allow conn-owned buffer (io_uring)?
	// usesReceiverBuffer() bool
	// receiveMsg(context.Context) (*message, error)
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

	boot BootstrapFuture

	// bootExportId is the export id of the bootstrap cap offered by the vat
	// on this conn.
	bootExportId ExportId

	outQueue chan outMsg

	crb CallContext

	// TODO: question and export IDs are set by local vat, answer and import
	// ids are set by the remote vat. Split table type into two
	// (incoming/outgoing table) to protect from remote misuse and restrict
	// API.
	mu        sync.Mutex
	questions questionsTable // table[QuestionId, question]
	answers   table[AnswerId, answer]
	imports   table[ImportId, imprt]
	exports   table[ExportId, export]

	// inBuf []byte

	ctx    context.Context
	cancel func(error) // Closes runningConn.
}

func (rc *runningConn) String() string {
	return rc.c.remoteName()
}

func (rc *runningConn) queue(ctx context.Context, m outMsg) error {
	// mr := m.mb.AsReader()
	// which := mr.Which().String()

	/*
		rc.log.Trace().
			Str("which", which.String()).
			Msg("Queueing outgoing message")
	*/

	select {
	case <-ctx.Done():
		return context.Cause(ctx)

	case rc.outQueue <- m:
		/*
			rc.log.Trace().
				Str("which", which).
				Msg("Queued outgoing message")
		*/
		return nil

	default:
		// Assume rc.outQueue is properly buffered. If the default case
		// is triggered, it means the buffer is full and sending is too
		// slow.
		return errors.New("outbound queue is full")
	}
}

func (rc *runningConn) send(ctx context.Context, outMsg outMsg) error {
	err := rc.c.send(ctx, OutMsg{Msg: outMsg.serMsg})
	outMsg.ackSent()
	if err != nil {
		return fmt.Errorf("send errored: %w", err)
	}

	rc.vat.mbp.put(outMsg.serMsg)
	return nil
}

type InboundReturn struct {
	SerMsg capnpser.Message
	Msg    types.Message
	Ret    types.Return
}

func (rc *runningConn) nextReturn(ctx context.Context, inRet *InboundReturn) error {
	msg, err := rc.c.receive(ctx)
	if err != nil {
		return err
	}

	err = inRet.Msg.ReadFromRoot(&msg.Msg)
	if err != nil {
		return err
	}

	if inRet.Msg.Which() != types.Message_Which_Return {
		return fmt.Errorf("message was not a return")
	}

	inRet.Ret, err = inRet.Msg.AsReturn()
	if err != nil {
		return err
	}

	return nil
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
			if errors.Is(err, io.EOF) {
				// EOF is considered a graceful conn shutdown.
				err = nil
			} else {
				err = fmt.Errorf("inLoop receive errored: %w", err)
			}
			return err
		}

		// Debug.
		/*
			logEvent := rc.log.Debug()
			msgRawData := msg.Msg.Arena().RawDataCopy()
			for i, data := range msgRawData {
				logEvent.Hex(fmt.Sprintf("msg.seg%d", i), data)
			}
			logEvent.Msg("debug processIn")
		*/

		var rpcMsg types.Message
		err = rpcMsg.ReadFromRoot(&msg.Msg)
		if err == nil {
			err = v.processInMessage(ctx, rc, rpcMsg, msg.Msg)
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

			err := c.send(ctx, OutMsg{Msg: outMsg.serMsg})
			outMsg.ackSent()
			if err != nil {
				return fmt.Errorf("outLoop send errored: %w", err)
			}

			v.mbp.put(outMsg.serMsg)

		case <-ctx.Done():
			return context.Cause(ctx)
		}
	}
}

func (rc *runningConn) waitToClose(ctx context.Context) error {
	<-ctx.Done()
	return rc.c.close()
}

func newRunningConn(c conn, v *Vat) *runningConn {
	log := v.log.With().Str("remote", c.remoteName()).Logger()

	rc := &runningConn{
		c:   c,
		vat: v,
		log: log,

		boot: BootstrapFuture(newRootFutureCap(v)),

		outQueue:  make(chan outMsg, 100), // TODO: Parametrize buffer size.
		questions: makeQuestionsTable(),   // makeTable[QuestionId, question](),
		answers:   makeTable[AnswerId, answer](),
		imports:   makeTable[ImportId, imprt](),
		exports:   makeTable[ExportId, export](),
	}

	rc.crb.vat = v
	rc.crb.rc = rc

	// Prepare boot message.
	rc.boot.step.value.Set(pipeStepStateBuilding, pipelineStepStateValue{conn: rc})

	return rc
}

func (rc *runningConn) Bootstrap() BootstrapFuture {
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
