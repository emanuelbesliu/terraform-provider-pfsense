data "pfsense_certificates" "all" {}

output "cert_names" {
  value = [for cert in data.pfsense_certificates.all.certificates : cert.descr]
}
