package services

import (
	"context"

	"github.com/null-channel/eddington/api/users/models"
	repositories "github.com/null-channel/eddington/api/users/repositories"
)

type OrgService struct {
	OrgRepository *repositories.OrgReposiotry
}

func (service *OrgService) CreateOrg(org *models.Org, ctx context.Context) error {
	return service.OrgRepository.Save(ctx, org)
}

func (service *OrgService) GetOrgByOwnerId(ctx context.Context, ownerID string) ([]*models.Org, error) {
	return service.OrgRepository.GetOrgByOwnerId(ctx, ownerID)
}
