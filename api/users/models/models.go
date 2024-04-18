package models

import (
	"fmt"
	"time"
)

type User struct {
	ID                string    `bun:",pk" json:"id"` // primary key, same as ory.
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	NewsLetterConsent bool      `bun:"newsletterConsent" json:"news_letter_consent"`
	DOB               time.Time `json:"dob"`
}

type Org struct {
	ID             int64            `bun:",pk,autoincrement" json:"id"`
	Name           string           `json:"name"`
	OwnerID        string           `bun:"owner_id" json:"owner_id"`
	Owner          *User            `bun:"rel:belongs-to,join:owner_id=id" json:"owner"`
	ResourceGroups []*ResourceGroup `bun:"rel:has-many,join:id=org_id" json:"resource_groups"`
}

type ResourceGroup struct {
	ID        int64        `bun:",pk,autoincrement" json:"id"`
	OrgID     int64        `bun:"org_id" json:"org_id"`
	Name      string       `bun:"name" json:"name"`
	Resources []*Resources `bun:"rel:has-many,join:id=resource_group_id" json:"resources"`
}

type Resources struct {
	ID              int64  `bun:",pk,autoincrement" json:"id"`
	CreatedAt       int64  `json:"created_at"`
	ResourceGroupID int64  `json:"resource_group_id"`
	Type            string `json:"type"`
}

func (u User) String() string {
	return fmt.Sprintf("User<%s %s>", u.ID, u.Email)
}
