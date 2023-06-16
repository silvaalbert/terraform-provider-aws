// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package rum

import (
	"context"

	aws_sdkv1 "github.com/aws/aws-sdk-go/aws"
	session_sdkv1 "github.com/aws/aws-sdk-go/aws/session"
	cloudwatchrum_sdkv1 "github.com/aws/aws-sdk-go/service/cloudwatchrum"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct {
	config map[string]any
}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceAppMonitor,
			TypeName: "aws_rum_app_monitor",
			Name:     "App Monitor",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceMetricsDestination,
			TypeName: "aws_rum_metrics_destination",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.RUM
}

func (p *servicePackage) Configure(ctx context.Context, config map[string]any) {
	p.config = config
}

// NewConn returns a new AWS SDK for Go v1 client for this service package's AWS API.
func (p *servicePackage) NewConn(ctx context.Context) (*cloudwatchrum_sdkv1.CloudWatchRUM, error) {
	sess := p.config["session"].(*session_sdkv1.Session)

	return cloudwatchrum_sdkv1.New(sess.Copy(&aws_sdkv1.Config{Endpoint: aws_sdkv1.String(p.config["endpoint"].(string))})), nil
}

var ServicePackage = &servicePackage{}
