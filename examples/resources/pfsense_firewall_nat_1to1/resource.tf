# Example 1: Basic 1:1 NAT rule
resource "pfsense_firewall_nat_1to1" "basic" {
  external            = "203.0.113.10"
  interface           = "wan"
  ipprotocol          = "inet"
  source_address      = "any"
  destination_address = "192.168.1.0/24"
  description         = "Basic 1:1 NAT"
}

# Example 2: 1:1 NAT with specific source and destination
resource "pfsense_firewall_nat_1to1" "specific" {
  external            = "203.0.113.20"
  interface           = "wan"
  ipprotocol          = "inet"
  source_address      = "10.0.0.0/8"
  destination_address = "192.168.1.100"
  description         = "1:1 NAT for specific subnet"
  nat_reflection      = "enable"
}

# Example 3: IPv6 1:1 NAT rule
resource "pfsense_firewall_nat_1to1" "ipv6" {
  external            = "2001:db8::10"
  interface           = "wan"
  ipprotocol          = "inet6"
  source_address      = "any"
  destination_address = "fd00::/64"
  description         = "IPv6 1:1 NAT"
}

# Example 4: Disabled 1:1 NAT rule
resource "pfsense_firewall_nat_1to1" "disabled" {
  external            = "203.0.113.30"
  interface           = "wan"
  ipprotocol          = "inet"
  source_address      = "any"
  destination_address = "192.168.1.50"
  description         = "Disabled 1:1 NAT"
  disabled            = true
}

# Example 5: 1:1 NAT with negation (no_binat)
resource "pfsense_firewall_nat_1to1" "negation" {
  external            = "203.0.113.40"
  interface           = "wan"
  ipprotocol          = "inet"
  source_address      = "any"
  destination_address = "192.168.1.0/24"
  description         = "Negated 1:1 NAT"
  no_binat            = true
}
