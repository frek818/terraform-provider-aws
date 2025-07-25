---
subcategory: "CloudWatch Observability Access Manager"
layout: "aws"
page_title: "AWS: aws_oam_link"
description: |-
  Terraform resource for managing an AWS CloudWatch Observability Access Manager Link.
---

# Resource: aws_oam_link

Terraform resource for managing an AWS CloudWatch Observability Access Manager Link.

~> **NOTE:** Creating an `aws_oam_link` may sometimes fail if the `aws_oam_sink_policy` for the attached `aws_oam_sink` is not created before the `aws_oam_link`. To prevent this, declare an explicit dependency using a [`depends_on`](https://developer.hashicorp.com/terraform/language/meta-arguments/depends_on) meta-argument.

## Example Usage

### Basic Usage

```terraform
resource "aws_oam_link" "example" {
  label_template  = "$AccountName"
  resource_types  = ["AWS::CloudWatch::Metric"]
  sink_identifier = aws_oam_sink.example.arn
  tags = {
    Env = "prod"
  }

  depends_on = [
    aws_oam_sink_policy.example
  ]
}

resource "aws_oam_sink" "example" {
  # ...
}

resource "aws_oam_sink_policy" "example" {
  sink_identifier = aws_oam_sink.example.arn
  # ...
}
```

### Log Group Filtering

```terraform
resource "aws_oam_link" "example" {
  label_template = "$AccountName"
  link_configuration {
    log_group_configuration {
      filter = "LogGroupName LIKE 'aws/lambda/%' OR LogGroupName LIKE 'AWSLogs%'"
    }
  }
  resource_types  = ["AWS::Logs::LogGroup"]
  sink_identifier = aws_oam_sink.example.arn

  depends_on = [
    aws_oam_sink_policy.example
  ]
}
```

### Metric Filtering

```terraform
resource "aws_oam_link" "example" {
  label_template = "$AccountName"
  link_configuration {
    metric_configuration {
      filter = "Namespace IN ('AWS/EC2', 'AWS/ELB', 'AWS/S3')"
    }
  }
  resource_types  = ["AWS::CloudWatch::Metric"]
  sink_identifier = aws_oam_sink.example.arn

  depends_on = [
    aws_oam_sink_policy.example
  ]
}
```

## Argument Reference

The following arguments are required:

* `label_template` - (Required) Human-readable name to use to identify this source account when you are viewing data from it in the monitoring account.
* `resource_types` - (Required) Types of data that the source account shares with the monitoring account.
* `sink_identifier` - (Required) Identifier of the sink to use to create this link.

The following arguments are optional:

* `region` - (Optional) Region where this resource will be [managed](https://docs.aws.amazon.com/general/latest/gr/rande.html#regional-endpoints). Defaults to the Region set in the [provider configuration](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#aws-configuration-reference).
* `link_configuration` - (Optional) Configuration for creating filters that specify that only some metric namespaces or log groups are to be shared from the source account to the monitoring account. See [`link_configuration` Block](#link_configuration-block) for details.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

### `link_configuration` Block

The `link_configuration` configuration block supports the following arguments:

* `log_group_configuration` - (Optional) Configuration for filtering which log groups are to send log events from the source account to the monitoring account. See [`log_group_configuration` Block](#log_group_configuration-block) for details.
* `metric_configuration` - (Optional) Configuration for filtering which metric namespaces are to be shared from the source account to the monitoring account. See [`metric_configuration` Block](#metric_configuration-block) for details.

### `log_group_configuration` Block

The `log_group_configuration` configuration block supports the following arguments:

* `filter` - (Required) Filter string that specifies which log groups are to share their log events with the monitoring account. See [LogGroupConfiguration](https://docs.aws.amazon.com/OAM/latest/APIReference/API_LogGroupConfiguration.html) for details.

### `metric_configuration` Block

The `metric_configuration` configuration block supports the following arguments:

* `filter` - (Required) Filter string that specifies  which metrics are to be shared with the monitoring account. See [MetricConfiguration](https://docs.aws.amazon.com/OAM/latest/APIReference/API_MetricConfiguration.html) for details.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `arn` - ARN of the link.
* `id` - ARN of the link. Use `arn` instead.
* `label` - Label that is assigned to this link.
* `link_id` - ID string that AWS generated as part of the link ARN.
* `sink_arn` - ARN of the sink that is used for this link.

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

* `create` - (Default `1m`)
* `update` - (Default `1m`)
* `delete` - (Default `1m`)

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import CloudWatch Observability Access Manager Link using the `arn`. For example:

```terraform
import {
  to = aws_oam_link.example
  id = "arn:aws:oam:us-west-2:123456789012:link/link-id"
}
```

Using `terraform import`, import CloudWatch Observability Access Manager Link using the `arn`. For example:

```console
% terraform import aws_oam_link.example arn:aws:oam:us-west-2:123456789012:link/link-id
```
