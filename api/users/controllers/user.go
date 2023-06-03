package controllers

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/null-channel/eddington/api/users/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

type UserController struct {
	database *bun.DB
}

func New() UserController {
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

	return UserController{database: db}
}

func (u *UserController) AddAllControllers(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/users", u.CreateUser())
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
func (u *UserController) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create new user in database

		user := models.User{
			Name:   c.PostForm("name"),
			Emails: []string{c.PostForm("email")},
		}

		res, err := u.database.NewInsert().Model(&user).Exec(context.Background())

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ownerId, err := res.LastInsertId()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		org := models.Org{
			OwnerID: ownerId,
		}
		res, err = u.database.NewInsert().Model(&org).Exec(context.Background())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		orgId, err := res.LastInsertId()
		resourceGroup := models.ResourceGroup{
			OrgID: orgId,
			Name:  "default",
		}
		_, err = u.database.NewInsert().Model(&resourceGroup).Exec(context.Background())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "User created successfully!"})
	}
}
