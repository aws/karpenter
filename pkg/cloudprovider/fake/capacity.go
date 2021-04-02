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

package fake

import (
	"context"
	"fmt"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/awslabs/karpenter/pkg/cloudprovider"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Capacity struct {
}

func (c *Capacity) Create(ctx context.Context, constraints *cloudprovider.Constraints) ([]cloudprovider.Packing, error) {
	name := strings.ToLower(randomdata.SillyName())
	return []cloudprovider.Packing{{
		Node: &v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name:   name,
				Labels: constraints.Labels,
			},
			Spec: v1.NodeSpec{
				ProviderID: fmt.Sprintf("fake:///%s", name),
				Taints:     constraints.Taints,
			},
			Status: v1.NodeStatus{
				Allocatable: v1.ResourceList{
					// The (fake) mythical mega-machine
					v1.ResourcePods:   resource.MustParse("1000000"),
					v1.ResourceCPU:    resource.MustParse("1000000"),
					v1.ResourceMemory: resource.MustParse("1000000Gi"),
				},
			},
		},
		Pods: constraints.Pods,
	}}, nil
}

func (c *Capacity) Delete(ctx context.Context, nodes []*v1.Node) (map[string]bool, error) {
	return nil, nil
}

func (c *Capacity) GetZones(ctx context.Context) ([]string, error) {
	return []string{
		"test-zone-1",
		"test-zone-2",
	}, nil
}

func (c *Capacity) GetInstanceTypes(ctx context.Context) ([]string, error) {
	return []string{
		"test-instance-type-1",
		"test-instance-type-2",
	}, nil
}

func (c *Capacity) GetArchitectures(ctx context.Context) ([]string, error) {
	return []string{
		"test-architecture-1",
		"test-architecture-2",
	}, nil
}

func (c *Capacity) GetOperatingSystems(ctx context.Context) ([]string, error) {
	return []string{
		"test-operating-system-1",
		"test-operating-system-2",
	}, nil
}
