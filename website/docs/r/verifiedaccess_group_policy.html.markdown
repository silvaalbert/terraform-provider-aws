---
subcategory: "Verified Access"
layout: "aws"
page_title: "AWS: aws_verifiedaccess_group_policy"
description: |-
  Terraform resource for managing an AWS Verified Access Group Policy.
---

# Resource: aws_verifiedaccess_group_policy

Terraform resource for managing an AWS Verified Access Group Policy.

## Example Usage

### Basic Usage

```terraform
resource "aws_verifiedaccess_group_policy" "example" {
  policy_document             = <<EOF
    permit(principal, action, resource) 
    when {
        context.http_request.method == "GET"
    };
  EOF
  verified_access_group_id = aws_verifiedaccess_group.example.id
}
```

## Argument Reference

The following arguments are required:

* `policy_document` - (Required) The Verified Access policy document.
* `verified_access_group_id` - (Required) The ID of the Verified Access group.

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

* `create` - (Default `5m`)
* `update` - (Default `5m`)
* `delete` - (Default `5m`)

## Import

AWS Verified Access Group Policy can be imported using the `verified_access_group_id`, e.g.,

```
$ terraform import aws_ec2_verifiedaccess_group_policy.example vagr-8012925589
```
