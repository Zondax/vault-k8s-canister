package operator

import (
	"context"
	"github.com/zondax/tororu-operator/operator/common"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func shouldRestart(ctx context.Context, pod *corev1.Pod) bool {
	logger := ctx.Value(common.ContextLoggerKey).(*zap.SugaredLogger)
	var check1, check2 bool

	// Check if the sidecar already exists
	if val, ok := pod.Annotations[common.TororuResourceReqRWAnnotation]; ok && len(val) > 0 {
		for _, container := range pod.Spec.Containers {
			if container.Name == common.SidecarName {
				logger.Infof("Sidecar already exists on: %s", pod.Name)
				check1 = true
			}
		}
	} else {
		check1 = true
	}

	if val, ok := pod.Annotations[common.TororuResourceReqROAnnotation]; ok && len(val) > 0 {
		if pod.Annotations[common.TororuEnvVarsAppliedAnnotation] == common.TororuAnnotationsTrueString {
			logger.Infof("EnvVarsAppliedAnnotaion found on pod: %s", pod.Name)
			check2 = true
		}
	} else {
		check2 = true
	}

	return !(check1 && check2)
}

func (o SidecarOperator) restartPod(ctx context.Context, pod *corev1.Pod) error {
	logger := ctx.Value(common.ContextLoggerKey).(*zap.SugaredLogger)
	if pod.Status.Phase == corev1.PodRunning {
		// Check if the sidecar can be created using the policy in the canister
		// Always recreating for now
		logger.Warnf("Sidecar/Env missing. Restarting pod %s", pod.Name)
		// Deleting pod so it gets recreated and admission controller adds the correct sidecar to it
		err := o.mgr.GetClient().Delete(ctx, pod)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func (o SidecarOperator) enqueueExistingPodsWithAnnotation(ctx context.Context) error {
	logger := ctx.Value(common.ContextLoggerKey).(*zap.SugaredLogger)
	podList := &corev1.PodList{}
	listOptions := []client.ListOption{
		client.MatchingLabels{common.TororuManagedAnnotation: common.TororuAnnotationsTrueString},
	}

	if err := o.mgr.GetClient().List(ctx, podList, listOptions...); err != nil {
		return err
	}

	for _, pod := range podList.Items {
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      pod.Name,
				Namespace: pod.Namespace,
			},
		}
		logger.Info("Enqueuing existing pod with annotation", "pod", req.NamespacedName)
		result, err := o.Reconcile(ctx, req)
		if err != nil {
			logger.Warnf("Failed to reconcile: %v", err)
		} else {
			logger.Infof("Reconciled: %v", result)
		}
	}

	return nil
}
