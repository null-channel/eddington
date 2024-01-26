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
	"context"
	"fmt"

	"github.com/pivotal/kpack/pkg/apis/build/v1alpha2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ImageReconciler reconciles a DontUseMe object
type ImageReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	BuildCache map[types.NamespacedName]v1alpha2.Image
}

//+kubebuilder:rbac:groups=deleteme.nullcloud.io,resources=dontusemes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=deleteme.nullcloud.io,resources=dontusemes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=deleteme.nullcloud.io,resources=dontusemes/finalizers,verbs=update
//+kubebuilder:rbac:groups=kpack.io,resources=images,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kpack.io,resources=images/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kpack.io,resources=images/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DontUseMe object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *ImageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	fmt.Println("Reconcile on image: " + req.Name)
	// TODO(user): your logic here
	var image v1alpha2.Image
	err := r.Get(ctx, req.NamespacedName, &image)

	if err != nil {
		fmt.Println("eeeeeeerrrrrrrrrorrr: " + err.Error())
		return ctrl.Result{}, err
	}

	oldValue, found := r.BuildCache[req.NamespacedName]

	if !found {
		r.BuildCache[req.NamespacedName] = image
		fmt.Println("initial cacheing")
		return ctrl.Result{}, nil
	}

	if oldValue.Status.Conditions[0].Reason != "UpToDate" {
		if image.Status.Conditions[0].Reason == "UpToDate" {
			fmt.Println("New image built: " + image.Status.LatestImage)
		}
	} else {
		fmt.Println("updates but no new image")
	}

	r.BuildCache[req.NamespacedName] = image

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ImageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha2.Image{}).
		Complete(r)
}
