//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
package main

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	godaddy_provider "github.com/myklst/terraform-provider-st-godaddy/godaddy"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	providerAddress := os.Getenv("PROVIDER_LOCAL_PATH")
	if providerAddress == "" {
		providerAddress = "registry.terraform.io/myklst/st-godaddy"
	}

	providerserver.Serve(context.Background(), godaddy_provider.New, providerserver.ServeOpts{
		Address: providerAddress,
	})
}
