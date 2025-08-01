// Code generated by internal/generate/servicepackage/main.go; DO NOT EDIT.

package ce

import (
	"context"
	"unique"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	inttypes "github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/internal/vcr"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*inttypes.ServicePackageFrameworkDataSource {
	return []*inttypes.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*inttypes.ServicePackageFrameworkResource {
	return []*inttypes.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*inttypes.ServicePackageSDKDataSource {
	return []*inttypes.ServicePackageSDKDataSource{
		{
			Factory:  dataSourceCostCategory,
			TypeName: "aws_ce_cost_category",
			Name:     "Cost Category",
			Region:   unique.Make(inttypes.ResourceRegionDisabled()),
		},
		{
			Factory:  dataSourceTags,
			TypeName: "aws_ce_tags",
			Name:     "Tags",
			Region:   unique.Make(inttypes.ResourceRegionDisabled()),
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*inttypes.ServicePackageSDKResource {
	return []*inttypes.ServicePackageSDKResource{
		{
			Factory:  resourceAnomalyMonitor,
			TypeName: "aws_ce_anomaly_monitor",
			Name:     "Anomaly Monitor",
			Tags: unique.Make(inttypes.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			}),
			Region: unique.Make(inttypes.ResourceRegionDisabled()),
			Identity: inttypes.GlobalARNIdentity(
				inttypes.WithIdentityDuplicateAttrs(names.AttrID),
				inttypes.WithV6_0SDKv2Fix(),
			),
			Import: inttypes.SDKv2Import{
				WrappedImport: true,
			},
		},
		{
			Factory:  resourceAnomalySubscription,
			TypeName: "aws_ce_anomaly_subscription",
			Name:     "Anomaly Subscription",
			Tags: unique.Make(inttypes.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			}),
			Region: unique.Make(inttypes.ResourceRegionDisabled()),
			Identity: inttypes.GlobalARNIdentity(
				inttypes.WithIdentityDuplicateAttrs(names.AttrID),
				inttypes.WithV6_0SDKv2Fix(),
			),
			Import: inttypes.SDKv2Import{
				WrappedImport: true,
			},
		},
		{
			Factory:  resourceCostAllocationTag,
			TypeName: "aws_ce_cost_allocation_tag",
			Name:     "Cost Allocation Tag",
			Region:   unique.Make(inttypes.ResourceRegionDisabled()),
		},
		{
			Factory:  resourceCostCategory,
			TypeName: "aws_ce_cost_category",
			Name:     "Cost Category",
			Tags: unique.Make(inttypes.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			}),
			Region: unique.Make(inttypes.ResourceRegionDisabled()),
			Identity: inttypes.GlobalARNIdentity(
				inttypes.WithIdentityDuplicateAttrs(names.AttrID),
				inttypes.WithV6_0SDKv2Fix(),
			),
			Import: inttypes.SDKv2Import{
				WrappedImport: true,
			},
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.CE
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*costexplorer.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws.Config))
	optFns := []func(*costexplorer.Options){
		costexplorer.WithEndpointResolverV2(newEndpointResolverV2()),
		withBaseEndpoint(config[names.AttrEndpoint].(string)),
		func(o *costexplorer.Options) {
			if region := config[names.AttrRegion].(string); o.Region != region {
				tflog.Info(ctx, "overriding provider-configured AWS API region", map[string]any{
					"service":         p.ServicePackageName(),
					"original_region": o.Region,
					"override_region": region,
				})
				o.Region = region
			}
		},
		func(o *costexplorer.Options) {
			if inContext, ok := conns.FromContext(ctx); ok && inContext.VCREnabled() {
				tflog.Info(ctx, "overriding retry behavior to immediately return VCR errors")
				o.Retryer = conns.AddIsErrorRetryables(cfg.Retryer().(aws.RetryerV2), retry.IsErrorRetryableFunc(vcr.InteractionNotFoundRetryableFunc))
			}
		},
		withExtraOptions(ctx, p, config),
	}

	return costexplorer.NewFromConfig(cfg, optFns...), nil
}

// withExtraOptions returns a functional option that allows this service package to specify extra API client options.
// This option is always called after any generated options.
func withExtraOptions(ctx context.Context, sp conns.ServicePackage, config map[string]any) func(*costexplorer.Options) {
	if v, ok := sp.(interface {
		withExtraOptions(context.Context, map[string]any) []func(*costexplorer.Options)
	}); ok {
		optFns := v.withExtraOptions(ctx, config)

		return func(o *costexplorer.Options) {
			for _, optFn := range optFns {
				optFn(o)
			}
		}
	}

	return func(*costexplorer.Options) {}
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
