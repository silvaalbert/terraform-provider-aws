package ec2_test

// import (
// 	"context"
// 	"fmt"
// 	"regexp"
// 	"strings"
// 	"testing"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/service/ec2"
// 	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
// 	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
// 	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
// 	"github.com/hashicorp/terraform-provider-aws/internal/conns"
// 	"github.com/hashicorp/terraform-provider-aws/internal/create"
// 	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
// )

// func TestVerifiedAccessEndpointExampleUnitTest(t *testing.T) {
// 	testCases := []struct {
// 		TestName string
// 		Input    string
// 		Expected string
// 		Error    bool
// 	}{
// 		{
// 			TestName: "empty",
// 			Input:    "",
// 			Expected: "",
// 			Error:    true,
// 		},
// 		{
// 			TestName: "descriptive name",
// 			Input:    "some input",
// 			Expected: "some output",
// 			Error:    false,
// 		},
// 		{
// 			TestName: "another descriptive name",
// 			Input:    "more input",
// 			Expected: "more output",
// 			Error:    false,
// 		},
// 	}

// 	for _, testCase := range testCases {
// 		t.Run(testCase.TestName, func(t *testing.T) {
// 			got, err := tfec2.FunctionFromResource(testCase.Input)

// 			if err != nil && !testCase.Error {
// 				t.Errorf("got error (%s), expected no error", err)
// 			}

// 			if err == nil && testCase.Error {
// 				t.Errorf("got (%s) and no error, expected error", got)
// 			}

// 			if got != testCase.Expected {
// 				t.Errorf("got %s, expected %s", got, testCase.Expected)
// 			}
// 		})
// 	}
// }

// func TestAccEC2VerifiedAccessEndpoint_basic(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping long-running test in short mode")
// 	}

