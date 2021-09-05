package provider

import (
	"context"

	"github.com/duffn/go-alertlogic/alertlogic"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "An Alert Logic User.",
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
			"webhook_url": {
				Description: "The user's webhook URL.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"notifications_only": {
				Description: "Make the user a notifications only user.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	var diags diag.Diagnostics

	createUser := alertlogic.CreateUserRequest{
		Name:              d.Get("name").(string),
		Email:             d.Get("email").(string),
		Active:            d.Get("active").(bool),
		MobilePhone:       d.Get("mobile_phone").(string),
		WebhookUrl:        d.Get("webhook_url").(string),
		NotificationsOnly: d.Get("notifications_only").(bool),
	}

	user, err := api.CreateUser(createUser, false)
	if err != nil {
		diag.FromErr(err)
	}

	d.SetId(user.ID)
	resourceUserRead(ctx, d, meta)
	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	var diags diag.Diagnostics

	userId := d.Id()

	user, err := api.GetUserDetailsById(userId, false, false, false)
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

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	userId := d.Id()

	userRequest := alertlogic.UpdateUserRequest{
		Name:              d.Get("name").(string),
		Email:             d.Get("email").(string),
		Active:            d.Get("active").(bool),
		MobilePhone:       d.Get("mobile_phone").(string),
		WebhookUrl:        d.Get("webhook_url").(string),
		NotificationsOnly: d.Get("notifications_only").(bool),
	}

	_, err := api.UpdateUserDetails(userId, userRequest, false)
	if err != nil {
		return diag.FromErr(err)
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
