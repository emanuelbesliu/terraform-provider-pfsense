data "pfsense_firewall_nat_npt_rules" "all" {}

output "rule_count" {
  value = length(data.pfsense_firewall_nat_npt_rules.all.rules)
}
