package netcup

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETCUP_WS_USER", nil),
				Description: "The user login name for API operations.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NETCUP_WS_PASSWORD", nil),
				Description: "The login password for API operations.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"netcup_vserver": dataSourceVServer(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

type ProviderConfig struct {
	LoginName string
	Password  string
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return ProviderConfig{
		LoginName: d.Get("login_name").(string),
		Password:  d.Get("password").(string),
	}, nil
}
