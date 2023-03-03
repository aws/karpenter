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

package cloudprovider_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/samber/lo"
	"k8s.io/client-go/tools/record"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aws/aws-sdk-go/aws"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clock "k8s.io/utils/clock/testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "knative.dev/pkg/logging/testing"

	"github.com/aws/karpenter/pkg/apis"
	"github.com/aws/karpenter/pkg/apis/settings"
	"github.com/aws/karpenter/pkg/apis/v1alpha1"
	"github.com/aws/karpenter/pkg/cloudprovider"
	awscontext "github.com/aws/karpenter/pkg/context"
	"github.com/aws/karpenter/pkg/fake"
	"github.com/aws/karpenter/pkg/test"

	coresettings "github.com/aws/karpenter-core/pkg/apis/settings"
	"github.com/aws/karpenter-core/pkg/apis/v1alpha5"
	corecloudprovider "github.com/aws/karpenter-core/pkg/cloudprovider"
	"github.com/aws/karpenter-core/pkg/controllers/provisioning"
	"github.com/aws/karpenter-core/pkg/controllers/state"
	"github.com/aws/karpenter-core/pkg/events"
	"github.com/aws/karpenter-core/pkg/operator/controller"
	"github.com/aws/karpenter-core/pkg/operator/injection"
	"github.com/aws/karpenter-core/pkg/operator/options"
	"github.com/aws/karpenter-core/pkg/operator/scheme"
	coretest "github.com/aws/karpenter-core/pkg/test"
	. "github.com/aws/karpenter-core/pkg/test/expectations"
	machineutil "github.com/aws/karpenter-core/pkg/utils/machine"
)

var ctx context.Context
var stop context.CancelFunc
var opts options.Options
var env *coretest.Environment
var awsCtx awscontext.Context
var fakeEC2API *fake.EC2API
var fakeSSMAPI *fake.SSMAPI
var prov *provisioning.Provisioner
var provisioningController controller.Controller
var cloudProvider *cloudprovider.CloudProvider
var cluster *state.Cluster
var fakeClock *clock.FakeClock
var provisioner *v1alpha5.Provisioner
var nodeTemplate *v1alpha1.AWSNodeTemplate

func TestAWS(t *testing.T) {
	ctx = TestContextWithLogger(t)
	RegisterFailHandler(Fail)
	RunSpecs(t, "CloudProvider/AWS")
}

var _ = BeforeSuite(func() {
	env = coretest.NewEnvironment(scheme.Scheme, coretest.WithCRDs(apis.CRDs...))
	ctx = coresettings.ToContext(ctx, coretest.Settings())
	ctx = settings.ToContext(ctx, test.Settings())
	ctx, stop = context.WithCancel(ctx)

	fakeEC2API = &fake.EC2API{}
	fakeSSMAPI = &fake.SSMAPI{}
	fakeClock = clock.NewFakeClock(time.Now())
	awsCtx = test.Context(ctx, fakeEC2API, fakeSSMAPI, env, fakeClock, test.ContextOptions{})

	cloudProvider = cloudprovider.New(awsCtx, awsCtx.InstanceTypesProvider, awsCtx.InstanceProvider, awsCtx.KubeClient, awsCtx.AMIProvider)
	cluster = state.NewCluster(fakeClock, awsCtx.KubeClient, cloudProvider)
	prov = provisioning.NewProvisioner(ctx, awsCtx.KubeClient, env.KubernetesInterface.CoreV1(), events.NewRecorder(&record.FakeRecorder{}), cloudProvider, cluster)
	provisioningController = provisioning.NewController(awsCtx.KubeClient, prov, events.NewRecorder(&record.FakeRecorder{}))
})

var _ = AfterSuite(func() {
	stop()
	Expect(env.Stop()).To(Succeed(), "Failed to stop environment")
})

