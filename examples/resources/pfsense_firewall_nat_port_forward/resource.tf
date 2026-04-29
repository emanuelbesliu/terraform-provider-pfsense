# Forward HTTP from WAN to internal web server
resource "pfsense_firewall_nat_port_forward" "web_server" {
  interface           = "wan"
  ipprotocol          = "inet"
  protocol            = "tcp"
  source_address      = "any"
  destination_address = "wanip"
  destination_port    = "8080"
  target              = "10.0.1.50"
  local_port          = "80"
  description         = "Forward HTTP to web server"
  associated_rule_id  = "pass"
}

# Forward HTTPS with NAT reflection enabled
resource "pfsense_firewall_nat_port_forward" "https_server" {
  interface           = "wan"
  ipprotocol          = "inet"
  protocol            = "tcp"
  destination_address = "wanip"
  destination_port    = "443"
  target              = "10.0.1.50"
  local_port          = "443"
  description         = "Forward HTTPS to web server"
  nat_reflection      = "purenat"
  associated_rule_id  = "pass"
}
