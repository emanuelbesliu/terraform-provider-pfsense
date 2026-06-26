package provider_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

func testAccCheckFirewallNATNPtDestroy(s *terraform.State) error {
	client, err := testAccNewClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pfsense_firewall_nat_npt" {
			continue
		}

		description := rs.Primary.Attributes["description"]

		_, err := client.GetNATNPt(context.Background(), description)
		if err == nil {
			return fmt.Errorf("NAT NPt rule %q still exists", description)
		}

		if !errors.Is(err, pfsense.ErrNotFound) {
			return err
		}
	}

	return nil
}

func TestAccFirewallNATNPtResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckFirewallNATNPtDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNATNPtConfig("wan", "fd00:1::/64", "2001:db8:1::/64", "tf-acc-test npt basic"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test", "interface", "wan"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test", "source_prefix", "fd00:1::/64"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test", "destination_prefix", "2001:db8:1::/64"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test", "source_not", "false"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test", "destination_not", "false"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test", "disabled", "false"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test", "description", "tf-acc-test npt basic"),
				),
			},
			{
				ResourceName:      "pfsense_firewall_nat_npt.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFirewallNATNPtConfig("wan", "fd00:2::/64", "2001:db8:2::/64", "tf-acc-test npt updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test", "source_prefix", "fd00:2::/64"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test", "destination_prefix", "2001:db8:2::/64"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test", "description", "tf-acc-test npt updated"),
				),
			},
		},
	})
}

func TestAccFirewallNATNPtResource_negated(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckFirewallNATNPtDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNATNPtConfig_negated("wan", "fd00:3::/64", "2001:db8:3::/64", "tf-acc-test npt negated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test_negated", "source_not", "true"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test_negated", "destination_not", "true"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test_negated", "description", "tf-acc-test npt negated"),
				),
			},
			{
				ResourceName:      "pfsense_firewall_nat_npt.test_negated",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFirewallNATNPtResource_disabled(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckFirewallNATNPtDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNATNPtConfig_disabled("wan", "fd00:4::/64", "2001:db8:4::/64", "tf-acc-test npt disabled"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test_disabled", "disabled", "true"),
					resource.TestCheckResourceAttr("pfsense_firewall_nat_npt.test_disabled", "description", "tf-acc-test npt disabled"),
				),
			},
		},
	})
}

func TestAccFirewallNATNPtRulesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckFirewallNATNPtDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNATNPtConfig("wan", "fd00:5::/64", "2001:db8:5::/64", "tf-acc-test npt ds") + `
data "pfsense_firewall_nat_npt_rules" "all" {
  depends_on = [pfsense_firewall_nat_npt.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.pfsense_firewall_nat_npt_rules.all", "rules.#"),
				),
			},
		},
	})
}

func TestAccFirewallNATNPtRuleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckFirewallNATNPtDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallNATNPtConfig("wan", "fd00:6::/64", "2001:db8:6::/64", "tf-acc-test npt single ds") + `
data "pfsense_firewall_nat_npt" "test" {
  description = "tf-acc-test npt single ds"
  depends_on  = [pfsense_firewall_nat_npt.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.pfsense_firewall_nat_npt.test", "interface", "wan"),
					resource.TestCheckResourceAttr("data.pfsense_firewall_nat_npt.test", "source_prefix", "fd00:6::/64"),
					resource.TestCheckResourceAttr("data.pfsense_firewall_nat_npt.test", "destination_prefix", "2001:db8:6::/64"),
					resource.TestCheckResourceAttr("data.pfsense_firewall_nat_npt.test", "description", "tf-acc-test npt single ds"),
				),
			},
		},
	})
}

func testAccFirewallNATNPtConfig(iface, srcPrefix, destPrefix, description string) string {
	return fmt.Sprintf(`
resource "pfsense_firewall_nat_npt" "test" {
  interface          = %q
  source_prefix      = %q
  destination_prefix = %q
  description        = %q
}
`, iface, srcPrefix, destPrefix, description)
}

func testAccFirewallNATNPtConfig_negated(iface, srcPrefix, destPrefix, description string) string {
	return fmt.Sprintf(`
resource "pfsense_firewall_nat_npt" "test_negated" {
  interface          = %q
  source_prefix      = %q
  source_not         = true
  destination_prefix = %q
  destination_not    = true
  description        = %q
}
`, iface, srcPrefix, destPrefix, description)
}

func testAccFirewallNATNPtConfig_disabled(iface, srcPrefix, destPrefix, description string) string {
	return fmt.Sprintf(`
resource "pfsense_firewall_nat_npt" "test_disabled" {
  interface          = %q
  source_prefix      = %q
  destination_prefix = %q
  description        = %q
  disabled           = true
}
`, iface, srcPrefix, destPrefix, description)
}
