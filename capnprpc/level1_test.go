package capnprpc

import (
	"context"
	"fmt"
	"testing"

	"github.com/canastic/chantest"
	"matheusd.com/depvendoredtestify/require"
	"matheusd.com/testctx"
)

type fixedUserMethodsImpl struct { // TODO: benchmark if this style of works.
	getAvatarDataCalled chan struct{}
	avatarData          string
}

func (f fixedUserMethodsImpl) GetAvatarData(cc testUserMethods_GetAvatarData_CallContext) error {
	cc.AddBytesToSizeHint(len(f.avatarData)) // Optional.
	rb := cc.InitResults()
	rb.SetAvatarData(f.avatarData)
	if f.getAvatarDataCalled != nil {
		f.getAvatarDataCalled <- struct{}{}
	}
	return nil
}

// TODO: Add a test sending a cap from Alice -> Bob, then calling this cap on
// Bob (full path for data is Alice -> Bob -> Alice).

// TestSenderHostedCapAsCallParams tests that sending a senderHosted capability
// as the sole Params in a Call works.
func TestSenderHostedCapAsCallParams(t *testing.T) {
	ctx := testctx.New(t)

	const avatarData string = "test avatar data"

	// Bootstrap handler for Bob. Only responds to the
	// CallGetAvatarFromUserMethods call.
	getAvatarReplyChan := make(chan string, 1)
	bobBootHandler := CallHandlerFunc(func(ctx context.Context, cc *CallContext) error {
		if err := cc.IsIidMid(testAPI_InterfaceID, testAPI_CallGetAvatarFromUserMethods_MethodId); err != nil {
			return err
		}

		methods, err := CallContextParamsCapability[testUserMethods](cc)
		if err != nil {
			return fmt.Errorf("unable to get 'methods' arg: %v", err)
		}

		// Decouple from processing to call the method (because we're
		// processing a Call from Alice right now, so trying to call
		// back will deadlock).
		// TODO: test case where we respond with a promise to the avatar
		// data.
		// TODO: test case where we decouple but don't respond anything
		// yet (respond from within the gorountine).
		go func() {
			res, _, err := methods.GetAvatarData().Wait(ctx)
			if err != nil {
				panic(err)
			}
			gotAvatarData := res.AvatarData()
			getAvatarReplyChan <- gotAvatarData

		}()
		return nil
	})

	// Setup harness.
	th := newTestHarness(t)
	alice := th.newVat("alice")
	bob := th.newVat("bob", WithBootstrapHandler(bobBootHandler))

	// Connect vats.
	aliceBobRc, _ := th.connectVatsWithTCP(alice, bob)

	// Alice asks Bob for the bootstrap cap. This is an API instance.
	aliceBobBoot := testAPIAsBootstrap(aliceBobRc.Bootstrap())

	// This is the testUserMethods cap implementation that Alice will
	// provide Bob.
	userMethodsCapImpl := fixedUserMethodsImpl{
		avatarData:          "test avatar data",
		getAvatarDataCalled: make(chan struct{}, 1),
	}

	// Alice asks bob to call the GetAvatar method of the passed
	// testUserMethods capability.
	bobReplyFuture := aliceBobBoot.CallGetAvatarFromUserMethods(userMethodsCapImpl)

	// The handler in Alice should be called.
	err := bobReplyFuture.Wait(ctx)
	require.NoError(t, err)
	chantest.AssertRecv(t, userMethodsCapImpl.getAvatarDataCalled)

	// Bob should get the reply sent by Alice.
	gotAvatarData := chantest.AssertRecv(t, getAvatarReplyChan).(string)
	require.Equal(t, avatarData, gotAvatarData)
}

// TestSenderHostedCapEmbeddedInCallParams tests that sending a senderHosted
// capability embedded in a struct of a Call.Params works.
func TestSenderHostedCapEmbeddedInCallParams(t *testing.T) {
}

