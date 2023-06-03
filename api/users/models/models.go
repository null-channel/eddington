package models

import "fmt"

type User struct {
	ID     int64 `bun:",pk,autoincrement"`
	Name   string
	Emails []string
}

type Org struct {
	ID     int64 `bun:",pk,autoincrement"`
	PlanID int64
}

type UsersToOrgs struct {
	UserID int64 `bun:",pk"`
	User   *User `bun:"rel:belongs-to,join:user_id=id"`
	OrgID  int64 `bun:",pk"`
	Org    *Org  `bun:"rel:belongs-to,join:org_id=id"`
}

type Application struct {
	ID      int64 `bun:",pk,autoincrement"`
	OwnerID int64
	Owner   *Org `bun:"rel:belongs-to,join:owner_id=id`
}

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %v>", u.ID, u.Name, u.Emails)
}

type Story struct {
	ID       int64 `bun:",pk,autoincrement"`
	Title    string
	AuthorID int64
	Author   *User `bun:"rel:belongs-to,join:author_id=id"`
}
