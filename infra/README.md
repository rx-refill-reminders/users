# infra

AWS infrastructure for the rx-refill-reminders backend, defined as Terraform modules and deployed per environment with Terragrunt stacks.

## Tech stack

- **Terraform** (>= 1.10) — resource definitions in `modules/`
- **Terragrunt** — stack orchestration, remote state, and generated units
- **AWS** (`us-east-1`) — Route 53, ACM, Cognito
- **Remote state** — S3 backend with lockfile (`root.hcl`, per-stack `stack.yml`)

## Directory structure

```
infra/
├── root.hcl              # Shared backend + AWS provider
├── Makefile              # fmt (terraform + terragrunt)
├── modules/              # Reusable Terraform modules
├── units/                # Terragrunt unit templates (referenced by stacks)
└── stacks/               # Per-environment stacks (dev, prd)
    ├── <env>/stack.yml           # Account, role, state bucket
    └── <env>/terragrunt.stack.hcl  # Unit blueprint
```

Terragrunt generates runnable units under `stacks/<env>/.terragrunt-stack/` (gitignored). See [`stacks/README.md`](stacks/README.md) for plan/apply commands.
