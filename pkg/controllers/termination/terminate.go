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

package termination

import (
	"context"
	"fmt"
	"time"

	provisioning "github.com/awslabs/karpenter/pkg/apis/provisioning/v1alpha3"
	"github.com/awslabs/karpenter/pkg/cloudprovider"
	"github.com/awslabs/karpenter/pkg/controllers/allocation/scheduling"
	"github.com/awslabs/karpenter/pkg/utils/functional"
	"github.com/awslabs/karpenter/pkg/utils/ptr"

	"knative.dev/pkg/logging"
	kptr "knative.dev/pkg/ptr"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Terminator struct {
	EvictionQueue *EvictionQueue
	KubeClient    client.Client
	CoreV1Client  corev1.CoreV1Interface
	CloudProvider cloudprovider.CloudProvider
}

// cordon cordons a node
func (t *Terminator) cordon(ctx context.Context, node *v1.Node) error {
	// 1. Check if node is already cordoned
	if node.Spec.Unschedulable {
		return nil
	}
	// 2. Cordon node
	persisted := node.DeepCopy()
	node.Spec.Unschedulable = true
	if err := t.KubeClient.Patch(ctx, node, client.MergeFrom(persisted)); err != nil {
		return fmt.Errorf("patching node %s, %w", node.Name, err)
	}
	logging.FromContext(ctx).Infof("Cordoned node %s", node.Name)
	return nil
}

// drain evicts pods from the node and returns true when all pods are evicted
func (t *Terminator) drain(ctx context.Context, node *v1.Node) (bool, error) {
	// 1. Get pods on node
	pods, err := t.getPods(ctx, node)
	if err != nil {
		return false, fmt.Errorf("listing pods for node %s, %w", node.Name, err)
	}

	// 2. Separate pods as non-critical and critical
	// https://kubernetes.io/docs/concepts/architecture/nodes/#graceful-node-shutdown
	nonCritical := []*v1.Pod{}
	critical := []*v1.Pod{}
	terminable := true

	for _, pod := range pods {
		if val := pod.Annotations[provisioning.DoNotEvictPodAnnotationKey]; val == "true" {
			logging.FromContext(ctx).Debugf("Unable to drain node %s, pod %s has do-not-evict annotation", node.Name, pod.Name)
			return false, nil
		}
		if scheduling.Tolerates(pod, v1.Taint{Key: v1.TaintNodeUnschedulable, Effect: v1.TaintEffectNoSchedule}) == nil {
			continue
		}
		// Don't attempt to evict a pod that's already evicting
		if !pod.DeletionTimestamp.IsZero() {
			// If the kubelet cannot register a pod as deleted, we need to force delete the pod in order to terminate the node
			if time.Since(pod.DeletionTimestamp.Time) > time.Duration(*pod.DeletionGracePeriodSeconds) {
				if err := t.KubeClient.Delete(ctx, pod, &client.DeleteOptions{GracePeriodSeconds: kptr.Int64(0)}); err != nil {
					return false, fmt.Errorf("force deleting pod %s for expired eviction", pod.Name)
				}
				continue
			}
			// If there are still pods that need to finish evicting, wait to terminate
			terminable = false
		}
		if pod.Spec.PriorityClassName == "system-cluster-critical" || pod.Spec.PriorityClassName == "system-node-critical" {
			critical = append(critical, pod)
		} else {
			nonCritical = append(nonCritical, pod)
		}
	}
	// 3. Evict non-critical pods
	if len(nonCritical) != 0 {
		t.EvictionQueue.Add(nonCritical)
		return false, nil
	}
	// 4. Evict critical pods once all non-critical pods are evicted
	if len(critical) != 0 {
		t.EvictionQueue.Add(critical)
		return false, nil
	}
	return terminable, nil
}

// terminate terminates the node then removes the finalizer to delete the node
func (t *Terminator) terminate(ctx context.Context, node *v1.Node) error {
	// 1. Terminate instance associated with node
	if err := t.CloudProvider.Terminate(ctx, node); err != nil {
		return fmt.Errorf("terminating cloudprovider instance, %w", err)
	}
	// 2. Remove finalizer from node in APIServer
	persisted := node.DeepCopy()
	node.Finalizers = functional.StringSliceWithout(node.Finalizers, provisioning.TerminationFinalizer)
	if err := t.KubeClient.Patch(ctx, node, client.MergeFrom(persisted)); err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("removing finalizer from node %s, %w", node.Name, err)
	}
	logging.FromContext(ctx).Infof("Deleted node %s", node.Name)
	return nil
}

// getPods returns a list of pods scheduled to a node based on some filters
func (t *Terminator) getPods(ctx context.Context, node *v1.Node) ([]*v1.Pod, error) {
	pods := &v1.PodList{}
	if err := t.KubeClient.List(ctx, pods, client.MatchingFields{"spec.nodeName": node.Name}); err != nil {
		return nil, fmt.Errorf("listing pods on node %s, %w", node.Name, err)
	}
	return ptr.PodListToSlice(pods), nil
}
