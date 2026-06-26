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

func testAccCheckInterfaceBridgeDestroy(s *terraform.State) error {
	client, err := testAccNewClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pfsense_interface_bridge" {
			continue
		}

		bridgeIf := rs.Primary.Attributes["bridge_if"]

		_, err := client.GetBridge(context.Background(), bridgeIf)
		if err == nil {
			return fmt.Errorf("bridge %q still exists", bridgeIf)
		}

		if !errors.Is(err, pfsense.ErrNotFound) {
			return err
		}
	}

	return nil
}

func TestAccInterfaceBridgeResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckInterfaceBridgeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceBridgeConfig_basic("tf-acc-test bridge"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("pfsense_interface_bridge.test", "bridge_if"),
					resource.TestCheckResourceAttr("pfsense_interface_bridge.test", "members.#", "1"),
					resource.TestCheckResourceAttr("pfsense_interface_bridge.test", "description", "tf-acc-test bridge"),
					resource.TestCheckResourceAttr("pfsense_interface_bridge.test", "enable_stp", "false"),
				),
			},
			{
				ResourceName:      "pfsense_interface_bridge.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInterfaceBridgeConfig_stp("tf-acc-test bridge stp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_interface_bridge.test", "description", "tf-acc-test bridge stp"),
					resource.TestCheckResourceAttr("pfsense_interface_bridge.test", "enable_stp", "true"),
					resource.TestCheckResourceAttr("pfsense_interface_bridge.test", "protocol", "rstp"),
					resource.TestCheckResourceAttr("pfsense_interface_bridge.test", "priority", "32768"),
				),
			},
		},
	})
}

func TestAccInterfaceBridgeDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckInterfaceBridgeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceBridgeConfig_basic("tf-acc-test bridge ds") + `
data "pfsense_interface_bridge" "test" {
  bridge_if  = pfsense_interface_bridge.test.bridge_if
  depends_on = [pfsense_interface_bridge.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.pfsense_interface_bridge.test", "description", "tf-acc-test bridge ds"),
					resource.TestCheckResourceAttr("data.pfsense_interface_bridge.test", "members.#", "1"),
				),
			},
		},
	})
}

func TestAccInterfaceBridgesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckInterfaceBridgeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceBridgeConfig_basic("tf-acc-test bridges ds") + `
data "pfsense_interface_bridges" "all" {
  depends_on = [pfsense_interface_bridge.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.pfsense_interface_bridges.all", "bridges.#"),
				),
			},
		},
	})
}

func testAccInterfaceBridgeConfig_basic(description string) string {
	return fmt.Sprintf(`
resource "pfsense_interface_bridge" "test" {
  members     = ["opt1"]
  description = %q
}
`, description)
}

func testAccInterfaceBridgeConfig_stp(description string) string {
	return fmt.Sprintf(`
resource "pfsense_interface_bridge" "test" {
  members     = ["opt1"]
  description = %q
  enable_stp  = true
  protocol    = "rstp"
  priority    = 32768
}
`, description)
}
