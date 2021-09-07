package provider

import (
	"context"

	"github.com/duffn/go-alertlogic/alertlogic"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAccount() *schema.Resource {
	return &schema.Resource{
		Description:   "An Alert Logic account. You cannot actually create or delete an Alert Logic account with this resource. The only thing you are allowed to do, per restrictions of the Alert Logic API, is set whether or not your organization requires MFA.",
		CreateContext: resourceAccountCreate,
		ReadContext:   resourceAccountRead,
		UpdateContext: resourceAccountUpdate,
		DeleteContext: resourceAccountDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The Alert Logic account ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"mfa_required": {
				Description: "Indicates whether or not the account has MFA required.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
		},
	}
}

func resourceAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Alert Logic account created in Terraform state",
		Detail:   "Your Alert Logic account was added to the Terraform state (but you didn't actually create a new Alert Logic account).",
	})

	account, err := api.GetAccountDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(account.ID)
	resourceAccountRead(ctx, d, meta)
	return diags
}

func resourceAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	var diags diag.Diagnostics

	account, err := api.GetAccountDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("mfa_required", account.Name); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	mfa_required := d.Get("mfa_required").(bool)

	_, err := api.UpdateAccountDetails(alertlogic.UpdateAccountDetailsRequest{MfaRequired: mfa_required})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceAccountRead(ctx, d, meta)
}

func resourceAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Alert Logic account removed from Terraform state",
		Detail:   "Your Alert Logic account was removed from Terraform state (but not actually deleted).",
	})

	d.SetId("")

	return diags
}
