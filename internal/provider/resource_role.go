package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"sort"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "User Role",

		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		//UpdateContext: resourceUserUpdate,
		//DeleteContext: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				// This description is used by the documentation generator and the language server.
				Description: "Name of role",
				Type:        schema.TypeString,
				Required:    true,
			},
			"permissions": {
				// This description is used by the documentation generator and the language server.
				Description: "Name of role",
				Type:        schema.TypeList,
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}

}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)

	roleName := d.Get("name").(string)
	permissions := d.Get("permissions").([]string)

	config := meta.(*Config)

	conn, err := Connect(config)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	sql := fmt.Sprintf("CREATE ROLE %s", roleName)
	err = conn.Exec(ctx, sql)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	tflog.Trace(ctx, "Role '%s' was created")

	sort.Strings(permissions)

	for _, permissionSql := range permissions {
		sql := fmt.Sprintf("GRANT %s TO %s;", permissionSql, roleName)
		err = conn.Exec(ctx, sql)

		if err != nil {
			return diag.Errorf(err.Error())
		}

		tflog.Trace(ctx, fmt.Sprintf("Role Grant '%s' to '%s' was created", permissionSql, roleName))
	}

	d.SetId(roleName)

	tflog.Trace(ctx, "created a resource")

	return nil
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	// client := meta.(*apiClient)
	userRole := d.Id()

	config := meta.(*Config)

	tflog.Trace(ctx, "test user read connecting")
	conn, err := Connect(config)

	if err != nil {
		return diag.Errorf(err.Error())
	}

	var result []struct {
		RoleName string,
		AccessType string,
		Database string,
		Table string,
		Column string,
	}

	sql := "SELECT " +
			"role_name as RoleName, " +
			"access_type as AccessType, " +
			"database as Database, " +
			"table as Table, " +
			"columnt as Column " +
			"FROM system.grants where role_name = $1"
		)
	err = conn.Select(ctx, &result, sql, userRole)

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
