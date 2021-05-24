package aws

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
	tfservicecatalog "github.com/terraform-providers/terraform-provider-aws/aws/internal/service/servicecatalog"
)

func dataSourceAwsServiceCatalogPortfolio() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAwsServiceCatalogPortfolioRead,

		Schema: map[string]*schema.Schema{
			"accept_language": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "en",
				ValidateFunc: validation.StringInSlice(tfservicecatalog.AcceptLanguage_Values(), false),
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func dataSourceAwsServiceCatalogPortfolioRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).scconn

	input := &servicecatalog.DescribePortfolioInput{
		Id: aws.String(d.Get("id").(string)),
	}

	if v, ok := d.GetOk("accept_language"); ok {
		input.AcceptLanguage = aws.String(v.(string))
	}

	output, err := conn.DescribePortfolio(input)

	if err != nil {
		return fmt.Errorf("error getting Service Catalog Portfolio: %w", err)
	}

	if output == nil || output.PortfolioDetail == nil {
		return fmt.Errorf("error getting Service Catalog Portfolio: empty response")
	}

	detail := output.PortfolioDetail

	d.SetId(aws.StringValue(detail.Id))

	if err := d.Set("created_time", detail.CreatedTime.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Error setting created_time: %s", err)
	}

	d.Set("arn", detail.ARN)
	d.Set("description", detail.Description)
	d.Set("name", detail.DisplayName)
	d.Set("provider_name", detail.ProviderName)

	ignoreTagsConfig := meta.(*AWSClient).IgnoreTagsConfig
	tags := keyvaluetags.ServicecatalogKeyValueTags(output.Tags)

	if err := d.Set("tags", tags.IgnoreAws().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	return nil
}
