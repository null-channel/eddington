package repositories

import (
	"context"

	"github.com/null-channel/eddington/api/infrastrucure"
	models "github.com/null-channel/eddington/api/users/models"
)

type IUserReposiotry interface {
	Seedable
	GetUserByID(id string, ctx context.Context) (*models.User, error)
	Save(user *models.User, ctx context.Context) error
}

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

func (repository *UserRepository) GetUserByID(id string, ctx context.Context) (*models.User, error) {
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

func (repository UserRepository) Save(user *models.User, ctx context.Context) error {
	_, err := repository.DB().NewInsert().
		Model(&user).
		On("CONFLICT (id) DO UPDATE").
		Exec(ctx)

	if err != nil {
		return err
	}
	return nil
}
