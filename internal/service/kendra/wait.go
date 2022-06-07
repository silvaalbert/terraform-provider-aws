package kendra

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kendra"
	"github.com/aws/aws-sdk-go-v2/service/kendra/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func waitQuerySuggestionsBlockListCreated(ctx context.Context, conn *kendra.Client, id, indexId string, timeout time.Duration) (*kendra.DescribeQuerySuggestionsBlockListOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{string(types.QuerySuggestionsBlockListStatusCreating)},
		Target:                    []string{string(types.QuerySuggestionsBlockListStatusActive)},
		Refresh:                   statusQuerySuggestionsBlockList(ctx, conn, id, indexId),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*kendra.DescribeQuerySuggestionsBlockListOutput); ok {
		if out.Status == types.QuerySuggestionsBlockListStatusFailed {
			tfresource.SetLastError(err, errors.New(aws.ToString(out.ErrorMessage)))
		}
		return out, err
	}

	return nil, err
}

func waitQuerySuggestionsBlockListUpdated(ctx context.Context, conn *kendra.Client, id, indexId string, timeout time.Duration) (*kendra.DescribeQuerySuggestionsBlockListOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{string(types.QuerySuggestionsBlockListStatusUpdating)},
		Target:                    []string{string(types.QuerySuggestionsBlockListStatusActive)},
		Refresh:                   statusQuerySuggestionsBlockList(ctx, conn, id, indexId),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*kendra.DescribeQuerySuggestionsBlockListOutput); ok {
		if out.Status == types.QuerySuggestionsBlockListStatusActiveButUpdateFailed || out.Status == types.QuerySuggestionsBlockListStatusFailed {
			tfresource.SetLastError(err, errors.New(aws.ToString(out.ErrorMessage)))
		}
		return out, err
	}

	return nil, err
}

func waitQuerySuggestionsBlockListDeleted(ctx context.Context, conn *kendra.Client, id, indexId string, timeout time.Duration) (*kendra.DescribeQuerySuggestionsBlockListOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{string(types.QuerySuggestionsBlockListStatusDeleting)},
		Target:  []string{},
		Refresh: statusQuerySuggestionsBlockList(ctx, conn, id, indexId),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*kendra.DescribeQuerySuggestionsBlockListOutput); ok {
		if out.Status == types.QuerySuggestionsBlockListStatusFailed {
			tfresource.SetLastError(err, errors.New(aws.ToString(out.ErrorMessage)))
		}
		return out, err
	}

	return nil, err
}
