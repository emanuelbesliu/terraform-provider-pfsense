resource "pfsense_vlan" "example" {
  parent_interface = "igb0"
  tag              = 100
  pcp              = 0
  description      = "Example VLAN"
}
