package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"database": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("HASHICUPS_DATABASE", nil),
				},
				"host": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   false,
					DefaultFunc: schema.EnvDefaultFunc("HASHICUPS_HOST", nil),
				},
				"port": &schema.Schema{
					Type:        schema.TypeInt,
					Required:    true,
					Sensitive:   false,
					DefaultFunc: schema.EnvDefaultFunc("HASHICUPS_PORT", nil),
				},
				"username": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   false,
					DefaultFunc: schema.EnvDefaultFunc("HASHICUPS_USERNAME", nil),
				},
				"password": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   false,
					DefaultFunc: schema.EnvDefaultFunc("HASHICUPS_PASSWORD", nil),
				},
				"timeout": &schema.Schema{
					Type:        schema.TypeInt,
					Required:    true,
					Sensitive:   false,
					DefaultFunc: schema.EnvDefaultFunc("HASHICUPS_TIMEOUT", nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"scaffolding_data_source": dataSourceScaffolding(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"clickhouse_user": resourceUser(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(cnx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		config := Config{
			d.Get("database").(string),
			d.Get("host").(string),
			d.Get("port").(int),
			d.Get("username").(string),
			d.Get("password").(string),
			d.Get("timeout").(int),
		}

		return &config, nil
	}
}
