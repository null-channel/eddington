package repositories

import (
	"context"

	"github.com/null-channel/eddington/api/infrastrucure"
	models "github.com/null-channel/eddington/api/users/models"
)

type OrgReposiotry struct {
	infrastrucure.Database
}

func (repository *OrgReposiotry) Seed() error {
	_, err := repository.DB().NewCreateTable().Model((*models.Org)(nil)).Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (repository *OrgReposiotry) GetOrgByOwnerId(ctx context.Context, userId string) ([]*models.Org, error) {
	var orgs []*models.Org
	err := repository.DB().NewSelect().
		Model(&orgs).
		Where("owner_id = ?", userId).
		Scan(ctx, &orgs)

	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (repository *OrgReposiotry) Save(ctx context.Context, org *models.Org) error {
	_, err := repository.DB().NewInsert().
		Model(org).
		On("CONFLICT (owner_id) DO UPDATE").
		Exec(ctx)

	return err
}
