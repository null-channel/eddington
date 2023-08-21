/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"bytes"
	"context"
	"text/template"

	app "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"

	nullappv1alpha1 "github.com/null-channel/eddington/container-runner/api/v1alpha1"
	"github.com/null-channel/eddington/container-runner/internal/templates"
)

//nolint:golint,unused
var (
	setupLog               = ctrl.Log.WithName("null-application-controller")
	istioVirtualServiceGVR = schema.GroupVersionResource{Group: "networking.istio.io", Version: "v1beta1", Resource: "virtualservice"}
	istioVirtualServiceGVK = schema.GroupVersionKind{Group: "networking.istio.io", Version: "v1beta1", Kind: "VirtualService"}
)

// NullApplicationReconciler reconciles a NullApplication object
type NullApplicationReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	DynoClient *dynamic.DynamicClient
}

//+kubebuilder:rbac:groups=nullapp.io.nullcloud,resources=nullapplications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nullapp.io.nullcloud,resources=nullapplications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nullapp.io.nullcloud,resources=nullapplications/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NullApplication object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *NullApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var nullApplication nullappv1alpha1.NullApplication
	if err := r.Get(ctx, req.NamespacedName, &nullApplication); err != nil {
		log.Error(err, "unable to fetch NullApplication")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	for _, workload := range nullApplication.Spec.Apps {
		// Create a kubernetes service for each microservice
		// service name should look like this: [nullapplication-name]-[microservice-name]
		workloadNamespacedName := types.NamespacedName{
			Name:      nullApplication.Spec.AppName + "-" + workload.Name,
			Namespace: req.Namespace,
		}

		err := r.CheckDeployment(ctx, workloadNamespacedName, nullApplication, workload)
		if err != nil {
			log.Error(err, "unable to update deployment for microservice")
			return ctrl.Result{}, err
		}

		err = r.CheckDeploymentService(ctx, workloadNamespacedName, nullApplication, workload)
		if err != nil {
			log.Error(err, "unable to update service for microservice")
			return ctrl.Result{}, err
		}

	}

	return ctrl.Result{}, nil
}

func (r *NullApplicationReconciler) CheckDeploymentService(ctx context.Context, workloadNamespacedName types.NamespacedName, nullApplication nullappv1alpha1.NullApplication, workload nullappv1alpha1.MicroserviceSpec) error {
	var service v1.Service
	if err := r.Get(ctx, workloadNamespacedName, &service); err != nil {
		if errors.IsNotFound(err) {
			serviceTemplate := templates.ServiceTemplate{
				NullApplicationName: nullApplication.Spec.AppName,
				AppName:             workload.Name,
				CustomerID:          workloadNamespacedName.Namespace,
			}
			t, err := template.New("service").Parse(templates.Service)
			if err != nil {
				setupLog.Info("Error parsing service template")
				return err
			}

			templateOutput := ""
			byteBuffer := bytes.NewBufferString(templateOutput)
			err = t.Execute(byteBuffer, serviceTemplate)
			if err != nil {
				setupLog.Info("Error executing service template")
				return err
			}
			appService, err := virtualServiceBytesToUnstructured(*byteBuffer)

			if err != nil {
				setupLog.Info("Error converting service template to unstructured")
				return err
			}

			setupLog.Info("Creating service for microservice: " + byteBuffer.String())
			return r.Create(ctx, appService)
		}
	}
	// TODO It is found!!! Need to update the service.
	return nil
}

func (r *NullApplicationReconciler) CheckDeployment(ctx context.Context, workloadNamespacedName types.NamespacedName, nullApplication nullappv1alpha1.NullApplication, workload nullappv1alpha1.MicroserviceSpec) error {

	var deployment app.Deployment
	if err := r.Get(ctx, workloadNamespacedName, &deployment); err != nil {
		if errors.IsNotFound(err) {
			deploymentVars := templates.DeploymentTemplate{
				NullApplicationName: nullApplication.Spec.AppName,
				AppName:             workload.Name,
				CustomerID:          workloadNamespacedName.Namespace,
				Image:               workload.Image,
			}
			deploymentTemplate, err := template.New("deployment").Parse(templates.Deployment)
			if err != nil {
				return err
			}

			templateOutput := ""
			byteBuffer := bytes.NewBufferString(templateOutput)
			err = deploymentTemplate.Execute(byteBuffer, deploymentVars)
			if err != nil {
				return err
			}
			deployment := app.Deployment{}
			err = yaml.Unmarshal(byteBuffer.Bytes(), &deployment)
			if err != nil {
				return err
			}
			err = r.Create(ctx, &deployment)
			if err != nil {
				return err
			}
		}
	}
	// TODO It is found!!! Need to update the deployment.
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NullApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nullappv1alpha1.NullApplication{}).
		Complete(r)
}

//nolint:golint,unused
func (r *NullApplicationReconciler) getObject(namespace, name string, gvr schema.GroupVersionResource, ctx context.Context) (*unstructured.Unstructured, error) {
	return r.DynoClient.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}

func virtualServiceBytesToUnstructured(data bytes.Buffer) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{Object: map[string]interface{}{}}
	if err := yaml.Unmarshal(data.Bytes(), &u); err != nil {
		return nil, err
	}

	return u, nil
}
