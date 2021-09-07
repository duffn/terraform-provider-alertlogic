data "alertlogic_account" "my_account" {}

output "my_account" {
  value = data.alertlogic_account.my_account
}
