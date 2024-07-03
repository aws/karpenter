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

package amifamily

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/karpenter/pkg/cloudprovider"
	"sigs.k8s.io/karpenter/pkg/scheduling"

	v1 "github.com/aws/karpenter-provider-aws/pkg/apis/v1"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/karpenter-provider-aws/pkg/providers/amifamily/bootstrap"
	"github.com/aws/karpenter-provider-aws/pkg/providers/ssm"
)

type AL2023 struct {
	DefaultFamily
	*Options
}

func (a AL2023) AMIQuery(ctx context.Context, ssmProvider ssm.Provider, k8sVersion string, amiVersion string) (AMIQuery, error) {
	query := AMIQuery{
		Filters: []*ec2.Filter{&ec2.Filter{
			Name: lo.ToPtr("image-id"),
		}},
		KnownRequirements: make(map[string][]scheduling.Requirements),
	}
	rootPath := fmt.Sprintf("/aws/service/eks/optimized-ami/%s/amazon-linux-2023", k8sVersion)
	results, err := ssmProvider.List(ctx, rootPath)
	if err != nil {
		return AMIQuery{}, fmt.Errorf("discovering AMIs from ssm")
	}
	for path, value := range results {
		pathComponents := strings.Split(path, "/")
		if len(pathComponents) != 11 || pathComponents[10] != "image_id" {
			continue
		}
		if av, err := a.extractAMIVersion(pathComponents[9]); err != nil || av != amiVersion {
			continue
		}
		variant, err := NewVariant(pathComponents[8])
		if err != nil {
			continue
		}
		query.Filters[0].Values = append(query.Filters[0].Values, lo.ToPtr(value))
		query.KnownRequirements[value] = []scheduling.Requirements{variant.Requirements()}
	}
	// Failed to discover any AMIs, we should short circuit AMI discovery
	if len(query.Filters[0].Values) == 0 {
		return AMIQuery{}, fmt.Errorf("failed to discover any AMIs for alias")
	}
	return query, nil
}

func (a AL2023) extractAMIVersion(versionStr string) (string, error) {
	if versionStr == "recommended" {
		return AMIVersionLatest, nil
	}
	rgx := regexp.MustCompile(`^.*(v\d+)$`)
	matches := rgx.FindStringSubmatch(versionStr)
	if len(matches) != 2 {
		return "", fmt.Errorf("failed to extract AMI version")
	}
	return matches[1], nil
}

func (a AL2023) UserData(kubeletConfig *v1.KubeletConfiguration, taints []corev1.Taint, labels map[string]string, caBundle *string, _ []*cloudprovider.InstanceType, customUserData *string, instanceStorePolicy *v1.InstanceStorePolicy) bootstrap.Bootstrapper {
	return bootstrap.Nodeadm{
		Options: bootstrap.Options{
			ClusterName:             a.Options.ClusterName,
			ClusterEndpoint:         a.Options.ClusterEndpoint,
			ClusterCIDR:             a.Options.ClusterCIDR,
			KubeletConfig:           kubeletConfig,
			Taints:                  taints,
			Labels:                  labels,
			CABundle:                caBundle,
			AWSENILimitedPodDensity: false,
			CustomUserData:          customUserData,
			InstanceStorePolicy:     instanceStorePolicy,
		},
	}
}

// DefaultBlockDeviceMappings returns the default block device mappings for the AMI Family
func (a AL2023) DefaultBlockDeviceMappings() []*v1.BlockDeviceMapping {
	return []*v1.BlockDeviceMapping{{
		DeviceName: a.EphemeralBlockDevice(),
		EBS:        &DefaultEBS,
	}}
}

func (a AL2023) EphemeralBlockDevice() *string {
	return lo.ToPtr("/dev/xvda")
}
