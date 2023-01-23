package fn

import (
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
)

const (
	Group                  = "fnctrlr.nephio.org"
	version                = "v1alpha1"
	ConditionedResourceKey = Group + "/" + "conditioned"
	ExternalResourceKey    = Group + "/" + "external"
)

// ResourceContextProcessor is implemented by configuration functions built with this framework
type ResourceContextProcessor interface {
	Process(fnCtx *ResourceContext) (bool, error)
}

// ResourceContextProcessorFunc converts a compatible function to a ResourceContextProcessor.
type ResourceContextProcessorFunc func(fnCtx *ResourceContext) (bool, error)

func (p ResourceContextProcessorFunc) Process(fnCtx *ResourceContext) (bool, error) {
	return p(fnCtx)
}

type ResourceContext struct {
	// fnconfig provides additional configuration for the function
	FunctionConfig map[string]runtime.RawExtension `json:"functionConfig,omitempty" yaml:"functionConfig,omitempty"`
	// Resources contain the resource on which this function operates
	Resources map[string][]runtime.RawExtension `json:"resources" yaml:"resources"`
	// results provide a structured
	Results *Results `json:"results,omitempty" yaml:"results,omitempty"`
}

func ParseResourceContext(input []byte) (*ResourceContext, error) {
	rCtx := &ResourceContext{}
	if err := json.Unmarshal(input, rCtx); err != nil {
		return nil, fmt.Errorf("error: %s with bytes: %s", err.Error(), string(input))
	}
	return rCtx, nil
}
