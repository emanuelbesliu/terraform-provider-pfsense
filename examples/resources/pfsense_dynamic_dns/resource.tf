resource "pfsense_dynamic_dns" "cloudflare" {
  type              = "cloudflare"
  interface         = "wan"
  host              = "myhost"
  domain_name       = "example.com"
  username          = "user@example.com"
  password          = "api-token-here"
  ttl               = "60"
  request_interface = "wan"
  description       = "Cloudflare DynDNS for myhost.example.com"
}
