package provider

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

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
							Type:     schema.TypeString,
							Optional: true,
						},
						"mfa_enabled": {
							Type:     schema.TypeString,
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
						// "created": {
						// 	Type:     schema.TypeList,
						// 	Optional: true,
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"at": {
						// 				Type:     schema.TypeInt,
						// 				Optional: true,
						// 			},
						// 			"by": {
						// 				Type:     schema.TypeString,
						// 				Optional: true,
						// 			},
						// 		},
						// 	},
						// },
						// "modified": {
						// 	Type:     schema.TypeList,
						// 	Optional: true,
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"at": {
						// 				Type:     schema.TypeInt,
						// 				Optional: true,
						// 			},
						// 			"by": {
						// 				Type:     schema.TypeString,
						// 				Optional: true,
						// 			},
						// 		},
						// 	},
						// },
					},
				},
			},
		},
	}
}

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*alertlogic.API)
	// fmt.Println(api)

	var diags diag.Diagnostics

	users, err := api.ListUsers(false, false, false, "")
	if err != nil {
		return diag.FromErr(err)
	}

	usersJson, err := json.Marshal(users.Users)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Println(string(usersJson))

	// if err := d.Set("users", usersJson); err != nil {
	// 	return diag.FromErr(err)
	// }

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
