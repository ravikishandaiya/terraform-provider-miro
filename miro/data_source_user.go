package miro

import (
	
	"strings"
	"terraform-provider-miro/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:        schema.TypeString, 
				Required:    true,
			},
			"role": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type":	&schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:	 true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:	 true,
			},
			"team_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:	 true,
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:	 true,
			},
			"company": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:	 true,
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:	 true,
			},
			"image_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:	 true,
			},
		},
	}	
}

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	email     := d.Get("email").(string)
	resp, err := apiClient.GetUser(email)
	if err != nil {
		if strings.Contains(err.Error(), "User Not Found") {
			d.SetId("")
		} else {
			return err
		}
	}
	d.SetId(resp.Email)
	d.Set("type",resp.Type)
	d.Set("name",resp.Name)
	d.Set("team_name",resp.TeamName)
	d.Set("created_at",resp.CreatedAt)
	d.Set("role",resp.Role)
	d.Set("state",resp.State)
	return nil
}