include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "git::github.com/rx-refill-reminders/terraform-modules//modules/cognito-user-pool?ref=user-pool-dummy&depth=0"
}

inputs = values
