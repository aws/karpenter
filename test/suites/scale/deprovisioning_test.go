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

package scale_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	"github.com/samber/lo"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/aws/karpenter-core/pkg/apis/v1alpha5"
	"github.com/aws/karpenter-core/pkg/test"
	"github.com/aws/karpenter/pkg/apis/settings"
	"github.com/aws/karpenter/pkg/apis/v1alpha1"
	awstest "github.com/aws/karpenter/pkg/test"
	"github.com/aws/karpenter/test/pkg/debug"
)

var _ = Describe("Deprovisioning", Label(debug.NoWatch), Label(debug.NoEvents), func() {
	var provisioner *v1alpha5.Provisioner
	var nodeTemplate *v1alpha1.AWSNodeTemplate
	var deployment *appsv1.Deployment
	var selector labels.Selector
	var dsCount int

	BeforeEach(func() {
		// Expect the Prometheus client to be up
		env.ExpectSettingsOverridden(map[string]string{
			"featureGates.driftEnabled": "true",
		})
		env.ExpectPrometheusQuery("karpenter_machines_created", nil)
		nodeTemplate = awstest.AWSNodeTemplate(v1alpha1.AWSNodeTemplateSpec{AWS: v1alpha1.AWS{
			SecurityGroupSelector: map[string]string{"karpenter.sh/discovery": settings.FromContext(env.Context).ClusterName},
			SubnetSelector:        map[string]string{"karpenter.sh/discovery": settings.FromContext(env.Context).ClusterName},
		}})
		provisioner = test.Provisioner(test.ProvisionerOptions{
			ProviderRef: &v1alpha5.MachineTemplateRef{
				Name: nodeTemplate.Name,
			},
			Requirements: []v1.NodeSelectorRequirement{
				{
					Key:      v1alpha1.LabelInstanceSize,
					Operator: v1.NodeSelectorOpIn,
					Values:   []string{"4xlarge"},
				},
				{
					Key:      v1alpha5.LabelCapacityType,
					Operator: v1.NodeSelectorOpIn,
					Values:   []string{v1alpha1.CapacityTypeOnDemand},
				},
				{
					Key:      v1.LabelOSStable,
					Operator: v1.NodeSelectorOpIn,
					Values:   []string{string(v1.Linux)},
				},
			},
			// No limits!!!
			// https://tenor.com/view/chaos-gif-22919457
			Limits: v1.ResourceList{},
		})
		deployment = test.Deployment()
		// Zonal topology spread to avoid exhausting IPs in each subnet
		deployment.Spec.Template.Spec.TopologySpreadConstraints = []v1.TopologySpreadConstraint{
			{
				LabelSelector:     deployment.Spec.Selector,
				TopologyKey:       v1.LabelTopologyZone,
				MaxSkew:           1,
				WhenUnsatisfiable: v1.DoNotSchedule,
			},
		}
		selector = labels.SelectorFromSet(deployment.Spec.Selector.MatchLabels)
		dsCount = env.GetDaemonSetCount(provisioner)
	})

	AfterEach(func() {
		env.Cleanup()
	})

	Context("Multiple Deprovisioners", func() {})
	Context("Consolidation", func() {})
	Context("Emptiness", func() {
		It("should deprovision all nodes when empty", func() {
			// Before Deprovisioning, we need to Provision the cluster to the state that we need.
			replicasPerNode := 1
			maxPodDensity := replicasPerNode + dsCount
			expectedNodeCount := 200
			replicas := replicasPerNode * expectedNodeCount

			deployment.Spec.Replicas = lo.ToPtr[int32](int32(replicas))
			provisioner.Spec.KubeletConfiguration = &v1alpha5.KubeletConfiguration{
				MaxPods: lo.ToPtr[int32](int32(maxPodDensity)),
			}

			By("waiting for the deployment to deploy all of its pods")
			env.ExpectCreated(deployment)
			env.EventuallyExpectPendingPodCount(selector, replicas)

			By("kicking off provisioning by applying the provisioner and nodeTemplate")
			env.ExpectCreated(provisioner, nodeTemplate)

			startTime := time.Now()

			env.EventuallyExpectCreatedMachineCount(">=", expectedNodeCount)
			env.EventuallyExpectCreatedNodeCount(">=", expectedNodeCount)
			env.EventuallyExpectInitializedNodeCount(">=", expectedNodeCount)
			env.EventuallyExpectHealthyPodCount(selector, replicas)
			fmt.Printf("It took %f for all pods to be healthy", time.Since(startTime).Seconds())

			createdNodes := env.Monitor.CreatedNodeCount()

			By(fmt.Sprintf("created %d nodes. resetting monitor for deprovisioning.", createdNodes))
			env.Monitor.Reset()
			By("waiting for all deployment pods to be deleted")
			// Fully scale down all pods to make nodes empty
			deployment.Spec.Replicas = lo.ToPtr[int32](0)
			env.ExpectDeleted(deployment)
			env.EventuallyExpectHealthyPodCount(selector, 0)

			By("kicking off deprovisioning by adding ttlSecondsAfterEmpty")
			startTime = time.Now()
			provisioner.Spec.TTLSecondsAfterEmpty = lo.ToPtr[int64](0)
			env.ExpectCreatedOrUpdated(provisioner)

			env.EventuallyExpectDeletedNodeCount("==", createdNodes)
			fmt.Printf("It took %f for all nodes to be deprovisioned", time.Since(startTime).Seconds())
		})
	})
	Context("Expiration", func() {
		It("should expire all nodes", func() {
			// Before Deprovisioning, we need to Provision the cluster to the state that we need.
			replicasPerNode := 1
			maxPodDensity := replicasPerNode + dsCount
			expectedNodeCount := 30
			replicas := replicasPerNode * expectedNodeCount

			deployment.Spec.Replicas = lo.ToPtr[int32](int32(replicas))
			provisioner.Spec.KubeletConfiguration = &v1alpha5.KubeletConfiguration{
				MaxPods: lo.ToPtr[int32](int32(maxPodDensity)),
			}

			By("waiting for the deployment to deploy all of its pods")
			env.ExpectCreated(deployment)
			env.EventuallyExpectPendingPodCount(selector, replicas)

			By("kicking off provisioning by applying the provisioner and nodeTemplate")
			env.ExpectCreated(provisioner, nodeTemplate)

			env.EventuallyExpectCreatedMachineCount(">=", expectedNodeCount)
			env.EventuallyExpectCreatedNodeCount(">=", expectedNodeCount)
			env.EventuallyExpectInitializedNodeCount(">=", expectedNodeCount)
			env.EventuallyExpectHealthyPodCount(selector, replicas)

			createdNodes := env.Monitor.CreatedNodeCount()

			By(fmt.Sprintf("Created %d nodes. Resetting monitor for deprovisioning.", createdNodes))
			env.Monitor.Reset()
			By("kicking off deprovisioning by adding expiration and another provisioner")
			// Change Provisioner limits so that replacement nodes will use another provisioner.
			provisioner.Spec.Limits = &v1alpha5.Limits{
				Resources: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse("0"),
					v1.ResourceMemory: resource.MustParse("0Gi"),
				},
			}
			// Enable Expiration
			provisioner.Spec.TTLSecondsUntilExpired = lo.ToPtr[int64](0)

			noExpireProvisioner := test.Provisioner(test.ProvisionerOptions{
				Requirements: []v1.NodeSelectorRequirement{
					{
						Key:      v1alpha1.LabelInstanceSize,
						Operator: v1.NodeSelectorOpIn,
						Values:   []string{"4xlarge"},
					},
					{
						Key:      v1alpha5.LabelCapacityType,
						Operator: v1.NodeSelectorOpIn,
						Values:   []string{v1alpha1.CapacityTypeOnDemand},
					},
				},
				ProviderRef: &v1alpha5.MachineTemplateRef{
					Name: nodeTemplate.Name,
				},
			})
			env.ExpectCreatedOrUpdated(provisioner, noExpireProvisioner)
			env.EventuallyExpectDeletedNodeCount("==", createdNodes)
		})
	})
	Context("Drift", func() {
		It("should drift all nodes", func() {
			// Before Deprovisioning, we need to Provision the cluster to the state that we need.
			replicasPerNode := 1
			maxPodDensity := replicasPerNode + dsCount
			expectedNodeCount := 30
			replicas := replicasPerNode * expectedNodeCount

			deployment.Spec.Replicas = lo.ToPtr[int32](int32(replicas))
			provisioner.Spec.KubeletConfiguration = &v1alpha5.KubeletConfiguration{
				MaxPods: lo.ToPtr[int32](int32(maxPodDensity)),
			}

			By("waiting for the deployment to deploy all of its pods")
			env.ExpectCreated(deployment)
			env.EventuallyExpectPendingPodCount(selector, replicas)

			By("kicking off provisioning by applying the provisioner and nodeTemplate")
			env.ExpectCreated(provisioner, nodeTemplate)

			env.EventuallyExpectCreatedMachineCount(">=", expectedNodeCount)
			env.EventuallyExpectCreatedNodeCount(">=", expectedNodeCount)
			env.EventuallyExpectInitializedNodeCount(">=", expectedNodeCount)
			env.EventuallyExpectHealthyPodCount(selector, replicas)

			createdNodes := env.Monitor.CreatedNodeCount()

			By(fmt.Sprintf("Created %d nodes. Resetting monitor for deprovisioning.", createdNodes))
			env.Monitor.Reset()
			By("waiting for all deployment pods to be deleted")
			// Fully scale down all pods to make nodes empty
			deployment.Spec.Replicas = lo.ToPtr[int32](0)
			env.ExpectDeleted(deployment)
			env.EventuallyExpectHealthyPodCount(selector, 0)

			By("kicking off deprovisioning by adding ttlSecondsAfterEmpty")
			provisioner.Spec.TTLSecondsAfterEmpty = lo.ToPtr[int64](0)
			env.ExpectCreatedOrUpdated(provisioner)

			env.EventuallyExpectDeletedNodeCount("==", createdNodes)
		})
	})
	Context("Interruption", func() {})
})
