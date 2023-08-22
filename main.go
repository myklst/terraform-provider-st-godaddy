//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"os"
	godaddy_provider "terraform-provider-st-godaddy/plugin/godaddy"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"

	// goreleaser can also pass the specific commit if you want
	// commit  string = ""
)

func main() {
	providerAddress := os.Getenv("PROVIDER_LOCAL_PATH")
	if providerAddress == "" {
		providerAddress = "registry.terraform.io/myklst/st-godaddy"
	}

	providerserver.Serve(context.Background(), godaddy_provider.New, providerserver.ServeOpts{
		Address: providerAddress,
	})
}
