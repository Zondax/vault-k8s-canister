package operator

import (
	"context"
	"github.com/zondax/vault-k8s-canister/operator/common"
	"github.com/zondax/vault-k8s-canister/operator/internal/conf"
	manager2 "github.com/zondax/vault-k8s-canister/operator/internal/k8s/manager"
	"os"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	k8sManager "sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type SidecarOperator struct {
	name   string
	config *conf.Config
	mgr    k8sManager.Manager
}

func NewSidecarOperator(config *conf.Config) *SidecarOperator {
	return &SidecarOperator{
		name:   "sidecar-operator",
		config: config,
	}
}

func (o SidecarOperator) Name() string {
	return o.name
}

func (o SidecarOperator) Start() error {
	managerKind := "sidecar"
	// Create a new manager to provide shared dependencies and start components
	mgr, err := manager2.NewManager(o.config, managerKind, ":8080")
	if err != nil {
		zap.S().Fatalf("[Operator Sidecar] Failed to create new manager: %v", err)
	}
	o.mgr = mgr

	err = builder.
		ControllerManagedBy(mgr).
		WithEventFilter(podAnnotationPredicate()).
		For(&corev1.Pod{}).
		Complete(o)
	if err != nil {
		return err
	}

	err = mgr.Add(manager2.NewCallbackInitializer(o.Initialize))
	if err != nil {
		return err
	}

	// Now we can start the manager, and block until it is stopped
	// This will run all the controllers, and start any reconciliations
	zap.S().Infof("[Operator Sidecar] Starting %s", managerKind)
	if err := mgr.Start(context.Background()); err != nil {
		zap.S().Fatalf("[Operator Sidecar] Failed to start manager: %v", err)
	}

	return nil
}

func (o SidecarOperator) Initialize(ctx context.Context) error {
	logger := zap.S().With("pod", "initialize")
	ctx = context.WithValue(ctx, common.ContextLoggerKey, logger)

	if err := o.enqueueExistingPodsWithAnnotation(ctx); err != nil {
		zap.S().Error(err, "[Operator Sidecar] Failed to enqueue existing Pods with annotation")
		os.Exit(1)
	}
	return nil
}

func (o SidecarOperator) Stop() error {
	zap.S().Warnf("[Operator Sidecar] Stop not implemented yet")
	return nil
}

// Reconcile handles the reconciliation of events
func (o SidecarOperator) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger := zap.S().With("pod", req.NamespacedName)
	ctx = context.WithValue(ctx, common.ContextLoggerKey, logger)

	// Fetch the Pod instance
	pod := &corev1.Pod{}
	err := o.mgr.GetClient().Get(ctx, req.NamespacedName, pod)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Error(err, "unable to fetch. Namespace not found")
			return reconcile.Result{}, nil
		}
		logger.Error(err, "unable to fetch Pod")
		return reconcile.Result{}, err
	}

	// Inject sidecar
	if shouldRestart(ctx, pod) {
		if err := o.restartPod(ctx, pod); err != nil {
			logger.Error(err, "unable to restart sidecars")
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
