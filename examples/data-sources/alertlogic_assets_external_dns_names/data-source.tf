data "alertlogic_assets_external_dns_names" "dns_names" {}

output "dns_names" {
  value = data.alertlogic_assets_external_dns_names.dns_names.external_dns_names
}
