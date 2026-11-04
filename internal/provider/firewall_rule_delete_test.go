package provider_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccCheckFirewallRuleDestroy confirms that every rule created during the
// test run - by tracker and by description - is really gone from pfSense
// afterwards, on any interface and in the floating ("any") section alike.
func testAccCheckFirewallRuleDestroy(s *terraform.State) error {
	client, err := testAccNewClient()
	if err != nil {
		return err
	}

	rules, err := client.GetFirewallRules(context.Background())
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pfsense_firewall_rule" {
			continue
		}

		description := rs.Primary.Attributes["description"]
		tracker := rs.Primary.Attributes["tracker"]

		for _, r := range *rules {
			if (description != "" && r.Description == description) ||
				(tracker != "" && r.Tracker == tracker) {
				return fmt.Errorf("firewall rule (tracker %s, description %q) still exists after destroy", tracker, description)
			}
		}
	}

	return nil
}

// testAccCheckFirewallRuleAbsent verifies out-of-band that no rule with the
// given description exists anywhere on the firewall. This is the point of the
// delete-by-position regression: a provider deleting by raw config array key
// removes neighbouring rules by mistake and silently leaves the "deleted"
// rule behind - which this check catches.
func testAccCheckFirewallRuleAbsent(description string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		client, err := testAccNewClient()
		if err != nil {
			return err
		}

		rules, err := client.GetFirewallRules(context.Background())
		if err != nil {
			return err
		}

		for _, r := range *rules {
			if r.Description == description {
				return fmt.Errorf("rule with description %q still exists on interface %q (tracker %s)", description, r.Interface, r.Tracker)
			}
		}

		return nil
	}
}

// testAccCheckFirewallRulePresent verifies out-of-band that a rule with the
// given description still exists on the firewall.
func testAccCheckFirewallRulePresent(description string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		client, err := testAccNewClient()
		if err != nil {
			return err
		}

		rules, err := client.GetFirewallRules(context.Background())
		if err != nil {
			return err
		}

		for _, r := range *rules {
			if r.Description == description {
				return nil
			}
		}

		return fmt.Errorf("rule with description %q is missing", description)
	}
}

// TestAccFirewallRule_deletePositionRegression reproduces the wrong-rule
// deletion class: one apply destroys two of three rules on the same
// interface. Deleting by the raw config.xml array index (pre-fix behaviour)
// hits the wrong positions once an earlier removal shifts the list;
// deleting by the GUI list position (fix) removes exactly what was planned.
func TestAccFirewallRule_deletePositionRegression(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckFirewallRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfig_deleteMulti(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_rule.del_a", "description", "tf-acc-del-a"),
					resource.TestCheckResourceAttr("pfsense_firewall_rule.del_b", "description", "tf-acc-del-b"),
					resource.TestCheckResourceAttr("pfsense_firewall_rule.del_c", "description", "tf-acc-del-c"),
					testAccCheckFirewallRulePresent("tf-acc-del-a"),
					testAccCheckFirewallRulePresent("tf-acc-del-b"),
					testAccCheckFirewallRulePresent("tf-acc-del-c"),
				),
			},
			{
				// Single apply, two removals: the exact conditions of the
				// original position-shift bug.
				Config: testAccFirewallRuleConfig_deleteSingle(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFirewallRuleAbsent("tf-acc-del-a"),
					testAccCheckFirewallRuleAbsent("tf-acc-del-b"),
					testAccCheckFirewallRulePresent("tf-acc-del-c"),
					resource.TestCheckResourceAttr("pfsense_firewall_rule.del_c", "description", "tf-acc-del-c"),
				),
			},
		},
	})
}

// accProviderConfig configures the provider explicitly. The lab appliances
// use self-signed certificates, so tls_skip_verify is enabled here; URL and
// credentials still come from the TF_PFSENSE_* environment variables.
func accProviderConfig() string {
	return "provider \"pfsense\" {\n  tls_skip_verify = true\n}\n\n"
}

func testAccFirewallRuleConfig_deleteMulti() string {
	return accProviderConfig() + `
resource "pfsense_firewall_rule" "del_a" {
  type        = "block"
  interface   = "wan"
  description = "tf-acc-del-a"
}

resource "pfsense_firewall_rule" "del_b" {
  type        = "block"
  interface   = "wan"
  description = "tf-acc-del-b"
}

resource "pfsense_firewall_rule" "del_c" {
  type        = "block"
  interface   = "wan"
  description = "tf-acc-del-c"
}
`
}

func testAccFirewallRuleConfig_deleteSingle() string {
	return accProviderConfig() + `
resource "pfsense_firewall_rule" "del_c" {
  type        = "block"
  interface   = "wan"
  description = "tf-acc-del-c"
}
`
}
