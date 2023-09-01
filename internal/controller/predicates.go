package controller

import (
	"github.com/tydanny/myapp-operator/internal/wdk"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var MyAppLabelPredicate = predicate.NewPredicateFuncs(func(o client.Object) bool {
	return o.GetLabels()[wdk.App] == wdk.MyApp
})
