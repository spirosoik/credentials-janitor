variable "s3_bucket" {
  type = string
  description = "The name of the bucket where the lambda artifact will exist"
}

variable "s3_key" {
  type = string
  description = "The key/path where the lambda artifact will exist"
}

variable "function_name" {
  type = string
  default     = "credentials-janitor"
  description = "The name of the lambda function"
}

variable "lambda_timeout" {
  default     = 120
  description = "The name of the lambda function"
}

variable "private_subet_ids" {
  type = list(string)
  description = ""
}

variable "environment" {
  type = string
  description = "The environment where the lambda function runs"
}

variable "max_expiration_days" {
  type = number
  default     = 90
  description = "The number of days for the rule about revoking credentials"
}

variable "janitor_lambda_schedule" {
  type = string
  default     = "cron(0 10 * * ? *)"
  description = "The schedule of triggering a cloudwatch event to invoke lambda"
}