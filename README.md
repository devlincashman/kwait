# kwait

🚀 A simple, yet powerful Go library and CLI tool to wait for Kubernetes resources to enter the desired state. 🚀

Entirely inspired by the amazing [k8s-wait-for](https://github.com/groundnuty/k8s-wait-for) project. The cli and featureset of kwait is designed to be a superset of [k8s-wait-for](https://github.com/groundnuty/k8s-wait-for)'s work with the power of Go!

[![Go Report Card](https://goreportcard.com/badge/github.com/devlincashman/kwait)](https://goreportcard.com/report/github.com/devlincashman/kwait)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Features

✅ Supports waiting for Pods, Services, Deployments and Jobs
✅ Provides additional wait modes for Pods and Jobs
✅ Can be used as a library or a standalone CLI tool
✅ Compatible with Kubernetes versions >= 1.24

## Installation

go get -u github.com/devlincashman/kwait/cmd/kwait

## Usage

### CLI

The kwait CLI tool allows you to wait for Kubernetes resources to become ready. Here's an overview of the available flags:

```bash
Usage of kwait:
  -kind string
        Kind of the Kubernetes resource (e.g., pod, service, job, deployment).
  -l string
        Label selector for the Kubernetes resource.
  -n string
        Kubernetes namespace for the resource. (default "default")
  -name string
        Name of the Kubernetes resource.
```

To wait for a resource to become ready, simply run the kwait program and specify the resource type, name, and (optionally) a label selector. For example, to wait for a deployment named my-deployment in the my-namespace namespace to become ready, run:

```bash
kwait -n my-namespace -kind deployment -name my-deployment
```

#### Examples

Wait for all pods with a following label to enter 'Ready' state:
`kwait pod -lapp=example`

Wait for all selected pods to enter the 'Ready' state:

`kwait pod -l"release in (develop)"`

Wait for all pods with a following label to enter 'Ready' or 'Error' state:

`kwait pod-we -lapp=example`

Wait for at least one pod to enter the 'Ready' state, even when the other ones are in 'Error' state:

`kwait pod-wr -lapp=example`

Wait for all the pods in that job to have a 'Succeeded' state:

`kwait job examplejob`

Wait for all the pods in that job to have a 'Succeeded' or 'Failed' state:

`kwait job-we examplejob`

Wait for at least one pod in that job to have 'Succeeded' state, does not mind some 'Failed' ones:

`kwait job-wr examplejob`

The program will wait until the resource exists and is ready or until it times out or fails a condition.

## Supported Resource Types

The kwait program supports the following Kubernetes resource types:

* 🚀 Pods
* 🚀 Services
* 🚀 Deployments
* 🚀 Jobs

To wait for a resource of a specific type, use the -kind flag followed by the name of the resource type (e.g., pod, service, etc.).

### Library

The `kwait` package can also be used as a library to wait for Kubernetes resources in your own Go programs. To use the package, simply import it and create a `k8s.Resource` object with the appropriate `k8s.ResourceWatcher` and `k8s.ResourceReadyChecker` implementations.

```go
import (
    "github.com/devlincashman/kwait/internal/k8s"
    "k8s.io/apimachinery/pkg/runtime"
)

func main() {
    // Create a Kubernetes client
    clientset, err := k8s.NewClient()
    if err != nil {
        panic(err)
    }

    // Create a DeploymentResource
    deploymentResource := k8s.NewDeploymentResource(clientset, "my-namespace")

    // Create a Resource object
    resource := k8s.Resource{
        Watcher:      deploymentResource,
        ReadyChecker: deploymentResource,
    }

    // Wait for the deployment to become ready
    err = resource.WaitForReady("my-namespace", "my-deployment", "")
    if err != nil {
        panic(err)
    }
}
```
