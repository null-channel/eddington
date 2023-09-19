package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/null-channel/eddington/api/users/models"
	pb "github.com/null-channel/eddington/proto/user"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"go.uber.org/zap"
)

// Mux Controller to handel user routes
type UserController struct {
	pb.UnimplementedUserServiceServer
	database *bun.DB
	logger   *zap.SugaredLogger
}

func New(logger *zap.Logger) (*UserController, error) {
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

	userServer := &UserController{database: db, logger: logger.Sugar()}

	return userServer, nil
}

func (u *UserController) GetUserContext(ctx context.Context, userId int64) (*models.Org, error) {

	// This assumes that the user is the owner. This is bad... but works for now.
	// This is probably not even going to be an indext column in the future.
	// Regrets future marek.
	var orgs []models.Org
	err := u.database.NewSelect().
		Model(&orgs).
		Where("owner_id = ?", userId).
		Scan(ctx, &orgs)

	if err != nil {
		u.logger.Errorw("Error getting user from database",
			"error", err)
		return nil, err
	}

	var resGroups []*models.ResourceGroup
	err = u.database.NewSelect().
		Model(&resGroups).
		Where("org_id = ?", &orgs[0].ID).
		Scan(ctx)

	orgs[0].ResourceGroups = resGroups

	fmt.Println(orgs)

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

func (u *UserController) CreateUser(user models.User) (int, error) {

	res, err := u.database.NewInsert().Model(&user).Exec(context.Background())

	if err != nil {
		return http.StatusInternalServerError, err
	}

	ownerId, err := res.LastInsertId()

	u.logger.Infow("New user created",
		"userId:", ownerId)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	orgAsString := strconv.FormatInt(ownerId, 10)
	org := models.Org{
		Name:    "default-" + orgAsString,
		OwnerID: ownerId,
	}
	res, err = u.database.NewInsert().Model(&org).Exec(context.Background())
	if err != nil {
		return http.StatusInternalServerError, err
	}

	u.logger.Infow("Org created", "orgId:", org.ID)

	resourceGroup := models.ResourceGroup{
		OrgID: org.ID,
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
	u.CreateUser(models.User{
		ID:   1234,
		Name: "Marek",
		Traits: &models.Traits{
			Emails:            []string{"none@none.com"},
			NewsLetterConsent: true,
		},
	})
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
