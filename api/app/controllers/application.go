package controllers

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"

	appmodels "github.com/null-channel/eddington/api/app/models"
	usercon "github.com/null-channel/eddington/api/users/controllers"
	"github.com/null-channel/eddington/api/users/models"
	pb "github.com/null-channel/eddington/proto/container"
)

//	@BasePath	/api/v1/

type ApplicationController struct {
	kube                   dynamic.Interface
	userController         *usercon.UserController
	database               *bun.DB
	containerServiceClient pb.ContainerServiceClient
}

func NewApplicationController(kube dynamic.Interface, userService *usercon.UserController, containerBuildingService pb.ContainerServiceClient) *ApplicationController {
	// Set up a connection to the server.
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	db := bun.NewDB(sqldb, sqlitedialect.New())
	if err != nil {
		panic(err)
	}
	_, err = db.NewCreateTable().Model((*appmodels.NullApplication)(nil)).Exec(context.Background())
	if err != nil {
		panic(err)
	}
	_, err = db.NewCreateTable().Model((*appmodels.NullApplicationService)(nil)).Exec(context.Background())
	if err != nil {
		panic(err)
	}

	return &ApplicationController{
		kube:                   kube,
		userController:         userService,
		database:               db,
		containerServiceClient: containerBuildingService,
	}
}

func (a ApplicationController) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/apps", a.AppPOST())
	routerGroup.GET("/apps/:id", a.AppGET())
}

type Application struct {
	Name          string `json:"name" binding:"required"`
	Image         string `json:"image"`
	GitRepo       string `json:"gitRepo"`
	RepoType      string `json:"repoType"`
	ResourceGroup string `json:"resourceGroup"`
}

// AppPOST godoc
//
//	@Summary	Create an application
//	@Schemes
//	@Description	create an application
//	@Tags			example
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Helloworld
//	@Router			/apps/ [post]
func (a ApplicationController) AppPOST() gin.HandlerFunc {
	return func(c *gin.Context) {
		app := Application{
			Name:          c.PostForm("name"),
			Image:         c.PostForm("image"),
			GitRepo:       c.PostForm("gitRepo"),
			RepoType:      c.PostForm("repoType"),
			ResourceGroup: c.PostForm("resourceGroup"),
		}

		// get user namespace
		userContext, err := a.userController.GetUserContext(context.Background(), 1)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}
		resourceGroup, rgId, err := getResourceGroupName(userContext.ResourceGroups, app.ResourceGroup)
		if err != nil {
			c.IndentedJSON(400, "Resource group not found")
		}

		//TODO: Build container image if given a repo

		ret, err := a.containerServiceClient.CreateContainer(c.Copy().Request.Context(), &pb.CreateContainerRequest{
			RepoURL:    app.GitRepo,
			Type:       app.RepoType,
			CustomerID: userContext.Owner.ID,
		})

		namespace := userContext.Name + resourceGroup
		nullApplication := getNullApplication(app, userContext, rgId, namespace, ret.BuildID)

		if err != nil {
			c.IndentedJSON(500, "Internal server error")
			return
		}

		//TODO: Save to database!
		_, err = a.database.NewInsert().Model(nullApplication).Exec(context.Background())

		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}

		go func() {
			keepChecking := true
			var status *pb.ContainerImageStatusReply
			for keepChecking {
				status, err = a.containerServiceClient.ImageStatus(context.Background(), &pb.BuildRequest{Id: ret.BuildID})
				if err != nil {
					panic("checking container build status failed")
				}

				if status.Status == pb.ContainerStatus_BUILT {
					keepChecking = false
				}
				time.Sleep(20 * time.Second)
			}

			deployment := getApplication(*nullApplication)
			_, err = a.kube.Resource(getDeploymentGVR()).Namespace(namespace).Apply(context.Background(), app.Name, deployment, v1.ApplyOptions{})
			if err != nil {
				c.IndentedJSON(500, "Internal server error")
			}
		}()
	}
}

func getNullApplication(app Application, org *models.Org, resourceGroupId int64, namespace string, buildID string) *appmodels.NullApplication {
	return &appmodels.NullApplication{
		OrgID:           org.ID,
		Name:            app.Name,
		ResourceGroupID: resourceGroupId,
		Namespace:       namespace,
		NullApplicationService: []*appmodels.NullApplicationService{
			{
				Type:    appmodels.ContainerImage,
				GitRepo: app.GitRepo,
				Name:    app.Name,
				Image:   app.Image,
				Cpu:     "100m",
				Memory:  "100Mi",
				Storage: "1Gi",
				BuildID: buildID,
			},
		},
	}
}

func getResourceGroupName(resourceGroups []*models.ResourceGroup, requested string) (string, int64, error) {
	if requested == "" {
		requested = "default"
	}
	for _, group := range resourceGroups {
		if group.Name == requested {
			return group.Name, group.ID, nil
		}
	}
	return "", 0, errors.New("Resource group not found")
}

func getDeploymentGVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
}

// AppGET godoc
//
//	@Summary	Get all applications created by the user
//	@Schemes
//	@Description	Get all applications created by the user
//	@Tags			example
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Helloworld
//	@Router			/apps/ [get]
func (a ApplicationController) AppGET() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.IndentedJSON(501, "Not implemented yet")
	}
}
