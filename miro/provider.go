package miro

import (
	"terraform-provider-miro/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			}, 
			"team_id": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"miro_user": resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"miro_user": dataSourceUser(),
		},
		ConfigureFunc:  providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	token 	:= d.Get("token").(string)
	team_id := d.Get("team_id").(string)
	return client.NewClient(token, team_id), nil
}
