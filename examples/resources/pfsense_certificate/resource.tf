resource "pfsense_certificate" "webserver" {
  descr       = "Web Server Certificate"
  type        = "server"
  certificate = file("${path.module}/certs/webserver.pem")
  private_key = file("${path.module}/certs/webserver-key.pem")
}

resource "pfsense_certificate" "vpn_user" {
  descr       = "VPN User Certificate"
  type        = "user"
  certificate = file("${path.module}/certs/vpn-user.pem")
  private_key = file("${path.module}/certs/vpn-user-key.pem")
  caref       = pfsense_certificate_authority.internal_ca.refid
}
