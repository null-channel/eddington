package models

import "fmt"

type User struct {
	ID     int64 `bun:",pk,autoincrement"`
	Name   string
	Emails []string
}

type Org struct {
	ID             int64 `bun:",pk,autoincrement"`
	Name           string
	OwnerID        int64
	Owner          *User            `bun:"rel:belongs-to,join:owner_id=id"`
	ResourceGroups []*ResourceGroup `bun:"rel:has-many,join:id=org_id"`
}

type ResourceGroup struct {
	ID        int64 `bun:",pk,autoincrement"`
	OrgID     int64
	Name      string
	Resources []*Resources `bun:"rel:has-many,join:id=resource_group_id"`
}

type Resources struct {
	ID              int64 `bun:",pk,autoincrement"`
	CreatedAt       int64
	ResourceGroupID int64
	Type            string
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v>", u.ID, u.Name, u.Emails)
}
