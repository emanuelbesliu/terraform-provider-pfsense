resource "pfsense_system_advanced_firewall" "example" {
  firewall_optimization  = "normal"
  bogons_update_interval = "monthly"
  nat_reflection_mode    = "proxy"
}
