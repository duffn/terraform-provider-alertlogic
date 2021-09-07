variable "access_key_id" {}
variable "secret_key" {}
variable "account_id" {}

provider "alertlogic" {
  access_key_id = var.access_key_id
  secret_key    = var.secret_key
  account_id    = var.account_id
}
