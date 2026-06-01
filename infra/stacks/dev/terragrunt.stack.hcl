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
    pool_name                  = "rx-refill-reminders"
    app_client_name            = "rx-refill-reminders"
    domain_prefix              = "rx-refill-reminders-dev"
    resource_server_identifier = "https://api.rx-refill-reminders.com"
    resource_server_name       = "Rx Refill Reminders API"
    resource_server_scopes     = local.cognito_resource_server_scopes

    enable_web_client     = true
    web_callback_urls = [
      "https://app.rx-refill-reminders.com/callback"
    ]
    web_logout_urls = [
      "https://app.rx-refill-reminders.com/logout"
    ]

    ios_callback_urls = [
      "rxrefillreminders://callback",
    ]
    ios_logout_urls = [
      "rxrefillreminders://logout",
    ]

    enable_service_client = true
    enable_apple_signin   = false
    enable_google_signin  = false

    domain = {
      zone_id  = local.hosted_zone_id
      hostname = "auth.${local.domain}"
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
