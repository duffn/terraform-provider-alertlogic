package provider

import (
	"context"

	"github.com/duffn/go-alertlogic/alertlogic"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "An Alert Logic user.",
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "A full name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"email": {
				Description: "An email address.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"active": {
				Description: "The user's status.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"mobile_phone": {
				Description: "A mobile telephone number.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"role_ids": {
				Description: "An array of role IDs to grant to the user.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	createUser := alertlogic.CreateUserRequest{
		Name:        d.Get("name").(string),
		Email:       d.Get("email").(string),
		Active:      d.Get("active").(bool),
		MobilePhone: d.Get("mobile_phone").(string),
	}

	user, err := api.CreateUser(createUser, false)
	if err != nil {
		return diag.FromErr(err)
	}

	roleIds := expandInterfaceToStringList(d.Get("role_ids"))
	for _, roleId := range roleIds {
		_, err := api.GrantUserRole(user.ID, roleId)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(user.ID)
	return resourceUserRead(ctx, d, meta)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	var diags diag.Diagnostics

	userId := d.Id()

	user, err := api.GetUserDetailsById(userId, false, false, false)
	if err != nil {
		return diag.FromErr(err)
	}

	roleIds, err := api.GetAssignedRoleIDs(user.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", user.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email", user.Email); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("active", user.Active); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("mobile_phone", user.MobilePhone); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("role_ids", roleIds.RoleIds); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	userId := d.Id()

	userRequest := alertlogic.UpdateUserRequest{
		Name:        d.Get("name").(string),
		Email:       d.Get("email").(string),
		Active:      d.Get("active").(bool),
		MobilePhone: d.Get("mobile_phone").(string),
	}

	_, err := api.UpdateUserDetails(userId, userRequest, false)
	if err != nil {
		return diag.FromErr(err)
	}

	// Take care of the user's roles.
	planRoleIds := expandInterfaceToStringList(d.Get("role_ids"))
	currentAssignedRoleIds, err := api.GetAssignedRoleIDs(userId)
	if err != nil {
		return diag.FromErr(err)
	}

	// Assign roles to the user that are in the plan, but that they currently don't
	// have assigned to them in Alert Logic.
	for _, roleId := range planRoleIds {
		if !contains(currentAssignedRoleIds.RoleIds, roleId) {
			_, err := api.GrantUserRole(userId, roleId)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// Delete roles that are currently assigned to the user in Alert Logic, but
	// are not in the plan.
	for _, roleId := range currentAssignedRoleIds.RoleIds {
		if !contains(planRoleIds, roleId) {
			_, err := api.RevokeUserRole(userId, roleId)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceUserRead(ctx, d, meta)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	var diags diag.Diagnostics

	userId := d.Id()

	_, err := api.DeleteUser(userId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
