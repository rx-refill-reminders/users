locals {
  hosted_zone_id = "Z08427401W2SCGIP77L8A"
  domain         = "dev.rx-refill-reminders.com"

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

    domain = {
      mode = "cognito-hosted"
      cognito_hosted = {
        hosted_zone_id  = local.hosted_zone_id
        domain          = "auth.${local.domain}"
        certificate_arn = "arn:aws:acm:us-east-1:339284817422:certificate/3ef26155-8494-4789-bae2-52d8299aa384"
      }
    }

    clients = {
      m2m = {
        automations = {}
      }
      app = {
        web = {
          callback_urls = ["https://app.${local.domain}/callback"]
          logout_urls   = ["https://app.${local.domain}/logout"]
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
    hash_key   = "userId"
    attributes = [
      { name = "userId", type = "S" },
    ]
  }
}
