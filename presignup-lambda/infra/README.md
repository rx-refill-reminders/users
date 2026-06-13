# api/infra

AWS infrastructure for the API Lambda, deployed per environment with Terragrunt stacks.

## Tech stack

- **Terraform** (>= 1.10) — Lambda function
- **Terragrunt** — stack orchestration, remote state, and generated units
- **AWS** (`us-east-1`) — Lambda, CloudWatch Logs
- **Remote state** — S3 backend with lockfile (`root.hcl`, per-stack `stack.yml`)

## Prerequisites

1. **Lambda IAM role** — apply [`infra` `lambda-role`](../infra/stacks/README.md) for the same env, then set `role-arn` in `stacks/<env>/stack.yml`.
2. **Code bucket** — set `code-bucket-id` in `stacks/<env>/stack.yml` (existing shared bucket).
3. **Deployment zip** — build `api/src/dist/api.zip` before apply (see `api/src`).

## Directory structure

```
api/infra/
├── root.hcl              # Shared backend + AWS provider
├── Makefile              # fmt (terragrunt)
├── units/                # Terragrunt unit templates (referenced by stacks)
└── stacks/               # Per-environment stacks (dev, prd)
    ├── <env>/stack.yml           # Account, state bucket, code bucket ID, role ARN
    └── <env>/terragrunt.stack.hcl  # Unit blueprint
```

Terragrunt generates runnable units under `stacks/<env>/.terragrunt-stack/` (gitignored). See [`stacks/README.md`](stacks/README.md) for plan/apply commands.
