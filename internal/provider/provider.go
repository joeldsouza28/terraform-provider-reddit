package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Reddit Bot Client Id",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Reddit Bot Client Secret",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Reddit UserName",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Reddit Password",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"reddit_post":    resourcePost(),
			"reddit_comment": resourceComment(),
		},
		ConfigureContextFunc: func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			config := map[string]string{
				"client_id":     d.Get("client_id").(string),
				"client_secret": d.Get("client_secret").(string),
				"username":      d.Get("username").(string),
				"password":      d.Get("password").(string),
			}
			return config, nil
		},
	}
}
