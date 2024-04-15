package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/beevik/guid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"

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
	GetOrgByOwnerId(ctx context.Context, ownerId string) (*models.Org, error)
	CreateResourceGroup(ctx context.Context, resourceGroup *models.ResourceGroup) error
	GetResourceGroupByID(ctx context.Context, id int64) (*models.ResourceGroup, error)
	UpdateResourceGroup(ctx context.Context, resourceGroup *models.ResourceGroup) error
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

	var userDTO types.NewUserRequest
	err := json.NewDecoder(r.Body).Decode(&userDTO)

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

	isUserExist, err := u.membersDatastore.GetUserByEmail(r.Context(), userDTO.Email)

	if err != nil {
		u.logger.Error(err)
		core.InternalErrorHandler(w)
		return
	}

	// Create or update the user based on userDTO
	user := &models.User{
		ID:                userDTO.ID,
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

	if isUserExist != nil {
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
	// TODO: implement
	fmt.Println("User ID: " + r.Context().Value("user-id").(string))
	w.WriteHeader(http.StatusNotImplemented)
}
