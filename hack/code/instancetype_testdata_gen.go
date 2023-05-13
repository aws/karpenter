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

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/samber/lo"
)

const packageHeader = `
package fake

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// GENERATED FILE. DO NOT EDIT DIRECTLY.
// Update hack/code/instancetype_testdata_gen.go and re-generate to edit
// You can add instance types by adding to the --instance-types CLI flag

`

var instanceTypesStr string
var outFile string

func init() {
	flag.StringVar(&instanceTypesStr, "instance-types", "", "comma-separated list of instance types to auto-generate static test data from")
	flag.StringVar(&outFile, "out-file", "zz_generated.describe_instance_types.go", "file to output the generated data")
	flag.Parse()
}

func main() {
	if err := os.Setenv("AWS_SDK_LOAD_CONFIG", "true"); err != nil {
		log.Fatalf("setting AWS_SDK_LOAD_CONFIG, %s", err)
	}
	if err := os.Setenv("AWS_REGION", "us-east-1"); err != nil {
		log.Fatalf("setting AWS_REGION, %s", err)
	}
	ctx := context.Background()
	sess := session.Must(session.NewSession())
	ec2Client := ec2.New(sess)
	instanceTypes := strings.Split(instanceTypesStr, ",")

	src := &bytes.Buffer{}
	fmt.Fprintln(src, "//go:build !ignore_autogenerated")
	license := lo.Must(os.ReadFile("hack/boilerplate.go.txt"))
	fmt.Fprintln(src, string(license))
	fmt.Fprint(src, packageHeader)
	fmt.Fprintln(src, getDescribeInstanceTypesOutput(ctx, ec2Client, instanceTypes))

	// Format and print to the file
	formatted, err := format.Source(src.Bytes())
	if err != nil {
		log.Fatalf("formatting generated source, %s", err)
	}
	if err := os.WriteFile(outFile, formatted, 0644); err != nil {
		log.Fatalf("writing output, %s", err)
	}
}

func getDescribeInstanceTypesOutput(ctx context.Context, ec2Client ec2iface.EC2API, instanceTypes []string) string {
	out, err := ec2Client.DescribeInstanceTypesWithContext(ctx, &ec2.DescribeInstanceTypesInput{
		InstanceTypes: aws.StringSlice(instanceTypes),
	})
	if err != nil {
		log.Fatalf("describing instance types, %s", err)
	}
	// Sort them by name so that we get a consistent ordering
	sort.SliceStable(out.InstanceTypes, func(i, j int) bool {
		return aws.StringValue(out.InstanceTypes[i].InstanceType) < aws.StringValue(out.InstanceTypes[j].InstanceType)
	})

	src := &bytes.Buffer{}
	fmt.Fprintln(src, "var defaultDescribeInstanceTypesOutput = &ec2.DescribeInstanceTypesOutput{")
	fmt.Fprintln(src, "InstanceTypes: []*ec2.InstanceTypeInfo{")
	for _, elem := range out.InstanceTypes {
		fmt.Fprintln(src, "{")
		data := getInstanceTypeInfo(elem)
		fmt.Fprintln(src, data)
		fmt.Fprintln(src, "},")
	}
	fmt.Fprintln(src, "},")
	fmt.Fprintln(src, "}")
	return src.String()
}

