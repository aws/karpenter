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

package scheduling

import (
	"context"
	"fmt"
	"sort"

	"knative.dev/pkg/logging"

	"github.com/aws/karpenter/pkg/utils/resources"

	"github.com/prometheus/client_golang/prometheus"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crmetrics "sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/aws/karpenter/pkg/apis/provisioning/v1alpha5"
	"github.com/aws/karpenter/pkg/cloudprovider"
	"github.com/aws/karpenter/pkg/metrics"
	"github.com/aws/karpenter/pkg/utils/injection"
)

var schedulingDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: metrics.Namespace,
		Subsystem: "allocation_controller",
		Name:      "scheduling_duration_seconds",
		Help:      "Duration of scheduling process in seconds. Broken down by provisioner and error.",
		Buckets:   metrics.DurationBuckets(),
	},
	[]string{metrics.ProvisionerLabel},
)

func init() {
	crmetrics.Registry.MustRegister(schedulingDuration)
}

type Scheduler struct {
	kubeClient client.Client
}

func NewScheduler(kubeClient client.Client) *Scheduler {
	return &Scheduler{
		kubeClient: kubeClient,
	}
}

func (s *Scheduler) Solve(ctx context.Context, provisioner *v1alpha5.Provisioner, instanceTypes []cloudprovider.InstanceType, pods []*v1.Pod) (nodes []*Node, err error) {
	defer metrics.Measure(schedulingDuration.WithLabelValues(injection.GetNamespacedName(ctx).Name))()
	constraints := provisioner.Spec.Constraints.DeepCopy()

	sort.Slice(pods, byCPUAndMemoryDescending(pods))
	sort.Slice(instanceTypes, byPrice(instanceTypes))

	nodeSet, err := NewNodeSet(ctx, constraints, instanceTypes, s.kubeClient)
	if err != nil {
		return nil, fmt.Errorf("constructing nodeset, %w", err)
	}

	if err := nodeSet.TrackTopologies(ctx, pods); err != nil {
		return nil, fmt.Errorf("tracking topology counts, %w", err)
	}

	type unschedulablePod struct {
		pod    *v1.Pod
		reason error
	}
	var unschedulablePods []unschedulablePod
	for _, p := range pods {
		unschedulablePods = append(unschedulablePods, unschedulablePod{pod: p, reason: nil})
	}
	previousUnschedulableCount := 0
	// We loop and retrying to schedule to unschedulable pods as long as we are making progress.  This solves a few
	// issues including pods with affinity to another pod in the batch. We could topo-sort to solve this, but it wouldn't
	// solve the problem of scheduling pods where a paritcular order is needed to prevent a max-skew violation. E.g. if we
	// had 5xA pods and 5xB pods were they have a zonal topology spread, but A can only go in one zone and B in another.
	// We need to schedule them alternating, A, B, A, B, .... and this solution also solves that as well.
	for {
		previousUnschedulableCount = len(unschedulablePods)
		var newUnschedulablePods []unschedulablePod
		for _, up := range unschedulablePods {
			if err := nodeSet.Schedule(ctx, up.pod); err != nil {
				newUnschedulablePods = append(newUnschedulablePods, unschedulablePod{pod: up.pod, reason: err})
			}
		}
		unschedulablePods = newUnschedulablePods
		// if there are no more pods to attempt scheduling, or we tried each pod in a scheduling round
		// and made no progress, we are finished
		if len(unschedulablePods) == 0 || len(unschedulablePods) == previousUnschedulableCount {
			break
		}
	}

	if len(unschedulablePods) != 0 {
		for _, up := range unschedulablePods {
			logging.FromContext(ctx).With("pod", client.ObjectKeyFromObject(up.pod)).Errorf("Scheduling pod, %s", up.reason)
		}
		logging.FromContext(ctx).Errorf("Failed to schedule %d pod(s)", len(unschedulablePods))
	}

	return nodeSet.nodes, nil
}

func byPrice(instanceTypes []cloudprovider.InstanceType) func(i int, j int) bool {
	return func(i, j int) bool {
		return instanceTypes[i].Price() < instanceTypes[j].Price()
	}
}

func byCPUAndMemoryDescending(pods []*v1.Pod) func(i int, j int) bool {
	return func(i, j int) bool {
		lhs := resources.RequestsForPods(pods[i])
		rhs := resources.RequestsForPods(pods[j])

		cpuCmp := resources.Cmp(lhs[v1.ResourceCPU], rhs[v1.ResourceCPU])
		if cpuCmp < 0 {
			// LHS has less CPU, so it should be sorted after
			return false
		} else if cpuCmp > 0 {
			return true
		}
		memCmp := resources.Cmp(lhs[v1.ResourceMemory], rhs[v1.ResourceMemory])

		if memCmp < 0 {
			return false
		} else if memCmp > 0 {
			return true
		}
		return false
	}
}
