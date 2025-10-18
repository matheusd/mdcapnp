// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"fmt"

	types "matheusd.com/mdcapnp/capnprpc/types"
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

func (api testAPI) VoidCall() futureVoid {
	return futureVoid(RemoteCall(
		CallFuture(api),
		CallSetup{
			InterfaceId: testAPI_InterfaceID,
			MethodId:    testAPI_Void_CallID,
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
			ParamsBuilder: func(types.MessageBuilder) error {
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
			ParamsBuilder: func(types.MessageBuilder) error {
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
