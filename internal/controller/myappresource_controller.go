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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	myappv1alpha1 "github.com/tydanny/myapp-operator/api/v1alpha1"
	"github.com/tydanny/myapp-operator/internal/wdk"
)

// MyAppResourceReconciler reconciles a MyAppResource object
type MyAppResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=myapp.example.com,resources=myappresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=myapp.example.com,resources=myappresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=myapp.example.com,resources=myappresources/finalizers,verbs=update

func (r *MyAppResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("name", req.Name, "namespace", req.Namespace)

	// Fetch the CR to be reconciled
	var myApp myappv1alpha1.MyAppResource
	if err := r.Get(ctx, req.NamespacedName, &myApp); err != nil {
		return ctrl.Result{}, err
	}

	// Build the desired Deployment and apply it to the cluster using server side apply
	desiredService := buildService(&myApp)
	desiredDeployment := buildDeployment(&myApp)

	if err := r.Apply(ctx, desiredService); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Apply(ctx, desiredDeployment); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("finished reconcile")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyAppResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myappv1alpha1.MyAppResource{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&appsv1.Deployment{}, builder.WithPredicates(MyAppLabelPredicate)).
		Complete(r)
}

func (r *MyAppResourceReconciler) Apply(ctx context.Context, obj client.Object, opts ...client.PatchOption) error {
	return r.Patch(
		ctx,
		obj,
		client.Apply,
		client.ForceOwnership,
		client.FieldOwner(wdk.MyAppController),
	)
}

func buildService(myAppRec *myappv1alpha1.MyAppResource) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: corev1.SchemeGroupVersion.Version,
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: myAppRec.Name,
			Namespace:    myAppRec.Namespace,
			Labels: map[string]string{
				wdk.MyApp: myAppRec.Name,
			},
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion:         myAppRec.APIVersion,
				Kind:               myAppRec.Kind,
				Name:               myAppRec.Name,
				UID:                myAppRec.UID,
				Controller:         pointer.Bool(true),
				BlockOwnerDeletion: pointer.Bool(true),
			}},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				wdk.MyApp: myAppRec.Name,
			},
			Ports: []corev1.ServicePort{{
				Name:       "http",
				Protocol:   corev1.ProtocolTCP,
				Port:       80,
				TargetPort: intstr.FromInt(80),
			}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}

func buildDeployment(myAppRec *myappv1alpha1.MyAppResource) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.Version,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      myAppRec.Name,
			Namespace: myAppRec.Namespace,
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion:         myAppRec.APIVersion,
				Kind:               myAppRec.Kind,
				Name:               myAppRec.Name,
				UID:                myAppRec.UID,
				Controller:         pointer.Bool(true),
				BlockOwnerDeletion: pointer.Bool(true),
			}},
			Labels: map[string]string{
				wdk.MyApp: myAppRec.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(myAppRec.Spec.ReplicaCount),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					wdk.MyApp: myAppRec.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						wdk.MyApp: myAppRec.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "my-app",
						Image: fmt.Sprintf("%s:%s", myAppRec.Spec.Image.Repository, myAppRec.Spec.Image.Tag),
						// Args:  []string{},
						// Ports: []corev1.ContainerPort{},
						Env: []corev1.EnvVar{
							{
								Name:  wdk.PODINFO_CACHE_SERVER,
								Value: "localhost",
							},
							{
								Name:  wdk.PODINFO_UI_COLOR,
								Value: myAppRec.Spec.UI.Color,
							},
							{
								Name:  wdk.PODINFO_UI_MESSAGE,
								Value: myAppRec.Spec.UI.Message,
							},
						},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: myAppRec.Spec.Resources.MemoryLimit,
							},
							Requests: corev1.ResourceList{
								corev1.ResourceCPU: myAppRec.Spec.Resources.CPURequest,
							},
						},
						VolumeMounts: []corev1.VolumeMount{},
					}},
					Volumes: []corev1.Volume{},
				},
			},
		},
	}
}
