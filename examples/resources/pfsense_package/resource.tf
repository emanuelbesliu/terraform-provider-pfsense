resource "pfsense_package" "example" {
  name = "pfSense-pkg-nmap"
}

# Install a third-party package from a URL
resource "pfsense_package" "custom" {
  name        = "pfSense-pkg-saml2-auth"
  package_url = "https://example.com/packages/pfSense-pkg-saml2-auth-0.4.txz"
}
