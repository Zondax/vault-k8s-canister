package sidecarPostgres

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zondax/tororu-operator/operator/common"
	v1 "github.com/zondax/tororu-operator/operator/common/v1"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type tResInfo struct {
	Name             string
	RotationDuration time.Duration
	commChan         chan string
	dynamicClient    *dynamic.DynamicClient
	kubeClient       *kubernetes.Clientset
}

func (t *tResInfo) shouldUpdate(tResSecret *corev1.Secret) bool {
	value := string(tResSecret.Data["value"])

	if value == "" {
		return true
	}

	managedFields := tResSecret.ObjectMeta.ManagedFields
	sideCarManagedField := metav1.ManagedFieldsEntry{}

	for _, managedField := range managedFields {
		if managedField.Manager == "sidecars" {
			sideCarManagedField = managedField
			break
		}
	}

	if sideCarManagedField.Manager == "" {
		return true
	}

	timeSinceUpdate := time.Since(sideCarManagedField.Time.Time)

	return timeSinceUpdate > t.RotationDuration
}

func (t *tResInfo) rotateAndUpdateOnce() error {
	ctx := context.TODO()
	nsName, err := common.GetNamespacedNameFromNameString(t.Name)
	if err != nil {
		t.commChan <- fmt.Sprintf("Counldn't convert namespaced name %s", t.Name)
		return err
	}

	tRes, err := common.GetTResFromName(ctx, nsName.Namespace, nsName.Name)
	if err != nil {
		t.commChan <- fmt.Sprintf("Failed to get %s", t.Name)
		return err
	}

	k8secret, err := t.kubeClient.CoreV1().Secrets(nsName.Namespace).Get(ctx, tRes.Name+"-secret", metav1.GetOptions{})
	if err != nil {
		t.commChan <- fmt.Sprintf("Failed to get secret with error %v", err)
		return err
	}

	if !t.shouldUpdate(k8secret) {
		t.commChan <- fmt.Sprintf("Won't update %s as recently updated", tRes.Name)
		return nil
	}

	// Update the secret on the resource
	// Generate postgressSecret
	postgressSecret, err := t.updateSecretOnPostgres(tRes)
	if err != nil {
		zap.S().Error(err)
		return err
	}

	k8secret.StringData = map[string]string{"value": postgressSecret}
	_, err = t.
		kubeClient.CoreV1().Secrets(nsName.Namespace).Update(ctx, k8secret, metav1.UpdateOptions{})
	if err != nil {
		t.commChan <- fmt.Sprintf("Failed to update secret with error %v", err)
		return err
	}

	// -- Updates --
	tRes.LastUpdated = time.Now().UTC().Format(http.TimeFormat)
	tRes.PodsRestartRequired = true
	// -- Updates end --

	tResRaw, err := tRes.ToUnstructured()
	if err != nil {
		t.commChan <- fmt.Sprintf("Failed convert updated to unstructures error %v", err)
		return err
	}

	if tResRaw.Object != nil {
		_, err = t.
			dynamicClient.
			Resource(v1.SchemeGroupVersion.WithResource(common.TororuResourceNamePlural)).
			Namespace(nsName.Namespace).
			Update(ctx, tResRaw, metav1.UpdateOptions{})

		if err != nil {
			t.commChan <- fmt.Sprintf("Failed to update crd with error %v", err)
			return err
		}

		t.commChan <- fmt.Sprintf("Updated %s", t.Name)
	}

	return nil
}

func (t *tResInfo) rotateAndUpdateForever() {
	for {
		err := t.rotateAndUpdateOnce()
		if err != nil {
			zap.S().Error(err)
		}

		zap.S().Infof("Sleeping for %s", t.RotationDuration)
		time.Sleep(t.RotationDuration)
	}
}
