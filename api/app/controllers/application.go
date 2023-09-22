package controllers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"go.uber.org/zap"
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
	logs                   *zap.SugaredLogger
}

func NewApplicationController(kube dynamic.Interface, userService *usercon.UserController, containerBuildingService pb.ContainerServiceClient, logger *zap.Logger) *ApplicationController {
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
		logs:                   logger.Sugar(),
	}
}

func (a *ApplicationController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("", a.AppPOST).Methods("POST")
	router.HandleFunc("/{id}", a.AppGET).Methods("GET")
}

type Application struct {
	Name          string `json:"name"`
	Image         string `json:"image"`
	GitRepo       string `json:"gitRepo"`
	RepoType      string `json:"repoType"`
	ResourceGroup string `json:"resourceGroup"`
	Directory     string `json:"directory"`
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
func (a ApplicationController) AppPOST(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	app := Application{
		Name:          r.Form.Get("name"),
		Image:         r.Form.Get("image"),
		GitRepo:       r.Form.Get("gitRepo"),
		RepoType:      r.Form.Get("repoType"),
		ResourceGroup: r.Form.Get("resourceGroup"),
		Directory:     r.Form.Get("directory"),
	}

	userId, err := strconv.ParseInt(r.Context().Value("user-id").(string), 10, 64)

	if err != nil {
		a.logs.Errorf("Failed to parse user id",
			"user-id:", r.Context().Value("user-id"))
	}
	// get user namespace
	userContext, err := a.userController.GetUserContext(r.Context(), userId)
	if err != nil {
		a.logs.Errorw("Failed to get user context for the user controller",
			"user-id", userId,
			"error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	fmt.Println(userContext)
	resourceGroup, rgId, err := getResourceGroupName(userContext.ResourceGroups, app.ResourceGroup)
	if err != nil {

		a.logs.Warnw("Resource Group not found", "Resource Group:", resourceGroup)
		http.Error(w, "Resource group not found", http.StatusNotFound)
		return
	}

	ret, err := a.containerServiceClient.CreateContainer(r.Context(), &pb.CreateContainerRequest{
		RepoURL:    app.GitRepo,
		Type:       pb.Language(pb.Language_value[app.RepoType]),
		CustomerID: userId,
		Directory:  app.Directory,
	})

	if err != nil {
		a.logs.Errorf("Failed to create container",
			"user-id", userId,
			"error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	namespace := userContext.Name + resourceGroup
	nullApplication := getNullApplication(app, userContext, rgId, namespace, ret.BuildID)

	//TODO: Save to database!
	_, err = a.database.NewInsert().Model(nullApplication).Exec(context.Background())

	if err != nil {
		a.logs.Errorf("Failed to create container",
			"user-id", userId,
			"error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	go func() {
		keepChecking := true
		var status *pb.BuildStatusResponse
		for keepChecking {
			status, err = a.containerServiceClient.BuildStatus(context.Background(), &pb.BuildStatusRequest{Id: ret.BuildID})
			if err != nil {
				a.logs.Errorw("Failed to get build status", "error", err)
				panic("checking container build status failed")
			}

			if status.Status == pb.ContainerStatus_BUILT {
				keepChecking = false
			}
			time.Sleep(20 * time.Second)
		}

		deployment := getApplication(app.Name, namespace, "nullchannel/"+app.Image)
		_, err = a.kube.Resource(getDeploymentGVR()).Namespace(namespace).Apply(context.Background(), app.Name, deployment, v1.ApplyOptions{})
		if err != nil {
			a.logs.Errorw("Failed to apply CRD for new application",
				"user-id", userId,
				"error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}()
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
	return schema.GroupVersionResource{Group: "nullapp.io.nullcloud", Version: "v1alpha1", Resource: "nullapplications"}
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
func (a ApplicationController) AppGET(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	w.WriteHeader(http.StatusNotImplemented)
}
