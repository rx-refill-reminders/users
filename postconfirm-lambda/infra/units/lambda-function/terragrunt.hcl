include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "git::github.com/rx-refill-reminders/terraform-modules//modules/lambda-function?ref=lambda-function%2Fv0&depth=0"
}

inputs = {
  function_name  = "cognito-postconfirm-handler"
  handler        = "bootstrap"
  executable_zip = "${get_repo_root()}/postconfirm-lambda/src/dist/api.zip"

  runtime         = "provided.al2023"
  timeout_seconds = 10

  code_bucket_id = values.code_bucket_id
  role_arn       = values.role_arn

  env_vars = {
    AWS_REGION  = values.aws_region
    USERS_TABLE = values.users_table_name
  }
}
