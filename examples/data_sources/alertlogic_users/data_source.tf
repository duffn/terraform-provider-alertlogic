data "alertlogic_users" "users" {}

output "users" {
  value = data.alertlogic_users.users.users
}
