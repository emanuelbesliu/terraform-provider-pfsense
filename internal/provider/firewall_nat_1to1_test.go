package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallNAT1to1Resource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccFirewallNAT1to1ResourceConfig("test-1to1-nat"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_nat_1to1.test", "external", "203.0.113.10"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_1to1.test", "interface", "wan"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_1to1.test", "ipprotocol", "inet"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_1to1.test", "source_address", "any"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_1to1.test", "destination_address", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_1to1.test", "description", "test-1to1-nat"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_1to1.test", "disabled", "false"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_1to1.test", "no_binat", "false"),
				),
			},
			// Update and Read testing
			{
				Config: testAccFirewallNAT1to1ResourceConfig("test-1to1-nat-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_nat_1to1.test", "description", "test-1to1-nat-updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccFirewallNAT1to1ResourceConfig(description string) string {
	return fmt.Sprintf(`
resource "pfsense_firewall_nat_1to1" "test" {
  external              = "203.0.113.10"
  interface             = "wan"
  ipprotocol            = "inet"
  source_address        = "any"
  destination_address   = "192.168.1.0/24"
  description           = "%s"
  disabled              = false
  no_binat              = false
  nat_reflection        = ""
}
`, description)
}
