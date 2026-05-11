data "pfsense_firewall_nat_one_to_one_rules" "all" {}

output "rule_count" {
  value = length(data.pfsense_firewall_nat_one_to_one_rules.all.rules)
}
