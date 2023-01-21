package fn

import "fmt"

// Severity indicates the severity of the Result
type Severity string

const (
	// Error indicates the result is an error.  Will cause the function to exit non-0.
	Error Severity = "error"
	// Info indicates the result is an informative message
	Info Severity = "info"
)

type Results []*Result

// Result defines a result for the fucntion execution
type Result struct {
	// Message is a human readable message. This field is required.
	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	// Severity is the severity of this result
	Severity Severity `yaml:"severity,omitempty" json:"severity,omitempty"`

	// ResourceRef is a reference to a resource.
	// Required fields: apiVersion, kind, name.
	ResourceRef *ResourceRef `json:"resourceRef,omitempty" yaml:"resourceRef,omitempty"`
}

// ResourceRef fills the ResourceRef field in Results
type ResourceRef struct {
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace  string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

func (r *Results) Errorf(format string, a ...any) {
	errResult := &Result{Severity: Error, Message: fmt.Sprintf(format, a...)}
	*r = append(*r, errResult)
}

func (r *Results) ErrorE(err error) {
	errResult := &Result{Message: err.Error()}
	*r = append(*r, errResult)
}

// Infof writes an Info level `result` to the results slice. It accepts arguments according to a format specifier.
func (r *Results) Infof(format string, a ...any) {
	infoResult := &Result{Severity: Info, Message: fmt.Sprintf(format, a...)}
	*r = append(*r, infoResult)
}