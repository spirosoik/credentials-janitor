## Providers

| Name | Version |
|------|---------|
| aws | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:-----:|
| environment | The environment where the lambda function runs | `string` | n/a | yes |
| function\_name | The name of the lambda function | `string` | `"credentials-janitor"` | no |
| janitor\_lambda\_schedule | The schedule of triggering a cloudwatch event to invoke lambda | `string` | `"cron(0 10 * * ? *)"` | no |
| lambda\_timeout | The name of the lambda function | `number` | `120` | no |
| max\_expiration\_days | The number of days for the rule about revoking credentials | `number` | `90` | no |
| private\_subet\_ids | n/a | `list(string)` | n/a | yes |
| s3\_bucket | The name of the bucket where the lambda artifact will exist | `string` | n/a | yes |
| s3\_key | The key/path where the lambda artifact will exist | `string` | n/a | yes |

## Outputs

No output.

