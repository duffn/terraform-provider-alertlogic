data "alertlogic_users" "all_users" {}

output "users" {
  value = data.alertlogic_users.all_users.users
}
