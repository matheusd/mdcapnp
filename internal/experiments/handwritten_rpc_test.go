// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package experiments

import (
	"context"
	"fmt"
)

type futureString struct {
	_futureString struct{}
	fc            futureCap
}

func (fs futureString) wait(ctx context.Context) (string, error) {
	return waitResult[string](ctx, fs.fc)
}

type testAPI struct {
	_testAPI struct{}
	fc       futureCap
}

func (api testAPI) GetUser(id string) testUser {
	pb := func(*msgBuilder) error {
		_ = id
		return nil
	}
	return testUser{fc: remoteCall(api.fc, 1000, 10, pb)}
}

func testAPIAsBootstrap(bt bootstrapCap) testAPI {
	return testAPI{fc: bt.fc}
}

type testUser struct {
	_testUser struct{}
	fc        futureCap
}

func (usr testUser) GetProfile() testUserProfile {
	pb := func(*msgBuilder) error {
		return nil
	}
	return testUserProfile{fc: remoteCall(usr.fc, 1000, 11, pb)}
}

type testUserProfile struct {
	_testUserProfile struct{}
	fc               futureCap
}

func (up testUserProfile) GetAvatarData() futureString {
	pb := func(*msgBuilder) error {
		return nil
	}
	return futureString{fc: remoteCall(up.fc, 1000, 12, pb)}
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
}
