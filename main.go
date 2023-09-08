package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/thiagoarrais/terraform-provider-sha256sum/internal/sha256sum"
)

var (
	// Example version string that can be overwritten by a release process
	version string = "dev"
)

func main() {
	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/thiagoarrais/sha256sum",
	}

	err := providerserver.Serve(context.Background(), sha256sum.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
