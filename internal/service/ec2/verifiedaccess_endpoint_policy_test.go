package ec2_test

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"testing"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/service/ec2"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
// 	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
// 	"github.com/hashicorp/terraform-provider-aws/internal/conns"
// 	"github.com/hashicorp/terraform-provider-aws/internal/create"
// 	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
// 	"github.com/hashicorp/terraform-provider-aws/names"
// )

// func TestAccVerifiedAccessEndpointPolicy_basic(t *testing.T) {
// 	ctx := context.Background()
// 	var verifiedaccessgrouppolicy ec2.GetVerifiedAccessEndpointPolicyOutput
// 	resourceName := "aws_verifiedaccess_endpoint_policy.test"
// 	policyDocument := "permit(principal, action, resource) \nwhen {\n    context.http_request.method == \"GET\"\n};"
// 	key := acctest.TLSRSAPrivateKeyPEM(t, 2048)
// 	certificate := acctest.TLSRSAX509SelfSignedCertificatePEM(t, key, "test.com")

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck: func() {
// 			acctest.PreCheck(ctx, t)
// 			testAccPreCheck(t)
// 		},
// 		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
// 		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
// 		CheckDestroy:             testAccCheckVerifiedAccessEndpointPolicyDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccVerifiedAccessEndpointPolicyConfig_basic(certificate, key, policyDocument),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckVerifiedAccessEndpointPolicyExists(resourceName, &verifiedaccessgrouppolicy),
// 					resource.TestCheckResourceAttr(resourceName, "policy_document", policyDocument),
// 					resource.TestCheckResourceAttrSet(resourceName, "verified_access_endpoint_id"),
// 				),
// 			},
// 			{
// 				ResourceName:            resourceName,
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"apply_immediately", "user"},
// 			},
// 		},
// 	})
// }

// func TestAccVerifiedAccessEndpointPolicy_disappears(t *testing.T) {
// 	ctx := context.Background()
// 	var verifiedaccessgrouppolicy ec2.GetVerifiedAccessEndpointPolicyOutput

