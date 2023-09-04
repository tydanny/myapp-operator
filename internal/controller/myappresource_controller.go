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
	"k8s.io/apimachinery/pkg/types"
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
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch

func (r *MyAppResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("name", req.Name, "namespace", req.Namespace)

	// Fetch the CR to be reconciled
	var myApp myappv1alpha1.MyAppResource
	if err := r.Get(ctx, req.NamespacedName, &myApp); err != nil {
		return ctrl.Result{}, err
	}

	// Fetch the resources needed to reach the desired state and
	// update the status of the CR.
	podInfoLookupKey := types.NamespacedName{
		Name:      "podinfo",
		Namespace: req.Namespace,
	}
	var currentPodInfo appsv1.Deployment
	if err := r.Get(ctx, podInfoLookupKey, &currentPodInfo); client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}

	redisLookupKey := types.NamespacedName{
		Name:      "redis",
		Namespace: myApp.Namespace,
	}
	var currentRedis appsv1.Deployment
	if err := r.Get(ctx, redisLookupKey, &currentRedis); client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}

	var currentRedisSvc corev1.Service
	if err := r.Get(ctx, redisLookupKey, &currentRedisSvc); client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	}

	// Update the status of the CR
	myApp.Status.AvailableReplicas = currentPodInfo.Status.AvailableReplicas
	myApp.Status.Conditions.PodInfoConditions = currentPodInfo.Status.Conditions
	myApp.Status.Conditions.RedisConditions = currentRedis.Status.Conditions

	// If the status update failes we don't want to block reconciliation,
	// so instead log the error and continue on.
	if err := r.Status().Update(ctx, &myApp); err != nil {
		log.Error(err, "failed to update status")
	}

	// The following code progesses the current cluster state towards the
	// desired cluster state.

	// If redis is enabled, create a redis service and deployment.
	// Otherwise, ensure that there is not redis service/deployment.
	desiredRedis := buildRedis(&myApp)
	desiredRedisSvc := buildRedisService(&myApp)
	if myApp.Spec.Redis.Enabled {
		if err := r.Apply(ctx, desiredRedis); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to apply service: %w", err)
		}

		if err := r.Apply(ctx, desiredRedisSvc); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to apply service: %w", err)
		}
	} else {
		if err := r.Delete(ctx, desiredRedisSvc); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, fmt.Errorf("failed to delete redis service: %w", err)
		}

		if err := r.Delete(ctx, desiredRedis); client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, fmt.Errorf("failed to delete redis deployment: %w", err)
		}
	}

	// If redis is enabled we need the service ready to build the
	// podinfo deployment. So, if the cluster IP is empty, requeue and try
	// again later.
	// Otherwise, create the deployment/service using the config specified in the
	// CR.
	if myApp.Spec.Redis.Enabled && currentRedisSvc.Spec.ClusterIP == "" {
		return ctrl.Result{Requeue: true}, nil
	}
	desiredPodInfo := buildPodInfo(&myApp, currentRedisSvc.Spec.ClusterIP)
	if err := r.Apply(ctx, desiredPodInfo); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to apply deployment: %w", err)
	}

	desiredPodInfoSvc := buildPodInfoService(&myApp)
	if err := r.Apply(ctx, desiredPodInfoSvc); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to apply service: %w", err)
	}

	log.Info("finished reconcile successfully")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyAppResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myappv1alpha1.MyAppResource{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&appsv1.Deployment{}, builder.WithPredicates(MyAppLabelPredicate)).
		Owns(&corev1.Service{}, builder.WithPredicates(MyAppLabelPredicate)).
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

func buildRedis(myAppRec *myappv1alpha1.MyAppResource) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis",
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
				wdk.App: wdk.MyApp,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					wdk.App: "redis",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						wdk.App: "redis",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "redis",
							Image: "redis:latest",
							Command: []string{
								"redis-server",
								"/redis-master/redis.conf",
							},
							Ports: []corev1.ContainerPort{{
								ContainerPort: 6379,
							}},
							Env: []corev1.EnvVar{{
								Name:  "MASTER",
								Value: "true",
							}},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data",
									MountPath: "/redis-master-data",
								},
								{
									Name:      "config",
									MountPath: "/redis-master",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "data",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "redis-conf",
									},
									Items: []corev1.KeyToPath{{
										Key:  "conf",
										Path: "redis.conf",
									}},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildRedisService(myAppRec *myappv1alpha1.MyAppResource) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: corev1.SchemeGroupVersion.Version,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis",
			Namespace: myAppRec.Namespace,
			Labels: map[string]string{
				wdk.App: wdk.MyApp,
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
				wdk.App: "redis",
			},
			Ports: []corev1.ServicePort{{
				Port:       6379,
				TargetPort: intstr.FromInt(6379),
			}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}

func buildPodInfoService(myAppRec *myappv1alpha1.MyAppResource) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: corev1.SchemeGroupVersion.Version,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "podinfo",
			Namespace: myAppRec.Namespace,
			Labels: map[string]string{
				wdk.App: wdk.MyApp,
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
				wdk.App: wdk.MyApp,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       8080,
					TargetPort: intstr.FromString("http"),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}

func buildPodInfo(myAppRec *myappv1alpha1.MyAppResource, redisIP string) *appsv1.Deployment {
	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "podinfo",
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
				wdk.App: wdk.MyApp,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(myAppRec.Spec.ReplicaCount),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					wdk.App: wdk.MyApp,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						wdk.App: wdk.MyApp,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "podinfo",
							Image: myAppRec.Spec.Image.String(),
							Command: []string{
								"./podinfo",
								"--port=8080",
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8080,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Env: []corev1.EnvVar{
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
						},
					},
					Volumes: []corev1.Volume{},
				},
			},
		},
	}

	// If the redis IP is not empty add the corresponding environment variable.
	// If it is an empty string then redis is not enabled and the environement variable
	// can be omitted.
	if redisIP != "" {
		deploy.Spec.Template.Spec.Containers[0].Env = append(
			deploy.Spec.Template.Spec.Containers[0].Env,
			corev1.EnvVar{
				Name:  wdk.PODINFO_CACHE_SERVER,
				Value: "tcp://" + redisIP + ":6379",
			},
		)
	}

	return deploy
}
