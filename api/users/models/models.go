package models

import "fmt"

type User struct {
	ID     int64 `bun:",pk,autoincrement"`
	Name   string
	Emails []string
}

type Org struct {
	ID      int64 `bun:",pk,autoincrement"`
	OwnerID int64
	Owner   *User `bun:"rel:belongs-to,join:owner_id=id"`
}

type ResourceGroup struct {
	ID    int64 `bun:",pk,autoincrement"`
	OrgID int64
	Org   *Org `bun:"rel:belongs-to,join:org_id=id"`
	Name  string
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v>", u.ID, u.Name, u.Emails)
}
