package k8s

import "k8s.io/apimachinery/pkg/runtime/schema"

var CockroachDBSchema = schema.GroupVersionResource{
	Group:    "ponglehub.co.uk",
	Version:  "v1alpha1",
	Resource: "cockroachdb",
}
