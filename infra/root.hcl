# Shared Terragrunt configuration for all stacks and units.

locals {
  stack = yamldecode(file(find_in_parent_folders("stack.yml")))

  project   = get_env("PIPELINE_PROJECT")
  component = get_env("PIPELINE_COMPONENT")
  
  providers_file_path    = try(find_in_parent_folders("providers.tf"), "")
  providers_file_exists  = local.providers_file_path != ""
  providers_file_content = local.providers_file_exists ? file(local.providers_file_path) : ""
}

remote_state {
  backend = "s3"

  config = {
    bucket       = local.stack["states-bucket"]
    key          = "${local.project}/${local.component}/${path_relative_to_include()}/terraform.tfstate"
    region       = "us-east-1"
    encrypt      = true
    use_lockfile = true
  }

  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }
}

generate "provider" {
  path      = "providers.tf"
  if_exists = "skip"
  disable   = !local.providers_file_exists
  contents  = local.providers_file_content
}
