package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sandpiper-io/terraform-provider-carina/carina"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: carina.Provider,
	})
}