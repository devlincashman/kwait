package k8s

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	"k8s.io/client-go/tools/cache"
)

type JobResource struct {
	Client batchv1.JobInterface
}

func NewJobResource(clientset *kubernetes.Clientset, namespace string) *JobResource {
	return &JobResource{
		Client: clientset.BatchV1().Jobs(namespace),
	}
}

func (jr *JobResource) Watch(ctx context.Context, namespace string, name string, labelSelector string) (<-chan runtime.Object, error) {
	jobList, err := jr.Client.List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", name),
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	jobInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return jr.Client.List(ctx, options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return jr.Client.Watch(ctx, options)
			},
		},
		&batchv1.Job{},
		0,
		cache.Indexers{},
	)

	jobChan := make(chan runtime.Object)

	jobInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			jobChan <- obj.(runtime.Object)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			jobChan <- newObj.(runtime.Object)
		},
		DeleteFunc: func(obj interface{}) {
			jobChan <- obj.(runtime.Object)
		},
	})

	stop := make(chan struct{})
	go jobInformer.Run(stop)

	return jobChan, nil
}

func (jr *JobResource) IsReady(resource runtime.Object) bool {
	job := resource.(*batchv1.Job)
	return job.Status.Succeeded > 0
}
