package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/null-channel/eddington/api/users/models"
	pb "github.com/null-channel/eddington/proto/user"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type UserController struct {
	pb.UnimplementedUserServiceServer
	database *bun.DB
}

func New() (*UserController, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	db := bun.NewDB(sqldb, sqlitedialect.New())
	if err != nil {
		panic(err)
	}
	_, err = db.NewCreateTable().Model((*models.User)(nil)).Exec(context.Background())
	if err != nil {
		panic(err)
	}
	_, err = db.NewCreateTable().Model((*models.ResourceGroup)(nil)).Exec(context.Background())
	if err != nil {
		panic(err)
	}
	_, err = db.NewCreateTable().Model((*models.Org)(nil)).Exec(context.Background())
	if err != nil {
		panic(err)
	}

	userServer := &UserController{database: db}

	return userServer, nil
}

func (u *UserController) GetUserContext(ctx context.Context, userId int64) (*models.Org, error) {

	// This assumes that the user is the owner. This is bad... but works for now.
	// This is probably not even going to be an indext column in the future.
	// Regrets future marek.
	var orgs []models.Org
	err := u.database.NewSelect().Model(orgs).Where(fmt.Sprintf("owner_id = %d", userId)).Scan(ctx, orgs)

	if err != nil {
		return nil, err
	}

	return &orgs[0], nil
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
	router.HandleFunc("", u.UpdateUser).Methods("POST")
	router.HandleFunc("/id", u.GetUserId).Methods("GET")
}

//	@BasePath	/api/v1/

// CreateUser godoc
//
//	@Summary	Create an user
//	@Schemes
//	@Description	create a user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Helloworld
//	@Router			/users/ [post]
func (u *UserController) CreateUser(user models.User) (int, error) {

	res, err := u.database.NewInsert().Model(&user).Exec(context.Background())

	if err != nil {
		return http.StatusInternalServerError, err
	}

	ownerId, err := res.LastInsertId()

	if err != nil {
		return http.StatusInternalServerError, err
	}

	org := models.Org{
		OwnerID: ownerId,
	}
	res, err = u.database.NewInsert().Model(&org).Exec(context.Background())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	orgId, err := res.LastInsertId()
	resourceGroup := models.ResourceGroup{
		OrgID: orgId,
		Name:  "default",
	}
	_, err = u.database.NewInsert().Model(&resourceGroup).Exec(context.Background())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
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
func (u *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	w.WriteHeader(http.StatusNotImplemented)
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

