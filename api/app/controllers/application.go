package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"

	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"

	appmodels "github.com/null-channel/eddington/api/app/models"
	pb "github.com/null-channel/eddington/api/proto/container"
	"github.com/null-channel/eddington/api/users/controllers"
	"github.com/null-channel/eddington/api/users/models"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	networkingapplyv1alpha3 "istio.io/client-go/pkg/applyconfiguration/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	coreapplyv1 "k8s.io/client-go/applyconfigurations/core/v1"
	coreapplymetav1 "k8s.io/client-go/applyconfigurations/meta/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kube "k8s.io/client-go/kubernetes"
)

//	@BasePath	/api/v1/

type ApplicationController struct {
	kube                   dynamic.Interface
	istioClient            *versionedclient.Clientset
	userController         *controllers.UserController
	database               *bun.DB
	containerServiceClient pb.ContainerServiceClient
	logs                   *zap.SugaredLogger
	kubeClientset          *kube.Clientset
}

func NewApplicationController(kube dynamic.Interface, istio *versionedclient.Clientset, kcs *kube.Clientset, userContoller *controllers.UserController, containerBuildingService pb.ContainerServiceClient, logger *zap.Logger) *ApplicationController {
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
		userController:         userContoller,
		database:               db,
		containerServiceClient: containerBuildingService,
		logs:                   logger.Sugar(),
		istioClient:            istio,
		kubeClientset:          kcs,
	}
}

func (a *ApplicationController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("", a.AppPOST).Methods("POST")
	router.HandleFunc("/{}", a.AppGET).Methods("GET")
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
		Name:          r.FormValue("name"),
		Image:         r.FormValue("image"),
		GitRepo:       r.FormValue("gitRepo"),
		RepoType:      r.FormValue("repoType"),
		ResourceGroup: r.FormValue("resourceGroup"),
		Directory:     r.FormValue("directory"),
	}

	fmt.Println("app: ", app)
	userId := r.Context().Value("user-id").(string)

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

	//TODO: accept the revision
	_, err = a.containerServiceClient.CreateContainer(r.Context(), &pb.CreateContainerRequest{
		RepoURL:       app.GitRepo,
		Type:          pb.Language(pb.Language_value[app.RepoType]),
		ResourceGroup: rgId,
		Directory:     app.Directory,
		Rev:           "main",
	})

	if err != nil {
		a.logs.Errorf("Failed to create container",
			"user-id", userId,
			"error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	namespace := userContext.Name + resourceGroup
	//TODO: short sha for build id?
	nullApplication := getNullApplication(app, userContext, rgId, namespace, "12345")

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
			status, err = a.containerServiceClient.BuildStatus(context.Background(), &pb.BuildStatusRequest{Repo: app.GitRepo, Directory: app.Directory})
			if err != nil {
				a.logs.Errorw("Failed to get build status", "error", err)
				time.Sleep(5 * time.Second)
				continue
			}

			if status.Status == pb.ContainerStatus_BUILT {
				keepChecking = false
			}
			time.Sleep(20 * time.Second)
		}

		// Create the services for our app.
		// TODO: This should all be in the null operator.
		// What was I thinking?!?!
		serviceName := app.Name + "-service"
		vitrualServiceName := app.Name + "-virtual-service"

		appService := createService(userId, serviceName, namespace)
		_, err = a.kubeClientset.CoreV1().Services(namespace).Apply(context.Background(), appService, metav1.ApplyOptions{FieldManager: "application/apply-patch"})

		if err != nil {
			a.logs.Errorw(err.Error())
		}

		virtualService := getVirtualService(userId, app.Name, serviceName, vitrualServiceName, namespace)
		_, err = a.istioClient.NetworkingV1alpha3().VirtualServices(namespace).Apply(context.TODO(), virtualService, metav1.ApplyOptions{FieldManager: "application/apply-patch"})

		if err != nil {
			a.logs.Errorw(err.Error())
		}
		// Create the operator object
		deployment := getApplication(app.Name, namespace, "nullchannel/"+app.Image)
		_, err = a.kube.Resource(getDeploymentGVR()).Namespace(namespace).Apply(context.Background(), app.Name, deployment, v1.ApplyOptions{FieldManager: "application/apply-patch"})
		if err != nil {
			a.logs.Errorw("Failed to apply CRD for new application",
				"user-id", userId,
				"error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}()
}

func createService(userId, serviceName, namespace string) *coreapplyv1.ServiceApplyConfiguration {
	srv := coreapplyv1.Service(serviceName, namespace)
	var port int32 = 80
	srv.Spec = coreapplyv1.ServiceSpec().WithPorts((&coreapplyv1.ServicePortApplyConfiguration{Port: &port}).WithTargetPort(intstr.IntOrString{IntVal: 8080}))
	return srv
}

func createApplyMeta(name, namespace string) coreapplymetav1.ObjectMetaApplyConfiguration {
	return coreapplymetav1.ObjectMetaApplyConfiguration{
		Name:      &name,
		Namespace: &namespace,
	}
}

func createMeta(name, namespace string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}
}

func getVirtualService(userId, appName, service, virtualServiceName, namespace string) *networkingapplyv1alpha3.VirtualServiceApplyConfiguration {
	return networkingapplyv1alpha3.VirtualService(virtualServiceName, namespace).WithSpec(networkingv1alpha3.VirtualService{
		Hosts:    []string{"*"},
		Gateways: []string{"nullcloud-gateway"},
		Http: []*networkingv1alpha3.HTTPRoute{
			&networkingv1alpha3.HTTPRoute{
				Match: []*networkingv1alpha3.HTTPMatchRequest{
					&networkingv1alpha3.HTTPMatchRequest{
						Uri: &networkingv1alpha3.StringMatch{
							MatchType: &networkingv1alpha3.StringMatch_Prefix{Prefix: "/dataplane/" + userId + "/" + appName},
						},
					},
				},
				Route: []*networkingv1alpha3.HTTPRouteDestination{
					&networkingv1alpha3.HTTPRouteDestination{
						Destination: &networkingv1alpha3.Destination{
							Host: service,
							Port: &networkingv1alpha3.PortSelector{
								Number: 80,
							},
						},
					},
				},
			},
		},
	})
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

func getIstioNetowrkGVR(resource string) schema.GroupVersionResource {
	return schema.GroupVersionResource{Group: "networking.istio.io", Version: "v1beta1", Resource: resource}
}

// AppGET godoc
//
//	@Summary	Get all applications in the org
//	@Schemes
//	@Description	Get all applications in the org
//	@Tags			example
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Helloworld
//	@Router			/apps/ [get]
func (a ApplicationController) AppGET(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user-id").(string)

	// get user org
	org, err := a.userController.GetUserContext(r.Context(), userId)
	if err != nil {
		a.logs.Errorw("Failed to get user org for the user controller",
			"user-id", userId,
			"error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	nullApplications, err := a.GetApplications(r.Context(), org.ID)
	if err != nil {
		a.logs.Errorw("Failed to get applications for the user",
			"user-id", userId,
			"error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nullApplications)
}

func (a ApplicationController) GetApplication(ctx context.Context, id int64) (*appmodels.NullApplication, error) {
	nullApplication := &appmodels.NullApplication{}
	err := a.database.NewSelect().Model(nullApplication).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return nullApplication, nil
}

func (a ApplicationController) GetApplications(ctx context.Context, orgId int64) ([]*appmodels.NullApplication, error) {
	nullApplications := []*appmodels.NullApplication{}
	err := a.database.NewSelect().Model(&nullApplications).Where("org_id = ?", orgId).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return nullApplications, nil
}
