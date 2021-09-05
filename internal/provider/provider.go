package provider

import (
	"context"

	"github.com/duffn/go-alertlogic/alertlogic"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"account_id": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("ALERTLOGIC_ACCOUNT_ID", nil),
					Description: "Your Alert Logic Account ID.",
				},
				"username": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("ALERTLOGIC_USERNAME", nil),
					Description: "Your Alert Logic username. While you can use your actual username here if you do not have MFA enabled, it is strongly recommended to use specific API access keys instead.",
				},
				"password": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("ALERTLOGIC_PASSWORD", nil),
					Description: "Your Alert Logic password. While you can use your actual password here if you do not have MFA enabled, it is strongly recommended to use specific API access keys instead.",
				},
			},
			ResourcesMap:   map[string]*schema.Resource{"alertlogic_user": resourceUser()},
			DataSourcesMap: map[string]*schema.Resource{"alertlogic_users": dataSourceUsers(), "alertlogic_roles": dataSourceRoles()},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics

		username := d.Get("username").(string)
		password := d.Get("password").(string)
		accountId := d.Get("account_id").(string)

		if username != "" && password != "" && accountId != "" {
			c, err := alertlogic.NewWithUsernameAndPassword(accountId, username, password)
			if err != nil {
				return nil, diag.FromErr(err)
			}

			return c, diags
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Alert Logic client",
				Detail:   "You must set username, password, and account_id",
			})
			return nil, diags
		}
	}
}
