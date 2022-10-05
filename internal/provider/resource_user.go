package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				// This description is used by the documentation generator and the language server.
				Description: "Name of user",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}

}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	userName := d.Get("name").(string)
	config := meta.(*Config)

	conn, err := Connect(config)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	sql := fmt.Sprintf("CREATE USER %s", userName)
	err = conn.Exec(ctx, sql)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	d.SetId(userName)

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "created a resource")

	return nil
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	userName := d.Id()
	tflog.Trace(ctx, "test user read create config")
	config := meta.(*Config)

	tflog.Trace(ctx, "test user read connecting")
	conn, err := Connect(config)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	var result []struct {
		Name string
	}

	sql := "SELECT name as Name FROM system.users WHERE name = $1"

	tflog.Trace(ctx, fmt.Sprintf("test user select user. query = %s", sql))
	err = conn.Select(ctx, &result, sql, userName)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	if len(result) != 1 {
		return diag.Errorf(fmt.Sprintf("User '%s' Does not exist", userName))
	}

	tflog.Trace(ctx, "test user set id")
	d.SetId(userName)
	return nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	oldUserName := d.Id()

	_, newUserName_ := d.GetChange("name")

	newUserName := newUserName_.(string)

	config := meta.(*Config)
	conn, err := Connect(config)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	err = conn.Exec(ctx, "ALTER USER $1 RENAME TO $2", oldUserName, newUserName)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	d.SetId(newUserName)

	return nil
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	userName := d.Get("name").(string)
	config := meta.(*Config)

	conn, err := Connect(config)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	sql := fmt.Sprintf("DROP USER %s", userName)
	err = conn.Exec(ctx, sql)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	// write logs using the tflog package
	// see https://pkg.go.dev/github.com/hashicorp/terraform-plugin-log/tflog
	// for more information
	tflog.Trace(ctx, "deleted a resource")

	return nil
}
