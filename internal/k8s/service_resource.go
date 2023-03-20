package k8s

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
)

type ServiceResource struct {
	Client *corev1.CoreV1Client
}

func NewServiceResource(client *corev1.CoreV1Client) *ServiceResource {
	return &ServiceResource{
		Client: client,
	}
}

func (sr *ServiceResource) Watch(ctx context.Context, namespace string, name string, labelSelector string) (chan runtime.Object, error) {
	serviceClient := sr.Client.Services(namespace)
	services, err := serviceClient.List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", name),
		LabelSelector: labelSelector,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	serviceInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return serviceClient.List(ctx, options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return serviceClient.Watch(ctx, options)
			},
		},
		&corev1.Service{},
		0,
		cache.Indexers{},
	)

	serviceChan := make(chan runtime.Object)

	serviceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			serviceChan <- obj.(runtime.Object)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			serviceChan <- newObj.(runtime.Object)
		},
		DeleteFunc: func(obj interface{}) {
			serviceChan <- obj.(runtime.Object)
		},
	})

	stop := make(chan struct{})
	go serviceInformer.Run(stop)

	return serviceChan, nil
}

func (sr *ServiceResource) IsReady(resource runtime.Object) bool {
	return true
}
