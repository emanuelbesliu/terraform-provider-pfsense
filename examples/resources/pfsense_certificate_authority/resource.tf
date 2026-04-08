# Import an external root CA certificate into pfSense
resource "pfsense_certificate_authority" "example" {
  descr       = "My Root CA"
  certificate = file("${path.module}/root-ca.crt")
  trust       = true
}

# Import a CA with private key for signing
resource "pfsense_certificate_authority" "signing_ca" {
  descr         = "Internal Signing CA"
  certificate   = file("${path.module}/signing-ca.crt")
  private_key   = file("${path.module}/signing-ca.key")
  trust         = true
  random_serial = true
}
