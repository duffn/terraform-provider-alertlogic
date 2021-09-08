package provider

import (
	"context"
	"fmt"

	"github.com/duffn/go-alertlogic/alertlogic"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountRead,
		Description: `Details about an Alert Logic account.

[API reference](https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_Account_Resources-GetAccountDetails)`,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The Alert Logic account ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"active": {
				Description: "The status of the account.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"name": {
				Description: "The account name.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"version": {
				Description: "The version number of the account.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"accessible_locations": {
				Description: "Locations that this account can access.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"default_location": {
				Description: "Default location of the account.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"created": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Information on when the record was created.",
			},
			"modified": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Information on when the record was modified.",
			},
		},
	}
}

func dataSourceAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	var diags diag.Diagnostics

	account, err := api.GetAccountDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", account.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("active", account.Active); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", account.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", account.Version); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("accessible_locations", account.AccessibleLocations); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_location", account.DefaultLocation); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created", map[string]interface{}{"at": fmt.Sprint(account.Created.At), "by": account.Created.By}); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("modified", map[string]interface{}{"at": fmt.Sprint(account.Modified.At), "by": account.Modified.By}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(account.ID)

	return diags
}
