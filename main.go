package main

import (
	"context"
	"log"

	"time"

	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/marshallford/terraform-provider-pfsense/internal/provider"
)

//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

// version is set by GoReleaser at build time via ldflags.
var version = "dev"

func main() {
	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/emanuelbesliu/pfsense",
	})
	// Providers are killed by Terraform after each operation: flush any queued
	// bulk apply ops before the process leaves the room.
	pfsense.FlushAllBulk(context.Background(), 20*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}
