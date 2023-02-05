package main

import (
	"context"
	"os"

	"github.com/fnrunner/fnsdk/go/fn"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

var _ fn.Runner = &implA{}

type implA struct{}

func main() {
	ctx := context.TODO()
	if err := fn.AsMain(fn.WithContext(ctx, &implA{})); err != nil {
		os.Exit(1)
	}
}

func (r *implA) Run(ctx *fn.Context, functionConfig map[string]runtime.RawExtension, resources *fn.Resources, results *fn.Results) bool {

	res := &unstructured.Unstructured{}
	res.SetAPIVersion("a.b.c/v1alpha1")
	res.SetKind("A")
	res.SetName("implA")
	res.SetNamespace("default")
	resources.AddResource(res, &fn.ResourceParameters{Conditioned: true, Internal: true})
	return true
}
