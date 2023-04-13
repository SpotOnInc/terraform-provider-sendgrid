/*
Provide a resource to manage SSO teammates.
Example Usage
```hcl
resource "sendgrid_sso_teammate" "teammate" {
	email       = "jane.doe@example.com"
	first_name  = "Jane"
	last_name   = "Doe"
	is_admin    = false
	persona     = "observer"
	scopes      = ["mail.send", "alerts.read"]
}
```
Import
An SSO teammate can be imported, e.g.
```sh
$ terraform import sendgrid_sso_teammate.teammate <email>
```
*/
package sendgrid

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	sendgrid "github.com/SpotOnInc/terraform-provider-sendgrid/sdk"
)

func resourceSendgridSSOTeammate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSendgridSSOTeammateCreate,
		ReadContext:   resourceSendgridSSOTeammateRead,
		UpdateContext: resourceSendgridSSOTeammateUpdate,
		DeleteContext: resourceSendgridSSOTeammateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
				Computed: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_admin": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"persona": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"accountant",
					"developer",
					"marketer",
					"observer",
				}, false),
			},
			"scopes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSendgridSSOTeammateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	email := d.Get("email").(string)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	isAdmin := d.Get("is_admin").(bool)
	persona := d.Get("persona").(string)
	scopes := ExpandStringList(d.Get("scopes").([]interface{}))

	teammateStruct, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.CreateSSOTeamMate(email, firstName, lastName, isAdmin, persona, scopes)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	teammate := teammateStruct.(*sendgrid.SSOTeammate)

	d.SetId(teammate.Email)

	return resourceSendgridSSOTeammateRead(ctx, d, m)
}

func resourceSendgridSSOTeammateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	email := d.Id()

	teammate, request := c.ReadSSOTeamMate(email)
	if request.Err != nil {
		return diag.FromErr(request.Err)
	}

	d.Set("email", teammate.Email)
	d.Set("first_name", teammate.FirstName)
	d.Set("last_name", teammate.LastName)
	d.Set("is_admin", teammate.IsAdmin)
	d.Set("persona", teammate.Persona)
	d.Set("scopes", teammate.Scopes)

	return nil
}

func resourceSendgridSSOTeammateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	email := d.Get("email").(string)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	isAdmin := d.Get("is_admin").(bool)
	persona := d.Get("persona").(string)
	scopes := ExpandStringList(d.Get("scopes").([]interface{}))

	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.UpdateSSOTeamMate(email, email, firstName, lastName, isAdmin, persona, scopes)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSendgridSSOTeammateRead(ctx, d, m)
}

func resourceSendgridSSOTeammateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*sendgrid.Client)

	email := d.Get("email").(string)

	_, err := sendgrid.RetryOnRateLimit(ctx, d, func() (interface{}, sendgrid.RequestError) {
		return c.DeleteSSOTeamMate(email)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ExpandStringList(d []interface{}) []string {
	vs := make([]string, 0, len(d))
	for _, v := range d {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}
