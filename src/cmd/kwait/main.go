package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/devlincashman/kwait/internal/k8s"
	"k8s.io/apimachinery/pkg/runtime"
)

func main() {
	var namespace, kind, name, labelSelector string
	flag.StringVar(&namespace, "n", "default", "Kubernetes namespace for the resource.")
	flag.StringVar(&kind, "kind", "", "Kind of the Kubernetes resource (e.g., pod, service, job, deployment).")
	flag.StringVar(&name, "name", "", "Name of the Kubernetes resource.")
	flag.StringVar(&labelSelector, "l", "", "Label selector for the Kubernetes resource.")
	flag.Parse()

	if kind == "" || (name == "" && labelSelector == "") {
		flag.Usage()
		os.Exit(1)
	}

	clientset, err := k8s.NewClient()
	if err != nil {
		fmt.Println("Error creating Kubernetes client:", err)
		os.Exit(1)
	}

	var resourceWatcher k8s.ResourceWatcher
	switch strings.ToLower(kind) {
	case "pod":
		resourceWatcher = k8s.NewPodResource(clientset, namespace)
	case "service":
		resourceWatcher = k8s.NewServiceResource(clientset.CoreV1(), namespace)
	case "deployment":
		resourceWatcher = k8s.NewDeploymentResource(clientset, namespace)
	default:
		fmt.Printf("Invalid resource type: %s\n", kind)
		os.Exit(1)
	}

	ctx := context.Background()
	resourceChan, err := resourceWatcher.Watch(ctx, namespace, name, labelSelector)
	if err != nil {
		fmt.Printf("Error watching resource: %v\n", err)
		os.Exit(1)
	}

	for resource := range resourceChan {
		if resourceWatcher.IsReady(resource.(runtime.Object)) {
			fmt.Println("Resource is ready.")
			os.Exit(0)
		}
	}

	fmt.Println("Resource is not ready or not found.")
	os.Exit(1)
}
