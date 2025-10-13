// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"fmt"
)

type futureString callFuture

func (fs futureString) wait(ctx context.Context) (string, error) {
	return castCallResultOrErr[string](waitResult(ctx, callFuture(fs)))
}

type testAPI callFuture

const testAPI_InterfaceID = 1000

type futureVoid callFuture

func (fv futureVoid) Wait(ctx context.Context) error {
	_, err := waitResult(ctx, callFuture(fv))
	return err
}

const testAPI_Void_CallID = 101
const testAPI_GetAnotherAPI_CallID = 102
const testAPI_GetUser_CallID = 103

func (api testAPI) VoidCall() futureVoid {
	return futureVoid(remoteCall(
		callFuture(api),
		callSetup{
			interfaceId: testAPI_InterfaceID,
			methodId:    testAPI_Void_CallID,
		},
	))
}

func (api testAPI) GetAnotherAPICap() testAPI {
	return testAPI(remoteCall(
		callFuture(api),
		callSetup{
			interfaceId: testAPI_InterfaceID,
			methodId:    testAPI_GetAnotherAPI_CallID,
		},
	))
}

func (api testAPI) GetUser(id string) testUser {
	return testUser(remoteCall(
		callFuture(api),
		callSetup{
			interfaceId: testAPI_InterfaceID,
			methodId:    testAPI_GetUser_CallID,
			paramsBuilder: func(*msgBuilder) error {
				_ = id
				return nil
			},
		},
	))
}

// Wait until this is resolved as a concrete, exported capability.
func (api testAPI) Wait(ctx context.Context) (capability, error) {
	return castCallResultOrErr[capability](waitResult(ctx, callFuture(api)))
}

func (api testAPI) WaitDiscardResult(ctx context.Context) error {
	_, err := waitResult(ctx, callFuture(api))
	return err
}

func testAPIAsBootstrap(bt bootstrapCap) testAPI {
	return testAPI(castBootstrap(bt))
}

type testUser callFuture

func (usr testUser) GetProfile() testUserProfile {
	return testUserProfile(remoteCall(
		callFuture(usr),
		callSetup{
			interfaceId: 1000,
			methodId:    11,
			paramsBuilder: func(*msgBuilder) error {
				return nil
			},
		},
	))
}

type testUserProfile callFuture

func (up testUserProfile) GetAvatarData() futureString {
	return futureString(remoteCall(
		callFuture(up),
		callSetup{
			interfaceId: 1000,
			methodId:    11,
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
	go waitResult(ctx, callFuture(user2)) // Dispatched fork parent after fork children.

	// _ = testUser(api).GetProfile() // Should not compile
}
