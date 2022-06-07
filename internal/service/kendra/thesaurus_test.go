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

func TestThesaurusParseResourceID(t *testing.T) {
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
			gotId, gotIndexId, err := tfkendra.ThesaurusParseResourceID(testCase.Input)

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

func TestAccKendraThesaurus_basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kendra_thesaurus.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.KendraEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:        acctest.ErrorCheck(t, names.KendraEndpointID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckThesaurusDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccThesaurusConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThesaurusExists(resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "index_id", "aws_kendra_index.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "role_arn", "aws_iam_role.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "source_s3_path.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "source_s3_path.0.bucket", "aws_s3_bucket.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "source_s3_path.0.key", "test/suggestions.txt"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "kendra", regexp.MustCompile(`thesaurus:+.`)),
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

func TestAccKendraThesaurus_disappears(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kendra_thesaurus.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.KendraEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:        acctest.ErrorCheck(t, names.KendraEndpointID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckThesaurusDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccThesaurusConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThesaurusExists(resourceName),
					acctest.CheckResourceDisappears(acctest.Provider, tfkendra.ResourceThesaurus(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccKendraThesaurus_tags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kendra_thesaurus.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.KendraEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:        acctest.ErrorCheck(t, names.KendraEndpointID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckThesaurusDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccThesaurusConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThesaurusExists(resourceName),
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
				Config: testAccThesaurusConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThesaurusExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccThesaurusConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckThesaurusExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckThesaurusDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).KendraConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_kendra_thesaurus" {
			continue
		}

		id, indexId, err := tfkendra.ThesaurusParseResourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		input := &kendra.DescribeThesaurusInput{
			Id:      aws.String(id),
			IndexId: aws.String(indexId),
		}

		_, err = conn.DescribeThesaurus(context.TODO(), input)
		if err != nil {
			var notFound *types.ResourceNotFoundException
			if errors.As(err, &notFound) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Expected Kendra Thesaurus to be destroyed, %s found", rs.Primary.ID)
	}

	return nil
}

func testAccCheckThesaurusExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Kendra Thesaurus is set")
		}

		id, indexId, err := tfkendra.ThesaurusParseResourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).KendraConn
		_, err = conn.DescribeThesaurus(context.TODO(), &kendra.DescribeThesaurusInput{
			Id:      aws.String(id),
			IndexId: aws.String(indexId),
		})

		if err != nil {
			return fmt.Errorf("Error describing Kendra Thesaurus: %s", err.Error())
		}

		return nil
	}
}

func testAccThesaurusBaseConfig(rName string) string {
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

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
  force_destroy = true
}

resource "aws_s3_object" "test" {
  bucket  = aws_s3_bucket.test.bucket
  content = "test"
  key     = "test/thesaurus.txt"
}
`, rName)
}

func testAccThesaurusConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccThesaurusBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_kendra_thesaurus" "test" {
  index_id = aws_kendra_index.id
  name = %[1]q
  role_arn = aws_iam_role.test_arn

  source_s3_path {
   bucket = aws_s3_bucket.test.id
   key = aws_s3_object.test.key
  }
}
`, rName))
}

func testAccThesaurusConfig_tags1(rName, tag, value string) string {
	return acctest.ConfigCompose(
		testAccThesaurusBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_kendra_thesaurus" "test" {
  index_id = aws_kendra_index.id
  name = %[1]q
  role_arn = aws_iam_role.test_arn

  source_s3_path {
   bucket = aws_s3_bucket.test.id
   key = aws_s3_object.test.key
  }

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tag, value))
}

func testAccThesaurusConfig_tags2(rName, tag1, value1, tag2, value2 string) string {
	return acctest.ConfigCompose(
		testAccThesaurusBaseConfig(rName),
		fmt.Sprintf(`
resource "aws_kendra_thesaurus" "test" {
  index_id = aws_kendra_index.id
  name = %[1]q
  role_arn = aws_iam_role.test_arn

  source_s3_path {
   bucket = aws_s3_bucket.test.id
   key = aws_s3_object.test.key
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tag1, value1, tag2, value2))
}
