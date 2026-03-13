data "pfsense_firewall_url_alias" "example" {
  name = "blocklist_ips"
}

output "url_alias_type" {
  value = data.pfsense_firewall_url_alias.example.type
}

output "url_alias_entries" {
  value = data.pfsense_firewall_url_alias.example.entries
}
