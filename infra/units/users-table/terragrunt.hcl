include "root" {
  path = find_in_parent_folders("root.hcl")
}

terraform {
  source = "git::github.com/rx-refill-reminders/terraform-modules//modules/dynamodb-table?ref=dynamodb-table%2Fv0&depth=0"
}

inputs = {
  table_name = "users"
  hash_key   = "id"
  attributes = [
    { name = "id", type = "S" },
  ]
  ttl_attribute = "ttl"
}
