data "pfsense_interface_bridge" "lan_bridge" {
  bridge_if = "bridge0"
}

output "bridge_members" {
  value = data.pfsense_interface_bridge.lan_bridge.members
}
