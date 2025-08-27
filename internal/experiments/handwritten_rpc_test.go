// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package experiments

import (
	"context"
	"fmt"
)

type futureString futureCap[string]

func (fs futureString) wait(ctx context.Context) (string, error) {
	return waitResult(ctx, futureCap[string](fs))
}

type testAPICap struct{}
type testAPI futureCap[testAPICap]

func (api testAPI) GetUser(id string) testUser {
	pb := func(*msgBuilder) error {
		_ = id
		return nil
	}
	return testUser(remoteCall[testAPICap, testUserCap](futureCap[testAPICap](api), 1000, 10, pb))
}

func testAPIAsBootstrap(bt bootstrapCap) testAPI {
	return testAPI(castBootstrap[testAPICap](bt))
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
	var v *vat
	var c *conn
	rc := v.RunConn(c)
	boot := rc.bootstrap()
	api := testAPIAsBootstrap(boot)
	user := api.GetUser("1000")
	prof := user.GetProfile()
	avatar := prof.GetAvatarData()
	avatarData, err := avatar.wait(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(avatarData)

	fmt.Println(api.
		GetUser("10000").
		GetProfile().
		GetAvatarData().
		wait(context.Background()))

	// _ = testUser(api).GetProfile() // Should not compile
}
