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

package packing

import (
	"github.com/awslabs/karpenter/pkg/utils/scheduling"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// TODO get this information from node-instance-selector
var (
	nodeCapacities = []*nodeCapacity{
		{
			instanceType: "m5.8xlarge",
			total: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("32000m"),
				v1.ResourceMemory: resource.MustParse("128Gi"),
			},
			reserved: v1.ResourceList{
				v1.ResourceCPU:    resource.Quantity{},
				v1.ResourceMemory: resource.Quantity{},
			},
		},
		{
			instanceType: "m5.2xlarge",
			total: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("8000m"),
				v1.ResourceMemory: resource.MustParse("32Gi"),
			},
			reserved: v1.ResourceList{
				v1.ResourceCPU:    resource.Quantity{},
				v1.ResourceMemory: resource.Quantity{},
			},
		},
		{
			instanceType: "m5.xlarge",
			total: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("4000m"),
				v1.ResourceMemory: resource.MustParse("16Gi"),
			},
			reserved: v1.ResourceList{
				v1.ResourceCPU:    resource.Quantity{},
				v1.ResourceMemory: resource.Quantity{},
			},
		},
		{
			instanceType: "m5.large",
			total: v1.ResourceList{
				v1.ResourceCPU:    resource.MustParse("2000m"),
				v1.ResourceMemory: resource.MustParse("8Gi"),
			},
			reserved: v1.ResourceList{
				v1.ResourceCPU:    resource.Quantity{},
				v1.ResourceMemory: resource.Quantity{},
			},
		},
	}
)

type nodeCapacity struct {
	instanceType string
	reserved     v1.ResourceList
	total        v1.ResourceList
}

func (nc *nodeCapacity) reserve(podSpec *v1.PodSpec) bool {
	resources := scheduling.GetResources(podSpec)
	cpu := nc.reserved.Cpu()
	cpu.Add(*resources.Cpu())
	targetMemory := nc.reserved.Memory()
	targetMemory.Add(*resources.Memory())
	targetPodCount := nc.reserved.Pods()
	targetPodCount.Add(*resource.NewQuantity(1, resource.BinarySI))
	// If pod fits reserve the capacity
	if nc.total.Cpu().Cmp(*targetCPU) >= 0 &&
		nc.total.Memory().Cmp(*targetMemory) >= 0 &&
		nc.total.Pods().Cmp(*targetPodCount) >= 0 {
		nc.reserved[v1.ResourceCPU] = *targetCPU
		nc.reserved[v1.ResourceMemory] = *targetMemory
		nc.reserved[v1.ResourcePods] = *targetPodCount
		return true
	}
	return false
}
