---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "alertlogic_account Data Source - terraform-provider-alertlogic"
subcategory: ""
description: |-
  Details about an Alert Logic account.
  API reference https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_Account_Resources-GetAccountDetails
---

# alertlogic_account (Data Source)

Details about an Alert Logic account.

[API reference](https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_Account_Resources-GetAccountDetails)

## Example Usage

```terraform
data "alertlogic_account" "my_account" {}

output "my_account" {
  value = data.alertlogic_account.my_account
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **accessible_locations** (List of String) Locations that this account can access.
- **active** (Boolean) The status of the account.
- **created** (Map of String) Information on when the record was created.
- **default_location** (String) Default location of the account.
- **modified** (Map of String) Information on when the record was modified.
- **name** (String) The account name.
- **version** (Number) The version number of the account.

### Read-Only

- **id** (String) The Alert Logic account ID.


