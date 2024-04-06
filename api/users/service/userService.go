package services

import (
	"context"
	"fmt"

	"github.com/null-channel/eddington/api/users/models"
	repositories "github.com/null-channel/eddington/api/users/repositories"
)

type UserService struct {
	UserRepository           repositories.IUserReposiotry
	OrgReposiotry            repositories.IOrgReposiotry
	ResourcesGroupRepository repositories.IResourcesGroupReposiotry
}

func (service *UserService) GetUserContext(ctx context.Context, userId string) (*models.Org, error) {
	// This assumes that the user is the owner. This is bad... but works for now.
	// This is probably not even going to be an indext column in the future.
	// Regrets future marek.

	orgs, err := service.OrgReposiotry.GetOrgByOwnerId(userId, ctx)

	if err != nil {
		return nil, err
	}

	if len(orgs) == 0 {
		return nil, fmt.Errorf("no orgs found for user")
	}

	var resGroups []*models.ResourceGroup
	resGroups, _ = service.ResourcesGroupRepository.GetResourceGroupByOrgID(&orgs[0].ID, ctx)

	orgs[0].ResourceGroups = resGroups

	fmt.Println(orgs)

	return orgs[0], nil
}

func (service *UserService) CreateOrUpdateUser(ctx context.Context, user *models.User) error {
	err := service.UserRepository.Save(user, ctx)

	if err != nil {
		return err
	}

	org := models.Org{
		Name:    user.Name,
		OwnerID: user.ID,
	}
	err = service.OrgReposiotry.Save(&org, ctx)
	if err != nil {
		return err
	}

	resourceGroup := models.ResourceGroup{
		OrgID: org.ID,
		Name:  "default",
	}
	err = service.ResourcesGroupRepository.Save(&resourceGroup, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (service *UserService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := service.UserRepository.GetUserByID(userID, ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}
