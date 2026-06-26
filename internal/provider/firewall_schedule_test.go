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

func testAccCheckFirewallScheduleDestroy(s *terraform.State) error {
	client, err := testAccNewClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pfsense_firewall_schedule" {
			continue
		}

		name := rs.Primary.Attributes["name"]

		_, err := client.GetSchedule(context.Background(), name)
		if err == nil {
			return fmt.Errorf("firewall schedule %q still exists", name)
		}

		if !errors.Is(err, pfsense.ErrNotFound) {
			return err
		}
	}

	return nil
}

func TestAccFirewallScheduleResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckFirewallScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallScheduleConfig_position("tfacctestworkhours", "Business hours", "1,2,3,4,5", "9:00", "17:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test", "name", "tfacctestworkhours"),
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test", "description", "Business hours"),
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test", "time_range.#", "1"),
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test", "time_range.0.position", "1,2,3,4,5"),
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test", "time_range.0.start_time", "9:00"),
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test", "time_range.0.stop_time", "17:00"),
					resource.TestCheckResourceAttrSet("pfsense_firewall_schedule.test", "label"),
				),
			},
			{
				ResourceName:      "pfsense_firewall_schedule.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFirewallScheduleConfig_position("tfacctestworkhours", "Updated hours", "1,2,3", "8:00", "18:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test", "description", "Updated hours"),
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test", "time_range.0.position", "1,2,3"),
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test", "time_range.0.start_time", "8:00"),
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test", "time_range.0.stop_time", "18:00"),
				),
			},
		},
	})
}

func TestAccFirewallScheduleResource_monthDay(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckFirewallScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallScheduleConfig_monthDay("tfacctestchristmas", "12", "25", "0:00", "23:59"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test_md", "name", "tfacctestchristmas"),
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test_md", "time_range.0.month", "12"),
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test_md", "time_range.0.day", "25"),
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test_md", "time_range.0.start_time", "0:00"),
					resource.TestCheckResourceAttr("pfsense_firewall_schedule.test_md", "time_range.0.stop_time", "23:59"),
				),
			},
			{
				ResourceName:      "pfsense_firewall_schedule.test_md",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFirewallScheduleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckFirewallScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallScheduleConfig_position("tfacctestds", "Single ds", "1,2,3,4,5", "9:00", "17:00") + `
data "pfsense_firewall_schedule" "test" {
  name       = "tfacctestds"
  depends_on = [pfsense_firewall_schedule.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.pfsense_firewall_schedule.test", "name", "tfacctestds"),
					resource.TestCheckResourceAttr("data.pfsense_firewall_schedule.test", "description", "Single ds"),
					resource.TestCheckResourceAttr("data.pfsense_firewall_schedule.test", "time_range.0.position", "1,2,3,4,5"),
				),
			},
		},
	})
}

func TestAccFirewallSchedulesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 func() { testAccPreCheck(t) },
		CheckDestroy:             testAccCheckFirewallScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallScheduleConfig_position("tfacctestschedds", "All ds", "1,2,3,4,5", "9:00", "17:00") + `
data "pfsense_firewall_schedules" "all" {
  depends_on = [pfsense_firewall_schedule.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.pfsense_firewall_schedules.all", "schedules.#"),
				),
			},
		},
	})
}

func testAccFirewallScheduleConfig_position(name, description, position, startTime, stopTime string) string {
	return fmt.Sprintf(`
resource "pfsense_firewall_schedule" "test" {
  name        = %q
  description = %q

  time_range = [
    {
      position   = %q
      start_time = %q
      stop_time  = %q
    },
  ]
}
`, name, description, position, startTime, stopTime)
}

func testAccFirewallScheduleConfig_monthDay(name, month, day, startTime, stopTime string) string {
	return fmt.Sprintf(`
resource "pfsense_firewall_schedule" "test_md" {
  name = %q

  time_range = [
    {
      month      = %q
      day        = %q
      start_time = %q
      stop_time  = %q
    },
  ]
}
`, name, month, day, startTime, stopTime)
}
