# units

Terragrunt unit templates — referenced by `source` in each stack’s `terragrunt.stack.hcl`. Terragrunt generates runnable units under `<stack>/.terragrunt-stack/`.

| Unit | Module |
|------|--------|
| `dns-hosted-zone` | [`modules/dns-hosted-zone`](../modules/dns-hosted-zone) |
| `cognito-user-pool` | [`modules/cognito-user-pool`](../modules/cognito-user-pool) |
| `api-gateway` | [`modules/api-gateway`](../modules/api-gateway) |
| `api-gateway-routes` | [`modules/api-gateway-routes`](https://github.com/rx-refill-reminders/terraform-modules/tree/main/modules/api-gateway-routes) |
| `lambda-role` | [`lambda-role`](https://github.com/rx-refill-reminders/terraform-modules/tree/main/modules/lambda-role) |
| `dynamodb-table` | [`dynamodb-table`](https://github.com/rx-refill-reminders/terraform-modules/tree/main/modules/dynamodb-table) |
