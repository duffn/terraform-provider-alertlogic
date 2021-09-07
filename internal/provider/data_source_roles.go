package provider

import (
	"context"
	"fmt"

	"github.com/duffn/go-alertlogic/alertlogic"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRolesRead,
		Description: `A list of Alert Logic roles, both account specific and global.

[API reference](https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_Role_Resources-ListRoles)`,
		Schema: map[string]*schema.Schema{
			"roles": {
				Type:        schema.TypeList,
				Description: "A list of roles.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The role's ID.",
						},
						"account_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Account ID that holds the role, or '*' if the role is global.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The role's name",
						},
						"permissions": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The role's permissions.",
						},
						"version": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The version number of the role.",
						},
						"global": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates whether or not the role is a global role.",
						},
						"legacy_permissions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Legacy permissions of this role.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
				},
			},
		},
	}
}

func dataSourceRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	var diags diag.Diagnostics

	roles, err := api.ListRoles()
	if err != nil {
		return diag.FromErr(err)
	}

	roleDetails := make([]interface{}, 0)
	roleIds := make([]string, 0)
	for _, v := range roles.Roles {

		roleDetails = append(roleDetails, map[string]interface{}{
			"id":                 v.ID,
			"account_id":         v.AccountID,
			"name":               v.Name,
			"permissions":        v.Permissions,
			"version":            v.Version,
			"global":             v.Global,
			"legacy_permissions": v.LegacyPermissions,
			"created":            map[string]interface{}{"at": fmt.Sprint(v.Created.At), "by": v.Created.By},
			"modified":           map[string]interface{}{"at": fmt.Sprint(v.Modified.At), "by": v.Modified.By},
		})
		roleIds = append(roleIds, v.ID)
	}

	if err := d.Set("roles", roleDetails); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(stringListChecksum(roleIds))

	return diags
}
