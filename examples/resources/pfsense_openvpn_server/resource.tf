# Remote access OpenVPN server using TLS + user authentication
resource "pfsense_openvpn_server" "remote_access" {
  mode        = "server_tls_user"
  dev_mode    = "tun"
  protocol    = "UDP4"
  interface   = "wan"
  local_port  = "1194"
  description = "Remote access VPN"

  auth_mode = ["Local Database"]

  ca_ref    = "5f3e...caref"
  cert_ref  = "5f3e...certref"
  dh_length = "2048"

  data_ciphers          = ["AES-256-GCM", "AES-128-GCM", "CHACHA20-POLY1305"]
  data_ciphers_fallback = "AES-256-CBC"
  digest                = "SHA256"

  tunnel_network = "10.8.0.0/24"
  topology       = "subnet"

  dns_server_enable = true
  dns_server1       = "10.0.0.1"

  gw_redir        = true
  max_clients     = "50"
  verbosity_level = "3"
}

# Site-to-site OpenVPN server using a shared key
resource "pfsense_openvpn_server" "site_to_site" {
  mode        = "p2p_shared_key"
  dev_mode    = "tun"
  protocol    = "UDP4"
  interface   = "wan"
  local_port  = "1195"
  description = "Site-to-site tunnel"

  shared_key     = file("${path.module}/openvpn_shared.key")
  tunnel_network = "10.9.0.0/30"
  remote_network = "192.168.20.0/24"

  data_ciphers_fallback = "AES-256-CBC"
  digest                = "SHA256"
}
