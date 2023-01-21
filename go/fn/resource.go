package fn

import (
	"encoding/json"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// An Object is a Kubernetes object.
type Object interface {
	metav1.Object
	runtime.Object
}

type Resources struct {
	Resources map[string][]runtime.RawExtension `json:"resources" yaml:"resources"`
}

func (r *Resources) AddResource(o Object) error {
	gvkString := GetGVKString(o)
	_, ok := r.Resources[gvkString]
	if !ok {
		r.Resources[gvkString] = []runtime.RawExtension{}
	}
	b, err := json.Marshal(o)
	if err != nil {
		return err
	}
	present, idx := isPresent(r.Resources[gvkString], o)
	if !present {
		r.Resources[gvkString] = append(r.Resources[gvkString], runtime.RawExtension{Raw: b})
	} else {
		// overwrite
		r.Resources[gvkString][idx] = runtime.RawExtension{Raw: b}
	}
	return nil
}

func (r *Resources) AddConditionalResource(o Object) error {
	gvkString := GetGVKString(o)
	_, ok := r.Resources[gvkString]
	if !ok {
		r.Resources[gvkString] = []runtime.RawExtension{}
	}
	labels := o.GetLabels()
	if len(labels) == 0 {
		labels = map[string]string{}
	}
	labels[ConditionedResourceKey] = "true"
	o.SetLabels(labels)
	b, err := json.Marshal(o)
	if err != nil {
		return err
	}
	present, idx := isPresent(r.Resources[gvkString], o)
	if !present {
		r.Resources[gvkString] = append(r.Resources[gvkString], runtime.RawExtension{Raw: b})
	} else {
		// overwrite
		r.Resources[gvkString][idx] = runtime.RawExtension{Raw: b}
	}
	return nil
}

func GetGVKString(o Object) string {
	gvk := o.GetObjectKind().GroupVersionKind()
	return GVKToString(&gvk)
}

func GVKToString(gvk *schema.GroupVersionKind) string {
	return fmt.Sprintf("%s.%s.%s", gvk.Kind, gvk.Version, gvk.Group)
}

func StringToGVK(s string) *schema.GroupVersionKind {
	var gvk *schema.GroupVersionKind
	if strings.Count(s, ".") >= 2 {
		s := strings.SplitN(s, ".", 3)
		gvk = &schema.GroupVersionKind{Group: s[2], Version: s[1], Kind: s[0]}
	}
	return gvk
}

func isPresent(slice []runtime.RawExtension, o Object) (bool, int) {
	for idx, v := range slice {
		u := &unstructured.Unstructured{}
		if err := json.Unmarshal(v.Raw, u); err != nil {
			return false, 0
		}
		if u.GetName() == o.GetName() && u.GetNamespace() == o.GetNamespace() {
			return true, idx
		}
	}
	return false, 0
}
