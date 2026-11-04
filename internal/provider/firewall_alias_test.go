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

// testAccCheckIPAliasAbsent verifies out-of-band that the named IP alias is
// gone from pfSense. A failure here after a combined destroy is the
// signature of the parallel-destroy race: independent aliases destroyed in
// one apply where the last-writer-wins config write silently keeps one of
// them (the delete skipped because the control position had shifted).
func testAccCheckIPAliasAbsent(name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		client, err := testAccNewClient()
		if err != nil {
			return err
		}

		_, err = client.GetFirewallIPAlias(context.Background(), name)
		if err == nil {
			return fmt.Errorf("alias %q still exists on the firewall - parallel destroy left residue", name)
		}

		if !errors.Is(err, pfsense.ErrNotFound) {
			return err
		}

		return nil
	}
}

// testAccCheckIPAliasPresent verifies out-of-band that the named IP alias
// exists on the firewall.
func testAccCheckIPAliasPresent(name string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		client, err := testAccNewClient()
		if err != nil {
			return err
		}

		_, err = client.GetFirewallIPAlias(context.Background(), name)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckFirewallIPAliasDestroy(s *terraform.State) error {
	client, err := testAccNewClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pfsense_firewall_ip_alias" {
			continue
		}

		name := rs.Primary.Attributes["name"]

		_, err := client.GetFirewallIPAlias(context.Background(), name)
		if err == nil {
			return fmt.Errorf("IP alias %q still exists after destroy", name)
		}

		if !errors.Is(err, pfsense.ErrNotFound) {
			return err
		}
	}

	return nil
}

// TestAccFirewallAlias_parallelDestroy targets the historic race: three
// independent aliases are created, then a single apply removes two of them.
// Terraform destroys independent resources concurrently; without the
// control-ID retry the concurrent form deletes raced and one alias
// routinely survived with an empty delete log. Run with
// 'go test -count=N' to sample flakiness.
func TestAccFirewallAlias_parallelDestroy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckFirewallIPAliasDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallIPAliasConfig_multi(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIPAliasPresent("tf_acc_alias_a"),
					testAccCheckIPAliasPresent("tf_acc_alias_b"),
					testAccCheckIPAliasPresent("tf_acc_alias_c"),
					resource.TestCheckResourceAttr("pfsense_firewall_ip_alias.acc_a", "type", "network"),
				),
			},
			{
				// Single apply, two concurrent independent destroys.
				Config: testAccFirewallIPAliasConfig_single(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIPAliasAbsent("tf_acc_alias_b"),
					testAccCheckIPAliasAbsent("tf_acc_alias_c"),
					testAccCheckIPAliasPresent("tf_acc_alias_a"),
				),
			},
		},
	})
}

func testAccFirewallIPAliasConfig_multi() string {
	return accProviderConfig() + `
resource "pfsense_firewall_ip_alias" "acc_a" {
  name        = "tf_acc_alias_a"
  type        = "network"
  entries     = [{ ip = "10.240.0.0/28" }]
  description = "tf-acc parallel destroy a"
}

resource "pfsense_firewall_ip_alias" "acc_b" {
  name        = "tf_acc_alias_b"
  type        = "network"
  entries     = [{ ip = "10.241.0.0/28" }]
  description = "tf-acc parallel destroy b"
}

resource "pfsense_firewall_ip_alias" "acc_c" {
  name        = "tf_acc_alias_c"
  type        = "network"
  entries     = [{ ip = "10.242.0.0/28" }]
  description = "tf-acc parallel destroy c"
}
`
}

func testAccFirewallIPAliasConfig_single() string {
	return accProviderConfig() + `
resource "pfsense_firewall_ip_alias" "acc_a" {
  name        = "tf_acc_alias_a"
  type        = "network"
  entries     = [{ ip = "10.240.0.0/28" }]
  description = "tf-acc parallel destroy a"
}
`
}
