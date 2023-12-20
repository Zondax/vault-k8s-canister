package admctrl_injector

import (
	"context"
	"fmt"
	common2 "github.com/zondax/vault-k8s-canister/operator/common"

	"github.com/samber/lo"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createServiceAccount(serviceAccountName, namespace string) (*corev1.ServiceAccount, error) {
	_, clientset := common2.GetKubernetesClients()

	// Create the ServiceAccount
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceAccountName,
		},
	}

	svcAcc, err := clientset.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), serviceAccount, metav1.CreateOptions{})
	if err != nil {
		if isAlreadyExistsError(err) {
			svcAcc, err = clientset.CoreV1().ServiceAccounts(namespace).Update(context.TODO(), serviceAccount, metav1.UpdateOptions{})
			if err != nil {
				return nil, err
			}

			zap.S().Infof("[ADM CTRL] ServiceAccount '%s' updated in namespace '%s'", serviceAccountName, namespace)
			return svcAcc, err
		}

		return nil, err
	}

	zap.S().Infof("[ADM CTRL] ServiceAccount '%s' created in namespace '%s'", serviceAccountName, namespace)
	return svcAcc, nil
}

func createRole(roleName string, namespace string, resNames []string) error {
	_, clientset := common2.GetKubernetesClients()

	// Define the Role that grants read permissions across the namespace
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name: roleName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{""},
				Resources:     []string{"secrets"},
				ResourceNames: lo.Map(resNames, func(item string, idx int) string { return fmt.Sprintf("%s-secret", item) }),
				Verbs:         []string{"get", "list", "watch", "create", "update", "patch"},
			},
			{
				APIGroups:     []string{common2.TororuResourcesGroup},
				Resources:     []string{common2.TororuResourceNamePlural},
				ResourceNames: resNames,
				Verbs:         []string{"get", "list", "update", "patch"},
			},
		},
	}
	_, err := clientset.RbacV1().Roles(namespace).Create(context.TODO(), role, metav1.CreateOptions{})
	if err != nil {
		if isAlreadyExistsError(err) {
			_, err := clientset.RbacV1().Roles(namespace).Update(context.TODO(), role, metav1.UpdateOptions{})
			if err != nil {
				return err
			}

			zap.S().Infof("[ADM CTRL] Role '%s' updated", roleName)
			return nil
		}

		return err
	}

	zap.S().Infof("[ADM CTRL] Role '%s' created", roleName)
	return nil
}

func createRoleBinding(roleBindingName, serviceAccountName, serviceAccountNamespace, roleName string) error {
	// Create the RoleBinding to associate the ServiceAccount with the Role
	RoleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: roleBindingName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: serviceAccountNamespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     roleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	_, clientset := common2.GetKubernetesClients()
	_, err := clientset.RbacV1().RoleBindings(serviceAccountNamespace).Create(context.TODO(), RoleBinding, metav1.CreateOptions{})
	if err != nil {
		if isAlreadyExistsError(err) {
			_, err := clientset.RbacV1().RoleBindings(serviceAccountNamespace).Update(context.TODO(), RoleBinding, metav1.UpdateOptions{})
			if err != nil {
				return err
			}

			zap.S().Infof("[ADM CTRL] RoleBinding '%s' updated", roleBindingName)
			return nil
		}

		return err
	}

	zap.S().Infof("[ADM CTRL] RoleBinding '%s' created", roleBindingName)
	return nil
}

// isAlreadyExistsError checks if the error indicates that a resource already exists.
func isAlreadyExistsError(err error) bool {
	if statusErr, ok := err.(*errors.StatusError); ok {
		return statusErr.ErrStatus.Reason == metav1.StatusReasonAlreadyExists
	}
	return false
}

func prepServiceAccountForPod(_ctx context.Context, pod corev1.Pod, scConfs []*SidecarConfig) (*corev1.ServiceAccount, error) {
	// TODO: clean service accounts for pods that do not exist anymore
	// Should we clean service accounts reactively or periodically?

	serviceAccountName := fmt.Sprintf("%s-svc-acc", pod.Name)
	svcAccount, err := createServiceAccount(serviceAccountName, pod.Namespace)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	tResNames := lo.FlatMap(scConfs, func(item *SidecarConfig, idx int) []string { return item.GetTororuResourceNames() })
	roleName := fmt.Sprintf("%s-role", pod.Name)
	err = createRole(roleName, pod.Namespace, tResNames)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	RoleBindingName := fmt.Sprintf("%s-role-bind", pod.Name)
	err = createRoleBinding(RoleBindingName, serviceAccountName, pod.Namespace, roleName)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}

	return svcAccount, nil
}
