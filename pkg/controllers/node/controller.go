/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package node

import (
	"context"
	"fmt"
	"time"

	"github.com/awslabs/karpenter/pkg/apis/provisioning/v1alpha3"
	"github.com/awslabs/karpenter/pkg/utils/functional"
	"github.com/awslabs/karpenter/pkg/utils/result"

	"go.uber.org/multierr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/logging"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Now is a time.Now() that may be mocked by tests.
var Now = time.Now

// NewController constructs a controller instance
func NewController(kubeClient client.Client) *Controller {
	return &Controller{
		kubeClient: kubeClient,
		liveness:   &Liveness{kubeClient: kubeClient},
		emptiness:  &Emptiness{kubeClient: kubeClient},
		expiration: &Expiration{kubeClient: kubeClient},
	}
}

// Controller manages a set of properites on karpenter provisioned nodes, such as
// taints, labels, finalizers.
type Controller struct {
	kubeClient client.Client
	readiness  *Readiness
	liveness   *Liveness
	emptiness  *Emptiness
	expiration *Expiration
	finalizer  *Finalizer
}

// Reconcile executes a reallocation control loop for the resource
func (c *Controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	ctx = logging.WithLogger(ctx, logging.FromContext(ctx).Named("Node"))
	// 1. Retrieve Node, ignore if not provisioned or terminating
	stored := &v1.Node{}
	if err := c.kubeClient.Get(ctx, req.NamespacedName, stored); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return result.RetryIfError(ctx, err)
	}
	if _, ok := stored.Labels[v1alpha3.ProvisionerNameLabelKey]; !ok {
		return reconcile.Result{}, nil
	}
	if !stored.DeletionTimestamp.IsZero() {
		return reconcile.Result{}, nil
	}

	// 2. Retrieve Provisioner
	provisioner := &v1alpha3.Provisioner{}
	if err := c.kubeClient.Get(ctx, types.NamespacedName{Name: stored.Labels[v1alpha3.ProvisionerNameLabelKey]}, provisioner); err != nil {
		return result.RetryIfError(ctx, err)
	}

	// 3. Execute reconcilers
	node := stored.DeepCopy()
	var backoffs []time.Duration
	var errs error
	for _, reconciler := range []interface {
		Reconcile(context.Context, *v1alpha3.Provisioner, *v1.Node) (reconcile.Result, error)
	}{
		c.readiness,
		c.liveness,
		c.expiration,
		c.emptiness,
		c.finalizer,
	} {
		res, err := reconciler.Reconcile(ctx, provisioner, node)
		errs = multierr.Append(errs, err)
		backoffs = append(backoffs, res.RequeueAfter)
	}

	// 4. Patch any changes, regardless of errors
	if !equality.Semantic.DeepEqual(node, stored) {
		if err := c.kubeClient.Patch(ctx, node, client.MergeFrom(stored)); err != nil {
			return reconcile.Result{}, fmt.Errorf("patching node %s, %w", node.Name, err)
		}
	}
	// 5. Requeue if error or if retryAfter is set
	if errs != nil {
		return result.RetryIfError(ctx, errs)
	}
	return reconcile.Result{RequeueAfter: functional.MinDuration(backoffs...)}, nil
}

func (c *Controller) Register(ctx context.Context, m manager.Manager) error {
	return controllerruntime.
		NewControllerManagedBy(m).
		Named("Node").
		For(&v1.Node{}).
		Watches(
			// Reconcile all nodes related to a provisioner when it changes.
			&source.Kind{Type: &v1alpha3.Provisioner{}},
			handler.EnqueueRequestsFromMapFunc(func(o client.Object) (requests []reconcile.Request) {
				nodes := &v1.NodeList{}
				if err := c.kubeClient.List(ctx, nodes, client.MatchingLabels(map[string]string{v1alpha3.ProvisionerNameLabelKey: o.GetName()})); err != nil {
					logging.FromContext(ctx).Errorf("Failed to list nodes when mapping expiration watch events, %s", err.Error())
					return requests
				}
				for _, node := range nodes.Items {
					requests = append(requests, reconcile.Request{NamespacedName: types.NamespacedName{Name: node.Name}})
				}
				return requests
			}),
		).
		Watches(
			// Reconcile node when a pod assigned to it changes.
			&source.Kind{Type: &v1.Pod{}},
			handler.EnqueueRequestsFromMapFunc(func(o client.Object) (requests []reconcile.Request) {
				if name := o.(*v1.Pod).Spec.NodeName; name != "" {
					requests = append(requests, reconcile.Request{NamespacedName: types.NamespacedName{Name: name}})
				}
				return requests
			}),
		).
		WithOptions(controller.Options{MaxConcurrentReconciles: 10}).
		Complete(c)
}
