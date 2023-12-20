package operator

import (
	"context"
	"fmt"
	common2 "github.com/zondax/tororu-operator/operator/common"
	"github.com/zondax/tororu-operator/operator/common/icp"
	"github.com/zondax/tororu-operator/operator/common/v1"
	"net/http"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// UpdateTororuResource updates a TororuResource in the Kubernetes cluster using the provided client.
// It takes a context, a client for managing Kubernetes resources, and a pointer to the TororuResource to be updated.
// The function returns an error if the update operation encounters any issues.
func UpdateTororuResource(ctx context.Context, mgrClient client.Client, tRes *v1.TororuResource) error {
	logger := ctx.Value(common2.ContextLoggerKey).(*zap.SugaredLogger)

	logger.Infof("Updating CRD: %s", tRes.Name)

	// Update the TororuResource using the provided client
	err := mgrClient.Update(ctx, tRes)

	logger.Infof("Updated CRD: %s", tRes.Name)

	return err
}

// RefreshTororuResource updates a TororuResource based on CRD status and other parameters.
// It takes a context and a pointer to the TororuResource object as input and returns an error if any operation fails.
func (o CRDOperator) RefreshTororuResource(ctx context.Context, tRes *v1.TororuResource) error {
	logger := ctx.Value(common2.ContextLoggerKey).(*zap.SugaredLogger)
	client := o.mgr.GetClient()

	// TODO: Define how the ID gets generated after CRD approved

	// Get the status of the CRD
	crdStatus, pendingReqs, err := common2.GetCRDStatus(common2.GetPodOrCRDId(tRes.Name, tRes.Namespace))
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Debugf("CRD Status: %v", crdStatus)

	if crdStatus == nil {
		skipCreate := false
		for _, req := range pendingReqs {
			if req.PendingType == icp.PendingOperationEnumCreate {
				skipCreate = true
			}
		}

		if !skipCreate {
			err = common2.CreateCRDOnboardRequest(common2.GetPodOrCRDId(tRes.Name, tRes.Namespace), uint32(tRes.Spec.Rotate))
			if err != nil {
				logger.Error(err)
				return err
			}
		}
	} else if int(crdStatus.TTL) != tRes.Spec.Rotate {
		skipUpdate := false
		for _, req := range pendingReqs {
			if req.PendingType == icp.PendingOperationEnumUpdate && req.Ttl == uint32(tRes.Spec.Rotate) {
				skipUpdate = true
			}
		}

		if !skipUpdate {
			err = common2.UpdateSecret(ctx, common2.GetPodOrCRDId(tRes.Name, tRes.Namespace), uint32(tRes.Spec.Rotate))
			if err != nil {
				logger.Error(err)
				return err
			}
		}
	}

	// Update the TororuResource with CRD status and other information
	tRes.LastUpdated = time.Now().UTC().Format(http.TimeFormat)

	// Create the secret based on the TororuResource name if doesn't exist
	secretNamespacedName := fmt.Sprintf("%s-secret", common2.GetPodOrCRDId(tRes.Name, tRes.Namespace))
	tRes.Secret = secretNamespacedName
	// Get the secret associated with the TororuResource
	var k8Secret corev1.Secret
	err = client.Get(
		ctx,
		types.NamespacedName{
			Namespace: tRes.Namespace,
			Name:      tRes.Name + "-secret",
		},
		&k8Secret,
	)
	if err != nil {
		// create secret
		err = client.Create(ctx, &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      tRes.Name + "-secret",
				Namespace: tRes.Namespace,
			},
			StringData: map[string]string{"value": ""},
		})

		if err != nil {
			logger.Error(err)
			return err
		}
	}

	if crdStatus == nil {
		tRes.Consumers = v1.TororuResourceConsumers{
			RW: "",
			RO: []string{},
		}
		tRes.Approved = false

		return nil
	}

	tRes.Consumers = v1.TororuResourceConsumers{
		RW: crdStatus.Consumers.RW,
		RO: crdStatus.Consumers.RO,
	}
	tRes.Approved = crdStatus.Status == icp.PendingOperationEnumNone
	tRes.Spec.Rotate = int(crdStatus.TTL)

	// Update the TororuResource in the Kubernetes API server
	err = UpdateTororuResource(ctx, client, tRes)
	return err

}

// restartPod restarts a given Pod by deleting it and allowing it to be recreated.
// It takes a context and a pointer to the Pod object as input and returns an error if any operation fails.
func (o CRDOperator) restartPod(ctx context.Context, pod *corev1.Pod) error {
	logger := ctx.Value(common2.ContextLoggerKey).(*zap.SugaredLogger)

	// Check if the Pod is in the Running phase
	if pod.Status.Phase == corev1.PodRunning {
		// Deleting the Pod to trigger a recreation, allowing the admission controller to add the correct sidecar
		err := o.mgr.GetClient().Delete(ctx, pod)
		if err != nil {
			// Log any errors encountered during deletion
			logger.Error(err)
			return err
		}
	}

	// Return nil if the Pod has been successfully restarted or didn't require restarting
	return nil
}

func (o CRDOperator) restartPodFromName(ctx context.Context, podName string) error {
	logger := ctx.Value(common2.ContextLoggerKey).(*zap.SugaredLogger)

	// Convert the pod name string to a NamespacedName
	nsName, err := common2.GetNamespacedNameFromNameString(podName)
	if err != nil {
		logger.Errorf("failed to find pod %s", podName)
		return err
	}

	// Initialize a variable to hold the Pod object
	pod := &corev1.Pod{}

	// Retrieve the Pod object from the Kubernetes API server
	err = o.mgr.GetClient().Get(ctx, *nsName, pod, &client.GetOptions{})
	if err != nil {
		logger.Errorf("failed to get pod %s", podName)
		return err
	}

	logger.Infof("Restarting Pod: %s in Namespace: %s", pod.Name, pod.Namespace)

	// Restart the Pod
	err = o.restartPod(ctx, pod)
	if err != nil {
		logger.Errorf("Error restarting Pod: %s in Namespace: %s - %v", pod.Name, pod.Namespace, err)
		return err
	}

	logger.Infof("Pod restarted: %s in Namespace: %s", pod.Name, pod.Namespace)
	return nil
}

// RestartROPods restarts Read-Only (RO) pods associated with a TororuResource.
// It takes a context, tRes, and returns an error if any operation fails.
func (o CRDOperator) RestartPods(ctx context.Context, pods []string) error {

	// fmt.Println("Pods to restart:", pods)

	// Iterate through the list of pod names
	for _, podName := range pods {
		if err := o.restartPodFromName(ctx, podName); err != nil {
			return err
		}
	}

	// Return nil if all pods have been successfully restarted
	return nil
}
