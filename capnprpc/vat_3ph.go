// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"matheusd.com/mdcapnp/capnpser"
)

type VatNetworkUniqueID [32]byte

func (id VatNetworkUniqueID) toString() string {
	return hex.EncodeToString(id[:])
}

type VatNetwork interface {
	introduce(src, target conn) (introductionInfo, error)
	connectToIntroduced(ctx context.Context, localVat *Vat, introducer conn, tcpd capnpser.AnyPointer) (conn, capnpser.AnyPointer, error)
	recipientIdUniqueKey(capnpser.AnyPointer) VatNetworkUniqueID
	provisionIdUniqueKey(capnpser.AnyPointer) VatNetworkUniqueID
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

// connectToIntroduced3rdParty connects to a third party, in order to
// (eventually) send an Accept. This runs on Alice, after receiving an
// introduction from Bob and is meant to connect to Carol.
func (v *Vat) connectToIntroduced3rdParty(ctx context.Context, introducer *runningConn,
	tpcd capnpser.AnyPointer) (*runningConn, capnpser.AnyPointer, error) {

	if v.cfg.net == nil {
		return nil, capnpser.AnyPointer{}, err3PHWithoutVatNetwork
	}

	// TODO: check if already connected or already trying to connect to the
	// target vat and reuse conn.

	c, provId, err := v.cfg.net.connectToIntroduced(ctx, v, introducer.c, tpcd)
	if err != nil {
		return nil, capnpser.AnyPointer{}, err
	}

	rc := v.RunConn(c)
	return rc, provId, nil
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
func (v *Vat) wasExpectingAccept(ctx context.Context, provId capnpser.AnyPointer) (expectedAccept, error) {
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
