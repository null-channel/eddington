package models

import (
	"fmt"

	"github.com/null-channel/eddington/api/users/types"
)

type User struct {
	ID     string  `bun:",pk"` // primary key, same as ory.
	Traits *Traits `bun:"rel:has-one,join:id=user_id"`
}

type Traits struct {
	ID                int64    `bun:",pk,autoincrement"`
	Emails            []string `bun:"email"`
	Name              string   `bun:"name"`
	NewsLetterConsent bool     `bun:"newsletterConsent"`
	UserID            string   `bun:"user_id"`
}

type Org struct {
	ID             int64 `bun:",pk,autoincrement"`
	Name           string
	OwnerID        int64
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
	return fmt.Sprintf("User<%s %s %s>", u.ID, u.Traits.Name, u.Traits.Emails)
}

func CreateUserRequestToDBModel(createUserRequest types.CreateUserRequest) User {
	return User{
		ID: createUserRequest.UserId,
		Traits: &Traits{
			Emails:            []string{createUserRequest.Traits.Email},
			Name:              createUserRequest.Traits.Name,
			NewsLetterConsent: createUserRequest.Traits.NewsLetterConsent,
		},
	}
}
