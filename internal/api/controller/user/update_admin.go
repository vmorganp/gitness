// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package user

import (
	"context"
	"fmt"
	"time"

	apiauth "github.com/harness/gitness/internal/api/auth"
	"github.com/harness/gitness/internal/api/usererror"
	"github.com/harness/gitness/internal/auth"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

type UpdateAdminInput struct {
	Admin bool `json:"admin"`
}

// UpdateAdmin updates the admin state of a user.
func (c *Controller) UpdateAdmin(ctx context.Context, session *auth.Session,
	userUID string, request *UpdateAdminInput) (*types.User, error) {
	user, err := findUserFromUID(ctx, c.principalStore, userUID)
	if err != nil {
		return nil, err
	}

	// Ensure principal has required permissions on parent.
	if err = apiauth.CheckUser(ctx, c.authorizer, session, user, enum.PermissionUserEditAdmin); err != nil {
		return nil, err
	}

	// Fail if the user being updated is the only admin in DB.
	if request.Admin == false && user.Admin == true {
		admUsrCount, err := c.principalStore.CountUsers(ctx, &types.UserFilter{Admin: true})
		if err != nil {
			return nil, fmt.Errorf("failed to check admin user count: %w", err)
		}

		if admUsrCount == 1 {
			return nil, usererror.BadRequest("system requires at least one admin user")
		}

		return user, nil
	}
	user.Admin = request.Admin
	user.Updated = time.Now().UnixMilli()

	err = c.principalStore.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
