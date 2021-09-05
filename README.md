# terraform-provider-alertlogic

This is a Terraform provider for the [Alert Logic Cloud Insights API](https://console.cloudinsight.alertlogic.com/api/#/).

This is in _very early_ development and only supports a single of the [myriad of endpoints](https://console.cloudinsight.alertlogic.com/api/#/) of the API. Expect the API here to break often during early development.

## Usage

```hcl
terraform {
  required_providers {
    alertlogic = {
      version = "0.0.1"
      source  = "github.com/duffn/alertlogic"
    }
  }
}

variable "username" {}
variable "password" {}
variable "account_id" {}

provider "alertlogic" {
  username   = var.username
  password   = var.password
  account_id = var.account_id
}
```

### Users

```hcl
resource "alertlogic_user" "user" {
  name         = "Bob Loblaw"
  email        = "bob@bobloblawlaw.com"
  mobile_phone = "234-555-5555"
}
```