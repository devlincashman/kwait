package k8s

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type DeploymentResource struct {
	client appsv1.DeploymentInterface
}

func NewDeploymentResource(clientset *kubernetes.Clientset, namespace string) *DeploymentResource {
	return &DeploymentResource{
		client: clientset.AppsV1().Deployments(namespace),
	}
}

func (r *DeploymentResource) Watch(ctx context.Context, namespace, name, labelSelector string) (chan metav1.Object, error) {
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
			resource := event.Object.(*appsv1.Deployment)
			resourceChan <- resource
		}
		close(resourceChan)
	}()

	return resourceChan, nil
}

func (r *DeploymentResource) IsReady(resource metav1.Object) bool {
	deployment := resource.(*appsv1.Deployment)
	return deployment.Status.AvailableReplicas == deployment.Status.Replicas
}
