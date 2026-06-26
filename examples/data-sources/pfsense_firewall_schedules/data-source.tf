data "pfsense_firewall_schedules" "all" {}

output "schedule_count" {
  value = length(data.pfsense_firewall_schedules.all.schedules)
}
