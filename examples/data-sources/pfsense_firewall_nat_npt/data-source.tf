data "pfsense_firewall_nat_npt" "ula_to_gua" {
  description = "ULA to GUA prefix translation"
}

output "npt_destination_prefix" {
  value = data.pfsense_firewall_nat_npt.ula_to_gua.destination_prefix
}
