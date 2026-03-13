resource "pfsense_interface" "example" {
  port         = "igb0.100"
  description  = "SERVERS"
  enabled      = true
  ipv4_type    = "staticv4"
  ipv4_address = "203.0.113.1"
  ipv4_subnet  = "24"
  ipv6_type    = "none"
}
