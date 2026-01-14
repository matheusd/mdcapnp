// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"

	"matheusd.com/mdcapnp/capnpser"
)

type futureString CallFuture

func (fs futureString) wait(ctx context.Context) (string, error) {
	return CastCallResultOrErr[string](WaitReturn(ctx, CallFuture(fs)))
}

type testAPI CallFuture

const testAPI_InterfaceID = 1000

type futureVoid CallFuture

func (fv futureVoid) Wait(ctx context.Context) error {
	_, err := WaitReturn(ctx, CallFuture(fv))
	return err
}

const testAPI_Void_CallID = 101
const testAPI_GetAnotherAPI_CallID = 102
const testAPI_GetUser_CallID = 103
const testAPI_Add_CallID = 104
const testAPI_GetTestUserMethods_MethodId = 105
const testAPI_CallGetAvatarFromUserMethods_MethodId = 106

func (api testAPI) VoidCall() futureVoid {
	return futureVoid(RemoteCall(
		CallFuture(api),
		SetupCallNoParams(CallFuture(api),
			testAPI_InterfaceID,
			testAPI_Void_CallID,
		),
	))
}

var addRequestSize = capnpser.StructSize{DataSectionSize: 2, PointerSectionSize: 0}

type addRequestBuilder capnpser.StructBuilder

func (b *addRequestBuilder) SetA(v int64) error {
	return (*capnpser.StructBuilder)(b).SetInt64(0, v)
}

func (b *addRequestBuilder) SetB(v int64) error {
	return (*capnpser.StructBuilder)(b).SetInt64(1, v)
}

func newAddRequestBuilder(serMsg *capnpser.MessageBuilder) (addRequestBuilder, error) {
	return capnpser.NewStructBuilder[addRequestBuilder](serMsg, addRequestSize)
}

type addRequest capnpser.Struct

func (s *addRequest) A() int64 {
	return (*capnpser.Struct)(s).Int64(0)
}

func (s *addRequest) B() int64 {
	return (*capnpser.Struct)(s).Int64(1)
}

var addResponseSize = capnpser.StructSize{DataSectionSize: 1, PointerSectionSize: 0}

type addResponseBuilder capnpser.StructBuilder

func (b *addResponseBuilder) SetC(v int64) error {
	return (*capnpser.StructBuilder)(b).SetInt64(0, v)
}

type addResponse capnpser.Struct

func (s *addResponse) C() int64 {
	return (*capnpser.Struct)(s).Int64(0)
}

type futureAddResult CallFuture

func (fut futureAddResult) wait(ctx context.Context) (res int64, err error) {
	r, rr, err := WaitShallowCopyReturnResultsStruct[addResponse](ctx, CallFuture(fut))
	if err != nil {
		return
	}
	res = r.C()
	rr.Release()
	return
}

func (api testAPI) Add(a int64, b int64) futureAddResult {
	cs, req := SetupCallWithStructParamsGeneric[addRequestBuilder](
		CallFuture(api),
		addRequestSize.TotalSize(),
		testAPI_InterfaceID,
		testAPI_Add_CallID,
		addRequestSize,
	)

	req.SetA(a)
	req.SetB(b)
	cs.WantShallowReturnCopy = true

	return futureAddResult(RemoteCall(
		CallFuture(api),
		cs,
	))
}

func (api testAPI) GetAnotherAPICap() testAPI {
	return testAPI(RemoteCall(
		CallFuture(api),
		CallSetup{
			InterfaceId: testAPI_InterfaceID,
			MethodId:    testAPI_GetAnotherAPI_CallID,
		},
	))
}

func (api testAPI) GetTestUserMethods() testUserMethods {
	return testUserMethods(RemoteCall(
		CallFuture(api),
		SetupCallNoParams(
			CallFuture(api),
			testAPI_InterfaceID,
			testAPI_GetTestUserMethods_MethodId,
		),
	))
}

func (api testAPI) CallGetAvatarFromUserMethods(s testUserMethodsCapServer) VoidFuture {
	cs := SetupCallWithNewCap(
		CallFuture(api),
		testAPI_InterfaceID,
		testAPI_CallGetAvatarFromUserMethods_MethodId,
		testUserMethodsCapServerAsCallHandler(s),
	)

	return VoidFuture(RemoteCall(
		CallFuture(api),
		cs,
	))
}

