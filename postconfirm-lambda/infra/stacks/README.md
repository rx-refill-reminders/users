# stacks

Terragrunt stacks (`dev`, `prd`). Each stack directory contains only:

| File | Purpose |
|------|---------|
| `stack.yml` | AWS account, state bucket, code bucket ID, Lambda role ARN |
| `terragrunt.stack.hcl` | Blueprint — generates the Lambda unit on demand |

Units are generated under `.terragrunt-stack/` (gitignored). Do not add per-unit folders under `prd/` or `dev/`.

Shared backend config: [`root.hcl`](../root.hcl) reads `stack.yml` (`use_lockfile = true`). Requires **Terraform >= 1.10**.

## Commands

```bash
cd api/infra/stacks/prd
terragrunt stack run plan    # generate + plan
terragrunt stack run apply
```

State key example: `prd/.terragrunt-stack/lambda-function/terraform.tfstate` (relative to `root.hcl`).

## Configuration

Set in `stack.yml`:

- `code-bucket-id` — existing shared Lambda code bucket
- `role-arn` — from [`infra` `lambda-role`](../../infra/stacks/README.md) output after apply

Unit template: [`../units/lambda-function/terragrunt.hcl`](../units/lambda-function/terragrunt.hcl).
