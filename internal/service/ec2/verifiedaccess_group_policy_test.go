package ec2_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccVerifiedAccessGroupPolicy_basic(t *testing.T) {
	var verifiedaccessgrouppolicy ec2.GetVerifiedAccessGroupPolicyOutput
	resourceName := "aws_verifiedaccess_group_policy.test"
	policyDocument := "permit(principal, action, resource) \nwhen {\n    context.http_request.method == \"GET\"\n};"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(ec2.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVerifiedAccessGroupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVerifiedAccessGroupPolicyConfig_basic(policyDocument),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVerifiedAccessGroupPolicyExists(resourceName, &verifiedaccessgrouppolicy),
					resource.TestCheckResourceAttr(resourceName, "policy_document", policyDocument),),
					resource.TestCheckResourceAttrSet(resourceName, "verified_access_group_id"),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"apply_immediately", "user"},
			},
		},
	})
}

func TestAccVerifiedAccessGroupPolicy_disappears(t *testing.T) {

	var verifiedaccessgrouppolicy ec2.GetVerifiedAccessGroupPolicyResponse
	
	resourceName := "aws_verifiedaccess_group_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(ec2.EndpointsID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVerifiedAccessGroupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVerifiedAccessGroupPolicyConfig_basic(rName, testAccVerifiedAccessGroupPolicyVersionNewer),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVerifiedAccessGroupPolicyExists(resourceName, &verifiedaccessgrouppolicy),
					acctest.CheckResourceDisappears(acctest.Provider, tfec2.ResourceVerifiedAccessGroupPolicy(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckVerifiedAccessGroupPolicyDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).Conn()
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_ec2_verifiedaccess_group_policy" {
			continue
		}

		input := &ec2.GetVerifiedAccessGroupPolicyInput{
			VerifiedAccessGroupId: aws.String(rs.Primary.ID),
		}
		_, err := conn.GetVerifiedAccessGroupPolicyWithContext(ctx, &ec2.DescribeVerifiedAccessGroupPolicyInput{
			VerifiedAccessGroupId: aws.String(rs.Primary.ID),
		})
		if err != nil {
			if tfawserr.ErrCodeEquals(err, ec2.ErrCodeNotFoundException) {
				return nil
			}
			return err
		}

		return create.Error(names.EC2, create.ErrActionCheckingDestroyed, tfec2.ResNameVerifiedAccessGroupPolicy, rs.Primary.ID, errors.New("not destroyed"))
	}

	return nil
}

func testAccCheckVerifiedAccessGroupPolicyExists(name string, verifiedaccessgrouppolicy *ec2.DescribeVerifiedAccessGroupPolicyResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.EC2, create.ErrActionCheckingExistence, tfec2.ResNameVerifiedAccessGroupPolicy, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.EC2, create.ErrActionCheckingExistence, tfec2.ResNameVerifiedAccessGroupPolicy, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).Conn()
		ctx := context.Background()
		resp, err := conn.DescribeVerifiedAccessGroupPolicyWithContext(ctx, &ec2.DescribeVerifiedAccessGroupPolicyInput{
			VerifiedAccessGroupPolicyId: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return create.Error(names., create.ErrActionCheckingExistence, tfec2.ResNameVerifiedAccessGroupPolicy, rs.Primary.ID, err)
		}

		*verifiedaccessgrouppolicy = *resp

		return nil
	}
}

func testAccCheckVerifiedAccessGroupPolicyNotRecreated(before, after *ec2.DescribeVerifiedAccessGroupPolicyResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.StringValue(before.VerifiedAccessGroupPolicyId), aws.StringValue(after.VerifiedAccessGroupPolicyId); before != after {
			return create.Error(names., create.ErrActionCheckingNotRecreated, tfec2.ResNameVerifiedAccessGroupPolicy, aws.StringValue(before.VerifiedAccessGroupPolicyId), errors.New("recreated"))
		}

		return nil
	}
}

func testAccVerifiedAccessGroupPolicyConfig_baseConfig() string {
	return fmt.Sprint(`
resource "aws_verifiedaccess_instance" "test" {}

resource "aws_verifiedaccess_trust_provider" "test" {
  policy_reference_name    = "test"
  trust_provider_type      = "user"
  user_trust_provider_type = "iam-identity-center"
}

resource "aws_verifiedaccess_trust_provider_attachment" "test" {
  verified_access_instance_id       = aws_verifiedaccess_instance.test.id
  verified_access_trust_provider_id = aws_verifiedaccess_trust_provider.test.id
}
resource "aws_verifiedaccess_group" "test" {
  verified_access_instance_id = aws_verifiedaccess_instance.test.id

  depends_on = [
    aws_verifiedaccess_trust_provider_attachment.test
  ]
}
`)
}

func testAccVerifiedAccessGroupPolicyConfig_basic(policy string) string {
	return acctest.ConfigCompose(testAccVerifiedAccessGroupConfig_baseConfig(policy), fmt.Sprintf(`
resource "aws_verifiedaccess_group_policy" "test" {
  policy_document          = %[1]q
  verified_access_group_id = aws_verifiedaccess_group.test.id
}
`, policy))
}
