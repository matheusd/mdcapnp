package capnprpc

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/canastic/chantest"
	"matheusd.com/depvendoredtestify/require"
	"matheusd.com/mdcapnp/capnpser"
	"matheusd.com/testctx"
)

type fixedUserMethodsImpl struct {
	getAvatarDataCalled chan struct{}
	avatarData          string
}

func (f fixedUserMethodsImpl) GetAvatarData(cc StringReplyCallContext) error {
	cc.SetResponse(f.avatarData)
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
			gotAvatarData, err := methods.GetAvatarData().Wait(ctx)
			if err != nil {
				panic(err)
			}
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
			gotAvatarData, err := aliceBobCap.GetAvatarData().Wait(testctx.New(t))
			require.NoError(t, err)
			require.Equal(t, userMethodsCapImpl.avatarData, gotAvatarData)
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
			gotAvatarData, err := aliceBobCap.GetAvatarData().Wait(testctx.New(t))
			require.NoError(t, err)
			require.Equal(t, userMethodsCapImpl.avatarData, gotAvatarData)
		})
	}

}

type complexMethodServerImpl struct {
	execs atomic.Int64
}

func (c *complexMethodServerImpl) ComplexMethod(params testCapWithComplexMethod_ComplexMethod_Params, cc testCapWithComplexMethod_ComplexMethod_CallContext) error {
	c.execs.Add(1)
	intVal, strVal := params.Int64Value(), params.UnsafeStrValue()
	cc.AddBytesToSizeHint(len(strVal))
	res := cc.InitResults()
	res.SetInt64Value(intVal)
	res.SetStrValue(strVal)
	return nil
}

// BenchmarkGenericStructCallHandler benchmarks the generic StructCallHandler
// pattern for handling calls. Ideally, it should not produce any additional
// heap allocs.
func BenchmarkGenericStructCallHandler(b *testing.B) {
	const wantIntValue = 1701
	const wantStrValue = "the hills are alive with the sound of music"

	// Server bootstrap cap implementation.
	complexCapImpl := complexMethodServerImpl{}
	serverBootHandler := testCapWithComplexMethodServerAsHandler(&complexCapImpl)

	tests := []struct {
		name     string
		isLevel0 bool
	}{
		{name: "level0", isLevel0: true},
		{name: "level1+", isLevel0: false},
	}

	for _, tc := range tests {
		b.Run(tc.name, func(b *testing.B) {
			// Setup test harness.
			th := newTestHarness(b)
			s := th.newVat("server", WithBootstrapHandler(serverBootHandler))
			ctx := testctx.New(b)

			// Setup client.
			var bt Boostrappable
			if tc.isLevel0 {
				bt = th.newLevel0VatTo(s, "client", nil)
			} else {
				c := th.newVat("client")
				_, bt = th.connectVatsWithTCP(s, c)
			}

			// Wait for bootstrap.
			api, _, err := WaitBootstrap[testCapWithComplexMethod](ctx, bt)
			require.NoError(b, err)

			var i int
			requestSizeHint, _ := capnpser.ByteCount(len(wantStrValue)).StorageWordCount()

			// Run test.
			b.ReportAllocs()
			for b.Loop() {
				// Setup the parameters of the call.
				rb, future := api.ComplexMethod(requestSizeHint)
				rb.SetInt64Value(wantIntValue)
				if err := rb.SetStrValue(wantStrValue); err != nil {
					panic(err)
				}

				// Wait for call to complete.
				res, rr, err := future.Wait(ctx)
				if err != nil {
					b.Fatalf("at %d: %v", i, err)
				}

				// Double check resulting values are correct.
				gotInt := res.Int64Value()
				if gotInt != wantIntValue {
					b.Fatalf("wrong intvalue: got %v, want %v", gotInt, wantIntValue)
				}
				gotStr := res.UnsafeStrValue()
				if gotStr != wantStrValue {
					b.Fatalf("wrong strValue: got %v, want %v", gotStr, wantStrValue)
				}

				// Release the results (not gonna be used anymore).
				rr.Release()
				i++
			}
		})
	}
}
