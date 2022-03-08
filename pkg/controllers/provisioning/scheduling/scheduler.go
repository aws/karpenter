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

	"github.com/prometheus/client_golang/prometheus"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crmetrics "sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/aws/karpenter/pkg/apis/provisioning/v1alpha5"
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
	KubeClient client.Client
	Topology   *Topology
}

type Schedule struct {
	*v1alpha5.Constraints
	// Pods is a set of pods that may schedule to the node; used for binpacking.
	Pods []*v1.Pod
}

func NewScheduler(kubeClient client.Client) *Scheduler {
	return &Scheduler{
		KubeClient: kubeClient,
		Topology:   &Topology{kubeClient: kubeClient},
	}
}

func (s *Scheduler) Solve(ctx context.Context, provisioner *v1alpha5.Provisioner, pods []*v1.Pod) ([]*Schedule, error) {
	defer metrics.Measure(schedulingDuration.WithLabelValues(injection.GetNamespacedName(ctx).Name))()
	constraints := provisioner.Spec.Constraints.DeepCopy()
	// Inject temporarily adds specific NodeSelectors to pods, which are then
	// used by scheduling logic. This isn't strictly necessary, but is a useful
	// trick to avoid passing topology decisions through the scheduling code. It
	// lets us to treat TopologySpreadConstraints as just-in-time NodeSelectors.
	if err := s.Topology.Inject(ctx, constraints, pods); err != nil {
		return nil, fmt.Errorf("injecting topology, %w", err)
	}
	// Separate pods into schedules of isomorphic scheduling constraints.
	schedules := s.getSchedules(constraints, pods)
	return schedules, nil
}

// getSchedules separates pods into a set of schedules. All pods in each group
// contain isomorphic scheduling constraints and can be deployed together on the
// same node, or multiple similar nodes if the pods exceed one node's capacity.
func (s *Scheduler) getSchedules(constraints *v1alpha5.Constraints, pods []*v1.Pod) []*Schedule {
	schedules := []*Schedule{}
	for _, pod := range pods {
		isCompatible := false
		for index, schedule := range schedules {
			if err := schedule.Requirements.Compatible(v1alpha5.NewPodRequirements(pod)); err == nil {
				//TODO: The cross product of the constraints values from cloud provider may not be supported by any instance type.
				// Need to test if there is any instance type that can support the combined constraints.
				schedules[index].Constraints = schedules[index].Tighten(pod)
				schedules[index].Pods = append(schedules[index].Pods, pod)
				isCompatible = true
				break
			}
		}
		if !isCompatible {
			schedules = append(schedules, &Schedule{Constraints: constraints.Tighten(pod), Pods: []*v1.Pod{pod}})
		}
	}
	return schedules
}
