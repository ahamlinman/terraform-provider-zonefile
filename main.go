package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/ahamlinman/terraform-provider-zonefile/internal/provider"
)

// version is set automatically during goreleaser builds.
var version string = "dev"

var flagDebug = flag.Bool("debug", false, "Set to true to run the provider with support for debuggers like delve")

func main() {
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/ahamlinman/zonefile",
		Debug:   *flagDebug,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
