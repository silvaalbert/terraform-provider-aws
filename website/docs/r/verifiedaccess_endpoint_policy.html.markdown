---
subcategory: "Verified Access"
layout: "aws"
page_title: "AWS: aws_verifiedaccess_endpoint_policy"
description: |-
  Terraform resource for managing an AWS Verified Access Endpoint Policy.
---

# Resource: aws_verifiedaccess_endpoint_policy

Terraform resource for managing an AWS Verified Access Endpoint Policy.

## Example Usage

### Basic Usage

```terraform
resource "aws_verifiedaccess_endpoint_policy" "example" {
  policy_document             = <<EOF
    permit(principal, action, resource) 
    when {
        context.http_request.method == "GET"
    };
  EOF
  verified_access_endpoint_id = aws_verifiedaccess_endpoint.example.id
}
```

## Argument Reference

The following arguments are required:

* `policy_document` - (Required) The Verified Access policy document.
* `verified_access_endpoint_id` - (Required) The ID of the Verified Access endpoint.

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

* `create` - (Default `5m`)
* `update` - (Default `5m`)
* `delete` - (Default `5m`)

## Import

AWS Verified Access Endpoint Policy can be imported using the `verified_access_endpoint_id`, e.g.,

```
$ terraform import aws_ec2_verifiedaccess_endpoint_policy.example vae-8012925589
```
