package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccFirewallAlias_parallelDestroySix is the stress variant of
// TestAccFirewallAlias_parallelDestroy: five independent aliases are
// destroyed concurrently in one apply (six created, one kept). The
// retry-based delete fix was designed for this window; this test samples
// it with -count>1 to expose flakiness.
func TestAccFirewallAlias_parallelDestroySix(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckFirewallIPAliasDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallIPAliasConfig_six(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIPAliasPresent("tf_acc_alias_a"),
					testAccCheckIPAliasPresent("tf_acc_alias_b"),
					testAccCheckIPAliasPresent("tf_acc_alias_c"),
					testAccCheckIPAliasPresent("tf_acc_alias_d"),
					testAccCheckIPAliasPresent("tf_acc_alias_e"),
					testAccCheckIPAliasPresent("tf_acc_alias_f"),
				),
			},
			{
				// Single apply, five concurrent independent destroys.
				Config: testAccFirewallIPAliasConfig_single(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIPAliasAbsent("tf_acc_alias_b"),
					testAccCheckIPAliasAbsent("tf_acc_alias_c"),
					testAccCheckIPAliasAbsent("tf_acc_alias_d"),
					testAccCheckIPAliasAbsent("tf_acc_alias_e"),
					testAccCheckIPAliasAbsent("tf_acc_alias_f"),
					testAccCheckIPAliasPresent("tf_acc_alias_a"),
				),
			},
		},
	})
}

func testAccFirewallIPAliasConfig_six() string {
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

resource "pfsense_firewall_ip_alias" "acc_d" {
  name        = "tf_acc_alias_d"
  type        = "network"
  entries     = [{ ip = "10.243.0.0/28" }]
  description = "tf-acc parallel destroy d"
}

resource "pfsense_firewall_ip_alias" "acc_e" {
  name        = "tf_acc_alias_e"
  type        = "network"
  entries     = [{ ip = "10.244.0.0/28" }]
  description = "tf-acc parallel destroy e"
}

resource "pfsense_firewall_ip_alias" "acc_f" {
  name        = "tf_acc_alias_f"
  type        = "network"
  entries     = [{ ip = "10.245.0.0/28" }]
  description = "tf-acc parallel destroy f"
}
`
}
