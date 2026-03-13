# Manage a static route to a remote network
resource "pfsense_static_route" "example" {
  network     = "10.10.0.0/24"
  gateway     = "WAN_GW"
  description = "Route to remote office network"
}

# Manage a static route in a disabled state
resource "pfsense_static_route" "disabled" {
  network     = "172.16.0.0/16"
  gateway     = "WAN_GW"
  description = "Maintenance network (disabled)"
  disabled    = true
}

# Reference a managed gateway resource
resource "pfsense_gateway" "backup" {
  name       = "BACKUP_GW"
  interface  = "opt1"
  ipprotocol = "inet"
  gateway    = "198.51.100.1"
}

resource "pfsense_static_route" "via_backup" {
  network     = "192.0.2.0/24"
  gateway     = pfsense_gateway.backup.name
  description = "Route via backup gateway"
}
