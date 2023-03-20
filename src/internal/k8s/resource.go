package k8s

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
)

type ResourceWatcher interface {
	Watch(namespace string, name string, labelSelector string) (<-chan runtime.Object, error)
}

type ResourceReadyChecker interface {
	IsReady(resource runtime.Object) bool
}

type Resource struct {
	Watcher      ResourceWatcher
	ReadyChecker ResourceReadyChecker
}

func (r *Resource) WaitForReady(namespace string, name string, labelSelector string) error {
	resourceChan, err := r.Watcher.Watch(namespace, name, labelSelector)
	if err != nil {
		return fmt.Errorf("failed to watch resource: %w", err)
	}

	for resource := range resourceChan {
		if r.ReadyChecker.IsReady(resource) {
			return nil
		}
	}

	return fmt.Errorf("resource not found or failed to become ready")
}
