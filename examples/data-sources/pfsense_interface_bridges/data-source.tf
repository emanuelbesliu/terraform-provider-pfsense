data "pfsense_interface_bridges" "all" {}

output "bridge_count" {
  value = length(data.pfsense_interface_bridges.all.bridges)
}
