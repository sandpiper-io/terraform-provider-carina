package carina

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CARINA_USERNAME", nil),
				Description: "Username for Carina connection",
			},
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CARINA_APIKEY", nil),
				Description: "API key for Carina connection",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"carina_cluster": resourceCarinaCluster(),
		},

		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Username: d.Get("username").(string),
		ApiKey: d.Get("api_key").(string),
	}

	client, err := config.NewClient()
	if err != nil {
		return nil, fmt.Errorf("Error initializing Carina client: %s", err)
	}

	return client, nil
}