package provider

import (
	"context"
	"fmt"

	"github.com/duffn/go-alertlogic/alertlogic"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"users": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"email": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"locked": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"mfa_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"version": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"linked_users": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_id": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"location": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"created": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"modified": {
							Type:     schema.TypeMap,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	var diags diag.Diagnostics

	users, err := api.ListUsers(false, false, false, "")
	if err != nil {
		return diag.FromErr(err)
	}

	userDetails := make([]interface{}, 0)
	userIds := make([]string, 0)
	for _, v := range users.Users {
		linkedUsersDetails := make([]interface{}, 0)

		for _, u := range v.LinkedUsers {
			linkedUsersDetails = append(linkedUsersDetails, map[string]interface{}{
				"user_id":  u.UserID,
				"location": u.Location,
			})
		}

		userDetails = append(userDetails, map[string]interface{}{
			"id":           v.ID,
			"account_id":   v.AccountID,
			"name":         v.Name,
			"username":     v.Username,
			"email":        v.Email,
			"active":       v.Active,
			"locked":       v.Locked,
			"mfa_enabled":  v.MfaEnabled,
			"version":      v.Version,
			"linked_users": linkedUsersDetails,
			"created":      map[string]interface{}{"at": fmt.Sprint(v.Created.At), "by": v.Created.By},
			"modified":     map[string]interface{}{"at": fmt.Sprint(v.Modified.At), "by": v.Modified.By},
		})
		userIds = append(userIds, v.ID)
	}

	if err := d.Set("users", userDetails); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(stringListChecksum(userIds))

	return diags
}
