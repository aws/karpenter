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

	provisioning "github.com/awslabs/karpenter/pkg/apis/provisioning/v1alpha3"
	"github.com/awslabs/karpenter/pkg/cloudprovider"
	"github.com/awslabs/karpenter/pkg/controllers/allocation/scheduling"
	"github.com/awslabs/karpenter/pkg/utils/functional"
	"github.com/awslabs/karpenter/pkg/utils/injectabletime"
	"github.com/awslabs/karpenter/pkg/utils/ptr"

	"knative.dev/pkg/logging"

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
	for _, pod := range pods {
		if val := pod.Annotations[provisioning.DoNotEvictPodAnnotationKey]; val == "true" {
			logging.FromContext(ctx).Debugf("Unable to drain node %s, pod %s has do-not-evict annotation", node.Name, pod.Name)
			return false, nil
		}
	}

	// 4. Get and evict pods
	evictable := t.getEvictablePods(pods)
	if len(evictable) == 0 {
		return true, nil
	}
	t.evict(ctx, evictable)
	return false, nil
}

// terminate calls cloud provider delete then removes the finalizer to delete the node
func (t *Terminator) terminate(ctx context.Context, node *v1.Node) error {
	// 1. Delete the instance associated with node
	if err := t.CloudProvider.Delete(ctx, node); err != nil {
		return fmt.Errorf("terminating cloudprovider instance, %w", err)
	}
	// 2. Remove finalizer from node in APIServer
	persisted := node.DeepCopy()
	node.Finalizers = functional.StringSliceWithout(node.Finalizers, provisioning.TerminationFinalizer)
	if err := t.KubeClient.Patch(ctx, node, client.MergeFrom(persisted)); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
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

func (t *Terminator) getEvictablePods(pods []*v1.Pod) []*v1.Pod {
	evictable := []*v1.Pod{}
	for _, pod := range pods {
		// Ignore if unschedulable is tolerated, since they will reschedule
		if scheduling.Tolerates(pod, v1.Taint{Key: v1.TaintNodeUnschedulable, Effect: v1.TaintEffectNoSchedule}) == nil {
			continue
		}
		// Ignore if kubelet is partitioned and pods are beyond graceful termination window
		if IsStuckTerminating(pod) {
			continue
		}
		evictable = append(evictable, pod)
	}
	return evictable
}

func (t *Terminator) evict(ctx context.Context, pods []*v1.Pod) {
	// 1. Prioritize noncritical pods https://kubernetes.io/docs/concepts/architecture/nodes/#graceful-node-shutdown
	critical := []*v1.Pod{}
	nonCritical := []*v1.Pod{}
	for _, pod := range pods {
		if !pod.DeletionTimestamp.IsZero() {
			continue
		}
		if pod.Spec.PriorityClassName != "system-cluster-critical" && pod.Spec.PriorityClassName != "system-node-critical" {
			critical = append(critical, pod)
		} else {
			nonCritical = append(nonCritical, pod)
		}
	}
	// 2. Evict critical pods if all noncritical are evicted
	if len(nonCritical) == 0 {
		t.EvictionQueue.Add(critical)
	} else {
		t.EvictionQueue.Add(nonCritical)
	}
}

func IsStuckTerminating(pod *v1.Pod) bool {
	if pod.DeletionTimestamp == nil {
		return false
	}
	return injectabletime.Now().After(pod.DeletionTimestamp.Time)
}
