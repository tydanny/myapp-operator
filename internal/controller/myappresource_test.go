package controller_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	myappv1alpha1 "github.com/tydanny/myapp-operator/api/v1alpha1"
	"github.com/tydanny/myapp-operator/internal/wdk"
)

var _ = Describe("Internal/Controller/Myappresource", func() {
	const (
		MyAppResourceName      = "myapp"
		MyAppResourceNamespace = "myapp"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating a MyAppResource", func() {
		It("should create podinfo pods", func() {
			By("Creating a namespace for the sample CR")
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: MyAppResourceNamespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).To(Succeed())

			By("Creating a new MyAppResource")
			ctx := context.Background()
			myapprec := &myappv1alpha1.MyAppResource{
				TypeMeta: metav1.TypeMeta{
					Kind:       "MyAppResource",
					APIVersion: myappv1alpha1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      MyAppResourceName,
					Namespace: MyAppResourceNamespace,
				},
				Spec: myappv1alpha1.MyAppResourceSpec{
					ReplicaCount: 3,
					Resources: myappv1alpha1.Resources{
						MemoryLimit: *resource.NewScaledQuantity(10, resource.Mega),
						CPURequest:  *resource.NewMilliQuantity(100, resource.BinarySI),
					},
					Image: myappv1alpha1.MyAppImage{
						Repository: "myimage",
						Tag:        "latest",
					},
					UI: myappv1alpha1.MyAppUI{
						Color:   "#34577c",
						Message: "this is a test",
					},
					Redis: myappv1alpha1.MyAppRedisConfig{
						Enabled: true,
					},
				},
			}

			Expect(k8sClient.Create(ctx, myapprec)).To(Succeed())

			myAppLookupKey := types.NamespacedName{
				Name:      MyAppResourceName,
				Namespace: MyAppResourceNamespace,
			}
			createdPodInfo := &appsv1.Deployment{}

			Eventually(ctx, func() error {
				if err := k8sClient.Get(ctx, myAppLookupKey, createdPodInfo); err != nil {
					return err
				}
				return nil
			}, timeout, interval).Should(Succeed())

			Expect(*createdPodInfo.Spec.Replicas).To(BeNumerically("==", 3))
			Expect(
				createdPodInfo.Spec.Template.Spec.Containers[0].Resources.Requests.Cpu().
					Equal(*resource.NewMilliQuantity(100, resource.BinarySI)),
			).To(BeTrue())
			Expect(
				createdPodInfo.Spec.Template.Spec.Containers[0].Resources.Limits.Memory().
					Equal(*resource.NewScaledQuantity(10, resource.Mega)),
			).To(BeTrue())
			Expect(createdPodInfo.Spec.Template.Spec.Containers[0].Image).To(Equal("myimage:latest"))
			Expect(createdPodInfo.Spec.Template.Spec.Containers[0].Env).To(ContainElements(
				corev1.EnvVar{
					Name:  wdk.PODINFO_UI_COLOR,
					Value: "#34577c",
				},
				corev1.EnvVar{
					Name:  wdk.PODINFO_UI_MESSAGE,
					Value: "this is a test",
				},
			))

			By("Checking the pod info service was created")
			Eventually(ctx, func() error {
				var createdPodInfoSvc corev1.Service
				return k8sClient.Get(ctx, myAppLookupKey, &createdPodInfoSvc)
			}).Should(Succeed())

			By("Checking the redis deployment was created")
			redisLookupKey := types.NamespacedName{
				Name:      "redis",
				Namespace: MyAppResourceNamespace,
			}
			Eventually(ctx, func() error {
				var createdRedis appsv1.Deployment
				return k8sClient.Get(ctx, redisLookupKey, &createdRedis)
			}).Should(Succeed())

			By("Checking a redis service was created")
			Eventually(ctx, func() error {
				var createdRedisSvc corev1.Service
				return k8sClient.Get(ctx, redisLookupKey, &createdRedisSvc)
			}).Should(Succeed())
		})
	})
})