/*
func (api testAPI) GetUser(id string) testUser {
	return testUser(RemoteCall(
		CallFuture(api),
		CallSetup{
			InterfaceId: testAPI_InterfaceID,
			MethodId:    testAPI_GetUser_CallID,
		},
	))
}
*/

// Wait until this is resolved as a concrete, exported capability.
func (api testAPI) Wait(ctx context.Context) (CapabilityClient, error) {
	return CastCallResultOrErr[CapabilityClient](WaitReturn(ctx, CallFuture(api)))
}

func (api testAPI) WaitDiscardResult(ctx context.Context) error {
	_, err := WaitReturn(ctx, CallFuture(api))
	return err
}

func testAPIAsBootstrap(bt BootstrapFuture) testAPI {
	return testAPI(bt)
}

const testUserMethods_InterfaceId InterfaceId = 1001
const testUserMethods_getAvatarData_methodId MethodId = 11

type testUserMethods CallFuture

func (up testUserMethods) Wait(ctx context.Context) (CapabilityClient, error) {
	return CastCallResultOrErr[CapabilityClient](WaitReturn(ctx, CallFuture(up)))
}

var testUserMethods_GetAvatarData_ResultsSize = capnpser.StructSize{DataSectionSize: 0, PointerSectionSize: 1}

type testUserMethods_GetAvatarData_Results capnpser.Struct

func (s *testUserMethods_GetAvatarData_Results) AvatarData() string {
	return (*capnpser.Struct)(s).String(0)
}

type testUserMethods_GetAvatarData_ResultsBuilder capnpser.StructBuilder

func (b *testUserMethods_GetAvatarData_ResultsBuilder) SetAvatarData(s string) error {
	return (*capnpser.StructBuilder)(b).SetString(0, s)
}

type testUserMethods_GetAvatarData_Results_future CallFuture

func (f testUserMethods_GetAvatarData_Results_future) Wait(ctx context.Context) (testUserMethods_GetAvatarData_Results, ReturnResults, error) {
	return WaitShallowCopyReturnResultsStruct[testUserMethods_GetAvatarData_Results](ctx, CallFuture(f))
}

// GetAvatarData @0 () -> (AvatarData :text)
func (up testUserMethods) GetAvatarData() testUserMethods_GetAvatarData_Results_future {
	cs := SetupCallNoParams(
		CallFuture(up),
		testUserMethods_InterfaceId,
		testUserMethods_getAvatarData_methodId,
	)
	cs.WantShallowReturnCopy = true

	return testUserMethods_GetAvatarData_Results_future(RemoteCall(
		CallFuture(up),
		cs,
	))
}

type testUserMethods_GetAvatarData_CallContext = StructReplyCallContext[testUserMethods_GetAvatarData_ResultsBuilder]

type testUserMethodsCapServer interface {
	GetAvatarData(ri testUserMethods_GetAvatarData_CallContext) error
}

func testUserMethodsCapServerAsCallHandler(s testUserMethodsCapServer) CallHandler {
	return CallHandlerFunc(func(ctx context.Context, cc *CallContext) error {
		iid, mid := cc.InterfaceId(), cc.MethodId()
		if iid != testUserMethods_InterfaceId {
			return ErrUnexpectedInterfaceId{Want: testUserMethods_InterfaceId, Got: iid}
		}
		switch mid {
		case testUserMethods_getAvatarData_methodId:
			return s.GetAvatarData(testUserMethods_GetAvatarData_CallContext(
				MakeStructReplyCallContext[testUserMethods_GetAvatarData_ResultsBuilder](
					cc, testUserMethods_GetAvatarData_ResultsSize,
				)))
		default:
			return ErrUnknownMethodId{IId: iid, Mid: mid}
		}
	})
}

type testUser capnpser.Struct

func (st *testUser) Username() string {
	return (*capnpser.Struct)(st).String(0)
}

func (st *testUser) Methods() testUserMethods {
	// return (*capnpser.Struct)(st).ReadAnyPointer()
	panic("boo")
}

/*
type testUser CallFuture

func (usr testUser) GetProfile() testUserProfile {
	return testUserProfile(RemoteCall(
		CallFuture(usr),
		CallSetup{
			InterfaceId: 1000,
			MethodId:    11,
		},
	))
}

type testUserProfile CallFuture

func (up testUserProfile) GetAvatarData() futureString {
	return futureString(RemoteCall(
		CallFuture(up),
		CallSetup{
			InterfaceId: 1000,
			MethodId:    11,
		},
	))
}
*/
