package types

type CreateUserRequest struct {
	Name   string `json:"name"`
	Emails string `json:"emails"`
}
