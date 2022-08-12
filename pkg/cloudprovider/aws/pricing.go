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

package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"knative.dev/pkg/ptr"

	"go.uber.org/multierr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/aws/aws-sdk-go/service/pricing/pricingiface"
	"github.com/samber/lo"
	"knative.dev/pkg/logging"
)

// PricingProvider provides actual pricing data to the AWS cloud provider to allow it to make more informed decisions
// regarding which instances to launch.  This is initialized at startup with a periodically updated static price list to
// support running in locations where pricing data is unavailable.  In those cases the static pricing data provides a
// relative ordering that is still more accurate than our previous pricing model.  In the event that a pricing update
// fails, the previous pricing information is retained and used which may be the static initial pricing data if pricing
// updates never succeed.
type PricingProvider struct {
	ec2     ec2iface.EC2API
	pricing pricingiface.PricingAPI
	region  string

	mu                 sync.RWMutex
	onDemandUpdateTime time.Time
	onDemandPrices     map[string]priceData
	spotUpdateTime     time.Time
	spotPrices         map[string]zonalPricing
}

type zonalPricing struct {
	defaultPrice float64 // Used until we get the spot pricing data
	prices       map[string]priceData
}

// priceData captures the last known price for an offering
// and tells us whether this data was retrieved in the last pricing update loop
// or whether it is stale data
type priceData struct {
	lastKnownPrice float64
	isCurrent      bool
}

func newZonalPricing(defaultPrice *float64) zonalPricing {
	z := zonalPricing{
		prices: map[string]priceData{},
	}
	if defaultPrice != nil {
		z.defaultPrice = ptr.Float64Value(defaultPrice)
	}
	return z
}

// pricingUpdatePeriod is how often we try to update our pricing information after the initial update on startup
const pricingUpdatePeriod = 12 * time.Hour

// NewPricingAPI returns a pricing API configured based on a particular region
func NewPricingAPI(sess *session.Session, region string) pricingiface.PricingAPI {
	if sess == nil {
		return nil
	}
	// pricing API doesn't have an endpoint in all regions
	pricingAPIRegion := "us-east-1"
	if strings.HasPrefix(region, "cn-") {
		pricingAPIRegion = "cn-north-1"
	} else if strings.HasPrefix(region, "ap-") {
		pricingAPIRegion = "ap-south-1"
	}
	return pricing.New(sess, &aws.Config{Region: aws.String(pricingAPIRegion)})
}

func NewPricingProvider(ctx context.Context, pricing pricingiface.PricingAPI, ec2Api ec2iface.EC2API, region string, isolatedVPC bool, startAsync <-chan struct{}) *PricingProvider {
	p := &PricingProvider{
		region:             region,
		onDemandUpdateTime: initialPriceUpdate,
		onDemandPrices:     populateInitialOnDemandPricing(initialOnDemandPrices),
		spotUpdateTime:     initialPriceUpdate,
		// default our spot pricing to the same as the on-demand pricing until a price update
		spotPrices: populateInitialSpotPricing(initialOnDemandPrices),
		ec2:        ec2Api,
		pricing:    pricing,
	}
	ctx = logging.WithLogger(ctx, logging.FromContext(ctx).Named("pricing"))

	if isolatedVPC {
		logging.FromContext(ctx).Infof("Assuming isolated VPC, pricing information will not be updated")
	} else {
		go func() {
			// perform an initial price update at startup
			p.updatePricing(ctx)

			startup := time.Now()
			// wait for leader election or to be signaled to exit
			select {
			case <-startAsync:
			case <-ctx.Done():
				return
			}
			// if it took many hours to be elected leader, we want to re-fetch pricing before we start our periodic
			// polling
			if time.Since(startup) > pricingUpdatePeriod {
				p.updatePricing(ctx)
			}

			for {
				select {
				case <-ctx.Done():
					return
				case <-time.After(pricingUpdatePeriod):
					p.updatePricing(ctx)
				}
			}
		}()
	}
	return p
}

// InstanceTypes returns the list of all instance types for which either a spot or on-demand price is known.
func (p *PricingProvider) InstanceTypes() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return lo.Union(lo.Keys(p.onDemandPrices), lo.Keys(p.spotPrices))
}

// OnDemandLastUpdated returns the time that the on-demand pricing was last updated
func (p *PricingProvider) OnDemandLastUpdated() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.onDemandUpdateTime
}

// SpotLastUpdated returns the time that the spot pricing was last updated
func (p *PricingProvider) SpotLastUpdated() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.spotUpdateTime
}

// OnDemandPrice returns the last known on-demand price for a given instance type, returning an error if there is no
// known on-demand pricing for the instance type.
func (p *PricingProvider) OnDemandPrice(instanceType string) (float64, bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	price, ok := p.onDemandPrices[instanceType]
	if !ok {
		return 0.0, false, fmt.Errorf("instance type %s not found", instanceType)
	}
	return price.lastKnownPrice, price.isCurrent, nil
}

