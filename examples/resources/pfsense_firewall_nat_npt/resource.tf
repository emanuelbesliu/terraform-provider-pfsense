resource "pfsense_firewall_nat_npt" "ula_to_gua" {
  interface          = "wan"
  source_prefix      = "fd00:1::/64"
  destination_prefix = "2001:db8:0:1::/64"
  description        = "ULA to GUA prefix translation"
}

resource "pfsense_firewall_nat_npt" "guest" {
  interface          = "wan"
  source_prefix      = "fd00:2::/64"
  destination_prefix = "2001:db8:0:2::/64"
  disabled           = true
  description        = "Guest network NPt (disabled)"
}
