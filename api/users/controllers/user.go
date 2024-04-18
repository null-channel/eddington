package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/beevik/guid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	pb "github.com/null-channel/eddington/proto/user"

	"github.com/null-channel/eddington/api/core"
	"github.com/null-channel/eddington/api/users/models"
	"github.com/null-channel/eddington/api/users/types"
)

type MembersDatastore interface {
	CreateOrUpdateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateOrg(ctx context.Context, org *models.Org) error
	GetOrgByID(ctx context.Context, id int64) (*models.Org, error)
	UpdateOrg(ctx context.Context, org *models.Org) error
	GetOrgByOwnerId(ctx context.Context, ownerId string) ([]*models.Org, error)
	CreateResourceGroup(ctx context.Context, resourceGroup *models.ResourceGroup) error
	GetResourceGroupByID(ctx context.Context, id int64) (*models.ResourceGroup, error)
	UpdateResourceGroup(ctx context.Context, resourceGroup *models.ResourceGroup) error
	GetResourceGroupByOrgID(ctx context.Context, orgID *int64) (resGroups []*models.ResourceGroup, err error)
}

// Mux Controller to handel user routes
type UserController struct {
	pb.UnimplementedUserServiceServer

	membersDatastore MembersDatastore

	logger *zap.SugaredLogger
}

func New(
	logger *zap.Logger,
	membersDatastore MembersDatastore,
) (*UserController, error) {

	userServer := &UserController{
		membersDatastore: membersDatastore,
		logger:           logger.Sugar(),
	}

	return userServer, nil
}

func modelToUserContextRequest(org models.Org, ownerId string) *pb.GetUserContextReply {
	return &pb.GetUserContextReply{
		Org: &pb.Org{
			ID:             org.ID,
			Name:           org.Name,
			OwnerID:        ownerId,
			ResourceGroups: resourceGroupModelToProto(org.ResourceGroups),
		},
	}
}

func resourceGroupModelToProto(resourceGroups []*models.ResourceGroup) []*pb.ResourceGroup {
	var ret []*pb.ResourceGroup
	for _, resourceGroup := range resourceGroups {
		rg := &pb.ResourceGroup{
			Id:    resourceGroup.ID,
			OrgId: resourceGroup.OrgID,
			Name:  resourceGroup.Name,
		}

		for _, resource := range resourceGroup.Resources {
			rg.Resources = append(rg.Resources, &pb.Resource{
				Id:   resource.ID,
				Type: resource.Type,
			})
		}
		ret = append(ret, rg)
	}
	return ret
}

func (u *UserController) AddAllControllers(router *mux.Router) {
	router.HandleFunc("", u.UpsertUser).Methods("POST")
	router.HandleFunc("", u.GetUserId).Methods("GET")
}

// UpdateUser godoc
//
// @Summary	update an user
// @Schemes
// @Description	update a user
// @Tags			Users
// @Accept			json
// @Produce		json
// @Success		200	{string}	Helloworld
// @Router			/users/ [post]
func (u *UserController) UpsertUser(w http.ResponseWriter, r *http.Request) {
	// TODO:
	userID := r.Context().Value("user-id").(string)

	var userDTO types.NewUserRequest
	err := json.NewDecoder(r.Body).Decode(&userDTO)

	u.logger.Info("Upserting user: ", userDTO)

	if err != nil {
		u.logger.Error(err)
		http.Error(w, "Decode error! Please check your JSON formatting.", http.StatusBadRequest)
		return
	}

	if err := userDTO.Validate(); err != nil {
		errorMessage := types.ConstructErrorMeesages(err)
		core.ValidationErrors(w, errorMessage)
		return
	}

	_, isNewUser := u.membersDatastore.GetUserByEmail(r.Context(), userDTO.Email)

	// Create or update the user based on userDTO
	user := &models.User{
		ID:                userID,
		Name:              userDTO.Name,
		Email:             userDTO.Email,
		NewsLetterConsent: userDTO.NewsletterConsent,
		DOB:               userDTO.DOB,
	}

	err = u.membersDatastore.CreateOrUpdateUser(r.Context(), user)

	if err != nil {
		u.logger.Error(err)
		core.InternalErrorHandler(w)
		return
	}

	if isNewUser == nil {
		// User already exist in the system it's an update. We can't update the org and default resource group here.
		w.WriteHeader(http.StatusCreated)
		return
	}

	// Create a new org for the user
	org := &models.Org{
		Name:    userDTO.Name + guid.New().String(),
		OwnerID: user.ID,
	}

	err = u.membersDatastore.CreateOrg(r.Context(), org)

	if err != nil {
		u.logger.Error(err)
		core.InternalErrorHandler(w)
		return
	}

	// Create a new resource group for the org
	resourceGroup := &models.ResourceGroup{
		OrgID: org.ID,
		Name:  "Default",
	}

	err = u.membersDatastore.CreateResourceGroup(r.Context(), resourceGroup)

	if err != nil {
		u.logger.Error(err)
		core.InternalErrorHandler(w)
		return
	}

	// Should we return the user context here?

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// UpdateUser godoc
//
// @Summary	update an user
// @Schemes
// @Description	update a user
// @Tags			Users
// @Accept			json
// @Produce		json
// @Success		200	{string}	Helloworld
// @Router			/users/ [post]
func (u *UserController) GetUserId(w http.ResponseWriter, r *http.Request) {

	// Get the user ID from the context
	userID := r.Context().Value("user-id").(string)

	user, err := u.membersDatastore.GetUserByID(r.Context(), userID)

	if err != nil {
		u.logger.Error(err)
		core.InternalErrorHandler(w)
		return
	}

	// Get the user context
	org, err := u.GetUserContext(r.Context(), user.ID)

	if err != nil {
		u.logger.Error(err)
		core.InternalErrorHandler(w)
		return
	}

	// return the UserGetResponse
	userContext := &types.UserGetResponse{
		User: user,
		Org:  org,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userContext)

}

func (u *UserController) GetUserContext(ctx context.Context, userId string) (*models.Org, error) {
	// This assumes that the user is the owner. This is bad... but works for now.
	// This is probably not even going to be an indext column in the future.
	// Regrets future marek.

	orgs, err := u.membersDatastore.GetOrgByOwnerId(ctx, userId)

	if err != nil {
		return nil, err
	}

	var resGroups []*models.ResourceGroup
	resGroups, _ = u.membersDatastore.GetResourceGroupByOrgID(ctx, &orgs[0].ID)

	orgs[0].ResourceGroups = resGroups

	fmt.Println(orgs)

	return orgs[0], nil
}

func (u *UserController) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := u.membersDatastore.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
