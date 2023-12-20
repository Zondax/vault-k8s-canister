package operator

import (
	"github.com/zondax/tororu-operator/operator/common"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func shouldReconcile(annotations map[string]string) bool {
	if val, ok := annotations[common.TororuManagedAnnotation]; ok && val == common.TororuAnnotationsTrueString {
		if val, ok := annotations[common.TororuResourceReqROAnnotation]; ok && len(val) > 0 {
			return true
		}

		if val, ok := annotations[common.TororuResourceReqRWAnnotation]; ok && len(val) > 0 {
			return true
		}
	}

	return true
}

func podAnnotationPredicate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			annotations := e.Object.GetAnnotations()
			return shouldReconcile(annotations)
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			annotations := e.ObjectNew.GetAnnotations()
			return shouldReconcile(annotations)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return false
		},
	}
}
