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

func waitThesaurusCreated(ctx context.Context, conn *kendra.Client, id, indexId string, timeout time.Duration) (*kendra.DescribeThesaurusOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{string(types.ThesaurusStatusCreating)},
		Target:                    []string{string(types.ThesaurusStatusActive)},
		Refresh:                   statusThesaurus(ctx, conn, id, indexId),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*kendra.DescribeThesaurusOutput); ok {
		if out.Status == types.ThesaurusStatusFailed {
			tfresource.SetLastError(err, errors.New(aws.ToString(out.ErrorMessage)))
		}
		return out, err
	}

	return nil, err
}

func waitThesaurusUpdated(ctx context.Context, conn *kendra.Client, id, indexId string, timeout time.Duration) (*kendra.DescribeThesaurusOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{string(types.QuerySuggestionsBlockListStatusUpdating)},
		Target:                    []string{string(types.QuerySuggestionsBlockListStatusActive)},
		Refresh:                   statusThesaurus(ctx, conn, id, indexId),
		Timeout:                   timeout,
		NotFoundChecks:            20,
		ContinuousTargetOccurence: 2,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*kendra.DescribeThesaurusOutput); ok {
		if out.Status == types.ThesaurusStatusActiveButUpdateFailed || out.Status == types.ThesaurusStatusFailed {
			tfresource.SetLastError(err, errors.New(aws.ToString(out.ErrorMessage)))
		}
		return out, err
	}

	return nil, err
}

func waitThesaurusDeleted(ctx context.Context, conn *kendra.Client, id, indexId string, timeout time.Duration) (*kendra.DescribeThesaurusOutput, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{string(types.QuerySuggestionsBlockListStatusDeleting)},
		Target:  []string{},
		Refresh: statusThesaurus(ctx, conn, id, indexId),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)
	if out, ok := outputRaw.(*kendra.DescribeThesaurusOutput); ok {
		if out.Status == types.ThesaurusStatusFailed {
			tfresource.SetLastError(err, errors.New(aws.ToString(out.ErrorMessage)))
		}
		return out, err
	}

	return nil, err
}
