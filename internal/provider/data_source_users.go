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
		Description: `A list of Alert Logic users.

[API reference](https://console.cloudinsight.alertlogic.com/api/aims/#api-AIMS_User_Resources-ListUsers)
		`,
		Schema: map[string]*schema.Schema{
			"users": {
				Type:        schema.TypeList,
				Description: "A list of users.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user's ID.",
						},
						"account_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Account ID that holds the user.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The user's full name",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The user's username.",
						},
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The user's email address.",
						},
						"active": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates whether or not the user is active.",
						},
						"locked": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates whether or not the user is allowed to log in.",
						},
						"mfa_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates the status of the users MFA.",
						},
						"version": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The version of the user's details; i.e. how many times has the user been updated.",
						},
						"linked_users": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Users linked to this user.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_id": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The ID of the user.",
									},
									"location": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The location of the user.",
									},
								},
							},
						},
						"role_ids": {
							Description: "Role IDs for the user.",
							Type:        schema.TypeList,
							Optional:    true,
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
		roleIds, err := api.GetAssignedRoleIDs(v.ID)
		if err != nil {
			return diag.FromErr(err)
		}

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
			"role_ids":     roleIds.RoleIds,
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
