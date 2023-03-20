package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type PodResource struct {
	client corev1.PodInterface
}

func NewPodResource(clientset *kubernetes.Clientset, namespace string) *PodResource {
	return &PodResource{
		client: clientset.CoreV1().Pods(namespace),
	}
}

func (r *PodResource) Watch(ctx context.Context, namespace, name, labelSelector string) (chan metav1.Object, error) {
	watchOptions := metav1.ListOptions{
		FieldSelector:  "metadata.name=" + name,
		LabelSelector:  labelSelector,
		TimeoutSeconds: 3600,
		Watch:          true,
	}

	watcher, err := r.client.Watch(ctx, watchOptions)
	if err != nil {
		return nil, err
	}

	resourceChan := make(chan metav1.Object)
	go func() {
		defer watcher.Stop()
		for event := range watcher.ResultChan() {
			resource := event.Object.(*corev1.Pod)
			resourceChan <- resource
		}
		close(resourceChan)
	}()

	return resourceChan, nil
}

func (r *PodResource) IsReady(resource metav1.Object) bool {
	pod := resource.(*corev1.Pod)
	for _, cond := range pod.Status.Conditions {
		if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}