// SpotPrice returns the last known spot price for a given instance type and zone, returning an error
// if there is no known spot pricing for that instance type or zone
func (p *PricingProvider) SpotPrice(instanceType string, zone string) (float64, bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if val, ok := p.spotPrices[instanceType]; ok {
		if p.spotUpdateTime.Equal(initialPriceUpdate) {
			return val.defaultPrice, true, nil
		}
		if price, ok := p.spotPrices[instanceType].prices[zone]; ok {
			return price.lastKnownPrice, price.isCurrent, nil
		}
		return 0.0, false, fmt.Errorf("instance type %s not found in zone %s", instanceType, zone)
	}
	return 0.0, false, fmt.Errorf("instance type %s not found", instanceType)
}

func (p *PricingProvider) updatePricing(ctx context.Context) {
	logging.FromContext(ctx).Infof("Updating EC2 pricing information")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.updateOnDemandPricing(ctx); err != nil {
			logging.FromContext(ctx).Errorf("updating on-demand pricing, %s, using existing pricing data from %s", err, p.onDemandUpdateTime.Format(time.RFC3339))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := p.updateSpotPricing(ctx); err != nil {
			logging.FromContext(ctx).Errorf("updating spot pricing, %s, using existing pricing data from %s", err, p.spotUpdateTime.Format(time.RFC3339))
		}
	}()

	wg.Wait()
}

