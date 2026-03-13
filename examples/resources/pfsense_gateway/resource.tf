# Manage a gateway with a static IP address
resource "pfsense_gateway" "example" {
  name       = "WAN_GW"
  interface  = "wan"
  ipprotocol = "inet"
  gateway    = "203.0.113.1"
}

# Manage a gateway for a DHCP interface
resource "pfsense_gateway" "dhcp" {
  name        = "WAN_DHCP"
  interface   = "wan"
  ipprotocol  = "inet"
  gateway     = "dynamic"
  description = "WAN DHCP gateway"
}

# Manage a gateway with custom monitoring thresholds
resource "pfsense_gateway" "advanced" {
  name        = "BACKUP_GW"
  interface   = "opt1"
  ipprotocol  = "inet"
  gateway     = "10.0.0.1"
  description = "Backup WAN with custom monitoring"

  # Use an external host for link monitoring instead of the gateway IP
  monitor = "198.51.100.1"

  # Tuning for a high-latency link (e.g., satellite)
  latency_low  = 300
  latency_high = 800
  loss_low     = 10
  loss_high    = 30

  # Probe settings
  interval       = 2000
  loss_interval  = 5000
  time_period    = 120000
  alert_interval = 5000
  data_payload   = 1

  weight            = 2
  non_local_gateway = true
}

# Manage an IPv6 gateway
resource "pfsense_gateway" "ipv6" {
  name       = "WANv6_GW"
  interface  = "wan"
  ipprotocol = "inet6"
  gateway    = "2001:db8::1"
}
