/*
Copyright 2025 raibagi.

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
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"

	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apiv1alpha1 "github.com/rvitz/Samples/api/v1alpha1"
)

// SampleReconciler reconciles a Sample object
type SampleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=api.raibagi.samples,resources=samples,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=api.raibagi.samples,resources=samples/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=api.raibagi.samples,resources=samples/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Sample object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *SampleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconcile started..!")
	sample := &apiv1alpha1.Sample{}
	err := r.Get(ctx, req.NamespacedName, sample)
	if err != nil {
		return ctrl.Result{}, nil
	}

	startTime := sample.Spec.Start
	endTime := sample.Spec.End
	replicas := sample.Spec.Replicas

	currHr := time.Now().Hour()
	if currHr >= startTime && endTime <= currHr {
		for _, deploy := range sample.Spec.Deployments {
			deployment := &v1.Deployment{}
			err := r.Get(ctx, types.NamespacedName{
				Namespace: deploy.Namespace,
				Name:      deploy.Name,
			}, deployment)
			if err != nil {
				return ctrl.Result{}, err
			}
			if deployment.Spec.Replicas != &replicas {
				deployment.Spec.Replicas = &replicas
				err := r.Update(ctx, deployment)
				if err != nil {
					return ctrl.Result{}, err
				}
			}
		}
	}

	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SampleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.Sample{}).
		Complete(r)
}
