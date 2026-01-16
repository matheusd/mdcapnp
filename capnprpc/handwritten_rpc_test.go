// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"

	"matheusd.com/mdcapnp/capnpser"
)

type testAPI CallFuture

const testAPI_InterfaceID = 1000

const testAPI_Void_CallID = 101
const testAPI_GetAnotherAPI_CallID = 102
const testAPI_GetUser_CallID = 103
const testAPI_Add_CallID = 104
const testAPI_GetTestUserMethods_MethodId = 105
const testAPI_CallGetAvatarFromUserMethods_MethodId = 106

func (api testAPI) VoidCall() VoidFuture {
	return VoidFuture(RemoteCall(
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
	r, rr, err := WaitReturnResultsStruct[addResponse](ctx, CallFuture(fut))
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
	return waitReturnResultsCapability(ctx, CallFuture(api))
}

func (api testAPI) WaitDiscardResult(ctx context.Context) error {
	_, err := api.Wait(ctx)
	return err
}

func testAPIAsBootstrap(bt BootstrapFuture) testAPI {
	return testAPI(bt)
}

const testUserMethods_InterfaceId InterfaceId = 1001
const testUserMethods_getAvatarData_methodId MethodId = 11

type testUserMethods CallFuture

func (up testUserMethods) Wait(ctx context.Context) (CapabilityClient, error) {
	return waitReturnResultsCapability(ctx, CallFuture(up))
}

var testUserMethods_GetAvatarData_ResultsSize = capnpser.StructSize{DataSectionSize: 0, PointerSectionSize: 1}

// GetAvatarData @0 () -> (AvatarData :text)
func (up testUserMethods) GetAvatarData() StringFuture {
	cs := SetupCallNoParams(
		CallFuture(up),
		testUserMethods_InterfaceId,
		testUserMethods_getAvatarData_methodId,
	)
	return StringFuture(RemoteCall(
		CallFuture(up),
		cs,
	))
}

type testUserMethodsCapServer interface {
	GetAvatarData(StringReplyCallContext) error
}

func testUserMethodsCapServerAsCallHandler(s testUserMethodsCapServer) CallHandler {
	return CallHandlerFunc(func(ctx context.Context, cc *CallContext) error {
		iid, mid := cc.InterfaceId(), cc.MethodId()
		if iid != testUserMethods_InterfaceId {
			return ErrUnexpectedInterfaceId{Want: testUserMethods_InterfaceId, Got: iid}
		}
		switch mid {
		case testUserMethods_getAvatarData_methodId:
			return s.GetAvatarData(MakeStringReplyCallContext(cc))
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

// *********** cap with complex method **********
//
// interface testCapWithComplexMethod {
//   ComplexMethod(int64Value :int64, strValue :Text) -> (int64Value :int64, strValue :Text)
// }

const (
	testCapWithComplexMethod_InterfaceId                     = 2001
	testCapWithComplexMethod_ComplexMethod_MethodId MethodId = 30
)

type testCapWithComplexMethod CallFuture

var testCapWithComplexMethod_ComplexMethod_ParamsSize = capnpser.StructSize{DataSectionSize: 2, PointerSectionSize: 2}

type testCapWithComplexMethod_ComplexMethod_Params capnpser.Struct

func (s *testCapWithComplexMethod_ComplexMethod_Params) Int64Value() int64 {
	return (*capnpser.Struct)(s).Int64(1) // Different than result.
}

func (s *testCapWithComplexMethod_ComplexMethod_Params) StrValue() string {
	return (*capnpser.Struct)(s).String(1) // Different than result.
}

func (s *testCapWithComplexMethod_ComplexMethod_Params) UnsafeStrValue() string {
	return (*capnpser.Struct)(s).UnsafeString(1) // Different than result.
}

type testCapWithComplexMethod_ComplexMethod_ParamsBuilder capnpser.StructBuilder

func (b *testCapWithComplexMethod_ComplexMethod_ParamsBuilder) SetInt64Value(v int64) error {
	return (*capnpser.StructBuilder)(b).SetInt64(1, v)
}

func (b *testCapWithComplexMethod_ComplexMethod_ParamsBuilder) SetStrValue(v string) error {
	return (*capnpser.StructBuilder)(b).SetString(1, v)
}

var testCapWithComplexMethod_ComplexMethod_ResultSize = capnpser.StructSize{DataSectionSize: 2, PointerSectionSize: 2}

type testCapWithComplexMethod_ComplexMethod_Result capnpser.Struct

func (s *testCapWithComplexMethod_ComplexMethod_Result) Int64Value() int64 {
	return (*capnpser.Struct)(s).Int64(0) // Different than params.
}

func (s *testCapWithComplexMethod_ComplexMethod_Result) StrValue() string {
	return (*capnpser.Struct)(s).String(0) // Different than result.
}

func (s *testCapWithComplexMethod_ComplexMethod_Result) UnsafeStrValue() string {
	return (*capnpser.Struct)(s).UnsafeString(0) // Different than result.
}

type testCapWithComplexMethod_ComplexMethod_ResultBuilder capnpser.StructBuilder

func (b *testCapWithComplexMethod_ComplexMethod_ResultBuilder) SetInt64Value(v int64) error {
	return (*capnpser.StructBuilder)(b).SetInt64(0, v)
}

func (b *testCapWithComplexMethod_ComplexMethod_ResultBuilder) SetStrValue(v string) error {
	return (*capnpser.StructBuilder)(b).SetString(0, v)
}

type futureTestCapWithComplexMethod_ComplexMethod CallFuture

func (cf futureTestCapWithComplexMethod_ComplexMethod) Wait(ctx context.Context) (testCapWithComplexMethod_ComplexMethod_Result, ReturnResults, error) {
	return WaitReturnResultsStruct[testCapWithComplexMethod_ComplexMethod_Result](ctx, CallFuture(cf))
}

func (parent testCapWithComplexMethod) ComplexMethod(payloadSizeHint capnpser.WordCount) (testCapWithComplexMethod_ComplexMethod_ParamsBuilder, futureTestCapWithComplexMethod_ComplexMethod) {
	return RemoteCallWithStructParams[
		testCapWithComplexMethod_ComplexMethod_ParamsBuilder,
		futureTestCapWithComplexMethod_ComplexMethod](
		CallWithStructParamsSetup{
			ParentCap:       CallFuture(parent),
			InterfaceId:     testCapWithComplexMethod_InterfaceId,
			MethodId:        testCapWithComplexMethod_ComplexMethod_MethodId,
			ParamsSize:      testCapWithComplexMethod_ComplexMethod_ParamsSize,
			PayloadSizeHint: payloadSizeHint,
		})
}

type testCapWithComplexMethod_ComplexMethod_CallContext = StructReplyCallContext[testCapWithComplexMethod_ComplexMethod_ResultBuilder]

type testCapWithComplexMethodServer interface {
	ComplexMethod(testCapWithComplexMethod_ComplexMethod_Params, testCapWithComplexMethod_ComplexMethod_CallContext) error
}

func testCapWithComplexMethodServerAsHandler(s testCapWithComplexMethodServer) CallHandler {
	return CallHandlerFunc(func(ctx context.Context, cc *CallContext) error {
		iid, mid := cc.InterfaceId(), cc.MethodId()
		if iid != testCapWithComplexMethod_InterfaceId {
			return ErrUnexpectedInterfaceId{Want: testCapWithComplexMethod_InterfaceId, Got: iid}
		}
		switch mid {
		case testCapWithComplexMethod_ComplexMethod_MethodId:
			params, err := CallContextParamsStruct[testCapWithComplexMethod_ComplexMethod_Params](cc)
			if err != nil {
				return err
			}
			return s.ComplexMethod(
				params,
				MakeStructReplyCallContext[testCapWithComplexMethod_ComplexMethod_ResultBuilder](
					cc, testCapWithComplexMethod_ComplexMethod_ResultSize))
		default:
			return ErrUnknownMethodId{IId: iid, Mid: mid}
		}
	})
}
