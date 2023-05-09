package ec2

import (
	"context"
	"errors"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_verifiedaccess_group_policy", name="Verified Access Group Policy")
func ResourceVerifiedAccessGroupPolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVerifiedAccessGroupPolicyCreate,
		ReadWithoutTimeout:   resourceVerifiedAccessGroupPolicyRead,
		UpdateWithoutTimeout: resourceVerifiedAccessGroupPolicyUpdate,
		DeleteWithoutTimeout: resourceVerifiedAccessGroupPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"policy_document": {
				Type:                  schema.TypeString,
				Required:              true,
				DiffSuppressFunc:      suppressEquivalentGroupPolicyDiffs,
				DiffSuppressOnRefresh: true,
				ValidateFunc:          validation.StringIsNotEmpty,
			},
			"verified_access_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

const (
	ResNameVerifiedAccessGroupPolicy = "Verified Access Group Policy"
)

func resourceVerifiedAccessGroupPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).EC2Conn()

	in := &ec2.ModifyVerifiedAccessGroupPolicyInput{
		PolicyEnabled:         aws.Bool(true),
		VerifiedAccessGroupId: aws.String(d.Get("verified_access_group_id").(string)),
	}

	if v, ok := d.GetOk("policy_document"); ok {
		in.PolicyDocument = aws.String(v.(string))
	}

	out, err := conn.ModifyVerifiedAccessGroupPolicyWithContext(ctx, in)
	if err != nil {
		return create.DiagError(names.EC2, create.ErrActionCreating, ResNameVerifiedAccessGroupPolicy, d.Get("verified_access_group_id").(string), err)
	}

	if out == nil || out.PolicyEnabled == nil {
		return create.DiagError(names.EC2, create.ErrActionCreating, ResNameVerifiedAccessGroupPolicy, d.Get("verified_access_group_id").(string), errors.New("empty output"))
	}

	d.SetId(d.Get("verified_access_group_id").(string))

	return resourceVerifiedAccessGroupPolicyRead(ctx, d, meta)
}

func resourceVerifiedAccessGroupPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).EC2Conn()

	in := &ec2.GetVerifiedAccessGroupPolicyInput{
		VerifiedAccessGroupId: aws.String(d.Id()),
	}

	out, err := conn.GetVerifiedAccessGroupPolicyWithContext(ctx, in)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 VerifiedAccessGroupPolicy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.EC2, create.ErrActionReading, ResNameVerifiedAccessGroupPolicy, d.Id(), err)
	}

	d.Set("policy_document", out.PolicyDocument)
	d.Set("verified_access_group_id", d.Id())

	return nil
}

func resourceVerifiedAccessGroupPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).EC2Conn()

	update := false

	in := &ec2.ModifyVerifiedAccessGroupPolicyInput{
		PolicyEnabled:         aws.Bool(true),
		VerifiedAccessGroupId: aws.String(d.Get("verified_access_group_id").(string)),
	}

	if d.HasChanges("policy_document") {
		in.PolicyDocument = aws.String(d.Get("policy_document").(string))
		update = true
	}

	if !update {
		return nil
	}

	log.Printf("[DEBUG] Updating EC2 VerifiedAccessGroupPolicy (%s): %#v", d.Id(), in)

	_, err := conn.ModifyVerifiedAccessGroupPolicyWithContext(ctx, in)

	if err != nil {
		return create.DiagError(names.EC2, create.ErrActionUpdating, ResNameVerifiedAccessGroupPolicy, d.Id(), err)
	}

	return resourceVerifiedAccessGroupPolicyRead(ctx, d, meta)
}

func resourceVerifiedAccessGroupPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).EC2Conn()

	log.Printf("[INFO] Deleting EC2 VerifiedAccessGroupPolicy %s", d.Id())

	in := &ec2.ModifyVerifiedAccessGroupPolicyInput{
		PolicyEnabled:         aws.Bool(false),
		VerifiedAccessGroupId: aws.String(d.Get("verified_access_group_id").(string)),
	}

	_, err := conn.ModifyVerifiedAccessGroupPolicyWithContext(ctx, in)

	if err != nil {
		return create.DiagError(names.EC2, create.ErrActionDeleting, ResNameVerifiedAccessGroupPolicy, d.Id(), err)
	}

	return nil
}

func suppressEquivalentGroupPolicyDiffs(k, old, new string, d *schema.ResourceData) bool {
	// ignore leading and trailing whitespace for policies
	old_policy := strings.TrimSpace(old)
	new_policy := strings.TrimSpace(new)

	return reflect.DeepEqual(old_policy, new_policy)
}
