# account-credentials-janitor
A lambda function which detects and remove unused IAM 
credentials for AWS users.

## Purpose

Respecting security we need to remove IAM user credentials that are not
used anymore and notify them back that we revoked their credentials (login profile + access keys).

## How To

We are going to schedule a CloudWatch event to invoke the lambda function periodically. Lambda function
will do the listing and will check when they used last time the login profile and access keys. The max time
is configurable as environment variable which can be passed in the lambda deployment.

... TODO