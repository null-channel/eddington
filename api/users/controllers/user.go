package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	pb "github.com/null-channel/eddington/proto/user"

	"github.com/null-channel/eddington/api/core"
	"github.com/null-channel/eddington/api/users/models"
	services "github.com/null-channel/eddington/api/users/service"
	"github.com/null-channel/eddington/api/users/types"
)

// Mux Controller to handel user routes
type UserController struct {
	pb.UnimplementedUserServiceServer

	userService *services.UserService
	// orgService            services.IOrgService
	// resourcesGroupService services.IResourcesGroupService
	logger *zap.SugaredLogger
}

func New(
	logger *zap.Logger,
	userService *services.UserService,
	// orgService services.IOrgService,
	// resourcesGroupService services.IResourcesGroupService,
) (*UserController, error) {

	userServer := &UserController{
		userService: userService,
		// orgService:            orgService,
		// resourcesGroupService: resourcesGroupService,
		logger: logger.Sugar(),
	}

	return userServer, nil
}

func modelToUserContextRequest(org models.Org, ownerId int64) *pb.GetUserContextReply {
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

	// Create or update the user based on userDTO
	user := &models.User{
		ID:                userDTO.ID,
		Name:              userDTO.Name,
		Email:             userDTO.Email,
		NewsLetterConsent: userDTO.NewsletterConsent,
		DOB:               userDTO.DOB,
	}
	err = u.userService.CreateOrUpdateUser(r.Context(), user)

	if err != nil {
		u.logger.Error(err)
		core.InternalErrorHandler(w)
		return
	}

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
