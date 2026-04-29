data "pfsense_certificate" "example" {
  descr = "Web Server Certificate"
}

output "cert_subject" {
  value = data.pfsense_certificate.example.subject
}

output "cert_valid_to" {
  value = data.pfsense_certificate.example.valid_to
}
