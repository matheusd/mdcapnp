// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"context"
	"fmt"
)

type futureString futureCap[string]

func (fs futureString) wait(ctx context.Context) (string, error) {
	return waitResult(ctx, futureCap[string](fs))
}

type testAPI futureCap[capability]

const testAPI_InterfaceID = 1000

type futureVoid futureCap[struct{}]

func (fv futureVoid) Wait(ctx context.Context) error {
	_, err := waitResult(ctx, futureCap[struct{}](fv))
	return err
}

const testAPI_Void_CallID = 101
const testAPI_GetAnotherAPI_CallID = 102

func (api testAPI) VoidCall() futureVoid {
	return futureVoid(remoteCall[capability, struct{}](futureCap[capability](api), testAPI_InterfaceID, testAPI_Void_CallID, nil))
}

func (api testAPI) GetAnotherAPI() testAPI {
	return testAPI(remoteCall[capability, capability](futureCap[capability](api), testAPI_InterfaceID, testAPI_GetAnotherAPI_CallID, nil))
}

func (api testAPI) GetUser(id string) testUser {
	pb := func(*msgBuilder) error {
		_ = id
		return nil
	}
	return testUser(remoteCall[capability, testUserCap](futureCap[capability](api), 1000, 10, pb))
}

// Wait until this is resolved as a concrete, exported capability.
func (api testAPI) Wait(ctx context.Context) error {
	_, err := waitResult(ctx, futureCap[capability](api))
	return err
}

func testAPIAsBootstrap(bt bootstrapCap) testAPI {
	return testAPI(castBootstrap[capability](bt))
}

type testUserCap struct{}
type testUser futureCap[testUserCap]

func (usr testUser) GetProfile() testUserProfile {
	pb := func(*msgBuilder) error {
		return nil
	}
	return testUserProfile(remoteCall[testUserCap, testUserProfileCap](futureCap[testUserCap](usr), 1000, 11, pb))
}

type testUserProfileCap struct{}
type testUserProfile futureCap[testUserProfileCap]

func (up testUserProfile) GetAvatarData() futureString {
	pb := func(*msgBuilder) error {
		return nil
	}
	return futureString(remoteCall[testUserProfileCap, string](futureCap[testUserProfileCap](up), 1000, 11, pb))
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
	go waitResult(ctx, futureCap[testUserCap](user2)) // Dispatched fork parent after fork children.

	// _ = testUser(api).GetProfile() // Should not compile
}
