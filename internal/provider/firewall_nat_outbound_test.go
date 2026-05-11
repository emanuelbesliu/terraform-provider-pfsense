package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallNATOutboundResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNATOutboundConfig("wan", "tcp", "192.168.1.0/24", "", "any", "", "(self)", "tf-acc-test nat out basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_nat_outbound.test", "interface", "wan"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_outbound.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_outbound.test", "source_address", "192.168.1.0/24"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_outbound.test", "destination_address", "any"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_outbound.test", "target", "(self)"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_outbound.test", "description", "tf-acc-test nat out basic"),
				),
			},
			{
				ResourceName:      "pfsense_firewall_nat_outbound.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFirewallNATOutboundConfig("wan", "udp", "10.0.0.0/8", "", "any", "", "(self)", "tf-acc-test nat out updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_nat_outbound.test", "protocol", "udp"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_outbound.test", "source_address", "10.0.0.0/8"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_outbound.test", "description", "tf-acc-test nat out updated"),
				),
			},
		},
	})
}

func TestAccFirewallNATOutboundResource_disabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNATOutboundConfig_disabled("wan", "", "192.168.2.0/24", "any", "(self)", "tf-acc-test nat out disabled"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_nat_outbound.test_disabled", "disabled", "true"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_outbound.test_disabled", "description", "tf-acc-test nat out disabled"),
				),
			},
		},
	})
}

func TestAccFirewallNATOutboundRulesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNATOutboundConfig("wan", "tcp", "192.168.3.0/24", "", "any", "", "(self)", "tf-acc-test nat out ds") + `
data "pfsense_firewall_nat_outbound_rules" "all" {
  depends_on = [pfsense_firewall_nat_outbound.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.pfsense_firewall_nat_outbound_rules.all", "mode"),
					resource.TestCheckResourceAttrSet("data.pfsense_firewall_nat_outbound_rules.all", "rules.#"),
				),
			},
		},
	})
}

func TestAccFirewallNATOutboundRuleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNATOutboundConfig("wan", "tcp", "192.168.4.0/24", "", "any", "", "(self)", "tf-acc-test nat out single ds") + `
data "pfsense_firewall_nat_outbound" "test" {
  description = "tf-acc-test nat out single ds"
  depends_on  = [pfsense_firewall_nat_outbound.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.pfsense_firewall_nat_outbound.test", "interface", "wan"),
					resource.TestCheckResourceAttr("data.pfsense_firewall_nat_outbound.test", "source_address", "192.168.4.0/24"),
					resource.TestCheckResourceAttr("data.pfsense_firewall_nat_outbound.test", "description", "tf-acc-test nat out single ds"),
				),
			},
		},
	})
}

func testAccFirewallNATOutboundConfig(iface, proto, srcAddr, srcPort, destAddr, destPort, target, description string) string {
	return fmt.Sprintf(`
resource "pfsense_firewall_nat_outbound" "test" {
  interface           = %q
  protocol            = %q
  source_address      = %q
  source_port         = %q
  destination_address = %q
  destination_port    = %q
  target              = %q
  description         = %q
}
`, iface, proto, srcAddr, srcPort, destAddr, destPort, target, description)
}

func testAccFirewallNATOutboundConfig_disabled(iface, proto, srcAddr, destAddr, target, description string) string {
	return fmt.Sprintf(`
resource "pfsense_firewall_nat_outbound" "test_disabled" {
  interface           = %q
  protocol            = %q
  source_address      = %q
  destination_address = %q
  target              = %q
  description         = %q
  disabled            = true
}
`, iface, proto, srcAddr, destAddr, target, description)
}
