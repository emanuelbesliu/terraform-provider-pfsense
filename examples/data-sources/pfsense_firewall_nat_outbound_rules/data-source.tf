data "pfsense_firewall_nat_outbound_rules" "all" {}

output "nat_outbound_mode" {
  value = data.pfsense_firewall_nat_outbound_rules.all.mode
}

output "nat_outbound_rule_count" {
  value = length(data.pfsense_firewall_nat_outbound_rules.all.rules)
}
