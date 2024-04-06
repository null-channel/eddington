package services

import (
	"context"

	"github.com/null-channel/eddington/api/users/models"
	repositories "github.com/null-channel/eddington/api/users/repositories"
)

type ResourcesGroupService struct {
	ResourcesGroupRepository *repositories.ResourcesGroupReposiotry
}

func (resourcesGroupService *ResourcesGroupService) CreateResourcesGroup(ctx context.Context, resourcesGroup *models.ResourceGroup) error {
	return resourcesGroupService.ResourcesGroupRepository.Save(ctx, resourcesGroup)
}
