package openfaas

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"

	openfaas_config "github.com/openfaas/faas-cli/config"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	log.Printf("[DEBUG] returning provider schema")
	// The actual provider
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "http://localhost:8080",
				Description: "OpenFaaS gateway uri",
			},
			"tls_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "if true, skip tls verification (not recommended)",
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "OpenFaaS gateway username",
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Default:     "",
				Description: "OpenFaaS gateway password",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"openfaas_function": dataSourceOpenFaaSFunction(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"openfaas_function": resourceOpenFaaSFunction(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Printf("[DEBUG] configuring provider")
	config := Config{
		GatewayURI:      d.Get("uri").(string),
		TLSInsecure:     d.Get("tls_insecure").(bool),
		GatewayUserName: d.Get("user_name").(string),
		GatewayPassword: d.Get("password").(string),
	}

	if config.GatewayUserName != "" && config.GatewayPassword != "" {
		openfaas_config.UpdateAuthConfig(config.GatewayURI, config.GatewayUserName, config.GatewayPassword)
	}

	return config, nil
}
