package models

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type User struct {
	ID                int64 `bun:",pk"` // primary key, same as ory.
	Name              string
	Email             string
	NewsLetterConsent bool `bun:"newsletterConsent"`
}

type Org struct {
	ID             int64 `bun:",pk,autoincrement"`
	Name           string
	OwnerID        int64            `bun:"owner_id"`
	Owner          *User            `bun:"rel:belongs-to,join:owner_id=id"`
	ResourceGroups []*ResourceGroup `bun:"rel:has-many,join:id=org_id"`
}

type ResourceGroup struct {
	ID        int64        `bun:",pk,autoincrement"`
	OrgID     int64        `bun:"org_id"`
	Name      string       `bun:"name"`
	Resources []*Resources `bun:"rel:has-many,join:id=resource_group_id"`
}

type Resources struct {
	ID              int64 `bun:",pk,autoincrement"`
	CreatedAt       int64
	ResourceGroupID int64
	Type            string
}

func (u User) String() string {
	return fmt.Sprintf("User<%s %s>", u.ID, u.Email)
}

// UpdateUser godoc
// @Summary	Get user info for user id
func GetUserForId(id int64, userDb *bun.DB) (*User, error) {
	var user User
	err := userDb.NewSelect().
		Model(&user).
		Where("id = ?", id).
		Scan(context.Background(), &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetOrgByOwnerId(user_id int64, userDb *bun.DB) (*Org, error) {
	var org Org
	err := userDb.NewSelect().
		Model(&org).
		Where("owner_id = ?", user_id).
		Scan(context.Background(), &org)

	if err != nil {
		return nil, err
	}

	return &org, nil
}