// 	resourceName := "aws_verifiedaccess_endpoint_policy.test"
// 	policyDocument := "permit(principal, action, resource) \nwhen {\n    context.http_request.method == \"GET\"\n};"
// 	key := acctest.TLSRSAPrivateKeyPEM(t, 2048)
// 	certificate := acctest.TLSRSAX509SelfSignedCertificatePEM(t, key, "test.com")

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck: func() {
// 			acctest.PreCheck(ctx, t)
// 			testAccPreCheck(t)
// 		},
// 		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
// 		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
// 		CheckDestroy:             testAccCheckVerifiedAccessEndpointPolicyDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccVerifiedAccessEndpointPolicyConfig_basic(certificate, key, policyDocument),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckVerifiedAccessEndpointPolicyExists(resourceName, &verifiedaccessgrouppolicy),
// 					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfec2.ResourceVerifiedAccessEndpointPolicy(), resourceName),
// 				),
// 				ExpectNonEmptyPlan: true,
// 			},
// 		},
// 	})
// }

// func testAccCheckVerifiedAccessEndpointPolicyDestroy(s *terraform.State) error {
// 	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn()
// 	ctx := context.Background()

// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "aws_ec2_verifiedaccess_endpoint_policy" {
// 			continue
// 		}

// 		input := &ec2.GetVerifiedAccessEndpointPolicyInput{
// 			VerifiedAccessEndpointId: aws.String(rs.Primary.ID),
// 		}
// 		_, err := conn.GetVerifiedAccessEndpointPolicyWithContext(ctx, input)

// 		if err != nil {
// 			return err
// 		}

// 		return create.Error(names.EC2, create.ErrActionCheckingDestroyed, tfec2.ResNameVerifiedAccessEndpointPolicy, rs.Primary.ID, errors.New("not destroyed"))
// 	}

// 	return nil
// }

// func testAccCheckVerifiedAccessEndpointPolicyExists(name string, verifiedaccessgrouppolicy *ec2.GetVerifiedAccessEndpointPolicyOutput) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[name]
// 		if !ok {
// 			return create.Error(names.EC2, create.ErrActionCheckingExistence, tfec2.ResNameVerifiedAccessEndpointPolicy, name, errors.New("not found"))
// 		}

// 		if rs.Primary.ID == "" {
// 			return create.Error(names.EC2, create.ErrActionCheckingExistence, tfec2.ResNameVerifiedAccessEndpointPolicy, name, errors.New("not set"))
// 		}

// 		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn()
// 		ctx := context.Background()
// 		resp, err := conn.GetVerifiedAccessEndpointPolicyWithContext(ctx, &ec2.GetVerifiedAccessEndpointPolicyInput{
// 			VerifiedAccessEndpointId: aws.String(rs.Primary.ID),
// 		})

// 		if err != nil {
// 			return create.Error(names.EC2, create.ErrActionCheckingExistence, tfec2.ResNameVerifiedAccessGroup, rs.Primary.ID, err)
// 		}

// 		*verifiedaccessgrouppolicy = *resp
// 		return nil
// 	}
// }

// func testAccVerifiedAccessEndpointPolicyConfig_baseConfig(certificate, key string) string {

// 	return fmt.Sprintf(`
// resource "aws_verifiedaccess_instance" "test" {}

// resource "aws_verifiedaccess_trust_provider" "test" {
//   policy_reference_name    = "test"
//   trust_provider_type      = "user"
//   user_trust_provider_type = "iam-identity-center"
// }

// resource "aws_verifiedaccess_trust_provider_attachment" "test" {
//   verified_access_instance_id       = aws_verifiedaccess_instance.test.id
//   verified_access_trust_provider_id = aws_verifiedaccess_trust_provider.test.id
// }

// resource "aws_verifiedaccess_group" "test" {
//   verified_access_instance_id = aws_verifiedaccess_instance.test.id

//   depends_on = [
//     aws_verifiedaccess_trust_provider_attachment.test
//   ]
// }

// resource "aws_verifiedaccess_endpoint" "test" {
// 	application_domain     = "test.com"
// 	attachment_type        = "vpc"
// 	description            = "test"
// 	domain_certificate_arn = aws_acm_certificate.test.arn
// 	endpoint_domain_prefix = "test"
// 	endpoint_type          = "load-balancer"
// 	load_balancer_options {
// 	  load_balancer_arn = aws_lb.test.arn
// 	  port              = 443
// 	  protocol          = "https"
// 	  subnet_ids        = [aws_subnet.test_01.id, aws_subnet.test_02.id]
// 	}
// 	security_group_ids       = [aws_security_group.test.id]
// 	verified_access_group_id = aws_verifiedaccess_group.test.id
//   }

//   data "aws_availability_zones" "available" {
// 	state = "available"

// 	filter {
// 	  name   = "opt-in-status"
// 	  values = ["opt-in-not-required"]
// 	}
//   }

//   resource "aws_vpc" "test" {
// 	assign_generated_ipv6_cidr_block = true
// 	cidr_block                       = "10.0.0.0/16"
//   }

//   resource "aws_subnet" "test_01" {
// 	availability_zone = data.aws_availability_zones.available.names[0]
// 	cidr_block        = "10.0.10.0/24"
// 	vpc_id            = aws_vpc.test.id
//   }

//   resource "aws_subnet" "test_02" {
// 	availability_zone = data.aws_availability_zones.available.names[1]
// 	cidr_block        = "10.0.20.0/24"
// 	vpc_id            = aws_vpc.test.id
//   }

//   resource "aws_security_group" "test" {
// 	vpc_id = aws_vpc.test.id
//   }

//   resource "aws_lb" "test" {
// 	internal                   = true
// 	load_balancer_type         = "application"
// 	security_groups            = [aws_security_group.test.id]
// 	subnets                    = [aws_subnet.test_01.id, aws_subnet.test_02.id]
// 	enable_deletion_protection = false
//   }

//   resource "aws_lb_listener" "test" {
// 	load_balancer_arn = aws_lb.test.arn
// 	port              = "443"
// 	protocol          = "HTTPS"
// 	ssl_policy        = "ELBSecurityPolicy-2016-08"
// 	certificate_arn   = aws_acm_certificate.test.arn

// 	default_action {
// 	  type             = "forward"
// 	  target_group_arn = aws_lb_target_group.test.arn
// 	}
//   }

//   resource "aws_lb_target_group" "test" {
// 	name        = "test"
// 	target_type = "lambda"
//   }

//   resource "aws_acm_certificate" "test" {
// 	certificate_body = "%[1]s"
// 	private_key      = "%[2]s"
//   }
// `, acctest.TLSPEMEscapeNewlines(certificate), acctest.TLSPEMEscapeNewlines(key))
// }

// func testAccVerifiedAccessEndpointPolicyConfig_basic(certificate, key, policy string) string {
// 	return acctest.ConfigCompose(testAccVerifiedAccessEndpointPolicyConfig_baseConfig(certificate, key), fmt.Sprintf(`
// resource "aws_verifiedaccess_endpoint_policy" "test" {
//   policy_document          = %[1]q
//   verified_access_endpoint_id = aws_verifiedaccess_endpoint.test.id
// }
// `, policy))
// }
