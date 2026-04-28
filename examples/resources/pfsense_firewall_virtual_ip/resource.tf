# IP Alias Virtual IP
resource "pfsense_firewall_virtual_ip" "web_alias" {
  mode        = "ipalias"
  interface   = "wan"
  subnet      = "203.0.113.10"
  subnet_bits = 32
  description = "Web server alias"
}

# CARP Virtual IP for high availability
resource "pfsense_firewall_virtual_ip" "ha_gateway" {
  mode        = "carp"
  interface   = "lan"
  subnet      = "10.0.1.1"
  subnet_bits = 24
  vhid        = 1
  advskew     = 0
  advbase     = 1
  password    = "secret"
  description = "LAN gateway CARP VIP"
}

# Proxy ARP Virtual IP
resource "pfsense_firewall_virtual_ip" "nat_pool" {
  mode        = "proxyarp"
  interface   = "wan"
  subnet      = "203.0.113.20"
  subnet_bits = 32
  description = "NAT pool address"
}