// 	var verifiedaccessendpoint ec2.DescribeVerifiedAccessEndpointResponse
// 	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
// 	resourceName := "aws_ec2_verifiedaccess_endpoint.test"

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck: func() {
// 			acctest.PreCheck(t)
// 			acctest.PreCheckPartitionHasService(ec2.EndpointsID, t)
// 			testAccPreCheck(t)
// 		},
// 		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
// 		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
// 		CheckDestroy:             testAccCheckVerifiedAccessEndpointDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccVerifiedAccessEndpointConfig_basic(rName),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckVerifiedAccessEndpointExists(resourceName, &verifiedaccessendpoint),
// 					resource.TestCheckResourceAttr(resourceName, "auto_minor_version_upgrade", "false"),
// 					resource.TestCheckResourceAttrSet(resourceName, "maintenance_window_start_time.0.day_of_week"),
// 					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "user.*", map[string]string{
// 						"console_access": "false",
// 						"groups.#":       "0",
// 						"username":       "Test",
// 						"password":       "TestTest1234",
// 					}),
// 					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "ec2", regexp.MustCompile(`verifiedaccessendpoint:+.`)),
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

// func TestAccEC2VerifiedAccessEndpoint_disappears(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skipping long-running test in short mode")
// 	}

// 	var verifiedaccessendpoint ec2.DescribeVerifiedAccessEndpointResponse
// 	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
// 	resourceName := "aws_ec2_verifiedaccess_endpoint.test"

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck: func() {
// 			acctest.PreCheck(t)
// 			acctest.PreCheckPartitionHasService(ec2.EndpointsID, t)
// 			testAccPreCheck(t)
// 		},
// 		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
// 		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
// 		CheckDestroy:             testAccCheckVerifiedAccessEndpointDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccVerifiedAccessEndpointConfig_basic(rName, testAccVerifiedAccessEndpointVersionNewer),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckVerifiedAccessEndpointExists(resourceName, &verifiedaccessendpoint),
// 					acctest.CheckResourceDisappears(acctest.Provider, tfec2.ResourceVerifiedAccessEndpoint(), resourceName),
// 				),
// 				ExpectNonEmptyPlan: true,
// 			},
// 		},
// 	})
// }

// func testAccCheckVerifiedAccessEndpointDestroy(s *terraform.State) error {
// 	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn()
// 	ctx := context.Background()

// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "aws_ec2_verifiedaccess_endpoint" {
// 			continue
// 		}

// 		input := &ec2.DescribeVerifiedAccessEndpointInput{
// 			VerifiedAccessEndpointId: aws.String(rs.Primary.ID),
// 		}
// 		_, err := conn.DescribeVerifiedAccessEndpointWithContext(ctx, &ec2.DescribeVerifiedAccessEndpointInput{
// 			VerifiedAccessEndpointId: aws.String(rs.Primary.ID),
// 		})
// 		if err != nil {
// 			if tfawserr.ErrCodeEquals(err, ec2.ErrCodeNotFoundException) {
// 				return nil
// 			}
// 			return err
// 		}

// 		return create.Error(names.EC2, create.ErrActionCheckingDestroyed, tfec2.ResNameVerifiedAccessEndpoint, rs.Primary.ID, errors.New("not destroyed"))
// 	}

// 	return nil
// }

// func testAccCheckVerifiedAccessEndpointExists(name string, verifiedaccessendpoint *ec2.DescribeVerifiedAccessEndpointResponse) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[name]
// 		if !ok {
// 			return create.Error(names.EC2, create.ErrActionCheckingExistence, tfec2.ResNameVerifiedAccessEndpoint, name, errors.New("not found"))
// 		}

// 		if rs.Primary.ID == "" {
// 			return create.Error(names.EC2, create.ErrActionCheckingExistence, tfec2.ResNameVerifiedAccessEndpoint, name, errors.New("not set"))
// 		}

// 		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn()
// 		ctx := context.Background()
// 		resp, err := conn.DescribeVerifiedAccessEndpointWithContext(ctx, &ec2.DescribeVerifiedAccessEndpointInput{
// 			VerifiedAccessEndpointId: aws.String(rs.Primary.ID),
// 		})

// 		if err != nil {
// 			return create.Error(names.EC2, create.ErrActionCheckingExistence, tfec2.ResNameVerifiedAccessEndpoint, rs.Primary.ID, err)
// 		}

// 		*verifiedaccessendpoint = *resp

// 		return nil
// 	}
// }

// func testAccPreCheck(t *testing.T) {
// 	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn()
// 	ctx := context.Background()

// 	input := &ec2.ListVerifiedAccessEndpointsInput{}
// 	_, err := conn.ListVerifiedAccessEndpointsWithContext(ctx, input)

// 	if acctest.PreCheckSkipError(err) {
// 		t.Skipf("skipping acceptance testing: %s", err)
// 	}

// 	if err != nil {
// 		t.Fatalf("unexpected PreCheck error: %s", err)
// 	}
// }

// func testAccCheckVerifiedAccessEndpointNotRecreated(before, after *ec2.DescribeVerifiedAccessEndpointResponse) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		if before, after := aws.StringValue(before.VerifiedAccessEndpointId), aws.StringValue(after.VerifiedAccessEndpointId); before != after {
// 			return create.Error(names.EC2, create.ErrActionCheckingNotRecreated, tfec2.ResNameVerifiedAccessEndpoint, aws.StringValue(before.VerifiedAccessEndpointId), errors.New("recreated"))
// 		}

// 		return nil
// 	}
// }

// func testAccVerifiedAccessEndpointConfig_basic(rName, version string) string {
// 	return fmt.Sprintf(`
// resource "aws_security_group" "test" {
//   name = %[1]q
// }

// resource "aws_ec2_verifiedaccess_endpoint" "test" {
//   verifiedaccess_endpoint_name             = %[1]q
//   engine_type             = "ActiveEC2"
//   engine_version          = %[2]q
//   host_instance_type      = "ec2.t2.micro"
//   security_groups         = [aws_security_group.test.id]
//   authentication_strategy = "simple"
//   storage_type            = "efs"

//   logs {
//     general = true
//   }

//   user {
//     username = "Test"
//     password = "TestTest1234"
//   }
// }
// `, rName, version)
// }
