resource "pfsense_firewall_nat_one_to_one" "webserver" {
  interface           = "wan"
  external            = "203.0.113.10"
  ipprotocol          = "inet"
  source_address      = "any"
  destination_address = "192.168.1.10"
  description         = "Web server 1:1 NAT"
}

resource "pfsense_firewall_nat_one_to_one" "mailserver" {
  interface           = "wan"
  external            = "203.0.113.11"
  ipprotocol          = "inet"
  source_address      = "any"
  destination_address = "192.168.1.11"
  disabled            = true
  no_binat            = true
  nat_reflection      = "enable"
  description         = "Mail server 1:1 NAT (disabled)"
}