var _ = BeforeEach(func() {
	ctx = injection.WithOptions(ctx, opts)
	ctx = coresettings.ToContext(ctx, coretest.Settings())
	ctx = settings.ToContext(ctx, test.Settings())
	nodeTemplate = &v1alpha1.AWSNodeTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Name: coretest.RandomName(),
		},
		Spec: v1alpha1.AWSNodeTemplateSpec{
			AWS: v1alpha1.AWS{
				AMIFamily:             aws.String(v1alpha1.AMIFamilyAL2),
				SubnetSelector:        map[string]string{"*": "*"},
				SecurityGroupSelector: map[string]string{"*": "*"},
			},
		},
	}
	nodeTemplate.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   v1alpha1.SchemeGroupVersion.Group,
		Version: v1alpha1.SchemeGroupVersion.Version,
		Kind:    "AWSNodeTemplate",
	})
	provisioner = test.Provisioner(coretest.ProvisionerOptions{
		Requirements: []v1.NodeSelectorRequirement{{
			Key:      v1alpha1.LabelInstanceCategory,
			Operator: v1.NodeSelectorOpExists,
		}},
		ProviderRef: &v1alpha5.ProviderRef{
			APIVersion: nodeTemplate.APIVersion,
			Kind:       nodeTemplate.Kind,
			Name:       nodeTemplate.Name,
		},
	})

	cluster.Reset()
	fakeEC2API.Reset()
	fakeSSMAPI.Reset()
	awsCtx.RestProviderCache()
	awsCtx.LaunchTemplateProvider.KubeDNSIP = net.ParseIP("10.0.100.10")
	awsCtx.LaunchTemplateProvider.ClusterEndpoint = "https://test-cluster"
})

var _ = AfterEach(func() {
	ExpectCleanedUp(ctx, awsCtx.KubeClient)
})

