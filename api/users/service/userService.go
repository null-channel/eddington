package services

import (
	"context"
	"fmt"

	"github.com/null-channel/eddington/api/users/models"
	repositories "github.com/null-channel/eddington/api/users/repositories"
)

type UserService struct {
	UserRepository           *repositories.UserRepository
	OrgReposiotry            *repositories.OrgReposiotry
	ResourcesGroupRepository *repositories.ResourcesGroupReposiotry
}

func (service *UserService) GetUserContext(ctx context.Context, userId string) (*models.Org, error) {
	// This assumes that the user is the owner. This is bad... but works for now.
	// This is probably not even going to be an indext column in the future.
	// Regrets future marek.

	orgs, err := service.OrgReposiotry.GetOrgByOwnerId(ctx, userId)

	if err != nil {
		return nil, err
	}

	if len(orgs) == 0 {
		return nil, fmt.Errorf("no orgs found for user")
	}

	var resGroups []*models.ResourceGroup
	resGroups, _ = service.ResourcesGroupRepository.GetResourceGroupByOrgID(ctx, &orgs[0].ID)

	orgs[0].ResourceGroups = resGroups

	fmt.Println(orgs)

	return orgs[0], nil
}

func (service *UserService) CreateOrUpdateUser(ctx context.Context, user *models.User) error {
	err := service.UserRepository.Save(ctx, user)

	if err != nil {
		return err
	}

	org := models.Org{
		Name:    user.Name,
		OwnerID: user.ID,
	}
	err = service.OrgReposiotry.Save(ctx, &org)
	if err != nil {
		return err
	}

	resourceGroup := models.ResourceGroup{
		OrgID: org.ID,
		Name:  "default",
	}
	err = service.ResourcesGroupRepository.Save(ctx, &resourceGroup)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := service.UserRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
