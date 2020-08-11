package horizontalautoscaler

import (
	"fmt"

	"github.com/ellistarn/karpenter/pkg/api/v1alpha1"
	cloudprovider "github.com/ellistarn/karpenter/pkg/cloudprovider"
	"go.uber.org/zap"
)

// Autoscaler applies a HorizontalAutoscaler policy to a NodeGroup
type Autoscaler struct {
	NodeGroup            cloudprovider.NodeGroup
	HorizontalAutoscaler *v1alpha1.HorizontalAutoscaler
}

// Reconcile executes a reconcilation loop for the Autoscaler's NodeGroup using the HorizontalAutoscaler policy
func (a *Autoscaler) Reconcile() error {
	zap.S().Infof("Executing reconciliation loop for %s.", a.HorizontalAutoscaler.ObjectMeta.SelfLink)
	metrics, err := a.ReadMetrics()
	if err != nil {
		return fmt.Errorf("while reading metrics, %v", err)
	}
	desiredReplicas, err := a.CalculateDesiredReplicas(metrics)
	if err != nil {
		return fmt.Errorf("while calculating desired replicas %v", err)
	}
	if err := a.SetReplicas(desiredReplicas); err != nil {
		return fmt.Errorf("while setting replicas %v", err)
	}
	return nil
}

// ReadMetrics for the NodeGroup
func (a *Autoscaler) ReadMetrics() (map[string]string, error) {
	zap.S().Infof("Reading metrics for %s.", a.HorizontalAutoscaler.ObjectMeta.SelfLink)
	return nil, nil
}

// CalculateDesiredReplicas for the NodeGroup
func (a *Autoscaler) CalculateDesiredReplicas(metrics map[string]string) (int, error) {
	zap.S().Infof("Calculating desired replicas for %s.", a.HorizontalAutoscaler.ObjectMeta.SelfLink)
	return 0, nil
}

// SetReplicas of the NodeGroup
func (a *Autoscaler) SetReplicas(replicas int) error {
	zap.S().Infof("Setting desired replicas of %s to %d.", a.HorizontalAutoscaler.ObjectMeta.SelfLink, 0)
	if err := a.NodeGroup.SetReplicas(replicas); err != nil {
		return fmt.Errorf("while setting replicas on node group %s, %v", a.NodeGroup.Name(), err)
	}
	return nil
}
