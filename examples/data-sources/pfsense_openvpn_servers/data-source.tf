data "pfsense_openvpn_servers" "all" {}

output "openvpn_server_count" {
  value = length(data.pfsense_openvpn_servers.all.all)
}
