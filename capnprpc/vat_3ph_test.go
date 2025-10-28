// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/canastic/chantest"
	"matheusd.com/depvendoredtestify/require"
	"matheusd.com/mdcapnp/capnpser"
	"matheusd.com/testctx"
)

type testVatId int

func (tvid testVatId) toBytes() []byte {
	return binary.LittleEndian.AppendUint32(nil, uint32(tvid))
}

type testProvisionIdNonce uint64

func (nonce testProvisionIdNonce) encodeProvisionFor(remoteId testVatId) []byte {
	res := make([]byte, 0, 12)
	res = binary.LittleEndian.AppendUint32(res, uint32(remoteId))
	res = binary.LittleEndian.AppendUint64(res, uint64(nonce))
	return res
}

func (nonce testProvisionIdNonce) encodeProvisionAsAnyPtr(remoteId testVatId) capnpser.AnyPointer {
	bytes := nonce.encodeProvisionFor(remoteId)
	mb, err := capnpser.NewMessageBuilder(capnpser.DefaultSimpleSingleAllocator)
	if err != nil {
		panic(err)
	}
	res, err := mb.CopyToNewByteList(bytes)
	if err != nil {
		panic(err)
	}
	return res.Reader()
}

func newTestProvisionIdNonce() testProvisionIdNonce {
	return testProvisionIdNonce(rand.Uint64())
}

type testProvisionId struct {
	remoteId testVatId
	nonce    testProvisionIdNonce
}

func (tpid testProvisionId) toRawBytes() []byte {
	return tpid.nonce.encodeProvisionFor(tpid.remoteId)
}

func testProvisionIdFromBytes(raw []byte) (res testProvisionId) {
	res.remoteId = testVatId(binary.LittleEndian.Uint32(raw))
	res.nonce = testProvisionIdNonce(binary.LittleEndian.Uint64(raw[4:]))
	return
}

func testProvisionIdFromAnyPointerBytes(ptr capnpser.AnyPointer) (res testProvisionId, err error) {
	if !ptr.IsList() {
		err = fmt.Errorf("testProvisionIdFromAnyPointerBytes: ptr is not a list")
		return
	}
	var rawTcpd []byte
	tcpdLs := ptr.AsList()
	rawTcpd, err = tcpdLs.Bytes()
	if err != nil {
		return
	}
	res = testProvisionIdFromBytes(rawTcpd)
	return
}

type testVatNetwork struct {
	th             *testHarness
	vatIdToTestVat map[testVatId]*testVat
	vatToTestVat   map[*Vat]*testVat
}

func newTestVatNetwork(t testing.TB) *testVatNetwork {
	th := newTestHarness(t)
	return &testVatNetwork{
		th:             th,
		vatIdToTestVat: make(map[testVatId]*testVat),
		vatToTestVat:   make(map[*Vat]*testVat),
	}
}

func (t *testVatNetwork) introduce(src conn, target conn) (introductionInfo, error) {
	srcIndex := testVatId(src.(*testPipeConn).remIndex)
	targetIndex := testVatId(target.(*testPipeConn).remIndex)
	nonce := newTestProvisionIdNonce()
	toRecProvId := nonce.encodeProvisionFor(targetIndex)
	toTargetProvId := nonce.encodeProvisionFor(srcIndex)

	// mb to use to store both the recipient and target
	// thirdPartyContactInfo.
	mb, err := capnpser.NewMessageBuilder(capnpser.DefaultSimpleSingleAllocator)
	if err != nil {
		return introductionInfo{}, err
	}
	sendToRec, err := mb.CopyToNewByteList(toRecProvId)
	if err != nil {
		return introductionInfo{}, err
	}
	sendToTarget, err := mb.CopyToNewByteList(toTargetProvId)
	if err != nil {
		return introductionInfo{}, err
	}

	return introductionInfo{
		sendToRecipientAlt: sendToRec.Reader(),
		sendToTargetAlt:    sendToTarget.Reader(),
	}, nil
}

func (t *testVatNetwork) connectToIntroduced(ctx context.Context, localVat *Vat,
	introducer conn, tcpd capnpser.AnyPointer) (conn, capnpser.AnyPointer, error) {

	// ProvisionId that Bob sent to Alice.
	provId, err := testProvisionIdFromAnyPointerBytes(tcpd) // testProvisionIdFromBytes(tcpd.id.st.rawData)
	if err != nil {
		return nil, capnpser.AnyPointer{}, fmt.Errorf("connectToIntroduced invalid provisionId: %v", err)
	}
	remoteVat, ok := t.vatIdToTestVat[provId.remoteId]
	if !ok {
		return nil, capnpser.AnyPointer{}, fmt.Errorf("could not find vat with remoteId %d", provId.remoteId)
	}

	localTestVat, ok := t.vatToTestVat[localVat]
	if !ok { // Should not happen.
		return nil, capnpser.AnyPointer{}, fmt.Errorf("bug: could not find local vat in network")
	}

	// ProvisionId that Alice sends to Carol. Replace with Alice's Id so
	// that it matches what Bob sent to Carol.
	remoteProvId := provId.nonce.encodeProvisionAsAnyPtr(testVatId(localTestVat.index))

	// Connect vats.
	localConn, remoteConn := t.th.twoVatsPipe(localTestVat, remoteVat)
	remoteVat.RunConn(remoteConn)
	return localConn, remoteProvId, nil
}

