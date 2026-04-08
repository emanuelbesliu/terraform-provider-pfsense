data "pfsense_certificate_authority" "example" {
  descr = "Homelab CA Root CA"
}

output "ca_subject" {
  value = data.pfsense_certificate_authority.example.subject
}

output "ca_valid_to" {
  value = data.pfsense_certificate_authority.example.valid_to
}
