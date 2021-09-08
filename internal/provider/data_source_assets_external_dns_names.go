package provider

import (
	"context"
	"fmt"

	"github.com/duffn/go-alertlogic/alertlogic"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAssetsExternalDNSNames() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetsExternalDNSNamesRead,
		Description: `A list of external DNS name assets.

[API reference](https://console.cloudinsight.alertlogic.com/api/assets_query/#api-Queries-QueryAccountAssets)`,
		Schema: map[string]*schema.Schema{
			"external_dns_names": {
				Type:        schema.TypeList,
				Description: "A list of external DNS name assets.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The version number of the asset.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of the asset.",
						},
						"threatiness": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "Threatiness?",
						},
						"threat_level": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The threat level of the external DNS name asset.",
						},
						"state": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The state of the asset.",
						},
						"native_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The native type of the asset.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The assets name.",
						},
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The asset key.",
						},
						"dns_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The asset's DNS name.",
						},
						"deployment_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the deployment that the external asset resides in.",
						},
						"deleted_on": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The epoch time the asset was deleted.",
						},
						"declared": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"created_on": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "When the asset was created in Alert Logic.",
						},
						"account_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Account ID that holds the asset.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAssetsExternalDNSNamesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	var diags diag.Diagnostics

	assets, err := api.GetExternalDNSNameAssets()
	if err != nil {
		return diag.FromErr(err)
	}

	assetDetails := make([]interface{}, 0)
	assetIds := make([]string, 0)
	// The API returns this as an array of arrays.
	for _, outerAssets := range assets.ExternalDNSAssets {
		for _, v := range outerAssets {
			assetDetails = append(assetDetails, map[string]interface{}{
				"version":       v.Version,
				"type":          v.Type,
				"threatiness":   v.Threatiness,
				"threat_level":  v.ThreatLevel,
				"state":         v.State,
				"native_type":   v.NativeType,
				"name":          v.Name,
				"key":           v.Key,
				"dns_name":      v.DNSName,
				"deployment_id": v.DeploymentID,
				"deleted_on":    v.DeletedOn,
				"declared":      v.Declared,
				"created_on":    v.CreatedOn,
				"account_id":    v.AccountID,
			})
			// There isn't an assigned ID from Alert Logic, but the
			// deployment and key of an asset can be the unique key.
			assetIds = append(assetIds, fmt.Sprintf("%s%s", v.DeploymentID, v.Key))
		}
	}

	if err := d.Set("external_dns_names", assetDetails); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(stringListChecksum(assetIds))

	return diags
}
