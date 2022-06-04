/*
Copyright 2022.

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

package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/pingcap/errors"
	crdv1beta1 "github.com/swh0318/sample-operator/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MyCRDReconciler reconciles a MyCRD object
type MyCRDReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// TODO(user):用户自定义逻辑
//+kubebuilder:rbac:groups=crd.github.com,resources=mycrds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crd.github.com,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crd.github.com,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=crd.github.com,resources=mycrds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=crd.github.com,resources=mycrds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyCRD object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *MyCRDReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	// TODO(user):用户自定义逻辑
	myCRD := &crdv1beta1.MyCRD{}
	if err := r.Get(ctx, req.NamespacedName, myCRD); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deploymentName := myCRD.Spec.DeploymentName
	if deploymentName == "" {
		return ctrl.Result{}, fmt.Errorf("%s: deployment name must be specified", myCRD.GetGenerateName)
	}
	deployment := &appsv1.Deployment{}
	err := r.Get(ctx, req.NamespacedName, deployment)
	if err == nil {
		if !metav1.IsControlledBy(deployment, myCRD) {
			msg := fmt.Sprintf("deployment is controlled by myCRD: %s", deployment.Name)
			return ctrl.Result{}, fmt.Errorf("%s", msg)
		}
		if myCRD.Spec.Replicas != nil && *myCRD.Spec.Replicas != *deployment.Spec.Replicas {
			klog.Infof("============mycrd replicas not equal deployment replicas===========")
			newDeployment := r.newDeployment(myCRD)
			if err := r.Update(ctx, newDeployment); err != nil && !errors.IsAlreadyExists(err) {
				return ctrl.Result{}, err
			}
		}
		// 更新CRD状态
		myCRD.Status.AvailableReplicas = deployment.Status.AvailableReplicas
		if err := r.Status().Update(ctx, myCRD); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil && errors.IsNotFound(err) {
		// 创建MyCRD对应的deployment
		deployment := r.newDeployment(myCRD)
		if err := controllerutil.SetControllerReference(myCRD, deployment, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, deployment); err != nil && !errors.IsAlreadyExists(err) {
			return ctrl.Result{}, err
		}
		klog.Infof("===========mycrd create deloyment success!!============")
	} else {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyCRDReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&crdv1beta1.MyCRD{}).
		Complete(r)
}

// newDeployment creates a new Deployment for a MyCRD resource
func (r *MyCRDReconciler) newDeployment(myCRD *crdv1beta1.MyCRD) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "mycrd",
		"controller": "mycrd",
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      myCRD.Spec.DeploymentName,
			Namespace: myCRD.Namespace,
			// OwnerReferences: []metav1.OwnerReference{
			// 	*metav1.NewControllerRef(myCRD, v1beta1.SchemeGroupVersion.WithKind("MyCRD")),
			// },
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: myCRD.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "doraemon-frontend",
							Image: "360cloud/doraemon-frontend",
						},
					},
				},
			},
		},
	}
}
