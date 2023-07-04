package controllers

import (
	"context"
	"errors"
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"

	usercon "github.com/null-channel/eddington/api/users/controllers"
	"github.com/null-channel/eddington/api/users/models"
)

//	@BasePath	/api/v1/

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

type ApplicationController struct {
	kube           dynamic.Interface
	userController *usercon.UserController
}

func NewApplicationController(kube dynamic.Interface, userService *usercon.UserController) *ApplicationController {
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	return &ApplicationController{
		kube:           kube,
		userController: userService,
	}
}

func (a ApplicationController) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/apps", a.AppPOST())
	routerGroup.GET("/apps/:id", a.AppGET())
}

type Application struct {
	Name          string `json:"name" binding:"required"`
	Image         string `json:"image" binding:"required"`
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
			ResourceGroup: c.PostForm("resourceGroup"),
		}
		//TODO: get user namespace

		userContext, err := a.userController.GetUserContext(context.Background(), 1)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}
		resourceGroup, err := getResourceGroupName(userContext.ResourceGroups, app.ResourceGroup)
		if err != nil {
			c.IndentedJSON(400, "Resource group not found")
		}

		namespace := userContext.Name + resourceGroup

		deployment := getApplication(app.Name, namespace, app.Image)
		_, err = a.kube.Resource(getDeploymentGVR()).Namespace(namespace).Apply(context.Background(), app.Name, deployment, v1.ApplyOptions{})
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}
	}
}

func getResourceGroupName(resourceGroups []*models.ResourceGroup, requested string) (string, error) {
	if requested == "" {
		return "default", nil
	}
	for _, group := range resourceGroups {
		if group.Name == requested {
			return group.Name, nil
		}
	}
	return "", errors.New("Resource group not found")
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
