package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/marshallford/terraform-provider-pfsense/internal/provider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"pfsense": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if os.Getenv("TF_PFSENSE_PASSWORD") == "" {
		t.Fatal("TF_PFSENSE_PASSWORD must be set for acceptance tests")
	}
}