// TestSenderHostedCapAsReturnResults tests that sending a senderHosted
// capability as the sole Results of a Return works.
func TestSenderHostedCapAsReturnResults(t *testing.T) {
	// Server side implementation of the testUserMethods capability
	// interface.
	userMethodsCapImpl := fixedUserMethodsImpl{avatarData: "test avatar data"}

	// Bootstrap handler for Bob. Only responds to the GetTestUserMethods
	// call.
	bobBootHandler := CallHandlerFunc(func(ctx context.Context, cc *CallContext) error {
		if err := cc.IsIidMid(testAPI_InterfaceID, testAPI_GetTestUserMethods_MethodId); err != nil {
			return err
		}

		// Respond with a single senderHosted cap (sole Payload contents
		// of Return.Results).
		userMethodsHandler := testUserMethodsCapServerAsCallHandler(userMethodsCapImpl)
		return cc.RespondAsSenderHostedCap(userMethodsHandler)
	})

	for _, wait := range []bool{true, false} {
		t.Run(fmt.Sprintf("wait=%v", wait), func(t *testing.T) {
			// Setup harness.
			th := newTestHarness(t)
			alice := th.newVat("alice")
			bob := th.newVat("bob", WithBootstrapHandler(bobBootHandler))

			// Connect vats.
			aliceBobRc, _ := th.connectVatsWithTCP(alice, bob)

			// Alice asks Bob for the bootstrap cap. This is an API instance.
			aliceBobBoot := testAPIAsBootstrap(aliceBobRc.Bootstrap())

			// Alice asks for the TestUserMethods() cap. This is the
			// actual test target.
			aliceBobCap := aliceBobBoot.GetTestUserMethods()

			// Test pipelining this call and not pipelining.
			if wait {
				_, err := aliceBobCap.Wait(testctx.New(t))
				require.NoError(t, err)
			}

			// Alice exercises a call on the returned capability
			// (double check it is working).
			gotAvatarData, _, err := aliceBobCap.GetAvatarData().Wait(testctx.New(t))
			require.NoError(t, err)
			require.Equal(t, userMethodsCapImpl.avatarData, gotAvatarData.AvatarData())
		})
	}
}

// TestSenderHostedCapEmbeddedInReturnResults tests that sending a senderHosted
// capability embedded in a struct of a Return.Results works.
func TestSenderHostedCapEmbeddedInReturnResults(t *testing.T) {
	// Server side implementation of the testUserMethods capability
	// interface.
	userMethodsCapImpl := fixedUserMethodsImpl{avatarData: "test avatar data"}

	// Bootstrap handler for Bob. Only responds to the GetTestUserMethods
	// call.
	bobBootHandler := CallHandlerFunc(func(ctx context.Context, cc *CallContext) error {
		if err := cc.IsIidMid(testAPI_InterfaceID, testAPI_GetTestUserMethods_MethodId); err != nil {
			return err
		}

		// Respond with a single senderHosted cap (sole Payload contents
		// of Return.Results).
		userMethodsHandler := testUserMethodsCapServerAsCallHandler(userMethodsCapImpl)
		return cc.RespondAsSenderHostedCap(userMethodsHandler)
	})

	for _, wait := range []bool{true, false} {
		t.Run(fmt.Sprintf("wait=%v", wait), func(t *testing.T) {
			// Setup harness.
			th := newTestHarness(t)
			alice := th.newVat("alice")
			bob := th.newVat("bob", WithBootstrapHandler(bobBootHandler))

			// Connect vats.
			aliceBobRc, _ := th.connectVatsWithTCP(alice, bob)

			// Alice asks Bob for the bootstrap cap. This is an API instance.
			aliceBobBoot := testAPIAsBootstrap(aliceBobRc.Bootstrap())

			// Alice asks for the TestUserMethods() cap. This is the
			// actual test target.
			aliceBobCap := aliceBobBoot.GetTestUserMethods()

			// Test pipelining this call and not pipelining.
			if wait {
				_, err := aliceBobCap.Wait(testctx.New(t))
				require.NoError(t, err)
			}

			// Alice exercises a call on the returned capability
			// (double check it is working).
			gotAvatarData, _, err := aliceBobCap.GetAvatarData().Wait(testctx.New(t))
			require.NoError(t, err)
			require.Equal(t, userMethodsCapImpl.avatarData, gotAvatarData.AvatarData())
		})
	}

}

// BenchmarkGenericStructCallHandler benchmarks the generic StructCallHandler
// pattern for handling calls. Ideally, it should not produce any additional
// heap allocs.
func BenchmarkGenericStructCallHandler(b *testing.B) {
	userMethodsCapImpl := fixedUserMethodsImpl{avatarData: "test avatar data"}
	userMethodsHandler := testUserMethodsCapServerAsCallHandler(userMethodsCapImpl)

	th := newTestHarness(b)
	s := th.newVat("server", WithBootstrapHandler(userMethodsHandler))

	io1, io2 := th.tcpTransportPair("client", "server")
	s.RunConn(io2)

	c := th.newLevel0Vat("client", &Level0ClientCfg{Conn: io1})

	ctx := testctx.New(b)

	// Wait for bootstrap.
	boot := c.Bootstrap()
	_, err := boot.Wait(ctx)
	require.NoError(b, err)
	api := testUserMethods(boot)

	var i int
	b.ReportAllocs()
	for b.Loop() {
		_, _, err = api.GetAvatarData().Wait(ctx)
		if err != nil {
			b.Fatalf("at %d: %v", i, err)
		}
		i++
	}
}
