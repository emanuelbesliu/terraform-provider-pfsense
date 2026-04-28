package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallVirtualIPResource_ipalias(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallVirtualIPConfig_ipalias("10.0.161.200", 32, "lan", "tf-acc-test ipalias"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_virtual_ip.test", "mode", "ipalias"),
					resource.TestCheckResourceAttr("pfsense_firewall_virtual_ip.test", "interface", "lan"),
					resource.TestCheckResourceAttr("pfsense_firewall_virtual_ip.test", "subnet", "10.0.161.200"),
					resource.TestCheckResourceAttr("pfsense_firewall_virtual_ip.test", "subnet_bits", "32"),
					resource.TestCheckResourceAttr("pfsense_firewall_virtual_ip.test", "description", "tf-acc-test ipalias"),
					resource.TestCheckResourceAttrSet("pfsense_firewall_virtual_ip.test", "unique_id"),
				),
			},
			{
				ResourceName:      "pfsense_firewall_virtual_ip.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFirewallVirtualIPConfig_ipalias("10.0.161.201", 32, "lan", "tf-acc-test ipalias updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_virtual_ip.test", "subnet", "10.0.161.201"),
					resource.TestCheckResourceAttr("pfsense_firewall_virtual_ip.test", "description", "tf-acc-test ipalias updated"),
				),
			},
		},
	})
}

func TestAccFirewallVirtualIPResource_carp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallVirtualIPConfig_carp("10.0.161.210", 32, "lan", 10, 0, 1, "carppass", "tf-acc-test carp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_virtual_ip.test_carp", "mode", "carp"),
					resource.TestCheckResourceAttr("pfsense_firewall_virtual_ip.test_carp", "interface", "lan"),
					resource.TestCheckResourceAttr("pfsense_firewall_virtual_ip.test_carp", "subnet", "10.0.161.210"),
					resource.TestCheckResourceAttr("pfsense_firewall_virtual_ip.test_carp", "vhid", "10"),
					resource.TestCheckResourceAttr("pfsense_firewall_virtual_ip.test_carp", "advskew", "0"),
					resource.TestCheckResourceAttr("pfsense_firewall_virtual_ip.test_carp", "advbase", "1"),
					resource.TestCheckResourceAttrSet("pfsense_firewall_virtual_ip.test_carp", "unique_id"),
				),
			},
			{
				ResourceName:            "pfsense_firewall_virtual_ip.test_carp",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccFirewallVirtualIPsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallVirtualIPConfig_ipalias("10.0.161.220", 32, "lan", "tf-acc-test ds") + `
data "pfsense_firewall_virtual_ips" "all" {
  depends_on = [pfsense_firewall_virtual_ip.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.pfsense_firewall_virtual_ips.all", "virtual_ips.#"),
				),
			},
		},
	})
}

func testAccFirewallVirtualIPConfig_ipalias(subnet string, subnetBits int, iface string, description string) string {
	return fmt.Sprintf(`
resource "pfsense_firewall_virtual_ip" "test" {
  mode        = "ipalias"
  interface   = %q
  subnet      = %q
  subnet_bits = %d
  description = %q
}
`, iface, subnet, subnetBits, description)
}

func testAccFirewallVirtualIPConfig_carp(subnet string, subnetBits int, iface string, vhid int, advskew int, advbase int, password string, description string) string {
	return fmt.Sprintf(`
resource "pfsense_firewall_virtual_ip" "test_carp" {
  mode        = "carp"
  interface   = %q
  subnet      = %q
  subnet_bits = %d
  vhid        = %d
  advskew     = %d
  advbase     = %d
  password    = %q
  description = %q
}
`, iface, subnet, subnetBits, vhid, advskew, advbase, password, description)
}
