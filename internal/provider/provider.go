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
				"access_key_id": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("ALERTLOGIC_ACCESS_KEY_ID", nil),
					Description: "Your Alert Logic API access key ID.",
				},
				"secret_key": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("ALERTLOGIC_SECRET_KEY", nil),
					Description: "Your Alert Logic API secret key.",
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"alertlogic_user":                     resourceUser(),
				"alertlogic_assets_external_dns_name": resourceAssetsExternalDnsName(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"alertlogic_users":                     dataSourceUsers(),
				"alertlogic_roles":                     dataSourceRoles(),
				"alertlogic_global_roles":              dataSourceGlobalRoles(),
				"alertlogic_account":                   dataSourceAccount(),
				"alertlogic_assets_external_dns_names": dataSourceAssetsExternalDNSNames(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics

		access_key_id := d.Get("access_key_id").(string)
		secret_key := d.Get("secret_key").(string)
		accountId := d.Get("account_id").(string)

		if access_key_id != "" && secret_key != "" && accountId != "" {
			c, err := alertlogic.NewWithAccessKey(accountId, access_key_id, secret_key)
			if err != nil {
				return nil, diag.FromErr(err)
			}

			return c, diags
		} else {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Alert Logic client",
				Detail:   "You must set access_key_id, secret_key, and account_id",
			})
			return nil, diags
		}
	}
}
