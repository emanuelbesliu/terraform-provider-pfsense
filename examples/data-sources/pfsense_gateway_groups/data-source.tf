data "pfsense_gateway_groups" "all" {}

output "gateway_group_names" {
  value = [for g in data.pfsense_gateway_groups.all.gateway_groups : g.name]
}
