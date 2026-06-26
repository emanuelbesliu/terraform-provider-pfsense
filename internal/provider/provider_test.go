package provider_test

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/marshallford/terraform-provider-pfsense/internal/provider"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
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

// testAccNewClient builds a pfSense API client from the same environment
// variables used by the provider. It is intended for use in CheckDestroy and
// other out-of-band verification helpers during acceptance tests.
func testAccNewClient() (*pfsense.Client, error) {
	opts := pfsense.Options{
		Password: os.Getenv("TF_PFSENSE_PASSWORD"),
	}

	urlValue := os.Getenv("TF_PFSENSE_URL")
	if urlValue == "" {
		urlValue = pfsense.DefaultURL
	}

	parsedURL, err := url.Parse(urlValue)
	if err != nil {
		return nil, fmt.Errorf("parsing pfSense URL: %w", err)
	}
	opts.URL = parsedURL

	if username := os.Getenv("TF_PFSENSE_USERNAME"); username != "" {
		opts.Username = username
	}

	return pfsense.NewClient(context.Background(), &opts)
}
