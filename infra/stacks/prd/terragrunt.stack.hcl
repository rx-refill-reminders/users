locals {
  hosted_zone_id = "Z07465232HRS85ZSQYRZY"
  domain         = "rx-refill-reminders.com"

  cognito_resource_server_scopes = [
    {
      name        = "access"
      description = "Full API access"
    },
  ]
}

unit "cognito_user_pool" {
  source = "${get_repo_root()}/infra/units/cognito-user-pool"
  path   = "cognito-user-pool"

  values = {
    pool_name = "rx-refill-reminders"

    resource_server_identifier = "https://api.rx-refill-reminders.com"
    resource_server_name       = "Rx Refill Reminders API"
    resource_server_scopes     = local.cognito_resource_server_scopes

    enable_apple_signin  = false
    enable_google_signin = false

    lambda_trigger_arns = {
      post_confirmation = "arn:aws:lambda:us-east-1:104875668206:function:cognito-postconfirm-handler"
    }

    domain = {
      mode = "user-hosted"
      user_hosted = {
        hosted_zone_id  = local.hosted_zone_id
        domain          = "auth.${local.domain}"
        certificate_arn = "arn:aws:acm:us-east-1:104875668206:certificate/481fc667-153a-4d14-8f9a-951e8a90cb36"
      }
    }

    clients = {
      m2m = {
        automations = {}
      }
      apps = {
        web = {
          callback_urls = ["https://${local.domain}/callback"]
          logout_urls   = ["https://${local.domain}/logout"]
        }
        ios = {
          callback_urls = ["rxrefillreminders://callback"]
          logout_urls   = ["rxrefillreminders://logout"]
        }
      }
    }
  }
}

unit "users_table" {
  source = "${get_repo_root()}/infra/units/users-table"
  path   = "users-table"

  values = {
    table_name = "users"
    hash_key   = "id"
    attributes = [
      { name = "id", type = "S" },
    ]
  }
}
