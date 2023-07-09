// Package main implements a server for Greeter service.
package main

import (
	"flag"
	"fmt"
	"log"

	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/null-channel/eddington/api/docs"
	marketing "github.com/null-channel/eddington/api/marketing/controllers"
	"github.com/null-channel/eddington/api/middleware"
	userController "github.com/null-channel/eddington/api/users/controllers"
	ory "github.com/ory/client-go"
	ginSwagger "github.com/swaggo/gin-swagger"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	app "github.com/null-channel/eddington/api/app/controllers"
	swaggerfiles "github.com/swaggo/files"
)

func main() {

	// ORY Stuff Not sure this is a good way to deal with this.
	proxyPort := os.Getenv("ORY_PROXY_PORT")
	if proxyPort == "" {
		proxyPort = "4000"
	}

	oryDomain := os.Getenv("ORY_DOMAIN")
	if oryDomain == "" {
		oryDomain = "http://localhost"
	}

	// register a new Ory client with the URL set to the Ory CLI Proxy
	// we can also read the URL from the env or a config file
	c := ory.NewConfiguration()
	c.Servers = ory.ServerConfigurations{{URL: fmt.Sprintf("%s:%s/.ory", oryDomain, proxyPort)}}

	oryMiddleware := &middleware.OryApp{
		Ory: ory.NewAPIClient(c),
	}

	fmt.Println("Starting server...")
	flag.Parse()
	//TODO: Never used gin. seems like mux is archived. going to try this out.

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(cors.Default())
	docs.SwaggerInfo.BasePath = "/api/v1"
	userController, err := userController.New()

	if err != nil {
		log.Fatal(err)
		return
	}

	kubeConfig := os.Getenv("KUBECONFIG")

	clusterConfig, err := getClusterConfig(kubeConfig)

	if err != nil {
		log.Fatal(err)
		return
	}

	config := dynamic.NewForConfigOrDie(clusterConfig)
	appController := app.NewApplicationController(config, userController)

	v1 := router.Group("api/v1")
	{
		// Apps
		apps := v1.Group("/apps")
		{
			apps.Use(oryMiddleware.SessionMiddleware())
			appController.RegisterRoutes(apps)
		}

		// Users
		users := v1.Group("/users")
		{
			users.Use(oryMiddleware.SessionMiddleware())
			userController.AddAllControllers(users)
		}
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

	host := os.Getenv("HOST")
	if host == "" {
		host = "http://localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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
