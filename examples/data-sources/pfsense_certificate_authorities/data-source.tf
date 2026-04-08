data "pfsense_certificate_authorities" "all" {}

output "ca_names" {
  value = [for ca in data.pfsense_certificate_authorities.all.certificate_authorities : ca.descr]
}
