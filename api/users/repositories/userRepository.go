package repositories

import (
	"context"

	"github.com/null-channel/eddington/api/infrastrucure"
	models "github.com/null-channel/eddington/api/users/models"
)

type UserRepository struct {
	infrastrucure.Database
}

func (repository *UserRepository) Seed() error {
	_, err := repository.DB().NewCreateTable().Model((*models.User)(nil)).Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (repository *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := repository.DB().NewSelect().
		Model(&user).
		Where("id = ?", id).
		Scan(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repository UserRepository) Save(ctx context.Context, user *models.User) error {
	_, err := repository.DB().NewInsert().
		Model(user).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name, email = EXCLUDED.email, newsletterConsent = EXCLUDED.newsletterConsent, dob = EXCLUDED.dob").
		Exec(ctx)

	if err != nil {
		return err
	}
	return nil
}