var _ = Describe("CloudProvider", func() {
	Context("Defaulting", func() {
		// Intent here is that if updates occur on the provisioningController, the Provisioner doesn't need to be recreated
		It("should not set the InstanceProfile with the default if none provided in Provisioner", func() {
			provisioner.SetDefaults(ctx)
			constraints, err := v1alpha1.DeserializeProvider(provisioner.Spec.Provider.Raw)
			Expect(err).ToNot(HaveOccurred())
			Expect(constraints.InstanceProfile).To(BeNil())
		})
		It("should default requirements", func() {
			provisioner.SetDefaults(ctx)
			Expect(provisioner.Spec.Requirements).To(ContainElement(v1.NodeSelectorRequirement{
				Key:      v1alpha5.LabelCapacityType,
				Operator: v1.NodeSelectorOpIn,
				Values:   []string{v1alpha5.CapacityTypeOnDemand},
			}))
			Expect(provisioner.Spec.Requirements).To(ContainElement(v1.NodeSelectorRequirement{
				Key:      v1.LabelArchStable,
				Operator: v1.NodeSelectorOpIn,
				Values:   []string{v1alpha5.ArchitectureAmd64},
			}))
		})
	})
	Context("EC2 Context", func() {
		It("should set context on the CreateFleet request if specified on the Provisioner", func() {
			provider, err := v1alpha1.DeserializeProvider(provisioner.Spec.Provider.Raw)
			Expect(err).ToNot(HaveOccurred())
			provider.Context = aws.String("context-1234")
			provider.SubnetSelector = map[string]string{"*": "*"}
			provider.SecurityGroupSelector = map[string]string{"*": "*"}
			provisioner = coretest.Provisioner(coretest.ProvisionerOptions{Provider: provider})
			provisioner.SetDefaults(ctx)
			ExpectApplied(ctx, awsCtx.KubeClient, provisioner)
			pod := coretest.UnschedulablePod()
			ExpectProvisioned(ctx, awsCtx.KubeClient, cluster, prov, pod)
			ExpectScheduled(ctx, awsCtx.KubeClient, pod)
			Expect(fakeEC2API.CreateFleetBehavior.CalledWithInput.Len()).To(Equal(1))
			createFleetInput := fakeEC2API.CreateFleetBehavior.CalledWithInput.Pop()
			Expect(aws.StringValue(createFleetInput.Context)).To(Equal("context-1234"))
		})
		It("should default to no EC2 Context", func() {
			provisioner.SetDefaults(ctx)
			ExpectApplied(ctx, awsCtx.KubeClient, provisioner, nodeTemplate)
			pod := coretest.UnschedulablePod()
			ExpectProvisioned(ctx, awsCtx.KubeClient, cluster, prov, pod)
			ExpectScheduled(ctx, awsCtx.KubeClient, pod)
			Expect(fakeEC2API.CreateFleetBehavior.CalledWithInput.Len()).To(Equal(1))
			createFleetInput := fakeEC2API.CreateFleetBehavior.CalledWithInput.Pop()
			Expect(createFleetInput.Context).To(BeNil())
		})
	})
	Context("Node Drift", func() {
		var validAMI string
		var selectedInstanceType *corecloudprovider.InstanceType
		var instance *ec2.Instance
		BeforeEach(func() {
			validAMI = fake.ImageID()
			fakeSSMAPI.GetParameterOutput = &ssm.GetParameterOutput{
				Parameter: &ssm.Parameter{Value: aws.String(validAMI)},
			}
			fakeEC2API.DescribeImagesOutput.Set(&ec2.DescribeImagesOutput{
				Images: []*ec2.Image{{ImageId: aws.String(validAMI)}},
			})
			ExpectApplied(ctx, awsCtx.KubeClient, provisioner, nodeTemplate)
			instanceTypes, err := cloudProvider.GetInstanceTypes(ctx, provisioner)
			Expect(err).ToNot(HaveOccurred())
			selectedInstanceType = instanceTypes[0]

			// Create the instance we want returned from the EC2 API
			instance = &ec2.Instance{
				ImageId:               aws.String(validAMI),
				PrivateDnsName:        aws.String(randomdata.IpV4Address()),
				InstanceType:          aws.String(selectedInstanceType.Name),
				SpotInstanceRequestId: aws.String(coretest.RandomName()),
				State: &ec2.InstanceState{
					Name: aws.String(ec2.InstanceStateNameRunning),
				},
				InstanceId: aws.String(fake.InstanceID()),
			}
			fakeEC2API.DescribeInstancesBehavior.Output.Set(&ec2.DescribeInstancesOutput{
				Reservations: []*ec2.Reservation{{Instances: []*ec2.Instance{instance}}},
			})
		})
		It("should not fail if node template does not exist", func() {
			ExpectDeleted(ctx, awsCtx.KubeClient, nodeTemplate)
			node := coretest.Node(coretest.NodeOptions{
				ProviderID: fake.ProviderID(lo.FromPtr(instance.InstanceId)),
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						v1alpha5.ProvisionerNameLabelKey: provisioner.Name,
						v1.LabelInstanceTypeStable:       selectedInstanceType.Name,
					},
				},
			})
			drifted, err := cloudProvider.IsMachineDrifted(ctx, machineutil.NewFromNode(node))
			Expect(err).ToNot(HaveOccurred())
			Expect(drifted).To(BeFalse())
		})
		It("should return false if providerRef is not defined", func() {
			provisioner.Spec.ProviderRef = nil
			ExpectApplied(ctx, awsCtx.KubeClient, provisioner)
			node := coretest.Node(coretest.NodeOptions{
				ProviderID: fake.ProviderID(lo.FromPtr(instance.InstanceId)),
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						v1alpha5.ProvisionerNameLabelKey: provisioner.Name,
						v1.LabelInstanceTypeStable:       selectedInstanceType.Name,
					},
				},
			})
			drifted, err := cloudProvider.IsMachineDrifted(ctx, machineutil.NewFromNode(node))
			Expect(err).ToNot(HaveOccurred())
			Expect(drifted).To(BeFalse())
		})
		It("should not fail if provisioner does not exist", func() {
			ExpectDeleted(ctx, awsCtx.KubeClient, provisioner)
			node := coretest.Node(coretest.NodeOptions{
				ProviderID: fake.ProviderID(lo.FromPtr(instance.InstanceId)),
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						v1alpha5.ProvisionerNameLabelKey: provisioner.Name,
						v1.LabelInstanceTypeStable:       selectedInstanceType.Name,
					},
				},
			})
			drifted, err := cloudProvider.IsMachineDrifted(ctx, machineutil.NewFromNode(node))
			Expect(err).ToNot(HaveOccurred())
			Expect(drifted).To(BeFalse())
		})
		It("should return drifted if the AMI is not valid", func() {
			node := coretest.Node(coretest.NodeOptions{
				ProviderID: fake.ProviderID(lo.FromPtr(instance.InstanceId)),
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						v1alpha5.ProvisionerNameLabelKey: provisioner.Name,
						v1.LabelInstanceTypeStable:       selectedInstanceType.Name,
					},
				},
			})
			// Instance is a reference to what we return in the GetInstances call
			instance.ImageId = aws.String(fake.ImageID())
			isDrifted, err := cloudProvider.IsMachineDrifted(ctx, machineutil.NewFromNode(node))
			Expect(err).ToNot(HaveOccurred())
			Expect(isDrifted).To(BeTrue())
		})
		It("should not return drifted if the AMI is valid", func() {
			node := coretest.Node(coretest.NodeOptions{
				ProviderID: fake.ProviderID(lo.FromPtr(instance.InstanceId)),
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						v1alpha5.ProvisionerNameLabelKey: provisioner.Name,
						v1.LabelInstanceTypeStable:       selectedInstanceType.Name,
					},
				},
			})
			isDrifted, err := cloudProvider.IsMachineDrifted(ctx, machineutil.NewFromNode(node))
			Expect(err).ToNot(HaveOccurred())
			Expect(isDrifted).To(BeFalse())
		})
		It("should error if the node doesn't have the instance-type label", func() {
			node := coretest.Node(coretest.NodeOptions{
				ProviderID: fake.ProviderID(lo.FromPtr(instance.InstanceId)),
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						v1alpha5.ProvisionerNameLabelKey: provisioner.Name,
					},
				},
			})
			_, err := cloudProvider.IsMachineDrifted(ctx, machineutil.NewFromNode(node))
			Expect(err).To(HaveOccurred())
		})
		It("should error drift if node doesn't have provider id", func() {
			node := coretest.Node(coretest.NodeOptions{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						v1alpha5.ProvisionerNameLabelKey: provisioner.Name,
						v1.LabelInstanceTypeStable:       selectedInstanceType.Name,
					},
				},
			})
			isDrifted, err := cloudProvider.IsMachineDrifted(ctx, machineutil.NewFromNode(node))
			Expect(err).To(HaveOccurred())
			Expect(isDrifted).To(BeFalse())
		})
	})
	Context("Provider Backwards Compatibility", func() {
		It("should launch a node using provider defaults", func() {
			provisioner = test.Provisioner(coretest.ProvisionerOptions{
				Provider: v1alpha1.AWS{
					AMIFamily:             aws.String(v1alpha1.AMIFamilyAL2),
					SubnetSelector:        map[string]string{"*": "*"},
					SecurityGroupSelector: map[string]string{"*": "*"},
				},
				Requirements: []v1.NodeSelectorRequirement{{
					Key:      v1alpha1.LabelInstanceCategory,
					Operator: v1.NodeSelectorOpExists,
				}},
			})
			ExpectApplied(ctx, awsCtx.KubeClient, provisioner)
			pod := coretest.UnschedulablePod()
			ExpectProvisioned(ctx, awsCtx.KubeClient, cluster, prov, pod)
			ExpectScheduled(ctx, awsCtx.KubeClient, pod)

			Expect(fakeEC2API.CalledWithCreateLaunchTemplateInput.Len()).To(Equal(1))
			firstLt := fakeEC2API.CalledWithCreateLaunchTemplateInput.Pop()
			Expect(fakeEC2API.CreateFleetBehavior.CalledWithInput.Len()).To(Equal(1))

			createFleetInput := fakeEC2API.CreateFleetBehavior.CalledWithInput.Pop()
			launchTemplate := createFleetInput.LaunchTemplateConfigs[0].LaunchTemplateSpecification
			Expect(createFleetInput.LaunchTemplateConfigs).To(HaveLen(1))

			Expect(*createFleetInput.LaunchTemplateConfigs[0].LaunchTemplateSpecification.LaunchTemplateName).
				To(Equal(*firstLt.LaunchTemplateName))
			Expect(firstLt.LaunchTemplateData.BlockDeviceMappings[0].Ebs.Encrypted).To(Equal(aws.Bool(true)))
			Expect(*launchTemplate.Version).To(Equal("$Latest"))
		})
		It("should discover security groups by ID", func() {
			provisioner = test.Provisioner(coretest.ProvisionerOptions{
				Provider: v1alpha1.AWS{
					AMIFamily:             aws.String(v1alpha1.AMIFamilyAL2),
					SubnetSelector:        map[string]string{"*": "*"},
					SecurityGroupSelector: map[string]string{"aws-ids": "sg-test1"},
				},
				Requirements: []v1.NodeSelectorRequirement{{
					Key:      v1alpha1.LabelInstanceCategory,
					Operator: v1.NodeSelectorOpExists,
				}},
			})
			ExpectApplied(ctx, awsCtx.KubeClient, provisioner)
			pod := coretest.UnschedulablePod()
			ExpectProvisioned(ctx, awsCtx.KubeClient, cluster, prov, pod)
			ExpectScheduled(ctx, awsCtx.KubeClient, pod)
			Expect(fakeEC2API.CalledWithCreateLaunchTemplateInput.Len()).To(Equal(1))
			input := fakeEC2API.CalledWithCreateLaunchTemplateInput.Pop()
			Expect(aws.StringValueSlice(input.LaunchTemplateData.SecurityGroupIds)).To(ConsistOf(
				"sg-test1",
			))
		})
		It("should discover subnets by ID", func() {
			provisioner = test.Provisioner(coretest.ProvisionerOptions{
				Provider: v1alpha1.AWS{
					AMIFamily:             aws.String(v1alpha1.AMIFamilyAL2),
					SubnetSelector:        map[string]string{"aws-ids": "subnet-test1"},
					SecurityGroupSelector: map[string]string{"*": "*"},
				},
				Requirements: []v1.NodeSelectorRequirement{{
					Key:      v1alpha1.LabelInstanceCategory,
					Operator: v1.NodeSelectorOpExists,
				}},
			})
			ExpectApplied(ctx, awsCtx.KubeClient, provisioner)
			pod := coretest.UnschedulablePod()
			ExpectProvisioned(ctx, awsCtx.KubeClient, cluster, prov, pod)
			ExpectScheduled(ctx, awsCtx.KubeClient, pod)
			createFleetInput := fakeEC2API.CreateFleetBehavior.CalledWithInput.Pop()
			Expect(fake.SubnetsFromFleetRequest(createFleetInput)).To(ConsistOf("subnet-test1"))
		})
		It("should use the instance profile on the Provisioner when specified", func() {
			provisioner = test.Provisioner(coretest.ProvisionerOptions{
				Provider: v1alpha1.AWS{
					AMIFamily:             aws.String(v1alpha1.AMIFamilyAL2),
					SubnetSelector:        map[string]string{"*": "*"},
					SecurityGroupSelector: map[string]string{"*": "*"},
					InstanceProfile:       aws.String("overridden-profile"),
				},
				Requirements: []v1.NodeSelectorRequirement{{
					Key:      v1alpha1.LabelInstanceCategory,
					Operator: v1.NodeSelectorOpExists,
				}},
			})
			ExpectApplied(ctx, awsCtx.KubeClient, provisioner)
			pod := coretest.UnschedulablePod()
			ExpectProvisioned(ctx, awsCtx.KubeClient, cluster, prov, pod)
			ExpectScheduled(ctx, awsCtx.KubeClient, pod)
			Expect(fakeEC2API.CalledWithCreateLaunchTemplateInput.Len()).To(Equal(1))
			input := fakeEC2API.CalledWithCreateLaunchTemplateInput.Pop()
			Expect(*input.LaunchTemplateData.IamInstanceProfile.Name).To(Equal("overridden-profile"))
		})
	})
	Context("Subnet Compatibility", func() {
		// Note when debugging these tests -
		// hard coded fixture data (ex. what the aws api will return) is maintained in fake/ec2api.go
		It("should default to the cluster's subnets", func() {
			ExpectApplied(ctx, awsCtx.KubeClient, provisioner, nodeTemplate)
			pod := coretest.UnschedulablePod(
				coretest.PodOptions{NodeSelector: map[string]string{v1.LabelArchStable: v1alpha5.ArchitectureAmd64}})
			ExpectProvisioned(ctx, awsCtx.KubeClient, cluster, prov, pod)
			ExpectScheduled(ctx, awsCtx.KubeClient, pod)
			Expect(fakeEC2API.CreateFleetBehavior.CalledWithInput.Len()).To(Equal(1))
			input := fakeEC2API.CreateFleetBehavior.CalledWithInput.Pop()
			Expect(input.LaunchTemplateConfigs).To(HaveLen(1))

			foundNonGPULT := false
			for _, v := range input.LaunchTemplateConfigs {
				for _, ov := range v.Overrides {
					if *ov.InstanceType == "m5.large" {
						foundNonGPULT = true
						Expect(v.Overrides).To(ContainElements(
							&ec2.FleetLaunchTemplateOverridesRequest{SubnetId: aws.String("subnet-test1"), InstanceType: aws.String("m5.large"), AvailabilityZone: aws.String("test-zone-1a")},
							&ec2.FleetLaunchTemplateOverridesRequest{SubnetId: aws.String("subnet-test2"), InstanceType: aws.String("m5.large"), AvailabilityZone: aws.String("test-zone-1b")},
							&ec2.FleetLaunchTemplateOverridesRequest{SubnetId: aws.String("subnet-test3"), InstanceType: aws.String("m5.large"), AvailabilityZone: aws.String("test-zone-1c")},
						))
					}
				}
			}
			Expect(foundNonGPULT).To(BeTrue())
		})
		It("should launch instances into subnet with the most available IP addresses", func() {
			fakeEC2API.DescribeSubnetsOutput.Set(&ec2.DescribeSubnetsOutput{Subnets: []*ec2.Subnet{
				{SubnetId: aws.String("test-subnet-1"), AvailabilityZone: aws.String("test-zone-1a"), AvailableIpAddressCount: aws.Int64(10),
					Tags: []*ec2.Tag{{Key: aws.String("Name"), Value: aws.String("test-subnet-1")}}},
				{SubnetId: aws.String("test-subnet-2"), AvailabilityZone: aws.String("test-zone-1a"), AvailableIpAddressCount: aws.Int64(100),
					Tags: []*ec2.Tag{{Key: aws.String("Name"), Value: aws.String("test-subnet-2")}}},
			}})
			ExpectApplied(ctx, awsCtx.KubeClient, provisioner, nodeTemplate)
			pod := coretest.UnschedulablePod(coretest.PodOptions{NodeSelector: map[string]string{v1.LabelTopologyZone: "test-zone-1a"}})
			ExpectProvisioned(ctx, awsCtx.KubeClient, cluster, prov, pod)
			ExpectScheduled(ctx, awsCtx.KubeClient, pod)
			createFleetInput := fakeEC2API.CreateFleetBehavior.CalledWithInput.Pop()
			Expect(fake.SubnetsFromFleetRequest(createFleetInput)).To(ConsistOf("test-subnet-2"))
		})
		It("should launch instances into subnet with the most available IP addresses in-between cache refreshes", func() {
			fakeEC2API.DescribeSubnetsOutput.Set(&ec2.DescribeSubnetsOutput{Subnets: []*ec2.Subnet{
				{SubnetId: aws.String("test-subnet-1"), AvailabilityZone: aws.String("test-zone-1a"), AvailableIpAddressCount: aws.Int64(10),
					Tags: []*ec2.Tag{{Key: aws.String("Name"), Value: aws.String("test-subnet-1")}}},
				{SubnetId: aws.String("test-subnet-2"), AvailabilityZone: aws.String("test-zone-1a"), AvailableIpAddressCount: aws.Int64(11),
					Tags: []*ec2.Tag{{Key: aws.String("Name"), Value: aws.String("test-subnet-2")}}},
			}})
			provisioner.Spec.KubeletConfiguration = &v1alpha5.KubeletConfiguration{MaxPods: aws.Int32(1)}
			ExpectApplied(ctx, awsCtx.KubeClient, provisioner, nodeTemplate)
			pod1 := coretest.UnschedulablePod(coretest.PodOptions{NodeSelector: map[string]string{v1.LabelTopologyZone: "test-zone-1a"}})
			pod2 := coretest.UnschedulablePod(coretest.PodOptions{NodeSelector: map[string]string{v1.LabelTopologyZone: "test-zone-1a"}})
			ExpectProvisioned(ctx, awsCtx.KubeClient, cluster, prov, pod1, pod2)
			ExpectScheduled(ctx, awsCtx.KubeClient, pod1)
			ExpectScheduled(ctx, awsCtx.KubeClient, pod2)
			createFleetInput := fakeEC2API.CreateFleetBehavior.CalledWithInput.Pop()
			Expect(fake.SubnetsFromFleetRequest(createFleetInput)).To(ConsistOf("test-subnet-2"))
			// Provision for another pod that should now use the other subnet since we've consumed some from the first launch.
			pod3 := coretest.UnschedulablePod(coretest.PodOptions{NodeSelector: map[string]string{v1.LabelTopologyZone: "test-zone-1a"}})
			ExpectProvisioned(ctx, awsCtx.KubeClient, cluster, prov, pod3)
			ExpectScheduled(ctx, awsCtx.KubeClient, pod3)
			createFleetInput = fakeEC2API.CreateFleetBehavior.CalledWithInput.Pop()
			Expect(fake.SubnetsFromFleetRequest(createFleetInput)).To(ConsistOf("test-subnet-1"))
		})
		It("should update in-flight IPs when a CreateFleet error occurs", func() {
			fakeEC2API.DescribeSubnetsOutput.Set(&ec2.DescribeSubnetsOutput{Subnets: []*ec2.Subnet{
				{SubnetId: aws.String("test-subnet-1"), AvailabilityZone: aws.String("test-zone-1a"), AvailableIpAddressCount: aws.Int64(10),
					Tags: []*ec2.Tag{{Key: aws.String("Name"), Value: aws.String("test-subnet-1")}}},
			}})
			pod1 := coretest.UnschedulablePod(coretest.PodOptions{NodeSelector: map[string]string{v1.LabelTopologyZone: "test-zone-1a"}})
			ExpectApplied(ctx, awsCtx.KubeClient, provisioner, nodeTemplate, pod1)
			fakeEC2API.CreateFleetBehavior.Error.Set(fmt.Errorf("CreateFleet synthetic error"))
			bindings := ExpectProvisioned(ctx, awsCtx.KubeClient, cluster, prov, pod1)
			Expect(len(bindings)).To(Equal(0))
		})
		It("should launch instances into subnets that are excluded by another provisioner", func() {
			fakeEC2API.DescribeSubnetsOutput.Set(&ec2.DescribeSubnetsOutput{Subnets: []*ec2.Subnet{
				{SubnetId: aws.String("test-subnet-1"), AvailabilityZone: aws.String("test-zone-1a"), AvailableIpAddressCount: aws.Int64(10),
					Tags: []*ec2.Tag{{Key: aws.String("Name"), Value: aws.String("test-subnet-1")}}},
				{SubnetId: aws.String("test-subnet-2"), AvailabilityZone: aws.String("test-zone-1b"), AvailableIpAddressCount: aws.Int64(100),
					Tags: []*ec2.Tag{{Key: aws.String("Name"), Value: aws.String("test-subnet-2")}}},
			}})
			nodeTemplate.Spec.SubnetSelector = map[string]string{"Name": "test-subnet-1"}
			ExpectApplied(ctx, awsCtx.KubeClient, provisioner, nodeTemplate)
			podSubnet1 := coretest.UnschedulablePod()
			ExpectProvisioned(ctx, awsCtx.KubeClient, cluster, prov, podSubnet1)
			ExpectScheduled(ctx, awsCtx.KubeClient, podSubnet1)
			createFleetInput := fakeEC2API.CreateFleetBehavior.CalledWithInput.Pop()
			Expect(fake.SubnetsFromFleetRequest(createFleetInput)).To(ConsistOf("test-subnet-1"))

			provisioner = test.Provisioner(coretest.ProvisionerOptions{Provider: &v1alpha1.AWS{
				SubnetSelector:        map[string]string{"Name": "test-subnet-2"},
				SecurityGroupSelector: map[string]string{"*": "*"},
			}})
			ExpectApplied(ctx, awsCtx.KubeClient, provisioner)
			podSubnet2 := coretest.UnschedulablePod(coretest.PodOptions{NodeSelector: map[string]string{v1alpha5.ProvisionerNameLabelKey: provisioner.Name}})
			ExpectProvisioned(ctx, awsCtx.KubeClient, cluster, prov, podSubnet2)
			ExpectScheduled(ctx, awsCtx.KubeClient, podSubnet2)
			createFleetInput = fakeEC2API.CreateFleetBehavior.CalledWithInput.Pop()
			Expect(fake.SubnetsFromFleetRequest(createFleetInput)).To(ConsistOf("test-subnet-2"))
		})
	})
})
