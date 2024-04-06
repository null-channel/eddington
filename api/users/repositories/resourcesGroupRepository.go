package repositories

import (
	"context"

	"github.com/null-channel/eddington/api/infrastrucure"
	models "github.com/null-channel/eddington/api/users/models"
)

type ResourcesGroupReposiotry struct {
	infrastrucure.Database
}

func (repository *ResourcesGroupReposiotry) Seed() error {
	_, err := repository.DB().NewCreateTable().Model((*models.ResourceGroup)(nil)).Exec(context.Background())
	return err
}

func (repository *ResourcesGroupReposiotry) GetResourceGroupByOrgID(ctx context.Context, orgID *int64) (resGroups []*models.ResourceGroup, err error) {

	err = repository.DB().NewSelect().
		Model(&resGroups).
		Where("org_id = ?", orgID).
		Scan(ctx)

	return
}

func (repository *ResourcesGroupReposiotry) Save(ctx context.Context, resourcesGroup *models.ResourceGroup) error {
	_, err := repository.DB().NewInsert().Model(&resourcesGroup).Exec(ctx)
	return err
}
