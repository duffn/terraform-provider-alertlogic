package provider

import (
	"context"

	"github.com/duffn/go-alertlogic/alertlogic"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGlobalRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalRolesRead,
		Description: `A list of global Alert Logic roles.

[API reference](https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_Role_Resources-ListGlobalRoles)`,
		Schema: map[string]*schema.Schema{
			"roles": {
				Type:        schema.TypeList,
				Description: "A list of roles.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: roleSchema,
				},
			},
		},
	}
}

func dataSourceGlobalRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	var diags diag.Diagnostics

	roles, err := api.ListGlobalRoles()
	if err != nil {
		return diag.FromErr(err)
	}

	return formatRolesResponse(diags, d, roles)
}
