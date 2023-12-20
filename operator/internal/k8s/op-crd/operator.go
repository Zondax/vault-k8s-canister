package operator

import (
	"context"
	"github.com/samber/lo"
	"github.com/zondax/vault-k8s-canister/operator/common"
	v12 "github.com/zondax/vault-k8s-canister/operator/common/v1"
	"github.com/zondax/vault-k8s-canister/operator/internal/conf"
	manager2 "github.com/zondax/vault-k8s-canister/operator/internal/k8s/manager"
	"time"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	k8sManager "sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type CRDOperator struct {
	name   string
	config *conf.Config
	mgr    k8sManager.Manager
}

func NewCRDOperator(config *conf.Config) *CRDOperator {
	return &CRDOperator{
		name:   "crd-operator",
		config: config,
	}
}

func (o CRDOperator) Name() string {
	return o.name
}

func (o CRDOperator) Start() error {
	// Create a new manager to provide shared dependencies and start components
	managerKind := "crd"
	mgr, err := manager2.NewManager(o.config, managerKind, ":8181")
	if err != nil {
		zap.S().Fatalf("[CRD Operator] Failed to create new manager: %v", err)
	}
	o.mgr = mgr

	err = v12.AddToScheme(scheme.Scheme)
	if err != nil {
		zap.S().Fatal(err)
	}

	err = builder.
		ControllerManagedBy(mgr).
		// WithEventFilter(podAnnotationPredicate()).
		For(&v12.TororuResource{}).
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
	zap.S().Infof("[CRD Operator] Starting %s", managerKind)
	if err := mgr.Start(context.Background()); err != nil {
		zap.S().Errorf("[CRD Operator] Failed to start manager: %v", err)
	}

	return nil
}

func (o CRDOperator) Initialize(ctx context.Context) error {
	return nil
}

func (o CRDOperator) Stop() error {
	zap.S().Warnf("[CRD Operator] Stop not implemented yet")
	return nil
}

// Reconcile handles the reconciliation of events for a TororuResource.
func (o CRDOperator) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	reconcileResult := reconcile.Result{RequeueAfter: 30 * time.Second}

	// Create a logger specific to the CRD being reconciled
	logger := zap.S().With("crd", req.NamespacedName)
	ctx = context.WithValue(ctx, common.ContextLoggerKey, logger)

	logger.Info("---------------------------------------")
	logger.Info("--------Reconciliation started --------")

	// Fetch the TororuResource associated with the request
	status := &v12.TororuResource{}
	if err := o.mgr.GetClient().Get(ctx, req.NamespacedName, status); err != nil {
		if errors.IsNotFound(err) {
			logger.Error(err, "TororuResource not found. Namespace does not exist")
			return reconcile.Result{}, nil // No action needed
		}
		logger.Error(err, "Failed to fetch TororuResource")
		return reconcileResult, err
	}

	newStatus := &v12.TororuResource{}
	if err := o.mgr.GetClient().Get(ctx, req.NamespacedName, newStatus); err != nil {
		logger.Error(err, "Failed to fetch TororuResource")
		return reconcileResult, err
	}

	if err := o.RefreshTororuResource(ctx, newStatus); err != nil {
		logger.Error(err, "Failed to update TororuResource on creation")
	}

	podsToRestart, _ := lo.Difference(status.Consumers.RO, newStatus.Consumers.RO)
	logger.Debugf("restart required on RO pods: %s", podsToRestart)

	if err := o.RestartPods(ctx, podsToRestart); err != nil {
		logger.Errorf("Failed to restart pods removed from current approved list: %v", err)
	}

	rwPodRestartRequired := status.Consumers.RW != "" && newStatus.Consumers.RW == ""
	logger.Debugf("restart required on RW pod [%t] -> old status: %s - new status %s", rwPodRestartRequired, status.Consumers.RW, newStatus.Consumers.RW)
	if rwPodRestartRequired {
		if err := o.restartPodFromName(ctx, status.Consumers.RW); err != nil {
			logger.Errorf("Failed to restart RW pods: %v", err)
		}
	}

	restartAllByTtlChange := status.Spec.Rotate != newStatus.Spec.Rotate
	if restartAllByTtlChange {
		if err := o.restartPodFromName(ctx, status.Consumers.RW); err != nil {
			logger.Errorf("Failed to restart RW pods because of TTL change: %v", err)
		}
		if err := o.RestartPods(ctx, status.Consumers.RO); err != nil {
			logger.Errorf("Failed to restart RO pods because of TTL change: %v", err)
		}
	}

	// Check if Pod restart is required
	if newStatus.PodsRestartRequired {
		if err := o.RestartPods(ctx, status.Consumers.RO); err != nil {
			logger.Errorf("Failed to restart Read-Only (RO) pods: %v", err)
		} else {
			// Reset the PodsRestartRequired flag and update the TororuResource
			newStatus.PodsRestartRequired = false
			if err := UpdateTororuResource(ctx, o.mgr.GetClient(), newStatus); err != nil {
				logger.Errorf("Failed to update TororuResource after restarting RO pods: %v", err)
			}
		}
	}

	logger.Info("--------Reconciliation ended --------")
	logger.Info("-------------------------------------")
	return reconcileResult, nil
}
