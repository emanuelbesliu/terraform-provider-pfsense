data "pfsense_firewall_nat_outbound" "example" {
  description = "LAN to WAN outbound NAT"
}

output "nat_outbound_interface" {
  value = data.pfsense_firewall_nat_outbound.example.interface
}
