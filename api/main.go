// Package main implements a server for Greeter service.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"os"

	"github.com/gorilla/mux"
	marketing "github.com/null-channel/eddington/api/marketing/controllers"
	"github.com/null-channel/eddington/api/middleware"
	"github.com/null-channel/eddington/api/notfound"
	userController "github.com/null-channel/eddington/api/users/controllers"
	ory "github.com/ory/client-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	pb "github.com/null-channel/eddington/proto/container"

	app "github.com/null-channel/eddington/api/app/controllers"
)

var (
	addr = flag.String("addr", "eddington-container-builder:50051", "the address to connect to")
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

	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(notfound.NotFoundHandler)

	router.Use(middleware.LoggingMiddleware)

	middleware.CreateCORSHandler(router)

	//TODO: Swagger
	//docs.SwaggerInfo.BasePath = "/api/v1"
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

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewContainerServiceClient(conn)

	config := dynamic.NewForConfigOrDie(clusterConfig)
	appController := app.NewApplicationController(config, userController, client)

	v1 := router.PathPrefix("/api/v1").Subrouter()
	{
		// Apps
		apps := v1.PathPrefix("/apps").Subrouter()
		{
			apps.Use(oryMiddleware.SessionMiddleware)
			appController.RegisterRoutes(apps)
		}

		// Users
		users := v1.PathPrefix("/users").Subrouter()
		{
			users.Use(oryMiddleware.SessionMiddleware)
			userController.AddAllControllers(users)
		}
		// AuthZ

		// AuthN

		// Space

		// Marketing
		marketingGroup := v1.PathPrefix("/marketing").Subrouter()
		{
			_ = marketing.New(os.Getenv("SENDGRID_API_KEY"), marketingGroup)
		}
	}

	//TODO: serve swagger
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	host := os.Getenv("HOST")
	if host == "" {
		host = "http://localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

// getClusterConfig return the config for k8s
func getClusterConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	klog.Info("Kubeconfig flag is empty")
	return rest.InClusterConfig()
}
