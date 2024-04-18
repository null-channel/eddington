package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"

	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"

	appmodels "github.com/null-channel/eddington/api/app/models"
	"github.com/null-channel/eddington/api/app/types"
	"github.com/null-channel/eddington/api/core"
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

type AppDatastore interface {
	CreateNullApplication(ctx context.Context, nullApplication *appmodels.NullApplication) error
	CreateNullApplicationService(ctx context.Context, nullApplication *appmodels.NullApplicationService) error
	GetApplicationsByOrgID(ctx context.Context, orgId int64) ([]*appmodels.NullApplication, error)
	GetApplicationByID(ctx context.Context, id int64) (*appmodels.NullApplication, error)
	GetApplicationServiceByAppID(ctx context.Context, nullApplicationID int64) ([]*appmodels.NullApplicationService, error)

	GetAllAppSvc(ctx context.Context) ([]*appmodels.NullApplicationService, error)
}

type ApplicationController struct {
	kube                   dynamic.Interface
	istioClient            *versionedclient.Clientset
	userController         *controllers.UserController
	appDatastore           AppDatastore
	containerServiceClient pb.ContainerServiceClient
	logger                 *zap.SugaredLogger
	kubeClientset          *kube.Clientset
}

func NewApplicationController(kube dynamic.Interface, istio *versionedclient.Clientset, kcs *kube.Clientset, appDatastore AppDatastore, userContoller *controllers.UserController, containerBuildingService pb.ContainerServiceClient, logger *zap.Logger) *ApplicationController {
	// Set up a connection to the server.

	return &ApplicationController{
		kube:                   kube,
		userController:         userContoller,
		appDatastore:           appDatastore,
		containerServiceClient: containerBuildingService,
		logger:                 logger.Sugar(),
		istioClient:            istio,
		kubeClientset:          kcs,
	}
}

func (a *ApplicationController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("", a.AppPOST).Methods("POST")
	router.HandleFunc("", a.AppGET).Methods("GET")
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
	userId := r.Context().Value("user-id").(string)

	var appDTO types.Application

	err := json.NewDecoder(r.Body).Decode(&appDTO)

	if err != nil {
		a.logger.Error(err)
		http.Error(w, "Decode error! Please check your JSON formatting.", http.StatusBadRequest)
		return
	}

	if err := appDTO.Validate(); err != nil {
		errorMessage := types.ConstructErrorMessages(err)
		core.ValidationErrors(w, errorMessage)
		return
	}

	// get user namespace
	org, err := a.userController.GetUserContext(r.Context(), userId)

	if err != nil {
		a.logger.Errorw("Failed to get user context for the user controller",
			"user-id", userId,
			"error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	resourceGroup, rgId, err := getResourceGroupName(org.ResourceGroups, appDTO.ResourceGroup)
	if err != nil {
		a.logger.Warnw("Resource Group not found", "Resource Group:", resourceGroup)
		http.Error(w, "Resource group not found", http.StatusNotFound)
		return
	}

	//TODO: accept the revision
	_, err = a.containerServiceClient.CreateContainer(r.Context(), &pb.CreateContainerRequest{
		RepoURL:       appDTO.GitRepo,
		Type:          pb.Language(pb.Language_value[appDTO.RepoType.String()]),
		ResourceGroup: rgId,
		Directory:     appDTO.Directory,
		Rev:           "main",
	})

	if err != nil {
		a.logger.Errorf("Failed to create container",
			"user-id", userId,
			"error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	namespace := org.Name + resourceGroup
	//TODO: short sha for build id?
	nullApplication, nullApplicationService := getNullApplication(appDTO, org, rgId, namespace, "12345")
	//TODO: Save to database!
	err = a.appDatastore.CreateNullApplication(r.Context(), nullApplication)
	if err != nil {
		a.logger.Errorf("Failed to create container",
			"user-id", userId,
			"error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	nullApplicationService.NullApplicationID = nullApplication.ID
	err = a.appDatastore.CreateNullApplicationService(r.Context(), nullApplicationService)
	if err != nil {
		a.logger.Errorf("Failed to create container",
			"user-id", userId,
			"error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	a.logger.Info(nullApplicationService.ID)
	go func() {
		keepChecking := true
		var status *pb.BuildStatusResponse
		for keepChecking {
			status, err = a.containerServiceClient.BuildStatus(context.Background(), &pb.BuildStatusRequest{Repo: appDTO.GitRepo, Directory: appDTO.Directory})
			if err != nil {
				a.logger.Errorw("Failed to get build status", "error", err)
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
		serviceName := appDTO.Name + "-service"
		vitrualServiceName := appDTO.Name + "-virtual-service"

		appService := createService(userId, serviceName, namespace)
		_, err = a.kubeClientset.CoreV1().Services(namespace).Apply(context.Background(), appService, metav1.ApplyOptions{FieldManager: "application/apply-patch"})

		if err != nil {
			a.logger.Errorw(err.Error())
		}

		virtualService := getVirtualService(userId, appDTO.Name, serviceName, vitrualServiceName, namespace)
		_, err = a.istioClient.NetworkingV1alpha3().VirtualServices(namespace).Apply(context.TODO(), virtualService, metav1.ApplyOptions{FieldManager: "application/apply-patch"})

		if err != nil {
			a.logger.Errorw(err.Error())
		}
		// Create the operator object
		deployment := getApplication(appDTO.Name, namespace, "nullchannel/"+appDTO.Image)
		_, err = a.kube.Resource(getDeploymentGVR()).Namespace(namespace).Apply(context.Background(), appDTO.Name, deployment, v1.ApplyOptions{FieldManager: "application/apply-patch"})
		if err != nil {
			a.logger.Errorw("Failed to apply CRD for new application",
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

func getNullApplication(app types.Application, org *models.Org, resourceGroupId int64, namespace string, buildID string) (*appmodels.NullApplication, *appmodels.NullApplicationService) {
	return &appmodels.NullApplication{
			OrgID:           org.ID,
			Name:            app.Name,
			ResourceGroupID: resourceGroupId,
			Namespace:       namespace,
		}, &appmodels.NullApplicationService{
			Type:    appmodels.ContainerImage,
			GitRepo: app.GitRepo,
			Name:    app.Name,
			Image:   app.Image,
			Cpu:     "100m",
			Memory:  "100Mi",
			Storage: "1Gi",
			BuildID: buildID,
		}

}

func getResourceGroupName(resourceGroups []*models.ResourceGroup, requested string) (string, int64, error) {
	if requested == "" && len(resourceGroups) > 0 {
		requested = "Default"
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
		a.logger.Errorw("Failed to get user org for the user controller",
			"user-id", userId,
			"error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	nullApplications, err := a.GetApplications(r.Context(), org.ID)
	if err != nil {
		a.logger.Errorw("Failed to get applications for the user",
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
	nullApplication, err := a.appDatastore.GetApplicationByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return nullApplication, nil
}

func (a ApplicationController) GetApplications(ctx context.Context, orgId int64) ([]*appmodels.NullApplication, error) {
	nullApplications, err := a.appDatastore.GetApplicationsByOrgID(ctx, orgId)
	if err != nil {
		return nil, err
	}
	return nullApplications, nil
}
