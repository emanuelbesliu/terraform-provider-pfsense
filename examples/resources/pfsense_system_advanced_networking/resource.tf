resource "pfsense_system_advanced_networking" "example" {
  dhcp_backend = "kea"
  ipv6_allow   = true
  prefer_ipv4  = true
}
