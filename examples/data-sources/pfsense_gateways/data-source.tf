data "pfsense_gateways" "all" {}

output "gateway_names" {
  value = [for gw in data.pfsense_gateways.all.gateways : gw.name]
}
