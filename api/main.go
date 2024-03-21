// Package main implements a server for Greeter service.
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	marketing "github.com/null-channel/eddington/api/marketing/controllers"
	"github.com/null-channel/eddington/api/middleware"
	"github.com/null-channel/eddington/api/notfound"
	userController "github.com/null-channel/eddington/api/users/controllers"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	kube "k8s.io/client-go/kubernetes"

	pb "github.com/null-channel/eddington/api/proto/container"

	app "github.com/null-channel/eddington/api/app/controllers"
)

var (
	addr  = flag.String("addr", "eddington-container-builder:50051", "the address to connect to")
	debug = flag.Bool("debug", false, "Run in debug mode")
)

func main() {
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	// ORY Stuff Not sure this is a good way to deal with this.
	proxyPort := os.Getenv("ORY_PROXY_PORT")
	if proxyPort == "" {
		proxyPort = "4000"
	}

	fmt.Println("Starting server...")
	flag.Parse()

	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(notfound.NotFoundHandler)

	router.Use(middleware.LoggingMiddleware)

	//Database stuff
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	userdb := bun.NewDB(sqldb, sqlitedialect.New())
	if err != nil {
		panic(err)
	}

	//TODO: Swagger
	//docs.SwaggerInfo.BasePath = "/api/v1"
	userController, err := userController.New(logger, userdb)

	//Add user middleware
	userMiddleware := middleware.NewUserMiddleware(userdb)

	authzMiddleware := middleware.NewAuthzMiddleware(userdb)

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
	kubeClient := kube.NewForConfigOrDie(clusterConfig)

	config := dynamic.NewForConfigOrDie(clusterConfig)
	istioClient := versionedclient.NewForConfigOrDie(clusterConfig)
	appController := app.NewApplicationController(config, istioClient, kubeClient, userController, client, logger)
	middlwares := []mux.MiddlewareFunc{
		middleware.AddJwtHeaders,
		authzMiddleware.CheckAuthz,
		userMiddleware.NewUserMiddlewareCheck,
	}

	v1 := router.PathPrefix("/api/v1").Subrouter()
	{
		// Apps
		apps := v1.PathPrefix("/apps").Subrouter()
		{
			addMiddleware(apps, middlwares...)
			appController.RegisterRoutes(apps)
		}

		// Users
		users := v1.PathPrefix("/users").Subrouter()
		{
			users.Use(middleware.AddJwtHeaders)
			users.Use(authzMiddleware.CheckAuthz)
			userController.AddAllControllers(users)
		}
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

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:9000", "http://localhost:5173", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	srv := &http.Server{
		Handler: handler,
		Addr:    "0.0.0.0:" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

func addMiddleware(router *mux.Router, middleware ...mux.MiddlewareFunc) *mux.Router {
	for _, m := range middleware {
		router.Use(m)
	}
	return router
}

// getClusterConfig return the config for k8s
func getClusterConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	klog.Info("Kubeconfig flag is empty")
	return rest.InClusterConfig()
}
