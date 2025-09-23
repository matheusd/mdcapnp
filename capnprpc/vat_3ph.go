// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"errors"
	"fmt"
)

type VatNetwork interface {
	introduce(src, target conn) (introductionInfo, error)
}

// getNetworkIntroduction generates a 3PH introduction. This runs on Bob, when
// generating an introduction that will be sent to Carol (as a Provide msg) and
// to Alice (inside a Return/Resolve that she will forward to Carol herself).
func (v *Vat) getNetworkIntroduction(src, target *runningConn) (introductionInfo, error) {
	if v.cfg.net == nil {
		return introductionInfo{}, errors.New("3PH introductions not supported without VatNetwork")
	}

	return v.cfg.net.introduce(src.c, target.c)
}

// startConnectToIntroduced3rdParty starts connection procedures to a third
// party, in order to (eventually) send an Accept. This runs on Alice, after
// receiving an introduction from Bob and is meant to connect to Carol.
func (v *Vat) startConnectToIntroduced3rdParty(introducer *runningConn, tpcd thirdPartyCapDescriptor) (connAndProvisionPromise, error) {
	// TODO: fulfill the conn request.
	cpp := connAndProvisionPromise{capId: tpcd}
	return cpp, nil
}

// expectConnAndAccept sets up the vat to wait for a connection in relation to a
// 3PH introduction. This runs on Carol while processing a Provide from Bob
// (setting things up to expect a future Accept from Alice).
func (v *Vat) expectConnAndAccept(introducer *runningConn, recipient recipientId) error {
	// TODO: wat?
	return nil
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
func (v *Vat) wasExpectingAccept(rc *runningConn, provId provisionId) (acceptedConnAndCap, error) {
	// TODO: find the matching recipientId. Check if it exists, which
	// srcConn it refers to and which capability.
	var srcConn *runningConn
	var target messageTarget
	var provideAid AnswerId

	// Check if the target still exists exported to srcConn. Determine if
	// this is a capability or a promise to a capability.
	//
	// TODO: not safe to lock srcConn here. Can lead to deadlocks. Maybe
	// this should be moved upwards, so that this information is stored and
	// then returned by whatever returns srcConn.
	var err error
	var handler callHandler
	srcConn.mu.Lock()
	if target.isImportedCap {
		exp, hasExp := srcConn.exports.get(ExportId(target.impcap))
		if !hasExp {
			err = fmt.Errorf("export not found %d", target.impcap)
		}

		handler = exp.handler
	} else if target.isPromisedAnswer {
		if !srcConn.answers.has(AnswerId(target.pans.qid)) {
			err = fmt.Errorf("answer not found %d", target.pans.qid)
		}
	} else {
		err = errors.New("unknown message target")
	}
	srcConn.mu.Unlock()
	if err != nil {
		return acceptedConnAndCap{}, err
	}

	return acceptedConnAndCap{
		srcConn:    srcConn,
		provideAid: provideAid,
		handler:    handler,
	}, nil
}
