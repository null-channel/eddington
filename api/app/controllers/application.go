package controllers

import (
	"context"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

//	@BasePath	/api/v1/

type ApplicationController struct {
	kube dynamic.Interface
}

func NewApplicationController(kube dynamic.Interface) ApplicationController {
	return ApplicationController{
		kube: kube,
	}
}

func (a ApplicationController) RegisterRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/apps", a.AppPOST())
	routerGroup.GET("/apps/:id", a.AppGET())
}

type Application struct {
	Name  string `json:"name" binding:"required"`
	Image string `json:"image" binding:"required"`
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
			Name:  c.PostForm("name"),
			Image: c.PostForm("image"),
		}
		//TODO: get user namespace
		deployment := getApplication(app.Name, "default", app.Image)
		a.kube.Resource(getDeploymentGVR()).Namespace("default").Apply(context.Background(), app.Name, deployment, v1.ApplyOptions{})
	}
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
