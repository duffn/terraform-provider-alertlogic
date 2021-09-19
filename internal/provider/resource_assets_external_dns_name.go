package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/duffn/go-alertlogic/alertlogic"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAssetsExternalDnsName() *schema.Resource {
	return &schema.Resource{
		Description: `An Alert Logic external DNS asset.
This is only the asset of the ` + "`external-dns-name`" + ` type. You cannot create any other asset type with this resource.

[API reference](https://console.cloudinsight.alertlogic.com/api/assets_write/#api-DeclareModify-DeclareAsset)`,
		CreateContext: resourceAssetsExternalDnsNameCreate,
		ReadContext:   resourceAssetsExternalDnsNameRead,
		DeleteContext: resourceAssetsExternalDnsNameDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				deploymentId, dnsName, err := parseAssetImportId(d.Id())
				if err != nil {
					return nil, err
				}

				d.Set("deployment_id", deploymentId)
				d.Set("dns_name", dnsName)
				d.SetId(getAssetId(deploymentId, dnsName))

				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"deployment_id": {
				Description: "The deployment ID of the external asset.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"dns_name": {
				Description: "The external DNS name of the asset.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceAssetsExternalDnsNameCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	deploymentId := d.Get("deployment_id").(string)
	dnsName := d.Get("dns_name").(string)

	assetId := getAssetId(deploymentId, dnsName)

	_, err := api.CreateExternalDNSNameAsset(deploymentId, dnsName)
	if err != nil {
		return diag.FromErr(err)
	}

	// The external assets takes a moment to create so we need to retry, waiting
	// for it to be available.
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		assets, err := api.GetExternalDNSNameAssets()
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error describing external DNS names: %s", err))
		}

		if getAssetFromList(assetId, assets) == (alertlogic.ExternalDNSNameAsset{}) {
			return resource.RetryableError(fmt.Errorf("expected asset to be created but it is not available yet"))
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(assetId)
	return resourceAssetsExternalDnsNameRead(ctx, d, meta)
}

func resourceAssetsExternalDnsNameRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	var diags diag.Diagnostics

	assetId := d.Id()

	assets, err := api.GetExternalDNSNameAssets()
	if err != nil {
		return diag.FromErr(err)
	}

	thisAsset := getAssetFromList(assetId, assets)

	if err := d.Set("deployment_id", thisAsset.DeploymentID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dns_name", thisAsset.DNSName); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAssetsExternalDnsNameDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)

	var diags diag.Diagnostics

	deploymentId := d.Get("deployment_id").(string)
	dnsName := d.Get("dns_name").(string)

	_, err := api.RemoveExternalDNSNameAsset(deploymentId, dnsName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

// assetId returns the asset ID for the external DNS names asset.
func getAssetId(deploymentId string, dnsName string) string {
	return fmt.Sprintf("%s/external-dns-name/%s", deploymentId, dnsName)
}

// getAssetFromList gets a single asset from a slice of external assets if it matches
// the assetId format.
func getAssetFromList(assetId string, assets alertlogic.ExternalDNSNameAssets) alertlogic.ExternalDNSNameAsset {
	var thisAsset alertlogic.ExternalDNSNameAsset
	for _, dnsAssets := range assets.ExternalDNSAssets {
		for _, asset := range dnsAssets {
			if getAssetId(asset.DeploymentID, asset.DNSName) == assetId {
				thisAsset = asset
			}
		}
	}

	return thisAsset
}

// parseAssetImportId parses an ID passed to the import function. The ID should be in
// the format `deploymentId/dnsName`
func parseAssetImportId(assetId string) (string, string, error) {
	parts := strings.SplitN(assetId, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected deploymentId/dnsName", assetId)
	}

	return parts[0], parts[1], nil
}
