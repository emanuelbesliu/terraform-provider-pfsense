data "pfsense_system_general" "this" {}

output "system_general" {
  value = data.pfsense_system_general.this
}
