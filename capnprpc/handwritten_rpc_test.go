// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"fmt"

	types "matheusd.com/mdcapnp/capnprpc/types"
	"matheusd.com/mdcapnp/capnpser"
)

type futureString CallFuture

func (fs futureString) wait(ctx context.Context) (string, error) {
	return CastCallResultOrErr[string](WaitResult(ctx, CallFuture(fs)))
}

type testAPI CallFuture

const testAPI_InterfaceID = 1000

type futureVoid CallFuture

func (fv futureVoid) Wait(ctx context.Context) error {
	_, err := WaitResult(ctx, CallFuture(fv))
	return err
}

const testAPI_Void_CallID = 101
const testAPI_GetAnotherAPI_CallID = 102
const testAPI_GetUser_CallID = 103
const testAPI_Add_CallID = 104

func (api testAPI) VoidCall() futureVoid {
	return futureVoid(RemoteCall(
		CallFuture(api),
		CallSetup{
			InterfaceId: testAPI_InterfaceID,
			MethodId:    testAPI_Void_CallID,
		},
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
	return CastCallResultOrErr[int64](WaitResult(ctx, CallFuture(fut)))
}

func (api testAPI) Add(a int64, b int64) futureAddResult {
	return futureAddResult(RemoteCall(
		CallFuture(api),
		CallSetup{
			InterfaceId: testAPI_InterfaceID,
			MethodId:    testAPI_Add_CallID,
			ParamsBuilder: func(payload types.PayloadBuilder) error {
				req, err := NewCallParamsStruct[addRequestBuilder](payload, addRequestSize)
				if err != nil {
					return err
				}
				req.SetA(a)
				req.SetB(b)
				return nil
			},
			ResultsParser: func(p types.Payload) (any, error) {
				res, err := ResultsStruct[addResponse](p)
				if err != nil {
					return nil, err
				}
				return res.C(), nil
			},
		},
	))
}

type futureAddAltResult CallFuture

func (fut futureAddAltResult) wait(ctx context.Context) error {
	_, err := WaitResult(ctx, CallFuture(fut))
	return err
}

func (api testAPI) AddAlt(a int64, b int64, c *int64) futureAddAltResult {
	return futureAddAltResult(RemoteCall(
		CallFuture(api),
		CallSetup{
			InterfaceId: testAPI_InterfaceID,
			MethodId:    testAPI_Add_CallID,
			ParamsBuilder: func(payload types.PayloadBuilder) error {
				req, err := NewCallParamsStruct[addRequestBuilder](payload, addRequestSize)
				if err != nil {
					return err
				}
				req.SetA(a)
				req.SetB(b)
				return nil
			},
			ResultsParser: func(p types.Payload) (any, error) {
				res, err := ResultsStruct[addResponse](p)
				if err != nil {
					return nil, err
				}
				*c = res.C()
				return c, nil
			},
		},
	))
}

type futureAddAlt2Result CallFuture

func (fut futureAddAlt2Result) wait(ctx context.Context) (res int64, err error) {
	// var resAny capnpser.AnyPointerBuilder
	// resAny, err = CastCallResultOrErr[capnpser.AnyPointerBuilder](WaitResult(ctx, CallFuture(fut)))
	var resMb *capnpser.MessageBuilder
	resMb, err = CastCallResultOrErr[*capnpser.MessageBuilder](WaitResult(ctx, CallFuture(fut)))
	if err != nil {
		return
	}
	resAnyReader := resMb.MessageReader()
	resStruct, err := resAnyReader.GetRoot()
	resAdd := addResponse(resStruct)
	res = resAdd.C()
	// fut.vat.mbp.put(resAny.MsgBuilder())
	fut.vat.mbp.put(resMb)
	return
}

func (api testAPI) AddAlt2(a int64, b int64) futureAddAlt2Result {
	callb := api.vat.GetCallMessageBuilder(addRequestSize.TotalSize())

	call, _ := callb.mb.NewCall()
	callb.builder = capnpser.StructBuilder(call)
	_ = call.SetInterfaceId(testAPI_InterfaceID)
	_ = call.SetMethodId(testAPI_Add_CallID)
	payload, _ := call.NewParams()
	req, err := NewCallParamsStruct[addRequestBuilder](payload, addRequestSize)
	if err != nil {
		panic(err)
	}
	req.SetA(a)
	req.SetB(b)

	return futureAddAlt2Result(RemoteCall(
		CallFuture(api),
		CallSetup{
			InterfaceId:       testAPI_InterfaceID,
			MethodId:          testAPI_Add_CallID,
			callOutMsg:        callb,
			copyReturnResults: true,
		},
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

func (api testAPI) GetUser(id string) testUser {
	return testUser(RemoteCall(
		CallFuture(api),
		CallSetup{
			InterfaceId: testAPI_InterfaceID,
			MethodId:    testAPI_GetUser_CallID,
			ParamsBuilder: func(types.PayloadBuilder) error {
				_ = id
				return nil
			},
		},
	))
}

// Wait until this is resolved as a concrete, exported capability.
func (api testAPI) Wait(ctx context.Context) (capability, error) {
	return CastCallResultOrErr[capability](WaitResult(ctx, CallFuture(api)))
}

func (api testAPI) WaitDiscardResult(ctx context.Context) error {
	_, err := WaitResult(ctx, CallFuture(api))
	return err
}

func testAPIAsBootstrap(bt BootstrapFuture) testAPI {
	return testAPI(bt)
}

type testUser CallFuture

func (usr testUser) GetProfile() testUserProfile {
	return testUserProfile(RemoteCall(
		CallFuture(usr),
		CallSetup{
			InterfaceId: 1000,
			MethodId:    11,
			ParamsBuilder: func(types.PayloadBuilder) error {
				return nil
			},
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

func example01() {
	var v *Vat
	var c conn
	rc := v.RunConn(c)
	boot := rc.Bootstrap()
	api := testAPIAsBootstrap(boot)
	user := api.GetUser("1000")
	prof := user.GetProfile()
	avatar := prof.GetAvatarData()
	avatarData, err := avatar.wait(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(avatarData)

	ctx := context.Background()

	fmt.Println(api.
		GetUser("10000").
		GetProfile().
		GetAvatarData().
		wait(context.Background()))

	_ = user.GetProfile() // Forked prior pipeline.

	user2 := api.GetUser("1000")
	prof2 := user2.GetProfile()
	prof2_2 := user2.GetProfile()        // Fork
	go prof2_2.GetAvatarData().wait(ctx) // Dispatched fork before original.
	go prof2.GetAvatarData().wait(ctx)
	go WaitResult(ctx, CallFuture(user2)) // Dispatched fork parent after fork children.

	// _ = testUser(api).GetProfile() // Should not compile
}
