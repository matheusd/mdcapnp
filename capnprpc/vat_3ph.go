// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"errors"
	"fmt"
)

type VatNetworkUniqueID [32]byte

type VatNetwork interface {
	introduce(src, target conn) (introductionInfo, error)
	connectToIntroduced(ctx context.Context, introducer conn, tcpd thirdPartyCapDescriptor) (conn, provisionId, error)
	recipientIdUniqueKey(recipientId) VatNetworkUniqueID
	provisionIdUniqueKey(provisionId) VatNetworkUniqueID
}

var err3PHWithoutVatNetwork = errors.New("3PH introductions not supported without VatNetwork")
var err3PHExpectedAcceptNotFound = errors.New("3PH prior Provide not found for Accept")

// getNetworkIntroduction generates a 3PH introduction. This runs on Bob, when
// generating an introduction that will be sent to Carol (as a Provide msg) and
// to Alice (inside a Return/Resolve that she will forward to Carol herself).
func (v *Vat) getNetworkIntroduction(src, target *runningConn) (introductionInfo, error) {
	if v.cfg.net == nil {
		return introductionInfo{}, err3PHWithoutVatNetwork
	}

	return v.cfg.net.introduce(src.c, target.c)
}

// startConnectToIntroduced3rdParty starts connection procedures to a third
// party, in order to (eventually) send an Accept. This runs on Alice, after
// receiving an introduction from Bob and is meant to connect to Carol.
func (v *Vat) startConnectToIntroduced3rdParty(ctx context.Context, introducer *runningConn,
	tpcd thirdPartyCapDescriptor) (connAndProvisionPromise, error) {

	if v.cfg.net == nil {
		return connAndProvisionPromise{}, err3PHWithoutVatNetwork
	}

	cpp := connAndProvisionPromise{capId: tpcd}

	go func() {
		c, provId, err := v.cfg.net.connectToIntroduced(ctx, introducer.c, tpcd)
		if err != nil {
			cpp.Fail(err)
			return
		}

		rc := v.RunConn(c)
		cpp.Fulfill(rc, provId)
	}()

	return cpp, nil
}

type expectedAccept struct {
	id VatNetworkUniqueID

	// srcConn is the original (aka provider, aka Bob) conn that sent a
	// Provide.
	srcConn *runningConn

	// provideAid is the question/answer associated with the original
	// Provide request.
	provideAid AnswerId

	// handler is the concrete capability that the provider (Bob) shared
	// with the initiating caller (Alice).
	handler callHandler
}

type getExpectedAccept struct {
	id        VatNetworkUniqueID
	replyChan chan any
}

type acceptedConnAndCap struct {
	// srcConn is the original (aka provider, aka Bob) conn that sent a
	// Provide.
	srcConn *runningConn

	// provideAid is the question/answer associated with the original
	// Provide request.
	provideAid AnswerId

	// handler is the concrete capability that the provider (Bob) shared
	// with the initiating caller (Alice).
	handler callHandler
}

// wasExpectingAccept determines if the given provisionId received on a conn is
// valid, and if so, what capability it refers to. This runs on Carol while
// processing an Accept received from Alice (rc), and so verifies if Bob
// previously sent a corresponding Provide.
func (v *Vat) wasExpectingAccept(ctx context.Context, provId provisionId) (expectedAccept, error) {
	if v.cfg.net == nil {
		return expectedAccept{}, err3PHWithoutVatNetwork
	}

	// Ask the vat run goroutine for the expected accept.
	id := v.cfg.net.provisionIdUniqueKey(provId)
	getAc := getExpectedAccept{id: id, replyChan: make(chan any, 1)}
	select {
	case v.getAccepts <- getAc:
	case <-ctx.Done():
		return expectedAccept{}, ctx.Err()
	}

	// Get the reply.
	var reply any
	select {
	case reply = <-getAc.replyChan:
	case <-ctx.Done():
		return expectedAccept{}, ctx.Err()
	}

	switch reply := reply.(type) {
	case error:
		return expectedAccept{}, reply
	case expectedAccept:
		return reply, nil
	default:
		panic(fmt.Sprintf("unsupported type in reply %T", reply))
	}
}