func (t *testVatNetwork) recipientIdUniqueKey(rid capnpser.AnyPointer) (res VatNetworkUniqueID) {
	if !rid.IsList() {
		panic("rid is not a List")
	}
	ls := rid.AsList()
	bytes, err := ls.Bytes()
	if err != nil {
		panic(err)
	}
	copy(res[:], bytes)
	return
}

func (t *testVatNetwork) provisionIdUniqueKey(pid capnpser.AnyPointer) (res VatNetworkUniqueID) {
	if !pid.IsList() {
		panic("rid is not a List")
	}
	ls := pid.AsList()
	bytes, err := ls.Bytes()
	if err != nil {
		panic(err)
	}
	copy(res[:], bytes)
	return
}

func (tvn *testVatNetwork) newVat(name string, opts ...VatOption) *testVat {
	opts = append(opts, WithVatNetwork(tvn))
	v := tvn.th.newVat(name, opts...)
	tvn.vatIdToTestVat[testVatId(v.index)] = v
	tvn.vatToTestVat[v.Vat] = v
	return v
}

func (tvn *testVatNetwork) connectVats(v1, v2 *testVat) (rc1, rc2 *runningConn) {
	return tvn.th.connectVats(v1, v2)
}

// Test3PHBasic tests a basic implementation of 3PH without cached or pipelined
// calls.
func Test3PHBasic(t *testing.T) {
	// Bob's bootstrap cap is a remote promise.
	bobHandlerCalledChan := make(chan answerPromise, 1)
	bobHandler := CallHandlerFunc(func(ctx context.Context, rb *CallContext) error {
		ap, err := rb.respondAsPromise()
		if err != nil {
			return err
		}

		bobHandlerCalledChan <- ap
		return nil
	})

	// Carol's bootstrap cap is a concrete handler that always returns
	// void.
	carolHandlerCalled := make(chan struct{}, 1)
	carolHandler := CallHandlerFunc(func(ctx context.Context, rb *CallContext) error {
		carolHandlerCalled <- struct{}{}
		return nil
	})

	stageLog := func(format string, args ...any) {
		time.Sleep(100 * time.Millisecond) // TODO: parametrize if it exists.
		t.Logf("================ "+format, args...)
	}

	// Setup harness.
	tnet := newTestVatNetwork(t)
	alice := tnet.newVat("alice")
	bob := tnet.newVat("bob", WithBootstrapHandler(bobHandler))
	carol := tnet.newVat("carol", WithBootstrapHandler(carolHandler))

	// Existing connections (before 3PH).
	aliceBobRc, _ := tnet.connectVats(alice, bob)
	bobCarolRc, _ := tnet.connectVats(bob, carol)
	stageLog("Initial connections done. Alice will ask for Bob's Bootstrap.")

	// Alice asks Bob for the bootstrap cap. This is an API instance.
	bobApiInAlice := testAPIAsBootstrap(aliceBobRc.Bootstrap())
	require.NoError(t, bobApiInAlice.WaitDiscardResult(testctx.New(t)))
	stageLog("Alice got Bob's bootstrap cap. Alice will ask for a sub-cap from boot.")

	// Alice asks for a capability from Bob. This will only complete after
	// 3PH completes, so wait for it in a goroutine.
	aliceGetCapErrChan := make(chan error, 1)
	aliceGetCapCall := bobApiInAlice.GetAnotherAPICap()
	go func() {
		aliceGetCapErrChan <- aliceGetCapCall.WaitDiscardResult(testctx.New(t))
	}()
	stageLog("Alice asked Bob for a new cap. Bob will fetch Carol's Boostrap cap.")

	// Bob realizes the cap that Alice wants is in Carol. So Bob will ask
	// Carol for it (it's her bootstrap cap). This is done OOB here, but in
	// a real application this would be a goroutine triggered by Bob's
	// original handler (bobHandler()).
	//
	// We explicitly wait for the bootstrap to be returned, but, presumably,
	// we could also fulfill it with just the promise of the bootstrap's
	// result.
	carolCapInBob, err := testAPIAsBootstrap(bobCarolRc.Bootstrap()).Wait(testctx.New(t))
	require.NoError(t, err)
	stageLog("Bob got Carol's bootstrap cap. Bob will fulfill Alice's promise with it.")

	// Bob will now answer its promise to Alice with a third party cap:
	// Carol's Bootstrap.
	bobsPromiseToAlice := chantest.AssertRecv(t, bobHandlerCalledChan).(answerPromise)
	bobsPromiseToAlice.resolveToThirdPartyCap(bobCarolRc, carolCapInBob)

	// Alice's call to get a new cap will complete: this means 3PH completed
	// and Alice now holds Carol's bootstrap cap as an explicit imported
	// cap.
	aliceGetCapErr := chantest.AssertRecv(t, aliceGetCapErrChan)
	require.Nil(t, aliceGetCapErr)
	stageLog("3PH completed! Alice will make a path-shortened call to Carol.")

	// Finally, Alice makes a call that will go to Carol. Note that
	// aliceGetCapCall was originally obtained from bobApiInAlice (which is
	// an interface to Bob's cap) but it resolved (after path-shortening
	// 3PH) to a Carol cap.
	chantest.AssertNoRecv(t, carolHandlerCalled) // Sanity precondition check.
	require.NoError(t, aliceGetCapCall.VoidCall().Wait(testctx.New(t)))
	stageLog("Path-shortened call completed!")
}
