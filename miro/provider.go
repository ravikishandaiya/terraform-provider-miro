package miro

import (
	"fmt"
	"terraform-provider-miro/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"miro_token": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				DefaultFunc: schema.EnvDefaultFunc("MIRO_TOKEN", ""),
			}, 
			"miro_team_id": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				DefaultFunc: schema.EnvDefaultFunc("MIRO_TEAM_ID", ""),
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
	miroToken 	:= d.Get("miro_token").(string)
	miroTeam_id := d.Get("miro_team_id").(string)
	if len(miroToken) == 0 || len(miroTeam_id) == 0 {
		return client.NewClient(miroToken, miroTeam_id), fmt.Errorf("Token or Team ID is not provided.")
	}
	return client.NewClient(miroToken, miroTeam_id), nil
}