func getInstanceTypeInfo(info *ec2.InstanceTypeInfo) string {
	src := &bytes.Buffer{}
	fmt.Fprintf(src, "InstanceType: aws.String(\"%s\"),\n", lo.FromPtr(info.InstanceType))
	fmt.Fprintf(src, "SupportedUsageClasses: aws.StringSlice([]string{%s}),\n", getStringSliceData(info.SupportedUsageClasses))
	fmt.Fprintf(src, "SupportedVirtualizationTypes: aws.StringSlice([]string{%s}),\n", getStringSliceData(info.SupportedVirtualizationTypes))
	fmt.Fprintf(src, "BurstablePerformanceSupported: aws.Bool(%t),\n", lo.FromPtr(info.BurstablePerformanceSupported))
	fmt.Fprintf(src, "BareMetal: aws.Bool(%t),\n", lo.FromPtr(info.BareMetal))
	fmt.Fprintf(src, "Hypervisor: aws.String(\"%s\"),\n", lo.FromPtr(info.Hypervisor))
	fmt.Fprintf(src, "ProcessorInfo: &ec2.ProcessorInfo{\n")
	fmt.Fprintf(src, "SupportedArchitectures: aws.StringSlice([]string{%s}),\n", getStringSliceData(info.ProcessorInfo.SupportedArchitectures))
	fmt.Fprintf(src, "},\n")
	fmt.Fprintf(src, "VCpuInfo: &ec2.VCpuInfo{\n")
	fmt.Fprintf(src, "DefaultCores: aws.Int64(%d),\n", lo.FromPtr(info.VCpuInfo.DefaultCores))
	fmt.Fprintf(src, "DefaultVCpus: aws.Int64(%d),\n", lo.FromPtr(info.VCpuInfo.DefaultVCpus))
	fmt.Fprintf(src, "},\n")
	fmt.Fprintf(src, "MemoryInfo: &ec2.MemoryInfo{\n")
	fmt.Fprintf(src, "SizeInMiB: aws.Int64(%d),\n", lo.FromPtr(info.MemoryInfo.SizeInMiB))
	fmt.Fprintf(src, "},\n")

	if info.InferenceAcceleratorInfo != nil {
		fmt.Fprintf(src, "InferenceAcceleratorInfo: &ec2.InferenceAcceleratorInfo{\n")
		fmt.Fprintf(src, "Accelerators: []*ec2.InferenceDeviceInfo{\n")
		for _, elem := range info.InferenceAcceleratorInfo.Accelerators {
			fmt.Fprintf(src, getInferenceAcceleratorDeviceInfo(elem))
		}
		fmt.Fprintf(src, "},\n")
		fmt.Fprintf(src, "},\n")
	}
	if info.GpuInfo != nil {
		fmt.Fprintf(src, "GpuInfo: &ec2.GpuInfo{\n")
		fmt.Fprintf(src, "Gpus: []*ec2.GpuDeviceInfo{\n")
		for _, elem := range info.GpuInfo.Gpus {
			fmt.Fprintf(src, getGPUDeviceInfo(elem))
		}
		fmt.Fprintf(src, "},\n")
		fmt.Fprintf(src, "},\n")
	}
	if info.InstanceStorageInfo != nil {
		fmt.Fprintf(src, "InstanceStorageInfo: &ec2.InstanceStorageInfo{")
		fmt.Fprintf(src, "NvmeSupport: aws.String(\"%s\"),\n", lo.FromPtr(info.InstanceStorageInfo.NvmeSupport))
		fmt.Fprintf(src, "TotalSizeInGB: aws.Int64(%d),\n", lo.FromPtr(info.InstanceStorageInfo.TotalSizeInGB))
		fmt.Fprintf(src, "},\n")
	}
	fmt.Fprintf(src, "NetworkInfo: &ec2.NetworkInfo{\n")
	fmt.Fprintf(src, "MaximumNetworkInterfaces: aws.Int64(%d),\n", lo.FromPtr(info.NetworkInfo.MaximumNetworkInterfaces))
	fmt.Fprintf(src, "Ipv4AddressesPerInterface: aws.Int64(%d),\n", lo.FromPtr(info.NetworkInfo.Ipv4AddressesPerInterface))
	fmt.Fprintf(src, "EncryptionInTransitSupported: aws.Bool(%t),\n", lo.FromPtr(info.NetworkInfo.EncryptionInTransitSupported))
	fmt.Fprintf(src, "DefaultNetworkCardIndex: aws.Int64(%d),\n", lo.FromPtr(info.NetworkInfo.DefaultNetworkCardIndex))
	fmt.Fprintf(src, "NetworkCards: []*ec2.NetworkCardInfo{\n")
	for _, networkCard := range info.NetworkInfo.NetworkCards {
		fmt.Fprintf(src, getNetworkCardInfo(networkCard))
	}
	fmt.Fprintf(src, "},\n")
	fmt.Fprintf(src, "},\n")
	return src.String()
}

func getNetworkCardInfo(info *ec2.NetworkCardInfo) string {
	src := &bytes.Buffer{}
	fmt.Fprintf(src, "{\n")
	fmt.Fprintf(src, "NetworkCardIndex: aws.Int64(%d),\n", lo.FromPtr(info.NetworkCardIndex))
	fmt.Fprintf(src, "MaximumNetworkInterfaces: aws.Int64(%d),\n", lo.FromPtr(info.MaximumNetworkInterfaces))
	fmt.Fprintf(src, "},\n")
	return src.String()
}

func getInferenceAcceleratorDeviceInfo(info *ec2.InferenceDeviceInfo) string {
	src := &bytes.Buffer{}
	fmt.Fprintf(src, "{\n")
	fmt.Fprintf(src, "Name: aws.String(\"%s\"),\n", lo.FromPtr(info.Name))
	fmt.Fprintf(src, "Manufacturer: aws.String(\"%s\"),\n", lo.FromPtr(info.Manufacturer))
	fmt.Fprintf(src, "Count: aws.Int64(%d),\n", lo.FromPtr(info.Count))
	fmt.Fprintf(src, "},\n")
	return src.String()
}

func getGPUDeviceInfo(info *ec2.GpuDeviceInfo) string {
	src := &bytes.Buffer{}
	fmt.Fprintf(src, "{\n")
	fmt.Fprintf(src, "Name: aws.String(\"%s\"),\n", lo.FromPtr(info.Name))
	fmt.Fprintf(src, "Manufacturer: aws.String(\"%s\"),\n", lo.FromPtr(info.Manufacturer))
	fmt.Fprintf(src, "Count: aws.Int64(%d),\n", lo.FromPtr(info.Count))
	fmt.Fprintf(src, "MemoryInfo: &ec2.GpuDeviceMemoryInfo{\n")
	fmt.Fprintf(src, "SizeInMiB: aws.Int64(%d),\n", lo.FromPtr(info.MemoryInfo.SizeInMiB))
	fmt.Fprintf(src, "},\n")
	fmt.Fprintf(src, "},\n")
	return src.String()
}

func getStringSliceData(slice []*string) string {
	return strings.Join(lo.Map(slice, func(s *string, _ int) string { return fmt.Sprintf(`"%s"`, lo.FromPtr(s)) }), ",")
}
