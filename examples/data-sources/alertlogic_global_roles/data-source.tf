data "alertlogic_global_roles" "global_roles" {}

output "roles" {
  value = data.alertlogic_roles.global_roles.roles
}
