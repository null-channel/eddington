package types

type CreateUserRequest struct {
	UserId int64      `json:"userId"`
	Traits UserTraits `json:"traits"`
}

type UserTraits struct {
	Email             string `json:"email"`
	Name              string `json:"name"`
	NewsLetterConsent bool   `json:"newsletterConsent"`
}
