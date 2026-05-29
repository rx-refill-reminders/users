# stacks

Terragrunt stacks (`dev`, `prd`). Each stack directory contains only:

| File | Purpose |
|------|---------|
| `stack.yml` | AWS account, assume-role ARN, S3 state bucket |
| `terragrunt.hcl` | Stack-level locals (e.g. domain) |
| `terragrunt.stack.hcl` | Blueprint — generates all units on demand |

Units are generated under `.terragrunt-stack/` (gitignored). Do not add per-unit folders under `prd/` or `dev/`.

Shared backend config: [`root.hcl`](root.hcl) reads `stack.yml` (`use_lockfile = true`). Requires **Terraform >= 1.10**.

## Commands

```bash
cd infra/stacks/prd
terragrunt stack run plan    # generate + plan
terragrunt stack run apply
```

State keys: `prd/.terragrunt-stack/dns-hosted-zone/terraform.tfstate` (relative to `root.hcl`). DynamoDB units use paths such as `dynamodb-users`, `dynamodb-rx-cycles`, `dynamodb-rx-cycle-instances`, and `dynamodb-reminders`.

## Adding a unit

1. Add a template under [`../units/<name>/terragrunt.hcl`](../units/).
2. Add a `unit` block to each env’s `terragrunt.stack.hcl`.

Unit templates: [`../units/dns-hosted-zone/terragrunt.hcl`](../units/dns-hosted-zone/terragrunt.hcl).

Stack domains: [`prd/terragrunt.hcl`](prd/terragrunt.hcl), [`dev/terragrunt.hcl`](dev/terragrunt.hcl).
