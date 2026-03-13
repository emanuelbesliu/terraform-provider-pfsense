# Manage system general settings with custom hostname, domain, and DNS servers
resource "pfsense_system_general" "this" {
  hostname = "firewall"
  domain   = "example.com"

  dns_servers = [
    { address = "203.0.113.10" },
    { address = "198.51.100.10" },
  ]

  dns_override = false

  timezone    = "America/New_York"
  timeservers = "0.pfsense.pool.ntp.org"
  language    = "en_US"

  webgui_theme      = "pfSense-dark-BETA.css"
  dashboard_columns = 2
}

# Manage system general settings with DNS over TLS and per-server gateways
resource "pfsense_system_general" "dot" {
  hostname = "firewall"
  domain   = "example.com"

  dns_servers = [
    {
      address  = "203.0.113.10"
      hostname = "dns.example.com"
      gateway  = "WAN_DHCP"
    },
    {
      address  = "198.51.100.10"
      hostname = "dns2.example.com"
    },
  ]

  timezone = "Etc/UTC"
}
