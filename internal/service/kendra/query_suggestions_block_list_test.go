package kendra_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kendra"
	"github.com/aws/aws-sdk-go-v2/service/kendra/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfkendra "github.com/hashicorp/terraform-provider-aws/internal/service/kendra"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestQuerySuggestionsBlockListParseID(t *testing.T) {
	testCases := []struct {
		TestName        string
		Input           string
		ExpectedId      string
		ExpectedIndexId string
		Error           bool
	}{
		{
			TestName:        "empty",
			Input:           "",
			ExpectedId:      "",
			ExpectedIndexId: "",
			Error:           true,
		},
		{
			TestName:        "Invalid ID",
			Input:           "abcdefg12345678/",
			ExpectedId:      "",
			ExpectedIndexId: "",
			Error:           true,
		},
		{
			TestName:        "Invalid ID separator",
			Input:           "abcdefg12345678:qwerty09876",
			ExpectedId:      "",
			ExpectedIndexId: "",
			Error:           true,
		},
		{
			TestName:        "Invalid ID with more than 1 separator",
			Input:           "abcdefg12345678/qwerty09876/zxcvbnm123456",
			ExpectedId:      "",
			ExpectedIndexId: "",
			Error:           true,
		},
		{
			TestName:        "Valid ID",
			Input:           "abcdefg12345678/qwerty09876",
			ExpectedId:      "abcdefg12345678",
			ExpectedIndexId: "qwerty09876",
			Error:           false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.TestName, func(t *testing.T) {
			gotId, gotIndexId, err := tfkendra.QuerySuggestionsBlockListParseResourceID(testCase.Input)

			if err != nil && !testCase.Error {
				t.Errorf("got error (%s), expected no error", err)
			}

			if err == nil && testCase.Error {
				t.Errorf("got (Id: %s, IndexId: %s) and no error, expected error", gotId, gotIndexId)
			}

			if gotId != testCase.ExpectedId {
				t.Errorf("got %s, expected %s", gotId, testCase.ExpectedIndexId)
			}

			if gotIndexId != testCase.ExpectedIndexId {
				t.Errorf("got %s, expected %s", gotIndexId, testCase.ExpectedIndexId)
			}
		})
	}
}

func TestAccKendraQuerySuggestionsBlockList_basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kendra_querysuggestionsblocklist.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.KendraEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:        acctest.ErrorCheck(t, names.KendraEndpointID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckQuerySuggestionsBlockListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccQuerySuggestionsBlockListConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQuerySuggestionsBlockListExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "index_id", "aws_kendra_index.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "role_arn", "aws_iam_role.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "source_s3_path.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "source_s3_path.0.bucket", "aws_s3_bucket.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "source_s3_path.0.key", "test/suggestions.txt"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "kendra", regexp.MustCompile(`querysuggestionsblocklist:+.`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKendraQuerySuggestionsBlockList_disappears(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kendra_querysuggestionsblocklist.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.KendraEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:        acctest.ErrorCheck(t, names.KendraEndpointID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckQuerySuggestionsBlockListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccQuerySuggestionsBlockListConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQuerySuggestionsBlockListExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfkendra.ResourceQuerySuggestionsBlockList(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccKendraQuerySuggestionsBlockList_tags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kendra_querysuggestionsblocklist.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.KendraEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:        acctest.ErrorCheck(t, names.KendraEndpointID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckQuerySuggestionsBlockListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccQuerySuggestionsBlockListConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQuerySuggestionsBlockListExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccQuerySuggestionsBlockListConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQuerySuggestionsBlockListExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccQuerySuggestionsBlockListConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckQuerySuggestionsBlockListExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckQuerySuggestionsBlockListDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).KendraConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_kendra_querysuggestionsblocklist" {
			continue
		}

		id, indexId, err := tfkendra.QuerySuggestionsBlockListParseResourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		input := &kendra.DescribeQuerySuggestionsBlockListInput{
			Id:      aws.String(id),
			IndexId: aws.String(indexId),
		}

		_, err = conn.DescribeQuerySuggestionsBlockList(context.TODO(), input)
		if err != nil {
			var notFound *types.ResourceNotFoundException

			if errors.As(err, &notFound) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Expected Kendra QuerySuggestionsBlockList to be destroyed, %s found", rs.Primary.ID)
	}

	return nil
}

func testAccCheckQuerySuggestionsBlockListExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Kendra QuerySuggestionsBlockList is set")
		}

		id, indexId, err := tfkendra.QuerySuggestionsBlockListParseResourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).KendraConn

		input := &kendra.DescribeQuerySuggestionsBlockListInput{
			Id:      aws.String(id),
			IndexId: aws.String(indexId),
		}

		_, err = conn.DescribeQuerySuggestionsBlockList(context.TODO(), input)

		if err != nil {
			return fmt.Errorf("Error describing Kendra QuerySuggestionsBlockList: %s", err.Error())
		}

		return nil
	}
}

func testAccPreCheck(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).KendraConn

	input := &kendra.ListQuerySuggestionsBlockListsInput{}

	_, err := conn.ListQuerySuggestionsBlockLists(context.TODO(), input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccQuerySuggestionsBlockListBaseConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "kendra.${data.aws_partition.current.dns_suffix}"
      }
    }]
  })

  name = %[1]q

  inline_policy {
    name = "test"
    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [{
        Action = [
          "kendra:*",
          "s3:GetBucketLocation",
          "s3:GetObject",
          "s3:ListBucket",
        ]
        Effect = "Allow"
        Resource = [
          aws_s3_bucket.test.arn,
          "${aws_s3_bucket.test.arn}/*"
        ]
      }]
    })
  }
}

resource "aws_kendra_index" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn
}
`, rName)
}
func testAccQuerySuggestionsBlockListConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccQuerySuggestionsBlockListBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_kendra_query_suggestions_block_list" "test" {
  index_id = aws_kendra_index.test.id
  name     = %[1]q
  role_arn = aws_iam_role.test.arn

  source_s3_path {
    bucket = aws_s3_bucket.test.id
    key    = "test/suggestions.txt"
  }
}
`, rName))
}

func testAccQuerySuggestionsBlockListConfig_tags1(rName, tag, value string) string {
	return acctest.ConfigCompose(
		testAccQuerySuggestionsBlockListBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_kendra_query_suggestions_block_list" "test" {
  index_id = aws_kendra_index.test.id
  name     = %[1]q
  role_arn = aws_iam_role.test.arn

  source_s3_path {
    bucket = aws_s3_bucket.test.id
    key    = "test/suggestions.txt"
  }

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tag, value))
}

func testAccQuerySuggestionsBlockListConfig_tags2(rName, tag1, value1, tag2, value2 string) string {
	return acctest.ConfigCompose(
		testAccQuerySuggestionsBlockListBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_kendra_query_suggestions_block_list" "test" {
  index_id = aws_kendra_index.test.id
  name     = %[1]q
  role_arn = aws_iam_role.test.arn

  source_s3_path {
    bucket = aws_s3_bucket.test.id
    key    = "test/suggestions.txt"
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tag1, value1, tag2, value2))
}
