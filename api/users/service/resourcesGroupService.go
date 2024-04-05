package services

import (
	"context"

	"github.com/null-channel/eddington/api/users/models"
	repositories "github.com/null-channel/eddington/api/users/repositories"
)

type IResourcesGroupService interface {
	CreateResourcesGroup(org *models.ResourceGroup, ctx context.Context) error
}

type ResourcesGroupService struct {
	ResourcesGroupRepository repositories.IResourcesGroupReposiotry
}

func (resourcesGroupService *ResourcesGroupService) CreateResourcesGroup(resourcesGroup *models.ResourceGroup, ctx context.Context) error {
	return resourcesGroupService.ResourcesGroupRepository.Save(resourcesGroup, ctx)
}
