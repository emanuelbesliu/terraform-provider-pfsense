data "pfsense_openvpn_server" "remote_access" {
  vpn_id = "1"
}

output "openvpn_server_tunnel_network" {
  value = data.pfsense_openvpn_server.remote_access.tunnel_network
}
