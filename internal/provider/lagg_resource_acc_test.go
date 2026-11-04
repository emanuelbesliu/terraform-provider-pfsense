package provider_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLaggResource_basic(t *testing.T) {
	t.Parallel()
	resourceName := "pfsense_lagg.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLaggDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLaggConfig("tf-acc round1", "lacp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile(`^lagg\d+$`)),
					resource.TestCheckResourceAttr(resourceName, "protocol", "lacp"),
					resource.TestCheckResourceAttr(resourceName, "members.0", "em0"),
					resource.TestCheckResourceAttr(resourceName, "members.1", "em2"),
					resource.TestCheckResourceAttr(resourceName, "description", "tf-acc round1"),
					resource.TestCheckResourceAttr(resourceName, "apply", "false"),
				),
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIgnore:              []string{"apply"},
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return state.RootModule().Resources[resourceName].Primary.Attributes["name"], nil
				},
			},
			{
				// in-place update of mutable fields (name must not change)
				Config: testAccLaggConfig("tf-acc round2", "failover"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "name", func(v string) error {
						return nil
					}),
					resource.TestCheckResourceAttr(resourceName, "protocol", "failover"),
					resource.TestCheckResourceAttr(resourceName, "description", "tf-acc round2"),
				),
			},
		},
	})
}

func testAccLaggConfig(description, protocol string) string {
	return fmt.Sprintf(`
provider "pfsense" {
  url             = "%s"
  username        = "%s"
  password        = "%s"
  tls_skip_verify = true
}

resource "pfsense_lagg" "test" {
  members     = ["em0", "em2"]
  protocol    = "%s"
  description = "%s"
  apply       = false
}
`, os.Getenv("TF_PFSENSE_URL"), os.Getenv("TF_PFSENSE_USERNAME"), os.Getenv("TF_PFSENSE_PASSWORD"), protocol, description)
}

func testAccCheckLaggDestroy(s *terraform.State) error {
	client, err := testAccNewClient()
	if err != nil {
		return err
	}
	laggs, err := client.GetLaggs(context.Background())
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pfsense_lagg" {
			continue
		}
		name := rs.Primary.Attributes["name"]
		for _, l := range laggs {
			if l.Name == name {
				return fmt.Errorf("lagg %s still exists after destroy", name)
			}
		}
	}
	return nil
}
