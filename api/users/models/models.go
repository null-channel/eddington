package models

import (
	"fmt"
	"time"
)

type User struct {
	ID                string `bun:",pk"` // primary key, same as ory.
	Name              string
	Email             string
	NewsLetterConsent bool `bun:"newsletterConsent"`
	DOB               time.Time
}

type Org struct {
	ID             int64 `bun:",pk,autoincrement"`
	Name           string
	OwnerID        string           `bun:"owner_id"`
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
