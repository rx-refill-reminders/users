# units

Terragrunt unit templates — referenced by `source` in each stack’s `terragrunt.stack.hcl`. Terragrunt generates runnable units under `<stack>/.terragrunt-stack/`.

| Unit | Module |
|------|--------|
| `lambda-function` | [`lambda-function`](https://github.com/rx-refill-reminders/terraform-modules/tree/main/modules/lambda-function) |

`code_bucket_id` and `role_arn` are set per stack in `stack.yml` (`code-bucket-id`, `role-arn`). The IAM role is provisioned by [`infra`](../../../infra/units/lambda-role).
