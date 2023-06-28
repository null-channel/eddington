// Package main implements a server for Greeter service.
package main

import (
	"log"

	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	app "github.com/null-channel/eddington/api/app/controllers"
	"github.com/null-channel/eddington/api/docs"
	marketing "github.com/null-channel/eddington/api/marketing/controllers"
	userController "github.com/null-channel/eddington/api/users/controllers"
	ginSwagger "github.com/swaggo/gin-swagger"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	swaggerfiles "github.com/swaggo/files"
)

func main() {

	//TODO: Never used gin. seems like mux is archived. going to try this out.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())
	docs.SwaggerInfo.BasePath = "/api/v1"
	userController := userController.New()

	kubeConfig := os.Getenv("KUBECONFIG")

	clusterConfig, err := getClusterConfig(kubeConfig)

	if err != nil {
		log.Fatal(err)
		return
	}
	config := dynamic.NewForConfigOrDie(clusterConfig)
	appController := app.NewApplicationController(config)

	v1 := router.Group("api/v1")
	{
		// Apps
		apps := v1.Group("/apps")
		{
			appController.RegisterRoutes(apps)
		}

		// Users
		userController.AddAllControllers(v1)

		// AuthZ

		// AuthN

		// Space

		// Marketing
		marketingGroup := v1.Group("/marketing")
		{
			_ = marketing.New(os.Getenv("SENDGRID_API_KEY"), marketingGroup)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	log.Fatal(router.Run(":" + port))

}

// getClusterConfig return the config for k8s
func getClusterConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	klog.Info("Kubeconfig flag is empty")
	return rest.InClusterConfig()
}
