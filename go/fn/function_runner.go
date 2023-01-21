package fn

import "k8s.io/apimachinery/pkg/runtime"

type Runner interface {
	// Run provides the entrypoint to allow you to process the resources. any crud operation
	// Args:
	//    fnConfig: the configuration parameters of the function
	//    resources: The KRM resources.
	//    results: You can use `ErrorE` `Errorf` `Infof` to add user message to `Results`.
	// Returns:
	//    return a boolean to tell whether the execution should be considered as PASS or FAIL.
	Run(context *Context, fnConfig map[string]runtime.RawExtension, resources *Resources, results *Results) bool
}
