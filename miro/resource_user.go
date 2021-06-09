package miro

import (
	"log"
	"context"
	"terraform-provider-miro/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

func resourceUserCreate(ctx context.Context,d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient 	:= m.(*client.Client)
	email 		:= d.Get("email").(string)
	err 		:= apiClient.CreateUser(email)
	if err != nil {
		log.Println("[ERROR]: ",err)
		return diag.FromErr(err)
	}
	d.SetId(email)
	resourceUserRead(ctx,d,m)
	return diags
}

func resourceUserRead(ctx context.Context,d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient 	:= m.(*client.Client)
	email 		:= d.Id()
	resp, err 	:= apiClient.GetUser(email)
	if err != nil {
		log.Println("[Error]: ",err)
	} else {
		d.SetId(resp.Email)
		d.Set("type",resp.Type)
		d.Set("email",resp.Email)
		d.Set("name",resp.Name)
		d.Set("team_name",resp.TeamName)
		d.Set("created_at",resp.CreatedAt)
		d.Set("role",resp.Role)
		d.Set("state",resp.State)
	}
	return diags
}

func resourceUserUpdate(ctx context.Context,d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	if d.HasChange("email") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "User not allowed to change email",
			Detail:   "User not allowed to change email",
		})
		return diags
	}
	email := d.Get("email").(string)
	role := d.Get("role").(string)
	err := apiClient.UpdateUser(email, role)
	if err != nil {
		log.Printf("[Error] Error updating user : %s", err)
		return diag.FromErr(err)
	}
	return resourceUserRead(ctx,d,m)
}

func resourceUserDelete(ctx context.Context,d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	if d.HasChange("email") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "User not allowed to change email",
			Detail:   "User not allowed to change email",
		})
		return diags
	}
	email := d.Id()
	err := apiClient.DeleteUser(email)
	if err != nil {
		log.Println("[Error]: ", err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}