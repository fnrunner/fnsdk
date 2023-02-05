package fn

import (
	"fmt"
	"strings"
)

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

func (i Result) Error() string {
	return (i).String()
}

// String provides a human-readable message for the result item
func (i Result) String() string {
	identifier := i.ResourceRef
	var idStringList []string
	if identifier != nil {
		if identifier.APIVersion != "" {
			idStringList = append(idStringList, identifier.APIVersion)
		}
		if identifier.Kind != "" {
			idStringList = append(idStringList, identifier.Kind)
		}
		if identifier.Namespace != "" {
			idStringList = append(idStringList, identifier.Namespace)
		}
		if identifier.Name != "" {
			idStringList = append(idStringList, identifier.Name)
		}
	}
	formatString := "[%s]"
	severity := i.Severity
	// We default Severity to Info when converting a result to a message.
	if i.Severity == "" {
		severity = Info
	}
	list := []interface{}{severity}
	if len(idStringList) > 0 {
		formatString += " %s"
		list = append(list, strings.Join(idStringList, "/"))
	}
	formatString += ": %s"
	list = append(list, i.Message)
	return fmt.Sprintf(formatString, list...)
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

func (r *Results) String() string {
	var results []string
	for _, result := range *r {
		results = append(results, result.String())
	}
	return strings.Join(results, "\n---\n")
}

// Error enables Results to be returned as an error
func (r Results) Error() string {
	var msgs []string
	for _, i := range r {
		msgs = append(msgs, i.String())
	}
	return strings.Join(msgs, "\n\n")
}
