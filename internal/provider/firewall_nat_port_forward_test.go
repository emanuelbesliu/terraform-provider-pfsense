package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallNATPortForwardResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNATPortForwardConfig("wan", "inet", "tcp", "wanip", "8080", "10.0.161.50", "80", "tf-acc-test nat pf basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test", "interface", "wan"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test", "ipprotocol", "inet"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test", "destination_address", "wanip"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test", "destination_port", "8080"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test", "target", "10.0.161.50"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test", "local_port", "80"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test", "description", "tf-acc-test nat pf basic"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test", "associated_rule_id", "pass"),
				),
			},
			{
				ResourceName:      "pfsense_firewall_nat_port_forward.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFirewallNATPortForwardConfig("wan", "inet", "tcp", "wanip", "8081", "10.0.161.51", "81", "tf-acc-test nat pf updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test", "destination_port", "8081"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test", "target", "10.0.161.51"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test", "local_port", "81"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test", "description", "tf-acc-test nat pf updated"),
				),
			},
		},
	})
}

func TestAccFirewallNATPortForwardResource_disabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNATPortForwardConfig_disabled("wan", "inet", "tcp", "wanip", "9090", "10.0.161.52", "90", "tf-acc-test nat pf disabled"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test_disabled", "disabled", "true"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_port_forward.test_disabled", "description", "tf-acc-test nat pf disabled"),
				),
			},
		},
	})
}

func TestAccFirewallNATPortForwardsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNATPortForwardConfig("wan", "inet", "tcp", "wanip", "7070", "10.0.161.53", "70", "tf-acc-test nat pf ds") + `
data "pfsense_firewall_nat_port_forwards" "all" {
  depends_on = [pfsense_firewall_nat_port_forward.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.pfsense_firewall_nat_port_forwards.all", "nat_port_forwards.#"),
				),
			},
		},
	})
}

func TestAccFirewallNATPortForwardDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNATPortForwardConfig("wan", "inet", "tcp", "wanip", "6060", "10.0.161.54", "60", "tf-acc-test nat pf single ds") + `
data "pfsense_firewall_nat_port_forward" "test" {
  description = "tf-acc-test nat pf single ds"
  depends_on  = [pfsense_firewall_nat_port_forward.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.pfsense_firewall_nat_port_forward.test", "interface", "wan"),
					resource.TestCheckResourceAttr("data.pfsense_firewall_nat_port_forward.test", "target", "10.0.161.54"),
					resource.TestCheckResourceAttr("data.pfsense_firewall_nat_port_forward.test", "description", "tf-acc-test nat pf single ds"),
				),
			},
		},
	})
}

func testAccFirewallNATPortForwardConfig(iface, ipproto, proto, destAddr, destPort, target, localPort, description string) string {
	return fmt.Sprintf(`
resource "pfsense_firewall_nat_port_forward" "test" {
  interface           = %q
  ipprotocol          = %q
  protocol            = %q
  destination_address = %q
  destination_port    = %q
  target              = %q
  local_port          = %q
  description         = %q
  associated_rule_id  = "pass"
}
`, iface, ipproto, proto, destAddr, destPort, target, localPort, description)
}

func testAccFirewallNATPortForwardConfig_disabled(iface, ipproto, proto, destAddr, destPort, target, localPort, description string) string {
	return fmt.Sprintf(`
resource "pfsense_firewall_nat_port_forward" "test_disabled" {
  interface           = %q
  ipprotocol          = %q
  protocol            = %q
  destination_address = %q
  destination_port    = %q
  target              = %q
  local_port          = %q
  description         = %q
  disabled            = true
  associated_rule_id  = "pass"
}
`, iface, ipproto, proto, destAddr, destPort, target, localPort, description)
}
