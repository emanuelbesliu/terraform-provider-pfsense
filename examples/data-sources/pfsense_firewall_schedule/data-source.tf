data "pfsense_firewall_schedule" "business_hours" {
  name = "BusinessHours"
}

output "business_hours_time_ranges" {
  value = data.pfsense_firewall_schedule.business_hours.time_range
}
