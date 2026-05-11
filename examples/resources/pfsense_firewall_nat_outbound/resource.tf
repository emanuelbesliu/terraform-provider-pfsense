# Outbound NAT rule using interface address
resource "pfsense_firewall_nat_outbound" "lan_to_wan" {
  interface           = "wan"
  protocol            = "tcp"
  source_address      = "192.168.1.0/24"
  destination_address = "any"
  target              = "(self)"
  description         = "LAN to WAN outbound NAT"
}

# Outbound NAT rule with specific translation IP
resource "pfsense_firewall_nat_outbound" "dmz_to_wan" {
  interface           = "wan"
  source_address      = "10.0.2.0/24"
  destination_address = "any"
  target              = "other-subnet"
  target_ip           = "203.0.113.10"
  target_ip_subnet    = "32"
  pool_options        = "round-robin"
  description         = "DMZ outbound NAT with specific IP"
}

# NAT exclusion rule (no NAT)
resource "pfsense_firewall_nat_outbound" "vpn_no_nat" {
  interface           = "wan"
  source_address      = "192.168.1.0/24"
  destination_address = "10.10.0.0/16"
  target              = ""
  no_nat              = true
  description         = "No NAT for VPN traffic"
}
