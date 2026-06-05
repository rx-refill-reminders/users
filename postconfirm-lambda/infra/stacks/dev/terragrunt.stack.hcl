locals {
  stack_config = yamldecode(file("${get_terragrunt_dir()}/stack.yml"))
}

unit "lambda_function" {
  source = "${get_repo_root()}/postconfirm-lambda/infra/units/lambda-function"
  path   = "lambda-function"

  values = {
    code_bucket_id = "lambda-source-code-339284817422-us-east-1-an"
    role_arn       = "arn:aws:iam::339284817422:role/backend-api-lambda"

    aws_region = "us-east-1"

    users_table_name = "users"
  }
}
