// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package experiments

type futureString future[string]

type testAPI future[capability]

func (api testAPI) GetUser(id string) testUser {
	return testUser(then[capability, testUser](future[capability](api)))
}

type testUser future[capability]

func (usr testUser) GetProfile() testUserProfile {
	return testUserProfile(then[capability, testUserProfile](future[capability](usr)))
}

type testUserProfile future[capability]

func (up testUserProfile) GetAvatarData() futureString {
	return futureString(then[capability, string](future[capability](up)))
}
