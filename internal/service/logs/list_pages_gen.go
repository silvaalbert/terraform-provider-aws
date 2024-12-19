// Code generated by "internal/generate/listpages/main.go -ListOps=DescribeIndexPolicies,DescribeQueryDefinitions,DescribeResourcePolicies"; DO NOT EDIT.

package logs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func describeIndexPoliciesPages(ctx context.Context, conn *cloudwatchlogs.Client, input *cloudwatchlogs.DescribeIndexPoliciesInput, fn func(*cloudwatchlogs.DescribeIndexPoliciesOutput, bool) bool) error {
	for {
		output, err := conn.DescribeIndexPolicies(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.ToString(output.NextToken) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextToken = output.NextToken
	}
	return nil
}
func describeQueryDefinitionsPages(ctx context.Context, conn *cloudwatchlogs.Client, input *cloudwatchlogs.DescribeQueryDefinitionsInput, fn func(*cloudwatchlogs.DescribeQueryDefinitionsOutput, bool) bool) error {
	for {
		output, err := conn.DescribeQueryDefinitions(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.ToString(output.NextToken) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextToken = output.NextToken
	}
	return nil
}
func describeResourcePoliciesPages(ctx context.Context, conn *cloudwatchlogs.Client, input *cloudwatchlogs.DescribeResourcePoliciesInput, fn func(*cloudwatchlogs.DescribeResourcePoliciesOutput, bool) bool) error {
	for {
		output, err := conn.DescribeResourcePolicies(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.ToString(output.NextToken) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.NextToken = output.NextToken
	}
	return nil
}
