//go:build !ignore_autogenerated

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

package pricing

// generated at 2024-12-02T13:14:37Z for us-east-1

import ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

var InitialOnDemandPricesUSGov = map[string]map[ec2types.InstanceType]float64{
	// us-gov-east-1
	"us-gov-east-1": {
		// c5 family
		"c5.12xlarge": 2.448000, "c5.18xlarge": 3.672000, "c5.24xlarge": 4.896000, "c5.2xlarge": 0.408000,
		"c5.4xlarge": 0.816000, "c5.9xlarge": 1.836000, "c5.large": 0.102000, "c5.metal": 4.896000,
		"c5.xlarge": 0.204000,
		// c5a family
		"c5a.12xlarge": 2.208000, "c5a.16xlarge": 2.944000, "c5a.24xlarge": 4.416000, "c5a.2xlarge": 0.368000,
		"c5a.4xlarge": 0.736000, "c5a.8xlarge": 1.472000, "c5a.large": 0.092000, "c5a.xlarge": 0.184000,
		// c5d family
		"c5d.18xlarge": 4.176000, "c5d.2xlarge": 0.464000, "c5d.4xlarge": 0.928000, "c5d.9xlarge": 2.088000,
		"c5d.large": 0.116000, "c5d.xlarge": 0.232000,
		// c5n family
		"c5n.18xlarge": 4.680000, "c5n.2xlarge": 0.520000, "c5n.4xlarge": 1.040000, "c5n.9xlarge": 2.340000,
		"c5n.large": 0.130000, "c5n.metal": 4.680000, "c5n.xlarge": 0.260000,
		// c6g family
		"c6g.12xlarge": 1.958400, "c6g.16xlarge": 2.611200, "c6g.2xlarge": 0.326400, "c6g.4xlarge": 0.652800,
		"c6g.8xlarge": 1.305600, "c6g.large": 0.081600, "c6g.medium": 0.040800, "c6g.metal": 2.767900,
		"c6g.xlarge": 0.163200,
		// c6gd family
		"c6gd.12xlarge": 2.227200, "c6gd.16xlarge": 2.969600, "c6gd.2xlarge": 0.371200, "c6gd.4xlarge": 0.742400,
		"c6gd.8xlarge": 1.484800, "c6gd.large": 0.092800, "c6gd.medium": 0.046400, "c6gd.metal": 2.969600,
		"c6gd.xlarge": 0.185600,
		// c6gn family
		"c6gn.12xlarge": 2.496000, "c6gn.16xlarge": 3.328000, "c6gn.2xlarge": 0.416000, "c6gn.4xlarge": 0.832000,
		"c6gn.8xlarge": 1.664000, "c6gn.large": 0.104000, "c6gn.medium": 0.052000, "c6gn.xlarge": 0.208000,
		// c6i family
		"c6i.12xlarge": 2.448000, "c6i.16xlarge": 3.264000, "c6i.24xlarge": 4.896000, "c6i.2xlarge": 0.408000,
		"c6i.32xlarge": 6.528000, "c6i.4xlarge": 0.816000, "c6i.8xlarge": 1.632000, "c6i.large": 0.102000,
		"c6i.metal": 6.528000, "c6i.xlarge": 0.204000,
		// c6in family
		"c6in.12xlarge": 3.276000, "c6in.16xlarge": 4.368000, "c6in.24xlarge": 6.552000, "c6in.2xlarge": 0.546000,
		"c6in.32xlarge": 8.736000, "c6in.4xlarge": 1.092000, "c6in.8xlarge": 2.184000, "c6in.large": 0.136500,
		"c6in.metal": 8.736000, "c6in.xlarge": 0.273000,
		// c7i family
		"c7i.12xlarge": 2.570400, "c7i.16xlarge": 3.427200, "c7i.24xlarge": 5.140800, "c7i.2xlarge": 0.428400,
		"c7i.48xlarge": 10.281600, "c7i.4xlarge": 0.856800, "c7i.8xlarge": 1.713600, "c7i.large": 0.107100,
		"c7i.metal-24xl": 5.654880, "c7i.metal-48xl": 10.281600, "c7i.xlarge": 0.214200,
		// d2 family
		"d2.2xlarge": 1.656000, "d2.4xlarge": 3.312000, "d2.8xlarge": 6.624000, "d2.xlarge": 0.828000,
		// g4dn family
		"g4dn.12xlarge": 4.931000, "g4dn.16xlarge": 5.486000, "g4dn.2xlarge": 0.948000, "g4dn.4xlarge": 1.518000,
		"g4dn.8xlarge": 2.743000, "g4dn.xlarge": 0.663000,
		// hpc6a family
		"hpc6a.48xlarge": 3.467000,
		// i3 family
		"i3.16xlarge": 6.016000, "i3.2xlarge": 0.752000, "i3.4xlarge": 1.504000, "i3.8xlarge": 3.008000,
		"i3.large": 0.188000, "i3.metal": 6.016000, "i3.xlarge": 0.376000,
		// i3en family
		"i3en.12xlarge": 6.552000, "i3en.24xlarge": 13.104000, "i3en.2xlarge": 1.092000, "i3en.3xlarge": 1.638000,
		"i3en.6xlarge": 3.276000, "i3en.large": 0.273000, "i3en.metal": 13.104000, "i3en.xlarge": 0.546000,
		// i4i family
		"i4i.12xlarge": 4.963000, "i4i.16xlarge": 6.618000, "i4i.24xlarge": 9.926400, "i4i.2xlarge": 0.827000,
		"i4i.32xlarge": 13.235200, "i4i.4xlarge": 1.654000, "i4i.8xlarge": 3.309000, "i4i.large": 0.207000,
		"i4i.metal": 13.235000, "i4i.xlarge": 0.414000,
		// inf1 family
		"inf1.24xlarge": 5.953000, "inf1.2xlarge": 0.456000, "inf1.6xlarge": 1.488000, "inf1.xlarge": 0.288000,
		// m5 family
		"m5.12xlarge": 2.904000, "m5.16xlarge": 3.872000, "m5.24xlarge": 5.808000, "m5.2xlarge": 0.484000,
		"m5.4xlarge": 0.968000, "m5.8xlarge": 1.936000, "m5.large": 0.121000, "m5.metal": 5.808000,
		"m5.xlarge": 0.242000,
		// m5a family
		"m5a.12xlarge": 2.616000, "m5a.16xlarge": 3.488000, "m5a.24xlarge": 5.232000, "m5a.2xlarge": 0.436000,
		"m5a.4xlarge": 0.872000, "m5a.8xlarge": 1.744000, "m5a.large": 0.109000, "m5a.xlarge": 0.218000,
		// m5d family
		"m5d.12xlarge": 3.432000, "m5d.16xlarge": 4.576000, "m5d.24xlarge": 6.864000, "m5d.2xlarge": 0.572000,
		"m5d.4xlarge": 1.144000, "m5d.8xlarge": 2.288000, "m5d.large": 0.143000, "m5d.metal": 6.864000,
		"m5d.xlarge": 0.286000,
		// m5dn family
		"m5dn.12xlarge": 4.104000, "m5dn.16xlarge": 5.472000, "m5dn.24xlarge": 8.208000, "m5dn.2xlarge": 0.684000,
		"m5dn.4xlarge": 1.368000, "m5dn.8xlarge": 2.736000, "m5dn.large": 0.171000, "m5dn.metal": 8.208000,
		"m5dn.xlarge": 0.342000,
		// m5n family
		"m5n.12xlarge": 3.576000, "m5n.16xlarge": 4.768000, "m5n.24xlarge": 7.152000, "m5n.2xlarge": 0.596000,
		"m5n.4xlarge": 1.192000, "m5n.8xlarge": 2.384000, "m5n.large": 0.149000, "m5n.metal": 7.152000,
		"m5n.xlarge": 0.298000,
		// m6g family
		"m6g.12xlarge": 2.323200, "m6g.16xlarge": 3.097600, "m6g.2xlarge": 0.387200, "m6g.4xlarge": 0.774400,
		"m6g.8xlarge": 1.548800, "m6g.large": 0.096800, "m6g.medium": 0.048400, "m6g.metal": 3.283500,
		"m6g.xlarge": 0.193600,
		// m6gd family
		"m6gd.12xlarge": 2.745600, "m6gd.16xlarge": 3.660800, "m6gd.2xlarge": 0.457600, "m6gd.4xlarge": 0.915200,
		"m6gd.8xlarge": 1.830400, "m6gd.large": 0.114400, "m6gd.medium": 0.057200, "m6gd.metal": 3.880400,
		"m6gd.xlarge": 0.228800,
		// m6i family
		"m6i.12xlarge": 2.904000, "m6i.16xlarge": 3.872000, "m6i.24xlarge": 5.808000, "m6i.2xlarge": 0.484000,
		"m6i.32xlarge": 7.744000, "m6i.4xlarge": 0.968000, "m6i.8xlarge": 1.936000, "m6i.large": 0.121000,
		"m6i.metal": 7.744000, "m6i.xlarge": 0.242000,
		// m7i-flex family
		"m7i-flex.2xlarge": 0.482800, "m7i-flex.4xlarge": 0.965600, "m7i-flex.8xlarge": 1.931200,
		"m7i-flex.large": 0.120700, "m7i-flex.xlarge": 0.241400,
		// m7i family
		"m7i.12xlarge": 3.049200, "m7i.16xlarge": 4.065600, "m7i.24xlarge": 6.098400, "m7i.2xlarge": 0.508200,
		"m7i.48xlarge": 12.196800, "m7i.4xlarge": 1.016400, "m7i.8xlarge": 2.032800, "m7i.large": 0.127050,
		"m7i.metal-24xl": 6.708240, "m7i.metal-48xl": 12.196800, "m7i.xlarge": 0.254100,
		// p3dn family
		"p3dn.24xlarge": 37.454000,
		// r5 family
		"r5.12xlarge": 3.624000, "r5.16xlarge": 4.832000, "r5.24xlarge": 7.248000, "r5.2xlarge": 0.604000,
		"r5.4xlarge": 1.208000, "r5.8xlarge": 2.416000, "r5.large": 0.151000, "r5.metal": 7.248000,
		"r5.xlarge": 0.302000,
		// r5a family
		"r5a.12xlarge": 3.264000, "r5a.16xlarge": 4.352000, "r5a.24xlarge": 6.528000, "r5a.2xlarge": 0.544000,
		"r5a.4xlarge": 1.088000, "r5a.8xlarge": 2.176000, "r5a.large": 0.136000, "r5a.xlarge": 0.272000,
		// r5d family
		"r5d.12xlarge": 4.152000, "r5d.16xlarge": 5.536000, "r5d.24xlarge": 8.304000, "r5d.2xlarge": 0.692000,
		"r5d.4xlarge": 1.384000, "r5d.8xlarge": 2.768000, "r5d.large": 0.173000, "r5d.metal": 8.304000,
		"r5d.xlarge": 0.346000,
		// r5dn family
		"r5dn.12xlarge": 4.824000, "r5dn.16xlarge": 6.432000, "r5dn.24xlarge": 9.648000, "r5dn.2xlarge": 0.804000,
		"r5dn.4xlarge": 1.608000, "r5dn.8xlarge": 3.216000, "r5dn.large": 0.201000, "r5dn.metal": 9.648000,
		"r5dn.xlarge": 0.402000,
		// r5n family
		"r5n.12xlarge": 4.296000, "r5n.16xlarge": 5.728000, "r5n.24xlarge": 8.592000, "r5n.2xlarge": 0.716000,
		"r5n.4xlarge": 1.432000, "r5n.8xlarge": 2.864000, "r5n.large": 0.179000, "r5n.metal": 8.592000,
		"r5n.xlarge": 0.358000,
		// r6g family
		"r6g.12xlarge": 2.899200, "r6g.16xlarge": 3.865600, "r6g.2xlarge": 0.483200, "r6g.4xlarge": 0.966400,
		"r6g.8xlarge": 1.932800, "r6g.large": 0.120800, "r6g.medium": 0.060400, "r6g.metal": 4.097500,
		"r6g.xlarge": 0.241600,
		// r6gd family
		"r6gd.12xlarge": 3.321600, "r6gd.16xlarge": 4.428800, "r6gd.2xlarge": 0.553600, "r6gd.4xlarge": 1.107200,
		"r6gd.8xlarge": 2.214400, "r6gd.large": 0.138400, "r6gd.medium": 0.069200, "r6gd.metal": 4.428800,
		"r6gd.xlarge": 0.276800,
		// r6i family
		"r6i.12xlarge": 3.624000, "r6i.16xlarge": 4.832000, "r6i.24xlarge": 7.248000, "r6i.2xlarge": 0.604000,
		"r6i.32xlarge": 9.664000, "r6i.4xlarge": 1.208000, "r6i.8xlarge": 2.416000, "r6i.large": 0.151000,
		"r6i.metal": 9.664000, "r6i.xlarge": 0.302000,
		// r7gd family
		"r7gd.12xlarge": 3.923500, "r7gd.16xlarge": 5.231400, "r7gd.2xlarge": 0.653900, "r7gd.4xlarge": 1.307800,
		"r7gd.8xlarge": 2.615700, "r7gd.large": 0.163500, "r7gd.medium": 0.081700, "r7gd.metal": 5.231400,
		"r7gd.xlarge": 0.327000,
		// r7i family
		"r7i.12xlarge": 3.805200, "r7i.16xlarge": 5.073600, "r7i.24xlarge": 7.610400, "r7i.2xlarge": 0.634200,
		"r7i.48xlarge": 15.220800, "r7i.4xlarge": 1.268400, "r7i.8xlarge": 2.536800, "r7i.large": 0.158550,
		"r7i.metal-24xl": 8.371440, "r7i.metal-48xl": 15.220800, "r7i.xlarge": 0.317100,
		// t3 family
		"t3.2xlarge": 0.390400, "t3.large": 0.097600, "t3.medium": 0.048800, "t3.micro": 0.012200,
		"t3.nano": 0.006100, "t3.small": 0.024400, "t3.xlarge": 0.195200,
		// t3a family
		"t3a.2xlarge": 0.351400, "t3a.large": 0.087800, "t3a.medium": 0.043900, "t3a.micro": 0.011000,
		"t3a.nano": 0.005500, "t3a.small": 0.022000, "t3a.xlarge": 0.175700,
		// t4g family
		"t4g.2xlarge": 0.313600, "t4g.large": 0.078400, "t4g.medium": 0.039200, "t4g.micro": 0.009800,
		"t4g.nano": 0.004900, "t4g.small": 0.019600, "t4g.xlarge": 0.156800,
		// u-12tb1 family
		"u-12tb1.112xlarge": 130.867000,
		// u-24tb1 family
		"u-24tb1.112xlarge": 261.730000,
		// u-6tb1 family
		"u-6tb1.112xlarge": 65.433000, "u-6tb1.56xlarge": 55.610750,
		// u-9tb1 family
		"u-9tb1.112xlarge": 98.150000,
		// x1 family
		"x1.16xlarge": 8.003000, "x1.32xlarge": 16.006000,
		// x1e family
		"x1e.16xlarge": 16.000000, "x1e.2xlarge": 2.000000, "x1e.32xlarge": 32.000000, "x1e.4xlarge": 4.000000,
		"x1e.8xlarge": 8.000000, "x1e.xlarge": 1.000000,
		// x2idn family
		"x2idn.16xlarge": 8.003000, "x2idn.24xlarge": 12.004500, "x2idn.32xlarge": 16.006000,
		"x2idn.metal": 16.006000,
		// x2iedn family
		"x2iedn.16xlarge": 16.006000, "x2iedn.24xlarge": 24.009000, "x2iedn.2xlarge": 2.000750,
		"x2iedn.32xlarge": 32.012000, "x2iedn.4xlarge": 4.001500, "x2iedn.8xlarge": 8.003000,
		"x2iedn.metal": 32.012000, "x2iedn.xlarge": 1.000380,
	},

	// us-gov-west-1
	"us-gov-west-1": {
		// c1 family
		"c1.medium": 0.157000, "c1.xlarge": 0.628000,
		// c3 family
		"c3.2xlarge": 0.504000, "c3.4xlarge": 1.008000, "c3.8xlarge": 2.016000, "c3.large": 0.126000,
		"c3.xlarge": 0.252000,
		// c4 family
		"c4.2xlarge": 0.479000, "c4.4xlarge": 0.958000, "c4.8xlarge": 1.915000, "c4.large": 0.120000,
		"c4.xlarge": 0.239000,
		// c5 family
		"c5.12xlarge": 2.448000, "c5.18xlarge": 3.672000, "c5.24xlarge": 4.896000, "c5.2xlarge": 0.408000,
		"c5.4xlarge": 0.816000, "c5.9xlarge": 1.836000, "c5.large": 0.102000, "c5.metal": 4.896000,
		"c5.xlarge": 0.204000,
		// c5a family
		"c5a.12xlarge": 2.208000, "c5a.16xlarge": 2.944000, "c5a.24xlarge": 4.416000, "c5a.2xlarge": 0.368000,
		"c5a.4xlarge": 0.736000, "c5a.8xlarge": 1.472000, "c5a.large": 0.092000, "c5a.xlarge": 0.184000,
		// c5d family
		"c5d.12xlarge": 2.784000, "c5d.18xlarge": 4.176000, "c5d.24xlarge": 5.568000, "c5d.2xlarge": 0.464000,
		"c5d.4xlarge": 0.928000, "c5d.9xlarge": 2.088000, "c5d.large": 0.116000, "c5d.metal": 5.568000,
		"c5d.xlarge": 0.232000,
		// c5n family
		"c5n.18xlarge": 4.680000, "c5n.2xlarge": 0.520000, "c5n.4xlarge": 1.040000, "c5n.9xlarge": 2.340000,
		"c5n.large": 0.130000, "c5n.metal": 4.680000, "c5n.xlarge": 0.260000,
		// c6g family
		"c6g.12xlarge": 1.958400, "c6g.16xlarge": 2.611200, "c6g.2xlarge": 0.326400, "c6g.4xlarge": 0.652800,
		"c6g.8xlarge": 1.305600, "c6g.large": 0.081600, "c6g.medium": 0.040800, "c6g.metal": 2.767900,
		"c6g.xlarge": 0.163200,
		// c6gd family
		"c6gd.12xlarge": 2.227200, "c6gd.16xlarge": 2.969600, "c6gd.2xlarge": 0.371200, "c6gd.4xlarge": 0.742400,
		"c6gd.8xlarge": 1.484800, "c6gd.large": 0.092800, "c6gd.medium": 0.046400, "c6gd.metal": 2.969600,
		"c6gd.xlarge": 0.185600,
		// c6gn family
		"c6gn.12xlarge": 2.496000, "c6gn.16xlarge": 3.328000, "c6gn.2xlarge": 0.416000, "c6gn.4xlarge": 0.832000,
		"c6gn.8xlarge": 1.664000, "c6gn.large": 0.104000, "c6gn.medium": 0.052000, "c6gn.xlarge": 0.208000,
		// c6i family
		"c6i.12xlarge": 2.448000, "c6i.16xlarge": 3.264000, "c6i.24xlarge": 4.896000, "c6i.2xlarge": 0.408000,
		"c6i.32xlarge": 6.528000, "c6i.4xlarge": 0.816000, "c6i.8xlarge": 1.632000, "c6i.large": 0.102000,
		"c6i.metal": 6.528000, "c6i.xlarge": 0.204000,
		// c6id family
		"c6id.12xlarge": 2.923200, "c6id.16xlarge": 3.897600, "c6id.24xlarge": 5.846400, "c6id.2xlarge": 0.487200,
		"c6id.32xlarge": 7.795200, "c6id.4xlarge": 0.974400, "c6id.8xlarge": 1.948800, "c6id.large": 0.121800,
		"c6id.metal": 7.795200, "c6id.xlarge": 0.243600,
		// c6in family
		"c6in.12xlarge": 3.276000, "c6in.16xlarge": 4.368000, "c6in.24xlarge": 6.552000, "c6in.2xlarge": 0.546000,
		"c6in.32xlarge": 8.736000, "c6in.4xlarge": 1.092000, "c6in.8xlarge": 2.184000, "c6in.large": 0.136500,
		"c6in.metal": 8.736000, "c6in.xlarge": 0.273000,
		// cc2 family
		"cc2.8xlarge": 2.250000,
		// d2 family
		"d2.2xlarge": 1.656000, "d2.4xlarge": 3.312000, "d2.8xlarge": 6.624000, "d2.xlarge": 0.828000,
		// d3 family
		"d3.2xlarge": 1.197000, "d3.4xlarge": 2.394000, "d3.8xlarge": 4.787760, "d3.xlarge": 0.598000,
		// f1 family
		"f1.16xlarge": 15.840000, "f1.2xlarge": 1.980000, "f1.4xlarge": 3.960000,
		// g3 family
		"g3.16xlarge": 5.280000, "g3.4xlarge": 1.320000, "g3.8xlarge": 2.640000,
		// g3s family
		"g3s.xlarge": 0.868000,
		// g4dn family
		"g4dn.12xlarge": 4.931000, "g4dn.16xlarge": 5.486000, "g4dn.2xlarge": 0.948000, "g4dn.4xlarge": 1.518000,
		"g4dn.8xlarge": 2.743000, "g4dn.metal": 9.862000, "g4dn.xlarge": 0.663000,
		// g6 family
		"g6.12xlarge": 5.800030, "g6.16xlarge": 4.281450, "g6.24xlarge": 8.413670, "g6.2xlarge": 1.232200,
		"g6.48xlarge": 16.827340, "g6.4xlarge": 1.667810, "g6.8xlarge": 2.539030, "g6.xlarge": 1.014400,
		// gr6 family
		"gr6.4xlarge": 1.940100, "gr6.8xlarge": 3.083590,
		// hpc6a family
		"hpc6a.48xlarge": 3.467000,
		// hpc6id family
		"hpc6id.32xlarge": 6.854400,
		// hpc7a family
		"hpc7a.12xlarge": 8.667400, "hpc7a.24xlarge": 8.667400, "hpc7a.48xlarge": 8.667400,
		"hpc7a.96xlarge": 8.667400,
		// hpc7g family
		"hpc7g.16xlarge": 2.026200, "hpc7g.4xlarge": 2.026200, "hpc7g.8xlarge": 2.026200,
		// hs1 family
		"hs1.8xlarge": 5.520000,
		// i2 family
		"i2.2xlarge": 2.046000, "i2.4xlarge": 4.092000, "i2.8xlarge": 8.184000, "i2.xlarge": 1.023000,
		// i3 family
		"i3.16xlarge": 6.016000, "i3.2xlarge": 0.752000, "i3.4xlarge": 1.504000, "i3.8xlarge": 3.008000,
		"i3.large": 0.188000, "i3.metal": 6.016000, "i3.xlarge": 0.376000,
		// i3en family
		"i3en.12xlarge": 6.552000, "i3en.24xlarge": 13.104000, "i3en.2xlarge": 1.092000, "i3en.3xlarge": 1.638000,
		"i3en.6xlarge": 3.276000, "i3en.large": 0.273000, "i3en.metal": 13.104000, "i3en.xlarge": 0.546000,
		// i3p family
		"i3p.16xlarge": 6.016000,
		// i4i family
		"i4i.12xlarge": 4.963000, "i4i.16xlarge": 6.618000, "i4i.24xlarge": 9.926400, "i4i.2xlarge": 0.827000,
		"i4i.32xlarge": 13.235200, "i4i.4xlarge": 1.654000, "i4i.8xlarge": 3.309000, "i4i.large": 0.207000,
		"i4i.metal": 13.235000, "i4i.xlarge": 0.414000,
		// inf1 family
		"inf1.24xlarge": 5.953000, "inf1.2xlarge": 0.456000, "inf1.6xlarge": 1.488000, "inf1.xlarge": 0.288000,
		// m1 family
		"m1.large": 0.211000, "m1.medium": 0.106000, "m1.small": 0.053000, "m1.xlarge": 0.423000,
		// m2 family
		"m2.2xlarge": 0.586000, "m2.4xlarge": 1.171000, "m2.xlarge": 0.293000,
		// m3 family
		"m3.2xlarge": 0.672000, "m3.large": 0.168000, "m3.medium": 0.084000, "m3.xlarge": 0.336000,
		// m4 family
		"m4.10xlarge": 2.520000, "m4.16xlarge": 4.032000, "m4.2xlarge": 0.504000, "m4.4xlarge": 1.008000,
		"m4.large": 0.126000, "m4.xlarge": 0.252000,
		// m5 family
		"m5.12xlarge": 2.904000, "m5.16xlarge": 3.872000, "m5.24xlarge": 5.808000, "m5.2xlarge": 0.484000,
		"m5.4xlarge": 0.968000, "m5.8xlarge": 1.936000, "m5.large": 0.121000, "m5.metal": 5.808000,
		"m5.xlarge": 0.242000,
		// m5a family
		"m5a.12xlarge": 2.616000, "m5a.16xlarge": 3.488000, "m5a.24xlarge": 5.232000, "m5a.2xlarge": 0.436000,
		"m5a.4xlarge": 0.872000, "m5a.8xlarge": 1.744000, "m5a.large": 0.109000, "m5a.xlarge": 0.218000,
		// m5ad family
		"m5ad.12xlarge": 3.144000, "m5ad.16xlarge": 4.192000, "m5ad.24xlarge": 6.288000, "m5ad.2xlarge": 0.524000,
		"m5ad.4xlarge": 1.048000, "m5ad.8xlarge": 2.096000, "m5ad.large": 0.131000, "m5ad.xlarge": 0.262000,
		// m5d family
		"m5d.12xlarge": 3.432000, "m5d.16xlarge": 4.576000, "m5d.24xlarge": 6.864000, "m5d.2xlarge": 0.572000,
		"m5d.4xlarge": 1.144000, "m5d.8xlarge": 2.288000, "m5d.large": 0.143000, "m5d.metal": 6.864000,
		"m5d.xlarge": 0.286000,
		// m5dn family
		"m5dn.12xlarge": 4.104000, "m5dn.16xlarge": 5.472000, "m5dn.24xlarge": 8.208000, "m5dn.2xlarge": 0.684000,
		"m5dn.4xlarge": 1.368000, "m5dn.8xlarge": 2.736000, "m5dn.large": 0.171000, "m5dn.metal": 8.208000,
		"m5dn.xlarge": 0.342000,
		// m5n family
		"m5n.12xlarge": 3.576000, "m5n.16xlarge": 4.768000, "m5n.24xlarge": 7.152000, "m5n.2xlarge": 0.596000,
		"m5n.4xlarge": 1.192000, "m5n.8xlarge": 2.384000, "m5n.large": 0.149000, "m5n.metal": 7.152000,
		"m5n.xlarge": 0.298000,
		// m6g family
		"m6g.12xlarge": 2.323200, "m6g.16xlarge": 3.097600, "m6g.2xlarge": 0.387200, "m6g.4xlarge": 0.774400,
		"m6g.8xlarge": 1.548800, "m6g.large": 0.096800, "m6g.medium": 0.048400, "m6g.metal": 3.283500,
		"m6g.xlarge": 0.193600,
		// m6gd family
		"m6gd.12xlarge": 2.745600, "m6gd.16xlarge": 3.660800, "m6gd.2xlarge": 0.457600, "m6gd.4xlarge": 0.915200,
		"m6gd.8xlarge": 1.830400, "m6gd.large": 0.114400, "m6gd.medium": 0.057200, "m6gd.metal": 3.880400,
		"m6gd.xlarge": 0.228800,
		// m6i family
		"m6i.12xlarge": 2.904000, "m6i.16xlarge": 3.872000, "m6i.24xlarge": 5.808000, "m6i.2xlarge": 0.484000,
		"m6i.32xlarge": 7.744000, "m6i.4xlarge": 0.968000, "m6i.8xlarge": 1.936000, "m6i.large": 0.121000,
		"m6i.metal": 7.744000, "m6i.xlarge": 0.242000,
		// m6id family
		"m6id.12xlarge": 3.604800, "m6id.16xlarge": 4.806400, "m6id.24xlarge": 7.209600, "m6id.2xlarge": 0.600800,
		"m6id.32xlarge": 9.612800, "m6id.4xlarge": 1.201600, "m6id.8xlarge": 2.403200, "m6id.large": 0.150200,
		"m6id.metal": 9.612800, "m6id.xlarge": 0.300400,
		// m6idn family
		"m6idn.12xlarge": 4.801680, "m6idn.16xlarge": 6.402240, "m6idn.24xlarge": 9.603360,
		"m6idn.2xlarge": 0.800280, "m6idn.32xlarge": 12.804480, "m6idn.4xlarge": 1.600560, "m6idn.8xlarge": 3.201120,
		"m6idn.large": 0.200070, "m6idn.metal": 12.804480, "m6idn.xlarge": 0.400140,
		// m6in family
		"m6in.12xlarge": 4.183920, "m6in.16xlarge": 5.578560, "m6in.24xlarge": 8.367840, "m6in.2xlarge": 0.697320,
		"m6in.32xlarge": 11.157120, "m6in.4xlarge": 1.394640, "m6in.8xlarge": 2.789280, "m6in.large": 0.174330,
		"m6in.metal": 11.157120, "m6in.xlarge": 0.348660,
		// m7i-flex family
		"m7i-flex.2xlarge": 0.482800, "m7i-flex.4xlarge": 0.965600, "m7i-flex.8xlarge": 1.931200,
		"m7i-flex.large": 0.120700, "m7i-flex.xlarge": 0.241400,
		// m7i family
		"m7i.12xlarge": 3.049200, "m7i.16xlarge": 4.065600, "m7i.24xlarge": 6.098400, "m7i.2xlarge": 0.508200,
		"m7i.48xlarge": 12.196800, "m7i.4xlarge": 1.016400, "m7i.8xlarge": 2.032800, "m7i.large": 0.127050,
		"m7i.metal-24xl": 6.708240, "m7i.metal-48xl": 12.196800, "m7i.xlarge": 0.254100,
		// p2 family
		"p2.16xlarge": 17.280000, "p2.8xlarge": 8.640000, "p2.xlarge": 1.080000,
		// p3 family
		"p3.16xlarge": 29.376000, "p3.2xlarge": 3.672000, "p3.8xlarge": 14.688000,
		// p3dn family
		"p3dn.24xlarge": 37.454000,
		// p4d family
		"p4d.24xlarge": 39.330000,
		// p5 family
		"p5.48xlarge": 117.984000,
		// r3 family
		"r3.2xlarge": 0.798000, "r3.4xlarge": 1.596000, "r3.8xlarge": 3.192000, "r3.large": 0.200000,
		"r3.xlarge": 0.399000,
		// r4 family
		"r4.16xlarge": 5.107200, "r4.2xlarge": 0.638400, "r4.4xlarge": 1.276800, "r4.8xlarge": 2.553600,
		"r4.large": 0.159600, "r4.xlarge": 0.319200,
		// r5 family
		"r5.12xlarge": 3.624000, "r5.16xlarge": 4.832000, "r5.24xlarge": 7.248000, "r5.2xlarge": 0.604000,
		"r5.4xlarge": 1.208000, "r5.8xlarge": 2.416000, "r5.large": 0.151000, "r5.metal": 7.248000,
		"r5.xlarge": 0.302000,
		// r5a family
		"r5a.12xlarge": 3.264000, "r5a.16xlarge": 4.352000, "r5a.24xlarge": 6.528000, "r5a.2xlarge": 0.544000,
		"r5a.4xlarge": 1.088000, "r5a.8xlarge": 2.176000, "r5a.large": 0.136000, "r5a.xlarge": 0.272000,
		// r5ad family
		"r5ad.12xlarge": 3.792000, "r5ad.16xlarge": 5.056000, "r5ad.24xlarge": 7.584000, "r5ad.2xlarge": 0.632000,
		"r5ad.4xlarge": 1.264000, "r5ad.8xlarge": 2.528000, "r5ad.large": 0.158000, "r5ad.xlarge": 0.316000,
		// r5d family
		"r5d.12xlarge": 4.152000, "r5d.16xlarge": 5.536000, "r5d.24xlarge": 8.304000, "r5d.2xlarge": 0.692000,
		"r5d.4xlarge": 1.384000, "r5d.8xlarge": 2.768000, "r5d.large": 0.173000, "r5d.metal": 8.304000,
		"r5d.xlarge": 0.346000,
		// r5dn family
		"r5dn.12xlarge": 4.824000, "r5dn.16xlarge": 6.432000, "r5dn.24xlarge": 9.648000, "r5dn.2xlarge": 0.804000,
		"r5dn.4xlarge": 1.608000, "r5dn.8xlarge": 3.216000, "r5dn.large": 0.201000, "r5dn.metal": 9.648000,
		"r5dn.xlarge": 0.402000,
		// r5n family
		"r5n.12xlarge": 4.296000, "r5n.16xlarge": 5.728000, "r5n.24xlarge": 8.592000, "r5n.2xlarge": 0.716000,
		"r5n.4xlarge": 1.432000, "r5n.8xlarge": 2.864000, "r5n.large": 0.179000, "r5n.metal": 8.592000,
		"r5n.xlarge": 0.358000,
		// r6g family
		"r6g.12xlarge": 2.899200, "r6g.16xlarge": 3.865600, "r6g.2xlarge": 0.483200, "r6g.4xlarge": 0.966400,
		"r6g.8xlarge": 1.932800, "r6g.large": 0.120800, "r6g.medium": 0.060400, "r6g.metal": 4.097500,
		"r6g.xlarge": 0.241600,
		// r6gd family
		"r6gd.12xlarge": 3.321600, "r6gd.16xlarge": 4.428800, "r6gd.2xlarge": 0.553600, "r6gd.4xlarge": 1.107200,
		"r6gd.8xlarge": 2.214400, "r6gd.large": 0.138400, "r6gd.medium": 0.069200, "r6gd.metal": 4.428800,
		"r6gd.xlarge": 0.276800,
		// r6i family
		"r6i.12xlarge": 3.624000, "r6i.16xlarge": 4.832000, "r6i.24xlarge": 7.248000, "r6i.2xlarge": 0.604000,
		"r6i.32xlarge": 9.664000, "r6i.4xlarge": 1.208000, "r6i.8xlarge": 2.416000, "r6i.large": 0.151000,
		"r6i.metal": 9.664000, "r6i.xlarge": 0.302000,
		// r6id family
		"r6id.12xlarge": 4.360800, "r6id.16xlarge": 5.814400, "r6id.24xlarge": 8.721600, "r6id.2xlarge": 0.726800,
		"r6id.32xlarge": 11.628800, "r6id.4xlarge": 1.453600, "r6id.8xlarge": 2.907200, "r6id.large": 0.181700,
		"r6id.metal": 11.628800, "r6id.xlarge": 0.363400,
		// r6idn family
		"r6idn.12xlarge": 5.644080, "r6idn.16xlarge": 7.525440, "r6idn.24xlarge": 11.288160,
		"r6idn.2xlarge": 0.940680, "r6idn.32xlarge": 15.050880, "r6idn.4xlarge": 1.881360, "r6idn.8xlarge": 3.762720,
		"r6idn.large": 0.235170, "r6idn.metal": 15.050880, "r6idn.xlarge": 0.470340,
		// r6in family
		"r6in.12xlarge": 5.026320, "r6in.16xlarge": 6.701760, "r6in.24xlarge": 10.052640, "r6in.2xlarge": 0.837720,
		"r6in.32xlarge": 13.403520, "r6in.4xlarge": 1.675440, "r6in.8xlarge": 3.350880, "r6in.large": 0.209430,
		"r6in.metal": 13.403520, "r6in.xlarge": 0.418860,
		// r7gd family
		"r7gd.12xlarge": 3.925000, "r7gd.16xlarge": 5.233300, "r7gd.2xlarge": 0.654200, "r7gd.4xlarge": 1.308300,
		"r7gd.8xlarge": 2.616600, "r7gd.large": 0.163500, "r7gd.medium": 0.081800, "r7gd.metal": 5.233300,
		"r7gd.xlarge": 0.327100,
		// r7i family
		"r7i.12xlarge": 3.805200, "r7i.16xlarge": 5.073600, "r7i.24xlarge": 7.610400, "r7i.2xlarge": 0.634200,
		"r7i.48xlarge": 15.220800, "r7i.4xlarge": 1.268400, "r7i.8xlarge": 2.536800, "r7i.large": 0.158550,
		"r7i.metal-24xl": 8.371440, "r7i.metal-48xl": 15.220800, "r7i.xlarge": 0.317100,
		// t1 family
		"t1.micro": 0.024000,
		// t2 family
		"t2.2xlarge": 0.435200, "t2.large": 0.108800, "t2.medium": 0.054400, "t2.micro": 0.013600,
		"t2.nano": 0.006800, "t2.small": 0.027200, "t2.xlarge": 0.217600,
		// t3 family
		"t3.2xlarge": 0.390400, "t3.large": 0.097600, "t3.medium": 0.048800, "t3.micro": 0.012200,
		"t3.nano": 0.006100, "t3.small": 0.024400, "t3.xlarge": 0.195200,
		// t3a family
		"t3a.2xlarge": 0.351400, "t3a.large": 0.087800, "t3a.medium": 0.043900, "t3a.micro": 0.011000,
		"t3a.nano": 0.005500, "t3a.small": 0.022000, "t3a.xlarge": 0.175700,
		// t4g family
		"t4g.2xlarge": 0.313600, "t4g.large": 0.078400, "t4g.medium": 0.039200, "t4g.micro": 0.009800,
		"t4g.nano": 0.004900, "t4g.small": 0.019600, "t4g.xlarge": 0.156800,
		// u-12tb1 family
		"u-12tb1.112xlarge": 130.867000,
		// u-24tb1 family
		"u-24tb1.112xlarge": 261.730000,
		// u-3tb1 family
		"u-3tb1.56xlarge": 32.716500,
		// u-6tb1 family
		"u-6tb1.112xlarge": 65.433000, "u-6tb1.56xlarge": 55.610750,
		// u-9tb1 family
		"u-9tb1.112xlarge": 98.150000,
		// x1 family
		"x1.16xlarge": 8.003000, "x1.32xlarge": 16.006000,
		// x1e family
		"x1e.16xlarge": 16.000000, "x1e.2xlarge": 2.000000, "x1e.32xlarge": 32.000000, "x1e.4xlarge": 4.000000,
		"x1e.8xlarge": 8.000000, "x1e.xlarge": 1.000000,
		// x2idn family
		"x2idn.16xlarge": 8.003000, "x2idn.24xlarge": 12.004500, "x2idn.32xlarge": 16.006000,
		"x2idn.metal": 16.006000,
		// x2iedn family
		"x2iedn.16xlarge": 16.006000, "x2iedn.24xlarge": 24.009000, "x2iedn.2xlarge": 2.000750,
		"x2iedn.32xlarge": 32.012000, "x2iedn.4xlarge": 4.001500, "x2iedn.8xlarge": 8.003000,
		"x2iedn.metal": 32.012000, "x2iedn.xlarge": 1.000380,
	},
}
