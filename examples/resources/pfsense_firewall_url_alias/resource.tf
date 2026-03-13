# URL Table example (single URL, periodically downloaded)
resource "pfsense_firewall_url_alias" "urltable_example" {
  name             = "blocklist_ips"
  description      = "IP blocklist updated daily"
  type             = "urltable"
  update_frequency = 1

  entries = [
    { url = "https://example.com/blocklist.txt", description = "Example blocklist" },
  ]
}

# URL (IPs) example (multiple URLs, contents merged)
resource "pfsense_firewall_url_alias" "url_example" {
  name        = "trusted_ips"
  description = "Trusted IP addresses from multiple sources"
  type        = "url"

  entries = [
    { url = "https://example.com/trusted-ips-1.txt", description = "Source 1" },
    { url = "https://example.com/trusted-ips-2.txt", description = "Source 2" },
  ]
}

# URL Table (Ports) example
resource "pfsense_firewall_url_alias" "urltable_ports_example" {
  name             = "blocked_ports"
  description      = "Blocked ports from URL"
  type             = "urltable_ports"
  update_frequency = 7

  entries = [
    { url = "https://example.com/blocked-ports.txt" },
  ]
}
