data "alertlogic_roles" "all_roles" {}

output "roles" {
  value = data.alertlogic_roles.all_roles.roles
}
