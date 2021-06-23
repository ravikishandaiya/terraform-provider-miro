package miro

import (
	"fmt"
	"time"
	"regexp"
	"strings"
	"context"
	"terraform-provider-miro/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func validateEmail(v interface{}, email string) (s []string, errs []error) {
	value := v.(string)
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !(emailRegex.MatchString(value)) {
		errs = append(errs, fmt.Errorf("Expected Email Id is not valid %s", email))
		return s,errs
	}
	return
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceUserImporter,
		},
		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:        schema.TypeString, 
				Required:    true,
				ValidateFunc: validateEmail,
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
			"industry":  &schema.Schema{
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
	var err error
	retryErr := resource.Retry(2*time.Second, func() *resource.RetryError {
		if err = apiClient.CreateUser(email); err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	resourceUserRead(ctx,d,m)
	return diags
}

func resourceUserRead(ctx context.Context,d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient 	:= m.(*client.Client)
	email 		:= d.Get("email").(string)
	retryErr := resource.Retry(2*time.Second, func() *resource.RetryError {
		resp, err := apiClient.GetUser(email)
		if err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.SetId(resp.Email)
		d.Set("type",resp.Type)
		d.Set("email",resp.Email)
		d.Set("name",resp.Name)
		d.Set("team_name",resp.TeamName)
		d.Set("created_at",resp.CreatedAt)
		d.Set("industry",resp.Industry)
		d.Set("company",resp.Company)
		d.Set("role",resp.Role)
		d.Set("state",resp.State)
		return nil
	})
	if retryErr!=nil {
		if strings.Contains(retryErr.Error(), "User Not Found")==true {
			d.SetId("")
			return diags
		}
		return diag.FromErr(retryErr)
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
	var err error
	retryErr := resource.Retry(2*time.Second, func() *resource.RetryError {
		if err = apiClient.UpdateUser(email, role); err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	if err != nil {
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
	var err error
	retryErr := resource.Retry(2*time.Second, func() *resource.RetryError {
		if err = apiClient.DeleteUser(email); err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceUserImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData ,error) {
	apiClient := m.(*client.Client)
	email     := d.Id()
	resp, err := apiClient.GetUser(email)
	if err != nil {
		if strings.Contains(err.Error(), "User Not Found") {
			d.SetId("")
			return nil, err
		}
	}
	d.SetId(resp.Email)
	d.Set("type",resp.Type)
	d.Set("name",resp.Name)
	d.Set("team_name",resp.TeamName)
	d.Set("created_at",resp.CreatedAt)
	d.Set("industry",resp.Industry)
	d.Set("company",resp.Company)
	d.Set("role",resp.Role)
	d.Set("state",resp.State)
	return []*schema.ResourceData{d}, nil
}
