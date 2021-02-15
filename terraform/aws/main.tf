data "aws_caller_identity" "current" {}

resource "aws_iam_role" "janitor_lambda_role" {
  name = "janitor_lambda_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "janitor_lambda_policy" {
  name = "janitor_lambda_policy"
  role = aws_iam_role.janitor_lambda_role.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
          "iam:ListUsers",
          "iam:ListAccessKeys",
          "iam:GetAccessKeyLastUsed",
          "iam:DeleteAccessKey",
          "iam:DeleteLoginProfile",
        ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "janitor_default_execution_role" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/DefaultLambdaExecutionRole"
  role       = aws_iam_role.janitor_lambda_role.name
}

resource "aws_iam_role_policy_attachment" "janitor_iam_executiron_role" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/IAMLambdaExecutironRole"
  role       = aws_iam_role.janitor_lambda_role.name
}

resource "aws_lambda_function" "janitor" {
  s3_bucket     = var.s3_bucket
  s3_key        = var.s3_key
  function_name = var.function_name
  role          = aws_iam_role.janitor_lambda_role.arn
  handler       = "main"
  timeout       = var.lambda_timeout
  runtime       = "go1.x"
  vpc_config {
    subnet_ids         = flatten(var.private_subnet_ids)
    security_group_ids = [aws_security_group.janitor_lambda_sg.id]
  }

  environment {
    variables = {
      JANITOR_ENVIRONMENT         = var.environment
      JANITOR_MAX_EXPIRATION_DAYS = var.max_expiration_days
    }
  }
}

resource "aws_security_group" "janitor_lambda_sg" {
  name        = "${var.deployment_name}-credentials-janitor-lambda-sg"
  description = "Credentials Janitor Lambda"
  vpc_id      = var.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.deployment_name}-credentials-janitor-lambda-sg"
  }
}

resource "aws_cloudwatch_event_rule" "janitor" {
  name                = "credentials-janitor"
  description         = "Runs based on the schedule expression"
  schedule_expression = var.janitor_lambda_schedule
}

resource "aws_cloudwatch_event_target" "janitor" {
  rule      = aws_cloudwatch_event_rule.janitor.name
  target_id = "credentials-janitor"
  arn       = aws_lambda_function.janitor.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_trigger_janitor" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.janitor.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.janitor.arn
}