func (p *PricingProvider) updateOnDemandPricing(ctx context.Context) error {
	// standard on-demand instances
	var wg sync.WaitGroup
	var onDemandPrices, onDemandMetalPrices map[string]float64
	var onDemandErr, onDemandMetalErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		onDemandPrices, onDemandErr = p.fetchOnDemandPricing(ctx,
			&pricing.Filter{
				Field: aws.String("tenancy"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String("Shared"),
			},
			&pricing.Filter{
				Field: aws.String("productFamily"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String("Compute Instance"),
			})
	}()

	// bare metal on-demand prices
	wg.Add(1)
	go func() {
		defer wg.Done()
		onDemandMetalPrices, onDemandMetalErr = p.fetchOnDemandPricing(ctx,
			&pricing.Filter{
				Field: aws.String("tenancy"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String("Dedicated"),
			},
			&pricing.Filter{
				Field: aws.String("productFamily"),
				Type:  aws.String("TERM_MATCH"),
				Value: aws.String("Compute Instance (bare metal)"),
			})
	}()

	wg.Wait()
	err := multierr.Append(onDemandErr, onDemandMetalErr)
	if err != nil {
		return err
	}

	if len(onDemandPrices) == 0 || len(onDemandMetalPrices) == 0 {
		return errors.New("no on-demand pricing found")
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	p.markODNotCurrent()
	currentODPrices := lo.Assign(onDemandPrices, onDemandMetalPrices)
	for k, v := range currentODPrices {
		p.onDemandPrices[k] = priceData{
			lastKnownPrice: v,
			isCurrent:      true,
		}
	}

	p.onDemandUpdateTime = time.Now()
	logging.FromContext(ctx).Infof("updated on-demand pricing with %d instance types", len(p.onDemandPrices))
	return nil
}

func (p *PricingProvider) fetchOnDemandPricing(ctx context.Context, additionalFilters ...*pricing.Filter) (map[string]float64, error) {
	prices := map[string]float64{}
	filters := append([]*pricing.Filter{
		{
			Field: aws.String("regionCode"),
			Type:  aws.String("TERM_MATCH"),
			Value: aws.String(p.region),
		},
		{
			Field: aws.String("serviceCode"),
			Type:  aws.String("TERM_MATCH"),
			Value: aws.String("AmazonEC2"),
		},
		{
			Field: aws.String("preInstalledSw"),
			Type:  aws.String("TERM_MATCH"),
			Value: aws.String("NA"),
		},
		{
			Field: aws.String("operatingSystem"),
			Type:  aws.String("TERM_MATCH"),
			Value: aws.String("Linux"),
		},
		{
			Field: aws.String("capacitystatus"),
			Type:  aws.String("TERM_MATCH"),
			Value: aws.String("Used"),
		},
		{
			Field: aws.String("marketoption"),
			Type:  aws.String("TERM_MATCH"),
			Value: aws.String("OnDemand"),
		}},
		additionalFilters...)
	if err := p.pricing.GetProductsPagesWithContext(ctx, &pricing.GetProductsInput{
		Filters:     filters,
		ServiceCode: aws.String("AmazonEC2")}, p.onDemandPage(prices)); err != nil {
		return nil, err
	}
	return prices, nil
}

// turning off cyclo here, it measures as a 12 due to all of the type checks of the pricing data which returns a deeply
// nested map[string]interface{}
// nolint: gocyclo
func (p *PricingProvider) onDemandPage(prices map[string]float64) func(output *pricing.GetProductsOutput, b bool) bool {
	// this isn't the full pricing struct, just the portions we care about
	type priceItem struct {
		Product struct {
			Attributes struct {
				InstanceType string
			}
		}
		Terms struct {
			OnDemand map[string]struct {
				PriceDimensions map[string]struct {
					PricePerUnit struct {
						USD string
					}
				}
			}
		}
	}

	return func(output *pricing.GetProductsOutput, b bool) bool {
		for _, outer := range output.PriceList {
			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			if err := enc.Encode(outer); err != nil {
				logging.FromContext(context.Background()).Errorf("encoding %s", err)
			}
			dec := json.NewDecoder(&buf)
			var pItem priceItem
			if err := dec.Decode(&pItem); err != nil {
				logging.FromContext(context.Background()).Errorf("decoding %s", err)
			}
			if pItem.Product.Attributes.InstanceType == "" {
				continue
			}
			for _, term := range pItem.Terms.OnDemand {
				for _, v := range term.PriceDimensions {
					price, err := strconv.ParseFloat(v.PricePerUnit.USD, 64)
					if err != nil || price == 0 {
						continue
					}
					prices[pItem.Product.Attributes.InstanceType] = price
				}
			}
		}
		return true
	}
}

// nolint: gocyclo
func (p *PricingProvider) updateSpotPricing(ctx context.Context) error {
	type pricingInfo struct {
		timestamp time.Time
		price     float64
	}

	totalOfferings := 0

	prices := map[string]map[string]*pricingInfo{}
	if err := p.ec2.DescribeSpotPriceHistoryPagesWithContext(ctx, &ec2.DescribeSpotPriceHistoryInput{
		ProductDescriptions: []*string{aws.String("Linux/UNIX")},
		// look for spot prices for the past day
		StartTime: aws.Time(time.Now().Add(24 * time.Hour)),
	}, func(output *ec2.DescribeSpotPriceHistoryOutput, b bool) bool {
		for _, sph := range output.SpotPriceHistory {
			spotPriceStr := aws.StringValue(sph.SpotPrice)
			spotPrice, err := strconv.ParseFloat(spotPriceStr, 64)
			// these errors shouldn't occur, but if pricing API does have an error, we ignore the record
			if err != nil {
				logging.FromContext(ctx).Debugf("unable to parse price record %#v", sph)
				continue
			}
			if sph.Timestamp == nil {
				continue
			}
			instanceType := aws.StringValue(sph.InstanceType)
			az := aws.StringValue(sph.AvailabilityZone)
			timeStamp := *sph.Timestamp

			// pricing can vary based on the sph.AvailabilityZone, but we just currently take the latest update
			// on a per-zone basis
			_, ok := prices[instanceType]
			if !ok {
				prices[instanceType] = make(map[string]*pricingInfo)
			}
			existing, ok := prices[instanceType][az]
			if !ok || timeStamp.After(existing.timestamp) {
				prices[instanceType][az] = &pricingInfo{
					timestamp: timeStamp,
					price:     spotPrice,
				}
			}
		}
		return true
	}); err != nil {
		return err
	}
	if len(prices) == 0 {
		return errors.New("no spot pricing found")
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	p.markSpotNotCurrent()
	for it, zoneData := range prices {
		if _, ok := p.spotPrices[it]; !ok {
			p.spotPrices[it] = newZonalPricing(nil)
		}
		for zone, data := range zoneData {
			p.spotPrices[it].prices[zone] = priceData{
				lastKnownPrice: data.price,
				isCurrent:      true,
			}
		}
		totalOfferings += len(zoneData)
	}

	p.spotUpdateTime = time.Now()
	logging.FromContext(ctx).Infof("updated spot pricing with %d instance types and %d offerings", len(p.spotPrices), totalOfferings)
	return nil
}

func (p *PricingProvider) markODNotCurrent() {
	for k, v := range p.onDemandPrices {
		p.onDemandPrices[k] = priceData{
			lastKnownPrice: v.lastKnownPrice,
			isCurrent:      false,
		}
	}
}

func (p *PricingProvider) markSpotNotCurrent() {
	for it := range p.spotPrices {
		m := map[string]priceData{}
		for z, v := range p.spotPrices[it].prices {
			m[z] = priceData{
				lastKnownPrice: v.lastKnownPrice,
				isCurrent:      false,
			}
		}
		p.spotPrices[it] = zonalPricing{
			prices: m,
		}
	}
}

func populateInitialOnDemandPricing(pricing map[string]float64) map[string]priceData {
	m := map[string]priceData{}
	for k, v := range pricing {
		m[k] = priceData{
			lastKnownPrice: v,
			isCurrent:      true,
		}
	}
	return m
}

func populateInitialSpotPricing(pricing map[string]float64) map[string]zonalPricing {
	m := map[string]zonalPricing{}
	for it, price := range pricing {
		m[it] = newZonalPricing(ptr.Float64(price))
	}
	return m
}
