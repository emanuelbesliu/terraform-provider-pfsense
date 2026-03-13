data "pfsense_system_advanced_admin" "this" {}

output "system_advanced_admin" {
  value = data.pfsense_system_advanced_admin.this
}
